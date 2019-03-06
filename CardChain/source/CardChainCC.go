package main

import (
    "fmt"
    "strconv"
    "bytes"
    "strings"

    pb "github.com/hyperledger/fabric/protos/peer"
    "github.com/hyperledger/fabric/core/chaincode/shim"
    // "github.com/hyperledger/fabric/common/flogging"
    "encoding/json"
    // "time"
)

// 卡记录前缀
const CARD_PREFIX = "Card_"

// 合约记录前缀
const AGREEMENT_PREFIX = "Agr_"

// 用户参与的合约记录前缀
const UAGREEMENT_PREFIX = "UAgr_"

// 交易记录前缀
const TX_PREFIX = "Tx_"

const INDEX_MAX = "999"
const INDEX_MIN = "100"

// 卡片结构
type Card struct{
    // Card information
    CardID          string  `json:CardID`           // 卡号
    CardBalance     int     `json:Balance`          // 余额
    CardPoints      int     `json:Points`           // 积分
    CardPointsBase  int     `json:CardPointsBase`   // 根据消费金额获得积分的基数
    CardCreateTime  string  `json:CardCreateTime`   // 开卡时间
    CardState       int     `json:State`            // 卡状态 1-active

    // Customer information
    CustomerID      string  `json:CustomerID`       // 顾客ID:身份证号
    CustomerGender  string  `json:CustomerGender`   // 性别
    CustomerAge     int     `json:CustomerAge`      // 年龄
    CustomerTEL     string  `json:CustomerTEL`      // 电话
}

// 合约
type Agreement struct{
    AgreementID             string  `json:AgreementID`          // 合约ID
    AgreementOrg            string  `json:AgreementOrg`         // 发布优惠活动的商家
    AgreementPrice          string  `json:AgreementPrice`       // 优惠价格
    AgreementCount          int     `json:AgreementCount`       // 优惠次数
    AgreementCreateTime     string  `json:AgreementCreateTime`  // 创建时间
    AgreementDeadline       string  `json:AgreementDeadline`    // 优惠截至时间。格式："20180101180000"，2018年01月01日18点00分00秒失效
    AgreementDescription    string  `json:AgreementDescription` // 描述
}

// 用户参与的合约
type UserAgreement struct{
    UAgrID          string  `json:UserAgrID`
    UAgrRemainCount int     `json:UAgrRemainCount`  // 合约剩余优惠次数
    UAgrRemainMoney int     `json:UAgrRemainMoney`  // 合约剩余优惠金额
    AgreementID     string  `json:AgreementID`
    CardID          string  `json:CardID`
}

// 交易
type Transaction struct{
    TxID                string  `json:TxID`
    TxTime              string  `json:TxTime`               // 消费时间
    TxContent           string  `json:TxContent`            // 消费内容
    TxTotalPrice        int     `json:TxTotalPrice`         // 消费金额
    TxIsUseAgreement    bool    `json:TxIsUseAgreement`     // 是否使用优惠次数
    TxIsUsePoint        bool    `json:TxIsUsePoint`         // 是否使用积分
    UAgreementID        string  `json:UAgreementID`
    CardID              string  `json:CardID`
    AgreementID         string  `json:AgreementID`          // 充值时使用的优惠
    NumOfPointsUsed     int     `json:NumOfPointsUsed`      // 使用积分的数量
    TxLocation          string  `json:TxLocation`           // 消费地点
    TxOrg               string  `json:TxOrg`
}

// CardChain
type CardChain struct {
}

// 发卡
// args: 0-{CardID},1-{CardBalance},2-{CardPoints},3-{CardPointsBase},4-{CardCreateTime},5-{CardState},6-{CustomerID},7-{CustomerGender},8-{CustomerAge},9-{CustomerTEL}
func (s *CardChain) issueCard(stub shim.ChaincodeStubInterface, args []string) pb.Response{
    if len(args) != 10 {
        return shim.Error("Cardchain Invoke issueCard args != 10");
    }

    balance, err := strconv.Atoi(args[1])
    points, err := strconv.Atoi(args[2])
    pointsBase, err := strconv.Atoi(args[3])
    state, err := strconv.Atoi(args[5])
    age, err := strconv.Atoi(args[8])

    // 根据卡ID查询是否已有该卡片，如果有则返回
    cardExist, err := stub.GetState(CARD_PREFIX + args[0])
    if err != nil {
        return shim.Error(err.Error())
    }
    if cardExist != nil {
        return shim.Error("The card is Exist")
    }

    card := Card{CARD_PREFIX + args[0], balance, points, pointsBase, args[4], state, args[6], args[7], age, args[9]}

    // 保存卡片
    c, _ := json.Marshal(card)
    err = stub.PutState(card.CardID, c)
    if err != nil {
        return shim.Error("CardChain issueCard putCard failed!")
    }

    return shim.Success(nil)
}

// 根据卡ID查询卡
// args: 0 - {CardID}
func (s *CardChain) queryCardByCardID(stub shim.ChaincodeStubInterface, args []string) pb.Response{
    if len(args) != 1{
        return shim.Error("CardChain queryCard args != 1")
    }

    cardByte, err := stub.GetState(CARD_PREFIX + args[0])
    if err != nil {
        return shim.Error("CardChain queryCard get card failed")
    }
    if cardByte == nil {
        return shim.Error("Card not Exist")
    }

    return shim.Success(cardByte)
}

// 新增合约
// args: 0-{OrgID},1-{AgreementOrg},2-{AgreementPrice},3-{AgreementCount},4-{AgreementCreateTime},5-{AgreementDeadline},6-{AgreementDescription}
func (s *CardChain) makeAgreement(stub shim.ChaincodeStubInterface, args []string) pb.Response{
    if len(args) != 7 {
        return shim.Error("Cardchain makeAgreement args != 7");
    }

    agreementCount, err := strconv.Atoi(args[3])

    agreementID := AGREEMENT_PREFIX + args[0]

    // 查找递增的数字
    startIndex := agreementID + INDEX_MIN
    endIndex := agreementID + INDEX_MAX
    resultsIterator, err := stub.GetStateByRange(startIndex, endIndex)
    if err != nil {
        return shim.Error("CardChain makeAgreement GetStateByRange failed")
    }
    defer resultsIterator.Close()

    count := 100
    for resultsIterator.HasNext() {
        _, err := resultsIterator.Next()
        if err != nil {
            shim.Error("CardChain queryUAgrByRange resultsIterator Next failed")
        }
        count ++
    }

    // 记录"AgreementLedger"账本
    agreement := Agreement{agreementID + strconv.Itoa(count), args[1], args[2], agreementCount, args[4], args[5], args[6]}
    a, _ := json.Marshal(agreement)
    err = stub.PutState(agreement.AgreementID, a)
    if err != nil {
        return shim.Error("CardChain makeAgreement PutState failed")
    }

    return shim.Success(nil)
}

// 查询商家所有的合约
// args:0-{OrgID}
func (s *CardChain) queryAgreementByRange(stub shim.ChaincodeStubInterface, args []string) pb.Response{
    if len(args) != 1 {
        return shim.Error("Cardchain queryAgreementByRange args != 1");
    }

    startIndex := AGREEMENT_PREFIX + args[0] + INDEX_MIN
    endIndex := AGREEMENT_PREFIX + args[0] + INDEX_MAX
    resultsIterator, _ := stub.GetStateByRange(startIndex, endIndex)
    defer resultsIterator.Close()

    var buffer bytes.Buffer
    buffer.WriteString("[")

    bArrayMemberAlreadyWritten := false
    count := 0
    for resultsIterator.HasNext() {
        agreement, err := resultsIterator.Next()
        if err != nil {
            shim.Error("CardChain queryAgreementByRange resultsIterator Next failed")
        }
        count ++
        if bArrayMemberAlreadyWritten == true {
            buffer.WriteString(",")
        }
        buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(agreement.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(agreement.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
    }
    buffer.WriteString("]")
    return shim.Success(buffer.Bytes())
}

// 查询已有的合约
// args: 0-{AgreementID}
func (s *CardChain) queryAgreementByAgreementID(stub shim.ChaincodeStubInterface, args []string) pb.Response{
    if len(args) != 1 {
        return shim.Error("Cardchain queryAgreementByOrgID args != 1");
    }
    agreementByte, err := stub.GetState(args[0])
    if err != nil {
        return shim.Error("CardChain queryAgreementByOrgID GetState failed")
    }
    if agreementByte == nil {
        return shim.Error("Agreement not Exist")
    }
    return shim.Success(agreementByte)
}

// 查询用户参加的所有合约
// agrs:0-{CardID}
func (s *CardChain) queryUAgrByRange(stub shim.ChaincodeStubInterface, args []string) pb.Response{
    if len(args) != 1 {
        return shim.Error("Cardchain queryUAgrByRange args != 1");
    }

    startIndex := UAGREEMENT_PREFIX + args[0] + "M" + INDEX_MIN
    endIndex := UAGREEMENT_PREFIX + args[0] + "M" + INDEX_MAX
    resultsIterator, _ := stub.GetStateByRange(startIndex, endIndex)
    defer resultsIterator.Close()

    var buffer bytes.Buffer
    buffer.WriteString("[")

    bArrayMemberAlreadyWritten := false
    for resultsIterator.HasNext() {
        agreement, err := resultsIterator.Next()
        if err != nil {
            shim.Error("CardChain queryUAgrByRange resultsIterator Next failed")
        }
        if bArrayMemberAlreadyWritten == true {
            buffer.WriteString(",")
        }
        buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(agreement.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(agreement.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
    }
    buffer.WriteString("]")
    return shim.Success(buffer.Bytes())
}

// 查询用户参加的指定合约
// args: 0 - {UAgrID}
func (s *CardChain) queryUAgrByUAgrID(stub shim.ChaincodeStubInterface, args []string) pb.Response{
    if len(args) != 1{
        return shim.Error("CardChain queryUAgrByUAgrID args != 1")
    }
    uAgrUsingBytes, err := stub.GetState(args[0])
    if err != nil {
        return shim.Error("CardChain queryUAgrByUAgrID GetState failed")
    }
    if uAgrUsingBytes == nil {
        return shim.Error("UAgrID not Exist")
    }
    return shim.Success(uAgrUsingBytes)
}

// 交易
// args:
// 0-{CardID},1-{TxTime},2-{TxContent},3-{TxTotalPrice},4-{TxIsUseAgreement},5-{TxIsUsePoint},6-{UserAgreementID},7-{AgreementID},8-{NumOfPointsUsed},9-{TxLocation},10-{TxOrg}
func (s *CardChain) makeTranscation(stub shim.ChaincodeStubInterface, args []string) pb.Response{
    if len(args) != 11{
        return shim.Error("CardChain makeTranscation args != 11")
    }

    cardBytes, err := stub.GetState(CARD_PREFIX + args[0])
    if err != nil {
        return shim.Error("CardChain makeTranscation get card " + args[0] + " failed")
    }

    card := Card{}
    err = json.Unmarshal(cardBytes, &card)
    if err != nil {
        return shim.Error("CardChain makeTranscation Unmarshal card failed")
    }


    totalPrice, err := strconv.Atoi(args[3])
    isUseAgreement, err := strconv.ParseBool(args[4])
    isUsePoint, err := strconv.ParseBool(args[5])
    numOfPointsUsed, err := strconv.Atoi(args[8])

    // 当充值时
    if args[2] == "charge" {
        fmt.Println("Charge")

        // 记录用户参与优惠
        agreementBytes, err := stub.GetState(args[7])
        if err != nil{
            return shim.Error("CardChain makeTranscation GetState when 'charge' failed")
        }

        agreement := Agreement{}
        err = json.Unmarshal(agreementBytes, &agreement)
        if err != nil {
            return shim.Error("CardChain makeTranscation Unmarshal when 'charge' failed")
        }

        // "消费预储值"类充值记录到"UserAgreementLedger"账本
        uAgrID := UAGREEMENT_PREFIX + args[0]
        uAgrRemainCount := agreement.AgreementCount
        uAgrRemainMoney := 0
        priceArray := strings.Split(agreement.AgreementPrice, ",")
        chargeValue, err := strconv.Atoi(priceArray[1])
        uAgrRemainMoney += chargeValue
        agreementID := agreement.AgreementID
        CardID := card.CardID

        startIndex := uAgrID + "M" + INDEX_MIN
        endIndex := uAgrID + "M" + INDEX_MAX
        resultsIterator, err := stub.GetStateByRange(startIndex, endIndex)
        if err != nil {
            return shim.Error("CardChain makeAgreement GetStateByRange failed")
        }
        defer resultsIterator.Close()

        count := 100
        for resultsIterator.HasNext() {
            _, err := resultsIterator.Next()
            if err != nil {
                shim.Error("CardChain makeTranscation resultsIterator Next failed")
            }
            count ++
        }

        uAgr := UserAgreement{uAgrID + "M" + strconv.Itoa(count), uAgrRemainCount, uAgrRemainMoney, agreementID, CardID}
        uAgrBytes, err := json.Marshal(uAgr)
        err = stub.PutState(uAgr.UAgrID, uAgrBytes)

        fmt.Println("uAgr")
        fmt.Println(uAgr)
        if err != nil {
            return shim.Error("CardChain makeTranscation PutState when 'charge' failed")
        }

        // "充100得120" 类直接充值到账户
        if uAgrRemainCount <= 0 {
            card.CardBalance += uAgrRemainMoney
        }
        fmt.Println(card)
    } else if isUseAgreement {
        fmt.Println("Use Agreement")

        // 使用“消费预储值”优惠时
        uAgrUsingBytes, err := stub.GetState(args[6])
        if err != nil {
            return shim.Error("CardChain makeTranscation GetState when 'useAgreement' failed")
        }
        uAgrUsing := UserAgreement{}
        err = json.Unmarshal(uAgrUsingBytes, &uAgrUsing)
        if err != nil {
            return shim.Error("CardChain makeTranscation Unmarshal when 'useAgreement' failed")
        }
        fmt.Println(uAgrUsing)
        uAgrUsing.UAgrRemainCount -= 1
        if uAgrUsing.UAgrRemainCount < 0 {
            return shim.Error("The Agreement has no count")
        }
        uAgrUsing.UAgrRemainMoney = 0
        uAgrUsingBytes, err = json.Marshal(uAgrUsing)
        fmt.Println(uAgrUsing.UAgrID)
        err = stub.PutState(uAgrUsing.UAgrID, uAgrUsingBytes)
        if err != nil {
            return shim.Error("CardChain makeTranscation PutState when 'useAgreement' failed")
        }
    } else if isUsePoint {
        fmt.Println("Use Points")

        // 使用积分时
        card.CardPoints = card.CardPoints - numOfPointsUsed + (totalPrice - numOfPointsUsed) / card.CardPointsBase
    } else {
        fmt.Println("use money")
        card.CardBalance -= totalPrice
        card.CardPoints += totalPrice / card.CardPointsBase
        if card.CardBalance < 0 {
            return shim.Error("Balance is not enough")
        }
    }

    // 保存卡账本
    cardBytes, err = json.Marshal(card)
    err = stub.PutState(card.CardID, cardBytes)
    fmt.Println(card.CardID)
    if err != nil {
        return shim.Error("CardChain makeTranscation PutState card failed")
    }

    // 保存交易账本
    tx := Transaction{TX_PREFIX + args[0], args[1], args[2], totalPrice, isUseAgreement, isUsePoint, args[6], args[0], args[7], numOfPointsUsed, args[9], args[10]}
    txBytes, err := json.Marshal(tx)
    err = stub.PutState(tx.TxID, txBytes)
    fmt.Println(tx.TxID)
    if err != nil {
        return shim.Error("CardChain makeAgreement PutState tx failed")
    }

    return shim.Success(nil)
}

// 按CardID查询交易历史
// args: 0 - {CardID}
func (s *CardChain) queryTranscationByCardID(stub shim.ChaincodeStubInterface, args []string) pb.Response{
    if len(args) != 1{
        return shim.Error("CardChain queryTranscationByCardID args != 1")
    }

    resultsIterator, err := stub.GetHistoryForKey(TX_PREFIX + args[0])
    if err != nil {
        return shim.Error("CardChain queryTranscationByCardID GetHistoryForKey failed")
    }
    defer resultsIterator.Close()

    var buffer bytes.Buffer
    buffer.WriteString("[")

    bArrayMemberAlreadyWritten := false
    for resultsIterator.HasNext() {
        tx, err := resultsIterator.Next()
        if err != nil {
            shim.Error("CardChain queryUAgrByRange resultsIterator Next failed")
        }
        if bArrayMemberAlreadyWritten == true {
            buffer.WriteString(",")
        }
        buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(tx.TxId)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(tx.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
    }
    buffer.WriteString("]")
    return shim.Success(buffer.Bytes())
}

// 按OrgID查询交易
// args: 0-{OrgID}
func (s *CardChain) queryTranscationByOrgID(stub shim.ChaincodeStubInterface, args []string) pb.Response{
    return shim.Success(nil)
}

func (s *CardChain) Init(stub shim.ChaincodeStubInterface) pb.Response{
    return shim.Success(nil)
}

func (s *CardChain) Invoke(stub shim.ChaincodeStubInterface) pb.Response{
    function, args := stub.GetFunctionAndParameters()

    if function == "init" {
        return s.Init(stub)
    } else if function == "issueCard" {
        return s.issueCard(stub, args)
    } else if function == "queryCardByCardID" {
        return s.queryCardByCardID(stub, args)
    } else if function == "makeAgreement" {
        return s.makeAgreement(stub, args)
    } else if function == "queryAgreementByRange" {
        return s.queryAgreementByRange(stub, args)
    } else if function == "queryAgreementByAgreementID" {
        return s.queryAgreementByAgreementID(stub, args)
    } else if function == "queryUAgrByRange" {
        return s.queryUAgrByRange(stub, args)
    } else if function == "queryUAgrByUAgrID" {
        return s.queryUAgrByUAgrID(stub, args)
    } else if function == "makeTranscation" {
        return s.makeTranscation(stub, args)
    } else if function == "queryTranscationByCardID" {
        return s.queryTranscationByCardID(stub, args)
    } else if function == "queryTranscationByOrgID" {
        return s.queryTranscationByOrgID(stub, args)
    }

    return shim.Error("CardChain Unkown method!")
}

func main(){
    if err := shim.Start(new(CardChain)); err != nil {
        fmt.Printf("Error starting CardChain: %s", err)
    }
}
