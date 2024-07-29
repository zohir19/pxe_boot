#!/bin/sh

# Identify the smallest disk
DISK=$(lsblk -b -o NAME,SIZE,TYPE | grep "disk" | sort -k 2 -n | head -n 1 | awk '{print $1}')
PART_PREFIX="/dev/${DISK}"

# Partition the disk
echo -n "Creating partitions..."
parted -s $PART_PREFIX mklabel gpt
parted -s -a opt $PART_PREFIX mkpart ESP fat32 1MiB 512MiB
parted -s $PART_PREFIX set 1 boot on
parted -s -a opt $PART_PREFIX mkpart primary ext4 512MiB 100%
echo "Done."

echo -n "Formatting..."
# Format partitions
if [[ $DISK == nvme* ]]; then
  EFI_PART="${PART_PREFIX}p1"
  ROOT_PART="${PART_PREFIX}p2"
else
  EFI_PART="${PART_PREFIX}1"
  ROOT_PART="${PART_PREFIX}2"
fi

mkfs.vfat -F32 $EFI_PART
mkfs.ext4 $ROOT_PART
echo "Done."

echo -n "Mounting filesystems..."
# Mount NFS root and new partitions
mkdir -p /mnt/nfs-root
mount -t nfs 192.168.0.1:/srv/nfs/jammy /mnt/nfs-root # Change the IP to match your server
mount $ROOT_PART /mnt/new-root-partition
mkdir -p /mnt/new-root-partition/boot/efi
mount $EFI_PART /mnt/new-root-partition/boot/efi
echo "Done."

echo -n "copying fs..."
# Copy NFS root to new partition
rsync -a --info=progress2 /mnt/nfs-root/ /mnt/new-root-partition/
echo "done"

echo -n "updating fstab..."


# Update fstab on new root partition
cat <<EOF > /mnt/new-root-partition/etc/fstab
/dev/sda2 / ext4 defaults 0 1
/dev/sda1 /boot/efi vfat defaults umask=0077 0 1
EOF
echo "done"



# Prepare for chroot
mount --bind /dev /mnt/new-root-partition/dev
mount --bind /proc /mnt/new-root-partition/proc
mount --bind /sys /mnt/new-root-partition/sys
mount -t efivarfs none /mnt/new-root-partition/sys/firmware/efi/efivars


echo -n "regenerating initramfs..."
# Regenerate initramfs to include necessary drivers
chroot /mnt/new-root-partition /bin/bash -c "echo sd_mod >> /etc/initramfs-tools/modules && echo ext4 >> /etc/initramfs-tools/modules && update-initramfs -u"
echo "done"


echo -n "installing grub..."
chroot /mnt/new-root-partition /bin/bash -c "grub-install --target=x86_64-efi --efi-directory=/boot/efi --bootloader-id=GRUB "
echo "done"
echo -n "configuring grub..."
chroot /mnt/new-root-partition /bin/bash -c "grub-mkdevicemap && grub-mkconfig -o /boot/grub/grub.cfg"
echo "done"

chroot /mnt/new-root-partition /bin/bash -c "systemctl disable install.service"
grub-mkconfig -o /boot/grub/grub.cfg
grub-mkconfig -o /boot/efi/EFI/GRUB/grub.cfg
echo -n "unmounting..."
# Clean up
#umount /mnt/new-root-partition/dev
#umount /mnt/new-root-partition/proc
#umount /mnt/new-root-partition/sys
#umount /mnt/nfs-root
#umount /mnt/new-root-partition/boot/efi
#umount /mnt/new-root-partition
echo "done"

reboot
