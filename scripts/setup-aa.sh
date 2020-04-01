#!/bin/bash

GRUB="/etc/default/grub"
GRUB_BAK="/etc/default/grub.bak"
ENABLED=$(aa-enabled 2>/dev/null)

if [[ $? -ne 0 || "$ENABLED" == No* ]]; then
	echo "AppArmor is disabled"
	echo "Install AppArmor"
	apt install -y apparmor-profiles apparmor-utils

	echo "Update grub"
	cp $GRUB $GRUB_BAK
	sed -i -e '/^GRUB_CMDLINE_LINUX_DEFAULT/s/"$/ apparmor=1 security=apparmor"/' $GRUB
	cat $GRUB | grep "^GRUB_CMDLINE_LINUX_DEFAULT"
	update-grub

	echo "Reboot system"
	reboot
	exit 0
else
	echo "AppArmor is enabled"
	exit 0
fi
