package io.mcc.goodmorn.config;

import com.zaxxer.hikari.HikariConfig;
import com.zaxxer.hikari.HikariDataSource;
import io.mcc.goodmorn.config.db.DatabaseProperties;
import io.mcc.goodmorn.config.db.GMDatabaseProperties;
import org.jasypt.encryption.StringEncryptor;
import org.mybatis.spring.annotation.MapperScan;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Qualifier;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.boot.context.properties.EnableConfigurationProperties;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.context.annotation.DependsOn;
import org.springframework.jdbc.datasource.DataSourceTransactionManager;
import org.springframework.transaction.PlatformTransactionManager;

import javax.sql.DataSource;

public abstract class DatabaseConfig {

	@Autowired
	@Qualifier("encryptorBean")
	private StringEncryptor stringEncryptor;

	@Value("${jasypt.encryptor.property.prefix:ENC@(}")
	private String prefixSE;
	@Value("${jasypt.encryptor.property.suffix:)}")
	private String suffixSE;

	public abstract DataSource dataSource();

	protected DataSource configureDataSource(DatabaseProperties databaseProperties) {
		HikariConfig jdbcConfig = new HikariConfig();
		jdbcConfig.setPoolName(databaseProperties.getPoolName());
		jdbcConfig.setDriverClassName(databaseProperties.getDriverClassName());
		jdbcConfig.setJdbcUrl(databaseProperties.getJdbcUrl());
		jdbcConfig.setUsername(databaseProperties.getUsername());

		String encPassword = databaseProperties.getPassword();
		if (encPassword.startsWith(prefixSE) && encPassword.endsWith(suffixSE)) {
			encPassword = encPassword.substring(prefixSE.length(), encPassword.length()-suffixSE.length());

			jdbcConfig.setPassword(stringEncryptor.decrypt(encPassword));
		} else
			jdbcConfig.setPassword(databaseProperties.getPassword());

		jdbcConfig.setAutoCommit(databaseProperties.isAutoCommit());
		jdbcConfig.setConnectionTimeout(databaseProperties.getConnectionTimeoutMs());
		jdbcConfig.setIdleTimeout(databaseProperties.getIdleTimeoutMs());
		jdbcConfig.setMinimumIdle(databaseProperties.getMinIdle());
		jdbcConfig.setMaximumPoolSize(databaseProperties.getMaximumPoolSize());
		jdbcConfig.setConnectionTestQuery(databaseProperties.getConnectionTestQuery());

		return new HikariDataSource(jdbcConfig);
	}
}

@Configuration
@EnableConfigurationProperties(GMDatabaseProperties.class)
class GMDatabaseConfig extends DatabaseConfig {

	@Autowired
	private GMDatabaseProperties gmDatabaseProperties;
	
	@Override
	@Bean(name = "gmDataSource")
	public DataSource dataSource() {
		return configureDataSource(gmDatabaseProperties);
	}

	@Bean(name = "gmTransactionManager")
	public PlatformTransactionManager userTransactionManager(@Qualifier("gmDataSource") DataSource mailDataSource) {
		DataSourceTransactionManager transactionManager = new DataSourceTransactionManager(mailDataSource);
		return transactionManager;
	}	
}

//@Configuration
//@EnableConfigurationProperties(MobileDatabaseProperties.class)
//class MobileDatabaseConfig extends DatabaseConfig {
//
//	@Autowired
//	private MobileDatabaseProperties mobileDatabaseProperties;
//
//	@Override
//	@Bean(name = "mobileDataSource")
//	@Primary
//	public DataSource dataSource() {
//		return configureDataSource(mobileDatabaseProperties);
//	}
//
//	@Bean(name="mobileTransactionManager")
//	@Primary
//	public PlatformTransactionManager mobileTransactionManager(@Qualifier("mobileDataSource") DataSource mobileDataSource) {
//		DataSourceTransactionManager transactionManager = new DataSourceTransactionManager(mobileDataSource);
//		return transactionManager;
//	}
//}