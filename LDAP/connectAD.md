<div align="center" style="text-align: center">
<a href="http://hpcme.com">
<img src="http://hpcme.com/wp-content/uploads/2021/10/cropped-Logo-HPCME-Systems-72x50.jpg" alt="HPCME logo"/>
</a>
<h3>HPCME Systems</h3>

</div>

# Connecting to AD
This Document shows how you can connect to AD and use it to authenticate users into your system


## Ubuntu & Base-view on ubuntu
### Install Packages
``` bash
apt update -y
apt install realmd sssd sssd-tools adcli samba-common-bin packagekit krb5-user libpam-sss libnss-sss oddjob oddjob-mkhomedir
```
### DNS config
``` bash
hostnamectl set-hostname <fullname>

```
Edit the /etc/resolv.conf
``` bash 
nameserver <AD IP >
search hpcme.com
```
Make the changes persistent
### Join the AD
``` bash 
ping hpcme.com
realm discover -v hpcme.com
kinit Administrator@HPCME.COM
realm join -v -U <user> hpcme.com
```
### User configs
``` bash
bash -c "cat > /usr/share/pam-configs/mkhomedir" << EOF
Name: activate mkhomedir
Default: yes
Priority: 900
Session-Type: Additional
Session:
 required pam_mkhomedir.so umask=0022 skel=/etc/skel
EOF

pam-auth-update

```

Modify the /etc/sssd/sssd.conf

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
ldap_idmap_range_min = 2000  #To limit the user assigned UID
ldap_idmap_range_max = 1000000
access_provider = simple
simple_allow_groups = IT-admins  # Only allow this group users
```
### Restart the service
``` bash
systemctl restart sssd
systemctl status sssd
```
If you encountered problems with sssd status when it comes to dynamic dns add these lines if not ignore
``` bash
dyndns_refresh_interval = 43200
dyndns_update_ptr = false
dyndns_ttl = 3600
```
### Log in with the username
``` bash
login <username>
```
If you see this error "Unable to create and initialize directory" just try again and it will work the directory is still being initialized.
### Enable SSH login

Modify /etc/ssh/sshd_config and uncomment the following:
```
PasswordAuthentication yes
```

## RHEL

### Install Packages
``` bash
dnf update -y
dnf install -y sssd sssd-tools realmd adcli krb5-workstation oddjob oddjob-mkhomedir
```
### DNS config
edit /etc/hosts
``` bash
< AD IP>  hpcme.com

```
Edit the /etc/resolv.conf
``` bash 
nameserver <AD IP >
search hpcme.com
```
After this you should be able to :
``` bash
ping hpcme.com 
```
## Sync the two servers
``` bash
dnf install -y chrony
systemctl enable --now chronyd
chronyc sources
timedatectl set-ntp true
```
## Modify krb config
vim /etc/krb5.conf
``` bash
[libdefaults]
  default_realm = HPCME.COM
  dns_lookup_realm = false
  dns_lookup_kdc = true

[realms]
  HPCME.COM = {
    kdc = <AD IP>
    admin_server = <AD IP>
  }

[domain_realm]
  .hpcme.com = HPCME.COM
  hpcme.com = HPCME.COM
```
Test Kerberos:
``` bash
kinit Administrator@HPCME.COM
klist
```
## Join the realm

### Important Notes
1. The local machine needs to be able to resolve the AD.
1. The AD and the local machine needs to be syncronized if they are not you need to install ntp and the AD the NTP server.
