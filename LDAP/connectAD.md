<div align="center" style="text-align: center">
<a href="http://hpcme.com">
<img src="http://hpcme.com/wp-content/uploads/2021/10/cropped-Logo-HPCME-Systems-72x50.jpg" alt="HPCME logo"/>
</a>
<h3>HPCME Systems</h3>

</div>

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

```

modify the /etc/sssd/sssd.conf

``` bash
[sssd]
domains = hpcme.com
config_file_version = 2
services = nss, pam

[domain/hpcme.com]
default_shell = /bin/bash
krb5_store_password_if_offline = True
cache_credentials = True
krb5_realm = HPCME.COM
realmd_tags = manages-system joined-with-adcli
id_provider = ad
fallback_homedir = /home/%u
ad_domain = hpcme.com
use_fully_qualified_names = False
ldap_id_mapping = True
access_provider = simple
simple_allow_groups = IT-admins  # Only allow this group users
```