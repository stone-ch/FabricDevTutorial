## CHFA 考试体验

#### CHFA 介绍

CHFA （Certified Hyperledger Fabric Administrator，认证 Hyperledger Fabric 管理员）。
是 Linux 基金会发放的认证证书，按照官方的说法，获得此证书的人，具备搭建一个安全的
可商用的 Hyperledger Fabric 网络的能力，其中包括对网络的节点进行安装、配置、操作、
管理和排错的能力。

证书有效期为 2 年，报名费 300 美元，报名有效期为 12 个月，有一次免费重考的机会，考
试形式为在线考试，考试时长两个小时。考试基于 Fabric 1.3 版本，操作系统环境为 Ubuntu 16。
这个考试有一个好处就是，考试开始前 24 小时，可以取消或者重新预约，离开了学校之后，
时间可能会随时被占用，这个还是挺人性化的。

更多详细信息，以及注册报名请访问 CHFA 官网 [https://training.linuxfoundation.org/certification/certified-hyperledger-fabric-administrator-chfa/](https://training.linuxfoundation.org/certification/certified-hyperledger-fabric-administrator-chfa/) 。

#### 考试前的准备

注册、报名和预约考试的过程直接在官网看就好了，不算太复杂，唯一一个要注意的是，需要
提前注册好 Linux Foundation 的账号，考试过程中最好统一使用这个账号登录。如果需要发
票的话需要提前和考务方联系，官网提供有邮箱，确定好开票信息。另外发票是美国那边开的，
和国内的不太一致，这些需要提前沟通好。

考试为在线考试，但是对考试的环境要求比较严格。首先说软件环境，必须使用　Chrome 或 
Chromium 浏览器；另外必须装一个插件（系统中有链接可以直接安装）；电脑必须有摄像头、
耳机和麦克风，因为考试全程都有监考官在看着你，同时你的屏幕也是处于分享状态；另外网
络带宽也有要求，具体值我忘了，要求不是很高，但是要 **F Q** ，否则进不去考试页面。不
用担心不知道自己的电脑是否达到了要求，在考试报名的时候，就有链接让你检查你的电脑配
置，达到配置要求也并不是很难。

另外还有你考试的环境，要求不能在公共区域，比如咖啡馆之类的，必须在一个比较安静、封
闭的环境中，考试期间不能有人打扰，而且考试开始之前考官会要求你将摄像头旋转 360 度，
查看周围环境，同时也会仔细查看你的桌子，要求桌面要干净，不能有其他东西。而且在考试
过程中要求你的眼睛要盯着屏幕，不能左顾右盼，也不能遮挡脸部，我中间无意识的捂了一下
嘴就被监考官警告了。

这些电脑和考试地点必须提前准备好，最好提前一天先确定一下。

#### 进入考试

##### 核对信息

我预约的是 2019.04.20 上午 10:00-12:00 考试，提前 15 分钟进入系统，监考官提前和你沟
通，确认你的信息和考试环境，沟通是通过页面上的一个聊天窗口，英文沟通。这时候有两个问
题需要注意，一个是你的用户名，监考官要核对你的身份信息，可能他们更习惯看护照，但是我
没有，我就给他看的身份证，身份证上是中文的名字，我的用户名是拼音，所以就要求我把 Linux 
Foundation 账号的用户名改成了和身份证一致的名字。还有一个就是监考官要你打开 System Monitor，
我用的 Ubuntu 系统，一开始没反应过来，后来打开 htop 给他看了一眼就 OK 了。确认完信息
之后，就让你旋转摄像头，查看屋子里的环境和桌面。在这个环节，监考官会把考试要求很详细
的讲给你。比如在考试过程中，在浏览器中只能打开考试页面外的一个 tab 标签，并且这个标签
中打开的网页只能是 Fabric 文档和 FabricCA 文档两个链接及其子页面，否则就算违规。确认
完成之后，监考官就会开始你的考试了。


##### 开始答题

一共 16 道考题，120 分钟，算下来每道题的时间还是很紧张的，毕竟每一个操作都有一长串的
命令要敲。每道题的网络都和 byfn 的网络类似，熟悉 byfn 就能感觉到轻松很多。而且整体来
说，题目并不算难，只要能熟练掌握 [tutorial](https://hyperledger-fabric.readthedocs.io/en/release-1.3/tutorials.html) 
和 CA 的使用就可以。

考试的界面是，左边是题目说明（英文），右边是 shell 窗口，每道题都需要 ssh 登录到不同
的 host 去操作，只有一个窗口，但是可以使用 tmux 进行分屏，可以使用 vim 进行文本编辑。
大概说一下题目的类型：

1. 链码的安装和升级：给定一个网络安装了 ```testcc``` 的 ```1.0``` 版本，让你升级为 
```2.0``` 版本，并将结果保存到一个 txt 文件。请参考 [Chaincode for Operators](https://hyperledger-fabric.readthedocs.io/en/release-1.3/chaincode4noah.html) 。
主要涉及到的命令有：
    ```
    # 安装、升级链码
    $ peer chaincode install ...... 
    $ peer chaincode upgrade ...... # **注意定义背书策略**

    # 输出已实例化的链码列表到 txt 文件
    $ peer chaincode list --instantiated -C mychannel > instantiated.txt
    ```

2. 修复网络中的问题：网络中使用 kafka 的 orderer 节点不能正常启动，要求你查看 docker 
日志并解决问题。一般是要修改 configtx.yaml 的配置，然后重启网络。请参考[Bringing up a Kafka-based Ordering Service](https://hyperledger-fabric.readthedocs.io/en/release-1.3/kafka.html) 。
主要涉及到的命令有：
    ```
    # 查看 orderer 错误日志
    $ docker logs -f orderer.example.com
    
    # 重启网络
    $ docker-compose -f docker-compose.yaml # yaml 文件名根据实际情况确定
    ```
3. 升级网络版本：参考 [Upgrading Your Network Components](https://hyperledger-fabric.readthedocs.io/en/release-1.3/upgrading_your_network_tutorial.html) 。

4. 使用 Fabric CA 注册一个账户：题目中给你了新建的账户的要求，按要求注册一个账户。
请参考 [Enrolling a peer identity](https://hyperledger-fabric-ca.readthedocs.io/en/latest/users-guide.html#registering-a-new-identity) 。
主要涉及到的命令有：
    ```
    # 注册用户信息
    $ fabric-ca-client register -d --id.name admin2 --id.affiliation org1.department1 --id.attrs '"hf.Registrar.Roles=peer,user"' --id.attrs hf.Revoker=true

    # 获得用户的密钥文件
    $ fabric-ca-client enroll -u http://peer1:peer1pw@localhost:7054 -M $FABRIC_CA_CLIENT_HOME/msp
    ```

5. fabric-ca-service 问题修复：主要是先查看 fabric-ca-service 启动时的报错，然后修
改 ```fabric-ca-server-config.yaml``` 文件，然后启动 fabric-ca-service 就好。请参
考 [Fabric CA Server](https://hyperledger-fabric-ca.readthedocs.io/en/latest/users-guide.html#initializing-the-server) 。
主要涉及的命令：
    ```
    # 启动 FabricCA Server
    $ fabric-ca-server start -b <admin>:<adminpw>
    ```

6. 创建通道，加入通道：请参考 [Create & Join Channel](https://hyperledger-fabric.readthedocs.io/en/release-1.3/build_network.html#create-join-channel) 。
主要涉及到的命令有：
    ```
    # 生成通道创世区块
    $ peer channel create -o orderer.example.com:7050 -c $CHANNEL_NAME -f path/to/channel_tx_file --tls --cafile path/to/cafile

    # 节点加入通道
    $ peer channel join -b mychannel.block
    ```

7. 修改区块中交易的数量和生成区块的时间：请参考[Editing a Config](https://hyperledger-fabric.readthedocs.io/en/release-1.3/config_update.html#editing-a-config) 。

8. 修改通道的控制权限：主要是修改 ```configtx.yaml``` 配置文件。请参考 [Access Control Lists (ACL)](https://hyperledger-fabric.readthedocs.io/en/release-1.3/access_control.html) 。

9. 配置 CA 的 HSM (Hardware Security Module)：请参考 [HSM](https://hyperledger-fabric-ca.readthedocs.io/en/latest/users-guide.html#configuring-fabric-ca-server-to-use-softhsm2) 。

大概是这些了，还有几个内容我这次的试题中没有涉及到，但是我觉得准备的时候也要注意一下：

10. 使用私有数据： 请参考 [Using Private Data in Fabric](https://hyperledger-fabric.readthedocs.io/en/release-1.3/private_data_tutorial.html) 。

11. 向通道中增加组织：请参考 [Adding an Org to a Channel](https://hyperledger-fabric.readthedocs.io/en/release-1.3/channel_update_tutorial.html) 。

12. 配置 peer 使用 CouchDB ： 请参考 [Using CouchDB](https://hyperledger-fabric.readthedocs.io/en/release-1.3/couchdb_tutorial.html) 。

13. 配置 CoucbDB 的索引： 请参考 [Create an index](https://hyperledger-fabric.readthedocs.io/en/release-1.3/couchdb_tutorial.html#create-an-index) 。

#### 结语

我这次准备的不是很充分，16 道题只做完了 8 道，过线是没什么希望了，还有一次重考的
机会，再考不过 300 美刀就要打水漂了，继续加油准备。
