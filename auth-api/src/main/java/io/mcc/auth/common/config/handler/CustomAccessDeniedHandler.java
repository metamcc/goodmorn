package io.mcc.auth.common.config.handler;

import java.io.IOException;
import java.util.Date;
import java.util.HashMap;
import java.util.Map;

import javax.servlet.ServletException;
import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;

import org.springframework.security.access.AccessDeniedException;
import org.springframework.security.oauth2.provider.error.OAuth2AccessDeniedHandler;
import org.springframework.stereotype.Component;


import lombok.extern.slf4j.Slf4j;

@Slf4j
@Component
public class CustomAccessDeniedHandler extends OAuth2AccessDeniedHandler {

	@Override
	public void handle(HttpServletRequest request, HttpServletResponse response,
			AccessDeniedException authException) throws IOException, ServletException {
		log.info("**********  CustomAccessDeniedHandler.handle()");		
		
		response.setContentType("application/json;charset=UTF-8");
        Map map = new HashMap();
        map.put("error", "400");
        map.put("message", authException.getMessage());
        map.put("path", request.getServletPath());
        map.put("timestamp", String.valueOf(new Date().getTime()));
        response.setContentType("application/json");
        response.setStatus(HttpServletResponse.SC_UNAUTHORIZED);
        response.getWriter().write("");
        
//        super.handle(request, response, authException);
//		
//	     // do something 
//	     response.sendRedirect("/error/deny"); 
		
	}

}
