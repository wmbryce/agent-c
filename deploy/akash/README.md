# Akash Deployment Guide

This guide covers deploying Agent-C to the Akash decentralized cloud.

## Prerequisites

1. **Akash CLI** installed
2. **AKT tokens** in your wallet (for deployment costs)
3. **Docker image** pushed to GitHub Container Registry

## Quick Start

### 1. Install Akash CLI

```bash
# macOS
brew install akash

# Linux
curl -sSfL https://raw.githubusercontent.com/akash-network/provider/main/install.sh | sh
```

### 2. Set Up Wallet

```bash
# Create new wallet
akash keys add deployer

# Or import existing wallet
akash keys add deployer --recover

# Get your address
akash keys show deployer -a
```

Fund your wallet with AKT tokens from an exchange or the Akash faucet.

### 3. Generate Secrets

```bash
cd deploy/akash
./generate-secrets.sh
```

This creates `secrets.env` with secure random values.

### 4. Prepare SDL with Secrets

Before deploying, substitute the secrets into the SDL:

```bash
# Load secrets
source secrets.env

# Create deployment SDL with secrets substituted
envsubst < deploy.yaml > deploy-final.yaml
```

### 5. Create Deployment

```bash
# Set environment
export AKASH_NODE=https://rpc.akashnet.net:443
export AKASH_CHAIN_ID=akashnet-2
export AKASH_FROM=deployer
export AKASH_KEYRING_BACKEND=os

# Create certificate (first time only)
akash tx cert create client --from deployer

# Create deployment
akash tx deployment create deploy-final.yaml --from deployer

# Note the DSEQ (deployment sequence) from the output
export AKASH_DSEQ=<your-dseq>
```

### 6. Accept a Bid

```bash
# List bids for your deployment
akash query market bid list --owner $(akash keys show deployer -a) --dseq $AKASH_DSEQ

# Accept a bid (use the provider address from the bid list)
akash tx market lease create \
  --dseq $AKASH_DSEQ \
  --gseq 1 \
  --oseq 1 \
  --provider <provider-address> \
  --from deployer
```

### 7. Send Manifest

```bash
akash provider send-manifest deploy-final.yaml \
  --dseq $AKASH_DSEQ \
  --provider <provider-address> \
  --from deployer
```

### 8. Get Deployment URL

```bash
akash provider lease-status \
  --dseq $AKASH_DSEQ \
  --provider <provider-address> \
  --from deployer
```

The output will include the URL to access your deployment.

### 9. Run Migrations

After deployment, run database migrations:

```bash
# Get a shell into the backend container (via provider)
# Then run:
./goose -dir ./migrations postgres "postgres://agentc:${DB_PASSWORD}@postgres:5432/agentc?sslmode=disable" up
```

## Using Akash Console (Alternative)

For a GUI experience, use [Akash Console](https://console.akash.network):

1. Connect your Keplr wallet
2. Click "Deploy"
3. Upload `deploy-final.yaml`
4. Select a provider and complete deployment

## Updating Deployment

To update your deployment with a new image:

```bash
# Update the deployment
akash tx deployment update deploy-final.yaml \
  --dseq $AKASH_DSEQ \
  --from deployer

# Send updated manifest
akash provider send-manifest deploy-final.yaml \
  --dseq $AKASH_DSEQ \
  --provider <provider-address> \
  --from deployer
```

## Closing Deployment

```bash
akash tx deployment close --dseq $AKASH_DSEQ --from deployer
```

## Cost Estimation

Based on the SDL configuration:
- **Backend**: ~2 CPU, 2GB RAM
- **Postgres**: ~1 CPU, 1GB RAM, 10GB persistent storage
- **Redis**: ~0.5 CPU, 512MB RAM

Estimated cost: **$15-25/month** in AKT (varies by provider and AKT price)

## Troubleshooting

### Check Deployment Status

```bash
akash query deployment get --owner $(akash keys show deployer -a) --dseq $AKASH_DSEQ
```

### View Logs

```bash
akash provider lease-logs \
  --dseq $AKASH_DSEQ \
  --provider <provider-address> \
  --from deployer \
  --service backend
```

### Common Issues

1. **Bid not received**: Increase pricing in SDL or try different region
2. **Container not starting**: Check logs, ensure image is publicly accessible
3. **Database connection failed**: Verify postgres service is healthy before backend starts

## Files

- `deploy.yaml` - Main SDL deployment file
- `secrets.env.example` - Template for secrets
- `generate-secrets.sh` - Script to generate secure secrets

## Security Notes

- Never commit `secrets.env` or `deploy-final.yaml` to git
- Rotate secrets periodically
- Use strong, unique passwords for all services
