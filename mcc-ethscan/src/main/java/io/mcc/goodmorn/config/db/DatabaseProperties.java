package io.mcc.goodmorn.config.db;

public interface DatabaseProperties {
	
	public String getPoolName();
	
	public String getJdbcUrl();
	
	public String getUsername();
	
	public String getPassword();
	
	public String getDriverClassName();
	
	public boolean isAutoCommit();
	
	public long getConnectionTimeoutMs();
	
	public long getIdleTimeoutMs();
	
	public int getMinIdle();
	
	public int getMaximumPoolSize();
	
	public String getConnectionTestQuery();
	 
}