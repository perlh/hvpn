# tunvitor
## install
``` bash
git clone xxx.git
cd tunvitor
go build
```
## deploy
### client
``` bash
sudo ./tunvitor -c \ client mode
    -socks "0.0.0.0:1080" \ listen client 1080 port
    -remote_server "222.204.52.222:9000" \ ss remote server address
    -stun "222.204.52.222:30000"  \ stun remote server address
    -subnet "192.168.1.111/24"    \ tun local vpn address                                              [14:11:34]
```

### server
``` bash
sudo ./tunvitor -s \ server mode
    -server "0.0.0.0:9000" \ ss server 9000 port
    -ltun "0.0.0.0:30000" \ tun listen port
    -subnet "192.168.123.1/24" \ tun local vpn address
```
