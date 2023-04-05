package ante

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypesv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
)

type ForbiddenProposalsDecorator struct {
	IsProposalWhitelisted func(govtypesv1beta1.Content) bool
}

func NewForbiddenProposalsDecorator(whiteListFn func(govtypesv1beta1.Content) bool) ForbiddenProposalsDecorator {
	return ForbiddenProposalsDecorator{IsProposalWhitelisted: whiteListFn}
}

func (decorator ForbiddenProposalsDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	currHeight := ctx.BlockHeight()

	for _, msg := range tx.GetMsgs() {
		submitProposalMgs, ok := msg.(*govtypesv1beta1.MsgSubmitProposal)
		//if the message is MsgSubmitProposal, check if proposal is whitelisted
		if ok {
			if !decorator.IsProposalWhitelisted(submitProposalMgs.GetContent()) {
				return ctx, fmt.Errorf("tx contains unsupported proposal message types at height %d", currHeight)
			}
		}
	}

	return next(ctx, tx, simulate)
}
