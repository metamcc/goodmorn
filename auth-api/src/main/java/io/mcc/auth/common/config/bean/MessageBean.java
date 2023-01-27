package io.mcc.auth.common.config.bean;

import java.util.Locale;

import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.MessageSource;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.context.support.ReloadableResourceBundleMessageSource;
import org.springframework.web.servlet.LocaleResolver;
import org.springframework.web.servlet.i18n.AcceptHeaderLocaleResolver;

@Configuration
public class MessageBean {

	@Bean
	public LocaleResolver localeResolver() {
		AcceptHeaderLocaleResolver alr = new AcceptHeaderLocaleResolver();
		alr.setDefaultLocale(Locale.US);
		return alr;
	}
	
	@Value("${app.messages.message-path}")
	private String messageMessagePath;
	
	@Value("${app.messages.error-path}")
	private String errorMessagePath;
	
	@Value("${app.messages.cache-second}")
	private int messageCacheSecond;
	
	@Bean
	public MessageSource messageSource() {
		ReloadableResourceBundleMessageSource messageSource = new ReloadableResourceBundleMessageSource();
		messageSource.setBasenames(messageMessagePath, errorMessagePath);
		messageSource.setDefaultEncoding("UTF-8");
		messageSource.setCacheSeconds(messageCacheSecond);
		return messageSource;
	}
}
