package utils

import (
	"sort"
	"time"

	abci "github.com/cometbft/cometbft/abci/types"
	tmprotocrypto "github.com/cometbft/cometbft/proto/tendermint/crypto"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
	host "github.com/cosmos/ibc-go/v7/modules/core/24-host"
	ccv "github.com/cosmos/interchain-security/provider/x/ccv/common_types"
)

func AccumulateChanges(currentChanges, newChanges []abci.ValidatorUpdate) []abci.ValidatorUpdate {
	m := make(map[string]abci.ValidatorUpdate)

	for i := 0; i < len(currentChanges); i++ {
		m[currentChanges[i].PubKey.String()] = currentChanges[i]
	}

	for i := 0; i < len(newChanges); i++ {
		m[newChanges[i].PubKey.String()] = newChanges[i]
	}

	var out []abci.ValidatorUpdate

	for _, update := range m {
		out = append(out, update)
	}

	// The list of tendermint updates should hash the same across all consensus nodes
	// that means it is necessary to sort for determinism.
	sort.Slice(out, func(i, j int) bool {
		if out[i].Power != out[j].Power {
			return out[i].Power > out[j].Power
		}
		return out[i].PubKey.String() > out[j].PubKey.String()
	})

	return out
}

// TMCryptoPublicKeyToConsAddr converts a TM public key to an SDK public key
// and returns the associated consensus address
func TMCryptoPublicKeyToConsAddr(k tmprotocrypto.PublicKey) (sdk.ConsAddress, error) {
	sdkK, err := cryptocodec.FromTmProtoPublicKey(k)
	if err != nil {
		return nil, err
	}
	return sdk.GetConsAddress(sdkK), nil
}

// SendIBCPacket sends an IBC packet with packetData
// over the source channelID and portID
func SendIBCPacket(
	ctx sdk.Context,
	scopedKeeper ccv.ScopedKeeper,
	channelKeeper ccv.ChannelKeeper,
	channelID string,
	portID string,
	packetData []byte,
	timeoutPeriod time.Duration,
) error {
	channelCap, ok := scopedKeeper.GetCapability(ctx, host.ChannelCapabilityPath(portID, channelID))
	if !ok {
		return sdkerrors.Wrap(channeltypes.ErrChannelCapabilityNotFound, "module does not own channel capability")
	}

	_, err := channelKeeper.SendPacket(ctx,
		channelCap,
		portID, channelID,
		clienttypes.Height{},
		uint64(ctx.BlockTime().Add(timeoutPeriod).UnixNano()),
		packetData,
	)
	return err
}

// AppendMany appends a variable number of byte slices together
func AppendMany(byteses ...[]byte) (out []byte) {
	for _, bytes := range byteses {
		out = append(out, bytes...)
	}
	return out
}