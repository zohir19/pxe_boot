#!/bin/#!/usr/bin/env bash
mkdir -p /srv/tftp
cat <<EOF >> /etc/dnsmasq.d/00-header.conf
port=0
dhcp-hostsfile=/etc/dnsmasq.d/01-test.hosts
interface=enp107s0                # Use the appropriate network interface
dhcp-range=192.168.0.100,192.168.0.150,12h
dhcp-boot=grubnetx64.efi.signed,linuxhint-s20,192.168.0.1
enable-tftp
tftp-root=/srv/tftp
EOF

cat <<EOF >> /etc/dnsmasq.d/01-test.hosts
dhcp-host=3c:ec:ef:7b:b4:94,client1,192.168.0.121,3600
dhcp-host=3c:ec:ef:7b:b3:68,client2,192.168.0.122,3600
dhcp-host=3c:ec:ef:7b:b3:48,client3,192.168.0.123,3600
EOF
cp /usr/lib/grub/x86_64-efi-signed/grubnetx64.efi.signed /srv/tftp/
cp /usr/lib/shim/shimx64.efi.signed /srv/tftp
mkdir /srv/nfs
debootstrap jammy /srv/nfs/jammy
echo "/srv/nfs/jammy *(rw,sync,no_subtree_check,no_root_squash)" >> /etc/exports
mkdir -p /srv/tftp/grub
cat <<EOF >> /srv/tftp/grub/grub.cfg
set timeout=5
timeout_style=menu
#debug=all
set net_default_server=192.168.0.1

menuentry 'DB overlay' {
    linux /jammy/vmlinuz root=/dev/nfs nfsroot=192.168.0.1:/srv/nfs/db_overlay rw BOOTIF=01-$net_default_mac BOOTIP=$net_default_ip console=tty0 console=ttyS0,115200 earlyprintk=ttyS0,115200
    initrd /jammy/initrd.img
}

menuentry 'Ubuntu 22.04' {
    linux /jammy/vmlinuz root=/dev/nfs nfsroot=192.168.0.1:/srv/nfs/jammy rw BOOTIF=01-$net_default_mac BOOTIP=$net_default_ip console=tty0 console=ttyS1,115200 earlyprintk=ttyS1,115200
    initrd /jammy/initrd.img
}
EOF
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
