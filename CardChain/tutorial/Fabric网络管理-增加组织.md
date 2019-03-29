1. 生成组织相关密钥文件

../../bin/cryptogen generate --config=./org3-crypto.yaml

2. 将 Org3 相关配置打印为 JSON 文件

export FABRIC_CFG_PATH=$PWD && ../../bin/configtxgen -printOrg Org3MSP > ../channel-artifacts/org3.json

3. 将排序节点的密钥复制到 Org3 的对应目录下

cd ../ && cp -r crypto-config/ordererOrganizations org3-artifacts/crypto-config/

4. 进入 cli 容器。

docker exec -it cli bash

5. 设置 cli 容器中的环境变量

export ORDERER_CA=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem  && export CHANNEL_NAME=mychannel

6. 获取最新版的通道配置

peer channel fetch config config_block.pb -o orderer.example.com:7050 -c $CHANNEL_NAME --tls --cafile $ORDERER_CA

7. 将配置转换到 JSON 并裁剪

configtxlator proto_decode --input config_block.pb --type common.Block | jq .data.data[0].payload.data.config > config.json

8. 追加 Org3 配置

jq -s '.[0] * {"channel_group":{"groups":{"Application":{"groups": {"Org3MSP":.[1]}}}}}' config.json ./channel-artifacts/org3.json > modified_config.json

9. 编码 protobuf 文件

configtxlator proto_encode --input config.json --type common.Config --output config.pb

configtxlator proto_encode --input modified_config.json --type common.Config --output modified_config.pb

10. 计算配置的差异

configtxlator compute_update --channel_id $CHANNEL_NAME --original config.pb --updated modified_config.pb --output org3_update.pb

11. 将包含差异的 proto 部分解码成 JSON 文件

configtxlator proto_decode --input org3_update.pb --type common.ConfigUpdate | jq . > org3_update.json

12. 使用信封消息包装差异部分

echo '{"payload":{"header":{"channel_header":{"channel_id":"mychannel", "type":2}},"data":{"config_update":'$(cat org3_update.json)'}}}' | jq . > org3_update_in_envelope.json

13. 格式化好的 JSON 文件，再次编码成 protobuf

configtxlator proto_encode --input org3_update_in_envelope.json --type common.Envelope --output org3_update_in_envelope.pb

14. Org1 签名 protobuf

peer channel signconfigtx -f org3_update_in_envelope.pb

15. Org2 执行通道配置更新（先切换但前使用的组织）


export CORE_PEER_LOCALMSPID="Org2MSP"

export CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt

export CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp

export CORE_PEER_ADDRESS=peer0.org2.example.com:7051

peer channel update -f org3_update_in_envelope.pb -c $CHANNEL_NAME -o orderer.example.com:7050 --tls --cafile $ORDERER_CA

16. 设置动态选举领导节点

CORE_PEER_GOSSIP_USELEADERELECTION=true
CORE_PEER_GOSSIP_ORGLEADER=false

17. 启动 Org3 的节点

docker-compose -f docker-compose-org3.yaml up -d

18. 进入 Org3 cli 容器

docker exec -it Org3cli bash

19. 配置环境变量

export ORDERER_CA=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem && export CHANNEL_NAME=mychannel

20. 获取排序组织区块

peer channel fetch 0 mychannel.block -o orderer.example.com:7050 -c $CHANNEL_NAME --tls --cafile $ORDERER_CA

21. 节点加入通道

peer channel join -b mychannel.block

22. 将第二个节点加入通道

export CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org3.example.com/peers/peer1.org3.example.com/tls/ca.crt && export CORE_PEER_ADDRESS=peer1.org3.example.com:7051

peer channel join -b mychannel.block

23. 分别在 Org1 Org2 Org3 的节点上安装和更新链码（需要配置环境变量来更改所连接的节点）

peer chaincode install -n mycc -v 2.0 -p github.com/chaincode/chaincode_example02/go/

peer chaincode upgrade -o orderer.example.com:7050 --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -C $CHANNEL_NAME -n mycc -v 2.0 -c '{"Args":["init","a","90","b","210"]}' -P "OR ('Org1MSP.peer','Org2MSP.peer','Org3MSP.peer')"

24. 查询链码 

peer chaincode query -C $CHANNEL_NAME -n mycc -c '{"Args":["query","a"]}'

25. 调用链码

peer chaincode invoke -o orderer.example.com:7050  --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -C $CHANNEL_NAME -n mycc -c '{"Args":["invoke","a","b","10"]}'

peer chaincode query -C $CHANNEL_NAME -n mycc -c '{"Args":["query","a"]}'
