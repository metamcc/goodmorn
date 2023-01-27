package io.mcc.mcctoken.service;

import io.mcc.common.MessageByLocaleService;
import io.mcc.common.entity.EthMccTxLog;
import io.mcc.common.entity.EthTransaction;
import io.mcc.common.util.RSATools;
import io.mcc.mcctoken.support.*;
import lombok.extern.slf4j.Slf4j;
import org.jasypt.encryption.StringEncryptor;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Qualifier;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.DependsOn;
import org.springframework.core.io.ClassPathResource;
import org.springframework.stereotype.Service;
import org.web3j.crypto.Credentials;
import org.web3j.crypto.RawTransaction;
import org.web3j.crypto.TransactionEncoder;
import org.web3j.crypto.WalletUtils;
import org.web3j.protocol.Web3j;
import org.web3j.protocol.core.DefaultBlockParameter;
import org.web3j.protocol.core.DefaultBlockParameterName;
import org.web3j.protocol.core.methods.response.*;
import org.web3j.tx.ClientTransactionManager;
import org.web3j.tx.Transfer;
import org.web3j.tx.gas.ContractGasProvider;
import org.web3j.tx.gas.StaticGasProvider;
import org.web3j.utils.Convert;
import org.web3j.utils.Numeric;
import rx.Subscription;

import java.io.File;
import java.io.IOException;
import java.math.BigDecimal;
import java.math.BigInteger;
import java.security.KeyPair;
import java.time.Instant;
import java.time.LocalDateTime;
import java.time.ZoneId;
import java.util.*;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.CountDownLatch;
import java.util.concurrent.TimeUnit;

/**
 * Sample service to demonstrate web3j bean being correctly injected.
 */
@Service
@Slf4j
public class Web3jService {

    @Autowired
    private Web3j web3j;

    @Autowired
    @Qualifier("MCC")
    private Web3j web3jMCC;

    @Autowired
    @Qualifier("encryptorBean")
    private StringEncryptor stringEncryptor;

    @Autowired
    private MessageByLocaleService messageByLocaleService;

    @Value("${mcc.contract.token.addr:@null}")
    private String mccTokenContractAddress;

    @Value("${mcc.contract.distrib.addr:@null}")
    private String mccDistribContractAddress;

    @Value("${mcc.contract.broker.addr:@null}")
    private String mccBrokerContractAddress;

    @Value("#{'${mcc.contract.broker.owners:@null}'.split(',')}")
    private String[] mccBrokerContractOwners;

    @Value("${mcc.eth.target.network}")
    private String targetNetwork;

    @Value("${mcc.keystore.pw}")
    private String keystorePW;

    @Value("${mcc.query.wallet.pw}")
    private String queryWalletFilePW;

    @Value("${mcc.query.wallet.addr}")
    private String queryWalletAddr;

    //ECO Wallet
    @Value("${mcc.eco.wallet.addr:@null}")
    private String ecoWalletAddr;

    @Value("${mcc.eco.wallet.pw:@null}")
    private String ecoWalletAddrPWEnc;

    @Value("${mcc.payreserv.wallet.addr:@null}")
    private String payReservWalletAddr;

    @Value("${mcc.payreserv.wallet.pw:@null}")
    private String payReservWalletAddrPWEnc;

    public String isWalletOwner(File walletFile, final String walletPassword) throws Exception {

        Credentials credentials = WalletUtils.loadCredentials(walletPassword, walletFile);

        return credentials.getAddress();
    }

    public BigInteger getBlockNumber() throws Exception {
        return web3j.ethBlockNumber().send().getBlockNumber();
    }
    //=============
    private MCCToken getMCCToken(Web3j _web3j, String _targetNetwork, Credentials credentials, ContractGasProvider contractGasProvider) {
        if ("mainnet".equalsIgnoreCase(_targetNetwork))
            return MCCTokenMain.load(mccTokenContractAddress, _web3j, credentials, contractGasProvider);
        else if ("rinkeby".equalsIgnoreCase(_targetNetwork))
            return MCCTokenRinkeby.load(mccTokenContractAddress, _web3j, credentials, contractGasProvider);

        return MCCTokenLocal.load(mccTokenContractAddress, _web3j, credentials, contractGasProvider);
    }


    private MCCToken getMCCToken(Credentials credentials, ContractGasProvider contractGasProvider) {
        return getMCCToken(web3j, targetNetwork, credentials, contractGasProvider);
    }

    private MCCToken getMCCToken(Credentials credentials, BigInteger GAS_PRICE, BigInteger GAS_LIMIT) {
        return getMCCToken(credentials, new StaticGasProvider(GAS_PRICE, GAS_LIMIT));
    }


    public String getClientVersion() throws IOException {
        Web3ClientVersion web3ClientVersion = web3j.web3ClientVersion().send();
        return web3ClientVersion.getWeb3ClientVersion();
    }

    public BigInteger balanceOfETH(String address) throws Exception {

        EthGetBalance web3ClientVersion = web3j.ethGetBalance(address, DefaultBlockParameter.valueOf("latest")).send();
        return web3ClientVersion.getBalance();
    }

    public TransactionReceipt transferETH(String senderPrivateKey, String toAddr, BigInteger ethValue) throws Exception {
        Credentials credentials = Credentials.create(senderPrivateKey);

        BigInteger GAS_PRICE = new BigInteger("5000000000");
        BigInteger GAS_LIMIT = new BigInteger("2000000");

        return transferETH(credentials, toAddr, ethValue, GAS_PRICE, GAS_LIMIT);
    }
    public TransactionReceipt transferETH(File ownerWalletFile, String ownerPW, String toAddr, BigInteger ethValue) throws Exception {
        //Credentials credentials = Credentials.create(ownerPVKey);
        Credentials credentials = WalletUtils.loadCredentials(ownerPW, ownerWalletFile);

        BigInteger GAS_PRICE = new BigInteger("5000000000");
        BigInteger GAS_LIMIT = new BigInteger("2000000");

        return transferETH(credentials, toAddr, ethValue, GAS_PRICE, GAS_LIMIT);
    }


    public TransactionReceipt transferETH(Credentials credentials, String toAddr, BigInteger ethValue, BigInteger GAS_PRICE, BigInteger GAS_LIMIT) throws Exception {
        log.trace("toAddr {}, mccValue {}", toAddr, ethValue);

//        EthGetTransactionCount ethGetTransactionCount = web3j.ethGetTransactionCount(
//                toAddr, DefaultBlockParameterName.LATEST).sendAsync().get();
//        BigInteger nonce = ethGetTransactionCount.getTransactionCount();
//
//        RawTransaction rawTransaction  = RawTransaction.createEtherTransaction(
//                nonce, GAS_PRICE, GAS_LIMIT, toAddr, ethValue);
//        byte[] signedMessage = TransactionEncoder.signMessage(rawTransaction, credentials);
//        String hexValue = Numeric.toHexString(signedMessage);
//        EthSendTransaction ethSendTransaction = web3j.ethSendRawTransaction(hexValue).sendAsync().get();
//        String transactionHash = ethSendTransaction.getTransactionHash();

        TransactionReceipt transactionReceipt = Transfer.sendFunds(web3j, credentials, toAddr, new BigDecimal(ethValue.toString()), Convert.Unit.WEI).send();

        //BigInteger txFees = txReceipt.getGasUsed().multiply(Web3jConstants.GAS_PRICE);

        return transactionReceipt;

    }

    public HashMap<String, Object> getMCCTokenInfo() throws Exception {
        String queryWalletFile = messageByLocaleService.getMessage("mcc.wallet." + queryWalletAddr + ".path");

        Credentials credentials = WalletUtils.loadCredentials(queryWalletFilePW, queryWalletFile);

        BigInteger GAS_PRICE = new BigInteger("5000000000");
        BigInteger GAS_LIMIT = new BigInteger("400000");

        MCCToken mccToken = getMCCToken(credentials, GAS_PRICE, GAS_LIMIT);

        if (!mccToken.isValid()) {
            //TODO..
            throw new Exception("Contracrt is Invalid");
        }

        HashMap<String, Object> mapInfo = new HashMap<>();
        mapInfo.put("address", mccToken.getContractAddress());

        BigInteger totalToken = mccToken.totalSupply().send();

        mapInfo.put("totalSupply", totalToken.toString(10));
        mapInfo.put("decimals", mccToken.decimals().send());
        mapInfo.put("name", mccToken.name().send());
        mapInfo.put("allowTransfers", mccToken.allowTransfers().send());
        mapInfo.put("isLimitEnabled", mccToken.isLimitEnabled().send());
        mapInfo.put("issuanceFinished", mccToken.issuanceFinished().send());
        mapInfo.put("eventListener", mccToken.eventListener().send());
        mapInfo.put("manager", mccToken.manager().send());
        mapInfo.put("owners", mccToken.getOwners().send());

        return mapInfo;
    }

    public void setMCCTokenOwners(String managerWalletFilePath, String managerWalletPW, String newOwnerWallets) throws Exception {

        Credentials credentials = WalletUtils.loadCredentials(managerWalletPW, managerWalletFilePath);

        BigInteger GAS_PRICE = new BigInteger("5000000000");
        BigInteger GAS_LIMIT = new BigInteger("400000");

        MCCToken mccToken = getMCCToken(credentials, new StaticGasProvider(GAS_PRICE, GAS_LIMIT));

        if (!mccToken.isValid()) {
            throw new Exception("Contracrt is Invalid");
            //LOGGER.debug("Contract is Invalid:{}", mccDistrib);
        }

        ArrayList<String> owerWallets = new ArrayList<String>();

        String[] owners = newOwnerWallets.split(",");
        for (String wallet : owners) {
            if (wallet != null && wallet.length() > 10)
                owerWallets.add(wallet.trim());
        }
        log.trace("setDistribOwner==>{}", owerWallets);
        TransactionReceipt trReceipt = mccToken.setOwners(owerWallets).send();

        log.trace("setDistribOwner==>{}", trReceipt);

    }

    public Boolean setMCCTokenAllowTrans(final String walletFile, final String walletPW, Boolean bAllow) throws Exception {
        Credentials credentials = WalletUtils.loadCredentials(walletPW, walletFile);


        BigInteger GAS_PRICE = new BigInteger("5000000000");
        BigInteger GAS_LIMIT = new BigInteger("400000");

        MCCToken mccToken = getMCCToken(credentials, GAS_PRICE, GAS_LIMIT);

        if (!mccToken.isValid()) {
            //TODO..
            throw new Exception("Contracrt is Invalid");
        }
        TransactionReceipt transactionReceipt = mccToken.setAllowTransfers(bAllow).send();
        log.trace("{}", transactionReceipt);

        return mccToken.allowTransfers().send();
    }

//    public Boolean setMCCTokenFinishIssue(final String contractOwnerPW) throws Exception {
//        Credentials credentials = WalletUtils.loadCredentials(contractOwnerPW, contractOwnerWalletFile);
//
//
//        BigInteger GAS_PRICE = new BigInteger("20000000000");
//        BigInteger GAS_LIMIT = new BigInteger("400000");
//
//        MCCToken mccToken = getMCCToken(credentials, GAS_PRICE, GAS_LIMIT);
//
//        if (!mccToken.isValid()) {
//            //TODO..
//            throw  new Exception("Contracrt is Invalid");
//        }
//        TransactionReceipt transactionReceipt = mccToken.finishIssuance().send();
//        log.trace("{}", transactionReceipt);
//
//
//        return mccToken.allowTransfers().send();
//    }

    public Boolean setMCCTokenTransLimit(final String walletFile, final String walletPW, Boolean bOnOff) throws Exception {
        Credentials credentials = WalletUtils.loadCredentials(walletPW, walletFile);


        BigInteger GAS_PRICE = new BigInteger("5000000000");
        BigInteger GAS_LIMIT = new BigInteger("400000");

        MCCToken mccToken = getMCCToken(credentials, GAS_PRICE, GAS_LIMIT);

        if (!mccToken.isValid()) {
            //TODO..
            throw new Exception("Contracrt is Invalid");
        }
        TransactionReceipt transactionReceipt = mccToken.setLimitEnabled(bOnOff).send();
        log.trace("{}", transactionReceipt);


        return mccToken.isLimitEnabled().send();
    }

    public Boolean setMCCTokenTransLimitWallet(final String walletFile, final String walletPW, String target, boolean bRemove) throws Exception {
        Credentials credentials = WalletUtils.loadCredentials(walletPW, walletFile);

        BigInteger GAS_PRICE = new BigInteger("5000000000");
        BigInteger GAS_LIMIT = new BigInteger("400000");

        MCCToken mccToken = getMCCToken(credentials, GAS_PRICE, GAS_LIMIT);

        if (!mccToken.isValid()) {
            //TODO..
            throw new Exception("Contracrt is Invalid");
        }
        TransactionReceipt transactionReceipt = null;
        if (bRemove)
            transactionReceipt = mccToken.delLimitedWalletAddress(target).send();
        else
            transactionReceipt = mccToken.addLimitedWalletAddress(target).send();
        log.trace("{}", transactionReceipt);


        return mccToken.isLimitedWalletAddress(target).send();
    }

    public Boolean isMCCTokenTransLimitWallet(String wallet) throws Exception {
        String queryWalletFile = messageByLocaleService.getMessage("mcc.wallet." + queryWalletAddr + ".path");

        Credentials credentials = WalletUtils.loadCredentials(queryWalletFilePW, queryWalletFile);

        BigInteger GAS_PRICE = new BigInteger("5000000000");
        BigInteger GAS_LIMIT = new BigInteger("400000");

        MCCToken mccToken = getMCCToken(credentials, GAS_PRICE, GAS_LIMIT);

        if (!mccToken.isValid()) {
            //TODO..
            throw new Exception("Contracrt is Invalid");
        }

        return mccToken.isLimitedWalletAddress(wallet).send();
    }

    public BigInteger balanceOfMCC(String address) throws Exception {
        return balanceOfMCC(queryWalletAddr, queryWalletFilePW, address);
    }

    public BigInteger balanceOfMCC(final String walletAddr, final String walletPW, String address) throws Exception {
        String payWalletFile = messageByLocaleService.getMessage("mcc.wallet."+walletAddr+".path");

        Credentials credentials = WalletUtils.loadCredentials(walletPW, payWalletFile);

        return balanceOfMCC(credentials, address);
    }
    public BigInteger balanceOfMCC(final Credentials credentials, String address) throws Exception {

        BigInteger GAS_PRICE = new BigInteger("5000000000");
        BigInteger GAS_LIMIT = new BigInteger("400000");

        MCCToken mccToken = getMCCToken(credentials, new StaticGasProvider(GAS_PRICE, GAS_LIMIT));

        if (!mccToken.isValid()) {
            //TODO..
            throw  new Exception("Contracrt is Invalid");
        }

        BigInteger balnace = mccToken.balanceOf(address).send();

        return balnace;
    }

    public TransactionReceipt transferMCC(String senderPrivateKey, String toAddr, BigInteger mccValue) throws Exception {
        Credentials credentials = Credentials.create(senderPrivateKey);

        BigInteger GAS_PRICE = new BigInteger("5000000000");
        BigInteger GAS_LIMIT = new BigInteger("2000000");

        return transferMCC(credentials, toAddr, mccValue, GAS_PRICE, GAS_LIMIT);
    }

    public TransactionReceipt transferMCC(File ownerWalletFile, String ownerPW, String toAddr, BigInteger mccValue) throws Exception {
        //Credentials credentials = Credentials.create(ownerPVKey);
        Credentials credentials = WalletUtils.loadCredentials(ownerPW, ownerWalletFile);

        BigInteger GAS_PRICE = new BigInteger("5000000000");
        BigInteger GAS_LIMIT = new BigInteger("2000000");

        return transferMCC(credentials, toAddr, mccValue, GAS_PRICE, GAS_LIMIT);
    }

    public TransactionReceipt transferMCC(String walletAddr, String walletPW, String toAddr, BigInteger mccValue) throws Exception {

        String payWalletFile = messageByLocaleService.getMessage("mcc.wallet."+walletAddr+".path");
        String gasPrice = messageByLocaleService.getMessage("mcc.wallet."+walletAddr+".gasprice");
        String gasLimit = messageByLocaleService.getMessage("mcc.wallet."+walletAddr+".gaslimit");

        Credentials credentials = WalletUtils.loadCredentials(walletPW, payWalletFile);


        BigInteger GAS_PRICE = new BigInteger(gasPrice);
        BigInteger GAS_LIMIT = new BigInteger(gasLimit);

        return transferMCC(credentials, toAddr, mccValue, GAS_PRICE, GAS_LIMIT);
    }


    public TransactionReceipt transferMCC(Credentials credentials, String toAddr, BigInteger mccValue, BigInteger GAS_PRICE, BigInteger GAS_LIMIT) throws Exception {
        MCCToken mccToken = getMCCToken(credentials, new StaticGasProvider(GAS_PRICE, GAS_LIMIT));
        log.trace("toAddr {}, mccValue {}", toAddr, mccValue);

        if (!mccToken.isValid()) {
            //TODO..
            throw  new Exception("Contracrt is Invalid");
        }

        TransactionReceipt transactionReceipt = mccToken.transfer(toAddr, mccValue).sendAsync().get();
        log.trace("{}", transactionReceipt);

        return transactionReceipt;
    }

    public CompletableFuture<TransactionReceipt> transferMCCAsync(String walletAddr, String walletPW, String toAddr, BigInteger mccValue) throws Exception {

        String payWalletFile = messageByLocaleService.getMessage("mcc.wallet."+walletAddr+".path");
        String gasPrice = messageByLocaleService.getMessage("mcc.wallet."+walletAddr+".gasprice");
        String gasLimit = messageByLocaleService.getMessage("mcc.wallet."+walletAddr+".gaslimit");

        Credentials credentials = WalletUtils.loadCredentials(walletPW, payWalletFile);


        BigInteger GAS_PRICE = new BigInteger(gasPrice);
        BigInteger GAS_LIMIT = new BigInteger(gasLimit);

        return transferMCCAsync(credentials, toAddr, mccValue, GAS_PRICE, GAS_LIMIT);
    }

    public CompletableFuture<TransactionReceipt> transferMCCAsync(Credentials credentials, String toAddr, BigInteger mccValue, BigInteger GAS_PRICE, BigInteger GAS_LIMIT) throws Exception {
        MCCToken mccToken = getMCCToken(credentials, new StaticGasProvider(GAS_PRICE, GAS_LIMIT));
        log.trace("toAddr {}, mccValue {}", toAddr, mccValue);

        if (!mccToken.isValid()) {
            //TODO..
            throw  new Exception("Contracrt is Invalid");
        }

        return mccToken.transfer(toAddr, mccValue).sendAsync();

    }

    public TransactionReceipt approveMCC(File ownerWalletFile, String ownerPW, String toAddr, BigInteger mccValue) throws Exception {
        Credentials credentials = WalletUtils.loadCredentials(ownerPW, ownerWalletFile);

        BigInteger GAS_PRICE = new BigInteger("5000000000");
        BigInteger GAS_LIMIT = new BigInteger("50000");

        MCCToken mccToken = getMCCToken(credentials, new StaticGasProvider(GAS_PRICE, GAS_LIMIT));
        log.trace("toAddr {}, ethValue {}", toAddr, mccValue);
        if (!mccToken.isValid()) {
            //TODO..
            throw  new Exception("Contracrt is Invalid");
        }

        TransactionReceipt transactionReceipt = mccToken.approve(toAddr, mccValue).send();
        log.trace("{}", transactionReceipt);

        return transactionReceipt;
    }
    public TransactionReceipt approveMCC(String senderPrivateKey, String toAddr, BigInteger mccValue) throws Exception {
        Credentials credentials = Credentials.create(senderPrivateKey);

        BigInteger GAS_PRICE = new BigInteger("5000000000");
        BigInteger GAS_LIMIT = new BigInteger("50000");

        MCCToken mccToken = getMCCToken(credentials, new StaticGasProvider(GAS_PRICE, GAS_LIMIT));
        log.trace("toAddr {}, ethValue {}", toAddr, mccValue);
        if (!mccToken.isValid()) {
            //TODO..
            throw  new Exception("Contracrt is Invalid");
        }

        TransactionReceipt transactionReceipt = mccToken.approve(toAddr, mccValue).send();
        log.trace("{}", transactionReceipt);

        return transactionReceipt;
    }

    public TransactionReceipt approveIncMCC(File ownerWalletFile, String ownerPW, String toAddr, BigInteger ethValue) throws Exception {
        Credentials credentials = WalletUtils.loadCredentials(ownerPW, ownerWalletFile);

        BigInteger GAS_PRICE = new BigInteger("5000000000");
        BigInteger GAS_LIMIT = new BigInteger("50000");

        MCCToken mccToken = getMCCToken(credentials, new StaticGasProvider(GAS_PRICE, GAS_LIMIT));
        log.trace("toAddr {}, ethValue {}", toAddr, ethValue);
        if (!mccToken.isValid()) {
            //TODO..
            throw  new Exception("Contracrt is Invalid");
        }

        TransactionReceipt transactionReceipt = mccToken.increaseApproval(toAddr, ethValue).send();
        log.trace("{}", transactionReceipt);

        return transactionReceipt;
    }

    public TransactionReceipt approveDecMCC(File ownerWalletFile, String ownerPW, String toAddr, BigInteger ethValue) throws Exception {
        Credentials credentials = WalletUtils.loadCredentials(ownerPW, ownerWalletFile);

        BigInteger GAS_PRICE = new BigInteger("5000000000");
        BigInteger GAS_LIMIT = new BigInteger("50000");

        MCCToken mccToken = getMCCToken(credentials, new StaticGasProvider(GAS_PRICE, GAS_LIMIT));
        log.trace("toAddr {}, ethValue {}", toAddr, ethValue);
        if (!mccToken.isValid()) {
            //TODO..
            throw  new Exception("Contracrt is Invalid");
        }

        TransactionReceipt transactionReceipt = mccToken.decreaseApproval(toAddr, ethValue).send();
        log.trace("{}", transactionReceipt);

        return transactionReceipt;
    }

    public BigInteger getAllowanceMCC(String owner, String toAddr) throws Exception {
        String queryWalletFile = messageByLocaleService.getMessage("mcc.wallet." + queryWalletAddr + ".path");
        Credentials credentials = WalletUtils.loadCredentials(queryWalletFilePW, queryWalletFile);

        BigInteger GAS_PRICE = new BigInteger("5000000000");
        BigInteger GAS_LIMIT = new BigInteger("2000000");

        MCCToken mccToken = getMCCToken(credentials, new StaticGasProvider(GAS_PRICE, GAS_LIMIT));
        log.trace("toAddr {}", toAddr);
        if (!mccToken.isValid()) {
            //TODO..
            throw  new Exception("Contracrt is Invalid");
        }

        BigInteger bigInteger = mccToken.allowance(owner, toAddr).send();

        return bigInteger;
    }

    public TransactionReceipt transferFromMCC(File ownerWalletFile, String ownerPW,
                                              String fromAddr, String toAddr, BigInteger ethValue) throws Exception {
        //Credentials credentials = Credentials.create(ownerPVKey);
        Credentials credentials = WalletUtils.loadCredentials(ownerPW, ownerWalletFile);

        BigInteger GAS_PRICE = new BigInteger("5000000000");
        BigInteger GAS_LIMIT = new BigInteger("2000000");

        MCCToken mccToken = getMCCToken(credentials, new StaticGasProvider(GAS_PRICE, GAS_LIMIT));
        log.trace("toAddr {}, ethValue {}", toAddr, ethValue);
        if (!mccToken.isValid()) {
            //TODO..
            throw  new Exception("Contracrt is Invalid");
        }

        TransactionReceipt transactionReceipt = mccToken.transferFrom(fromAddr, toAddr, ethValue).send();
        log.trace("{}", transactionReceipt);

        return transactionReceipt;
    }
    public TransactionReceipt transferFromMCC(String ownerPrivateKey,
                                              String fromAddr, String toAddr, BigInteger ethValue) throws Exception {
        Credentials credentials = Credentials.create(ownerPrivateKey);

        BigInteger GAS_PRICE = new BigInteger("5000000000");
        BigInteger GAS_LIMIT = new BigInteger("2000000");

        MCCToken mccToken = getMCCToken(credentials, new StaticGasProvider(GAS_PRICE, GAS_LIMIT));
        log.trace("toAddr {}, ethValue {}", toAddr, ethValue);
        if (!mccToken.isValid()) {
            //TODO..
            throw  new Exception("Contracrt is Invalid");
        }

        TransactionReceipt transactionReceipt = mccToken.transferFrom(fromAddr, toAddr, ethValue).send();
        log.trace("{}", transactionReceipt);

        return transactionReceipt;
    }

    /**
     * 매일 10000개의 DAD 토큰을 이관한다.
     * @return Eco Wallet의 현재 잔고
     */
    public TransactionReceipt transferDailyDADTokens() throws Exception {

        KeyPair pair = RSATools.getKeyPairFromKeyStore(new ClassPathResource("mcc_batch.jks").getInputStream(), decryptKeystorePW(), "batch_private");

        log.debug("KeyPair:{}", pair);

        String ecoWalletPW = RSATools.decrypt(ecoWalletAddrPWEnc, pair.getPrivate());

        TransactionReceipt receipt = transferMCC(ecoWalletAddr, ecoWalletPW, payReservWalletAddr, new BigInteger("10000000000000000000000"));

        return receipt;
    }

    /**
     * Fruit 을 토큰으로 변환 처리... pay 계정에서 사용자 계정으로 토큰 이동
     * @return
     * @throws Exception
     */
    public TransactionReceipt exchangeFruit2Token(String toAddr, BigInteger amount) throws Exception {
        KeyPair pair = RSATools.getKeyPairFromKeyStore(new ClassPathResource("mcc_batch.jks").getInputStream(), decryptKeystorePW(), "pay_wallet");

        String payReservWalletAddrPW = RSATools.decrypt(payReservWalletAddrPWEnc, pair.getPrivate());

        TransactionReceipt receipt = transferMCC(payReservWalletAddr, payReservWalletAddrPW, toAddr, amount);

        return receipt;
    }

    public List<TransactionReceipt> exchangeFruit2Tokens(List<String> toAddrs, List<BigInteger> amounts) throws Exception {
        KeyPair pair = RSATools.getKeyPairFromKeyStore(new ClassPathResource("mcc_batch.jks").getInputStream(), decryptKeystorePW(), "pay_wallet");

        String payReservWalletAddrPW = RSATools.decrypt(payReservWalletAddrPWEnc, pair.getPrivate());

        String payReservWalletAddrFilePath = messageByLocaleService.getMessage("mcc.wallet."+payReservWalletAddr+".path");

        List<TransactionReceipt> receipts = transferMCCByDistrib(payReservWalletAddrFilePath, payReservWalletAddrPW, toAddrs, amounts,
                true, true);

        return receipts;
    }

    public CompletableFuture<TransactionReceipt> exchangeFruit2TokenAsync(String toAddr, BigInteger amount) throws Exception {
        KeyPair pair = RSATools.getKeyPairFromKeyStore(new ClassPathResource("mcc_batch.jks").getInputStream(), decryptKeystorePW(), "pay_wallet");

        String payReservWalletAddrPW = RSATools.decrypt(payReservWalletAddrPWEnc, pair.getPrivate());

        return transferMCCAsync(payReservWalletAddr, payReservWalletAddrPW, toAddr, amount);

    }

    public CompletableFuture<TransactionReceipt> exchangeFruit2TokensAsync(List<String> toAddrs, List<BigInteger> amounts) throws Exception {
        KeyPair pair = RSATools.getKeyPairFromKeyStore(new ClassPathResource("mcc_batch.jks").getInputStream(), decryptKeystorePW(), "pay_wallet");

        String payReservWalletAddrPW = RSATools.decrypt(payReservWalletAddrPWEnc, pair.getPrivate());

        String payReservWalletAddrFilePath = messageByLocaleService.getMessage("mcc.wallet."+payReservWalletAddr+".path");

        return transferMCCByDistribAsync(payReservWalletAddrFilePath, payReservWalletAddrPW, toAddrs, amounts,
                false, false);
    }

    public String isWalletOwner(final String walletAddr, final String walletPW) throws Exception {
        String payWalletFile = messageByLocaleService.getMessage("mcc.wallet."+walletAddr+".path");

        Credentials credentials = WalletUtils.loadCredentials(walletPW, payWalletFile);

        return credentials.getAddress();
    }

    //===================
    //MCC Distrib Info

    private MCCDistrib getMCCDistrib(Credentials credentials, ContractGasProvider contractGasProvider) {
        if ("mainnet".equalsIgnoreCase(targetNetwork))
            return MCCDistribMain.load(mccDistribContractAddress, web3j, credentials, contractGasProvider);
        else if ("rinkeby".equalsIgnoreCase(targetNetwork))
            return MCCDistribRinkeby.load(mccDistribContractAddress, web3j, credentials, contractGasProvider);

        return MCCDistribRinkeby.load(mccDistribContractAddress, web3j, credentials, contractGasProvider);

    }

    public HashMap<String, Object> getMCCDistribInfo() throws Exception {
        String queryWalletFile = messageByLocaleService.getMessage("mcc.wallet." + queryWalletAddr + ".path");

        Credentials credentials = WalletUtils.loadCredentials(queryWalletFilePW, queryWalletFile);


        BigInteger GAS_PRICE = new BigInteger("5000000000");
        BigInteger GAS_LIMIT = new BigInteger("400000");

        MCCDistrib mccDistrib = getMCCDistrib(credentials, new StaticGasProvider(GAS_PRICE, GAS_LIMIT));

        if (!mccDistrib.isValid()) {
            throw new Exception("Contracrt is Invalid");
            //LOGGER.debug("Contract is Invalid:{}", mccDistrib);
        }

        HashMap<String, Object> mapInfo = new HashMap<>();
        mapInfo.put("address", mccDistrib.getContractAddress());
        mapInfo.put("master", mccDistrib.manager().send());
        mapInfo.put("owners", mccDistrib.getOwners().send());

        return mapInfo;
    }

    public void setMCCDistribOwners(String managerWalletFilePath, String managerWalletPW, String newOwnerWallets) throws Exception {

        Credentials credentials = WalletUtils.loadCredentials(managerWalletPW, managerWalletFilePath);

        BigInteger GAS_PRICE = new BigInteger("5000000000");
        BigInteger GAS_LIMIT = new BigInteger("400000");

        MCCDistrib mccDistrib = getMCCDistrib(credentials, new StaticGasProvider(GAS_PRICE, GAS_LIMIT));

        if (!mccDistrib.isValid()) {
            throw new Exception("Contracrt is Invalid");
            //LOGGER.debug("Contract is Invalid:{}", mccDistrib);
        }

        ArrayList<String> owerWallets = new ArrayList<String>();

        String[] owners = newOwnerWallets.split(",");
        for (String wallet : owners) {
            if (wallet != null && wallet.length() > 10)
                owerWallets.add(wallet.trim());
        }
        log.trace("setDistribOwner==>{}", owerWallets);
        TransactionReceipt trReceipt = mccDistrib.setOwners(owerWallets).send();

        log.trace("setDistribOwner==>{}", trReceipt);

    }

    public List<TransactionReceipt> transferMCCByDistrib(String ownerWalletFilePath, String ownerWalletPW,
                                                   List<String> receivers, List<BigInteger> values,
                                            boolean bTransMCCFromOwner, boolean onlyInsufficientAmtSend) throws Exception {

        if (receivers == null || values == null || receivers.size() != values.size())
            throw new Exception("Invalid request");

        Credentials credentials = WalletUtils.loadCredentials(ownerWalletPW, ownerWalletFilePath);

        BigInteger GAS_PRICE = new BigInteger("5000000000");

        ArrayList<TransactionReceipt> resultTR = new ArrayList<>();
        if (bTransMCCFromOwner) {
            BigInteger GAS_LIMIT = new BigInteger("400000");

            MCCToken mccToken = getMCCToken(credentials, new StaticGasProvider(GAS_PRICE, GAS_LIMIT));

            if (!mccToken.isValid()) {
                //TODO..
                throw  new Exception("Contracrt Token is Invalid");
            }
            BigInteger bgTotalValue = BigInteger.ZERO;

            for (BigInteger v: values) {
                bgTotalValue = bgTotalValue.add(v);
            }

            if (onlyInsufficientAmtSend) {
                BigInteger remainToken = mccToken.balanceOf(mccDistribContractAddress).send();
                log.trace("MCCDistb remainToken:{}", remainToken.toString());

                if (remainToken.compareTo(bgTotalValue) < 0) {
                    bgTotalValue = bgTotalValue.subtract(remainToken);

                }
            }
            //x2ho.validation check...
            BigInteger payWHasMCC = mccToken.balanceOf(credentials.getAddress()).send();
            if (payWHasMCC.compareTo(bgTotalValue) < 0) {
                throw new Exception("Insufficient balance of "+credentials.getAddress());
            }

            TransactionReceipt trReceipt = mccToken.transfer(mccDistribContractAddress, bgTotalValue).send();
            log.trace("token transfer to MCCDistrib:{} ==> {}", bgTotalValue.toString(), trReceipt);

            resultTR.add(trReceipt);
        }
        BigInteger GAS_LIMIT = new BigInteger("100000");
        GAS_LIMIT = GAS_LIMIT.add(new BigInteger("10000")).multiply(new BigInteger(String.valueOf(receivers.size())));

        MCCDistrib mccDistrib = getMCCDistrib(credentials, new StaticGasProvider(GAS_PRICE, GAS_LIMIT));

        if (!mccDistrib.isValid()) {
            throw new Exception("Contracrt is Invalid");
            //LOGGER.debug("Contract is Invalid:{}", mccDistrib);
        }

        TransactionReceipt trReceipt = mccDistrib.transferMulti(mccTokenContractAddress, receivers, values).send();
        log.trace("token transfer by MCCDistrib:{} ==> {}", trReceipt);

        resultTR.add(trReceipt);

        return resultTR;
    }

    public CompletableFuture<TransactionReceipt> transferMCCByDistribAsync(String ownerWalletFilePath, String ownerWalletPW,
                                                         List<String> receivers, List<BigInteger> values,
                                                         boolean bTransMCCFromOwner, boolean onlyInsufficientAmtSend) throws Exception {

        if (receivers == null || values == null || receivers.size() != values.size())
            throw new Exception("Invalid request");

        Credentials credentials = WalletUtils.loadCredentials(ownerWalletPW, ownerWalletFilePath);

        BigInteger GAS_PRICE = new BigInteger("5000000000");

        if (bTransMCCFromOwner) {
            BigInteger GAS_LIMIT = new BigInteger("400000");

            MCCToken mccToken = getMCCToken(credentials, new StaticGasProvider(GAS_PRICE, GAS_LIMIT));

            if (!mccToken.isValid()) {
                //TODO..
                throw  new Exception("Contracrt Token is Invalid");
            }
            BigInteger bgTotalValue = BigInteger.ZERO;

            for (BigInteger v: values) {
                bgTotalValue = bgTotalValue.add(v);
            }

            if (onlyInsufficientAmtSend) {
                BigInteger remainToken = mccToken.balanceOf(mccDistribContractAddress).send();
                log.trace("MCCDistb remainToken:{}", remainToken.toString());

                if (remainToken.compareTo(bgTotalValue) < 0) {
                    bgTotalValue = bgTotalValue.subtract(remainToken);

                }
            }
            TransactionReceipt trReceipt = mccToken.transfer(mccDistribContractAddress, bgTotalValue).send();
            log.trace("token transfer to MCCDistrib:{} ==> {}", bgTotalValue.toString(), trReceipt);

        }
        BigInteger GAS_LIMIT = (new BigInteger("100000")).multiply(new BigInteger(String.valueOf(receivers.size())));

        MCCDistrib mccDistrib = getMCCDistrib(credentials, new StaticGasProvider(GAS_PRICE, GAS_LIMIT));

        if (!mccDistrib.isValid()) {
            throw new Exception("Contracrt is Invalid");
            //LOGGER.debug("Contract is Invalid:{}", mccDistrib);
        }

        return mccDistrib.transferMulti(mccTokenContractAddress, receivers, values).sendAsync();

    }

    public Subscription registTokenTxObserve(DefaultBlockParameter startBlock, DefaultBlockParameter endBlock,
                                            MCCTokenTxListener listener) throws Exception {

        String queryWalletFile = messageByLocaleService.getMessage("mcc.wallet." + queryWalletAddr + ".path");

        Credentials credentials = WalletUtils.loadCredentials(queryWalletFilePW, queryWalletFile);

        BigInteger GAS_PRICE = new BigInteger("5000000000");
        BigInteger GAS_LIMIT = new BigInteger("400000");

        MCCToken mccToken = getMCCToken(web3jMCC, targetNetwork, credentials, new StaticGasProvider(GAS_PRICE, GAS_LIMIT));

        if (!mccToken.isValid()) {
            //TODO..
            throw new Exception("Contracrt is Invalid");
        }

        Subscription subscription = mccToken.transferEventObservable(startBlock, endBlock)
                .subscribe(tx -> {
                    String toAddress = tx.to.toString();
                    String fromAddress = tx.from.toString();
                    String txHash = tx.log.getTransactionHash();
                    BigInteger value = tx.value;

                    log.debug(">>>>>>>>>>> MCC LOG:{}=>{} value:{} log:{}", fromAddress, toAddress, value, tx.log);

                    try {

                        if (listener!=null) {
                            TransactionReceipt transactionReceipt =
                                    web3jMCC.ethGetTransactionReceipt(txHash).send().getTransactionReceipt().get();

                            Transaction transaction = web3jMCC.ethGetTransactionByHash(txHash).send().getTransaction().get();

                            EthBlock.Block ethBlock = web3jMCC.ethGetBlockByNumber(
                                    DefaultBlockParameter.valueOf(tx.log.getBlockNumber()), false).send().getBlock();

                            EthMccTxLog mccTxLog = new EthMccTxLog();

                            mccTxLog.setHash(txHash);
                            mccTxLog.setBlockNumber(tx.log.getBlockNumber().longValue());
                            mccTxLog.setTransactionIndex(tx.log.getTransactionIndex().intValue());
                            mccTxLog.setFromAddr(tx.from);
                            mccTxLog.setToAddr(tx.to);
                            mccTxLog.setValue(tx.value);
                            mccTxLog.setTxInternal((byte)0);
                            mccTxLog.setTimeStamp(ethBlock.getTimestamp().longValue());

                            mccTxLog.setGas(transactionReceipt.getGasUsed());
                            mccTxLog.setGasPrice(transaction.getGasPrice());
                            mccTxLog.setCumulativeGasUsed(transactionReceipt.getCumulativeGasUsed());

                            listener.receiveTx(mccTxLog);
                        }

                    } catch (Exception e) {
                        log.warn("{}", e);
                    }

                });

        return subscription;
    }

        //web3j.transactionObservable().filter(tx -> isWantedTransaction(tx)).subscribe(processTransaction,errorHandler,OnCompleteHandler);
        /*
        Subscription subscription = web3j
                .catchUpToLatestTransactionObservable(DefaultBlockParameterName.EARLIEST)
                //.replayTransactionsObservable(startBlock, endBlock)
                .filter(transaction -> {
                    return transaction.get;
                })
                .doOnError(throwable -> LOGGER.error("Error occurred", throwable))
                .doOnCompleted(() -> LOGGER.info("Completed"))
                .doOnEach(notification -> LOGGER.info("OnEach"))
                .subscribe(transaction -> {
                    LOGGER.info("Transaction {} of block {}. From {} to {}",
                            transaction.getHash(),
                            transaction.getBlockNumber(),
                            transaction.getFrom(),
                            transaction.getTo());
                });
*/

//        Subscription subscription = null;
//        subscription = web3j.transactionObservable()
//                .subscribe(tx -> {
//                    String toAddress = tx.getTo();
//                    String fromAddress = tx.getFrom();
//                    //String txHash = tx..getTransactionHash();
//
//                    tx.get
//                    log.info("txLOG:{}", tx);
//                    countDownLatch.countDown();
//                });
//        subscription = web3j.blockObservable(true).subscribe(
//                block -> {
//                    EthBlock.Block b = block.getBlock();
//
//                    LocalDateTime timestamp = Instant.ofEpochSecond(
//                            b.getTimestamp()
//                                    .longValueExact()).atZone(ZoneId.of("UTC")).toLocalDateTime();
//                    int transactionCount = b.getTransactions().size();
//                    String hash = b.getHash();
//                    String parentHash = b.getParentHash();
//
//                    log.info(
//                            timestamp + " " +
//                                    "Tx count: " + transactionCount + ", " +
//                                    "Hash: " + hash + ", " +
//                                    "Parent hash: " + parentHash
//                    );
//
//                    b.getTransactions()
//                    countDownLatch.countDown();
//                }
//        );
//        countDownLatch.await(10, TimeUnit.MINUTES);
//        subscription.unsubscribe();


    public String createUserInputBroker() throws Exception {
        String queryWalletFile = messageByLocaleService.getMessage("mcc.wallet." + queryWalletAddr + ".path");

        Credentials credentials = WalletUtils.loadCredentials(queryWalletFilePW, queryWalletFile);

        BigInteger GAS_PRICE = new BigInteger("5000000000");
        BigInteger GAS_LIMIT = new BigInteger("1200000");

        List<String> owners = Arrays.asList(mccBrokerContractOwners);

        MCCBroker mccBrokerInst = MCCBroker.deploy(web3j, credentials,  new StaticGasProvider(GAS_PRICE, GAS_LIMIT), owners).send();

        log.info("new Instanace:{}", mccBrokerInst);
        return mccBrokerInst.getContractAddress();

    }

    //===================
    //MCC Distrib Info

    private MCCBroker getMCCBroker(Credentials credentials, ContractGasProvider contractGasProvider) {
        return MCCBroker.load(mccBrokerContractAddress, web3j, credentials, contractGasProvider);

        //return MCCBroker.load("0x0b05ed61d3a3a4ad763f64ea687790f858c07a8e", web3j, credentials, contractGasProvider);
    }

    public HashMap<String, Object> getMCBrokerInfo() throws Exception {
        String queryWalletFile = messageByLocaleService.getMessage("mcc.wallet." + queryWalletAddr + ".path");

        Credentials credentials = WalletUtils.loadCredentials(queryWalletFilePW, queryWalletFile);


        BigInteger GAS_PRICE = new BigInteger("5000000000");
        BigInteger GAS_LIMIT = new BigInteger("400000");

        MCCBroker mccBroker = getMCCBroker(credentials, new StaticGasProvider(GAS_PRICE, GAS_LIMIT));

        if (!mccBroker.isValid()) {
            throw new Exception("Contracrt is Invalid");
            //LOGGER.debug("Contract is Invalid:{}", mccDistrib);
        }

        HashMap<String, Object> mapInfo = new HashMap<>();
        mapInfo.put("address", mccBroker.getContractAddress());
        mapInfo.put("master", mccBroker.manager().send());
        mapInfo.put("owners", mccBroker.getOwners().send());

        return mapInfo;
    }

    public void setMCCBrokerOwners(String managerWalletFilePath, String managerWalletPW, String newOwnerWallets) throws Exception {

        Credentials credentials = WalletUtils.loadCredentials(managerWalletPW, managerWalletFilePath);

        BigInteger GAS_PRICE = new BigInteger("5000000000");
        BigInteger GAS_LIMIT = new BigInteger("400000");

        MCCBroker mccBroker = getMCCBroker(credentials, new StaticGasProvider(GAS_PRICE, GAS_LIMIT));

        if (!mccBroker.isValid()) {
            throw new Exception("Contracrt is Invalid");
            //LOGGER.debug("Contract is Invalid:{}", mccDistrib);
        }

        ArrayList<String> owerWallets = new ArrayList<String>();

        String[] owners = newOwnerWallets.split(",");
        for (String wallet : owners) {
            if (wallet != null && wallet.length() > 10)
                owerWallets.add(wallet.trim());
        }
        log.trace("setDistribOwner==>{}", owerWallets);
        TransactionReceipt trReceipt = mccBroker.setOwners(owerWallets).send();

        log.trace("setDistribOwner==>{}", trReceipt);

    }

    public List<TransactionReceipt> transferMCCByBroker(String ownerWalletFilePath, String ownerWalletPW,
                                                         List<String> receivers, List<BigInteger> values,
                                                         boolean bTransMCCFromOwner, boolean onlyInsufficientAmtSend) throws Exception {

        if (receivers == null || values == null || receivers.size() != values.size())
            throw new Exception("Invalid request");

        Credentials credentials = WalletUtils.loadCredentials(ownerWalletPW, ownerWalletFilePath);

        BigInteger GAS_PRICE = new BigInteger("5000000000");

        ArrayList<TransactionReceipt> resultTR = new ArrayList<>();
        if (bTransMCCFromOwner) {
            BigInteger GAS_LIMIT = new BigInteger("400000");

            MCCToken mccToken = getMCCToken(credentials, new StaticGasProvider(GAS_PRICE, GAS_LIMIT));

            if (!mccToken.isValid()) {
                //TODO..
                throw  new Exception("Contracrt Token is Invalid");
            }
            BigInteger bgTotalValue = BigInteger.ZERO;

            for (BigInteger v: values) {
                bgTotalValue = bgTotalValue.add(v);
            }

            if (onlyInsufficientAmtSend) {
                BigInteger remainToken = mccToken.balanceOf(mccBrokerContractAddress).send();
                log.trace("MCCDistb remainToken:{}", remainToken.toString());

                if (remainToken.compareTo(bgTotalValue) < 0) {
                    bgTotalValue = bgTotalValue.subtract(remainToken);

                }
            }
            TransactionReceipt trReceipt = mccToken.transfer(mccBrokerContractAddress, bgTotalValue).send();
            log.trace("token transfer to MCCDistrib:{} ==> {}", bgTotalValue.toString(), trReceipt);

            resultTR.add(trReceipt);
        }
        BigInteger GAS_LIMIT = (new BigInteger("100000")).multiply(new BigInteger(String.valueOf(receivers.size())));

        MCCBroker mccBroker = getMCCBroker(credentials, new StaticGasProvider(GAS_PRICE, GAS_LIMIT));

        if (!mccBroker.isValid()) {
            throw new Exception("Contracrt is Invalid");
            //LOGGER.debug("Contract is Invalid:{}", mccDistrib);
        }

        TransactionReceipt trReceipt = mccBroker.transferMulti(mccTokenContractAddress, receivers, values).send();
        log.trace("token transfer by MCCDistrib:{} ==> {}", trReceipt);

        resultTR.add(trReceipt);

        return resultTR;
    }

    public CompletableFuture<TransactionReceipt> transferMCCByBrokerAsync(String ownerWalletFilePath, String ownerWalletPW,
                                                                           List<String> receivers, List<BigInteger> values,
                                                                           boolean bTransMCCFromOwner, boolean onlyInsufficientAmtSend) throws Exception {

        if (receivers == null || values == null || receivers.size() != values.size())
            throw new Exception("Invalid request");

        Credentials credentials = WalletUtils.loadCredentials(ownerWalletPW, ownerWalletFilePath);

        BigInteger GAS_PRICE = new BigInteger("5000000000");

        if (bTransMCCFromOwner) {
            BigInteger GAS_LIMIT = new BigInteger("400000");

            MCCToken mccToken = getMCCToken(credentials, new StaticGasProvider(GAS_PRICE, GAS_LIMIT));

            if (!mccToken.isValid()) {
                //TODO..
                throw  new Exception("Contracrt Token is Invalid");
            }
            BigInteger bgTotalValue = BigInteger.ZERO;

            for (BigInteger v: values) {
                bgTotalValue = bgTotalValue.add(v);
            }

            if (onlyInsufficientAmtSend) {
                BigInteger remainToken = mccToken.balanceOf(mccBrokerContractAddress).send();
                log.trace("MCCDistb remainToken:{}", remainToken.toString());

                if (remainToken.compareTo(bgTotalValue) < 0) {
                    bgTotalValue = bgTotalValue.subtract(remainToken);

                }
            }
            TransactionReceipt trReceipt = mccToken.transfer(mccBrokerContractAddress, bgTotalValue).send();
            log.trace("token transfer to MCCDistrib:{} ==> {}", bgTotalValue.toString(), trReceipt);

        }
        BigInteger GAS_LIMIT = (new BigInteger("100000")).multiply(new BigInteger(String.valueOf(receivers.size())));

        MCCBroker mccBroker = getMCCBroker(credentials, new StaticGasProvider(GAS_PRICE, GAS_LIMIT));

        if (!mccBroker.isValid()) {
            throw new Exception("Contracrt is Invalid");
            //LOGGER.debug("Contract is Invalid:{}", mccDistrib);
        }

        return mccBroker.transferMulti(mccTokenContractAddress, receivers, values).sendAsync();

    }
    //=============
    public String decryptKeystorePW() {
        if (keystorePW.startsWith("ENC@[") && keystorePW.endsWith("]")) {
            keystorePW = stringEncryptor.decrypt(keystorePW.substring(5, keystorePW.length()-1));
        }
        return keystorePW;
    }


    public void jksReadTest() throws Exception {
        KeyPair pair = RSATools.getKeyPairFromKeyStore(new ClassPathResource("mcc_batch.jks").getInputStream(), decryptKeystorePW(), "pay_wallet");

        log.trace("{}, {}", payReservWalletAddr, payReservWalletAddrPWEnc);

        String payReservWalletAddrPW = RSATools.decrypt(payReservWalletAddrPWEnc, pair.getPrivate());


    }
}
