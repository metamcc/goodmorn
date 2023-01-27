package io.mcc.mcctoken.support;

import io.mcc.common.entity.EthMccTxLog;

public interface MCCTokenTxListener {
    public void receiveTx(EthMccTxLog mccTxLog);
}
