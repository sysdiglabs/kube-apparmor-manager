#!/bin/bash
set -eu

SSH_KEY="~/.ssh/id_rsa"
SSH_USER="admin"
WORK_DIR="/var/lib/aa"
AA_PROFILE_DIR="/etc/apparmor.d"
AA_PROFILES="$PWD/profiles/*"

# get worker nodes IPs
IPs=($(kubectl get nodes -l kubernetes.io/role=node -o json | jq '.items[] | .status.addresses | .[] | select(.type == "ExternalIP") | .address' -r | xargs))

for NODE_IP in "${IPs[@]}"
do
	echo "***Start copying AppArmor profile to IP: $NODE_IP***"
	for PROFILE in $AA_PROFILES
	do
		scp -i $SSH_KEY $PROFILE $"$SSH_USER@$NODE_IP:$WORK_DIR/profiles"
	done
	ssh -i $SSH_KEY "$SSH_USER@$NODE_IP" "sudo cp $WORK_DIR/profiles/* $AA_PROFILE_DIR"
	echo "***End copying AppArmor profile to IP: $NODE_IP***"
done
