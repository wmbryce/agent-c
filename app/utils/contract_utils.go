package utils

import (
	"math/big"
	"regexp"
)

// IsValidEthereumAddress checks if a string is a valid Ethereum address
func IsValidEthereumAddress(address string) bool {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	return re.MatchString(address)
}

// WeiToEth converts wei (smallest unit) to ETH
func WeiToEth(wei *big.Int) *big.Float {
	if wei == nil {
		return big.NewFloat(0)
	}
	
	// 1 ETH = 10^18 wei
	ethValue := new(big.Float).SetInt(wei)
	ethDenominator := new(big.Float).SetInt(big.NewInt(1e18))
	
	return new(big.Float).Quo(ethValue, ethDenominator)
}

// EthToWei converts ETH to wei (smallest unit)
func EthToWei(eth *big.Float) *big.Int {
	if eth == nil {
		return big.NewInt(0)
	}
	
	// 1 ETH = 10^18 wei
	weiMultiplier := new(big.Float).SetInt(big.NewInt(1e18))
	weiValue := new(big.Float).Mul(eth, weiMultiplier)
	
	result := new(big.Int)
	weiValue.Int(result)
	return result
}

// GweiToWei converts Gwei to wei
func GweiToWei(gwei int64) *big.Int {
	// 1 Gwei = 10^9 wei
	return big.NewInt(gwei * 1e9)
}

// IsValidTransactionHash checks if a string is a valid transaction hash
func IsValidTransactionHash(hash string) bool {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{64}$")
	return re.MatchString(hash)
}

