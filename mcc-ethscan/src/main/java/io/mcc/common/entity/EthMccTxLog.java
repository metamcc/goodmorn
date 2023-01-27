package io.mcc.common.entity;

import lombok.Data;

import java.sql.Timestamp;

@Data
public class EthMccTxLog extends EthTransaction {
    private Long txSeq;
    private String txTime;

    public String getFromAddr() {
        return this.getFrom();
    }
    public String getToAddr() {
        return this.getTo();
    }
    public void setFromAddr(String _fromAddr) {
        this.setFrom(_fromAddr);
    }
    public void setToAddr(String _toAddr) {
        this.setTo(_toAddr);
    }

    public static EthMccTxLog getInstance(io.api.etherscan.model.Tx txLog) {
        EthMccTxLog mccTxLog = new EthMccTxLog();

        mccTxLog.setHash(txLog.getHash());
        mccTxLog.setBlockNumber(txLog.getBlockNumber());
        mccTxLog.setTransactionIndex(txLog.getTransactionIndex());
        mccTxLog.setFromAddr(txLog.getFrom());
        mccTxLog.setToAddr(txLog.getTo());
        mccTxLog.setValue(txLog.getValue());
        mccTxLog.setTxInternal((byte)0);
        mccTxLog.setTimeStamp(Timestamp.valueOf(txLog.getTimeStamp().now()).getTime());
        mccTxLog.setGas(txLog.getGas());
        mccTxLog.setGasPrice(txLog.getGasPrice());
        mccTxLog.setCumulativeGasUsed(txLog.getCumulativeGasUsed());
        return mccTxLog;
    }
}
