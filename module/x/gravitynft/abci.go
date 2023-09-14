package gravitynft

import (
	"github.com/Gravity-Bridge/Gravity-Bridge/module/x/gravitynft/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// EndBlocker is called at the end of every block
func EndBlocker(ctx sdk.Context, k keeper.Keeper) {
	// TODO: Implement me
	panic("EndBlocker not implemented yet")
}
