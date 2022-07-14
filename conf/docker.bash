#!/usr/bin/env bash

# setup sysctl max udp
echo "net.core.rmem_max=2500000" >> /etc/sysctl.conf
sysctl -p

# run
docker run -d --restart=always --name="lan-expose-proxy" \
  -v /volume1/docker/lan-expose-proxy/proxy.ini:/config/proxy.ini \
  -v /volume1/homes/syj/syno-acme/acme.sh/ **#YOUR DOMAIN#** /fullchain.cer:/config/ssl.crt \
  -v /volume1/homes/syj/syno-acme/acme.sh/ **#YOUR DOMAIN#** / **#YOUR DOMAIN#** .key:/config/ssl.key \
  -p 690:690/udp \
  ghcr.io/shiyunjin/lan-expose-proxy:test \
  -c /config/proxy.ini

