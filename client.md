<div align="center" style="text-align: center">
<a href="http://hpcme.com">
<img src="http://hpcme.com/wp-content/uploads/2021/10/cropped-Logo-HPCME-Systems-72x50.jpg" alt="HPCME logo"/>
</a>
<h3>HPCME Systems</h3>

</div>

# Client side configurations
Now it's time to customize the image
## Accessing the image
Mount the filesystem
```bash
mount --bind /dev /srv/nfs/jammy/dev
mount --bind /proc/ /srv/nfs/jammy/proc/
mount --bind /sys /srv/nfs/jammy/sys
```
Access the image
```bash
chroot /srv/nfs/jammy/
```
# Install the required packages
Update and install the following packages
```bash
cp /etc/apt/sources.list /srv/nfs/jammy/etc/apt/sources.list
apt update
apt install linux-image-generic
apt install vim
apt install parted
apt install dosfstools
apt install rsync
apt install nfs-common
apt install grub-pc-lib
apt install grub-pc-bin
```
# Create the root user password or another management user
```bash
chroot /srv/nfs/jammy
passwd
```
