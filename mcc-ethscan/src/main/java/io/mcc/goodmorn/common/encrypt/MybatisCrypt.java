package io.mcc.goodmorn.common.encrypt;

import com.github.trang.typehandlers.crypt.Crypt;
import com.github.trang.typehandlers.util.ConfigUtil;

public enum MybatisCrypt implements Crypt {
	INSTANCE;

	public String encrypt(String content) {
		return encrypt(content, ConfigUtil.getPrivateKey());
	}

	public String decrypt(String content) {
		return decrypt(content, ConfigUtil.getPrivateKey());
	}

	public String encrypt(String content, String password) {
		try {
			return AESCrypt.encrypt(password, content);
		} catch (Exception e) {
			throw new SecurityException(e);
		}
	}

	public String decrypt(String content, String password) {
		try {
			return AESCrypt.decrypt(password, content);
		} catch (Exception e) {
			throw new SecurityException(e);
		}
	}
}