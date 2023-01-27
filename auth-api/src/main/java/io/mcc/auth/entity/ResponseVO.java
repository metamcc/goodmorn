package io.mcc.auth.entity;


import lombok.Data;

@Data
public class ResponseVO {
	
	
	String resultCode;
	String resultMessage;
	Object resultObject;
	
	public ResponseVO() {}
    
    public ResponseVO(String code, String message) {
    	this.resultCode = code;
        this.resultMessage = message;
    }
    
    public ResponseVO(String code, String message, String result) {
    	this.resultCode = code;
        this.resultMessage = message;
        this.resultObject = result;
    }
	
	public ResponseVO(ResultCode resultCode) {
		this.resultCode = resultCode.getCode();
		this.resultMessage = resultCode.getMsg();
	}
	
	public void setResponseStatus(final ResultCode resultCode) {
		this.resultCode = resultCode.getCode();
		this.resultMessage = resultCode.getMsg();	
	}
	
}
