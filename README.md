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
copy the pxelinux.0 to the tftp directory if not found you can download the syslinux package.
``` bash
cp /usr/lib/syslinux/modules/bios/pxelinux.0 .
```