package io.mcc.common.entity;

import lombok.Data;

import java.math.BigInteger;

@Data
public class EthTransaction {
    private Long blockNumber;
    private String hash;
    private Integer transactionIndex;
    private String from;
    private String to;
    private Long timeStamp;
    private BigInteger value;
    private byte txInternal;
    private char xchgFruit;

    private BigInteger gas;
    private BigInteger gasUsed;
    private BigInteger gasPrice;
    private BigInteger cumulativeGasUsed;
}
