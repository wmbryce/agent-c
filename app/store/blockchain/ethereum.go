package blockchain

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// EthereumClient represents a connection to an Ethereum node
type EthereumClient struct {
	Client     *ethclient.Client
	ChainID    *big.Int
	PrivateKey *ecdsa.PrivateKey
	Address    common.Address
}

// NewEthereumClient creates a new Ethereum client connection
func NewEthereumClient() (*EthereumClient, error) {
	// Get Ethereum RPC URL from environment
	rpcURL := os.Getenv("ETHEREUM_RPC_URL")
	if rpcURL == "" {
		return nil, fmt.Errorf("ETHEREUM_RPC_URL environment variable is not set")
	}

	// Connect to Ethereum node
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ethereum node: %v", err)
	}

	// Get chain ID
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get chain ID: %v", err)
	}

	ethClient := &EthereumClient{
		Client:  client,
		ChainID: chainID,
	}

	// Load private key if available (optional for read-only operations)
	privateKeyHex := os.Getenv("ETHEREUM_PRIVATE_KEY")
	if privateKeyHex != "" {
		privateKey, err := crypto.HexToECDSA(privateKeyHex)
		if err != nil {
			return nil, fmt.Errorf("failed to load private key: %v", err)
		}

		publicKey := privateKey.Public()
		publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
		if !ok {
			return nil, fmt.Errorf("failed to cast public key to ECDSA")
		}

		address := crypto.PubkeyToAddress(*publicKeyECDSA)

		ethClient.PrivateKey = privateKey
		ethClient.Address = address
	}

	return ethClient, nil
}

// GetBalance returns the ETH balance of an address
func (ec *EthereumClient) GetBalance(address string) (*big.Int, error) {
	addr := common.HexToAddress(address)
	balance, err := ec.Client.BalanceAt(context.Background(), addr, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get balance: %v", err)
	}
	return balance, nil
}

// GetTransactionReceipt returns the receipt of a transaction
func (ec *EthereumClient) GetTransactionReceipt(txHash string) (*types.Receipt, error) {
	hash := common.HexToHash(txHash)
	receipt, err := ec.Client.TransactionReceipt(context.Background(), hash)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction receipt: %v", err)
	}
	return receipt, nil
}

// GetTransactOpts creates transaction options for sending transactions
func (ec *EthereumClient) GetTransactOpts() (*bind.TransactOpts, error) {
	if ec.PrivateKey == nil {
		return nil, fmt.Errorf("private key not loaded")
	}

	nonce, err := ec.Client.PendingNonceAt(context.Background(), ec.Address)
	if err != nil {
		return nil, fmt.Errorf("failed to get nonce: %v", err)
	}

	gasPrice, err := ec.Client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get gas price: %v", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(ec.PrivateKey, ec.ChainID)
	if err != nil {
		return nil, fmt.Errorf("failed to create transactor: %v", err)
	}

	auth.Nonce = new(big.Int).SetUint64(nonce)
	auth.Value = big.NewInt(0)     // in wei
	auth.GasLimit = uint64(300000) // default gas limit
	auth.GasPrice = gasPrice

	return auth, nil
}

// GetCallOpts creates call options for read-only contract calls
func (ec *EthereumClient) GetCallOpts() *bind.CallOpts {
	return &bind.CallOpts{
		Pending: false,
		Context: context.Background(),
	}
}

// Close closes the Ethereum client connection
func (ec *EthereumClient) Close() {
	ec.Client.Close()
}

// GetBlockNumber returns the latest block number
func (ec *EthereumClient) GetBlockNumber() (uint64, error) {
	blockNumber, err := ec.Client.BlockNumber(context.Background())
	if err != nil {
		return 0, fmt.Errorf("failed to get block number: %v", err)
	}
	return blockNumber, nil
}

// EstimateGas estimates gas for a transaction
func (ec *EthereumClient) EstimateGas(to common.Address, data []byte) (uint64, error) {
	msg := ethereum.CallMsg{
		From: ec.Address,
		To:   &to,
		Data: data,
	}

	gasLimit, err := ec.Client.EstimateGas(context.Background(), msg)
	if err != nil {
		return 0, fmt.Errorf("failed to estimate gas: %v", err)
	}
	return gasLimit, nil
}

