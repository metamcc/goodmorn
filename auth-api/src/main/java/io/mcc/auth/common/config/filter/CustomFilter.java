package io.mcc.auth.common.config.filter;

import java.io.IOException;

import javax.servlet.Filter;
import javax.servlet.FilterChain;
import javax.servlet.FilterConfig;
import javax.servlet.ServletException;
import javax.servlet.ServletRequest;
import javax.servlet.ServletResponse;
import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;

import lombok.extern.slf4j.Slf4j;

@Slf4j
//@Component
public class CustomFilter implements Filter  {
//	private final LoggingService loggingService;
//	private final CommonService commonService;
//	
//	@Autowired
//	public CustomFilter(LoggingService loggingService, CommonService commonService) {
//		this.loggingService=loggingService;
//		this.commonService=commonService;
//	}

	@Override
	public void init(FilterConfig filterConfig) throws ServletException {
		// TODO Auto-generated method stub
		
	}

	@Override
	public void doFilter(ServletRequest request, ServletResponse response, FilterChain chain)
			throws IOException, ServletException {
		
		HttpServletRequest req = (HttpServletRequest) request;
		HttpServletResponse res = (HttpServletResponse) response;
		
        log.info("^^^^^ [preFilter] Starting a transaction for req : {}", req.getRequestURI());

        
        //-server process
        chain.doFilter(request, response);
        
        
        log.info("^^^^^ [afterFilter] Committing a transaction for req : {}", req.getRequestURI());                
        
        if(res.getStatus() != 200) {
        	log.info("^^^^^ [afterFilter] Response Status :{}",  res.getStatus());
	   		
        }
        
        String authenticate = res.getHeader("Connection");
   		log.info("^^^^^ [afterFilter] Connection :{}",  authenticate);
		
	}

	@Override
	public void destroy() {
		// TODO Auto-generated method stub
		
	}

	

}
