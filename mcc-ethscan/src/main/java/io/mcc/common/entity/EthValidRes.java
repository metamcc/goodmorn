package io.mcc.common.entity;

public class EthValidRes {
    private String blockchain;
    private Boolean valid_address;


    public String getBlockchain() {
        return blockchain;
    }

    public void setBlockchain(String blockchain) {
        this.blockchain = blockchain;
    }

    public Boolean getValid_address() {
        return valid_address;
    }

    public void setValid_address(Boolean valid_address) {
        this.valid_address = valid_address;
    }
}
