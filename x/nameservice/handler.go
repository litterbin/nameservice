package nameservice

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler returns a handler for "nameservice" type messages
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgSetName:
			return handleMsgSetName(ctx, keeper, msg)
		case MsgBuyName:
			return handleMsgBuyName(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized nameservice Msg type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle MsgSetName
func handleMsgSetName(ctx sdk.Context, keeper Keeper, msg MsgSetName) sdk.Result {
	if !msg.Owner.Equals(keeper.GetOwner(ctx, msg.Name)) { // Checks if the msg sender is the same as the current owner
		return sdk.ErrUnauthorized("Incorrect Owner").Result() // If not, throw an error
	}
	keeper.SetName(ctx, msg.Name, msg.Value)
	return sdk.Result{}
}

func handleMsgBuyName(ctx sdk.Context, keeper Keeper, msg MsgBuyName) sdk.Result {
	if keeper.GetPrice(ctx, msg.Name).IsAllGT(msg.Bid) { // Checks if the bid price is greater than the price paid by the current owner
		return sdk.ErrInsufficientCoins("Bid not high enough").Result() // If not, throw an error
	}
	if keeper.HasOwner(ctx, msg.Name) {
		_, err := keeper.coinKeeper.SendCoins(ctx, msg.Buyer, keeper.GetOwner(ctx, msg.Name), msg.Bid)
		if err != nil {
			return sdk.ErrInsufficientCoins("Buyer does not have enough coins").Result()
		}
	} else {
		_, _, err := keeper.coinKeeper.SubtractCoins(ctx, msg.Buyer, msg.Bid)
		if err != nil {
			return sdk.ErrInsufficientCoins("Buyer does not have enough coins").Result()
		}
	}
	keeper.SetOwner(ctx, msg.Name, msg.Buyer)
	keeper.SetPrice(ctx, msg.Name, msg.Bid)

	return sdk.Result{}
}
