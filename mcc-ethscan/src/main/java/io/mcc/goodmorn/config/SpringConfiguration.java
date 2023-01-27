package io.mcc.goodmorn.config;

import io.mcc.common.MessageByLocaleService;
import io.mcc.common.MessageByLocaleServiceImpl;
import lombok.extern.slf4j.Slf4j;
import org.mybatis.spring.annotation.MapperScan;
import org.springframework.beans.factory.annotation.Qualifier;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.ComponentScan;
import org.springframework.context.annotation.Configuration;
import org.springframework.context.support.ReloadableResourceBundleMessageSource;
import org.springframework.core.task.AsyncTaskExecutor;
import org.springframework.scheduling.annotation.EnableAsync;
import org.springframework.scheduling.concurrent.ThreadPoolTaskExecutor;
import org.web3j.protocol.Web3j;
import org.web3j.protocol.http.HttpService;

import java.util.concurrent.Callable;
import java.util.concurrent.Executor;
import java.util.concurrent.Future;

@Configuration
@EnableAsync
@ComponentScan("io.mcc")
@Slf4j
public class SpringConfiguration {

    @Value("${mcc.ethscan.workpool.max:10}")
    private int maxPoolSize;

    @Value("${mcc.ethscan.workpool.core:3}")
    private int corePoolSize;

    @Value("${mcc.ethscan.workpool.queue:3}")
    private int queueSize;

    @Value("${mcc.ethscan.workpool.enabled:false}")
    private boolean bThreadPoolEanabled;

    @Value("${web3j.client-address:'https://mainnet.infura.io/zEkQtfWnZUsXuN1DPUmA'}")
    private String web3jProvider;

    @Value("${web3j.mcc-client-address:'http://211.232.21.47:38546'}")
    private String web3jMCCProvider;

    @Value("${project.base-dir}")
    private String projectBaseDir;

    @Bean(name = "threadPoolTaskExecutor")
    public Executor threadPoolTaskExecutor() {
        log.info(">>>>>>>>>>>>>>>>>>>>>>>> threadPool:{}", bThreadPoolEanabled);

        if (bThreadPoolEanabled) {
            ThreadPoolTaskExecutor taskExecutor = new ThreadPoolTaskExecutor();
            taskExecutor.setCorePoolSize(corePoolSize);
            taskExecutor.setMaxPoolSize(maxPoolSize);
            taskExecutor.setQueueCapacity(queueSize);
            taskExecutor.setThreadNamePrefix("Executor-");
            //taskExecutor.setRejectedExecutionHandler(new ThreadPoolExecutor.CallerRunsPolicy());
            taskExecutor.initialize();
            return new HandlingExecutor(taskExecutor); // HandlingExecutor로 wrapping 합니다.
        }
        return null;
    }

    public class HandlingExecutor implements AsyncTaskExecutor {
        private AsyncTaskExecutor executor;

        public HandlingExecutor(AsyncTaskExecutor executor) {
            this.executor = executor;
        }

        @Override
        public void execute(Runnable task) {
            executor.execute(task);
        }

        @Override
        public void execute(Runnable task, long startTimeout) {
            executor.execute(createWrappedRunnable(task), startTimeout);
        }

        @Override
        public Future<?> submit(Runnable task) {
            return executor.submit(createWrappedRunnable(task));
        }

        @Override
        public <T> Future<T> submit(final Callable<T> task) {
            return executor.submit(createCallable(task));
        }

        private <T> Callable<T> createCallable(final Callable<T> task) {
            return new Callable<T>() {
                @Override
                public T call() throws Exception {
                    try {
                        return task.call();
                    } catch (Exception ex) {
                        handle(ex);
                        throw ex;
                    }
                }
            };
        }

        private Runnable createWrappedRunnable(final Runnable task) {
            return new Runnable() {
                @Override
                public void run() {
                    try {
                        task.run();
                    } catch (Exception ex) {
                        handle(ex);
                    }
                }
            };
        }

        private void handle(Exception ex) {
            log.info("Failed to execute task. : {}", ex.getMessage());
        }

    }


    @Bean
    public Web3j web3j() {
        Web3j web3j = Web3j.build(new HttpService(web3jProvider));
        return web3j;
    }

    @Bean
    @Qualifier("MCC")
    public Web3j web3jMCC() {
        Web3j web3j = Web3j.build(new HttpService(web3jMCCProvider));
        return web3j;
    }
}
