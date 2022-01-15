# option-dance

## 构建

####  构建 engine
```shell
cd cmd/engine
go build -o od-engine
```

####  构建 api server
```shell
cd cmd/api
go build -o od-api
```

## 配置

eg. config/node_sample.yaml

## 运行与部署

### 多签与撮合引擎 od-engine
```
./od-engine -h

Available Commands:
  help        Help about any command  
  migrate     migrate option dance database table struct  # 迁移数据库表结构

Flags:
  -c, --config string   config file (default is $HOME/config.yaml) 
  -h, --help            help for od-engine
      --notify          enable message notify (default false) # 是否需要推送通知，只需主体机器人开启，其它参与多签的机器人无需开启
  -p, --port string     server bind port (default "6028")  # http server 端口

# eg.
> ./od-engine  -c config.yaml -p 6028
```


### API

```
./od-api -h 

Usage:
  od-api [flags]

Flags:
  -c, --config string   config file (default is $HOME/config.yaml)
  -h, --help            help for od-api
  -p, --port string     api server bind port (default "6024")

# eg.
./od-api -c config.yaml -p 6024
```

### 其它

[OptionDance介绍](docs/introduction.md)