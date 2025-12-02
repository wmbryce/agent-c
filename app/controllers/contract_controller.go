package controllers

import (
	"fmt"
	"math/big"

	"github.com/gofiber/fiber/v2"
	"github.com/wmbryce/agent-c/app/store/blockchain"
	"github.com/wmbryce/agent-c/app/types"
	"github.com/wmbryce/agent-c/app/utils"
)

// GetBalance func gets the ETH balance of an address.
// @Description Get ETH balance of an address.
// @Summary get ETH balance
// @Tags Blockchain
// @Accept json
// @Produce json
// @Param address body string true "Ethereum address"
// @Success 200 {object} models.GetBalanceResponse
// @Security ApiKeyAuth
// @Router /v1/blockchain/balance [post]
func GetBalance(c *fiber.Ctx) error {
	// Create new GetBalanceRequest struct
	request := &types.GetBalanceRequest{}

	// Check, if received JSON data is valid.
	if err := c.BodyParser(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Validate request fields.
	validate := utils.NewValidator()
	if err := validate.Struct(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   utils.ValidatorErrors(err),
		})
	}

	// Create Ethereum client
	ethClient, err := blockchain.NewEthereumClient()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   fmt.Sprintf("failed to connect to Ethereum: %v", err),
		})
	}
	defer ethClient.Close()

	// Get balance
	balance, err := ethClient.GetBalance(request.Address)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Convert to ETH
	balanceEth := utils.WeiToEth(balance)

	response := types.GetBalanceResponse{
		Address:    request.Address,
		Balance:    balance.String(),
		BalanceEth: balanceEth.Text('f', 18),
	}

	return c.JSON(fiber.Map{
		"error":    false,
		"msg":      nil,
		"response": response,
	})
}

// GetTransactionReceipt func gets the receipt of a transaction.
// @Description Get transaction receipt by hash.
// @Summary get transaction receipt
// @Tags Blockchain
// @Accept json
// @Produce json
// @Param tx_hash body string true "Transaction hash"
// @Success 200 {object} models.TransactionReceiptResponse
// @Security ApiKeyAuth
// @Router /v1/blockchain/receipt [post]
func GetTransactionReceipt(c *fiber.Ctx) error {
	// Create new TransactionReceiptRequest struct
	request := &types.TransactionReceiptRequest{}

	// Check, if received JSON data is valid.
	if err := c.BodyParser(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Validate request fields.
	validate := utils.NewValidator()
	if err := validate.Struct(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   utils.ValidatorErrors(err),
		})
	}

	// Create Ethereum client
	ethClient, err := blockchain.NewEthereumClient()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   fmt.Sprintf("failed to connect to Ethereum: %v", err),
		})
	}
	defer ethClient.Close()

	// Get transaction receipt
	receipt, err := ethClient.GetTransactionReceipt(request.TxHash)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	contractAddress := ""
	if receipt.ContractAddress.Hex() != "0x0000000000000000000000000000000000000000" {
		contractAddress = receipt.ContractAddress.Hex()
	}

	// Note: From and To are not in the receipt, would need to fetch the transaction
	// For now, we'll leave them empty
	response := types.TransactionReceiptResponse{
		TxHash:          receipt.TxHash.Hex(),
		BlockNumber:     receipt.BlockNumber.Uint64(),
		BlockHash:       receipt.BlockHash.Hex(),
		GasUsed:         receipt.GasUsed,
		Status:          receipt.Status,
		ContractAddress: contractAddress,
		From:            "",
		To:              "",
	}

	return c.JSON(fiber.Map{
		"error":    false,
		"msg":      nil,
		"response": response,
	})
}

// GetBlockInfo func gets current blockchain information.
// @Description Get current block number and chain ID.
// @Summary get blockchain info
// @Tags Blockchain
// @Accept json
// @Produce json
// @Success 200 {object} models.BlockInfoResponse
// @Security ApiKeyAuth
// @Router /v1/blockchain/info [get]
func GetBlockInfo(c *fiber.Ctx) error {
	// Create Ethereum client
	ethClient, err := blockchain.NewEthereumClient()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   fmt.Sprintf("failed to connect to Ethereum: %v", err),
		})
	}
	defer ethClient.Close()

	// Get block number
	blockNumber, err := ethClient.GetBlockNumber()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	response := types.BlockInfoResponse{
		BlockNumber: blockNumber,
		ChainID:     ethClient.ChainID.Int64(),
	}

	return c.JSON(fiber.Map{
		"error":    false,
		"msg":      nil,
		"response": response,
	})
}

// SendEther func sends ETH to an address.
// @Description Send ETH to an address.
// @Summary send ETH
// @Tags Blockchain
// @Accept json
// @Produce json
// @Param to body string true "Recipient address"
// @Param amount body string true "Amount in ETH"
// @Success 200 {object} map[string]interface{}
// @Security ApiKeyAuth
// @Router /v1/blockchain/send-ether [post]
func SendEther(c *fiber.Ctx) error {
	// Define request structure
	var request struct {
		To     string `json:"to" validate:"required,eth_addr"`
		Amount string `json:"amount" validate:"required"`
	}

	// Check, if received JSON data is valid.
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Validate request fields.
	validate := utils.NewValidator()
	if err := validate.Struct(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   utils.ValidatorErrors(err),
		})
	}

	// Create Ethereum client
	ethClient, err := blockchain.NewEthereumClient()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   fmt.Sprintf("failed to connect to Ethereum: %v", err),
		})
	}
	defer ethClient.Close()

	// Parse amount
	amountFloat := new(big.Float)
	amountFloat, ok := amountFloat.SetString(request.Amount)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   "invalid amount format",
		})
	}

	// Convert to wei
	amountWei := utils.EthToWei(amountFloat)

	// Note: This is a placeholder. In production, you'd implement the actual transaction sending logic
	// using ethClient.GetTransactOpts() and creating a transaction

	return c.JSON(fiber.Map{
		"error": false,
		"msg":   "Transaction would be sent (implementation needed)",
		"response": fiber.Map{
			"to":     request.To,
			"amount": amountWei.String() + " wei",
		},
	})
}

