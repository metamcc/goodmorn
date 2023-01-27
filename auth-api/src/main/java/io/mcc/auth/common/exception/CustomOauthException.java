package io.mcc.auth.common.exception;

import org.springframework.security.oauth2.common.exceptions.OAuth2Exception;

import org.springframework.security.oauth2.common.exceptions.*;

@org.codehaus.jackson.map.annotate.JsonSerialize(using = OAuth2ExceptionJackson1Serializer.class)
@org.codehaus.jackson.map.annotate.JsonDeserialize(using = OAuth2ExceptionJackson1Deserializer.class)
@com.fasterxml.jackson.databind.annotation.JsonSerialize(using = OAuth2ExceptionJackson2Serializer.class)
@com.fasterxml.jackson.databind.annotation.JsonDeserialize(using = OAuth2ExceptionJackson2Deserializer.class)
public class CustomOauthException extends OAuth2Exception {
	
	private String oaut_error_code;
  private int http_error_code;

  public CustomOauthException(String msg, Throwable t) {
      super(msg, t);
  }

  public CustomOauthException(String msg) {
      super(msg);
  }


  public void setOaut_error_code(String oaut_error_code) {
      this.oaut_error_code = oaut_error_code;
  }

  public void setHttp_error_code(int http_error_code) {
      this.http_error_code = http_error_code;
  }

	
}
