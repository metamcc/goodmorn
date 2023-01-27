package io.mcc.auth.common.config;

import java.util.List;

import io.mcc.auth.common.util.GenerateLogTransactionKey;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.mobile.device.DeviceHandlerMethodArgumentResolver;
import org.springframework.mobile.device.DeviceResolverHandlerInterceptor;
import org.springframework.web.method.support.HandlerMethodArgumentResolver;
import org.springframework.web.servlet.config.annotation.InterceptorRegistry;
import org.springframework.web.servlet.config.annotation.ResourceHandlerRegistry;
import org.springframework.web.servlet.config.annotation.WebMvcConfigurationSupport;


@Configuration
public class WebMvcConfigurer extends WebMvcConfigurationSupport{


	@Bean
	public DeviceResolverHandlerInterceptor 
	        deviceResolverHandlerInterceptor() {
	    return new DeviceResolverHandlerInterceptor();
	}

	@Bean
	public DeviceHandlerMethodArgumentResolver 
	        deviceHandlerMethodArgumentResolver() {
	    return new DeviceHandlerMethodArgumentResolver();
	}

	@Override
	public void addArgumentResolvers( List<HandlerMethodArgumentResolver> argumentResolvers) {
	   argumentResolvers.add(deviceHandlerMethodArgumentResolver());
	}


	@Override
	public void addInterceptors(InterceptorRegistry registry) {
//		registry.addInterceptor(interceptor);
		registry.addInterceptor(deviceResolverHandlerInterceptor());
		registry.addInterceptor(new GenerateLogTransactionKey()).addPathPatterns("/**");
	}
	
	protected void addResourceHandlers(ResourceHandlerRegistry registry) {
    	
    	registry.addResourceHandler("swagger-ui.html")
    			.addResourceLocations("classpath:/META-INF/resources/");
    	
    	registry.addResourceHandler("/webjars/**")
    			.addResourceLocations("classpath:/META-INF/resources/webjars/");
	}
}
