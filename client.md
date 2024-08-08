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
cp /etc/apt/sources.list /srv/nfs/jammy/etc/apt/sources.list
chroot /srv/nfs/jammy/
```
# Install the required packages
Update and install the following packages
```bash
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
copy the kernel and the initrd to the tftp directory
```bash
cp /srv/nfs/jammy/boot/vmlinuz /srv/tftp/jammy/vmlinuz
cp /srv/nfs/jammy/boot/initrd.img /srv/tftp/jammy/initrd.img
#chown -R tftp:tftp /srv/tftp
#chmod -R 755 /srv/tftp
```
# Modify the image and set the initial layout
Create the password for the root user or create another management user:
```bash
chroot /srv/nfs/jammy
passwd
```
Copy the contents of the initial_setup.sh and initial_setup.service
```bash
exit #to exit from the chroot
cp initial_setup.sh /srv/nfs/jammy/usr/local/bin
cp initial_setup.service /srv/nfs/jammy/etc/systemd/system
systemctl enable initial_setup.service
```
# Start your machine and boot from PXE
