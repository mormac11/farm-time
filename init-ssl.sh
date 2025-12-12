#!/bin/bash

# SSL Initialization Script for donaldson.farm
# Run this once on first deployment to get SSL certificates

set -e

DOMAIN="donaldson.farm"
EMAIL="${1:-admin@donaldson.farm}"  # Pass email as first argument or use default

echo "=== SSL Certificate Setup for $DOMAIN ==="
echo "Using email: $EMAIL"
echo ""

# Create required directories
mkdir -p certbot/conf certbot/www

# Step 1: Start nginx with init config (HTTP only)
echo "Step 1: Starting nginx for certificate challenge..."
cp nginx/nginx-init.conf nginx/nginx-active.conf

docker-compose -f docker-compose.prod.yml run -d --rm --name nginx-init \
    -v $(pwd)/nginx/nginx-active.conf:/etc/nginx/nginx.conf:ro \
    -v $(pwd)/certbot/www:/var/www/certbot:ro \
    -p 80:80 \
    nginx nginx

sleep 3

# Step 2: Get SSL certificate
echo ""
echo "Step 2: Requesting SSL certificate from Let's Encrypt..."
docker run --rm \
    -v $(pwd)/certbot/conf:/etc/letsencrypt \
    -v $(pwd)/certbot/www:/var/www/certbot \
    certbot/certbot certonly \
    --webroot \
    --webroot-path=/var/www/certbot \
    --email "$EMAIL" \
    --agree-tos \
    --no-eff-email \
    -d "$DOMAIN" \
    -d "www.$DOMAIN"

# Step 3: Stop init nginx
echo ""
echo "Step 3: Stopping temporary nginx..."
docker stop nginx-init || true

# Step 4: Clean up
rm nginx/nginx-active.conf

echo ""
echo "=== SSL Setup Complete! ==="
echo ""
echo "Now start the full stack with:"
echo "  docker-compose -f docker-compose.prod.yml up -d --build"
echo ""
echo "Your site will be available at: https://www.$DOMAIN"
