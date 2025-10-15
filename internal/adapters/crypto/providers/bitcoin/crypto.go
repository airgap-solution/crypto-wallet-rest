package bitcoin

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil"
	hd "github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/lamengao/go-electrum/electrum"
)

const (
	SatoshiPerBTC = 1e8
	AccountDepth  = 3
	ChainDepth    = 4
	ExternalChain = 0
	ChangeChain   = 1
)

var (
	ErrIndexOutOfRange   = errors.New("index out of range for uint32")
	BitcoinMainNetParams = &chaincfg.MainNetParams
	BitcoinTestNetParams = &chaincfg.TestNet3Params
)

func addressToScripthash(addr string, isTestnet bool) (string, error) {
	params := BitcoinMainNetParams
	if isTestnet {
		params = BitcoinTestNetParams
	}

	a, err := btcutil.DecodeAddress(addr, params)
	if err != nil {
		return "", fmt.Errorf("failed to decode address: %w", err)
	}
	script, err := txscript.PayToAddrScript(a)
	if err != nil {
		return "", fmt.Errorf("failed to create script: %w", err)
	}
	h := sha256.Sum256(script)
	for i, j := 0, len(h)-1; i < j; i, j = i+1, j-1 {
		h[i], h[j] = h[j], h[i]
	}
	return hex.EncodeToString(h[:]), nil
}

func deriveTaprootAddresses(
	xpub string, externalCount, changeCount int, isTestnet bool,
) ([]btcutil.Address, []btcutil.Address, error) {
	key, err := hd.NewKeyFromString(xpub)
	if err != nil {
		return nil, nil, fmt.Errorf("bad xpub: %w", err)
	}

	extRoot, chRoot, err := deriveChainKeys(key)
	if err != nil {
		return nil, nil, err
	}

	external, err := deriveAddresses(extRoot, externalCount, isTestnet)
	if err != nil {
		return nil, nil, fmt.Errorf("derive external addresses: %w", err)
	}

	change, err := deriveAddresses(chRoot, changeCount, isTestnet)
	if err != nil {
		return nil, nil, fmt.Errorf("derive change addresses: %w", err)
	}

	return external, change, nil
}

func deriveChainKeys(key *hd.ExtendedKey) (*hd.ExtendedKey, *hd.ExtendedKey, error) {
	var extRoot *hd.ExtendedKey
	var chRoot *hd.ExtendedKey
	var err error

	switch key.Depth() {
	case AccountDepth:
		extRoot, err = key.Derive(ExternalChain)
		if err != nil {
			return nil, nil, fmt.Errorf("derive external chain: %w", err)
		}
		chRoot, err = key.Derive(ChangeChain)
		if err != nil {
			return nil, nil, fmt.Errorf("derive change chain: %w", err)
		}
	case ChainDepth:
		extRoot = key
		chRoot = nil
	default:
		extRoot = key
		chRoot = nil
	}

	return extRoot, chRoot, nil
}

func deriveAddresses(root *hd.ExtendedKey, count int, isTestnet bool) ([]btcutil.Address, error) {
	if root == nil {
		return nil, nil
	}

	addresses := make([]btcutil.Address, 0, count)
	for i := range count {
		if i < 0 || i > 0x7FFFFFFF { // Check for uint32 overflow
			return nil, ErrIndexOutOfRange
		}

		child, err := root.Derive(uint32(i))
		if err != nil {
			return nil, fmt.Errorf("derive child %d: %w", i, err)
		}

		pub, err := child.ECPubKey()
		if err != nil {
			return nil, fmt.Errorf("get public key for child %d: %w", i, err)
		}

		addr, err := makeTaprootAddress(pub, isTestnet)
		if err != nil {
			return nil, fmt.Errorf("make taproot address for child %d: %w", i, err)
		}

		addresses = append(addresses, addr)
	}

	return addresses, nil
}

func makeTaprootAddress(pub *btcec.PublicKey, isTestnet bool) (*btcutil.AddressTaproot, error) {
	params := BitcoinMainNetParams
	if isTestnet {
		params = BitcoinTestNetParams
	}

	outputKey := txscript.ComputeTaprootKeyNoScript(pub)
	ser := schnorr.SerializePubKey(outputKey)
	addr, err := btcutil.NewAddressTaproot(ser, params)
	if err != nil {
		return nil, fmt.Errorf("create taproot address: %w", err)
	}
	return addr, nil
}

func getXpubBalance(node *electrum.Client, addresses []btcutil.Address, isTestnet bool) (float64, error) {
	totalSats := int64(0)

	for _, addr := range addresses {
		sats, err := getAddressBalance(node, addr, isTestnet)
		if err != nil {
			return 0, err
		}
		totalSats += sats
	}

	return float64(totalSats) / SatoshiPerBTC, nil
}

func getAddressBalance(node *electrum.Client, addr btcutil.Address, isTestnet bool) (int64, error) {
	sh, err := addressToScripthash(addr.EncodeAddress(), isTestnet)
	if err != nil {
		return 0, err
	}

	balResp, err := node.GetBalance(context.Background(), sh)
	if err != nil {
		return 0, fmt.Errorf("get balance from electrum: %w", err)
	}

	confirmed := balResp.Confirmed
	unconfirmed := balResp.Unconfirmed
	return int64(confirmed) + int64(unconfirmed), nil
}
