package io.mcc.auth.common.config;

import java.util.ArrayList;
import java.util.Arrays;
import java.util.HashSet;
import java.util.List;
import java.util.Locale;
import java.util.Set;

import javax.sql.DataSource;

import io.mcc.auth.service.CommonService;
import io.mcc.auth.service.LoggingService;
import io.mcc.auth.service.token.CustomTokenEnhancer;
import io.mcc.auth.service.token.CustomTokenGranter;
import io.mcc.auth.service.token.CustomTokenServices;
import org.apache.commons.lang3.StringUtils;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.autoconfigure.web.ServerProperties;
import org.springframework.context.MessageSource;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.core.env.Environment;
import org.springframework.http.HttpHeaders;
import org.springframework.http.ResponseEntity;
import org.springframework.security.authentication.AuthenticationManager;
import org.springframework.security.core.Authentication;
import org.springframework.security.core.GrantedAuthority;
import org.springframework.security.core.authority.SimpleGrantedAuthority;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.security.core.userdetails.UserDetailsService;
import org.springframework.security.crypto.bcrypt.BCryptPasswordEncoder;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.security.oauth2.common.exceptions.OAuth2Exception;
import org.springframework.security.oauth2.config.annotation.configurers.ClientDetailsServiceConfigurer;
import org.springframework.security.oauth2.config.annotation.web.configuration.AuthorizationServerConfigurerAdapter;
import org.springframework.security.oauth2.config.annotation.web.configuration.EnableAuthorizationServer;
import org.springframework.security.oauth2.config.annotation.web.configurers.AuthorizationServerEndpointsConfigurer;
import org.springframework.security.oauth2.config.annotation.web.configurers.AuthorizationServerSecurityConfigurer;
import org.springframework.security.oauth2.provider.ClientDetails;
import org.springframework.security.oauth2.provider.ClientDetailsService;
import org.springframework.security.oauth2.provider.ClientRegistrationException;
import org.springframework.security.oauth2.provider.CompositeTokenGranter;
import org.springframework.security.oauth2.provider.TokenGranter;
import org.springframework.security.oauth2.provider.client.BaseClientDetails;
import org.springframework.security.oauth2.provider.error.DefaultOAuth2ExceptionRenderer;
import org.springframework.security.oauth2.provider.error.OAuth2AccessDeniedHandler;
import org.springframework.security.oauth2.provider.error.OAuth2AuthenticationEntryPoint;
import org.springframework.security.oauth2.provider.error.WebResponseExceptionTranslator;
import org.springframework.security.oauth2.provider.token.AccessTokenConverter;
import org.springframework.security.oauth2.provider.token.TokenEnhancer;
import org.springframework.security.oauth2.provider.token.TokenStore;
import org.springframework.security.oauth2.provider.token.store.JdbcTokenStore;
import org.springframework.security.web.authentication.WebAuthenticationDetails;

import io.mcc.auth.common.config.handler.CustomAccessDeniedHandler;
import io.mcc.auth.common.config.handler.CustomAuthenticationEntryPoint;
import io.mcc.auth.common.config.interceptor.LoggingInterceptor;
import io.mcc.auth.common.exception.CustomOAuthExceptionRenderer;
import io.mcc.auth.common.exception.CustomOauthException;
import io.mcc.auth.entity.ResultCode;

import lombok.extern.slf4j.Slf4j;

import org.springframework.security.oauth2.provider.error.DefaultWebResponseExceptionTranslator;

@Slf4j
@Configuration
@EnableAuthorizationServer
public class OAuth2ServerConfiguration extends AuthorizationServerConfigurerAdapter{

	@Autowired
    private LoggingService loggingService;

    @Autowired
    private CommonService commonService;
    
    @Autowired
	private MessageSource messageSource;
	
	@Autowired
    private DataSource datasource;
	
	@Autowired
    private AuthenticationManager authenticationManager;
	
	@Autowired
    private AccessTokenConverter accessTokenConverter;
	
	@Autowired
    private UserDetailsService userDetailsService;
	
	@Autowired
    private Environment env;
	
    /**
     * Client 정보
     * application.yml 
     */
	@Bean
	public ClientDetailsService clientDetailsService() {
//        return new JdbcClientDetailsService(dataSource);
        return new ClientDetailsService() {
        	@Override
            public ClientDetails loadClientByClientId(String clientId) throws ClientRegistrationException {
        		
        		BaseClientDetails details = new BaseClientDetails();
                details.setClientId(clientId);
        		
        		String detailStr = env.getProperty("client.details."  + clientId);
        		
        		if(!StringUtils.isEmpty(detailStr)) {        		
	        		String[] detailArr = detailStr.split("::");
	        		for(String keyValStr : detailArr) {
	        			
	        			String[]  keyValArr = keyValStr.split("=");
	        			if("resourceIds".equals(keyValArr[0])) {
	        				details.setResourceIds(Arrays.asList(keyValArr[1]));
	        			} else if("clientSecret".equals(keyValArr[0])) {
	        				details.setClientSecret(keyValArr[1]);
	        			} else if("scope".equals(keyValArr[0])) {
	        				
	        				String[] scopeArr = keyValArr[1].split(",");
	        				List scopeList = Arrays.asList(scopeArr);
	        				
	        				details.setScope(scopeList);
	        			} else if("authorizedGrantTypes".equals(keyValArr[0])) {
	        				
	        				String[] grantArr = keyValArr[1].split(",");
	        				List grantList = Arrays.asList(grantArr);
	        				
	        				details.setAuthorizedGrantTypes(grantList );			
	        			} else if("authorities".equals(keyValArr[0])) {
	        				Set<GrantedAuthority> authorities = new HashSet<GrantedAuthority>();
	        				authorities.add(new SimpleGrantedAuthority(keyValArr[1]));
	        				details.setAuthorities(authorities);
	        			} else if("accessTokenValidity".equals(keyValArr[0])) {
	        				details.setAccessTokenValiditySeconds(Integer.parseInt(keyValArr[1]));
	        			} else if("refreshTokenValidity".equals(keyValArr[0])) {
	        				 details.setRefreshTokenValiditySeconds(Integer.parseInt(keyValArr[1]));
	        			} else {
	        				log.info(">>>>>>>>>>>>    keyValStr={}", keyValStr); 
	        			}
	        			
	        		}
        		}
                
                return details;
            }
        };
    }
	
    @Override
    public void configure(ClientDetailsServiceConfigurer clients) throws Exception {
        clients.withClientDetails(clientDetailsService());
    }
    
    @Bean
    public TokenStore jdbcTokenStore(DataSource dataSource) {
        return new JdbcTokenStore(dataSource); 
    }
    
    @Bean
    public LoggingInterceptor loggingInterceptor() {
//        return new LoggingInterceptor(); 
        return new LoggingInterceptor(loggingService, commonService);
    }
    
	@Bean
	public WebResponseExceptionTranslator customExceptionTranslator() {
		return new DefaultWebResponseExceptionTranslator() {
	        @Override
	        public ResponseEntity<OAuth2Exception> translate(Exception e) throws Exception {
	        	
	            ResponseEntity<OAuth2Exception> responseEntity = super.translate(e);
	            OAuth2Exception body = responseEntity.getBody();
	            HttpHeaders headers = new HttpHeaders();
	            headers.setAll(responseEntity.getHeaders().toSingleValueMap());

	           // translate the exception
	            int status = body.getHttpErrorCode();
	            String exCode = body.getOAuth2ErrorCode();
	            String exMsg = body.getMessage();
	            String exLocaleCode = "";
	            String exLocaleMsg = "";
	    		
	    		//Locale locale = new Locale("en");
	            //언어정보 획득
	            Authentication auth = SecurityContextHolder.getContext().getAuthentication();
	    		WebAuthenticationDetails details = (WebAuthenticationDetails) auth.getDetails();
	    		String ssionid = details.toString();
	    		Locale locale = null;
	    		
	    		try {
	    			String ln      = ssionid == null || "".equals(ssionid) ? "en" : System.getProperty(ssionid); 
		    		System.clearProperty(ssionid);
		    		locale = new Locale(ln);	
	    		} catch (Exception l) {
	    			locale = new Locale("en");
	    		}
	    		
//	    		log.info(">>>>>>>>>>>>>>>>>     OAuth2Exception code : {}", exCode);
//	    		log.info(">>>>>>>>>>>>>>>>>     OAuth2Exception mag : {}", exMsg);
	    		
	    		if("invalid_request".equals(exCode)) {
	    			exLocaleCode = "AE02";
	    			exLocaleMsg = messageSource.getMessage(ResultCode.AE02.getCode(), null, locale);
	    		} else if("invalid_client".equals(exCode)) {
	    			exLocaleCode = "AE03";
	    			exLocaleMsg = messageSource.getMessage(ResultCode.AE03.getCode(), null, locale);
	    		} else if("invalid_grant".equals(exCode)){
	    			exLocaleCode = "AE04";
	    			exLocaleMsg = messageSource.getMessage(ResultCode.AE04.getCode(), null, locale);
	    		} else if("unauthorized_client".equals(exCode)){
	    			exLocaleCode = "AE05";
	    			exLocaleMsg = messageSource.getMessage(ResultCode.AE05.getCode(), null, locale);
	    		} else if("invalid_token".equals(exCode)){
	    			exLocaleCode = "AE06";
	    			exLocaleMsg = messageSource.getMessage(ResultCode.AE06.getCode(), null, locale);
	    		} else if("unsupported_grant_type".equals(exCode)){
	    			exLocaleCode = "AE07";
	    			exLocaleMsg = messageSource.getMessage(ResultCode.AE07.getCode(), null, locale);
	    		} else if("invalid_scope".equals(exCode)){
	    			exLocaleCode = "AE08";
	    			exLocaleMsg = messageSource.getMessage(ResultCode.AE08.getCode(), null, locale);
	    		} else if("insufficient_scope".equals(exCode)){
	    			exLocaleCode = "AE09";
	    			exLocaleMsg = messageSource.getMessage(ResultCode.AE09.getCode(), null, locale);
	    		} else if("redirect_uri_mismatch".equals(exCode)){
	    			exLocaleCode = "AE10";
	    			exLocaleMsg = messageSource.getMessage(ResultCode.AE10.getCode(), null, locale);
	    		} else if("unsupported_response_type".equals(exCode)){
	    			exLocaleCode = "AE11";
	    			exLocaleMsg = messageSource.getMessage(ResultCode.AE11.getCode(), null, locale);
	    		} else if("access_denied".equals(exCode)){
	    			exLocaleCode = "AE12";
	    			exLocaleMsg = messageSource.getMessage(ResultCode.AE12.getCode(), null, locale);
	    		} else {
	    			exLocaleCode = "AE01";
	    			exLocaleMsg = messageSource.getMessage(ResultCode.AE01.getCode(), null, locale);
	    		}
	    		
	            CustomOauthException ex = new CustomOauthException(body.getMessage());
	            ex.addAdditionalInformation("code", body.getOAuth2ErrorCode());
	            ex.addAdditionalInformation("msg", body.getMessage());
				ex.addAdditionalInformation("resultCode", exLocaleCode);
				ex.addAdditionalInformation("resultMessage", exLocaleMsg);
	            
	            return new ResponseEntity<>(ex, headers, responseEntity.getStatusCode());
	        }
	    };
	}
	
	@Bean
	public OAuth2AccessDeniedHandler accessDeniedHandler(){
		return new CustomAccessDeniedHandler();
	}
	
	@Bean
	public OAuth2AuthenticationEntryPoint authenticationEntryPoint(){
		
		CustomAuthenticationEntryPoint ep = new CustomAuthenticationEntryPoint();
//		ep.setExceptionTranslator(customExceptionTranslator());
		ep.setExceptionRenderer(oAuthExceptionRenderer());
		ep.setRealmName("oauth2/client");
		
		return ep;
	}
	
	@Bean
    public PasswordEncoder passwordEncoder() {
        return new BCryptPasswordEncoder();    	
//		return new BCryptPasswordEncoder(8);
    }
	
	@Bean
    public DefaultOAuth2ExceptionRenderer oAuthExceptionRenderer() {
		return new CustomOAuthExceptionRenderer();
	}
	
	@Bean
	public TokenEnhancer accessTokenEnhancer() {
	    return new CustomTokenEnhancer();
	}
	
	@Bean
	public CustomTokenServices createCustomTokenServices() {
		CustomTokenServices service = new CustomTokenServices();
		return service;
	} 
	

	@Autowired
    private ServerProperties server;

	
    /**
     * 인증 서버 token 발급!
     */
    @Override
    public void configure(AuthorizationServerEndpointsConfigurer endpoints)throws Exception {
    	
        endpoints
//		        .pathMapping("/oauth/authorize", "/api/v1/oauth/authorize")// here 
//		        .pathMapping("/oauth/check_token", "/api/v1/oauth/check_token")// here 
//		        .pathMapping("/oauth/confirm_access", "/api/v1/oauth/confirm_access")// here 
//		        .pathMapping("/oauth/error", "/api/v1/oauth/error")// here 
//		        .pathMapping("/oauth/token", "/proc/oauth/token")
                .tokenStore(jdbcTokenStore(datasource))
                .authenticationManager(authenticationManager)
                .tokenGranter(tokenGranter(endpoints))
                .accessTokenConverter(accessTokenConverter)
                .userDetailsService(userDetailsService)
                .tokenServices(createCustomTokenServices())
                .approvalStoreDisabled()                
                .exceptionTranslator(customExceptionTranslator())
//                .exceptionTranslator(exception  -> {
//                	if (exception instanceof OAuth2Exception) {
//                        OAuth2Exception oAuth2Exception = (OAuth2Exception) exception;
//                        return ResponseEntity
//                                .status(oAuth2Exception.getHttpErrorCode())
//                                .body(new CustomOauthException22(oAuth2Exception.getMessage()));
//                    } else {
//                        throw exception;
//                    }
//                })
                .addInterceptor(loggingInterceptor())
        ;

    }
    
    /**
     * 토큰을 확인
     * Basic realm="oauth2/client"
     */
    @Override
    public void configure(AuthorizationServerSecurityConfigurer security) throws Exception{
        security
                .tokenKeyAccess("permitAll()")
                .checkTokenAccess("isAuthenticated()")
                .realm("oauth2/client")
                .accessDeniedHandler(accessDeniedHandler())
                .authenticationEntryPoint(authenticationEntryPoint())
                .allowFormAuthenticationForClients()
                .passwordEncoder(passwordEncoder())// Token 정보를 API(/oauth/check_token)를 활성화 시킨다.
                ;
    }
    
    private TokenGranter tokenGranter(final AuthorizationServerEndpointsConfigurer endpoints) {
        List<TokenGranter> granters = new ArrayList<TokenGranter>(Arrays.asList(endpoints.getTokenGranter()));
        granters.add(new CustomTokenGranter(endpoints.getTokenServices(), endpoints.getClientDetailsService(), endpoints.getOAuth2RequestFactory(), "custom"));
        return new CompositeTokenGranter(granters);
    }
}
