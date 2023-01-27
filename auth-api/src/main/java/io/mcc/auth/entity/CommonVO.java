package io.mcc.auth.entity;

import java.io.Serializable;
import java.util.LinkedHashMap;

import org.apache.ibatis.type.Alias;

import lombok.Data;
import lombok.Getter;
import lombok.Setter;

@Data
@Alias("commonVO")
public class CommonVO<K, V> extends LinkedHashMap<K, V> implements Serializable{
	
	public String getInfo(String key) {
		if ( get(key ) == null )
			return "";

		return (String)get(key);
	}
	
	public Object getObjectInfo(String key) {
		if ( get(key ) == null )
			return "";

		return get(key);
	}

	public String getString(final Object key) {
		final Object valueObj = this.get(key);
		if (valueObj == null) {
			return null;
		} else {
			return valueObj.toString();
		}
	}

	public int getInt(final Object key) {
		final Object valueObj = this.get(key);
		if (valueObj == null) {
			return 0;
		} else {
			final String value = valueObj.toString();
			if (value.isEmpty()) {
				return 0;
			} else {
				return Integer.parseInt(value);
			}
		}
	}

	public long getLong(final Object key) {
		final Object valueObj = this.get(key);
		if (valueObj == null) {
			return 0L;
		} else {
			final String value = valueObj.toString();
			if (value.isEmpty()) {
				return 0L;
			} else {
				return Long.parseLong(value);
			}
		}
	}
	
	
	
}
