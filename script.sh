#!/bin/#!/usr/bin/env bash
mkdir -p /srv/tftp
vim /etc/dnsmasq.d/00-header.conf
vim /etc/dnsmasq.d/01-test.hosts
cp /usr/lib/grub/x86_64-efi-signed/grubnetx64.efi.signed /srv/tftp/
cp /usr/lib/shim/shimx64.efi.signed /srv/tftp
mkdir /srv/nfs
debootstrap jammy /srv/nfs/jammy
echo "/srv/nfs/jammy *(rw,sync,no_subtree_check,no_root_squash)" >> /etc/exports
mkdir -p /srv/tftp/grub
vim /srv/tftp/grub/grub.cfg
cp -r /boot/grub/x86_64-efi/ /srv/tftp/grub/
systemctl restart dnsmasq
systemctl restart nfs-kernel-server
mount --bind /dev /srv/nfs/jammy/dev
mount --bind /proc/ /srv/nfs/jammy/proc/
mount --bind /sys /srv/nfs/jammy/sys
cp /etc/apt/sources.list /srv/nfs/jammy/etc/apt/sources.list
cp initial_setup.sh /srv/nfs/jammy/usr/local/bin
cp initial_setup.service /srv/nfs/jammy/etc/systemd/system
chroot /srv/nfs/jammy/
apt update
apt install linux-image-generic
apt install vim
apt install parted
apt install dosfstools
apt install rsync
apt install nfs-common
apt install grub-pc-lib
apt install grub-pc-bin
passwd
systemctl enable initial_setup.service
exit
cp /srv/nfs/jammy/boot/vmlinuz /srv/tftp/jammy/vmlinuz
cp /srv/nfs/jammy/boot/initrd.img /srv/tftp/jammy/initrd.img
