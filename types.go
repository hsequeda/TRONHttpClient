package tronhttpClient

type Transaction struct {
	Visible    bool        `json:"visible"`
	TxId       string      `json:"txID"`
	RawData    interface{} `json:"raw_data"`
	RawDataHex string      `json:"raw_data_hex"`
	Signature  []string    `json:"signature"`
}

type Address struct {
	PrivateKey string `json:"privateKey"`
	Address    string `json:"address"`
	HexAddress string `json:"hexAddress"`
}
