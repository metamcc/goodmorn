package constants

/*
const SettleDADLog = "date~keyIndex~wallet~fruit~dad"

const IndexUserID = "idxtype~userid~wallet~date" //index for find user wallet by userid

const IndexUnqWalID = "idxtype~ecowltid"

const IndexTokenWalID = "idxtype~tokenwltid"










//const GethHost = "http://192.168.1.67:8545"
const GethHost = "http://192.168.1.98:8545"
*/


// -------------------------------------------------------------------------------

const IndexUserID = "idxtype~userid~wallet~date" //index for find user wallet by userid
const EcoWallet = "eco"
const UserWallet = "user"
const SystemWallet = "system"
const DADLog_CHAINCODE = "dadlog"
const MCCChannelID = "mycreditchain"
const Config_CHAINCODE = "confcc"




const IndexCoreCheck = "YY~MM~DD~svcCD~refID~txID" // 코어 체크 인덱스

const CORE_SEED_PENALTY_YN = "CORE_SEED_PENALTY_YN" // 코어 체크 여부

const DAD_LIMIT_AMNT = "DAD_LIMIT_AMNT"				// DAD 정산 량

const SEED_LIMIT_CNT = "SEED_LIMIT_CNT"				// 씨앗 전송 횟수

const SND_DATA_SYNC_CNT = "SND_DATA_SYNC_CNT"				// 결제 동기화 기준 갯수

const RCV_DATA_SYNC_CNT = "RCV_DATA_SYNC_CNT"				// 결제 동기화 기준 갯수 - 받는 지갑용

const IndexFruitPayment = "myWallet~YY~MM~DD~refID~refWallet~amount~txID" // 열매 결제 인덱스

//PAYMENT_DATA_SYNC_CNT