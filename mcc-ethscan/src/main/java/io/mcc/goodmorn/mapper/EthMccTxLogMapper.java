package io.mcc.goodmorn.mapper;


import io.api.etherscan.model.Tx;
import io.mcc.common.vo.CommonVO;
import io.mcc.common.entity.EthMccTxLog;
import org.apache.ibatis.annotations.Mapper;
import org.apache.ibatis.annotations.Param;

import java.util.List;

@Mapper
public interface EthMccTxLogMapper {

	int countLog(final EthMccTxLog log);
	int insertLog(final EthMccTxLog log);

	Long getLastBlockNumber(@Param("wallet") final String wallet);

	List<EthMccTxLog> getTransLog(@Param("fromAddr") final String fromAddr, @Param("toAddr") final String toAddr);

	List<EthMccTxLog> getTransLogByMap(CommonVO params);

	int setTransLogByM2FDone(@Param("hash") String hash, @Param("blockNumber") Integer blockNumber);

	List<EthMccTxLog> getTransactionLog(@Param("hash") String hash);


	int countMccTx(CommonVO param);
	List<EthMccTxLog> findMccTx(CommonVO param);

	long getMaxTxSEQ();

}
