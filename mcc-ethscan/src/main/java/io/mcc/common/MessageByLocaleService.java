package io.mcc.common;

import java.util.Locale;

public interface MessageByLocaleService {
	public String getMessage(String id);
	public String getMessage(String id, Locale locale);
}
