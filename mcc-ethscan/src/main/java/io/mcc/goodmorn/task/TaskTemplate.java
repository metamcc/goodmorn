package io.mcc.goodmorn.task;

import io.mcc.common.MessageByLocaleService;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.scheduling.annotation.Scheduled;

/**
 * Created by wisehouse on 2017. 7. 12..
 */
@Slf4j
public abstract class TaskTemplate {

    @Autowired
    private MessageByLocaleService messageByLocaleService;


    @Scheduled
    abstract void run();

    boolean isServiceStop() {
        String stopAll = null;
        try {
            stopAll = messageByLocaleService.getMessage("mcc.batch.servicestop");
        } catch (Exception e){
            log.trace(e.toString());
        }

        if (stopAll!=null && "on".equalsIgnoreCase(stopAll))
            return true;
        return false;
    }

    boolean isKafkaServiceStop() {
        String stopKafka = null;
        try {
            stopKafka = messageByLocaleService.getMessage("mcc.batch.kafkastop");
        } catch (Exception e){
            log.trace(e.toString());
        }

        if (stopKafka!=null && "on".equalsIgnoreCase(stopKafka))
            return true;
        return false;
    }
}
