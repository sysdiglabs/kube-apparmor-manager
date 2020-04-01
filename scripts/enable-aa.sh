#!/bin/bash
set -eu

SSH_KEY="~/.ssh/id_rsa"
SSH_USER="admin"
WORK_DIR="/var/lib/aa"
AA_SCRIPT="setup-aa.sh"

# get worker nodes IPs
IPs=($(kubectl get nodes -l kubernetes.io/role=node -o json | jq '.items[] | .status.addresses | .[] | select(.type == "ExternalIP") | .address' -r | xargs))

for NODE_IP in "${IPs[@]}"
do
	echo "***Start enabling AppArmor at IP: $NODE_IP***"
	ssh -i $SSH_KEY "$SSH_USER@$NODE_IP" "sudo mkdir -p $WORK_DIR/profiles && sudo chown -R $SSH_USER:$SSH_USER $WORK_DIR"
	scp -i $SSH_KEY $AA_SCRIPT "$SSH_USER@$NODE_IP:$WORK_DIR"
	ssh -o ConnectTimeout=10 -i $SSH_KEY "$SSH_USER@$NODE_IP" "sudo bash $WORK_DIR/$AA_SCRIPT"
	echo "***End configuring AppArmor at IP: $NODE_IP***"
done
