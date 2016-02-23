ORIGINAL_GW     ?= $$(ip route  | grep default | awk '{print $$3}')

DEVICE_NAME     ?= tun0
DEVICE_IP       ?= 10.0.0.1
DEVICE_GW_IP    ?= 10.0.0.2

DNS1            ?= 8.8.8.8
DNS2            ?= 8.8.4.4

PROXY_IP        ?= 10.0.0.66

up:
	sudo ip tuntap del $(DEVICE_NAME) mode tun && \
	sudo ip tuntap add $(DEVICE_NAME) mode tun &&& \
	sudo ifconfig $(DEVICE_NAME) $(DEVICE_IP) netmask 255.255.255.0 && \
	sudo route add $(PROXY_IP) gw $(ORIGINAL_GW) metric 5 && \
	sudo route add default gw $(DEVICE_GW_IP) metric 6 && \
	print "nameserver $(DNS1)\r\nnameserver $(DNS2)\r\n" | sudo tee /etc/resolv.conf
