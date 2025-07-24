package taiko_genesis

import (
	_ "embed"
)

//go:embed moonchain_mainnet.json
var MoonchainMainnetGenesisAllocJSON []byte

//go:embed moonchain_hudson.json
var MoonchainHudsonGenesisAllocJSON []byte
