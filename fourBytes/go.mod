module github.com/TrueBlocks/trueblocks-core/examples/fourBytes/v5

// Go Version
go 1.25.1

require (
	github.com/ethereum/go-ethereum v1.15.10
	github.com/panjf2000/ants/v2 v2.11.3
	github.com/spf13/cobra v1.9.1
)

replace github.com/TrueBlocks/trueblocks-core/src/apps/chifra => ../../src/apps/chifra

require (
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.4.0 // indirect
	github.com/holiman/uint256 v1.3.2 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/pflag v1.0.6 // indirect
	golang.org/x/crypto v0.37.0 // indirect
	golang.org/x/sync v0.13.0 // indirect
	golang.org/x/sys v0.32.0 // indirect
)
