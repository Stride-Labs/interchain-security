package types

import (
	ccv "github.com/cosmos/interchain-security/provider/x/ccv/common_types"
	consumertypes "github.com/cosmos/interchain-security/provider/x/consumer/types"
)

func NewConsumerStates(
	chainID,
	clientID,
	channelID string,
	initialHeight uint64,
	genesis consumertypes.GenesisState,
	unbondingOpsIndexes []VscUnbondingOps,
	pendingValsetChanges []ccv.ValidatorSetChangePacketData,
	slashDowntimeAck []string,
) ConsumerState {
	return ConsumerState{
		ChainId:              chainID,
		ClientId:             clientID,
		ChannelId:            channelID,
		InitialHeight:        initialHeight,
		UnbondingOpsIndex:    unbondingOpsIndexes,
		PendingValsetChanges: pendingValsetChanges,
		ConsumerGenesis:      genesis,
		SlashDowntimeAck:     slashDowntimeAck,
	}
}
