package keeper

import (
	store2 "github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/nft"
)

// SaveClass defines a method for creating a new nft class
func (k Keeper) SaveClass(ctx sdk.Context, class nft.Class) error {
	if k.HasClass(ctx, class.Id) {
		return sdkerrors.Wrap(nft.ErrClassExists, class.Id)
	}
	bz, err := k.cdc.Marshal(&class)
	if err != nil {
		return sdkerrors.Wrap(err, "Marshal nft.Class failed")
	}
	store := k.getStore(ctx)
	store.Set(classStoreKey(class.Id), bz)
	return nil
}

// UpdateClass defines a method for updating a exist nft class
func (k Keeper) UpdateClass(ctx sdk.Context, class nft.Class) error {
	if !k.HasClass(ctx, class.Id) {
		return sdkerrors.Wrap(nft.ErrClassNotExists, class.Id)
	}
	bz, err := k.cdc.Marshal(&class)
	if err != nil {
		return sdkerrors.Wrap(err, "Marshal nft.Class failed")
	}
	store := k.getStore(ctx)
	store.Set(classStoreKey(class.Id), bz)
	return nil
}

func (k Keeper) decodeClass(bz []byte) (nft.Class, bool) {
	var class nft.Class
	if len(bz) == 0 {
		return class, false
	}
	k.cdc.MustUnmarshal(bz, &class)
	return class, true
}

// GetClass defines a method for returning the class information of the specified id
func (k Keeper) GetClass(ctx sdk.Context, classID string) (nft.Class, bool) {
	store := ctx.KVStore(k.storeKey)
	class, boolval := store2.GetAndDecodeWithBool(store, k.decodeClass, classStoreKey(classID))
	if !boolval {
		return class, boolval
	}
	return class, boolval
}

// GetClasses defines a method for returning all classes information
func (k Keeper) GetClasses(ctx sdk.Context) (classes []*nft.Class) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, ClassKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var class nft.Class
		k.cdc.MustUnmarshal(iterator.Value(), &class)
		classes = append(classes, &class)
	}
	return
}

// HasClass determines whether the specified classID exist
func (k Keeper) HasClass(ctx sdk.Context, classID string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(classStoreKey(classID))
}
