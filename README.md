# pxe_boot
 booting ubuntu 22.04 from pxe

 ## Introduction
 ## Installation
 ```bash
apt install tftpd-hpa nfs-kernel-server
apt install syslinux-common
apt install debootstrap
apt-get install isc-dhcp-server
```
 ## configuration
modify the /etc/default/tftpd-hpa as follows:
```
# /etc/default/tftpd-hpa

TFTP_USERNAME="tftp"
TFTP_DIRECTORY="/srv/tftp"
TFTP_ADDRESS="0.0.0.0:69"
TFTP_OPTIONS="--secure --create --listen"
```
Create the TFTP_DIRECTORY 
``` bash
mkdir -p /srv/tftp
cd /srv/tftp
```
Copy the configuration files to the tftp directory .
``` bash
cp /boot/vmlinuz-$(uname -r) /srv/tftp/
cp /boot/initrd.img-$(uname -r) /srv/tftp/
cp /usr/lib/syslinux/modules/bios/pxelinux.0 .
```
### if the pxelinux.0 deosn't exist
download the syslinux package and then copy the pxelinux.0 to the tftp directory
``` bash
wget https://mirrors.edge.kernel.org/pub/linux/utils/boot/syslinux/6.xx/syslinux-6.03.tar.xz
tar xvf syslinux-6.03.tar.xz
cp syslinux-6.03/bios/core/pxelinux.0 /srv/tftp/
cp /usr/lib/syslinux/modules/bios/ldlinux.c32 /srv/tftp/
```