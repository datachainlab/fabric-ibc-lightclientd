package main

import (
	"github.com/cosmos/cosmos-sdk/store/dbadapter"
	sdk "github.com/cosmos/cosmos-sdk/types"
	fabibctypes "github.com/datachainlab/fabric-ibc/x/ibc/light-clients/xx-fabric/types"
	dbm "github.com/tendermint/tm-db"
)

type Lightclientd struct {
	store sdk.KVStore
	cs    *fabibctypes.ClientState
}

func NewLightclientd(name, dir string) *Lightclientd {
	db, err := dbm.NewDB(name, dbm.GoLevelDBBackend, dir)
	if err != nil {
		panic(err)
	}

	return &Lightclientd{
		store: dbadapter.Store{DB: db},
	}
}
