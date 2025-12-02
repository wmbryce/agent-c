#!/bin/bash

# Script to compile Solidity contracts and generate Go bindings
# Usage: ./scripts/compile_contracts.sh [contract_name]

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if solc is installed
if ! command -v solc &> /dev/null; then
    echo -e "${RED}Error: solc (Solidity compiler) is not installed${NC}"
    echo "Install with: brew install solidity"
    exit 1
fi

# Check if abigen is installed
if ! command -v abigen &> /dev/null; then
    echo -e "${RED}Error: abigen is not installed${NC}"
    echo "Install with: go install github.com/ethereum/go-ethereum/cmd/abigen@latest"
    exit 1
fi

# Navigate to contracts directory
cd "$(dirname "$0")/../contracts"

# Create build directory if it doesn't exist
mkdir -p build
mkdir -p ../app/contracts

# If no argument provided, compile all .sol files
if [ -z "$1" ]; then
    echo -e "${YELLOW}Compiling all contracts...${NC}"
    for contract in *.sol; do
        if [ -f "$contract" ]; then
            contract_name="${contract%.sol}"
            echo -e "${GREEN}Compiling $contract_name...${NC}"
            
            # Compile contract
            solc --abi --bin "$contract" -o build/ --overwrite
            
            # Find the main contract name (usually matches filename)
            for abi_file in build/*.abi; do
                base_name=$(basename "$abi_file" .abi)
                bin_file="build/${base_name}.bin"
                
                if [ -f "$bin_file" ]; then
                    echo -e "${GREEN}Generating Go bindings for $base_name...${NC}"
                    
                    # Convert to snake_case for Go filename
                    go_file=$(echo "$base_name" | sed 's/\([A-Z]\)/_\L\1/g' | sed 's/^_//')
                    
                    # Generate Go bindings
                    abigen --abi="$abi_file" \
                           --bin="$bin_file" \
                           --pkg=contracts \
                           --type="$base_name" \
                           --out="../app/contracts/${go_file}.go"
                    
                    echo -e "${GREEN}✓ Generated: app/contracts/${go_file}.go${NC}"
                fi
            done
        fi
    done
else
    # Compile specific contract
    contract="$1"
    if [ ! -f "$contract.sol" ]; then
        echo -e "${RED}Error: Contract $contract.sol not found${NC}"
        exit 1
    fi
    
    echo -e "${GREEN}Compiling $contract...${NC}"
    
    # Compile contract
    solc --abi --bin "$contract.sol" -o build/ --overwrite
    
    # Generate Go bindings for the specific contract
    for abi_file in build/${contract}*.abi; do
        if [ -f "$abi_file" ]; then
            base_name=$(basename "$abi_file" .abi)
            bin_file="build/${base_name}.bin"
            
            if [ -f "$bin_file" ]; then
                echo -e "${GREEN}Generating Go bindings for $base_name...${NC}"
                
                # Convert to snake_case for Go filename
                go_file=$(echo "$base_name" | sed 's/\([A-Z]\)/_\L\1/g' | sed 's/^_//')
                
                # Generate Go bindings
                abigen --abi="$abi_file" \
                       --bin="$bin_file" \
                       --pkg=contracts \
                       --type="$base_name" \
                       --out="../app/contracts/${go_file}.go"
                
                echo -e "${GREEN}✓ Generated: app/contracts/${go_file}.go${NC}"
            fi
        fi
    done
fi

echo -e "${GREEN}✓ All done!${NC}"
echo ""
echo "Next steps:"
echo "1. Import the generated contract in your controller:"
echo "   import \"github.com/create-go-app/fiber-go-template/app/contracts\""
echo "2. Use the contract methods in your code"
echo "3. See BLOCKCHAIN_SETUP.md for examples"

