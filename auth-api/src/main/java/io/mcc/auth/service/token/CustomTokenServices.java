package io.mcc.auth.service.token;

import java.util.Date;
import java.util.Set;
import java.util.UUID;

import org.springframework.beans.factory.InitializingBean;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.security.authentication.AuthenticationManager;
import org.springframework.security.core.Authentication;
import org.springframework.security.core.AuthenticationException;
import org.springframework.security.oauth2.common.DefaultExpiringOAuth2RefreshToken;
import org.springframework.security.oauth2.common.DefaultOAuth2AccessToken;
import org.springframework.security.oauth2.common.DefaultOAuth2RefreshToken;
import org.springframework.security.oauth2.common.ExpiringOAuth2RefreshToken;
import org.springframework.security.oauth2.common.OAuth2AccessToken;
import org.springframework.security.oauth2.common.OAuth2RefreshToken;
import org.springframework.security.oauth2.common.exceptions.InvalidGrantException;
import org.springframework.security.oauth2.common.exceptions.InvalidScopeException;
import org.springframework.security.oauth2.common.exceptions.InvalidTokenException;
import org.springframework.security.oauth2.provider.ClientDetails;
import org.springframework.security.oauth2.provider.ClientDetailsService;
import org.springframework.security.oauth2.provider.ClientRegistrationException;
import org.springframework.security.oauth2.provider.OAuth2Authentication;
import org.springframework.security.oauth2.provider.OAuth2Request;
import org.springframework.security.oauth2.provider.TokenRequest;
import org.springframework.security.oauth2.provider.token.AuthorizationServerTokenServices;
import org.springframework.security.oauth2.provider.token.ConsumerTokenServices;
import org.springframework.security.oauth2.provider.token.ResourceServerTokenServices;
import org.springframework.security.oauth2.provider.token.TokenEnhancer;
import org.springframework.security.oauth2.provider.token.TokenStore;
import org.springframework.security.web.authentication.preauth.PreAuthenticatedAuthenticationToken;
import org.springframework.stereotype.Component;
import org.springframework.transaction.annotation.Transactional;

import lombok.extern.slf4j.Slf4j;

@Slf4j
@Component
public class CustomTokenServices implements AuthorizationServerTokenServices, ResourceServerTokenServices, ConsumerTokenServices, InitializingBean{

	private int refreshTokenValiditySeconds = 60 * 60 * 24 * 30; // default 30 days.

	private int accessTokenValiditySeconds = 60 * 60 * 24 * 7; // default 7 days.

	private boolean supportRefreshToken = true;

	private boolean reuseRefreshToken = true;

	private boolean updateRefreshTokenLastLoginTime = true;
	
	@Autowired
	private TokenStore jdbcTokenStore;

	@Autowired
	private ClientDetailsService clientDetailsService;

	@Autowired
	private TokenEnhancer accessTokenEnhancer;

	@Autowired
	private AuthenticationManager authenticationManager;
	
	@Override
	public OAuth2Authentication loadAuthentication(String accessTokenValue) throws AuthenticationException, InvalidTokenException {
		
		OAuth2AccessToken accessToken = jdbcTokenStore.readAccessToken(accessTokenValue);
	    if (accessToken == null) {
	        throw new InvalidTokenException("Invalid access token: " + accessTokenValue);
	    }
	    else if (accessToken.isExpired()) {
	    	jdbcTokenStore.removeAccessToken(accessToken);
	        throw new InvalidTokenException("Access token expired: " + accessTokenValue);
	    }

	    OAuth2Authentication result = jdbcTokenStore.readAuthentication(accessToken);
	    if (result == null) {
	        // in case of race condition
	        throw new InvalidTokenException("Invalid access token: " + accessTokenValue);
	    }
	    if (clientDetailsService != null) {
	        String clientId = result.getOAuth2Request().getClientId();
	        try {
	            clientDetailsService.loadClientByClientId(clientId);
	        }
	        catch (ClientRegistrationException e) {
	            throw new InvalidTokenException("Client not valid: " + clientId, e);
	        }
	    }
	    return result;
	    
	}

	@Override
	public OAuth2AccessToken readAccessToken(String accessToken) {
		log.debug(">>>> called readAccessToken:{}", accessToken);

		return jdbcTokenStore.readAccessToken(accessToken);
	}

	@Override
	@Transactional
	public OAuth2AccessToken createAccessToken(OAuth2Authentication authentication) throws AuthenticationException {
		log.info("+++++++++++++++++  [CustomTokenServices] createAccessToken(): {} ", authentication);
		OAuth2AccessToken existingAccessToken = jdbcTokenStore.getAccessToken(authentication);
	    OAuth2RefreshToken refreshToken = null;
	    if (existingAccessToken != null) {
	        if (existingAccessToken.isExpired()) {
				log.info("+++++++++++++++++  existingAccessToken:{}", existingAccessToken);
	            if (existingAccessToken.getRefreshToken() != null) {
	                refreshToken = existingAccessToken.getRefreshToken();
	                //-????????? ????????? ???????????? ?????? ??????????????? ?????? ?????? ????????? ?????? ??? ??? ????????? ????????????????????? ...
	                jdbcTokenStore.removeRefreshToken(refreshToken);
					log.info("+++++++++++++++++  removeRefreshToken:{}", refreshToken);
					//x2ho.bugifx..??????
					//TODO.???
					//remove refreshToken?????? ??? ???????????? refreshToken????????? ??????????? null ????????? ???????????? ??????????
					//refreshToken = null;
	            }
	            jdbcTokenStore.removeAccessToken(existingAccessToken);
				log.info("+++++++++++++++++  removeAccessToken:{}", existingAccessToken);
	        } else {
	            //-????????? ????????? ?????? ????????? ????????? ?????? ??????.
	        	jdbcTokenStore.storeAccessToken(existingAccessToken, authentication);
	            return existingAccessToken;
	        }
	    }
		log.info("+++++++++++++++++  refreshToken:{} existingAccessToken:{}", refreshToken, existingAccessToken);
	    // ?????? ?????? ?????? ????????? ????????? ?????? ?????? ?????? ??? ??????.
	    // ?????? ??? ????????? ????????? ?????? ???.
	    // ?????????????????? ????????? ?????? ?????? ????????? ?????????????????? ??? ???????????? ?????? ????????? ????????? ?????? ??? ?????? ?????? ?????????.
	    if (refreshToken == null) {
	        refreshToken = createRefreshToken(authentication);
	    }
	    // ????????? ?????? ?????? ?????? ????????? ?????? ??? ?????? ?????? ???????????? ??? ?????? ??????.
	    else if (refreshToken instanceof ExpiringOAuth2RefreshToken) {
	        ExpiringOAuth2RefreshToken expiring = (ExpiringOAuth2RefreshToken) refreshToken;
	        if (System.currentTimeMillis() > expiring.getExpiration().getTime()) {
	            refreshToken = createRefreshToken(authentication);
	        }
	    }
		log.info("+++++++++++++++++  refreshToken 2:{} ", refreshToken);
	    OAuth2AccessToken accessToken = createAccessToken(authentication, refreshToken);
	    
	    jdbcTokenStore.storeAccessToken(accessToken, authentication);
	    // ?????? ??? ??????,
	    refreshToken = accessToken.getRefreshToken();
	    if (refreshToken != null) {
	    	jdbcTokenStore.storeRefreshToken(refreshToken, authentication);
	    }
		log.info("+++++++++++++++++  accessToken:{} ", accessToken);
	    return accessToken;
	}

	@Override
	@Transactional(noRollbackFor={InvalidTokenException.class, InvalidGrantException.class})
	public OAuth2AccessToken refreshAccessToken(String refreshTokenValue, TokenRequest tokenRequest)	throws AuthenticationException {
		log.info("*********************  [CustomTokenServices] refreshAccessToken() /{}", refreshTokenValue);
		
		if (!supportRefreshToken) {
	        throw new InvalidGrantException("Invalid refresh token1: " + refreshTokenValue);
	    }

	    OAuth2RefreshToken refreshToken = jdbcTokenStore.readRefreshToken(refreshTokenValue);
	    if (refreshToken == null) {
	        throw new InvalidGrantException("Invalid refresh token2: " + refreshTokenValue);
	    }
	    
	    OAuth2Authentication authentication = jdbcTokenStore.readAuthenticationForRefreshToken(refreshToken);
//	    log.info("---------- auth name : {}", authentication.getName());
//	    if (this.authenticationManager != null && !authentication.isClientOnly()) {
//	        // ?????????????????? ?????? ?????????????????? ????????? ????????? ??????????????? ??? ???????????? ?????? ?????? ??? ????????? ????????? ??????????????????.
//	        Authentication user = new PreAuthenticatedAuthenticationToken(authentication.getUserAuthentication(), "", authentication.getAuthorities());
//	        user = authenticationManager.authenticate(user);
//	        Object details = authentication.getDetails();
//	        authentication = new OAuth2Authentication(authentication.getOAuth2Request(), user);
//	        authentication.setDetails(details);
//	    }
	    String clientId = authentication.getOAuth2Request().getClientId();
//	    log.info("---------- clientId : {}", clientId);
	    if (clientId == null || !clientId.equals(tokenRequest.getClientId())) {
	        throw new InvalidGrantException("Wrong client for this refresh token: " + refreshTokenValue);
	    }

	    // ?????? ????????? ?????? ????????? ?????? ????????? ????????? ?????????.
	    jdbcTokenStore.removeAccessTokenUsingRefreshToken(refreshToken);

	    // ?????? ?????? ??????
	    if (isExpired(refreshToken)) {
	    	jdbcTokenStore.removeRefreshToken(refreshToken);
	        throw new InvalidTokenException("Invalid refresh token (expired): " + refreshToken);
	    }
	    
	    // ??????
	    authentication = createRefreshedAuthentication(authentication, tokenRequest);

	    // ???????????? ????????? ????????? ??????
	    if (!reuseRefreshToken) {
	    	// ????????? ?????? ?????????. (?????? -> ??????)
	    	jdbcTokenStore.removeRefreshToken(refreshToken);
	        refreshToken = createRefreshToken(authentication);
	    } else if(updateRefreshTokenLastLoginTime) {
	    	// ?????? ?????? ?????? ??????
	        String refreshTokn = refreshToken.getValue();
	        jdbcTokenStore.removeRefreshToken(refreshToken);
	        refreshToken = updateRefreshToken(authentication, refreshTokn);
	    }
	    
	    OAuth2AccessToken accessToken = createAccessToken(authentication, refreshToken);
	    jdbcTokenStore.storeAccessToken(accessToken, authentication);
	    if (!reuseRefreshToken) {
	    	jdbcTokenStore.storeRefreshToken(accessToken.getRefreshToken(), authentication);
	    } else if(updateRefreshTokenLastLoginTime) {
	    	// ?????? ?????? ??? ?????? 
	    	jdbcTokenStore.storeRefreshToken(accessToken.getRefreshToken(), authentication);
	    }
	    
	    log.info("*********************  [CustomTokenServices] refreshAccessToken() end...");
	    return accessToken;
	}

	@Override
	public OAuth2AccessToken getAccessToken(OAuth2Authentication authentication) {
		return jdbcTokenStore.getAccessToken(authentication);
	}

	@Override
	public void afterPropertiesSet() throws Exception {
		// TODO Auto-generated method stub
		
	}

	@Override
	public boolean revokeToken(String tokenValue) {
		OAuth2AccessToken accessToken = jdbcTokenStore.readAccessToken(tokenValue);
	    if (accessToken == null) {
	        return false;
	    }
	    if (accessToken.getRefreshToken() != null) {
	    	jdbcTokenStore.removeRefreshToken(accessToken.getRefreshToken());
	    }
	    jdbcTokenStore.removeAccessToken(accessToken);
	    return true;
	}
	
	/**
	 * Create a refreshed authentication.
	 * 
	 * @param authentication The authentication.
	 * @param request The scope for the refreshed token.
	 * @return The refreshed authentication.
	 * @throws InvalidScopeException If the scope requested is invalid or wider than the original scope.
	 */
	private OAuth2Authentication createRefreshedAuthentication(OAuth2Authentication authentication, TokenRequest request) {
		log.info("*********************  [CustomTokenServices] createRefreshedAuthentication() ");
	    OAuth2Authentication narrowed = authentication;
	    Set<String> scope = request.getScope();
	    OAuth2Request clientAuth = authentication.getOAuth2Request().refresh(request);
	    if (scope != null && !scope.isEmpty()) {
	        Set<String> originalScope = clientAuth.getScope();
	        if (originalScope == null || !originalScope.containsAll(scope)) {
	            throw new InvalidScopeException("Unable to narrow the scope of the client authentication to " + scope + ".", originalScope);
	        } else {
	            clientAuth = clientAuth.narrowScope(scope);
	        }
	    }
	    narrowed = new OAuth2Authentication(clientAuth, authentication.getUserAuthentication());
	    return narrowed;
	}
	
	private OAuth2RefreshToken updateRefreshToken(OAuth2Authentication authentication,String ref_tokn) {
	    if (!isSupportRefreshToken(authentication.getOAuth2Request())) {
	        return null;
	    }
	    
	    int validitySeconds = getRefreshTokenValiditySeconds(authentication.getOAuth2Request());
	    String value = ref_tokn;
	    if (validitySeconds > 0) {
	        return new DefaultExpiringOAuth2RefreshToken(value, new Date(System.currentTimeMillis() + (validitySeconds * 1000L)));
	    }
	    return new DefaultOAuth2RefreshToken(value);
	}
	
	private OAuth2AccessToken createAccessToken(OAuth2Authentication authentication, OAuth2RefreshToken refreshToken) {
		log.info("+++++++++++++++++  [CustomTokenServices] createAccessToken() 2 ");
	    DefaultOAuth2AccessToken token = new DefaultOAuth2AccessToken(UUID.randomUUID().toString());
	    int validitySeconds = getAccessTokenValiditySeconds(authentication.getOAuth2Request());
	    if (validitySeconds > 0) {
	        token.setExpiration(new Date(System.currentTimeMillis() + (validitySeconds * 1000L)));
	    }
	    token.setRefreshToken(refreshToken);
	    token.setScope(authentication.getOAuth2Request().getScope());

	    return accessTokenEnhancer != null ? accessTokenEnhancer.enhance(token, authentication) : token;
	}
	
	private OAuth2RefreshToken createRefreshToken(OAuth2Authentication authentication) {
	    if (!isSupportRefreshToken(authentication.getOAuth2Request())) {
	        return null;
	    }
	    int validitySeconds = getRefreshTokenValiditySeconds(authentication.getOAuth2Request());
	    
	    String value = UUID.randomUUID().toString();
	    if (validitySeconds > 0) {
	        return new DefaultExpiringOAuth2RefreshToken(value, new Date(System.currentTimeMillis() + (validitySeconds * 1000L)));
	    }
	    return new DefaultOAuth2RefreshToken(value);
	}
	
	/**
	 * The refresh token validity period in seconds
	 * 
	 * @param clientAuth the current authorization request
	 * @return the refresh token validity period in seconds
	 */
	protected int getRefreshTokenValiditySeconds(OAuth2Request clientAuth) {
	    if (clientDetailsService != null) {
	        ClientDetails client = clientDetailsService.loadClientByClientId(clientAuth.getClientId());
	        Integer validity = client.getRefreshTokenValiditySeconds();
	        if (validity != null) {
	        	log.info("*********************  [CustomTokenServices] getRefreshTokenValiditySeconds() validitySeconds = {}", validity);
	            return validity;
	        }
	    }
	    return refreshTokenValiditySeconds;
	}
	
	/**
	 * The access token validity period in seconds
	 * 
	 * @param clientAuth the current authorization request
	 * @return the access token validity period in seconds
	 */
	protected int getAccessTokenValiditySeconds(OAuth2Request clientAuth) {
	    if (clientDetailsService != null) {
	        ClientDetails client = clientDetailsService.loadClientByClientId(clientAuth.getClientId());
	        Integer validity = client.getAccessTokenValiditySeconds();
	        if (validity != null) {
	            return validity;
	        }
	    }
	    return accessTokenValiditySeconds;
	}
	
	/**
	 * Is a refresh token supported for this client 
	 * 
	 * @param clientAuth the current authorization request
	 * @return boolean to indicate if refresh token is supported
	 */
	protected boolean isSupportRefreshToken(OAuth2Request clientAuth) {
	    if (clientDetailsService != null) {
	        ClientDetails client = clientDetailsService.loadClientByClientId(clientAuth.getClientId());
	        return client.getAuthorizedGrantTypes().contains("refresh_token");
	    }
	    return this.supportRefreshToken;
	}

	protected boolean isExpired(OAuth2RefreshToken refreshToken) {
	    if (refreshToken instanceof ExpiringOAuth2RefreshToken) {
	        ExpiringOAuth2RefreshToken expiringToken = (ExpiringOAuth2RefreshToken) refreshToken;
	        return expiringToken.getExpiration() == null
	                || System.currentTimeMillis() > expiringToken.getExpiration().getTime();
	    }
	    return false;
	}
}
