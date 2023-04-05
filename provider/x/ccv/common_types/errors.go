package common_types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrInvalidVSCMaturedId      = sdkerrors.Register(ModuleName, 11, "invalid vscId for VSC packet")
	ErrInvalidVSCMaturedTime    = sdkerrors.Register(ModuleName, 12, "invalid maturity time for VSC packet")
	ErrInvalidPacketData        = sdkerrors.Register(ModuleName, 13, "invalid CCV packet data")
	ErrInvalidVersion           = sdkerrors.Register(ModuleName, 14, "invalid CCV version")
	ErrInvalidChannelFlow       = sdkerrors.Register(ModuleName, 15, "invalid message sent to channel end")
	ErrInvalidGenesis           = sdkerrors.Register(ModuleName, 16, "invalid genesis state")
	ErrDuplicateChannel         = sdkerrors.Register(ModuleName, 17, "CCV channel already exists")
	ErrInvalidConsumerClient    = sdkerrors.Register(ModuleName, 18, "ccv channel is not built on correct client")
	ErrInvalidHandshakeMetadata = sdkerrors.Register(ModuleName, 19, "invalid provider handshake metadata")
	ErrChannelNotFound          = sdkerrors.Register(ModuleName, 20, "channel not found")
	ErrClientNotFound           = sdkerrors.Register(ModuleName, 21, "client not found")
	ErrDuplicateConsumerChain   = sdkerrors.Register(ModuleName, 22, "consumer chain already exists")
	ErrConsumerChainNotFound    = sdkerrors.Register(ModuleName, 23, "consumer chain not found")
)
