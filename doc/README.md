# looklook-on-huaweicloud

本项目是基于go-zero微服务框架的最佳实践looklook开源项目使用华为云服务进行改造，帮助用户在华为云上构建高可用的Go微服务项目

项目地址 :  [looklook-on-huaweicloud](https://codehub-g.huawei.com/DTSE-GO/looklook-on-huaweicloud/files?ref=release)

## 1 项目说明

### 1.1 技术选型

- 微服务框架：go-zero
- 数据库：mysql
- 缓存：redis
- 消息队列：asynq
- 日志：filebeat+kafka+go-stash+es+kibana

### 1.2 项目结构

- app：业务代码
  -   api：api服务，对外提供服务，暴露http接口给外部客户端，为内部rpc服务实现api的聚合处理
  -   rpc：rpc服务，对内提供服务，用于内部rpc服务交互，使用protobuf+grpc进行rpc通信，rpc服务只需关注业务逻辑的处理
- common：通用方法error，middleware等
- deploy：组件部署配置
- modd.conf : modd热加载配置文件

### 1.3 系统架构

![pic](doc/架构图.png)

整个请求链路调用流程：

客户端发送请求到nginx网关，nginx会根据接口路由到对应的api服务，api服务作为BFF层，负责处理http接口请求，调用下游微服务并对返回的数据结果做聚合处理，最后返回统一的http响应给前端。api服务通过grpc与rpc服务进行通信，rpc服务负责具体的业务逻辑的处理，直接跟数据层进行交互。mq服务负责处理不同模块发布的消息队列的任务，下游主要跟各种消息队列对接，sheduler服务负责定时任务的定制，发布任务到asynq队列，job服务则负责获取asynq的任务并进行处理。日志系统则负责从各个微服务收集日志并展示到kibana上。

## 2 快速开始

下面演示如何搭建开发环境

### 2.1 必备程序

-   Go 1.16
-   Docker
-   Docker-compose
-   goctl

### 2.2 运行环境搭建

**1.克隆当前项目**

```sh
git clone https://codehub-dg-g.huawei.com/p30036525/looklook-on-huaweicloud.git
```

**2.打开项目并更新依赖**

```sh
cd go-zero-looklook-for-huaweicloud
go mod tidy
```

**3.启动项目环境**

```sh
docker network create go-zero-looklook_looklook_net
docker-compose -f docker-compose-env.yml up -d
```

**4.创建kafka topic**

```sh
docker exec -it kafka /bin/sh
cd /opt/kafka/bin/

#创建日志的topic
./kafka-topics.sh --create --zookeeper zookeeper:2181 --replication-factor 1 -partitions 1 --topic looklook-log
#创建支付服务的topic
./kafka-topics.sh --create --zookeeper zookeeper:2181 --replication-factor 1 -partitions 1 --topic payment-update-paystatus-topic
```

**5.数据库设置**

修改app目录的4个模块yaml配置文件的db连接：

```sh
DataSource: username:password@tcp(host:port)/database?charset=utf8mb4&parseTime=true&loc=Asia%2FShanghai
```

其中

- username: 用户名
- password: 用户密码
- host: 数据库IP地址
- port: 数据库端口
- database: 数据库名

创建数据库looklook，并运行/deploy/sql/下的sql脚本创建5个模块的数据表

**6.启动服务**

```sh
docker-compose up -d 
```

**7.检查服务运行状态**

访问  [http://127.0.0.1:9090/](http://127.0.0.1:9090/)，单击Status-\>Targets

![](doc/prometheus.png)

可以看到服务的启动状况

## 3 生产环境部署

下面演示通过华为云的多种服务来部署项目，构建一个企业高可用的线上服务

### 3.1 基础环境搭建

-   [云数据库RDS](https://support.huaweicloud.com/rds/index.html)，对4个不同模块的服务分别创建4个独立的数据库，根据不同业务进行分库，将/deploy/sql/下的sql脚本分别在不同数据库中运行
-   [分布式消息服务Kafka](https://support.huaweicloud.com/kafka/index.html)，创建主题looklook-log和payment-update-paystatus-topic
-   [云搜索服务CSS](https://support.huaweicloud.com/css/index.html)，作为es集群，用于存储日志，通过内置的kibana进行搜索
-   [分布式缓存服务DCS](https://support.huaweicloud.com/dcs/index.html，)，用做缓存处理和asynq消息队列的支持
-   [云容器引擎CCE](https://support.huaweicloud.com/cce/index.html)，基于k8s集群部署api和rpc微服务
-   [弹性云服务器ECS](https://gitee.com/link?target=https%3A%2F%2Fwww.huaweicloud.com%2Fproduct%2Fecs.html)，用于部署go-stash
-   [弹性负载均衡ELB](https://support.huaweicloud.com/elb/index.html)，为微服务配置网关，根据规则进行路由转发，并实现负载均衡

### 3.2 部署服务

![部署图](doc/部署图.png)

1. **克隆项目的release分支**

   ```
   git clone -b release https://codehub-dg-g.huawei.com/p30036525/looklook-on-huaweicloud.git 
   ```

2. **修改项目配置**

   -   中间件连接：将app目录下的5个模块的中间件\(kafka、redis、mysql\)的连接配置项改为对应组件的ip地址，使之能够正常连接
   -   服务发现：将rpc和api的服务配置改成k8s的endpoints

3. **cce镜像构建**

   使用goctl工具生成每个服务的Dockerfile

   ```
   goctl docker --version 1.16 -go xxx.go
   ```

   将镜像推送到华为云SWR上

   ```
   #使用临时指令登录
   docker login -u ${用户名} -p ${密码} ${区域}.myhuaweicloud.com
   docker build -t ${镜像名称}:${版本名称} .
   docker tag ${镜像名称}:${版本名称} ${区域}.myhuaweicloud.com/${组织名称}/${镜像名称}:${版本名称}
   docker push ${区域}.myhuaweicloud.com/${组织名称}/${镜像名称}:${版本名称}
   ```

4. **cce服务部署**

   使用goctl工具生成每个服务的k8s yaml文件

   ```
   goctl kube deploy -secret default-secret -replicas 2 -nodePort 3${端口} -requestCpu 200 -requestMem 50 -limitCpu 300 -limitMem 100 -name ${服务}-${类型} -namespace go-zero-looklook -image ${区域}.myhuaweicloud.com/${组织名称}/${镜像名称}:${版本名称} -o ${服务名} -port ${端口} --serviceAccount find-endpoints
   ```

   >**说明：** 
   >default-secret是默认的SWR镜像仓库凭证
   >指定find-endpoints为serviceAccount，以获取集群Endpoints的节点信息来实现k8s的服务发现

   进入cce控制台，在工作负载-\>无状态负载页面上，单击YAML创建，导入刚才生成的yaml创建deployment和service

5. **配置网关**

   点击服务发现-\>路由页面，创建路由，在转发策略配置中填写每个服务的路径和端口，通过华为云的elb进行路由转发，实现负载均衡

   ![](pic\路由.png)

6. **连接外部中间件**

   在cce集群信息中查找容器网段，在VPC的安全组中添加规则，在入方向规则中放通容器网段172.16.0.0/16访问，让k8s的容器能够连接上外部的中间件服务\(kafka、redis、mysql\)

   ![](pic\安全组.png)

7. **访问项目**

   通过elb绑定的弹性公网ip访问

   ```
   curl -X POST "http://139.xx.xx.xx/usercenter/v1/user/register" -H "Content-Type: application/json" -d "{\"mobile\":\"18888888888\",\"password\":\"123456\"}"
   ```

   成功则返回：

   ```
   {"code":200,"msg":"OK","data":{"accessToken":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NzM5NjY0MjUsImlhdCI6MTY0MjQzMDQyNSwiand0VXNlcklkIjo1fQ.E5-yMF0OvNpBcfr0WyDxuTq1SRWGC3yZb9_Xpxtzlyw","accessExpire":1673966425,"refreshAfter":1658198425}}
   ```

# 相关开源项目参考 
- 微服务框架: [go-zero](https://github.com/zeromicro/go-zero)
- 最佳实践: [go-zero-looklook]( https://github.com/Mikaelemmmm/go-zero-looklook) 
