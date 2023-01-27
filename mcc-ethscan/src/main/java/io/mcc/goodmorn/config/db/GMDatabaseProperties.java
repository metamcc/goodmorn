package io.mcc.goodmorn.config.db;

import lombok.Getter;
import lombok.Setter;
import lombok.ToString;
import org.springframework.boot.context.properties.ConfigurationProperties;

@Getter
@Setter
@ToString(exclude="password")
@ConfigurationProperties(prefix = "goodmorn.datasource")
public class GMDatabaseProperties implements DatabaseProperties {

	private String poolName;
	private String jdbcUrl;
	private String username;
	private String password;
	private String driverClassName;
	private boolean isAutoCommit;
	private long connectionTimeoutMs;
	private long idleTimeoutMs;
	private int minIdle;
	private int maximumPoolSize;
	private String connectionTestQuery;
	
}
