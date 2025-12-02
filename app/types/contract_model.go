package types

// ContractCallRequest struct for read-only contract calls
type ContractCallRequest struct {
	ContractAddress string   `json:"contract_address" validate:"required,eth_addr"`
	MethodName      string   `json:"method_name" validate:"required"`
	Parameters      []string `json:"parameters"`
}

// ContractTransactionRequest struct for state-changing contract calls
type ContractTransactionRequest struct {
	ContractAddress string   `json:"contract_address" validate:"required,eth_addr"`
	MethodName      string   `json:"method_name" validate:"required"`
	Parameters      []string `json:"parameters"`
	GasLimit        *uint64  `json:"gas_limit,omitempty"`
	Value           *string  `json:"value,omitempty"` // in wei
}

// GetBalanceRequest struct to get ETH balance
type GetBalanceRequest struct {
	Address string `json:"address" validate:"required,eth_addr"`
}

// GetBalanceResponse struct for balance response
type GetBalanceResponse struct {
	Address string `json:"address"`
	Balance string `json:"balance"` // in wei
	BalanceEth string `json:"balance_eth"` // in ETH
}

// TransactionReceiptRequest struct to get transaction receipt
type TransactionReceiptRequest struct {
	TxHash string `json:"tx_hash" validate:"required"`
}

// TransactionReceiptResponse struct for transaction receipt
type TransactionReceiptResponse struct {
	TxHash          string `json:"tx_hash"`
	BlockNumber     uint64 `json:"block_number"`
	BlockHash       string `json:"block_hash"`
	GasUsed         uint64 `json:"gas_used"`
	Status          uint64 `json:"status"` // 1 = success, 0 = failure
	ContractAddress string `json:"contract_address,omitempty"`
	From            string `json:"from"`
	To              string `json:"to"`
}

// DeployContractRequest struct for deploying new contracts
type DeployContractRequest struct {
	Bytecode     string   `json:"bytecode" validate:"required"`
	ABI          string   `json:"abi" validate:"required"`
	ConstructorArgs []string `json:"constructor_args"`
	GasLimit     *uint64  `json:"gas_limit,omitempty"`
}

// DeployContractResponse struct for deployment response
type DeployContractResponse struct {
	TxHash          string `json:"tx_hash"`
	ContractAddress string `json:"contract_address"`
}

// BlockInfoResponse struct for blockchain info
type BlockInfoResponse struct {
	BlockNumber uint64 `json:"block_number"`
	ChainID     int64  `json:"chain_id"`
}

// ERC20TransferRequest struct for ERC20 token transfers
type ERC20TransferRequest struct {
	TokenAddress string `json:"token_address" validate:"required,eth_addr"`
	To           string `json:"to" validate:"required,eth_addr"`
	Amount       string `json:"amount" validate:"required"` // in token units
}

// ERC20BalanceRequest struct for ERC20 token balance
type ERC20BalanceRequest struct {
	TokenAddress string `json:"token_address" validate:"required,eth_addr"`
	Address      string `json:"address" validate:"required,eth_addr"`
}

// ERC20BalanceResponse struct for ERC20 balance response
type ERC20BalanceResponse struct {
	TokenAddress string `json:"token_address"`
	Address      string `json:"address"`
	Balance      string `json:"balance"`
	Symbol       string `json:"symbol,omitempty"`
	Decimals     uint8  `json:"decimals,omitempty"`
}

