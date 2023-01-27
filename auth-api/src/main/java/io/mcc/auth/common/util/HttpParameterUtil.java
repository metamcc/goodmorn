package io.mcc.auth.common.util;

import java.io.IOException;
import java.util.Enumeration;
import java.util.HashMap;

import javax.servlet.http.HttpServletRequest;

import org.apache.commons.lang3.StringUtils;

import uap_clj.java.api.*;

import lombok.extern.slf4j.Slf4j;

@Slf4j
public class HttpParameterUtil {

	public static String getParameters(HttpServletRequest request) {
		StringBuffer posted = new StringBuffer();
	    Enumeration<?> e = request.getParameterNames();
	    if (e != null) {
	        posted.append("?");
	    }
	    while (e.hasMoreElements()) {
	        if (posted.length() > 1) {
	            posted.append("&");
	        }
	        String curr = (String) e.nextElement();
	        posted.append(curr + "=");
	        if (curr.contains("password") 
	          || curr.contains("pass")
	          || curr.contains("pwd")) {
	            posted.append("*****");
	        } else {
	            posted.append(request.getParameter(curr));
	        }
	    }
	    String ip = request.getHeader("X-FORWARDED-FOR");
	    String ipAddr = (ip == null) ? getRemoteAddr(request) : ip;
	    if (ipAddr!=null && !ipAddr.equals("")) {
	        posted.append("&_psip=" + ipAddr); 
	    }
	    return posted.toString();
	}
	
	
	public static String getUserAgentInfo(HttpServletRequest request, String key) throws IOException {
		
		String agent = request.getHeader("User-Agent");
		//-ex) "Good Morn/1.0.2 (iPhone; iOS 12.1.2; Scale/3.00)"
		
		String value = "";
		log.info("----- [User Agent] getUserAgentInfo agent: {}", agent);
		if(!StringUtils.isEmpty(agent)) {
		
			HashMap b = Browser.lookup(agent);
		    HashMap o = OS.lookup(agent);
		    HashMap d = Device.lookup(agent);

		    String ua_device = d.get("family").toString();
		    String ua_os = o.get("family").toString();
		    String ua_browser = b.get("family").toString();
			
			if("os".equals(key)) {
				value = ua_os;
			} else if("device".equals(key)) {
				value = ua_device;
			} else if("brwsr".equals(key)) {
				value = ua_browser;
			}
			log.info("----- [User Agent] value: {}", value);
		}
		
		return value;
	}
	
	public static String getUserOS(HttpServletRequest request) {
		String agent = request.getHeader("User-Agent");
		
		String os = "";
		
		if(!StringUtils.isEmpty(agent)) {
			HashMap o = OS.lookup(agent);
			os = o.get("family").toString();
		}
		
		return os;
	}
	
	public static String getUserBrwsr(HttpServletRequest request) {

		String agent = request.getHeader("User-Agent");
		
		String brower = "";
		
		if(!StringUtils.isEmpty(agent)) {
			HashMap b = Browser.lookup(agent);
			brower = b.get("family").toString();
		}
		
		return brower;

	}
	
	public static String getUserDevice(HttpServletRequest request) {

		String agent = request.getHeader("User-Agent");
		
		String device = "";
		
		if(!StringUtils.isEmpty(agent)) {
			HashMap d = Device.lookup(agent);
			device = d.get("family").toString();
		}
		
		return device;

	}
	
	
	private static String getRemoteAddr(HttpServletRequest request) {
	    String ipFromHeader = request.getHeader("X-FORWARDED-FOR");
	    if (ipFromHeader != null && ipFromHeader.length() > 0) {
	        log.debug("ip from proxy - X-FORWARDED-FOR : " + ipFromHeader);
	        return ipFromHeader;
	    }
	    
	    return request.getRemoteAddr();
	}
	
	public static String getOsVer(HttpServletRequest request){
		String userAgent = request.getHeader("User-Agent");
		String osVer = "";

		if (userAgent!=null)
			userAgent = userAgent.toLowerCase();

		int idxF = -1;
		if((idxF=userAgent.toLowerCase().indexOf("android")) >= 0){
			int idxSep = userAgent.lastIndexOf(";");
			int idxstr = userAgent.lastIndexOf("android")+7;
			osVer = userAgent.substring(idxstr, idxSep).trim();
		} else if((idxF=userAgent.toLowerCase().indexOf("ios")) >= 0){
			int idxSep = userAgent.indexOf(";", idxF+4);
			osVer = userAgent.substring(idxF+4, idxSep).trim();
		} else {
			osVer = "UnKnown";
		}
		return osVer;
	}
	

	public static String getDiviceName(HttpServletRequest request){
		String userAgent = request.getHeader("User-Agent");
		String diviceName = "";

		if (userAgent!=null)
			userAgent = userAgent.toLowerCase();

		if((userAgent.indexOf("ios")) >= 0){
			diviceName = "ios";
		} else if(userAgent.indexOf("android") >= 0){
			//GoodMorn/1.2.0 (SM-G930S; Android 26;)"; ==> SM-G930S
			int idx    = userAgent.indexOf("(");		
			int idxEnd = userAgent.indexOf(";");
			if(idx > -1 &&  idxEnd > 0) {
				diviceName = userAgent.substring(idx+1, idxEnd);	
			}
		} else {
			diviceName = "UnKnown";
		}
		return diviceName;
	}
}
