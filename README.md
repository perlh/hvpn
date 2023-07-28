# tunvitor
A simple tun VPN and ss tunnel, tun uses udp as the tunnel.

一个简单的ss和vpn程序，vpn程序使用vpn over udp模式运行。
## install
``` bash
git clone https://github.com/perlh/tunvitor.git
cd tunvitor
make
```
## 部署/Deploy
### 客户端/client
``` bash
sudo ./tunvitor -c -socks "0.0.0.0:1081" -remote_server "222.204.52.222:9000" -stun "222.204.52.222:30000" -subnet "192.168.1.111/24"
```
#### 详细信息
| 默认值           | 参数                     | 描述                         |
|---------------|------------------------|----------------------------|
| -c            |                        | 客户端模式                      |
| -socks        | 0.0.0.0:1080           | 创建一个socks5服务器，并且监听本地1080端口 |
| remote_server | 222.204.52.222:9000    | 远程ss服务器地址                  |
| -stun         | "222.204.52.222:30000" | 远程vpn服务器地址                 |
| -subnet       | 192.168.1.111/24       | vpn虚拟网卡网络信息                |
### server
``` bash
sudo ./tunvitor -s -server "0.0.0.0:9000" -ltun "0.0.0.0:30000" -subnet "192.168.123.1/24"
```

#### 详细信息
| 参数      | 参数值              | 描述              |
|---------|------------------|-----------------|
| -s      || 服务模式             |
| -server | 0.0.0.0:9000     | ss服务器监听地址       |
| -ltun   | 0.0.0.0:30000    | vpn服务器UDP隧道监听地址 |
| -subnet | 192.168.123.1/24 | vpn虚拟网卡网络信息     |