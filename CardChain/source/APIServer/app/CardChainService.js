var path = require('path')
var co = require('co')
var fs = require('fs')
var util = require('util')
var Client = require('fabric-client')
var Peer = require('fabric-client/lib/Peer.js')
var User = require('fabric-client/lib/User.js')
var crypto = require('crypto')
var helper = require('./helper.js');

var log4js = require('log4js')
var logger = log4js.getLogger('Helper')
logger.setLevel('DEBUG')

var tempdir = "$HOME/fabric-client-kvs"

// let client = new Client()
// var channel = client.newChannel('cardchainchannel')
// var orderer = client.newOrderer('grpc://orderer.cystone.me:7050')
// channel.addOrderer(orderer)
// var peer0OrgBarber = client.newPeer('grpc://peer0.barber.cystone.me:7051')
// channel.addPeer(peer0OrgBarber)

var getBlockchainInfo = async function(name){
    console.log(name)
    return name
}

var query = async function(chaincodeName, fcn, queryArgs){
    console.log(chaincodeName)
    console.log(fcn)
    console.log(queryArgs)

    let member = await getOrgUser4Local()
    var request = {
        chaincodeId: chaincodeName,
        fcn: fcn,
        args: queryArgs
    }
    console.log(request)
    let response_payloads = await channel.queryByChaincode(request)
    
    if(response_payloads) {
        return response_payloads
    } else {
        return 'response_payloads is null'
    }
}

var invoke = function(chaincodeId, fcn, cardAgrs){
    return getOrgUser4Local().then((user) => {
        tx_id = client.newTransactionID()
        var request = {
            chaincodeId: chaincodeId,
            txId: tx_id,
            fcn: fcn,
            args: cardAgrs
        }
        return channel.sendTransactionProposal(request)
    }, (err) => {
        console.log('error', e)
    }).then((result) => {
        var proposalResponses = result[0]
        var proposal = result[1]
        var all_good = true;
        for (var i in proposalResponses) {
        	let one_good = false;
        	if (proposalResponses && proposalResponses[i].response &&
        		proposalResponses[i].response.status === 200) {
        		one_good = true;
        		logger.info('invoke chaincode proposal was good');
        	} else {
        		logger.error('invoke chaincode proposal was bad');
        	}
        	all_good = all_good & one_good;
        }

        if (all_good) {
            console.info(util.format(
        		'Successfully sent Proposal and received ProposalResponse: Status - %s, message - "%s", metadata - "%s", endorsement signature: %s',
        		proposalResponses[0].response.status, proposalResponses[0].response.message,
        		proposalResponses[0].response.payload, proposalResponses[0].endorsement.signature
            ))

            var orderer_request = {
                proposalResponses: proposalResponses,
                proposal: proposal,
                txID: tx_id
            }
            return channel.sendTransaction(orderer_request)
        }
    }, (err) => {
        console.log('error', e)
    }).then((sendtransactionresult) => {
        return sendtransactionresult
    }, (err) => {
        console.log('error', e)
    })
}

function getOrgUser4Local() {
    var keyPath = "/home/stone/Documents/TianChi/CardChain/artifacts/channel/crypto-config/peerOrganizations/barber.cystone.me/users/Admin@barber.cystone.me/msp/keystore"
    var keyPEM = Buffer.from(readAllFiles(keyPath)[0]).toString()
    var certPath = "/home/stone/Documents/TianChi/CardChain/artifacts/channel/crypto-config/peerOrganizations/barber.cystone.me/users/Admin@barber.cystone.me/msp/signcerts"
    var certPEM = readAllFiles(certPath)[0].toString()

    return Client.newDefaultKeyValueStore({
        path:tempdir
    }).then((store)=>{
        client.setStateStore(store)

        return client.createUser({
            username: 'Admin',
            mspid: 'OrgBarberMSP',
            cryptoContent:{
                privateKeyPEM: keyPEM,
                signedCertPEM:certPEM
            }
        })
    })
}

function readAllFiles(dir){
    var files = fs.readdirSync(dir)
    var certs = []
    files.forEach((file_name) => {
        let file_path = path.join(dir, file_name)
        let data = fs.readFileSync(file_path)
        certs.push(data)
    })
    return certs
}

exports.invoke = invoke
exports.query = query
