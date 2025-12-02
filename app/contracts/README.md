# Smart Contracts

This directory contains Solidity smart contracts for your project.

## How to Generate Go Bindings

### Prerequisites

Install `abigen` (part of go-ethereum):

```bash
go install github.com/ethereum/go-ethereum/cmd/abigen@latest
```

Install Solidity compiler:

```bash
# macOS
brew install solidity

# Or download from https://soliditylang.org
```

### Step 1: Compile Contract

Compile your Solidity contract to generate ABI and bytecode:

```bash
cd contracts

# Compile to get ABI and bytecode
solc --abi --bin ERC20.sol -o build/
```

This creates:
- `build/ERC20Token.abi` - Contract interface
- `build/ERC20Token.bin` - Contract bytecode

### Step 2: Generate Go Bindings

Use `abigen` to generate Go code from the ABI:

```bash
# Generate Go bindings
abigen --abi=build/ERC20Token.abi \
       --bin=build/ERC20Token.bin \
       --pkg=contracts \
       --type=ERC20Token \
       --out=../app/contracts/erc20_token.go
```

This creates a Go file with type-safe contract bindings!

### Step 3: Use in Your Code

```go
package controllers

import (
    "github.com/create-go-app/fiber-go-template/app/contracts"
    "github.com/create-go-app/fiber-go-template/platform/blockchain"
    "github.com/ethereum/go-ethereum/common"
)

func InteractWithERC20(contractAddress string) {
    // Connect to blockchain
    ethClient, _ := blockchain.NewEthereumClient()
    defer ethClient.Close()
    
    // Load contract
    address := common.HexToAddress(contractAddress)
    token, err := contracts.NewERC20Token(address, ethClient.Client)
    
    // Read contract state (no gas cost)
    name, _ := token.Name(ethClient.GetCallOpts())
    balance, _ := token.BalanceOf(ethClient.GetCallOpts(), common.HexToAddress("0x123..."))
    
    // Write to contract (costs gas)
    auth, _ := ethClient.GetTransactOpts()
    tx, _ := token.Transfer(auth, common.HexToAddress("0x456..."), big.NewInt(1000))
}
```

## Example: Full Workflow

### 1. Deploy a New Contract

```go
// Get transaction options
auth, _ := ethClient.GetTransactOpts()

// Deploy contract
address, tx, instance, err := contracts.DeployERC20Token(
    auth,
    ethClient.Client,
    "MyToken",      // name
    "MTK",          // symbol
    18,             // decimals
    big.NewInt(1000000), // total supply
)
```

### 2. Interact with Deployed Contract

```go
// Load existing contract
token, _ := contracts.NewERC20Token(
    common.HexToAddress("0xContractAddress..."),
    ethClient.Client,
)

// Read operations (free)
symbol, _ := token.Symbol(ethClient.GetCallOpts())
totalSupply, _ := token.TotalSupply(ethClient.GetCallOpts())

// Write operations (costs gas)
auth, _ := ethClient.GetTransactOpts()
tx, _ := token.Transfer(
    auth,
    common.HexToAddress("0xRecipient..."),
    big.NewInt(100),
)
```

## Directory Structure

```
contracts/
├── README.md           # This file
├── ERC20.sol          # Example ERC20 token contract
├── YourContract.sol   # Your custom contracts here
└── build/             # Compiled contracts (generated)
    ├── *.abi
    └── *.bin

../app/contracts/      # Generated Go bindings go here
└── erc20_token.go     # Generated from ERC20.sol
```

## Adding Your Own Contracts

1. Create `YourContract.sol` in this directory
2. Compile: `solc --abi --bin YourContract.sol -o build/`
3. Generate bindings: `abigen --abi=build/YourContract.abi --bin=build/YourContract.bin --pkg=contracts --out=../app/contracts/your_contract.go`
4. Import and use in controllers!

## Common Contract Types

- **ERC20**: Fungible tokens (already included as example)
- **ERC721**: NFTs
- **ERC1155**: Multi-token standard
- **Custom**: Your own business logic

## Testing Contracts

Use local development networks:
- **Hardhat**: `npx hardhat node`
- **Ganache**: `ganache-cli`
- **Anvil** (Foundry): `anvil`

Then set `ETHEREUM_RPC_URL=http://localhost:8545` in your `.env`

