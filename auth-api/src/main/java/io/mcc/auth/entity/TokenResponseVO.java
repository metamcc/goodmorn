package io.mcc.auth.entity;


import lombok.Data;

@Data
public class TokenResponseVO extends ResponseVO{
	
	String access_token;
	String token_type;
	String refresh_token;
	int expires_in;
	String scope;

	public TokenResponseVO(String code, String message) {
		this.setResultCode(code);
		this.setResultMessage(message);
	}
	
}
