package io.mcc.goodmorn.task;

import java.math.BigDecimal;
import java.math.BigInteger;
import java.util.Calendar;
import java.util.HashMap;

import io.mcc.goodmorn.service.MCCTxService;
import io.mcc.mcctoken.service.Web3jService;
import org.apache.http.client.utils.URIBuilder;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.scheduling.annotation.Scheduled;
import org.springframework.stereotype.Component;
import org.web3j.protocol.core.DefaultBlockParameter;
import org.web3j.protocol.core.DefaultBlockParameterName;
import org.web3j.utils.Convert;

import lombok.extern.slf4j.Slf4j;
import rx.Subscription;

/**
 * ICO용 스케쥴러
 */
@Component
@Slf4j
public class ServiceMonitorScheduler extends TaskTemplate {

	private String _simpleName = this.getClass().getSimpleName();

	private boolean txObserverEnabled = true;

	private Subscription subscriptionMCCTx = null;


	@Autowired
	private Web3jService web3jService;

	@Autowired
	private MCCTxService mccTxService;

	public ServiceMonitorScheduler() {
		log.info("{} Created!!", _simpleName);
	}


	@Scheduled(fixedDelayString = "60000")
	@Override
	public void run() {

		log.debug("run {} subscriptionMCCTx:{}", _simpleName, subscriptionMCCTx);

		if (txObserverEnabled && (subscriptionMCCTx == null || subscriptionMCCTx.isUnsubscribed())) {
			Long lastBlock = mccTxService.getLastMccTxBlock(null);

			if (lastBlock == null || lastBlock == 0)
				lastBlock = 7500000L;

			try {
				subscriptionMCCTx = web3jService.registTokenTxObserve(DefaultBlockParameter.valueOf(BigInteger.valueOf(lastBlock)),
						DefaultBlockParameterName.LATEST, mccTxService);
			} catch(Exception e) {
				log.warn("{}", e);
			}
		}

	}
}
