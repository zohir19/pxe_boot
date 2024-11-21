# Connecting to AD
This Document shows how you can connect to AD and use it to authenticate users into your system


## Ubuntu
### Install Packages
``` bash
apt update -y
apt install realmd sssd sssd-tools adcli samba-common-bin packagekit krb5-user libpam-sss libnss-sss oddjob oddjob-mkhomedir
```
### DNS config
``` bash
hostnamectl set-hostname fullname 
systemctl disable systemd-resolved.service
systemctl stop systemd-resolved.service
```
Edit the /etc/resolv.conf
``` bash 
nameserver <AD IP >
search <AD Domain >
```
### Join the AD
``` bash 
ping <realm>
realm discover -v <realm>
kinit Administrator@<realm>
realm join -v -U <user> <realm>
```
### User configs
``` bash
bash -c "cat > /usr/share/pam-configs/mkhomedir" << EOF
Name: activate mkhomedir
Default: yes
Priority: 900
Session-Type: Additional
Session:
 required pam_mkhomedir.so umask=0022 skel=/et/skel
EOF

pam-auth-update

vim /etc/sssd/sssd.conf
```