package io.mcc.auth.common.exception;

import org.springframework.security.oauth2.common.exceptions.OAuth2Exception;

/**
 * Base exception
 * 
 * @author Ryan Heaton
 * @author Dave Syer
 */
@SuppressWarnings("serial")
public abstract class ClientAuthenticationException extends OAuth2Exception {

	public ClientAuthenticationException(String msg, Throwable t) {
		super(msg, t);
	}

	public ClientAuthenticationException(String msg) {
		super(msg);
	}

	@Override
	public int getHttpErrorCode() {
		// The spec says this is a bad request (not unauthorized)
		return 400;
	}

	@Override
	public abstract String getOAuth2ErrorCode();
}
