package io.mcc.auth.controller;

import java.io.IOException;
import java.util.Arrays;
import java.util.Base64;
import java.util.List;
import java.util.Locale;
import java.util.Map;

import javax.servlet.http.HttpServletRequest;

import java.util.Base64.Decoder;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.context.MessageSource;
import org.springframework.core.env.Environment;
import org.springframework.http.HttpEntity;
import org.springframework.http.HttpHeaders;
import org.springframework.http.HttpMethod;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.mobile.device.Device;
import org.springframework.mobile.device.DevicePlatform;
import org.springframework.security.oauth2.client.registration.ClientRegistration;
import org.springframework.security.oauth2.client.registration.ClientRegistrationRepository;
import org.springframework.security.oauth2.client.registration.ClientRegistration.ProviderDetails;
import org.springframework.util.StringUtils;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestHeader;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestMethod;
import org.springframework.web.bind.annotation.RestController;
import org.springframework.web.client.RestTemplate;

import io.mcc.auth.common.util.HttpParameterUtil;
import io.mcc.auth.entity.CommonVO;
import io.mcc.auth.entity.ResponseVO;
import io.mcc.auth.entity.ResultCode;
import io.mcc.auth.entity.User;
import io.mcc.auth.service.ClientService;
import io.mcc.auth.service.CommonService;
import io.mcc.auth.service.LoggingService;
import io.mcc.auth.service.UserService;
import com.google.api.client.googleapis.auth.oauth2.GoogleIdToken;
import com.google.api.client.googleapis.auth.oauth2.GoogleIdToken.Payload;
import com.google.api.client.googleapis.auth.oauth2.GoogleIdTokenVerifier;
import com.google.api.client.http.HttpTransport;
import com.google.api.client.http.javanet.NetHttpTransport;
import com.google.api.client.json.JsonFactory;
import com.google.api.client.json.jackson2.JacksonFactory;

import lombok.extern.slf4j.Slf4j;

@Slf4j
@RestController
@RequestMapping("/client")
public class AuthClientController {
	
	        	
    private String googleClientIdList= "google.client.id.list";
    private String accessTokenProofUrlKakao = "access.token.proof.url.kakao";
//    private static final String accessTokenProofUrlFacebook = "https://graph.accountkit.com/v1.3/me/?access_token=";
	
	@Autowired
    private Environment env;

	@Autowired
	private ClientRegistrationRepository clientRegistrationRepository;

	@Autowired
	private UserService userService;
	
	@Autowired
	private CommonService commonService;
	
	@Autowired
	private ClientService clientService;
	
	@Autowired
	private LoggingService loggingService;
	
	@Autowired
	private MessageSource messageSource;
	
	private static final HttpTransport transport = new NetHttpTransport();
	private static final  JsonFactory jsonFactory = new JacksonFactory();
	
	@RequestMapping(value = "/token/{registrationId}", method = RequestMethod.POST)
	private ResponseVO getClientLogin2(HttpServletRequest request,
									@PathVariable String registrationId, 
									@RequestHeader HttpHeaders reqHeader, 
									@RequestBody CommonVO paramVo, Locale locale, Device device) throws IOException {
		log.info("===== /token/registrationId getClientLogin2() IN {}, registrationId {}", paramVo, registrationId);
		ResponseVO res = new ResponseVO(ResultCode.SUCCESS);
		
		String idTokenStr = paramVo.getInfo("idtoken");
		String accessToken = paramVo.getInfo("token");
//		String fcm = paramVo.getInfo("FCM");
		
		log.info("----- idToken : " + idTokenStr);
		log.info("----- accessToken : " + accessToken);
		log.info("----- RemoteAddr : " + request.getRemoteAddr());
		
		//-ClientRegistration
		ClientRegistration clientRegistration = clientRegistrationRepository.findByRegistrationId(registrationId);
		ProviderDetails clientDetails = clientRegistration.getProviderDetails();
		String userInfoEndpointUri = clientDetails.getUserInfoEndpoint().getUri();
		
		//-RestTemplate
		RestTemplate restTemplate = new RestTemplate();
		HttpHeaders headers = new HttpHeaders();
		headers.add(HttpHeaders.AUTHORIZATION, "Bearer " + accessToken);
		
		String sns_id = "";
		
		if( registrationId.equals("google") ) {
			String[] clientIdList = env.getProperty(googleClientIdList).split("::");
			List clientIds = Arrays.asList(clientIdList);
			log.info("----- clientIds.size() : {} ", clientIds.size());
			
			GoogleIdTokenVerifier verifier = new GoogleIdTokenVerifier.Builder(transport, jsonFactory)
				    .setAudience(clientIds)
				    .build();
			
			GoogleIdToken idToken;
			try {
				idToken = verifier.verify(idTokenStr);
			} catch (Exception e) {
			      return new ResponseVO(ResultCode.IVCT.getCode(), messageSource.getMessage(ResultCode.IVCT.getCode(), null, locale));
		    }
			
			if (idToken != null) {
				Payload payload = idToken.getPayload();
				sns_id = payload.getSubject();
			} else {
				return new ResponseVO(ResultCode.AE61.getCode(), messageSource.getMessage(ResultCode.AE61.getCode(), null, locale));
			}		
		} else if( registrationId.equals("kakao") ) {
			String tokenProofUrl = env.getProperty(accessTokenProofUrlKakao);

			HttpEntity<String> entity = new HttpEntity<String>("", headers);
			
			ResponseEntity<Map> response = restTemplate.exchange(tokenProofUrl, HttpMethod.GET, entity, Map.class);
			Map userAttributes = response.getBody();
			
			if(response.getStatusCodeValue() == 200) {
				log.info("----- kakao userAttributes = " + userAttributes.toString());
				String app_id = String.valueOf( userAttributes.get("appId"));
				
				sns_id = String.valueOf( userAttributes.get("id"));
			} else {
				return new ResponseVO(ResultCode.AE62.getCode(), messageSource.getMessage(ResultCode.AE62.getCode(), null, locale));
			}
			
		} else if( registrationId.equals("facebook") ) {
			
			HttpEntity<String> entity = new HttpEntity<String>("", headers);
			
			ResponseEntity<Map> response = restTemplate.exchange(userInfoEndpointUri, HttpMethod.GET, entity, Map.class);
			Map userAttributes = response.getBody();
			if(response.getStatusCodeValue() == 200) {
				log.info("----- facebook response code = " + response.getStatusCodeValue());
				
				sns_id = (String)userAttributes.get("id");
			} else {
				return new ResponseVO(ResultCode.AE62.getCode(), messageSource.getMessage(ResultCode.AE62.getCode(), null, locale));
			}
		} 
		log.info("----- sns_id = " + sns_id);
		// 기존 회원 가입 여부 확인(DB 조회)
		CommonVO user = userService.selectUserBySnsId(sns_id);
		if (user == null) {		//x2ho.bugfix.
			// 비가입 회원(가입 유도)
			return new ResponseVO(ResultCode.NEU.getCode(), messageSource.getMessage(ResultCode.NEU.getCode(), null, locale));
		}

		//- 회원 상태 코드 확인
		String userStatusCd = user.getString("user_status_cd");
		if (user != null && "01".equals(userStatusCd)) { 
			//- 회원 상태 코드 확인
			if("02".equals(userStatusCd)) {
				//TODO 휴면 계정입니다.
				return new ResponseVO(ResultCode.NEU.getCode(), messageSource.getMessage(ResultCode.NEU.getCode(), null, locale));
			} else if("04".equals(userStatusCd) || "05".equals(userStatusCd)) {
				//TODO 탈퇴 처리 된 회원 입니다.
				return new ResponseVO(ResultCode.NEU.getCode(), messageSource.getMessage(ResultCode.NEU.getCode(), null, locale));
			}
			
			String username = user.getString("username");
			log.info("----- Username = " + user.getString("username"));
			log.info("----- Password = " + sns_id);
			
			User param = new User();
			param.setUsername((String) user.get("username"));
			param.setPassword(sns_id);
			
			//-Get Client Account from HttpHeaders
			List<String> authCreds = reqHeader.get(HttpHeaders.AUTHORIZATION);
			String[] client = null;
			if( authCreds!=null && authCreds.size()>0 ) {
				for(String authCred : authCreds) {
					if(authCred.contains("Basic")) {
						log.info("----- authCred : " + authCred);
						String cred = authCred.split("Basic ")[1];
						
						Decoder decoder = Base64.getDecoder(); 
						byte[] decodedBytes = decoder.decode(cred);
						client = new String(decodedBytes).split(":");
					}
				}
			}
			
			List<String> deviceIds = reqHeader.get("device-id");
			String reqDeviceId = "";
			if( deviceIds!=null && deviceIds.size()>0 ) {
				for(String deviceId : deviceIds) {
					reqDeviceId = deviceId;
				}
			} else {
				return new ResponseVO(ResultCode.M704.getCode(), messageSource.getMessage(ResultCode.M704.getCode(), null, locale));
			}
			
			if(client != null) {
				//deviceId 변경 여부 확인 : 변경시 기존 MCC 토큰 삭제
				String existDeviceId = user.getString("uuid");
				if(!StringUtils.isEmpty(reqDeviceId) && !reqDeviceId.equals(existDeviceId)) {
					int rst = loggingService.deleteAccessToken(username);
						rst += commonService.updateUserDevice(user.getLong("user_seq"), reqDeviceId);
				}
				
				// MCC 토큰 발행
				CommonVO resVo = clientService.getResourceOwnerPasswordToken(client, param);
				
				//디바이스 정보 등록 및 업데이트
				//신규 유저 및 기존 유저의 새로운 로그인시 디바이스 정보를 업데이트 합니다.
				//Device device = DeviceUtils.getCurrentDevice(request);
				DevicePlatform platform = device.getDevicePlatform();
				log.info("===== /token/registrationId DevicePlatform name {}, os {}", platform.name(), (platform == DevicePlatform.IOS) ); 
				log.info("===== /token/registrationId DevicePlatform isReqDeviceId {}, isexistDeviceId {}",!StringUtils.isEmpty(reqDeviceId), !reqDeviceId.equals(existDeviceId) ); 
				log.info("===== /token/registrationId DevicePlatform headerDeviceId {}, existDeviceId {}", reqDeviceId, existDeviceId); 
				commonService.updateUserDevicePlatform(user.getLong("user_seq"), HttpParameterUtil.getOsVer(request), HttpParameterUtil.getDiviceName(request), reqDeviceId, platform.name());//유저 일련번호 , 
				
				res.setResultObject(resVo.get("accessToken"));
			} else {
				res = new ResponseVO(ResultCode.AE64.getCode(), messageSource.getMessage(ResultCode.AE64.getCode(), null, locale));
			}
		} else { 
			//- 회원 상태 코드 확인
			if("02".equals(userStatusCd)) {
				//휴면 계정입니다.
				res = new ResponseVO(ResultCode.M110.getCode(), messageSource.getMessage(ResultCode.M110.getCode(), null, locale));
			} else if("04".equals(userStatusCd) || "05".equals(userStatusCd)) {
				//탈퇴 처리 된 회원 입니다.
				res = new ResponseVO(ResultCode.M110.getCode(), messageSource.getMessage(ResultCode.M110.getCode(), null, locale));
			}
		}
		
		log.info("===== /token/registrationId getClientLogin2() IN {}, registrationId {} OUT {}", paramVo, registrationId, res.getResultObject());
		return res;
	}

	@RequestMapping(value = "/info/{registrationId}", method = RequestMethod.POST)
	private ResponseVO getClientInfo(@PathVariable String registrationId, @RequestHeader HttpHeaders reqHeader, @RequestBody CommonVO paramVo, Locale locale) {
		log.info("===== /info/registrationId getClientInfo() IN {}, registrationId {} ", paramVo, registrationId);
		ResponseVO res = new ResponseVO(ResultCode.SUCCESS);
		
		String accessToken = paramVo.getInfo("token");
		log.info("----- accessToken = " + accessToken);
		
		ClientRegistration clientRegistration = clientRegistrationRepository.findByRegistrationId(registrationId);
		ProviderDetails clientDetails = clientRegistration.getProviderDetails();
		String userInfoEndpointUri = clientDetails.getUserInfoEndpoint().getUri();

		if (!StringUtils.isEmpty(userInfoEndpointUri)) {
			RestTemplate restTemplate = new RestTemplate();
			HttpHeaders headers = new HttpHeaders();
			headers.setContentType(MediaType.APPLICATION_FORM_URLENCODED);
			headers.add(HttpHeaders.AUTHORIZATION, "Bearer " + accessToken);

			HttpEntity<String> entity = new HttpEntity<String>("", headers);
			
			if(registrationId.equals("facebook")) {
				userInfoEndpointUri+= "?fields=email,name,gender,id,installed,last_name,first_name,birthday";
			}
			
			ResponseEntity<Map> response = restTemplate.exchange(userInfoEndpointUri, HttpMethod.GET, entity, Map.class);
			
			Map userAttributes = response.getBody();
			log.info("----- userAttributes = " + userAttributes.toString());

			if (userAttributes != null) {
				res.setResultObject(userAttributes);
			} else {
				res = new ResponseVO(ResultCode.NEU.getCode(), messageSource.getMessage(ResultCode.NEU.getCode(), null, locale));
			}
			
		} else {
			res = new ResponseVO(ResultCode.AE63.getCode(), messageSource.getMessage(ResultCode.AE63.getCode(), null, locale));
		}
		
		log.info("===== /info/registrationId getClientInfo() IN {}, registrationId {} OUT {}", paramVo, registrationId, res.getResultObject());
		return res;
	}
			

}

