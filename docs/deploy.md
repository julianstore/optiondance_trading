### 安装golang 1.16.x+, mysql8.0

### 时间同步
```shell
# centos8
dnf install chrony -y
systemctl enable chronyd.service
timedatectl set-timezone Asia/Shanghai
timedatectl set-ntp true
```

### 创建用户
```shell
useradd optiondance
passwd optiondance
```
加入sudoer列表

### 克隆仓库
   
```shell
su optiondance
cd /home/optiondance
git clone https://github.com/24h-purewater/option-dance
```



### 初始化配置文件

```shell
cat << EOF > /home/optiondance/node_config.yaml
mysql:
  username: ''
  password: ''
  path: ''
  db-name: ''
  config: 'charset=utf8mb4&parseTime=True&loc=Asia%2fShanghai'
  max-idle-conns: 100
  max-open-conns: 1000
  log-mode: true
  

zap:
  level: 'info'
  filename: '/home/optiondance/log/option_dance_engine.log'
  max-size: 200
  max-age: 7
  max-backups: 10

dapp:
  receivers:
    - '1b9195c8-e03b-448c-a931-69526f879332'
    - '98dd6252-1ae4-4d36-a8a0-8b4f9797632c'
    - ''
  threshold: 2
  app_id: ''
  session_id: ''
  secret: ''
  pin: ''
  pin_token: ''
  private_key: ''
EOF
```

### 构建
```shell
cd /home/optiondance/option-dance/cmd/engine
go build -o od-engine
```

### 初始化数据库
```shell
./od-engine migrate -c /home/optiondance/node_config.yaml
```


### 初始化engine systemd配置文件
```shell
cat << EOF > /etc/systemd/system/od-engine.service
[Unit]
Description=Option Dance Engine Daemon
After=network.target

[Service]
User=optiondance
Type=simple
ExecStart= /bin/bash -c "/home/optiondance/option-dance/cmd/engine/od-engine -c /home/optiondance/node_config.yaml  -p 6100"
Restart=on-failure
LimitNOFILE=65536

[Install]
WantedBy=multi-user.target
EOF
```

### 启动 engine
```shell
sudo systemctl start od-engine
```



### （可选）构建api
```shell
cd /home/optiondance/option-dance/cmd/api
go build -o od-api
```

```shell
cat << EOF > /etc/systemd/system/od-api.service
[Unit]
Description=Option Dance Api Server
After=network.target

[Service]
User=optiondance
Type=simple
ExecStart=/bin/bash -c "/home/optiondance/option-dance/cmd/api/od-api -c /home/optiondance/node_api.yaml"
Restart=on-failure
LimitNOFILE=65536

[Install]
WantedBy=multi-user.target
EOF
```