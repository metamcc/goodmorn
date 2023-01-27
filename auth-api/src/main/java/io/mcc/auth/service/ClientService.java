package io.mcc.auth.service;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.security.oauth2.client.token.DefaultAccessTokenRequest;
import org.springframework.security.oauth2.client.token.grant.client.ClientCredentialsAccessTokenProvider;
import org.springframework.security.oauth2.client.token.grant.client.ClientCredentialsResourceDetails;
import org.springframework.security.oauth2.client.token.grant.password.ResourceOwnerPasswordAccessTokenProvider;
import org.springframework.security.oauth2.client.token.grant.password.ResourceOwnerPasswordResourceDetails;
import org.springframework.security.oauth2.common.OAuth2AccessToken;
import org.springframework.security.oauth2.common.OAuth2RefreshToken;
import org.springframework.stereotype.Service;

import io.mcc.auth.common.config.property.ConfigProperties;
import io.mcc.auth.common.util.Util;
import io.mcc.auth.entity.CommonVO;
import io.mcc.auth.entity.User;

import lombok.extern.slf4j.Slf4j;

@Slf4j
@Service
public class ClientService {
	public static final Logger logger = LoggerFactory.getLogger(ClientService.class);

	@Autowired
	public ConfigProperties configProperties;
	
	@Autowired
	public PasswordEncoder passwordEncoder;

	private String compareStr = "";
	public CommonVO getResourceOwnerPasswordToken(String[] client, User user){
		logger.info("===== getResourceOwnerPasswordToken() start...");
		String test = passwordEncoder.encode(compareStr);
		
		logger.info("----- Username = " + user.getUsername());
		logger.info("----- Password = " + user.getPassword());
		logger.info("----- ClientId = " + client[0]);
		logger.info("----- ClientSecret = " + client[1]);
		logger.info("----- snsid = " + compareStr);
		logger.info("----- ClientSecret22 = " +test);
		logger.info("----- ClientSecret 11= " +("BCrypt 비교: " + passwordEncoder.matches(compareStr, test)));
		
		CommonVO rst = new CommonVO();
		
		ResourceOwnerPasswordResourceDetails resource = new ResourceOwnerPasswordResourceDetails();
		resource.setAccessTokenUri(configProperties.getAccessTokenUri());
		resource.setGrantType("client_credentials");
		resource.setClientId(client[0]);
		resource.setClientSecret(client[1]);
		resource.setUsername(user.getUsername());
		resource.setPassword(user.getPassword());	//$2a$10$g8gDU1LKfZZVHwcclkpL5uFufVUqHniguD9DGJODPz8Rzsx/A302C
		
		
		ResourceOwnerPasswordAccessTokenProvider provider = new ResourceOwnerPasswordAccessTokenProvider();
		OAuth2AccessToken accessToken =  provider.obtainAccessToken(resource, new DefaultAccessTokenRequest());
		logger.info("----- accessToken : " + accessToken.getValue());
		logger.info("----- accessToken : " + Util.toDayText(accessToken.getExpiration().getTime()));
		
		int expires_in = 0;
		OAuth2RefreshToken refreshToken = null;
		
		if(accessToken!=null && !accessToken.isExpired()) {
			expires_in = accessToken.getExpiresIn();
			logger.info("----- expires_in : " + expires_in);
			
			refreshToken = accessToken.getRefreshToken(); 
			if(refreshToken == null) {
				logger.info("----- refreshToken null.");
			} else {
				logger.info("----- refreshToken : " + refreshToken.getValue());
			}
		}
		
		
		rst.put("accessToken", accessToken);
		
		return rst;
	}
	
	public CommonVO getClientCredentialsToken(String[] client) {
		logger.info("===== getClientCredentialsToken() start...");
		//authorization_code vs implicit
		//password,refresh_token,client_credentials
		
		CommonVO rst = new CommonVO();
		
		ClientCredentialsResourceDetails resource = new ClientCredentialsResourceDetails();
		resource.setAccessTokenUri(configProperties.getAccessTokenUri());
		resource.setGrantType("client_credentials");
		resource.setClientId(client[0]);
		resource.setClientSecret(client[1]);
		
		ClientCredentialsAccessTokenProvider provider = new ClientCredentialsAccessTokenProvider();
		OAuth2AccessToken accessToken =  provider.obtainAccessToken(resource, new DefaultAccessTokenRequest());
		
		logger.info("----- accessToken : " + accessToken.getValue());
		
		int expires_in = 0;
		OAuth2RefreshToken refreshToken = null;
		
		if(accessToken!=null && !accessToken.isExpired()) {
			expires_in = accessToken.getExpiresIn();
			logger.info("----- expires_in : " + expires_in);
			
			refreshToken = accessToken.getRefreshToken(); 
			if(refreshToken == null) {
				logger.info("----- refreshToken null.");
			} else {
				logger.info("----- refreshToken : " + refreshToken.getValue());
			}
		}
		
		rst.put("accessToken", accessToken);
		
		return rst;
	}
}
