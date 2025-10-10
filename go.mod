module github.com/airgap-solution/crypto-wallet-rest

go 1.25.1

require (
	github.com/airgap-solution/cmc-rest/openapi v1.0.1
	github.com/airgap-solution/crypto-wallet-rest/openapi/servergen/go v0.0.0-00010101000000-000000000000
	github.com/btcsuite/btcd v0.24.2
	github.com/btcsuite/btcd/btcec/v2 v2.3.5
	github.com/btcsuite/btcd/btcutil v1.1.6
	github.com/lamengao/go-electrum v0.0.0-20231031090039-0e19b90480c4
	github.com/restartfu/gophig v0.0.2
	github.com/samber/lo v1.52.0
	github.com/stretchr/testify v1.11.1
	go.uber.org/mock v0.6.0
)

replace github.com/airgap-solution/crypto-wallet-rest/openapi/servergen/go => ./openapi/servergen/go

require (
	github.com/btcsuite/btcd/chaincfg/chainhash v1.1.0 // indirect
	github.com/btcsuite/btclog v0.0.0-20170628155309-84c8d2346e9f // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/decred/dcrd/crypto/blake256 v1.0.0 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.0.1 // indirect
	github.com/goccy/go-json v0.10.2 // indirect
	github.com/gorilla/mux v1.8.1 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/pelletier/go-toml/v2 v2.0.7 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9 // indirect
	golang.org/x/sys v0.0.0-20200814200057-3d37ad5750ed // indirect
	golang.org/x/text v0.22.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
