package io.mcc.goodmorn.config;

import org.apache.ibatis.annotations.Mapper;
import org.apache.ibatis.session.SqlSessionFactory;
import org.mybatis.spring.SqlSessionFactoryBean;
import org.mybatis.spring.annotation.MapperScan;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Qualifier;
import org.springframework.context.ApplicationContext;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.context.annotation.Primary;
import org.springframework.core.io.support.PathMatchingResourcePatternResolver;

import javax.sql.DataSource;
import java.io.IOException;

public abstract class MyBatisConfig {
	
	public static final String BASE_PACKAGE_PREFIX = "io.mcc";
	
	public static final String CONFIG_LOCATION_PATH = "classpath:mybatis-config.xml";
	
	protected void configureSqlSessionFactory(SqlSessionFactoryBean sessionFactoryBean, DataSource dataSource) throws IOException {
		PathMatchingResourcePatternResolver pathResolver = new PathMatchingResourcePatternResolver();
		
		sessionFactoryBean.setDataSource(dataSource);
		sessionFactoryBean.setConfigLocation(pathResolver.getResource(CONFIG_LOCATION_PATH));
	}
}

@Configuration
@MapperScan(annotationClass=Mapper.class, basePackages=MyBatisConfig.BASE_PACKAGE_PREFIX, sqlSessionFactoryRef = "gmSqlSessionFactory")
class GMMybatisConfig extends MyBatisConfig {

	@Autowired
	private ApplicationContext applicationContext;

	@Bean(name = "gmSqlSessionFactory")
	@Primary
	public SqlSessionFactory gmSqlSessionFactory(@Qualifier("gmDataSource") DataSource gmDataSource) throws Exception {
		SqlSessionFactoryBean sessionFactoryBean = new SqlSessionFactoryBean();
		//sessionFactoryBean.setMapperLocations(applicationContext.getResources("classpath:mybatis/mapper/*.xml"));
		configureSqlSessionFactory(sessionFactoryBean, gmDataSource);
		return sessionFactoryBean.getObject();
	}
}
