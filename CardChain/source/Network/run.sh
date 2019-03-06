#!/bin/sh

export PATH=$PATH:./bin

if [ $1 != clean ]
then
    echo "### export CHANNEL_NAME=cardchannel ###"
    export CHANNEL_NAME=cardchannel
    
    echo "### export CORE_PEER_LOCALMSPID=OrgMarket ###"
    export CORE_PEER_LOCALMSPID=OrgMarket
     
    echo "### export CORE_PEER_MSPCONFIGPATH=./crypto-config/peerOrganizations/market.cardchain.com/users/Admin@market.cardchain.com/msp ###"
    export CORE_PEER_MSPCONFIGPATH=./crypto-config/peerOrganizations/market.cardchain.com/users/Admin@market.cardchain.com/msp
fi

if [ $1 = clean ]
then
    rm cardchannel.tx cardchannel.block genesisblock OrgMarketAnchors.tx

    rm -rf ./production

    # rm -rf ./crypto-config

elif [ $1 = s1 ]
then
    # echo "### configtxgen -profile OrdererOrg -outputBlock orderer.genesis.block ###"
    # cryptogen generate --config=crypto-config.yaml

    echo "### configtxgen -profile OrdererChannel -outputBlock genesisblock ###"
    configtxgen -profile OrdererChannel -outputBlock genesisblock -channelID ordererchannel

    echo "### nohup orderer start > orderer.log 2>&1 & ###"
    nohup orderer start > orderer.log 2>&1 &

    echo "### nohup peer node start > peer.log 2>&1 & ###"
    nohup peer node start > peer.log 2>&1 &

    echo "### Finish stage 1 ###"

elif [ $1 = s2 ]
then
    echo "### configtxgen -profile CardChannel -outputCreateChannelTx atlchannel.tx -channelID $CHANNEL_NAME ###"
    configtxgen -profile CardChannel -outputCreateChannelTx cardchannel.tx -channelID $CHANNEL_NAME

    echo "### configtxgen -profile CardChannel -outputAnchorPeersUpdate OrgMarketAnchors.tx -channelID $CHANNEL_NAME -asOrg $CORE_PEER_LOCALMSPID ###"
    configtxgen -profile CardChannel -outputAnchorPeersUpdate OrgMarketAnchors.tx -channelID $CHANNEL_NAME -asOrg $CORE_PEER_LOCALMSPID

    echo "### peer channel create -o 127.0.0.1:7050 -c cardchannel -f cardchannel.tx ###"
    peer channel create -o 127.0.0.1:7050 -c cardchannel -f cardchannel.tx

    echo "### peer channel join -b cardchannel.block ###"
    peer channel join -b cardchannel.block
    
    echo "###  peer channel update -o 127.0.0.1:7050 -c $CHANNEL_NAME -f OrgMarketAnchors.tx ###"
    peer channel update -o 127.0.0.1:7050 -c $CHANNEL_NAME -f OrgMarketAnchors.tx

    echo "### peer chaincode install -n cc -v 1.0 -p github.com/hyperledger/fabric/examples/chaincode/go/example02/cmd/ ###"
    peer chaincode install -n cc -v 1.0 -p github.com/hyperledger/fabric/examples/chaincode/go/example02/cmd/ 
    
    echo "### peer chaincode instantiate -o 127.0.0.1:7050 -C $CHANNEL_NAME -n cc -v 1.0 -c '{"Args":["init","a","100","b","200"]}' ###"
    peer chaincode instantiate -o 127.0.0.1:7050 -C $CHANNEL_NAME -n cc -v 1.0 -c '{"Args":["init","a","100","b","200"]}'
     
    echo "### Finish stage 2 ###"
elif [ $1 = query ]
then 
    echo "### query ###"
    peer chaincode query -C $CHANNEL_NAME -n cc -c '{"Args":["query","a"]}'
elif [ $1 = invoke ]
then
    echo "### invoke ###"
    peer chaincode invoke -o 127.0.0.1:7050 -C $CHANNEL_NAME -n cc -c '{"Args":["invoke","a","b","1"]}'
elif [ $1 = stop ]
then
    echo "### ###"
fi
