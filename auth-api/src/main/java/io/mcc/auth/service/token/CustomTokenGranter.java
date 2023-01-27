package io.mcc.auth.service.token;

import java.util.List;
import java.util.Map;

import lombok.extern.slf4j.Slf4j;
import org.springframework.security.authentication.UsernamePasswordAuthenticationToken;
import org.springframework.security.core.Authentication;
import org.springframework.security.core.GrantedAuthority;
import org.springframework.security.core.authority.AuthorityUtils;
import org.springframework.security.oauth2.common.OAuth2AccessToken;
import org.springframework.security.oauth2.common.util.OAuth2Utils;
import org.springframework.security.oauth2.provider.ClientDetails;
import org.springframework.security.oauth2.provider.ClientDetailsService;
import org.springframework.security.oauth2.provider.OAuth2Authentication;
import org.springframework.security.oauth2.provider.OAuth2RequestFactory;
import org.springframework.security.oauth2.provider.TokenRequest;
import org.springframework.security.oauth2.provider.token.AbstractTokenGranter;
import org.springframework.security.oauth2.provider.token.AuthorizationServerTokenServices;

@Slf4j
public class CustomTokenGranter extends AbstractTokenGranter {

	private static final String GRANT_TYPE = "refresh_token";
	
	public CustomTokenGranter(AuthorizationServerTokenServices tokenServices, ClientDetailsService clientDetailsService,
			OAuth2RequestFactory requestFactory, String grantType) {
		super(tokenServices, clientDetailsService, requestFactory, grantType);
	}

	protected OAuth2Authentication getOAuth2Authentication(ClientDetails client, TokenRequest tokenRequest) {
		Map<String, String> params = tokenRequest.getRequestParameters();
		String username = params.containsKey("username") ? params.get("username") : "guest";

		log.debug("getOAuth2Authentication:{}", username);

		List<GrantedAuthority> authorities = params.containsKey("authorities") ? AuthorityUtils
				.createAuthorityList(OAuth2Utils.parseParameterList(params.get("authorities")).toArray(new String[0]))
				: AuthorityUtils.NO_AUTHORITIES;
		Authentication user = new UsernamePasswordAuthenticationToken(username, "N/A", authorities);
		OAuth2Authentication authentication = new OAuth2Authentication(tokenRequest.createOAuth2Request(client), user);
		return authentication;
	}
	
	@Override
	protected OAuth2AccessToken getAccessToken(ClientDetails client, TokenRequest tokenRequest) {
	    String refreshToken = tokenRequest.getRequestParameters().get("refresh_token");

		log.debug("getAccessToken:refreshToken={}", refreshToken);

	    return getTokenServices().refreshAccessToken(refreshToken, tokenRequest);
	}


}
