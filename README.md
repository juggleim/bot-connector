# 说明

这是一个为 JuggleIM 扩展机器人能力的桥接项目，当前代码示例，展示了如何与 TG Bot 进行对接

# 使用前提

1. 部署 JuggleIM 服务，参考：https://github.com/juggleim/im-server/blob/master/README.md
2. 部署 JuggleChat-Server 服务，参考：https://github.com/juggleim/jugglechat-server/blob/master/README.md

# 部署说明

### 1. 安装并初始化 MySQL

#### 1）安装 MySQL
略

#### 2）创建DB实例

```
CREATE SCHEMA `jbot_db`;
```

#### 3）初始化表结构
初始化表结构的SQL文件在 bot-connector/docs/bot.sql，导入命令如下：
```
mysql -u{db_user} -p{db_password} jbot_db < bot.sql
```

### 2. 启动 bot-connector

#### 1）运行目录
运行目录即根目录 bot-connector， 其中 conf 目录下存放配置文件。

#### 2）编辑配置文件
配置文件位置：bot-connector/conf/config.yml
```
port: 8070

log:
  logPath: ./logs
  logName: bot-connector

mysql:
  user: root
  password: 123456
  address: 127.0.0.1:3306
  name: jbot_db

imApiDomain: https://api.juggle.im   # JuggleIM 服务对应的服务端api访问地址

botConnector:
  domain: http://127.0.0.1:8070      # bot-connector 对外访问地址，IM 服务会使用这个接口来访问 bot-connector
  apiKeySecret: <secret string(length=16)>   # bot-connector 对外接口的鉴权 secret，注意使用 16 位的密钥串
```

#### 3）启动 bot-connector
在根目录，执行如下命令：
```
go run main.go
```

