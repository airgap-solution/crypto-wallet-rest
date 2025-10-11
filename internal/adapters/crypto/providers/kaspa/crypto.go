package kaspa

import (
	"fmt"
	"log"

	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/kaspanet/kaspad/util"
)

func deriveAddresses(xpub string, nRecv, nChange int) (recvAddrs, changeAddrs []string, err error) {
	key, err := hdkeychain.NewKeyFromString(xpub)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid xpub: %w", err)
	}

	prefix := util.Bech32PrefixKaspa

	recvBranch, err := key.Derive(0)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to derive receive branch: %w", err)
	}
	changeBranch, err := key.Derive(1)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to derive change branch: %w", err)
	}

	deriveAddrs := func(branch *hdkeychain.ExtendedKey, count int) []string {
		addrs := make([]string, 0, count)
		for i := range count {
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
			if len(pubKeyBytes) != 33 {
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
		return addrs
	}

	recvAddrs = deriveAddrs(recvBranch, nRecv)
	changeAddrs = deriveAddrs(changeBranch, nChange)

	return recvAddrs, changeAddrs, nil
}
