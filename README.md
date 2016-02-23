# tunio

The `tunio` package captures and forwards raw TCP packets to a Go's
[net.Dialer](https://golang.org/pkg/net/#Dialer) using a [TUN device][3].

`tunio` is able to forward UDP packets as well, an external [badvpn-udpgw][4]
server is required.

## How to compile and run?

Clone the source

```
mkdir -p $GOPATH/src/github.com/getlantern
cd $GOPATH/src/github.com/getlantern
git clone https://github.com/getlantern/tunio.git
cd tunio
```

Build and run the docker container that provides the `badvpn-udpgw` +
[Lantern][2] bundle:

```
make -C scripts/server docker-run
# ...
docker ps
# ...
# a6fca81fffee        getlantern/tunio-proxy   "/usr/bin/start.sh"   16 minutes ago      Up 16 minutes       0.0.0.0:2099->2099/tcp, 0.0.0.0:5353->5353/tcp   tunio-proxy
# ...
```

Using a docker container is fine for creating and running a proxy, but in order
to work with tunio we need more room to experiment, let's create a full virtual
machine with centos-7.2:

```
mkdir -p ~/projects/tunio-vm
cd ~/projects/tunio-vm
vagrant init bento/centos-7.2
vagrant up
# ...
```

Log in into the newly created virtual machine and install some required
packages:

```sh
vagrant ssh
sudo yum install -y git gcc glibc-static

# Go
curl --silent https://storage.googleapis.com/golang/go1.6.linux-amd64.tar.gz \
	| sudo tar -xzv -C /usr/local/

echo 'export GOROOT=/usr/local/go'    >> $HOME/.bashrc
echo 'export PATH=$PATH:$GOROOT/bin'  >> $HOME/.bashrc
echo 'export GOPATH=$HOME/go'         >> $HOME/.bashrc

source $HOME/.bashrc
```

Clone the `tunio` package (again, into the vm) with `git` and change directory
to the `tunio`'s source path:

```sh
mkdir -p $GOPATH/src/github.com/getlantern
cd $GOPATH/src/github.com/getlantern
git clone https://github.com/getlantern/tunio.git
cd tunio
```

Compile the `tunio` binary with:

```sh
make binary
```

After this you should have a `tunio` binary on `$GOPATH/bin`.

Create a new TUN device, let's name it `tun0` and assign the `10.0.0.1` IP
address to it.

```sh
export ORIGINAL_GW=$(ip route  | grep default | awk '{print $3}')

export DEVICE_NAME=tun0
export DEVICE_IP=10.0.0.1

sudo ip tuntap del $DEVICE_NAME mode tun
sudo ip tuntap add $DEVICE_NAME mode tun
sudo ifconfig $DEVICE_NAME $DEVICE_IP netmask 255.255.255.0
```

Replace the virtual machine's name servers with 8.8.8.8 and 8.8.4.4.

```
echo "nameserver 8.8.8.8" | sudo tee /etc/resolv.conf
echo "nameserver 8.8.4.4" | sudo tee -a /etc/resolv.conf
```

Modify the routing table to only allow direct traffic with the proxy server
(the docker host's IP)

```sh
export PROXY_IP=10.0.0.66

sudo route add $PROXY_IP gw $ORIGINAL_GW metric 5
sudo route add default gw 10.0.0.2 metric 6
```

That was a lot of stuff, you can also do all that by using this script:

```
PROXY_IP=10.0.0.66 make -C scripts/tunconfig up
```

If everything is OK you should not be able to ping external hosts:

```
ping google.com
PING google.com (74.125.227.165) 56(84) bytes of data.
^C
--- google.com ping statistics ---
5 packets transmitted, 0 received, 100% packet loss, time 4001ms
```

But you should be able to ping `$PROXY_IP`.

```
ping $PROXY_IP
PING 10.0.0.66 (10.0.0.66) 56(84) bytes of data.
64 bytes from 10.0.0.66: icmp_seq=1 ttl=63 time=3.45 ms
64 bytes from 10.0.0.66: icmp_seq=2 ttl=63 time=0.403 ms
^C
--- 10.0.0.66 ping statistics ---
2 packets transmitted, 2 received, 0% packet loss, time 1001ms
rtt min/avg/max/mdev = 0.403/1.928/3.454/1.526 ms
```

Finally, run `tunio` with the `--proxy-addr` parameter pointing to Lantern and
with `--udpgw-remote-server-addr` pointing to `127.0.0.1:5353` (which is the
address of the udpgw server as the docker container sees it).

```sh
./tunio --tundev tun0 \
  --netif-ipaddr 10.0.0.2 \
  --netif-netmask 255.255.255.0 \
  --proxy-addr $PROXY_IP:2099 \
  --udpgw-remote-server-addr 127.0.0.1:5353
```

You should be able to browse any site now, the request will be captured by tun0
and it will be forwarded to tunio which is connected to Lantern using a
`net.Dialer`.

```
curl google.com
# <HTML><HEAD><meta http-equiv="content-type" content="text/html;charset=utf-8">
# <TITLE>302 Moved</TITLE></HEAD><BODY>
# <H1>302 Moved</H1>
# The document has moved
# <A HREF="http://www.google.com.mx/?gfe_rd=cr&amp;ei=lbvLVtb0FM_E8ge_zouwDg">here</A>.
# </BODY></HTML>
```

Hurray!

[1]: https://github.com/ambrop72/badvpn/tree/master/tun2socks
[2]: https://getlantern.org
[3]: https://en.wikipedia.org/wiki/TUN/TAP
[4]: https://felixc.at/BadVPN)
