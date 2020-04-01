#!/bin/bash
set -eu

PROFILE=${1:-""}

if [ -z $PROFILE ]; then
	echo "profile cannot be emtpy"
	exit 1
fi

SSH_KEY="~/.ssh/id_rsa"
SSH_USER="admin"
AA_PROFILE_DIR="/etc/apparmor.d"

# get worker nodes IPs
IPs=($(kubectl get nodes -l kubernetes.io/role=node -o json | jq '.items[] | .status.addresses | .[] | select(.type == "ExternalIP") | .address' -r | xargs))

for NODE_IP in "${IPs[@]}"
do
	echo "***Start enabling AppArmor at IP: $NODE_IP***"
	ssh -o ConnectTimeout=10 -i $SSH_KEY "$SSH_USER@$NODE_IP" "sudo aa-enforce /etc/apparmor.d/$PROFILE && sudo systemctl reload apparmor.service"
	echo "***End configuring AppArmor at IP: $NODE_IP***"
done
