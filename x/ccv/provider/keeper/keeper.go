package keeper

import (
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	clienttypes "github.com/cosmos/ibc-go/v3/modules/core/02-client/types"
	conntypes "github.com/cosmos/ibc-go/v3/modules/core/03-connection/types"
	channeltypes "github.com/cosmos/ibc-go/v3/modules/core/04-channel/types"
	host "github.com/cosmos/ibc-go/v3/modules/core/24-host"
	ibcexported "github.com/cosmos/ibc-go/v3/modules/core/exported"
	ibctmtypes "github.com/cosmos/ibc-go/v3/modules/light-clients/07-tendermint/types"

	consumertypes "github.com/cosmos/interchain-security/x/ccv/consumer/types"
	"github.com/cosmos/interchain-security/x/ccv/provider/types"
	ccv "github.com/cosmos/interchain-security/x/ccv/types"

	"github.com/tendermint/tendermint/libs/log"
)

// Keeper defines the Cross-Chain Validation Provider Keeper
type Keeper struct {
	storeKey         sdk.StoreKey
	cdc              codec.BinaryCodec
	paramSpace       paramtypes.Subspace
	scopedKeeper     ccv.ScopedKeeper
	channelKeeper    ccv.ChannelKeeper
	portKeeper       ccv.PortKeeper
	connectionKeeper ccv.ConnectionKeeper
	accountKeeper    ccv.AccountKeeper
	clientKeeper     ccv.ClientKeeper
	stakingKeeper    ccv.StakingKeeper
	slashingKeeper   ccv.SlashingKeeper
	feeCollectorName string
}

// NewKeeper creates a new provider Keeper instance
func NewKeeper(
	cdc codec.BinaryCodec, key sdk.StoreKey, paramSpace paramtypes.Subspace, scopedKeeper ccv.ScopedKeeper,
	channelKeeper ccv.ChannelKeeper, portKeeper ccv.PortKeeper,
	connectionKeeper ccv.ConnectionKeeper, clientKeeper ccv.ClientKeeper,
	stakingKeeper ccv.StakingKeeper, slashingKeeper ccv.SlashingKeeper,
	accountKeeper ccv.AccountKeeper, feeCollectorName string,
) Keeper {
	// set KeyTable if it has not already been set
	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		cdc:              cdc,
		storeKey:         key,
		paramSpace:       paramSpace,
		scopedKeeper:     scopedKeeper,
		channelKeeper:    channelKeeper,
		portKeeper:       portKeeper,
		connectionKeeper: connectionKeeper,
		accountKeeper:    accountKeeper,
		clientKeeper:     clientKeeper,
		stakingKeeper:    stakingKeeper,
		slashingKeeper:   slashingKeeper,
		feeCollectorName: feeCollectorName,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+host.ModuleName+"-"+types.ModuleName)
}

// IsBound checks if the CCV module is already bound to the desired port
func (k Keeper) IsBound(ctx sdk.Context, portID string) bool {
	_, ok := k.scopedKeeper.GetCapability(ctx, host.PortPath(portID))
	return ok
}

// BindPort defines a wrapper function for the port Keeper's function in
// order to expose it to module's InitGenesis function
func (k Keeper) BindPort(ctx sdk.Context, portID string) error {
	cap := k.portKeeper.BindPort(ctx, portID)
	return k.ClaimCapability(ctx, cap, host.PortPath(portID))
}

// GetPort returns the portID for the CCV module. Used in ExportGenesis
func (k Keeper) GetPort(ctx sdk.Context) string {
	store := ctx.KVStore(k.storeKey)
	return string(store.Get(types.PortKey()))
}

// SetPort sets the portID for the CCV module. Used in InitGenesis
func (k Keeper) SetPort(ctx sdk.Context, portID string) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.PortKey(), []byte(portID))
}

// AuthenticateCapability wraps the scopedKeeper's AuthenticateCapability function
func (k Keeper) AuthenticateCapability(ctx sdk.Context, cap *capabilitytypes.Capability, name string) bool {
	return k.scopedKeeper.AuthenticateCapability(ctx, cap, name)
}

// ClaimCapability allows the transfer module that can claim a capability that IBC module
// passes to it
func (k Keeper) ClaimCapability(ctx sdk.Context, cap *capabilitytypes.Capability, name string) error {
	return k.scopedKeeper.ClaimCapability(ctx, cap, name)
}

// SetChainToChannel sets the mapping from a consumer chainID to the CCV channel ID for that consumer chain.
func (k Keeper) SetChainToChannel(ctx sdk.Context, chainID, channelID string) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.ChainToChannelKey(chainID), []byte(channelID))
}

// GetChainToChannel gets the CCV channelID for the given consumer chainID
func (k Keeper) GetChainToChannel(ctx sdk.Context, chainID string) (string, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.ChainToChannelKey(chainID))
	if bz == nil {
		return "", false
	}
	return string(bz), true
}

// DeleteChainToChannel deletes the CCV channel ID for the given consumer chain ID
func (k Keeper) DeleteChainToChannel(ctx sdk.Context, chainID string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.ChainToChannelKey(chainID))
}

// IterateConsumerChains iterates over the IDs of all clients to consumer chains.
// Consumer chains with created clients are also referred to as registered.
//
// Note that the registered consumer chains are stored under keys with the following format:
// ChainToClientBytePrefix | chainID
// Thus, the iteration is in ascending order of chainIDs.
func (k Keeper) IterateConsumerChains(ctx sdk.Context, cb func(ctx sdk.Context, chainID, clientID string) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte{types.ChainToClientBytePrefix})
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		// remove 1 byte prefix from key to retrieve chainID
		chainID := string(iterator.Key()[1:])
		clientID := string(iterator.Value())

		stop := cb(ctx, chainID, clientID)
		if stop {
			return
		}
	}
}

// SetChannelToChain sets the mapping from the CCV channel ID to the consumer chainID.
func (k Keeper) SetChannelToChain(ctx sdk.Context, channelID, chainID string) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.ChannelToChainKey(channelID), []byte(chainID))
}

// GetChannelToChain gets the consumer chainID for a given CCV channelID
func (k Keeper) GetChannelToChain(ctx sdk.Context, channelID string) (string, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.ChannelToChainKey(channelID))
	if bz == nil {
		return "", false
	}
	return string(bz), true
}

// DeleteChannelToChain deletes the consumer chain ID for a given CCV channe lID
func (k Keeper) DeleteChannelToChain(ctx sdk.Context, channelID string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.ChannelToChainKey(channelID))
}

// IterateChannelToChain iterates over the channel to chain mappings.
//
// Note that mapping from CCV channel IDs to consumer chainIDs is stored under keys with the following format:
// ChannelToChainBytePrefix | channelID
// Thus, the iteration is in ascending order of channelIDs.
func (k Keeper) IterateChannelToChain(ctx sdk.Context, cb func(ctx sdk.Context, channelID, chainID string) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte{types.ChannelToChainBytePrefix})
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		// remove prefix from key to retrieve channelID
		channelID := string(iterator.Key()[1:])

		chainID := string(iterator.Value())

		stop := cb(ctx, channelID, chainID)
		if stop {
			break
		}
	}
}

func (k Keeper) SetConsumerGenesis(ctx sdk.Context, chainID string, gen consumertypes.GenesisState) error {
	store := ctx.KVStore(k.storeKey)
	bz, err := gen.Marshal()
	if err != nil {
		return err
	}
	store.Set(types.ConsumerGenesisKey(chainID), bz)

	return nil
}

func (k Keeper) GetConsumerGenesis(ctx sdk.Context, chainID string) (consumertypes.GenesisState, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.ConsumerGenesisKey(chainID))
	if bz == nil {
		return consumertypes.GenesisState{}, false
	}

	var data consumertypes.GenesisState
	if err := data.Unmarshal(bz); err != nil {
		panic(fmt.Errorf("consumer genesis could not be unmarshaled: %w", err))
	}
	return data, true
}

func (k Keeper) DeleteConsumerGenesis(ctx sdk.Context, chainID string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.ConsumerGenesisKey(chainID))
}

// VerifyConsumerChain verifies that the chain trying to connect on the channel handshake
// is the expected consumer chain.
func (k Keeper) VerifyConsumerChain(ctx sdk.Context, channelID string, connectionHops []string) error {
	if len(connectionHops) != 1 {
		return sdkerrors.Wrap(channeltypes.ErrTooManyConnectionHops, "must have direct connection to provider chain")
	}
	connectionID := connectionHops[0]
	clientID, tmClient, err := k.getUnderlyingClient(ctx, connectionID)
	if err != nil {
		return err
	}
	ccvClientId, found := k.GetConsumerClientId(ctx, tmClient.ChainId)
	if !found {
		return sdkerrors.Wrapf(ccv.ErrClientNotFound, "cannot find client for consumer chain %s", tmClient.ChainId)
	}
	if ccvClientId != clientID {
		return sdkerrors.Wrapf(ccv.ErrInvalidConsumerClient, "CCV channel must be built on top of CCV client. expected %s, got %s", ccvClientId, clientID)
	}

	// Verify that there isn't already a CCV channel for the consumer chain
	if prevChannel, ok := k.GetChainToChannel(ctx, tmClient.ChainId); ok {
		return sdkerrors.Wrapf(ccv.ErrDuplicateChannel, "CCV channel with ID: %s already created for consumer chain %s", prevChannel, tmClient.ChainId)
	}
	return nil
}

// SetConsumerChain ensures that the consumer chain has not already been
// set by a different channel, and then sets the consumer chain mappings
// in keeper, and set the channel status to validating.
// If there is already a CCV channel between the provider and consumer
// chain then close the channel, so that another channel can be made.
//
// SetConsumerChain is called by OnChanOpenConfirm.
func (k Keeper) SetConsumerChain(ctx sdk.Context, channelID string) error {
	channel, ok := k.channelKeeper.GetChannel(ctx, ccv.ProviderPortID, channelID)
	if !ok {
		return sdkerrors.Wrapf(channeltypes.ErrChannelNotFound, "channel not found for channel ID: %s", channelID)
	}
	if len(channel.ConnectionHops) != 1 {
		return sdkerrors.Wrap(channeltypes.ErrTooManyConnectionHops, "must have direct connection to consumer chain")
	}
	connectionID := channel.ConnectionHops[0]
	clientID, tmClient, err := k.getUnderlyingClient(ctx, connectionID)
	if err != nil {
		return err
	}
	// Verify that there isn't already a CCV channel for the consumer chain
	chainID := tmClient.ChainId
	if prevChannelID, ok := k.GetChainToChannel(ctx, chainID); ok {
		return sdkerrors.Wrapf(ccv.ErrDuplicateChannel, "CCV channel with ID: %s already created for consumer chain %s", prevChannelID, chainID)
	}

	// the CCV channel is established:
	// - set channel mappings
	k.SetChainToChannel(ctx, chainID, channelID)
	k.SetChannelToChain(ctx, channelID, chainID)
	// - set current block height for the consumer chain initialization
	k.SetInitChainHeight(ctx, chainID, uint64(ctx.BlockHeight()))
	// - remove init timeout timestamp
	k.DeleteInitTimeoutTimestamp(ctx, chainID)

	// emit event on successful addition
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			ccv.EventTypeChannelEstablished,
			sdk.NewAttribute(sdk.AttributeKeyModule, consumertypes.ModuleName),
			sdk.NewAttribute(ccv.AttributeChainID, chainID),
			sdk.NewAttribute(conntypes.AttributeKeyClientID, clientID),
			sdk.NewAttribute(channeltypes.AttributeKeyChannelID, channelID),
			sdk.NewAttribute(conntypes.AttributeKeyConnectionID, connectionID),
		),
	)
	return nil
}

// SetUnbondingOp sets the UnbondingOp by its unique ID
func (k Keeper) SetUnbondingOp(ctx sdk.Context, unbondingOp ccv.UnbondingOp) {
	store := ctx.KVStore(k.storeKey)
	bz, err := unbondingOp.Marshal()
	if err != nil {
		panic(fmt.Errorf("unbonding op could not be marshaled: %w", err))
	}
	store.Set(types.UnbondingOpKey(unbondingOp.Id), bz)
}

// GetUnbondingOp gets a UnbondingOp by its unique ID
func (k Keeper) GetUnbondingOp(ctx sdk.Context, id uint64) (ccv.UnbondingOp, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.UnbondingOpKey(id))
	if bz == nil {
		return ccv.UnbondingOp{}, false
	}

	return types.MustUnmarshalUnbondingOp(k.cdc, bz), true
}

// DeleteUnbondingOp deletes a UnbondingOp given its ID
func (k Keeper) DeleteUnbondingOp(ctx sdk.Context, id uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.UnbondingOpKey(id))
}

// IterateUnbondingOps iterates over UnbondingOps.
//
// Note that UnbondingOps are stored under keys with the following format:
// UnbondingOpBytePrefix | ID
// Thus, the iteration is in ascending order of IDs.
func (k Keeper) IterateUnbondingOps(ctx sdk.Context, cb func(id uint64, ubdOp ccv.UnbondingOp) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte{types.UnbondingOpBytePrefix})

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		id := sdk.BigEndianToUint64(iterator.Key()[1:])
		bz := iterator.Value()
		if bz == nil {
			panic(fmt.Errorf("unbonding operation is nil for id %d", id))
		}
		ubdOp := types.MustUnmarshalUnbondingOp(k.cdc, bz)

		stop := cb(id, ubdOp)
		if stop {
			break
		}
	}
}

// SetUnbondingOpIndex sets an unbonding index,
// i.e., the IDs of unbonding operations that are waiting for
// a VSCMaturedPacket with vscID from a consumer with chainID
func (k Keeper) SetUnbondingOpIndex(ctx sdk.Context, chainID string, vscID uint64, IDs []uint64) {
	store := ctx.KVStore(k.storeKey)

	index := ccv.UnbondingOpsIndex{
		Ids: IDs,
	}
	bz, err := index.Marshal()
	if err != nil {
		panic("Failed to marshal UnbondingOpsIndex")
	}

	store.Set(types.UnbondingOpIndexKey(chainID, vscID), bz)
}

// IterateUnbondingOpIndex iterates over the unbonding indexes for a given chainID.
//
// Note that the unbonding indexes for a given chainID are stored under keys with the following format:
// UnbondingOpIndexBytePrefix | len(chainID) | chainID | vscID
// Thus, the iteration is in ascending order of vscIDs.
func (k Keeper) IterateUnbondingOpIndex(
	ctx sdk.Context,
	chainID string,
	cb func(vscID uint64, ubdIndex []uint64) (stop bool),
) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.ChainIdWithLenKey(types.UnbondingOpIndexBytePrefix, chainID))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		// parse key to get the current VSC ID
		_, vscID, err := types.ParseUnbondingOpIndexKey(iterator.Key())
		if err != nil {
			panic(fmt.Errorf("failed to parse UnbondingOpIndexKey: %w", err))
		}

		var index ccv.UnbondingOpsIndex
		if err = index.Unmarshal(iterator.Value()); err != nil {
			panic("Failed to unmarshal JSON")
		}

		stop := cb(vscID, index.GetIds())
		if stop {
			return
		}
	}
}

// GetUnbondingOpIndex gets an unbonding index,
// i.e., the IDs of unbonding operations that are waiting for
// a VSCMaturedPacket with vscID from a consumer with chainID
func (k Keeper) GetUnbondingOpIndex(ctx sdk.Context, chainID string, vscID uint64) ([]uint64, bool) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.UnbondingOpIndexKey(chainID, vscID))
	if bz == nil {
		return []uint64{}, false
	}

	var idx ccv.UnbondingOpsIndex
	if err := idx.Unmarshal(bz); err != nil {
		panic("Failed to unmarshal UnbondingOpsIndex")
	}

	return idx.GetIds(), true
}

// DeleteUnbondingOpIndex deletes an unbonding index,
// i.e., the IDs of unbonding operations that are waiting for
// a VSCMaturedPacket with vscID from a consumer with chainID
func (k Keeper) DeleteUnbondingOpIndex(ctx sdk.Context, chainID string, vscID uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.UnbondingOpIndexKey(chainID, vscID))
}

// GetUnbondingOpsFromIndex gets the unbonding ops waiting for a given chainID and vscID
func (k Keeper) GetUnbondingOpsFromIndex(
	ctx sdk.Context,
	chainID string,
	vscID uint64,
) (entries []ccv.UnbondingOp, found bool) {
	ids, found := k.GetUnbondingOpIndex(ctx, chainID, vscID)
	if !found {
		return entries, false
	}
	for _, id := range ids {
		entry, found := k.GetUnbondingOp(ctx, id)
		if !found {
			panic("did not find UnbondingOp according to index- index was probably not correctly updated")
		}
		entries = append(entries, entry)
	}

	return entries, true
}

// GetMaturedUnbondingOps returns the list of matured unbonding operation ids
func (k Keeper) GetMaturedUnbondingOps(ctx sdk.Context) (ids []uint64, err error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.MaturedUnbondingOpsKey())
	if bz == nil {
		return nil, nil
	}

	var ops ccv.MaturedUnbondingOps
	if err := ops.Unmarshal(bz); err != nil {
		return nil, err
	}
	return ops.GetIds(), nil
}

// AppendMaturedUnbondingOps adds a list of ids to the list of matured unbonding operation ids
func (k Keeper) AppendMaturedUnbondingOps(ctx sdk.Context, ids []uint64) {
	if len(ids) == 0 {
		return
	}
	existingIds, err := k.GetMaturedUnbondingOps(ctx)
	if err != nil {
		panic(fmt.Sprintf("failed to get matured unbonding operations: %s", err))
	}

	maturedOps := ccv.MaturedUnbondingOps{
		Ids: append(existingIds, ids...),
	}

	store := ctx.KVStore(k.storeKey)
	bz, err := maturedOps.Marshal()
	if err != nil {
		panic(fmt.Sprintf("failed to marshal matured unbonding operations: %s", err))
	}
	store.Set(types.MaturedUnbondingOpsKey(), bz)
}

// ConsumeMaturedUnbondingOps empties and returns list of matured unbonding operation ids (if it exists)
func (k Keeper) ConsumeMaturedUnbondingOps(ctx sdk.Context) ([]uint64, error) {
	ids, err := k.GetMaturedUnbondingOps(ctx)
	if err != nil {
		return nil, err
	}
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.MaturedUnbondingOpsKey())
	return ids, nil
}

// Retrieves the underlying client state corresponding to a connection ID.
func (k Keeper) getUnderlyingClient(ctx sdk.Context, connectionID string) (
	clientID string, tmClient *ibctmtypes.ClientState, err error) {

	conn, ok := k.connectionKeeper.GetConnection(ctx, connectionID)
	if !ok {
		return "", nil, sdkerrors.Wrapf(conntypes.ErrConnectionNotFound,
			"connection not found for connection ID: %s", connectionID)
	}
	clientID = conn.ClientId
	clientState, ok := k.clientKeeper.GetClientState(ctx, clientID)
	if !ok {
		return "", nil, sdkerrors.Wrapf(clienttypes.ErrClientNotFound,
			"client not found for client ID: %s", conn.ClientId)
	}
	tmClient, ok = clientState.(*ibctmtypes.ClientState)
	if !ok {
		return "", nil, sdkerrors.Wrapf(clienttypes.ErrInvalidClientType,
			"invalid client type. expected %s, got %s", ibcexported.Tendermint, clientState.ClientType())
	}
	return clientID, tmClient, nil
}

// chanCloseInit defines a wrapper function for the channel Keeper's function
func (k Keeper) chanCloseInit(ctx sdk.Context, channelID string) error {
	capName := host.ChannelCapabilityPath(ccv.ProviderPortID, channelID)
	chanCap, ok := k.scopedKeeper.GetCapability(ctx, capName)
	if !ok {
		return sdkerrors.Wrapf(channeltypes.ErrChannelCapabilityNotFound, "could not retrieve channel capability at: %s", capName)
	}
	return k.channelKeeper.ChanCloseInit(ctx, ccv.ProviderPortID, channelID, chanCap)
}

func (k Keeper) IncrementValidatorSetUpdateId(ctx sdk.Context) {
	validatorSetUpdateId := k.GetValidatorSetUpdateId(ctx)
	k.SetValidatorSetUpdateId(ctx, validatorSetUpdateId+1)
}

func (k Keeper) SetValidatorSetUpdateId(ctx sdk.Context, vscID uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.ValidatorSetUpdateIdKey(), sdk.Uint64ToBigEndian(vscID))
}

func (k Keeper) GetValidatorSetUpdateId(ctx sdk.Context) (vscID uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.ValidatorSetUpdateIdKey())

	if bz == nil {
		vscID = 0
	} else {
		// Unmarshal
		vscID = sdk.BigEndianToUint64(bz)
	}

	return vscID
}

// SetValsetUpdateBlockHeight sets the block height for a given vscID
func (k Keeper) SetValsetUpdateBlockHeight(ctx sdk.Context, vscID, blockHeight uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.ValsetUpdateBlockHeightKey(vscID), sdk.Uint64ToBigEndian(blockHeight))
}

// GetValsetUpdateBlockHeight gets the block height for a given vscID
func (k Keeper) GetValsetUpdateBlockHeight(ctx sdk.Context, vscID uint64) (uint64, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.ValsetUpdateBlockHeightKey(vscID))
	if bz == nil {
		return 0, false
	}
	return sdk.BigEndianToUint64(bz), true
}

// IterateValsetUpdateBlockHeight iterates through the mapping from vscIDs to block heights.
//
// Note that the mapping from vscIDs to block heights is stored under keys with the following format:
// ValsetUpdateBlockHeightBytePrefix | vscID
// Thus, the iteration is in ascending order of vscIDs.
func (k Keeper) IterateValsetUpdateBlockHeight(ctx sdk.Context, cb func(vscID, height uint64) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte{types.ValsetUpdateBlockHeightBytePrefix})

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		vscID := sdk.BigEndianToUint64(iterator.Key()[1:])
		height := sdk.BigEndianToUint64(iterator.Value())

		stop := cb(vscID, height)
		if stop {
			return
		}
	}
}

// DeleteValsetUpdateBlockHeight deletes the block height value for a given vscID
func (k Keeper) DeleteValsetUpdateBlockHeight(ctx sdk.Context, vscID uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.ValsetUpdateBlockHeightKey(vscID))
}

// SetSlashAcks sets the slash acks under the given chain ID
func (k Keeper) SetSlashAcks(ctx sdk.Context, chainID string, acks []string) {
	store := ctx.KVStore(k.storeKey)

	sa := types.SlashAcks{
		Addresses: acks,
	}
	bz, err := sa.Marshal()
	if err != nil {
		panic("failed to marshal SlashAcks")
	}
	store.Set(types.SlashAcksKey(chainID), bz)
}

// GetSlashAcks returns the slash acks stored under the given chain ID
func (k Keeper) GetSlashAcks(ctx sdk.Context, chainID string) []string {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.SlashAcksKey(chainID))
	if bz == nil {
		return nil
	}
	var acks types.SlashAcks
	if err := acks.Unmarshal(bz); err != nil {
		panic(fmt.Errorf("failed to decode json: %w", err))
	}

	return acks.GetAddresses()
}

// ConsumeSlashAcks empties and returns the slash acks for a given chain ID
func (k Keeper) ConsumeSlashAcks(ctx sdk.Context, chainID string) (acks []string) {
	acks = k.GetSlashAcks(ctx, chainID)
	if len(acks) < 1 {
		return
	}
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.SlashAcksKey(chainID))
	return
}

// IterateSlashAcks iterates through the slash acks set in the store.
//
// Note that the slash acks are stored under keys with the following format:
// SlashAcksBytePrefix | chainID
// Thus, the iteration is in ascending order of chainIDs.
//
// Note: This method is only used in testing
func (k Keeper) IterateSlashAcks(ctx sdk.Context, cb func(chainID string, acks []string) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte{types.SlashAcksBytePrefix})

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {

		chainID := string(iterator.Key()[1:])

		var sa types.SlashAcks
		err := sa.Unmarshal(iterator.Value())
		if err != nil {
			panic(fmt.Errorf("failed to unmarshal SlashAcks: %w", err))
		}

		stop := cb(chainID, sa.GetAddresses())
		if stop {
			return
		}
	}
}

// AppendSlashAck appends the given slash ack to the given chain ID slash acks in store
func (k Keeper) AppendSlashAck(ctx sdk.Context, chainID, ack string) {
	acks := k.GetSlashAcks(ctx, chainID)
	acks = append(acks, ack)
	k.SetSlashAcks(ctx, chainID, acks)
}

// SetInitChainHeight sets the provider block height when the given consumer chain was initiated
func (k Keeper) SetInitChainHeight(ctx sdk.Context, chainID string, height uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.InitChainHeightKey(chainID), sdk.Uint64ToBigEndian(height))
}

// GetInitChainHeight returns the provider block height when the given consumer chain was initiated
func (k Keeper) GetInitChainHeight(ctx sdk.Context, chainID string) (uint64, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.InitChainHeightKey(chainID))
	if bz == nil {
		return 0, false
	}

	return sdk.BigEndianToUint64(bz), true
}

// DeleteInitChainHeight deletes the block height value for which the given consumer chain's channel was established
func (k Keeper) DeleteInitChainHeight(ctx sdk.Context, chainID string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.InitChainHeightKey(chainID))
}

// GetPendingVSCPackets returns the list of pending ValidatorSetChange packets stored under chain ID
func (k Keeper) GetPendingVSCPackets(ctx sdk.Context, chainID string) []ccv.ValidatorSetChangePacketData {
	var packets ccv.ValidatorSetChangePackets

	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.PendingVSCsKey(chainID))
	if bz == nil {
		return []ccv.ValidatorSetChangePacketData{}
	}
	if err := packets.Unmarshal(bz); err != nil {
		panic(fmt.Errorf("cannot unmarshal pending validator set changes: %w", err))
	}
	return packets.GetList()
}

// AppendPendingVSCPackets adds the given ValidatorSetChange packet to the list
// of pending ValidatorSetChange packets stored under chain ID
func (k Keeper) AppendPendingVSCPackets(ctx sdk.Context, chainID string, newPackets ...ccv.ValidatorSetChangePacketData) {
	pds := append(k.GetPendingVSCPackets(ctx, chainID), newPackets...)

	store := ctx.KVStore(k.storeKey)
	packets := ccv.ValidatorSetChangePackets{List: pds}
	buf, err := packets.Marshal()
	if err != nil {
		panic(fmt.Errorf("cannot marshal pending validator set changes: %w", err))
	}
	store.Set(types.PendingVSCsKey(chainID), buf)
}

// DeletePendingVSCPackets deletes the list of pending ValidatorSetChange packets for chain ID
func (k Keeper) DeletePendingVSCPackets(ctx sdk.Context, chainID string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.PendingVSCsKey(chainID))
}

// SetConsumerClientId sets the client ID for the given chain ID
func (k Keeper) SetConsumerClientId(ctx sdk.Context, chainID, clientID string) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.ChainToClientKey(chainID), []byte(clientID))
}

// GetConsumerClientId returns the client ID for the given chain ID.
func (k Keeper) GetConsumerClientId(ctx sdk.Context, chainID string) (string, bool) {
	store := ctx.KVStore(k.storeKey)
	clientIdBytes := store.Get(types.ChainToClientKey(chainID))
	if clientIdBytes == nil {
		return "", false
	}
	return string(clientIdBytes), true
}

// DeleteConsumerClientId removes from the store the clientID for the given chainID.
func (k Keeper) DeleteConsumerClientId(ctx sdk.Context, chainID string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.ChainToClientKey(chainID))
}

// SetInitTimeoutTimestamp sets the init timeout timestamp for the given chain ID
func (k Keeper) SetInitTimeoutTimestamp(ctx sdk.Context, chainID string, ts uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.InitTimeoutTimestampKey(chainID), sdk.Uint64ToBigEndian(ts))
}

// GetInitTimeoutTimestamp returns the init timeout timestamp for the given chain ID.
// This method is used only in testing.
func (k Keeper) GetInitTimeoutTimestamp(ctx sdk.Context, chainID string) (uint64, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.InitTimeoutTimestampKey(chainID))
	if bz == nil {
		return 0, false
	}
	return sdk.BigEndianToUint64(bz), true
}

// DeleteInitTimeoutTimestamp removes from the store the init timeout timestamp for the given chainID.
func (k Keeper) DeleteInitTimeoutTimestamp(ctx sdk.Context, chainID string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.InitTimeoutTimestampKey(chainID))
}

// IterateInitTimeoutTimestamp iterates through the init timeout timestamps in the store.
//
// Note that the init timeout timestamps are stored under keys with the following format:
// InitTimeoutTimestampBytePrefix | chainID
// Thus, the iteration is in ascending order of chainIDs (NOT in timestamp order).
func (k Keeper) IterateInitTimeoutTimestamp(ctx sdk.Context, cb func(chainID string, ts uint64) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte{types.InitTimeoutTimestampBytePrefix})

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		chainID := string(iterator.Key()[1:])
		ts := sdk.BigEndianToUint64(iterator.Value())

		stop := cb(chainID, ts)
		if stop {
			return
		}
	}
}

// SetVscSendTimestamp sets the VSC send timestamp
// for a VSCPacket with ID vscID sent to a chain with ID chainID
func (k Keeper) SetVscSendTimestamp(
	ctx sdk.Context,
	chainID string,
	vscID uint64,
	timestamp time.Time,
) {
	store := ctx.KVStore(k.storeKey)

	// Convert timestamp into bytes for storage
	timeBz := sdk.FormatTimeBytes(timestamp)

	store.Set(types.VscSendingTimestampKey(chainID, vscID), timeBz)
}

// GetVscSendTimestamp returns a VSC send timestamp by chainID and vscID
//
// Note: This method is used only for testing.
func (k Keeper) GetVscSendTimestamp(ctx sdk.Context, chainID string, vscID uint64) (time.Time, bool) {
	store := ctx.KVStore(k.storeKey)

	timeBz := store.Get(types.VscSendingTimestampKey(chainID, vscID))
	if timeBz == nil {
		return time.Time{}, false
	}

	ts, err := sdk.ParseTimeBytes(timeBz)
	if err != nil {
		return time.Time{}, false
	}
	return ts, true
}

// DeleteVscSendTimestamp removes from the store a specific VSC send timestamp
// for the given chainID and vscID.
func (k Keeper) DeleteVscSendTimestamp(ctx sdk.Context, chainID string, vscID uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.VscSendingTimestampKey(chainID, vscID))
}

// IterateVscSendTimestamps iterates over the vsc send timestamps of the given chainID.
//
// Note that the vsc send timestamps of a given chainID are stored under keys with the following format:
// VscSendTimestampBytePrefix | len(chainID) | chainID | vscID
// Thus, the iteration is in ascending order of vscIDs, and as a result in send timestamp order.
func (k Keeper) IterateVscSendTimestamps(
	ctx sdk.Context,
	chainID string,
	cb func(vscID uint64, ts time.Time) (stop bool),
) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.ChainIdWithLenKey(types.VscSendTimestampBytePrefix, chainID))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		key := iterator.Key()
		_, vscID, err := types.ParseVscSendingTimestampKey(key)
		if err != nil {
			panic(fmt.Errorf("failed to parse VscSendTimestampKey: %w", err))
		}
		ts, err := sdk.ParseTimeBytes(iterator.Value())
		if err != nil {
			panic(fmt.Errorf("failed to parse timestamp value: %w", err))
		}

		stop := cb(vscID, ts)
		if stop {
			return
		}
	}
}
