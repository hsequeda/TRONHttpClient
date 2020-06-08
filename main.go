package tronhttpClient

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
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
func (c *Client) BroadcastTx(tx *Transaction) (*Transaction, error) {
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

	err = json.NewDecoder(resp).Decode(tx)
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
	encodeData, err := json.Marshal(map[string]string{
		"value": hex.EncodeToString([]byte(password)),
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
func (c *Client) BroadcastHex(txHex string) (*Transaction, error) {
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
	err = json.NewDecoder(resp).Decode(&tx)
	if err != nil {
		return nil, err
	}

	return &tx, nil
}

// EasyTransfer Easily transfer from an address using the password string.
// Only works with accounts created from createAddress,integrated getransactionsign and broadcasttransaction.
func (c *Client) EasyTransfer(password, toAddress string, amount int) (*Transaction, error) {
	encodeData, err := json.Marshal(
		map[string]interface{}{
			"passPhrase": hex.EncodeToString([]byte(password)),
			"toAddress":  toAddress,
			"amount":     amount,
		})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", testNet+"/wallet/easytransfer",
		bytes.NewBuffer(encodeData))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	resp, err := c.client.CallRetryable(req)
	if err != nil {
		return nil, err
	}

	var result struct {
		Result struct {
			Ok      bool   `json:"result"`
			Code    string `json:"code"`
			Message string `json:"message"`
		}
		Transaction Transaction `json:"transaction"`
	}
	err = json.NewDecoder(resp).Decode(&result)
	if err != nil {
		return nil, err
	}

	if result.Result.Ok {
		return &result.Transaction, nil
	} else {
		b, err := hex.DecodeString(result.Result.Message)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("%s: %s", result.Result.Code, string(b))
	}
}

// EasyTransferByPrivate Easily transfer from an address using the private key.
func (c *Client) EasyTransferByPrivate(privateKey, toAddress string, amount int) (*Transaction, error) {
	encodeData, err := json.Marshal(
		map[string]interface{}{
			"privateKey": privateKey,
			"toAddress":  toAddress,
			"amount":     amount,
		})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", testNet+"/wallet/easytransferbyprivate",
		bytes.NewBuffer(encodeData))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	resp, err := c.client.CallRetryable(req)
	if err != nil {
		return nil, err
	}

	var result struct {
		Result struct {
			Ok      bool   `json:"result"`
			Code    string `json:"code"`
			Message string `json:"message"`
		}
		Transaction Transaction `json:"transaction"`
	}
	err = json.NewDecoder(resp).Decode(&result)
	if err != nil {
		return nil, err
	}

	if result.Result.Ok {
		return &result.Transaction, nil
	} else {
		b, err := hex.DecodeString(result.Result.Message)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("%s: %s", result.Result.Code, string(b))
	}
}

// CreateAccount Create an account. Uses an already activated account to create a new account
// Note: The expiration time of the http api creation transaction is 1 minute,
//       so to complete the on-chain, you need to complete gettransactionsign and
//       broadcasttransaction within 1 minute after the creation.
func (c *Client) CreateAccount(ownerAddr, accountAddr string, visible bool, permissionID int) (*Transaction, error) {
	encodeData, err := json.Marshal(
		map[string]interface{}{
			"owner_address":   ownerAddr,
			"account_address": accountAddr,
			"visible":         visible,
			"permission_id":   visible,
		})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", testNet+"/wallet/createaccount",
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

// GetAccount Query information about an account,Including balances, freezes, votes and time, etc
func (c *Client) GetAccount(address string, visible bool) (*Account, error) {
	encodeData, err := json.Marshal(
		map[string]interface{}{
			"address": address,
			"visible": visible,
		})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", testNet+"/wallet/getaccount",
		bytes.NewBuffer(encodeData))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	resp, err := c.client.CallRetryable(req)
	if err != nil {
		return nil, err
	}

	var account Account
	err = json.NewDecoder(resp).Decode(&account)
	if err != nil {
		return nil, err
	}

	return &account, err
}
