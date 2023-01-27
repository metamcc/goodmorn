package io.mcc.auth.common.exception;

import java.io.IOException;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.Collections;
import java.util.List;

import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;

import org.springframework.http.HttpEntity;
import org.springframework.http.HttpHeaders;
import org.springframework.http.HttpInputMessage;
import org.springframework.http.HttpOutputMessage;
import org.springframework.http.HttpStatus;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.http.converter.HttpMessageConverter;
import org.springframework.http.server.ServerHttpResponse;
import org.springframework.http.server.ServletServerHttpRequest;
import org.springframework.http.server.ServletServerHttpResponse;
import org.springframework.security.oauth2.common.exceptions.OAuth2Exception;
import org.springframework.security.oauth2.http.converter.jaxb.JaxbOAuth2ExceptionMessageConverter;
import org.springframework.security.oauth2.provider.error.DefaultOAuth2ExceptionRenderer;
import org.springframework.web.HttpMediaTypeNotAcceptableException;
import org.springframework.web.client.RestTemplate;
import org.springframework.web.context.request.NativeWebRequest;
import org.springframework.web.context.request.ServletWebRequest;

import lombok.extern.slf4j.Slf4j;

@Slf4j
public class CustomOAuthExceptionRenderer extends DefaultOAuth2ExceptionRenderer {

	
	
	@Override
	public void handleHttpEntityResponse(HttpEntity<?> responseEntityE, ServletWebRequest webRequest) throws Exception {
		if (responseEntityE == null) {
			return;
		}
		
		OAuth2Exception exception = (OAuth2Exception) responseEntityE.getBody();
		
		if(exception!=null){
			String errorCode = exception.getOAuth2ErrorCode();
			HttpStatus httpStatus = ((ResponseEntity<?>) responseEntityE).getStatusCode();
			CustomOauthException customException  = new CustomOauthException(exception.getMessage(),exception.getCause());
			customException.setOaut_error_code(errorCode);
			customException.setHttp_error_code(httpStatus.value());
			HttpEntity<OAuth2Exception> responseEntity = new HttpEntity<OAuth2Exception>(customException,responseEntityE.getHeaders());
			HttpInputMessage inputMessage = createHttpInputMessage(webRequest);
			HttpOutputMessage outputMessage = createHttpOutputMessage(webRequest,httpStatus.value());
			if (responseEntity instanceof ResponseEntity && outputMessage instanceof ServerHttpResponse) {
				((ServerHttpResponse) outputMessage).setStatusCode(httpStatus);
			}
			HttpHeaders entityHeaders = responseEntity.getHeaders();
			if (!entityHeaders.isEmpty()) {
				outputMessage.getHeaders().putAll(entityHeaders);
			}
			outputMessage.getHeaders().put("content-type", Arrays.asList("application/json;charset=UTF-8"));
			Object body = responseEntity.getBody();
			if (body != null) {
				writeWithMessageConverters(body, inputMessage, outputMessage);
			} else {
				outputMessage.getBody();
			}
		} 
		
	}
	@SuppressWarnings({ "unchecked", "rawtypes" })
	private void writeWithMessageConverters(Object returnValue, HttpInputMessage inputMessage,
			HttpOutputMessage outputMessage) throws IOException, HttpMediaTypeNotAcceptableException {
		List<MediaType> acceptedMediaTypes = inputMessage.getHeaders().getAccept();
		if (acceptedMediaTypes.isEmpty()) {
			acceptedMediaTypes = Collections.singletonList(MediaType.ALL);
		}
		MediaType.sortByQualityValue(acceptedMediaTypes);
		Class<?> returnValueType = returnValue.getClass();
		List<MediaType> allSupportedMediaTypes = new ArrayList<MediaType>();
		for (MediaType acceptedMediaType : acceptedMediaTypes) {
			for (HttpMessageConverter messageConverter : geDefaultMessageConverters()) {
				if (messageConverter.canWrite(returnValueType, acceptedMediaType)) {
					messageConverter.write(returnValue, acceptedMediaType, outputMessage);
					return;
				}
			}
		}
		for (HttpMessageConverter messageConverter : geDefaultMessageConverters()) {
			allSupportedMediaTypes.addAll(messageConverter.getSupportedMediaTypes());
		}
		throw new HttpMediaTypeNotAcceptableException(allSupportedMediaTypes);
	}

	private List<HttpMessageConverter<?>> geDefaultMessageConverters() {
		List<HttpMessageConverter<?>> result = new ArrayList<HttpMessageConverter<?>>();
		result.addAll(new RestTemplate().getMessageConverters());
		result.add(new JaxbOAuth2ExceptionMessageConverter());
		return result;
	}

	private HttpInputMessage createHttpInputMessage(NativeWebRequest webRequest) throws Exception {
		HttpServletRequest servletRequest = webRequest.getNativeRequest(HttpServletRequest.class);
		return new ServletServerHttpRequest(servletRequest);
	}

	private HttpOutputMessage createHttpOutputMessage(NativeWebRequest webRequest,int status) throws Exception {
		HttpServletResponse servletResponse = (HttpServletResponse) webRequest.getNativeResponse();
		servletResponse.setStatus(status);
		return new ServletServerHttpResponse(servletResponse);
	}
}
