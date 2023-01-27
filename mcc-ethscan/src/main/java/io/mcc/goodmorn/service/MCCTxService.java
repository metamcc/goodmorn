package io.mcc.goodmorn.service;

import io.api.etherscan.core.impl.EtherScanApi;
import io.api.etherscan.model.EthNetwork;
import io.api.etherscan.model.Tx;
import io.mcc.common.entity.EthMccTxLog;
import io.mcc.common.entity.EthTransRes;

import io.mcc.common.vo.CommonVO;
import io.mcc.goodmorn.mapper.*;
import io.mcc.mcctoken.service.Web3jService;
import io.mcc.mcctoken.support.MCCTokenTxListener;
import lombok.extern.slf4j.Slf4j;
import org.apache.commons.lang3.StringUtils;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;
import org.web3j.crypto.ECKeyPair;
import org.web3j.crypto.Keys;
import org.web3j.crypto.Wallet;
import org.web3j.crypto.WalletFile;
import org.web3j.protocol.core.methods.response.TransactionReceipt;

import java.io.File;
import java.math.BigDecimal;
import java.math.BigInteger;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.UUID;

@Service
@Slf4j
public class MCCTxService implements MCCTokenTxListener {

    @Autowired
    private EthMccTxLogMapper ethMccTxLogMapper;

    //---------------------------------------------------
    // MCCTokenTxListener
    public void receiveTx(EthMccTxLog mccTxLog) {
        ethMccTxLogMapper.insertLog(mccTxLog);
    }
   //---------------------------------------------------

    public Long getLastMccTxBlock(String wallet) {
        Long lastBlock = null;
        try {
            lastBlock = ethMccTxLogMapper.getLastBlockNumber(wallet);
        } catch (Exception e) {
            log.warn("{}", e);
        }

        if (lastBlock == null)
            lastBlock = 0L;
        else
            lastBlock = lastBlock + 1;

        return lastBlock;
    }

    public EthTransRes findMccTxLog(final long startBlock, final long endBlock, String orderby) {
        CommonVO params = new CommonVO();
        params.put("startblock", startBlock);
        params.put("endblock", endBlock);
        params.put("orderby", orderby.toLowerCase());

        List<EthMccTxLog> resultItems = ethMccTxLogMapper.findMccTx(params);

        EthTransRes txRes = new EthTransRes();
        txRes.setResult(resultItems);
        txRes.setStatus(1);
        txRes.setMessage(String.valueOf(resultItems.size()));
        return txRes;
    }
}
