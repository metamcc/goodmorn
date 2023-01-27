package io.mcc.auth.entity;

import lombok.Getter;

@Getter
public enum ResultCode {

    SUCCESS("0", "Success"),
    
    REQUIRED_PARAMETER_FAIL("S001", "Required parameter fail..."),
    
    AE01("AE01", "Error Occured."),
	AE02("AE02", "Invalid request to the api resource."),
	AE03("AE03", "Invalid client ID"),
	AE04("AE04", "Invalid Grant(Username,Password login fail)"),
	AE05("AE05", "Unauthorized User."),
	AE06("AE06", "InvalidTokenException, Token was not recognised."),
	AE07("AE07", "Unsupported grant type."),
	AE08("AE08", "Invalid scope request."),
	AE09("AE09", "insufficient_scope."),
	AE10("AE10", "Redirect Url Mismatch."),
	AE11("AE11", "Unsupported response type."),
	AE12("AE12", "User authorization Denied."),
	
	AE61("AE61", "Invalid ID token."),
	AE62("AE62", "Invalid Access token."),
	AE63("AE63", "No UserInfo EndPoint."),
	AE64("AE64", "Client Info Does Not Exist ."),
	
	IVCT("IVCT", "Invalid Client Access token."),
     
	NEU("NEU", "Not Exist User."),
	
	M006("M006","The account is registered through Social Media. Please Sign in through relevent Social Media account."),	
	M023("M023","The password error count has exceeded five times. Then please email authentication simple to reset your password."),
	M110("M110","Please check your Login details."),	
	M704("M704","Device ID was required on request header."),
	M701("M701","The system is currently under maintenance.")
	;

    private String code;
    private String msg;

    ResultCode(String code, String msg) {
        this.code = code;
        this.msg = msg;
    }
}
