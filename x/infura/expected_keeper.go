package infura

import evm "github.com/fibonacci-chain/fbc/x/evm/watcher"

type EvmKeeper interface {
	SetObserverKeeper(keeper evm.InfuraKeeper)
}
