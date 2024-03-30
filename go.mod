module github.com/dig-coins/btcconnect

go 1.20

require (
	github.com/0xb10c/rawtx v1.5.0
	github.com/btcsuite/btcd v0.24.0
	github.com/dig-coins/hd-wallet v0.0.0-20240312121459-3253c499232c
	github.com/gin-contrib/gzip v1.0.0
	github.com/gin-contrib/requestid v1.0.0
	github.com/gin-gonic/gin v1.9.1
	github.com/okx/go-wallet-sdk/coins/bitcoin v0.0.0-00010101000000-000000000000
	github.com/redis/rueidis v1.0.31
	github.com/sgostarter/i v0.1.16
	github.com/sgostarter/libconfig v0.0.2
	github.com/sgostarter/libeasygo v0.1.70
	github.com/shopspring/decimal v1.3.1
	github.com/spf13/cast v1.6.0
	github.com/stretchr/testify v1.9.0
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
	github.com/bytedance/sonic v1.11.3 // indirect
	github.com/chenzhuoyu/base64x v0.0.0-20230717121745-296ad89f973d // indirect
	github.com/chenzhuoyu/iasm v0.9.1 // indirect
	github.com/consensys/bavard v0.1.13 // indirect
	github.com/consensys/gnark-crypto v0.12.1 // indirect
	github.com/crate-crypto/go-kzg-4844 v0.3.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/decred/dcrd/crypto/blake256 v1.0.0 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.0.1 // indirect
	github.com/ethereum/c-kzg-4844 v0.3.1 // indirect
	github.com/ethereum/go-ethereum v1.13.4 // indirect
	github.com/gabriel-vasile/mimetype v1.4.3 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.19.0 // indirect
	github.com/go-stack/stack v1.8.1 // indirect
	github.com/goccy/go-json v0.10.2 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/holiman/uint256 v1.2.3 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/kelseyhightower/envconfig v1.4.0 // indirect
	github.com/klauspost/cpuid/v2 v2.2.7 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mmcloughlin/addchain v0.4.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/okx/go-wallet-sdk/crypto v0.0.1 // indirect
	github.com/okx/go-wallet-sdk/util v0.0.1 // indirect
	github.com/pelletier/go-toml/v2 v2.2.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/shengdoushi/base58 v1.0.0 // indirect
	github.com/supranational/blst v0.3.11 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/tyler-smith/go-bip39 v1.1.0 // indirect
	github.com/ugorji/go/codec v1.2.12 // indirect
	golang.org/x/arch v0.7.0 // indirect
	golang.org/x/crypto v0.21.0 // indirect
	golang.org/x/net v0.22.0 // indirect
	golang.org/x/sync v0.3.0 // indirect
	golang.org/x/sys v0.18.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	google.golang.org/protobuf v1.33.0 // indirect
	rsc.io/tmplfunc v0.0.3 // indirect
)

//replace github.com/okx/go-wallet-sdk/coins/bitcoin => github.com/dig-coins/go-wallet-sdk/coins/bitcoin v0.0.0-20240326121340-9543ab2b2240

//replace github.com/sgostarter/libeasygo => ../../work_sgostarter/libeasygo

//replace github.com/sgostarter/libconfig => ../../work_sgostarter/libconfig

replace github.com/okx/go-wallet-sdk/coins/bitcoin => ../../work_dig-coins/go-wallet-sdk/coins/bitcoin
