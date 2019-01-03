## Fabric 分布式部署

#### 部署方式
三台主机，两个Org，每个Org各两个Peer、各一个Orderer，三个Kafka节点，一个Zookeeper。


| Host1       | Host2       | Host3       |
| :---------: | :---------: | :---------: |
|             | Zookeeper   |             |
| Kafka0      | Kafka1      | Kafka2      |
| OrgB Peer0  |	OrgA Peer0  | OrgA Peer1  |
|             |	Orderer     |             |
|             |	            |             |

Host资源有限，就这样分配了。

#### 一、部署kafka

部署过程可参考[Kafka官方文档](http://kafka.apache.org/quickstart)。

1. 下载

    从[这里(http://kafka.apache.org/downloads)](http://kafka.apache.org/downloads)下载kafka。三个主机都要下载。

2. 解压

    三个主机分别解压。
    ```
    $ tar -xzf kafka_2.11-2.1.0.tgz
    $ cd kafka_2.11-2.1.0
    ```

3. 启动

    三个主机都需要启动。Zookeeper 和 Kafka0 在 Host2 启动，Kafka1 在 Host1 启动， Kafka2 在 Host3 启动。启动前需要先修改配置：

    - ** config/server.properties 三个主机都要修改: **

    - broker.id = id   // 每个配置文件的id必须唯一，且为整数

    - listeners=PLAINTEXT://:9092  // 删除这一行前边的“#”

    - advertised.listeners=PLAINTEXT://your.host.name:9092 // 删除这一行前边的“#”,并且加上主机ip

    - log.dirs=/tmp/kafka-logs   // 将路径改为其他你认为合适的位置，保证路径有写入权限，比如你的home目录。/tmp 下的文件可能会在每次重启或者固定时间被清除。

    - zookeeper.connect=localhost:2181   // 需要改为zookeeper运行的主机的ip，此处为 Host2 的ip

    - ** config/zookeeper.properties 只改 Host2 即可: **

    - dataDir=/tmp/zookeeper    // 将路径改为其他你认为合适的位置，保证路径有写入权限，比如你的home目录。/tmp 下的文件可能会在每次重启或者固定时间被清除。

    然后，在Host2上启动Zookeeper：
    ```
    $ bin/zookeeper-server-start.sh config/zookeeper.properties   // 启动Zookeeper
    ```

    在三个主机上分别启动Kafka：
    ```
    $ bin/kafka-server-start.sh config/server.properties          // 启动Kafka
    ```

4. 测试

    在本地主机即可测试Kafka是否搭建成功。进入kafka解压目录，进行如下操作。
    ```
    $ bin/kafka-topics.sh --create --zookeeper 148.70.109.243:2181 --replication-factor 3 --partitions 1 --topic my-replicated-topic // 创建topic
    $ bin/kafka-topics.sh --describe --zookeeper 148.70.109.243:2181 --topic my-replicated-topic    // 查看topic信息
    ```

    向生产者发送消息：
    ```
    $ bin/kafka-console-producer.sh --broker-list 148.70.109.243:9092 --topic my-replicated-topic
    > my test message 1
    > my test message 2
    ^C
    ```

    消费者消费消息（可在一个新终端中执行该命令）：
    ```
    $ bin/kafka-console-consumer.sh --bootstrap-server 148.70.109.243:9092 --from-beginning --topic my-replicated-topic

    my test message 1
    my test message 2
    ^C
    ```

    能正常收发消息即代表kafka正常运行。

Kafka 搭建完毕。

进行下边操作之前，先把代码从github同步下来。仓库地址：[https://github.com/SuperMap/ATLab-ATLChain](https://github.com/SuperMap/ATLab-ATLChain)

#### 二、部署Orderer

登录Host2。先从 configs/ 拷贝并修改配置文件：
```
$ cd path/to/ATLab-ATLChain/config
$ cp configs/orderer-samples.yaml orderer.yaml
$ cp configs/configtx-samples.yaml configtx.yaml
```

需要修改的配置文件有：

1. configtx.yaml->Organizations->AnchorPeers  修改锚节点的Host

2. configtx.yaml->Orderer->Addresses 修改Orderer的ip地址

3. configtx.yaml->Orderer->Kafka->Brokers 添加kafka Broker的ip和Port

然后启动orderer：
```
$ configtxgen -profile OrdererChannel -outputBlock genesisblock -channelID ordererchannel   // 生成orderer创世块
$ nohup orderer start > orderer.log 2>&1 &
```

当orderer.log中出现类似下边的内容时，表明启动成功
```
xxxxx Channel consumer set up successfully
xxxxx Start phase completed successfully
```

#### 三、部署Peers

1. 登录Host2。先从 configs/ 拷贝并修改配置文件：
    ```
    $ cd path/to/ATLab-ATLChain/config
    $ cp configs/core-samples.yaml core.yaml
    ```

    需要修改的配置文件有：
    1. core.yaml->peer->id 改成peer0.orga
    2. core.yaml->peer->gossip->externalEndpoint 改成本机公网ip

    然后启动peer：
    ```
    $ nohup peer node start > peer.log 2>&1 &
    $ export CORE_PEER_LOCALMSPID=OrgA
    $ export CORE_PEER_MSPCONFIGPATH=$PWD/crypto-config/peerOrganizations/orga.atlchain.com/users/Admin@orga.atlchain.com/msp
    $ configtxgen -profile TxChannel -outputCreateChannelTx atlchannel.tx -channelID atlchannel
    $ configtxgen -profile TxChannel -outputAnchorPeersUpdate OrgAanchors.tx -channelID atlchannel -asOrg OrgA
    $ peer channel create -o Orderer-IP:7050 -c atlchannel -f atlchannel.tx
    $ peer channel join -b atlchannel.block
    $ peer channel update -o Orderer-IP:7050 -c atlchannel -f OrgAanchors.tx

    // 多机部署是要先打包链码，然后每个节点都安装链码包，不然会报链“chaincode fingerprint mismatch: data mismatch”这个错误。
    $ peer chaincode package -n cc -p github.com/hyperledger/fabric/examples/chaincode/go/example02/cmd -v 1.0 ccpack.out

    $ peer chaincode install sacc/ccpack.out
    ```

    测试peer是否运行成功：
    ```
    $ ./run query   // 查询a的值，应为 100
    $ ./run invoke  // a向b转移1
    $ ./run query   // 查询a的值，应为 99
    ```

OK，一个节点部署完毕。

2. 登录Host1。先从 configs/ 拷贝并修改配置文件：
    ```
    $ cd path/to/ATLab-ATLChain/config
    $ cp configs/core-samples.yaml core.yaml
    ```

    需要修改的配置文件有：
    1. core.yaml->peer->id 改成peer1.orga

    2. core.yaml->peer->gossip->externalEndpoint 改成本机公网ip

    3. core.yaml->peer->gossip->bootstrap 改成Host2公网ip，为了更快找到peer0

    4. core.yaml->peer->mspConfigPath 改成peer1.orga对应的msp路径

    然后启动peer:
    ```
    $ nohup peer node start > peer.log 2>&1 &
    $ export CORE_PEER_LOCALMSPID=OrgA
    $ export CORE_PEER_MSPCONFIGPATH=$PWD/crypto-config/peerOrganizations/orga.atlchain.com/users/Admin@orga.atlchain.com/msp
    $ scp Host2:/path/to/atlchannel.block ./    // 将Host2中的通道创世块拷过来
    $ peer channel join -b atlchannel.block     // 加入通道
    $ peer chaincode install ccpack.out         // 安装链码
    ```

    然后按上边的步骤测试peer是否运行成功。
    > 这个地方要稍微等一会儿，让peer的账本数据同步完成，不然查询不到数据。

3. 登录Host3。同样从 configs/ 拷贝并修改配置文件。
    ```
    $ cd path/to/ATLab-ATLChain/config
    $ cp configs/core-samples.yaml core.yaml
    $ cp configs/configtx-samples.yaml configtx.yaml    // 因为要通知锚节点，所以要这个
    ```

    需要修改的配置文件有：

    1. configtx.yaml->Organizations->AnchorPeers  修改锚节点的Host。这里不需要处理Orderer的配置，只修改OrgB的锚节点ip即可

    2. core.yaml->peer->id 改成peer0.orgb

    3. core.yaml->peer->gossip->externalEndpoint 改成本机公网ip

    4. core.yaml->peer->mspConfigPath 改成peer0.orgb对应的msp路径

    然后启动peer：
    ```
    $ nohup peer node start > peer.log 2>&1 &
    $ export CORE_PEER_LOCALMSPID=OrgB
    $ export CORE_PEER_MSPCONFIGPATH=$PWD/crypto-config/peerOrganizations/orgb.atlchain.com/users/Admin@orgb.atlchain.com/msp
    $ configtxgen -profile TxChannel -outputCreateChannelTx atlchannel.tx -channelID atlchannel
    $ configtxgen -profile TxChannel -outputAnchorPeersUpdate OrgBanchors.tx -channelID atlchannel -asOrg OrgB
    $ peer channel create -o Orderer-IP:7050 -c atlchannel -f atlchannel.tx
    $ peer channel join -b atlchannel.block
    $ peer channel update -o Orderer-IP:7050 -c atlchannel -f OrgBanchors.tx

    $ peer chaincode install sacc/ccpack.out
    ```

    然后按上边的步骤测试peer是否运行成功。
    > 这个地方要稍微等一会儿，让peer的账本数据同步完成，不然查询不到数据。
