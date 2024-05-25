# pxe_boot
 booting ubuntu 22.04 from pxe

 ## Introduction
 ## DHCP
 Install dhcp and modify it's configuration file as follows:
 ``` bash
apt-get install isc-dhcp-server
vim /etc/dhcp/dhcpd.conf
```
```bash
allow booting;
allow bootp;

subnet 192.168.56.0 netmask 255.255.255.0 {
        range 192.168.56.122 192.168.56.125;
        option domain-name "example.com";
        option domain-name-servers 8.8.8.8, 8.8.4.4;
        option broadcast-address 192.168.56.255;
        option routers 192.168.56.1;
        next-server 192.168.56.121;
        option subnet-mask 255.255.255.0;

        filename "/pxelinux.0";
}

# force the client to this ip for pxe.
# This isn't strictly necessary but forces each computer to always have the same IP address
host node21 {
        hardware ethernet 01:23:45:a8:50:26;
        fixed-address 192.168.56.122;
}
```
 ## TFTP
Install the tftpd and modify it's config file
 ```bash
apt install tftpd-hpa syslinux-common
vim /etc/default/tftpd-hpa
```

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
#### if the pxelinux.0 deosn't exist
download the syslinux package and then copy the pxelinux.0 to the tftp directory
``` bash
wget https://mirrors.edge.kernel.org/pub/linux/utils/boot/syslinux/6.xx/syslinux-6.03.tar.xz
tar xvf syslinux-6.03.tar.xz
cp syslinux-6.03/bios/core/pxelinux.0 /srv/tftp/
cp /usr/lib/syslinux/modules/bios/ldlinux.c32 /srv/tftp/
```
Create the default file and modify it
```bash
mkdir /srv/tftp/pxelinux.cfg
vim /srv/tftp/pxelinux.cfg/default
```

```
DEFAULT linux
LABEL linux
KERNEL vmlinuz-5.15.0-107-generic
APPEND root=/dev/nfs initrd=initrd.img-5.15.0-107-generic nfsroot=192.168.56.121:/clusternfs,ro ip=dhcp ro
IPAPPEND 2
```
## Creating the worker node filesystem
```bash
apt install debootstrap
mkdir /clusternfs
debootstrap jammy /clusternfs/
cp -a /lib/modules /clusternfs/lib/
```
 ## NFS
 Install the required packages
 ```bash
apt install nfs-kernel-server
```
Modify the /etc/exports file 
```
/clusternfs 192.168.56.0/24(rw,sync,no_root_squash,no_subtree_check)
```
Edit the /clusternfs/etc/fstab file
```
proc            /proc         proc   defaults       0      0
/dev/nfs        /             nfs    defaults       0      0
```
