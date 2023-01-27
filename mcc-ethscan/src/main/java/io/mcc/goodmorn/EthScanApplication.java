package io.mcc.goodmorn;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.EnableAutoConfiguration;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.cache.annotation.EnableCaching;
import org.springframework.scheduling.annotation.EnableScheduling;

@SpringBootApplication
@EnableScheduling
@EnableCaching
public class EthScanApplication {

    public static void main(String... args) throws Exception {
        SpringApplication.run(EthScanApplication.class, args);
    }
}
