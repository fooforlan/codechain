# codechain

​	codechain基于 Tendermint共识机制的代码链，采用go语言实现，并且依赖于Tendermint环境，建议在ubuntu16.04下开发，首先 [建立开发环境](https://github.com/little51/codechain/blob/master/getting-start.md) ，为了验证Tendermint的四个节点的运行情况，可以参考 [运行环境组网](https://github.com/little51/codechain/blob/master/making-testnet.md) 。

## Codechain核心

### 安装Mongodb

​	资产数据采用Mongodb保存，所以请先安装Mongodb。

```shell
sudo apt-get install mongodb
```

### 编译

```shell
export GO111MODULE=on
export GOPROXY=https://goproxy.io
cd core
go build
```

### 运行

#### 链初始化

```shell
tendermint init
```

链初始化只做一次，执行完成后，会在~/.tendermint目录下生成配置文件和数据文件。

#### 运行应用

```shell
./core
```

core程序会在26658端口监听tendermint进程发过来的交易。

#### 运行NODE

​	再开一个命令行，执行：

```shell
tendermint node
```

tendermint会在26657端口监听。

### 测试

```shell
curl -s 'localhost:26657/broadcast_tx_commit?tx="key1=value1"'
curl -s 'localhost:26657/abci_query?data="key=key1"'
```

## Web应用

​	Web应用目前提供了账户管理和简单的资产管理，可实现账户密钥对的创建，资产的签名、登记。

### 编译运行

```shell
cd webserver
go build
./webserver
```

### 测试

​	新建账户，返回公私钥对，对资产字符串进行私钥签名，用公钥验证。

#### 新建账户

```shell
curl -X POST http://localhost:3000/account/new
```

结果如下：

```json
{"address":"89FC3A4172D79264D3ECA93DDF988D678EB2EC08","error":"","privateKey":"a53ea605169339e66891464a3de05c415e05814087fddb6fa41074044a43b8fa61e21af00e674610ddc63c3975a118c39639dea848df0baacdaf61f604f5d9a5","publicKey":"61E21AF00E674610DDC63C3975A118C39639DEA848DF0BAACDAF61F604F5D9A5"}
```

#### 资产签名

```shell
curl  -H "Content-Type: application/json" -d '{"privatekey":"a53ea605169339e66891464a3de05c415e05814087fddb6fa41074044a43b8fa61e21af00e674610ddc63c3975a118c39639dea848df0baacdaf61f604f5d9a5","msg":"{\"key\":\"myasset2\",\"value\":\"string111\"}"}' -X POST http://localhost:3000/account/sign
```

结果如下：

```json
{"error":"","sign":"f3b2ff6c01455afc350c67037135d60a61ff049d81ec4a111d54fac07578fb0af1dc049205af16c6cf5ec77575dfa37a5d0e89991eca167c38a5d2c9c50b3308"}
```

#### 资产登记

```shell
curl  -H "Content-Type: application/json" -d '{"publickey":"61E21AF00E674610DDC63C3975A118C39639DEA848DF0BAACDAF61F604F5D9A5","sign":"f3b2ff6c01455afc350c67037135d60a61ff049d81ec4a111d54fac07578fb0af1dc049205af16c6cf5ec77575dfa37a5d0e89991eca167c38a5d2c9c50b3308","msg":"{\"key\":\"myasset2\",\"value\":\"string111\"}"}' -X POST http://localhost:3000/assets/new
```

结果如下：

```json
{"error":"","info":"{\n  \"jsonrpc\": \"2.0\",\n  \"id\": -1,\n  \"result\": {\n    \"check_tx\": {\n      \"code\": 0,\n      \"data\": null,\n      \"log\": \"\",\n      \"info\": \"\",\n      \"gasWanted\": \"1\",\n      \"gasUsed\": \"0\",\n      \"events\": [],\n      \"codespace\": \"\"\n    },\n    \"deliver_tx\": {\n      \"code\": 0,\n      \"data\": null,\n      \"log\": \"\",\n      \"info\": \"\",\n      \"gasWanted\": \"0\",\n      \"gasUsed\": \"0\",\n      \"events\": [],\n      \"codespace\": \"\"\n    },\n    \"hash\": \"917E8691662385EA143F65DEE660A99D3DE04D08B5C1CADC99695ADBE04C5A05\",\n    \"height\": \"10\"\n  }\n}","result":true}
```

### 数据库查询

​	资产会保存到Mongodb中，通过Mongodb客户端可以查询chain数据库的assets集合。

## 路线图

| 序号 | 类别  | 任务                                         | 完成情况 |
| :--: | ----- | -------------------------------------------- | -------- |
|  1   | 基础  | 开发环境                                     | 完成     |
|  2   |       | 运行环境                                     | 完成     |
|  3   |       | CURL一键自动安装应用                         |          |
|  4   | 资产  | 基本key-value(mongledb)                      | 完成     |
|  5   |       | 账户account与签名                            | 完成     |
|  6   |       | 可分隔资产                                   |          |
|  7   | API   | API规划                                      | 完成     |
|  8   | WEBUI | 块浏览器                                     | 开发中   |
|  9   |       | 资产管理                                     |          |
|  10  |       | 账户管理                                     |          |
|  11  | 网络  | 协调节点                                     |          |
|  12  |       | 节点管理（增减节点）                         |          |
|  13  | 应用  | 生态系统（代码仓库、分布式应用、资产、价值） | 规划中   |

