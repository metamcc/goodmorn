package io.mcc.auth.common.config;

import java.util.Arrays;
import java.util.HashSet;
import java.util.List;
import java.util.Map;
import java.util.Set;
import java.util.stream.Collectors;

import javax.sql.DataSource;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.core.env.Environment;
import org.springframework.security.authentication.AuthenticationManager;
import org.springframework.security.config.annotation.authentication.builders.AuthenticationManagerBuilder;
import org.springframework.security.config.annotation.web.builders.HttpSecurity;
import org.springframework.security.config.annotation.web.builders.WebSecurity;
import org.springframework.security.config.annotation.web.configuration.EnableWebSecurity;
import org.springframework.security.config.annotation.web.configuration.WebSecurityConfigurerAdapter;
import org.springframework.security.config.oauth2.client.CommonOAuth2Provider;
import org.springframework.security.core.GrantedAuthority;
import org.springframework.security.core.authority.mapping.GrantedAuthoritiesMapper;
import org.springframework.security.oauth2.client.registration.ClientRegistration;
import org.springframework.security.oauth2.client.registration.ClientRegistrationRepository;
import org.springframework.security.oauth2.client.registration.InMemoryClientRegistrationRepository;
import org.springframework.security.oauth2.core.oidc.OidcIdToken;
import org.springframework.security.oauth2.core.oidc.OidcUserInfo;
import org.springframework.security.oauth2.core.oidc.user.OidcUserAuthority;
import org.springframework.security.oauth2.core.user.OAuth2UserAuthority;
import org.springframework.security.oauth2.provider.error.OAuth2AccessDeniedHandler;
import org.springframework.security.oauth2.provider.token.AccessTokenConverter;
import org.springframework.security.oauth2.provider.token.DefaultAccessTokenConverter;
import org.springframework.security.web.AuthenticationEntryPoint;

import lombok.extern.slf4j.Slf4j;

@Slf4j
@Configuration
@EnableWebSecurity
//@EnableOAuth2Client
public class WebSecurityConfiguration extends WebSecurityConfigurerAdapter {
	
	private static List<String> clients = Arrays.asList("google", "facebook", "kakao");
	
	private static String CLIENT_PROPERTY_KEY = "spring.security.oauth2.client.registration.";
	private static String KAKAO_REGISTRATION = "spring.security.oauth2.client.registration.kakao.";

	@Autowired 
	private DataSource dataSource;
	
	@Autowired
    private Environment env;
	
	@Autowired
	private OAuth2AccessDeniedHandler accessDeniedHandler;
	
	@Autowired
	private AuthenticationEntryPoint authenticationEntryPoint;
	
	@Bean
	public AuthenticationManager authenticationManager() throws Exception {
		return super.authenticationManager();
	}
	
	@Bean
	public AccessTokenConverter accessTokenConverter() throws Exception{
        return new DefaultAccessTokenConverter();
    }
	
//	@Bean
//	public FilterRegistrationBean<CustomFilter> loggingFilter(){
//	    FilterRegistrationBean<CustomFilter> registrationBean = new FilterRegistrationBean<>();
//	         
//	    registrationBean.setFilter(new CustomFilter());
//	    registrationBean.addUrlPatterns("/oauth/*");
//	         
//	    return registrationBean;    
//	}
	
	@Override 
	protected void configure(final AuthenticationManagerBuilder auth) throws Exception { 
		auth.jdbcAuthentication().dataSource(dataSource)
				.usersByUsernameQuery("select username, password, enabled from t_user where username = ?")		
				.authoritiesByUsernameQuery("select username, role from t_user_roles where username = ?")
				;
	}
	
	@Override
	public void configure(HttpSecurity http) throws Exception {
		
		http
    		.headers()
        		.frameOptions().sameOrigin()				//X-Frame-Options 
        		.httpStrictTransportSecurity().disable()	//HSTS
        		.cacheControl()								//no-cache(필수)
        ;
		
		http.csrf().disable();
		
		http
            .authorizeRequests()
	            .antMatchers("/**").permitAll()					//인증제외
	            .antMatchers("/revoketoken").permitAll()	
	            .antMatchers("/oauth/**").permitAll()
	            .antMatchers("/client/**").permitAll()
	            .antMatchers("/proc/**").permitAll()
	            .antMatchers("/swagger**").permitAll()
	            .antMatchers("/health").permitAll()	            
	            .anyRequest().authenticated()
	            .and()  
	            	.exceptionHandling()
	            		.accessDeniedHandler(accessDeniedHandler)
	            		.authenticationEntryPoint(authenticationEntryPoint)
            ;
	}
	
	/**
	 * Client 설정
	 * 
	 * 
	 */
	@Bean
	public ClientRegistrationRepository clientRegistrationRepository() {
        List<ClientRegistration> registrations = clients.stream()
            .map(c -> getRegistration(c))
            .filter(registration -> registration != null)
            .collect(Collectors.toList());
        
        registrations.add(CustomOAuth2Provider.KAKAO.getBuilder("kakao")
        	.clientId(env.getProperty(KAKAO_REGISTRATION  + "client-id"))
        	.clientSecret(env.getProperty(KAKAO_REGISTRATION  + "client-secret"))
        	.build());
        
        return new InMemoryClientRegistrationRepository(registrations);
    }
	
	@Bean
	public GrantedAuthoritiesMapper userAuthoritiesMapper() {
		return (authorities) -> {
			Set<GrantedAuthority> mappedAuthorities = new HashSet<>();

			authorities.forEach(authority -> {
				if (OidcUserAuthority.class.isInstance(authority)) {
					OidcUserAuthority oidcUserAuthority = (OidcUserAuthority)authority;

					OidcIdToken idToken = oidcUserAuthority.getIdToken();
					OidcUserInfo userInfo = oidcUserAuthority.getUserInfo();

					// Map the claims found in idToken and/or userInfo
					// to one or more GrantedAuthority's and add it to mappedAuthorities

				} else if (OAuth2UserAuthority.class.isInstance(authority)) {
					OAuth2UserAuthority oauth2UserAuthority = (OAuth2UserAuthority)authority;

					Map<String, Object> userAttributes = oauth2UserAuthority.getAttributes();

					// Map the attributes found in userAttributes
					// to one or more GrantedAuthority's and add it to mappedAuthorities
				}
			});

			return mappedAuthorities;
		};
	}
	
	@Override
	public void configure(WebSecurity web) throws Exception {
		web.ignoring().antMatchers("/v2/api-docs", "/configuration/ui", "/swagger-resources", "/configuration/security", "/swagger-ui.html", "/webjars/**");
	}

	private ClientRegistration getRegistration(String client) {
		String clientId = env.getProperty(CLIENT_PROPERTY_KEY + client + ".client-id");

        if (clientId == null) {
            return null;
        }

        String clientSecret = env.getProperty(CLIENT_PROPERTY_KEY + client + ".client-secret");
        
        if (client.equals("google")) {
			return CommonOAuth2Provider.GOOGLE.getBuilder(client)
                .clientId(clientId)
                .clientSecret(clientSecret)
                .build();
        }
        
        if (client.equals("facebook")) {
            return CommonOAuth2Provider.FACEBOOK.getBuilder(client)            		
                .clientId(clientId)
                .clientSecret(clientSecret)
//                .userInfoUri("")
                .build();
        }
        
		return null;
	}
}
