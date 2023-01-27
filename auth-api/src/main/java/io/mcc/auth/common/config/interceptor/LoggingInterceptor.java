package io.mcc.auth.common.config.interceptor;

import java.io.PrintWriter;
import java.util.Enumeration;
import java.util.Locale;

import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;

import io.mcc.auth.common.util.HttpParameterUtil;
import io.mcc.auth.entity.CommonVO;
import io.mcc.auth.entity.ResultCode;
import io.mcc.auth.service.CommonService;
import io.mcc.auth.service.LoggingService;
import org.apache.commons.lang3.StringUtils;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.context.MessageSource;
import org.springframework.http.HttpHeaders;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.web.servlet.HandlerInterceptor;
import org.springframework.web.servlet.ModelAndView;
import org.springframework.security.core.Authentication;
import org.springframework.security.web.authentication.WebAuthenticationDetails;

import lombok.extern.slf4j.Slf4j;

@Slf4j
public class LoggingInterceptor implements HandlerInterceptor{

	private final LoggingService loggingService;
	private final CommonService commonService;
	
	@Autowired
	private MessageSource messageSource;


	@Autowired
	public LoggingInterceptor(LoggingService loggingService, CommonService commonService) {
		this.loggingService=loggingService;
		this.commonService=commonService;
	}

	@Override
	public boolean preHandle(HttpServletRequest request, HttpServletResponse response, Object handler)
		throws Exception {

		
		//- 시스템 점검중 접근 차단
		Locale locale = Locale.ENGLISH;
		/*
		String prtConfig = messageSource.getMessage("MCC_SYSTEM_CHECK_YN", null, locale);	
		log.info("----- [preHandle] MCC_SYSTEM_CHECK_YN : {}", prtConfig);
		if("Y".equals(prtConfig)) {
			String msg = "The system is currently under maintenance";//messageSource.getMessage("M701", null, locale);
			
			try {
				response.setContentType("text/html; charset=UTF-8");
				PrintWriter out = response.getWriter();
				out.println("{\r\n" + 
						"    \"resultCode\": \"M701\",\r\n" + 
						"    \"resultMessage\": \""+msg+"\"\r\n" + 
						"}");
	            out.flush();
	            out.close();
			} catch (Exception e) {
				log.error(e.getMessage(), e);
			}
			return false;
		}*/
		
		log.info("+++++ preHandle Method Execution +++++");
		String reqUri = request.getRequestURI();
		Authentication auth = SecurityContextHolder.getContext().getAuthentication();
		WebAuthenticationDetails details = (WebAuthenticationDetails) auth.getDetails(); //사용자 정보를 가져온다
		String ssionid = details.toString();
		System.setProperty(ssionid, request.getLocale().toString());
		
		log.info("----- [preHandle] Req URI : {}", reqUri);
		
		if("/oauth/token".equals(reqUri)) {
			//- refresh_token 관련 정보 확인
			String grantType = request.getParameter("grant_type");
			if("refresh_token".equals(grantType)) {
				log.info("----- [preHandle] grant_type : {} / refresh_token : {}", request.getParameter("grant_type"), request.getParameter("refresh_token"));
			} 
		}
		
		String lang = request.getHeader(HttpHeaders.ACCEPT_LANGUAGE);
		if (!StringUtils.isBlank(lang)) {
			locale = new Locale(lang);
		}
		
		Enumeration<?> ele = request.getHeaderNames();
		 while (ele.hasMoreElements()) {
			 String curr = (String) ele.nextElement();
		 }
		 
		String reqDeviceId = request.getHeader("device-id");
		log.info("----- [header] reqDeviceId: {}", reqDeviceId);
		
		String[] params = HttpParameterUtil.getParameters(request).split("&");
		long err_cnt = 0;
		String exist_uuid = "";
		for(String param : params) {
			 if(param.contains("username")){
				
				String username = param.split("=")[1];
				
				log.info("----- username : {}", username);
				
				if(!StringUtils.isEmpty(username)){
					
					CommonVO userInfo = loggingService.getConnentUserInfo(username);
					
					if(userInfo!=null) {
						err_cnt = userInfo.getLong("conect_error_cnt");
						exist_uuid = userInfo.getString("uuid");
						log.info("----- [DB] con_err_cnt : {}", err_cnt);
						log.info("----- [DB] reqDeviceId: {} / uuid: {}", reqDeviceId, exist_uuid);
						
						if(!StringUtils.isEmpty(reqDeviceId) && !reqDeviceId.equals(exist_uuid)) {
							log.info("----- [header2] device-id : {}", reqDeviceId);
							
							
							int rst = loggingService.deleteAccessToken(username);
							
								rst += commonService.updateUserDevice(userInfo.getLong("user_seq"), reqDeviceId);
							
						}
					} else {
						log.info("----- userInfo is null.");
					}
				}
			}
		}
		
		String value = "5";
		log.info("-----   CONNECT_ERROR_CNT = {}", err_cnt);
		if(err_cnt >= 5) {
			log.info("-----   MAX_CONNECT_ERROR_OVER !!! GO TO HELL");
				
			String message = messageSource.getMessage(ResultCode.M023.getCode(), null, locale);
			
			try {
				response.setContentType("text/html; charset=UTF-8");
				PrintWriter out = response.getWriter();
				out.println("{\r\n" + 
						"    \"resultCode\": \"M023\",\r\n" + 
						"    \"resultMessage\": \""+message+"\"\r\n" + 
						"}");
	            out.flush();
	            out.close();
			} catch (Exception e) {
				log.error(e.getMessage(), e);
			}
			return false;
		}
		
		return true;
	}
	
	@Override
	public void postHandle(	HttpServletRequest request, HttpServletResponse response,
			Object handler, ModelAndView modelAndView) throws Exception {
		log.info("+++++ postHandle Method Execution +++++");
		
		log.info("[postHandle][" + request + "]");
	}
	
	@Override
	public void afterCompletion(HttpServletRequest request, HttpServletResponse response,
			Object handler, Exception ex) throws Exception {		
		log.info("+++++   afterCompletion method Execution +++++ [exception: " + ex + "]");
		if (ex != null){
			log.info("+++++   [ex.getMessage : " + ex.getMessage() + "]");
			log.info("+++++   [ex.SimpleName : " + ex.getClass().getSimpleName() + "]");
			String exSimpleName = ex.getClass().getSimpleName();
	    } 
    	
    	CommonVO logParam = new CommonVO();
    	
    	String username = "";
        
    	//필요없는 정보는 삭제
    	Authentication auth = SecurityContextHolder.getContext().getAuthentication();
		WebAuthenticationDetails details = (WebAuthenticationDetails) auth.getDetails();
		String ssionid = details.toString();
    	System.clearProperty(ssionid);
    	
		String[] params = HttpParameterUtil.getParameters(request).split("&");
		for(String param : params) {
			if(param.contains("_psip")) {
				logParam.put("conect_ip_addr", param.split("=")[1]);
			} else if(param.contains("username")){
				username = param.split("=")[1];
				logParam.put("username", param.split("=")[1]);
			}
		}
		
		logParam.put("conect_trmnl", HttpParameterUtil.getUserAgentInfo(request, "device"));
		logParam.put("conect_os", HttpParameterUtil.getUserAgentInfo(request, "os"));
		logParam.put("conect_brwsr", HttpParameterUtil.getUserAgentInfo(request, "brwsr"));
		
		if(!StringUtils.isEmpty(username)) {
			log.info("----- response.getStatus() : {}", response.getStatus());
			if(response.getStatus()==200) {
				//
			} else {
				
				logParam.put("hist_type_cd", "02");	//로그인 오류
	
				//접속 오류 횟수 ++
				loggingService.increaseConectErrorCnt(username);
				loggingService.insertConecHistory(logParam);
				
			}
		}
    }
	
}


