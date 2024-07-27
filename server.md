<div align="center" style="text-align: center">
<a href="http://hpcme.com">
<img src="http://hpcme.com/wp-content/uploads/2021/10/cropped-Logo-HPCME-Systems-72x50.jpg" alt="HPCME logo"/>
</a>
<h3>HPCME Systems</h3>
# Server side configurations
</div>


After perfornming the below steps you should be able to boot a basic linux image (jammy) through serving the filesystem using nfs.
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
chown -R tftp:tftp /srv/tftp
chmod -R 755 /srv/tftp
```
modify the nfs exports
```bash
vim /etc/exports
```
```bash
/srv/nfs/jammy *(rw,sync,no_subtree_check,no_root_squash)
```
## grub configurations
```bash
mkdir -p /srv/tftp/grub
```
copy the contents of grub/grub.cfg to your server and copy the required modules
``` bash
vim /srv/tftp/grub/grub.cfg
cp -r /boot/grub/x86_64-efi/ /srv/tftp/grub/
```
## restart the services
```bash
systemctl restart dnsmasq
systemctl restart nfs-kernel-server
```
