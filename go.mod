module github.com/dig-coins/btcconnect

go 1.20

require (
	github.com/0xb10c/rawtx v1.5.0
	github.com/btcsuite/btcd v0.24.0
	github.com/dig-coins/hd-wallet v0.0.0-20240312121459-3253c499232c
	github.com/okx/go-wallet-sdk/coins/bitcoin v0.0.0-00010101000000-000000000000
	github.com/redis/rueidis v1.0.31
	github.com/sgostarter/i v0.1.16
	github.com/sgostarter/libconfig v0.0.2
	github.com/sgostarter/libeasygo v0.1.70
	github.com/spf13/cast v1.6.0
	github.com/stretchr/testify v1.8.4
	golang.org/x/exp v0.0.0-20231110203233-9a3e6036ecaa
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/bits-and-blooms/bitset v1.7.0 // indirect
	github.com/btcsuite/btcd/btcec/v2 v2.3.2 // indirect
	github.com/btcsuite/btcd/btcutil v1.1.5 // indirect
	github.com/btcsuite/btcd/btcutil/psbt v1.1.8 // indirect
	github.com/btcsuite/btcd/chaincfg/chainhash v1.1.0 // indirect
	github.com/btcsuite/btclog v0.0.0-20170628155309-84c8d2346e9f // indirect
	github.com/btcsuite/btcutil v1.0.2 // indirect
	github.com/consensys/bavard v0.1.13 // indirect
	github.com/consensys/gnark-crypto v0.12.1 // indirect
	github.com/crate-crypto/go-kzg-4844 v0.3.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/decred/dcrd/crypto/blake256 v1.0.0 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.0.1 // indirect
	github.com/ethereum/c-kzg-4844 v0.3.1 // indirect
	github.com/ethereum/go-ethereum v1.13.4 // indirect
	github.com/go-stack/stack v1.8.1 // indirect
	github.com/holiman/uint256 v1.2.3 // indirect
	github.com/kelseyhightower/envconfig v1.4.0 // indirect
	github.com/mmcloughlin/addchain v0.4.0 // indirect
	github.com/okx/go-wallet-sdk/crypto v0.0.1 // indirect
	github.com/okx/go-wallet-sdk/util v0.0.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/shengdoushi/base58 v1.0.0 // indirect
	github.com/supranational/blst v0.3.11 // indirect
	github.com/tyler-smith/go-bip39 v1.1.0 // indirect
	golang.org/x/crypto v0.16.0 // indirect
	golang.org/x/sync v0.3.0 // indirect
	golang.org/x/sys v0.17.0 // indirect
	google.golang.org/protobuf v1.30.0 // indirect
	rsc.io/tmplfunc v0.0.3 // indirect
)

replace github.com/okx/go-wallet-sdk/coins/bitcoin => github.com/dig-coins/go-wallet-sdk/coins/bitcoin v0.0.0-20240326121340-9543ab2b2240

//replace github.com/sgostarter/libeasygo => ../../work_sgostarter/libeasygo

//replace github.com/sgostarter/libconfig => ../../work_sgostarter/libconfig

//replace github.com/okx/go-wallet-sdk/coins/bitcoin => ../../work_dig-coins/go-wallet-sdk/coins/bitcoin
