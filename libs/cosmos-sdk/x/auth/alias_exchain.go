package auth

import (
	"github.com/fibonacci-chain/fbc/libs/cosmos-sdk/x/auth/exported"
	"github.com/fibonacci-chain/fbc/libs/cosmos-sdk/x/auth/keeper"
)

type (
	Account       = exported.Account
	ModuleAccount = exported.ModuleAccount
	ObserverI     = keeper.ObserverI
)
