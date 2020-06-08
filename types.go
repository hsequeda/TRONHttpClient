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

type AddressWithoutPrivKey struct {
	Base58CheckAddress string `json:"base58checkAddress"`
	Value              string `json:"value"`
}

type Account struct {
	Address               string          `json:"address"`
	Balance               int             `json:"balance"`
	Frozen                []Frozen        `json:"frozen"`
	CreateTime            int             `json:"create_time"`
	LatestOperationTime   int             `json:"latest_opration_time"`
	FreeNetUsage          int             `json:"free_net_usage"`
	LatestConsumeFreeTime int             `json:"latest_consume_free_time"`
	AccountResource       AccountResource `json:"account_resource"`
	OwnerPermission       Permission      `json:"owner_permission"`
	ActivePermission      []Permission    `json:"active_permission"`
	AssetV2               []Asset         `json:"assetV2"`
	FreeAssetNetUsageV2   []Asset         `json:"free_asset_net_usageV2"`
}

type Permission struct {
	Id             int    `json:"id"`
	Type           string `json:"type"`
	Operations     string `json:"operations"`
	PermissionName string `json:"permission_name"`
	Threshold      int    `json:"threshold"`
	Keys           []Key  `json:"keys"`
}

type Key struct {
	Address string `json:"address"`
	Weight  int    `json:"weight"`
}

type Frozen struct {
	FrozenBalance int `json:"frozen_balance"`
	ExpireTime    int `json:"expire_time"`
}

type AccountResource struct {
	FrozenBalanceForEnergy     Frozen `json:"frozen_balance_for_energy"`
	LatestConsumeTimeForEnergy int    `json:"latest_consume_time_for_energy"`
}

type Asset struct {
	Key   string `json:"key"`
	Value int    `json:"value"`
}
