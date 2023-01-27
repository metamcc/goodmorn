package io.mcc.auth.common.config.tx;


import java.io.IOException;

import javax.sql.DataSource;

import org.apache.ibatis.session.SqlSessionFactory;
import org.mybatis.spring.SqlSessionFactoryBean;
import org.mybatis.spring.SqlSessionTemplate;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.context.ApplicationContext;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Primary;
import org.springframework.core.io.support.PathMatchingResourcePatternResolver;
import org.springframework.jdbc.datasource.DataSourceTransactionManager;
import org.springframework.transaction.PlatformTransactionManager;


public class TransactionConfig {
	
	public static final String CONFIG_LOCATION_PATH = "classpath:mybatis-config.xml";
	
	@Autowired
	private ApplicationContext applicationContext;
	
	@Autowired
	private DataSource dataSource;

	@Bean
	@Primary
	public PlatformTransactionManager transactionManager() {
		return new DataSourceTransactionManager(dataSource);
	}
	
	@Bean
    public SqlSessionTemplate sqlSessionTemplate(SqlSessionFactory sqlSessionFactory) throws Exception {
        return new SqlSessionTemplate(sqlSessionFactory);
    }
	
    @Bean
    public SqlSessionFactoryBean sqlSessionFactoryBean() throws IOException {
    	PathMatchingResourcePatternResolver pathResolver = new PathMatchingResourcePatternResolver();
    	
    	SqlSessionFactoryBean sessionFactoryBean = new SqlSessionFactoryBean();
        
    	sessionFactoryBean.setDataSource(dataSource);
    	sessionFactoryBean.setConfigLocation(pathResolver.getResource(CONFIG_LOCATION_PATH));
        return sessionFactoryBean;
    }
    
}
