package kaspa

import (
	"errors"
	"fmt"
	"log"

	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/kaspanet/kaspad/util"
)

const (
	SompiPerKAS = 1e8
)

var (
	ErrIndexOutOfRange = errors.New("index out of range for uint32")
)

func deriveAddresses(xpub string, nRecv, nChange int) ([]string, []string, error) {
	key, err := hdkeychain.NewKeyFromString(xpub)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid xpub: %w", err)
	}

	recvBranch, err := key.Derive(0)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to derive receive branch: %w", err)
	}
	changeBranch, err := key.Derive(1)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to derive change branch: %w", err)
	}

	recvAddrs, err := derive(recvBranch, nRecv)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to derive receive addresses: %w", err)
	}
	changeAddrs, err := derive(changeBranch, nChange)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to derive change addresses: %w", err)
	}

	return recvAddrs, changeAddrs, nil
}
func derive(branch *hdkeychain.ExtendedKey, count int) ([]string, error) {
	prefix := util.Bech32PrefixKaspa
	pubKeyLength := 33
	addrs := make([]string, 0, count)
	for i := range count {
		if i < 0 || i > 0x7FFFFFFF { // Check for uint32 overflow
			return nil, ErrIndexOutOfRange
		}

		child, err := branch.Derive(uint32(i))
		if err != nil {
			log.Printf("Skipping index %d due to derivation error: %s\n", i, err)
			continue
		}

		pubKey, err := child.ECPubKey()
		if err != nil {
			log.Printf("Skipping index %d due to pubkey error: %s\n", i, err)
			continue
		}

		pubKeyBytes := pubKey.SerializeCompressed()
		if len(pubKeyBytes) != pubKeyLength {
			log.Printf("Skipping index %d, pubkey length %d not 33\n", i, len(pubKeyBytes))
			continue
		}

		schnorrPubKey := pubKeyBytes[1:]

		addr, err := util.NewAddressPublicKey(schnorrPubKey, prefix)
		if err != nil {
			log.Printf("Skipping index %d, failed to create address: %s\n", i, err)
			continue
		}

		addrs = append(addrs, addr.EncodeAddress())
	}
	return addrs, nil
}
