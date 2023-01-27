package io.mcc.auth.controller;

import java.io.IOException;
import java.security.Principal;
import java.util.Base64;
import java.util.Locale;
import java.util.Map;
import java.util.Base64.Decoder;

import javax.servlet.http.HttpServletRequest;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.context.MessageSource;
import org.springframework.http.HttpHeaders;
import org.springframework.http.MediaType;
import org.springframework.mobile.device.Device;
import org.springframework.mobile.device.DevicePlatform;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.security.oauth2.common.OAuth2AccessToken;
import org.springframework.security.oauth2.common.exceptions.OAuth2Exception;
import org.springframework.security.oauth2.provider.endpoint.TokenEndpoint;
import org.springframework.util.StringUtils;
import org.springframework.web.HttpRequestMethodNotSupportedException;
import org.springframework.web.bind.annotation.RequestHeader;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestMethod;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;

import io.mcc.auth.common.util.HttpParameterUtil;
import io.mcc.auth.entity.CommonVO;
import io.mcc.auth.entity.ResponseVO;
import io.mcc.auth.entity.ResultCode;
import io.mcc.auth.entity.TokenResponseVO;
import io.mcc.auth.entity.User;
import io.mcc.auth.service.ClientService;
import io.mcc.auth.service.CommonService;
import io.mcc.auth.service.LoggingService;

import lombok.extern.slf4j.Slf4j;

@Slf4j
@RestController
@RequestMapping("/proc")
public class AuthProcessController {
	public static final Logger logger = LoggerFactory.getLogger(AuthProcessController.class);
	
	@Autowired
	private CommonService commonService;
	
	@Autowired
	private ClientService clientService;
	
	@Autowired
	private LoggingService loggingService;
	
	@Autowired
	private MessageSource messageSource;
	
	@Autowired
	private PasswordEncoder passwordEncoder;
	
	@Autowired
    private TokenEndpoint tokenEndpoint;
	
	@RequestMapping(value = "/oauth/login", 
					method = RequestMethod.POST, 
					consumes = {MediaType.APPLICATION_FORM_URLENCODED_VALUE})
    public ResponseVO postAccessToken(
    				@RequestParam Map<String, String> params, 
    				HttpServletRequest request,
    				@RequestHeader HttpHeaders reqHeader,
    				Locale locale, Device device) throws HttpRequestMethodNotSupportedException, IOException {
	
		log.info("===== /proc/oauth/login postAccessToken() IN {} ", params);
		
		ResponseVO res = new ResponseVO(ResultCode.SUCCESS);
        
        if(StringUtils.isEmpty(params.get("grant_type")) 
        		|| StringUtils.isEmpty(params.get("username")) 
        		|| StringUtils.isEmpty(params.get("password")) ) {
        	return new ResponseVO(ResultCode.REQUIRED_PARAMETER_FAIL.getCode(), messageSource.getMessage(ResultCode.REQUIRED_PARAMETER_FAIL.getCode(), null, locale));
        }
        
        String grantType = params.get("grant_type");
        String username = params.get("username");
        String password = params.get("password");
        
        String headerAuthCred = request.getHeader(HttpHeaders.AUTHORIZATION);
        String headerDeviceId = request.getHeader("device-id");
               
        if(StringUtils.isEmpty(headerAuthCred) 
        		|| StringUtils.isEmpty(headerDeviceId) ) {
        	return new ResponseVO(ResultCode.REQUIRED_PARAMETER_FAIL.getCode(), messageSource.getMessage(ResultCode.REQUIRED_PARAMETER_FAIL.getCode(), null, locale));
        }
        
      	//-Get Client Account from HttpHeaders
        String[] clientInfo = null;
        if(!StringUtils.isEmpty(headerAuthCred) && headerAuthCred.contains("Basic")) {
        	String cred = headerAuthCred.split("Basic ")[1];
			
			Decoder decoder = Base64.getDecoder(); 
			byte[] decodedBytes = decoder.decode(cred);
			clientInfo = new String(decodedBytes).split(":");
        } else {
        	return new ResponseVO(ResultCode.AE64.getCode(), messageSource.getMessage(ResultCode.AE64.getCode(), null, locale));
        }
        
		//- 회원 정보 조회
		CommonVO userInfo = null;
		userInfo = loggingService.getConnentUserInfo(username);
		if(userInfo==null || userInfo.isEmpty()) {
			return new ResponseVO(ResultCode.NEU.getCode(), messageSource.getMessage(ResultCode.NEU.getCode(), null, locale));
		}
		
		//- 회원 상태 코드 확인
		String userStatusCd = userInfo.getString("user_status_cd");
		log.info("----- userStatusCd : {}", userStatusCd);
		if(!"01".equals(userStatusCd)) {
			if("02".equals(userStatusCd)) {
				//휴면 계정입니다.
				return new ResponseVO(ResultCode.M110.getCode(), messageSource.getMessage(ResultCode.M110.getCode(), null, locale));
			} else if("04".equals(userStatusCd) || "05".equals(userStatusCd)) {
				//탈퇴 처리 된 회원 입니다.
				return new ResponseVO(ResultCode.M110.getCode(), messageSource.getMessage(ResultCode.M110.getCode(), null, locale));
			}
		}
		
		//sns user 일반 로그인 차단
		String snsUserYn = userInfo.getString("snsUserYn");
		
		if( StringUtils.pathEquals(snsUserYn, "Y")){
			String msg = messageSource.getMessage(ResultCode.M006.getCode(), null, locale);
			return new ResponseVO(ResultCode.M006.getCode(), msg);
		}
		
		//- device 변경 여부 확인
		String existDeviceId = userInfo.getString("uuid");
		
		if(!StringUtils.isEmpty(headerDeviceId) && !headerDeviceId.equals(existDeviceId)) {
			int rst = loggingService.deleteAccessToken(username);
				rst += commonService.updateUserDevice(userInfo.getLong("user_seq"), headerDeviceId);
		}
		
		//디바이스 정보 등록 및 업데이트
		//신규 유저 및 기존 유저의 새로운 로그인시 디바이스 정보를 업데이트 합니다.
		//Device device = DeviceUtils.getCurrentDevice(request);
		DevicePlatform platform = device.getDevicePlatform();
		log.info("===== /proc/oauth/login DevicePlatform name {}, os {}", platform.name(), (platform == DevicePlatform.IOS) ); 
		log.info("===== /proc/oauth/login DevicePlatform isheaderDeviceId {}, isexistDeviceId {}",!StringUtils.isEmpty(headerDeviceId), !headerDeviceId.equals(existDeviceId) ); 
		log.info("===== /proc/oauth/login DevicePlatform headerDeviceId {}, existDeviceId {}", headerDeviceId, existDeviceId); 
		commonService.updateUserDevicePlatform(userInfo.getLong("user_seq"), HttpParameterUtil.getOsVer(request), HttpParameterUtil.getDiviceName(request), headerDeviceId, platform.name());//유저 일련번호 , 
		
		
		//- conect_error_cnt 확인
		long err_cnt = userInfo.getLong("conect_error_cnt");
		if(err_cnt >= 5) {
			String msg = messageSource.getMessage(ResultCode.M023.getCode(), null, locale);
			return new ResponseVO(ResultCode.M023.getCode(), msg);
		}
        
		//- Connect Log
		CommonVO logParam = new CommonVO();
		logParam.put("conect_trmnl", HttpParameterUtil.getUserAgentInfo(request, "device"));
		logParam.put("conect_os", HttpParameterUtil.getUserAgentInfo(request, "os"));
		logParam.put("conect_brwsr", HttpParameterUtil.getUserAgentInfo(request, "brwsr"));
		
		
		//- password 검증

		String existPw = userInfo.getString("password");
		if(!passwordEncoder.matches(password, existPw)) {
			logParam.put("hist_type_cd", "02");	//로그인 오류
			
			//접속 오류 횟수 ++
			loggingService.increaseConectErrorCnt(username);
			loggingService.insertConecHistory(logParam);
			
			String msg = messageSource.getMessage(ResultCode.AE04.getCode(), null, locale);
			return new ResponseVO(ResultCode.AE04.getCode(), msg);
		}
		
        User param = new User();
		param.setUsername(username);
		param.setPassword(password);
		

		// MCC 토큰 발행
		try {
			CommonVO resVo = clientService.getResourceOwnerPasswordToken(clientInfo, param);

	        final OAuth2AccessToken token = (OAuth2AccessToken) resVo.get("accessToken");
	        
	        if (token==null || token.isExpired() ) {
	        	String msg = messageSource.getMessage(ResultCode.AE12.getCode(), null, locale);
				return new ResponseVO(ResultCode.AE12.getCode(), msg);
	        }
	        
	        TokenResponseVO tokenRes = new TokenResponseVO(ResultCode.SUCCESS.getCode(), ResultCode.SUCCESS.getMsg());
	        tokenRes.setAccess_token(token.getValue());
	        tokenRes.setToken_type(token.getTokenType());
	        tokenRes.setRefresh_token(token.getRefreshToken().getValue());
	        tokenRes.setExpires_in(token.getExpiresIn());
	        tokenRes.setScope(token.getScope().toString());
	        
	        log.info("===== /proc/oauth/login postAccessToken() IN {} OUT {}", params, tokenRes.getAccess_token());
	        return tokenRes;
	        	
		} catch (OAuth2Exception e) {
			
			String exCode = e.getOAuth2ErrorCode();
            String exMsg = e.getMessage();
			
			return this.getResponseVO(exCode, exMsg, locale);
		}

    }
	
	private ResponseVO getResponseVO(String exCode, String exMsg, Locale locale) {
        log.info("----- OAuth2Exception.exCode : " + exCode);
		log.info("----- OAuth2Exception.exMsg : " + exMsg);
		
		String exLocaleCode = "";
        String exLocaleMsg = "";
		
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
			exLocaleMsg = messageSource.getMessage(ResultCode.AE12.getCode(), null, locale) + exMsg;
			
			//TODO 접속 오류 횟수 ++
		} else {
			exLocaleCode = "AE01";
			exLocaleMsg = messageSource.getMessage(ResultCode.AE01.getCode(), null, locale);
		}
		
		return new ResponseVO(exLocaleCode, exLocaleMsg);
		
	}
	
	
	@RequestMapping(value = "/oauth/token", method = RequestMethod.POST)
	private ResponseVO getProcessOauthToken(									
									Principal principal,
						            @RequestParam Map<String, String> parameters) throws HttpRequestMethodNotSupportedException {
		log.info("=====  getProcessOauthToken() start... ");
		ResponseVO res = new ResponseVO(ResultCode.SUCCESS);
		
		tokenEndpoint.getAccessToken(principal, parameters);
		
		
		log.info("=====  getProcessOauthToken() end... ");
		return res;
	}
	
	
}
