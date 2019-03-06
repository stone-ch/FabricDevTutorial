var co = require('co')
var cardchainservice = require('./app/CardChainService.js')
var express= require('express')
// var helper = require('./app/helper.js')
// var jwt = require('jsonwebtoken') 
// var config = require('./app/connection-profile-standard.json');
// var cors = require('cors');
// var bodyParser = require('body-parser');
// var bearerToken = require('express-bearer-token');

var app = express()

///////////////////////////////////////////////////////////////////////////////
//////////////////////////////// SET CONFIGURATONS ////////////////////////////
///////////////////////////////////////////////////////////////////////////////
app.options('*', cors());
app.use(cors());
//support parsing of application/json type post data
app.use(bodyParser.json());
//support parsing of application/x-www-form-urlencoded post data
app.use(bodyParser.urlencoded({
    extended: false
}));

// set secret variable
app.set('secret', 'thisismysecret');
// app.use(expressJWT({
//  secret: 'thisismysecret'
// }).unless({
//  path: ['/users']
// }));
app.use(bearerToken());

process.TOKENS = [];

app.use(function(req, res, next) {
    if (req.originalUrl.indexOf('/users') >= 0 || req.originalUrl.indexOf('/ng') >= 0 || req.originalUrl.indexOf('/login') >= 0) {
        return next();
    }

    var token = req.token || req.body.token;
    jwt.verify(token, app.get('secret'), function(err, decoded) {
        if (err || process.TOKENS[decoded.username] == null) {
            res.send({
                success: false,
                message: 'Failed to authenticate token. Make sure to include the ' +
                'token returned from /createAccount call in the authorization header ' +
                ' as a Bearer token'
            });
            return;
        } else {
            if (process.TOKENS[decoded.username] != null && process.TOKENS[decoded.username] != token) {
                res.send({
                    success: false,
                    message: 'token has expired'
                });
                return;
            }
            // add the decoded user name and org name to the request object
            // for the downstream code to use
            req.username = decoded.username;
            req.orgname = decoded.orgName;
            logger.info(util.format('Decoded from JWT token: username - %s, orgname - %s', decoded.username, decoded.orgName));
            return next();
        }
    });

});

// 1. issueCard args: 0-{CardID},1-{CardBalance},2-{CardPoints},3-{CardPointsBase},4-{CardCreateTime},5-{CardState},6-{CustomerID},7-{CustomerGender},8-{CustomerAge},9-{CustomerTEL}
// http://localhost:3000/issuecard?cardid=OrgBarberMSP00000000&cardbalance=10&cardpoints=0&cardpointsbase=10&cardcreatetime=20180831170800&cardstate=1&customerid=00000000&customergender=M&customerage=18&customertel=13900001111
app.get('/issuecard', function(req, res){
    co (function *(){        
        var cardAgrs = [req.query.cardid, req.query.cardbalance, req.query.cardpoints, req.query.cardpointsbase, req.query.cardcreatetime, req.query.cardstate, req.query.customerid, req.query.customergender, req.query.customerage, req.query.customertel]
        var result = yield cardchainservice.invoke("cardchaincc","issueCard",cardAgrs)    
        res.send("successCallback()")
    }).catch((err) => {
        res.send(err)
    })
})

// 2. queryCardByCardID args: 0 - {CardID}
// http://localhost:3000/querycardbycardid?cardid=OrgBarberMSP1234567890
app.get('/querycardbycardid', function(req, res){
    co(function *(){
        var queryArgs = [req.query.cardid]
        var result = yield cardchainservice.query("cardchaincc", "queryCardByCardID", queryArgs)
        res.send("successCallback(" + result.toString() + ")")
    }).catch((err) => {
        res.send(err)
    })
})

// 3. makeAgreement args: 0-{OrgID},1-{AgreementOrg},2-{AgreementPrice},3-{AgreementCount},4-{AgreementCreateTime},5-{AgreementDeadline},6-{AgreementDescription}
// ["OrgBarberMSP", "OrgBarberMSP", "120", "0", "20180902151200", "20190902151200", "充100得120"]
// http://localhost:3000/makeagreement?orgid=OrgBarberMSP&agreementorg=OrgBarberMSP&agreementprice=120&agreementcount=0&agreementcreatetime=20180902151200&agreementdeadline=agreementdeadline&agreementdescription=%E5%85%85100%E5%BE%97120

// ["OrgBarberMSP", "OrgBarberMSP", "40,40,10,10", "4", "20180902151200", "20190902151200", "100元享受4次原价30元洗头"]
// http://localhost:3000/makeagreement?orgid=OrgBarberMSP&agreementorg=OrgBarberMSP&agreementprice=40,40,10,10&agreementcount=4&agreementcreatetime=20180902151200&agreementdeadline=agreementdeadline&agreementdescription=100%E5%85%83%E4%BA%AB%E5%8F%974%E6%AC%A1%E5%8E%9F%E4%BB%B730%E5%85%83%E6%B4%97%E5%A4%B4
app.get('/makeagreement', function(req, res){
    co (function *(){
        var agreementArgs = [req.query.orgid, req.query.agreementorg, req.query.agreementprice, req.query.agreementcount, req.query.agreementcreatetime, req.query.agreementdeadline, req.query.agreementdescription]
        var result = yield cardchainservice.invoke("cardchaincc", "makeAgreement", agreementArgs)
        res.send("successCallback()")
        // res.send("successCallback(" + result.toString() + ")")
    }).catch((err) => {
        res.send(err)
    })
})

// 4. queryAgreementByRange args:0-{OrgID} "OrgBarberMSP"
// http://localhost:3000/queryagreementbyrange?orgid=OrgBarberMSP
app.get('/queryagreementbyrange', function(req, res){
    co (function *(){
        var queryArgs = [req.query.orgid]
        var result = yield cardchainservice.query("cardchaincc", "queryAgreementByRange", queryArgs)
        res.send("successCallback(" + result.toString() + ")")
    }).catch((err) => {
        res.send(err)
    })
})

// 5. queryAgreementByAgreementID args: 0-{AgreementID} "Agr_OrgBarberMSP109"
// http://localhost:3000/queryagreementbyagreementid?agreementid=Agr_OrgBarberMSP109
app.get('/queryagreementbyagreementid', function(req, res){
    co (function *(){
        var queryArgs = [req.query.agreementid]
        var result = yield cardchainservice.query("cardchaincc", "queryAgreementByAgreementID", queryArgs)
        res.send("successCallback(" + result.toString() + ")")
    }).catch((err) => {
        res.send(err)
    })
})

// 6. queryUAgrByRange agrs:0-{CardID}
// http://localhost:3000/queryuagrbyrange?cardid=OrgBarberMSP1234567890
app.get('/queryuagrbyrange', function(req, res){
    co (function *(){
        var queryArgs = [req.query.cardid]
        var result = yield cardchainservice.query("cardchaincc", "queryUAgrByRange", queryArgs)
        res.send("successCallback(" + result.toString() + ")")
    }).catch((err) => {
        res.send(err)
    })
})

// 7. queryUAgrByUAgrID args: 0 - {UAgrID} "UAgr_OrgBarberMSP1234567890M102"
// http://localhost:3000/queryuagrbyuagrid?uagrid=UAgr_OrgBarberMSP1234567890M102
app.get('/queryuagrbyuagrid', function(req, res){
    co (function *(){
        var queryArgs = [req.query.uagrid]
        var result = yield cardchainservice.query("cardchaincc", "queryUAgrByUAgrID", queryArgs)
        res.send("successCallback(" + result.toString() + ")")
    }).catch((err) => {
        res.send(err)
    })
})

// 8. makeTranscation 0-{CardID},1-{TxTime},2-{TxContent},3-{TxTotalPrice},4-{TxIsUseAgreement},5-{TxIsUsePoint},6-{UserAgreementID},7-{AgreementID},8-{NumOfPointsUsed},9-{TxLocation},10-{TxOrg}

// charge "充100得120"
// ["OrgBarberMSP1234567890", "20180902151200", "charge", "100", "false", "false", "", "Agr_OrgBarberMSP108", "0", "Barber", "OrgBarberMSP"]
// http://localhost:3000/maketranscation?cardid=OrgBarberMSP1234567890&txtime=20180902151200&txcontent=charge&txtotalprice=100&txisuseagreement=false&txisusepoint=false&useragreementid=&agreementid=Agr_OrgBarberMSP108&numofpointsused=0&txlocation=Barber&txorg=OrgBarberMSP

// charge "100元享受4次原价30元洗头"
// ["OrgBarberMSP1234567890", "20180902151200", "charge", "0", "false", "false", "", "Agr_OrgBarberMSP109", "0", "Barber", "OrgBarberMSP"]
// http://localhost:3000/maketranscation?cardid=OrgBarberMSP1234567890&txtime=20180902151200&txcontent=charge&txtotalprice=100&txisuseagreement=false&txisusepoint=false&useragreementid=&agreementid=Agr_OrgBarberMSP109&numofpointsused=0&txlocation=Barber&txorg=OrgBarberMSP

// consume use agreement
// ["OrgBarberMSP1234567890", "20180902151200", "cut hair", "30", "true", "false", "UAgr_OrgBarberMSP1234567890M102", "", "0", "Barber", "OrgBarberMSP"]
// http://localhost:3000/maketranscation?cardid=OrgBarberMSP1234567890&txtime=20180902151200&txcontent=cut hair&txtotalprice=30&txisuseagreement=true&txisusepoint=false&useragreementid=UAgr_OrgBarberMSP1234567890M102&agreementid=&numofpointsused=0&txlocation=Barber&txorg=OrgBarberMSP

// consume use point
// ["OrgBarberMSP1234567890", "20180902151200", "cut hair", "30", "false", "false", "", "", "10", "Barber", "OrgBarberMSP"]
// http://localhost:3000/maketranscation?cardid=OrgBarberMSP1234567890&txtime=20180902151200&txcontent=cut hair&txtotalprice=30&txisuseagreement=false&txisusepoint=false&useragreementid=UAgr_OrgBarberMSP1234567890M102&agreementid=&numofpointsused=10&txlocation=Barber&txorg=OrgBarberMSP

app.get('/maketranscation', function(req, res){
    co (function *(){
        var transcationArgs = [req.query.cardid, req.query.txtime, req.query.txcontent, req.query.txtotalprice, req.query.txisuseagreement, req.query.txisusepoint, req.query.useragreementid, req.query.agreementid, req.query.numofpointsused, req.query.txlocation, req.query.txorg]
        var result = yield cardchainservice.invoke("cardchaincc", "makeTranscation", transcationArgs)
        res.send("successCallback()")
    }).catch((err) => {
        res.send(err)
    })
})

// 9. queryTranscationByCardID args: 0 - {CardID}
// http://localhost:3000/querytranscationbycardid?cardid=OrgBarberMSP1234567890
app.get('/querytranscationbycardid', function(req, res){
    co (function *(){
        var queryArgs = [req.query.cardid]
        var result = yield cardchainservice.query("cardchaincc", "queryTranscationByCardID", queryArgs)
        res.send("successCallback(" + result.toString() + ")")
    }).catch((err) => {
        res.send(err)
    })
})

function getMockupUserInfo(userName) {
    var mockupUser = {
        "username": "",
        "name": "",
        "passwd": "",
        "cmID": "",
        "Acct": ""
    };
    for (var user of config.certificateAuthorities.registrar) {
        if (user.username == userName) {
            mockupUser = user;
            break;
        }
    }
    return mockupUser;
}

// manage token
function updateToken(username, token) {
    console.log(JSON.stringify(process.TOKENS));
    try {
        console.log(token);
        process.TOKENS[username] = token;
        return token;
    } catch(err) {
        return null;
    }
}

// 10. user args: 0-{username},1-{orgName}
// http://localhost:3000/users?username=barberuser1&orgName=OrgBarberMSP
// Register and enroll user
app.get('/users', function(req, res) {
    var username = req.query.username;
    var orgName = req.query.orgName;
    if (!username) {
        res.json(getErrorMessage('\'username\''));
        return;
    }
    if (!orgName) {
        res.json(getErrorMessage('\'orgName\''));
        return;
    }
    var token = jwt.sign({
        exp: Math.floor(Date.now() / 1000) + 36000,
        username: username,
        orgName: orgName
    }, app.get('secret'));
    helper.getRegisteredUsers(username, orgName, true).then(function(response) {
        if (response && typeof response !== 'string') {
            updateToken(username, token);
            response.token = token;
            res.json(response);
        } else {
            res.json({
                success: false,
                message: response
            });
        }
    });
});

// 11. login user
app.get('/login', function(req, res) {
    var username = req.query.username;
    var orgName = req.query.orgName;
    var password = req.query.password;

    var mockupUser = null;
    for (var user of config['certificateAuthorities']['ca1.barbermsp.aliyunbaas.com']['registrar']) {
        if (user.enrollId == username) {
            mockupUser = user;
            break;
        }
    }
    if (mockupUser == null) {
        res.json({
            success: false,
            message: "用户不存在"
        });
    }

    new Promise(function (resolve, reject) {
        let oldtoken = process.TOKENS[username];

        if (oldtoken != null) {

            jwt.verify(oldtoken, app.get('secret'), function(err, decoded) {
                if (decoded.enrollSecret != password) {
                    reject({});
                }
                resolve({});
            });
        }else {
            resolve({});
        }
    }).then(function () {
        console.log("log:"+ app.get('secret'));
        var token = jwt.sign({
            exp: Math.floor(Date.now() / 1000) + 36000,
            username: username,
            orgName: orgName,
            password: password
        }, app.get('secret'));
        helper.getRegisteredUsers(username, orgName, true).then(function(response) {
            if (response && typeof response !== 'string') {
                updateToken(username, token);

                var mockupUser = getMockupUserInfo(username)
                res.status(200).json({
                    success : true,
                    secret : response.secret,
                    message : response.message,
                    token : token,
                    user : mockupUser
                });
            } else {
                res.json({
                    success: false,
                    message: response
                });
            }
        });

    }, function () {
        res.json({
            success: false,
            message: 'password incorrect, please retry'
        });
        return;
    });

});

var server = app.listen(3000, function(){
    var host = server.address().address
    var port = server.address().port

    console.log('App listening at http://%s:%s', host, port)
})

process.on('uhandleRejection', function(err){
    console.error(err.stack)
})

process.on('uncaughtException', console.error)
