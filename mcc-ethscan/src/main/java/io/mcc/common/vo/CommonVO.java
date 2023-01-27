package io.mcc.common.vo;

import lombok.ToString;

import java.util.HashMap;

@ToString
public class CommonVO<K, V> extends HashMap<K, V> {

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
}
