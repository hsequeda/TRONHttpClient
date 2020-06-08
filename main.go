package tronhttpClient

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	httpClient "github.com/stdevHsequeda/TRONHttpClient/client"
	"net/http"
)

const testNet = "https://api.shasta.trongrid.io"
const mainNet = "https://api.trongrid.io"

type Client struct {
	client  *httpClient.Client
	network string
}

// NewClient returns a new instance of Client
func NewClient(network string) *Client {
	httpClient.MaxRetry = 5
	return &Client{client: httpClient.NewClient(), network: network}
}

// CreateTx Create a TRX transfer transaction.
// If toAddr does not exist, then create the account on the blockchain.
func (c *Client) CreateTx(toAddr, ownerAddr string, amount int) (*Transaction, error) {
	encodeData, err := json.Marshal(
		map[string]interface{}{
			"to_address":    toAddr,
			"owner_address": ownerAddr,
			"amount":        amount,
		})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", testNet+"/wallet/createtransaction",
		bytes.NewBuffer(encodeData))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	resp, err := c.client.CallRetryable(req)
	if err != nil {
		return nil, err
	}

	var tx Transaction
	err = json.NewDecoder(resp).Decode(&tx)
	if err != nil {
		return nil, err
	}

	return &tx, err
}

// GetTxSign Sign the transaction, the api has the risk of leaking the private key,
// please make sure to call the api in a secure environment
func (c *Client) GetTxSign(tx *Transaction, privKey string) (*Transaction, error) {
	encodeData, err := json.Marshal(
		struct {
			Transaction *Transaction `json:"transaction"`
			PrivateKey  string       `json:"privateKey"`
		}{
			Transaction: tx,
			PrivateKey:  privKey,
		})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", testNet+"/wallet/gettransactionsign",
		bytes.NewBuffer(encodeData))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	resp, err := c.client.CallRetryable(req)
	if err != nil {
		return nil, err
	}

	err = json.NewDecoder(resp).Decode(tx)
	if err != nil {
		return nil, err
	}

	return tx, err
}

// BroadcastTx  Broadcast the signed transaction
func (c *Client) BroadcastTx(tx *Transaction) (*Transaction,error) {
	encodeData, err := json.Marshal(tx)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", testNet+"/wallet/broadcasttransaction",
		bytes.NewBuffer(encodeData))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	resp, err := c.client.CallRetryable(req)
	if err != nil {
		return nil, err
	}

   err= json.NewDecoder(resp).Decode(tx)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

// GenerateAddress Generates a random private key and address pair. Returns a private key,
// the corresponding address in hex, and base58.
func (c *Client) GenerateAddress() (*Address, error) {
	req, err := http.NewRequest("GET", testNet+"/wallet/generateaddress",
		nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	resp, err := c.client.CallRetryable(req)
	if err != nil {
		return nil, err
	}

	var addr Address
	err = json.NewDecoder(resp).Decode(&addr)
	if err != nil {
		return nil, err
	}

	return &addr, nil
}

// CreateAddress Create address from a specified password string (NOT PRIVATE KEY)
func (c *Client) CreateAddress(password string) (*AddressWithoutPrivKey, error) {
	hexPass, err := hex.DecodeString(password)
	if err != nil {
		return nil, err
	}

	encodeData, err := json.Marshal(map[string]string{
		"value": string(hexPass),
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", testNet+"/wallet/createaddress", bytes.NewBuffer(encodeData))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	resp, err := c.client.CallRetryable(req)
	if err != nil {
		return nil, err
	}

	var addr AddressWithoutPrivKey
	err = json.NewDecoder(resp).Decode(&addr)
	if err != nil {
		return nil, err
	}

	return &addr, nil
}

// ValidateAddress Validates address, returns either true or false.
func (c *Client) ValidateAddress(address string) (bool, error) {
	encodeData, err := json.Marshal(map[string]string{
		"address": string(address),
	})
	if err != nil {
		return false, err
	}

	req, err := http.NewRequest("GET", testNet+"/wallet/validateaddress", bytes.NewBuffer(encodeData))
	if err != nil {
		return false, err
	}

	req.Header.Add("Content-Type", "application/json")
	resp, err := c.client.CallRetryable(req)
	if err != nil {
		return false, err
	}

	var addr struct {
		ok bool `json:"result"`
	}
	err = json.NewDecoder(resp).Decode(&addr)
	if err != nil {
		return false, err
	}

	return addr.ok, nil
}

// BroadcastHex Broadcast the protobuf encoded transaction hex string after sign
func (c *Client) BroadcastHex(txHex string) (*Transaction,error){
 	encodeData, err := json.Marshal(
      map[string]string{
         "transaction": txHex,
      })
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", testNet+"/wallet/broadcasthex",
		bytes.NewBuffer(encodeData))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	resp, err := c.client.CallRetryable(req)
	if err != nil {
		return nil, err
	}

   var tx Transaction
   err= json.NewDecoder(resp).Decode(&tx)
	if err != nil {
		return nil, err
	}

	return &tx, nil
}



