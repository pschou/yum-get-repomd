#!/bin/bash

# Just a huge key downloader!

keys=(
# Rocky
https://dl.rockylinux.org/pub/rocky/RPM-GPG-KEY-rockyofficial
https://dl.rockylinux.org/pub/rocky/RPM-GPG-KEY-rockytesting
https://dl.rockylinux.org/pub/rocky/RPM-GPG-KEY-rockyinfra
)

# Download them all to the keys folder!
mkdir -p ${0%/*}/keys; cd ${0%/*}/keys
for key in "${keys[@]}"; do
  curl -o "${key##*/}.gpg" "$key"
done
