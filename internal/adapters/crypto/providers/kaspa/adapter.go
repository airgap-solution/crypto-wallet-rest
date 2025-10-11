package kaspa

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
)

type Adapter struct {
	explorerURL string
	cache       map[string][]string
	mu          sync.RWMutex
}

func NewAdapter(explorerURL string) *Adapter {
	return &Adapter{
		explorerURL: explorerURL,
		cache:       make(map[string][]string),
	}
}

func (a *Adapter) GetBalance(kpub string) (float64, error) {
	a.mu.RLock()
	addresses, ok := a.cache[kpub]
	a.mu.RUnlock()

	if !ok {
		recv, change, err := deriveAddresses(kpub, 1000, 1000)
		if err != nil {
			return 0, err
		}

		addresses = append(recv, change...)

		a.mu.Lock()
		a.cache[kpub] = addresses
		a.mu.Unlock()
	}

	res, err := a.fetchBalances(addresses)
	if err != nil {
		return 0, err
	}

	var bal float64
	for _, r := range res {
		bal += r.Balance / 1e8
	}
	return bal, nil
}

type balanceResponse struct {
	Address string  `json:"address"`
	Balance float64 `json:"balance"`
}

func (a *Adapter) fetchBalances(addresses []string) ([]balanceResponse, error) {
	payload := map[string][]string{
		"addresses": addresses,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	respBody, err := postJSON(a.explorerURL+"/addresses/balances", data)
	if err != nil {
		return nil, err
	}

	var result []balanceResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return result, nil
}

func postJSON(url string, data []byte) ([]byte, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}

	return io.ReadAll(resp.Body)
}
