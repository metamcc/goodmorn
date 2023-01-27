package io.mcc.auth.common.exception;

import org.springframework.security.oauth2.common.exceptions.OAuth2Exception;

import com.fasterxml.jackson.databind.annotation.JsonSerialize;

@JsonSerialize(using = CustomOauthExceptionSerializer.class)
public class CustomOauthException22 extends OAuth2Exception {
    public CustomOauthException22(String msg) {
        super(msg);
    }
}
