package slashing

import (
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/slashing/exported"
	"github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	"github.com/cosmos/cosmos-sdk/x/slashing/types"
)

var (
	_ module.AppModule           = AppModule{}
	_ module.AppModuleBasic      = AppModuleBasic{}
	_ module.AppModuleSimulation = AppModule{}
)

// AppModule embeds the Cosmos SDK's x/slashing AppModuleBasic.
type AppModuleBasic struct {
	slashing.AppModuleBasic
}

// AppModule embeds the Cosmos SDK's x/slashing AppModule where we only override
// specific methods.
type AppModule struct {
	// embed the Cosmos SDK's x/slashing AppModule
	slashing.AppModule
}

// NewAppModule creates a new AppModule object using the native x/slashing module
// AppModule constructor.
func NewAppModule(cdc codec.Codec,
	keeper keeper.Keeper,
	ak types.AccountKeeper,
	bk types.BankKeeper, sk types.StakingKeeper, subspace exported.Subspace) AppModule {
	slashingAppMod := slashing.NewAppModule(cdc, keeper, ak, bk, sk, subspace)

	return AppModule{
		AppModule: slashingAppMod,
	}
}

// BeginBlock returns the begin blocker for the slashing module.
func (am AppModule) BeginBlock(ctx sdk.Context, req abci.RequestBeginBlock) {
}
