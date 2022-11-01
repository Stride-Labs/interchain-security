package keeper

import (
	"fmt"
	"time"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
	providertypes "github.com/cosmos/interchain-security/x/ccv/provider/types"
	ccvtypes "github.com/cosmos/interchain-security/x/ccv/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

// This file contains functionality relevant to the throttling of slash and vsc matured packets, aka circuit breaker logic.

// HandlePendingSlashPackets handles all or some portion of pending slash packets depending on circuit breaker logic.
// This method executes every end block routine
func (k Keeper) HandlePendingSlashPackets(ctx sdktypes.Context) {

	meter := k.GetSlashGasMeter(ctx)

	handledEntries := []providertypes.SlashPacketEntry{}
	k.IteratePendingSlashPacketEntries(ctx, func(entry providertypes.SlashPacketEntry) bool {

		// TODO Get data from other queue, handle the vsc matured that're there too

		_, err := k.HandleSlashPacket(ctx, entry.ConsumerChainID,
			ccvtypes.SlashPacketData{}, // TODO
		)
		if err != nil {
			panic(fmt.Sprintf("failed to handle slash packet: %s", err.Error()))
		}

		valPower := k.stakingKeeper.GetLastValidatorPower(ctx, entry.ValAddr)
		meter.Sub(sdktypes.NewInt(valPower))
		handledEntries = append(handledEntries, entry)

		// Do not handle anymore slash packets if the meter has 0 or negative gas
		return !meter.IsPositive()
	})

	k.DeletePendingSlashPacketEntries(ctx, handledEntries...)
	k.SetSlashGasMeter(ctx, meter)
}

// TODO: Make an e2e test that asserts that the order of endblockers is correct between staking and ccv
// TODO: ie. the staking updates to voting power need to occur before circuit breaker logic, so circuit breaker has most up to date val powers.

// CheckForSlashMeterReplenishment checks if the slash gas meter should be replenished, and if so, replenishes it.
// This method executes every end block routine.
// TODO: hook this into endblocker, unit and e2e tests, tests must include odd time formats, since UTC is is used
func (k Keeper) CheckForSlashMeterReplenishment(ctx sdktypes.Context) {
	// TODO: Need to set initial replenishment time
	if ctx.BlockTime().UTC().After(k.GetLastSlashGasReplenishTime(ctx).Add(time.Hour)) {
		// TODO: Use param for replenish period, allowance, etc.
		// TODO: change code and documentation to reflect that this is a string fraction param
		slashGasAllowanceFraction := sdktypes.NewDec(5).Quo(sdktypes.NewDec(100)) // This will be a string param, ex: "0.05"

		// Compute slash gas allowance in units of tendermint voting power (integer)
		// TODO: total voting power would change as validators are jailed, is there a timing guarantee we can
		// make on maximum slash packet delay? Perhaps we use a static "total voting power" like the tm maximum.

		// TODO: Maybe the param could itself be an amount of voting power? This would be easier to reason about
		totalPower := k.stakingKeeper.GetLastTotalPower(ctx)
		slashGasAllowance := sdktypes.NewInt(slashGasAllowanceFraction.MulInt(totalPower).RoundInt64())

		meter := k.GetSlashGasMeter(ctx)

		// Replenish gas up to gas allowance per period. That is, if meter was negative
		// before being replenished, it'll gain some additional gas. However, if the meter
		// was 0 or positive in value, it'll be replenished only up to it's allowance for the period.
		meter = meter.Add(slashGasAllowance)
		if meter.GT(slashGasAllowance) {
			meter = slashGasAllowance
		}
		k.SetSlashGasMeter(ctx, meter)
		k.SetLastSlashGasReplenishTime(ctx, ctx.BlockTime())
	}
}

//
// CRUD section
//

// TODO: Maybe this method just goes in the on recv method
func (k Keeper) HandleOrQueueVSCMaturedPacket(ctx sdktypes.Context, consumerChainID string, data ccvtypes.VSCMaturedPacketData) {
	// TODO: if queue for this chain is empty (no pending slash packets), handle vsc matured packet immediately
	// else queue it
	k.QueuePendingPacketData(ctx, consumerChainID, data)
}

// TODO: maybe this goes in the on recv method too
func (k Keeper) QueuePendingSlashPacket(
	ctx sdktypes.Context, recvTime time.Time, consumerChainID string, data ccvtypes.SlashPacketData) {
	k.QueuePendingSlashPacketEntry(ctx, providertypes.NewSlashPacketEntry(recvTime, consumerChainID, data.Validator.Address))
	k.QueuePendingPacketData(ctx, consumerChainID, data)
}

// Highest level "parent" queue
// Note: this will overwrite the existing entry if a malicious consumer sends duplicate slash packets in the same block.
// TODO: unit test edge case where duplicate slash packet entries are added
func (k Keeper) QueuePendingSlashPacketEntry(ctx sdktypes.Context, entry providertypes.SlashPacketEntry) {
	store := ctx.KVStore(k.storeKey)
	key := providertypes.PendingSlashPacketEntryKey(entry)
	// Note: Val address is stored as value to assist in debugging.
	store.Set(key, entry.ValAddr)
}

// GetAllPendingSlashPacketEntries returns all pending slash packet entries in the queue
// This method is used for testing purposes only
func (k Keeper) GetAllPendingSlashPacketEntries(ctx sdktypes.Context) (entries []providertypes.SlashPacketEntry) {
	k.IteratePendingSlashPacketEntries(ctx, func(entry providertypes.SlashPacketEntry) bool {
		entries = append(entries, entry)
		return false
	})
	return entries
}

// IteratePendingSlashPackets iterates over the pending slash packet entry queue and calls the provided callback
func (k Keeper) IteratePendingSlashPacketEntries(ctx sdktypes.Context, cb func(providertypes.SlashPacketEntry) bool) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdktypes.KVStorePrefixIterator(store, []byte{providertypes.PendingSlashPacketEntryBytePrefix})
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		recvTime, chainID := providertypes.ParsePendingSlashPacketEntryKey(iterator.Key())
		valAddr := iterator.Value()
		entry := providertypes.NewSlashPacketEntry(recvTime, chainID, valAddr)
		if cb(entry) {
			break
		}
	}
}

// DeletePendingSlashPackets deletes the given entries from the pending slash packet entry queue
func (k Keeper) DeletePendingSlashPacketEntries(ctx sdktypes.Context, entries ...providertypes.SlashPacketEntry) {
	store := ctx.KVStore(k.storeKey)
	for _, entry := range entries {
		store.Delete(providertypes.PendingSlashPacketEntryKey(entry))
	}
}

// Needs to be true FIFO
// Handling a slash packet at the head of the queue handles all vsc matured packets after the head
// VSC matured packets at the head are handled immediately
func (k Keeper) QueuePendingPacketData(ctx sdktypes.Context, consumerChainID string, data interface{}) {
	switch data.(type) {
	case ccvtypes.SlashPacketData:
		// TODO: add slash packet data to queue
	case ccvtypes.VSCMaturedPacketData:
		// TODO: add vsc packet data to queue
	default:
		panic("invalid packet data type")
	}
}

// GetSlashGasMeter returns a meter (persisted as a signed int) which stores "slash gas",
// ie. an amount of voting power corresponding to an allowance of validators (with non-zero voting power)
// that can be jailed at a given time.
//
// Note: the value of this decimal should always be in the range of tendermint's [-MaxVotingPower, MaxVotingPower]
// TODO: If you keep slash gas meter as a percent, make sure it's clear that the param is a percent (put in name)
func (k Keeper) GetSlashGasMeter(ctx sdktypes.Context) sdktypes.Int {
	// TODO: is this the standard way to set a signed int?
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(providertypes.SlashGasMeterKey())
	if bz == nil {
		panic("slash gas meter not set")
	}
	value := sdktypes.ZeroInt()
	err := value.Unmarshal(bz)
	if err != nil {
		panic(fmt.Sprintf("failed to unmarshal slash gas meter: %v", err))
	}
	return value
}

// SetSlashGasMeter sets the "slash gas" meter to the given signed int value
//
// Note: the value of this decimal should always be in the range of tendermint's [-MaxVotingPower, MaxVotingPower]
func (k Keeper) SetSlashGasMeter(ctx sdktypes.Context, value sdktypes.Int) {
	if value.GT(sdktypes.NewInt(tmtypes.MaxTotalVotingPower)) {
		panic("slash gas meter value cannot be greater than tendermint's MaxTotalVotingPower")
	}
	if value.LT(sdktypes.NewInt(-tmtypes.MaxTotalVotingPower)) {
		panic("slash gas meter value cannot be less than negative tendermint's MaxTotalVotingPower")
	}
	store := ctx.KVStore(k.storeKey)
	bz, err := value.Marshal()
	if err != nil {
		panic(fmt.Sprintf("failed to marshal slash gas meter: %v", err))
	}
	store.Set(providertypes.SlashGasMeterKey(), bz)
}

// GetLastSlashGasReplenishTime returns the last UTC time the slash gas meter was replenished
func (k Keeper) GetLastSlashGasReplenishTime(ctx sdktypes.Context) time.Time {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(providertypes.LastSlashGasReplenishTimeKey())
	if bz == nil {
		panic("last slash gas replenish time not set")
	}
	time, err := sdktypes.ParseTimeBytes(bz)
	if err != nil {
		panic(fmt.Sprintf("failed to parse last slash gas replenish time: %s", err))
	}
	return time.UTC()
}

// SetLastSlashGasReplenishTime sets the last time the slash gas meter was replenished
func (k Keeper) SetLastSlashGasReplenishTime(ctx sdktypes.Context, time time.Time) {
	store := ctx.KVStore(k.storeKey)
	store.Set(providertypes.LastSlashGasReplenishTimeKey(), sdktypes.FormatTimeBytes(time.UTC()))
}
