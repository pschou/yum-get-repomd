#!/bin/bash

# Just a huge key downloader!

keys=(
# CentOS 8 and beyond
https://www.centos.org/keys/RPM-GPG-KEY-CentOS-Official
https://www.centos.org/keys/RPM-GPG-KEY-CentOS-Testing

# current CentOS 7 starter set...
https://www.centos.org/keys/RPM-GPG-KEY-CentOS-7
https://www.centos.org/keys/RPM-GPG-KEY-CentOS-Debug-7
https://www.centos.org/keys/RPM-GPG-KEY-CentOS-Testing-7

# old CentOS 6
https://www.centos.org/keys/RPM-GPG-KEY-CentOS-6
https://www.centos.org/keys/RPM-GPG-KEY-CentOS-Debug-6
https://www.centos.org/keys/RPM-GPG-KEY-CentOS-Testing-6
https://www.centos.org/keys/RPM-GPG-KEY-CentOS-Security-6

# old CentOS 5
https://www.centos.org/keys/RPM-GPG-KEY-CentOS-5

#  Beginning of the SIG (Special Interest Groups?)
https://www.centos.org/keys/RPM-GPG-KEY-CentOS-SIG-Extras 
https://www.centos.org/keys/RPM-GPG-KEY-CentOS-SIG-Atomic
https://www.centos.org/keys/RPM-GPG-KEY-CentOS-SIG-Automotive
https://www.centos.org/keys/RPM-GPG-KEY-CentOS-SIG-Cloud
https://www.centos.org/keys/RPM-GPG-KEY-CentOS-SIG-ConfigManagement
https://www.centos.org/keys/RPM-GPG-KEY-CentOS-SIG-Core
https://www.centos.org/keys/RPM-GPG-KEY-CentOS-SIG-HyperScale
https://www.centos.org/keys/RPM-GPG-KEY-CentOS-SIG-Infra
https://www.centos.org/keys/RPM-GPG-KEY-CentOS-SIG-Kmods
https://www.centos.org/keys/RPM-GPG-KEY-CentOS-SIG-Messaging
https://www.centos.org/keys/RPM-GPG-KEY-CentOS-SIG-NFV
https://www.centos.org/keys/RPM-GPG-KEY-CentOS-SIG-OpsTools
https://www.centos.org/keys/RPM-GPG-KEY-CentOS-SIG-PaaS
https://www.centos.org/keys/RPM-GPG-KEY-CentOS-SIG-SCLo
https://www.centos.org/keys/RPM-GPG-KEY-CentOS-SIG-Storage
https://www.centos.org/keys/RPM-GPG-KEY-CentOS-SIG-Virtualization

# Platform ones. Alt Architectures...
https://www.centos.org/keys/RPM-GPG-KEY-CentOS-7-aarch64
https://www.centos.org/keys/RPM-GPG-KEY-CentOS-SIG-AltArch-Arm32
https://www.centos.org/keys/RPM-GPG-KEY-CentOS-SIG-AltArch-7-ppc64
)

# Download them all to the keys folder!
mkdir -p ${0%/*}/keys; cd ${0%/*}/keys
for key in "${keys[@]}"; do
  curl -o "${key##*/}.gpg" "$key"
done
