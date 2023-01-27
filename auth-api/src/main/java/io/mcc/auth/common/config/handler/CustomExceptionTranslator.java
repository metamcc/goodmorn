package io.mcc.auth.common.config.handler;

import org.springframework.http.HttpHeaders;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.security.oauth2.common.OAuth2AccessToken;
import org.springframework.security.oauth2.common.exceptions.InsufficientScopeException;
import org.springframework.security.oauth2.common.exceptions.OAuth2Exception;
import org.springframework.security.oauth2.provider.error.DefaultWebResponseExceptionTranslator;
import org.springframework.stereotype.Component;

import io.mcc.auth.common.exception.CustomOauthException;

import lombok.extern.slf4j.Slf4j;

@Slf4j
@Component(value = "customExceptionTranslator")
public class CustomExceptionTranslator extends DefaultWebResponseExceptionTranslator  {
	
	@Override
	public ResponseEntity<OAuth2Exception> translate(Exception e) throws Exception {
		log.info("**********  CustomExceptionTranslator.translate() START");
		
		ResponseEntity<OAuth2Exception> responseEntity = super.translate(e);
        OAuth2Exception body = responseEntity.getBody();
        
        
		int status = body.getHttpErrorCode();
		
		
		HttpHeaders headers = new HttpHeaders();
		headers.setAll(responseEntity.getHeaders().toSingleValueMap());
		
		
		CustomOauthException coe = new CustomOauthException( body.getMessage() );
		coe.setHttp_error_code(400);
		coe.setOaut_error_code(body.getOAuth2ErrorCode());
		
		
		log.info("**********  Summary = {}", coe.getSummary());
		
		if (status == HttpStatus.UNAUTHORIZED.value() || (e instanceof InsufficientScopeException)) {
			headers.set("WWW-Authenticate", String.format("%s %s", OAuth2AccessToken.BEARER_TYPE, coe.getSummary()));
		}
		
		responseEntity.status(responseEntity.getStatusCode()).body(new CustomOauthException(body.getMessage()));
		
		return new ResponseEntity<>(coe, headers, responseEntity.getStatusCode());		
		
	}

}
