package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/supply/types"
)

// Keeper
type Keeper struct {
	*bank.BaseViewKeeper

	cdc        *codec.Codec
	storeKey   sdk.StoreKey
	ak         auth.AccountKeeper
	paramSpace params.Subspace
}

// NewKeeper creates a new Keeper instance
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, ak auth.AccountKeeper, codespace sdk.CodespaceType, paramSpace params.Subspace) Keeper {
	baseViewKeeper := bank.NewBaseViewKeeper(ak, codespace)
	return Keeper{
		&baseViewKeeper,
		cdc,
		key,
		ak,
		paramSpace,
	}
}

// GetSupply retrieves the Supply from store
func (k Keeper) GetSupply(ctx sdk.Context) (supply types.Supply) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(supplyKey)
	if b == nil {
		panic("Stored supply should not have been nil")
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &supply)
	return
}

// SetSupply sets the Supply to store
func (k Keeper) SetSupply(ctx sdk.Context, supply types.Supply) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshalBinaryLengthPrefixed(supply)
	store.Set(supplyKey, b)
}

// Inflate increases the total supply amount
func (k Keeper) Inflate(ctx sdk.Context, amount sdk.Coins) {
	supply := k.GetSupply(ctx)
	supply.Inflate(amount)
	k.SetSupply(ctx, supply)
}

// Deflate reduces the total supply amount
func (k Keeper) Deflate(ctx sdk.Context, amount sdk.Coins) {
	supply := k.GetSupply(ctx)
	supply.Deflate(amount)
	k.SetSupply(ctx, supply)
}