# pxe_boot
 booting ubuntu 22.04 from pxe with dnsmasq
## Install the required packages
you can find all the required packages under requirements.yml
## create the tftp directory and modify the dnsmasq configurations
```bash
mkdir -p /srv/tftp
# copy the contents of dnsmasq.d/00-header.conf to your server
vim /etc/dnsmasq.d/00-header.conf
# copy the contents of dnsmasq.d/01-test.hosts to your server
vim /etc/dnsmasq.d/01-test.hosts
```
