#!/bin/bash
# Generate secure secrets for Akash deployment

set -e

generate_secret() {
    openssl rand -base64 32 | tr -d '=/+' | head -c 32
}

echo "Generating secure secrets for Akash deployment..."
echo ""

DB_PASSWORD=$(generate_secret)
REDIS_PASSWORD=$(generate_secret)
JWT_SECRET_KEY=$(generate_secret)
JWT_REFRESH_KEY=$(generate_secret)

cat > secrets.env << EOF
# Akash Deployment Secrets
# Generated on $(date)
# DO NOT commit this file to version control

DB_PASSWORD=${DB_PASSWORD}
REDIS_PASSWORD=${REDIS_PASSWORD}
JWT_SECRET_KEY=${JWT_SECRET_KEY}
JWT_REFRESH_KEY=${JWT_REFRESH_KEY}
EOF

echo "Secrets written to secrets.env"
echo ""
echo "Generated values:"
echo "  DB_PASSWORD:     ${DB_PASSWORD}"
echo "  REDIS_PASSWORD:  ${REDIS_PASSWORD}"
echo "  JWT_SECRET_KEY:  ${JWT_SECRET_KEY}"
echo "  JWT_REFRESH_KEY: ${JWT_REFRESH_KEY}"
echo ""
echo "IMPORTANT: Keep these secrets safe and never commit secrets.env to git!"
