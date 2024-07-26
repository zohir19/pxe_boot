# pxe_boot
 booting ubuntu 22.04 from pxe with dnsmasq
## Install the required packages
you can find all the required packages under requirements.yml
## create the tftp directory and modify the dnsmasq configurations
```bash
mkdir -p /srv/tftp
```
copy the contents of dnsmasq.d/00-header.conf to your server
```bash
vim /etc/dnsmasq.d/00-header.conf
```
copy the contents of dnsmasq.d/01-test.hosts to your server
```bash
vim /etc/dnsmasq.d/01-test.hosts
```
copy the required files to tftp directory
```bash
cp /usr/lib/grub/x86_64-efi-signed/grubnetx64.efi.signed /srv/tftp/
cp /usr/lib/shim/shimx64.efi.signed /srv/tftp/
```
## create the worker node filesystem
create the filesystem with debootstrap
```bash
debootstrap jammy /srv/nfs/jammy
```
copy the kernel and the initrd to the tftp directory
```bash
cp /srv/nfs/jammy/boot/vmlinuz /srv/tftp/jammy/vmlinuz
cp /srv/nfs/jammy/boot/initrd.img /srv/tftp/jammy/initrd.img
```
modify the nfs exports
```bash
vim /etc/exports
```
```bash

## grub configurations
```bash
mkdir -p /srv/tftp/grub
