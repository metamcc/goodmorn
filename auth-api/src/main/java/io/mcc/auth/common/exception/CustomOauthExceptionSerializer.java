package io.mcc.auth.common.exception;

import java.io.IOException;
import java.util.Arrays;
import java.util.Map;

import com.fasterxml.jackson.core.JsonGenerator;
import com.fasterxml.jackson.databind.SerializerProvider;
import com.fasterxml.jackson.databind.ser.std.StdSerializer;

import lombok.extern.slf4j.Slf4j;

@Slf4j
public class CustomOauthExceptionSerializer extends StdSerializer<CustomOauthException22>{
	
	public CustomOauthExceptionSerializer() {
        super(CustomOauthException22.class);
    }

	@Override
	public void serialize(CustomOauthException22 value, JsonGenerator  gen, SerializerProvider provider) throws IOException {
		 log.info(">>>>>>>>>>>>>>>>>  CustomOauthExceptionSerializer.serialize() called...");
		
		 gen.writeStartObject();
	     gen.writeNumberField("resultCode", value.getHttpErrorCode());
	     gen.writeStringField("resultMsg", value.getMessage());
	     gen.writeBooleanField("status", false);
	     gen.writeObjectField("data", null);
	     gen.writeObjectField("errors", Arrays.asList(value.getOAuth2ErrorCode(),value.getMessage()));
	        
	     if (value.getAdditionalInformation()!=null) {
	            for (Map.Entry<String, String> entry : value.getAdditionalInformation().entrySet()) {
	                String key = entry.getKey();
	                String add = entry.getValue();
	                
	                gen.writeStringField(key, add);
	            }
	     }
	        
	     gen.writeEndObject();
		
	}

}
