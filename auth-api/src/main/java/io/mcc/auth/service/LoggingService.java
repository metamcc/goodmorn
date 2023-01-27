package io.mcc.auth.service;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Component;
import org.springframework.stereotype.Service;

import io.mcc.auth.entity.CommonVO;
import io.mcc.auth.mapper.LoggingMapper;
import io.mcc.auth.mapper.UserMapper;

import lombok.extern.slf4j.Slf4j;

@Slf4j
@Service
@Component("loggingService")
public class LoggingService {
	
	@Autowired
	private UserMapper userMapper;
	
	@Autowired
	private LoggingMapper loggingMapper;
	
	public int insertConecHistory(CommonVO param) {
		log.info("===== insertConecHistory() called... :{}", param.getString("username"));
		
		int rst = 0;
		
		CommonVO user = userMapper.findOneByUsername(param.getString("username"));
		
		if(user!=null) {
			log.info("----- user_seq : " + user.get("user_seq"));
			
			param.put("user_seq", user.get("user_seq"));
			
//			userMapper.updateLastConectTime(param.getString("username"));
			
			rst = loggingMapper.insertConecHistory(param);
			
		}else {
			log.info("----- user is null.");
		}
				
		return rst;
	}
	
	public CommonVO getConnentUserInfo(String username) {
		return userMapper.selectUserByUsername(username);
	}
	
	public long getUserConnentErrorCnt(String username) {
		CommonVO user = userMapper.findOneByUsername(username);
		
		long cnt = 0;
		
		if(user!=null) {
			cnt = (long)user.get("conect_error_cnt");
		} else {
			log.info("----- user is null.");
		}
		
		return cnt;
	}

	public void increaseConectErrorCnt(String username) {
		CommonVO user = userMapper.findOneByUsername(username);
		
		if(user!=null) {
			userMapper.increaseConectErrorCnt(username);
		} else {
			log.info("----- user is null.");
		}
	}

	public int deleteAccessToken(String username) {
		
		int rst = userMapper.deleteRefreshToken(username);
		
		rst += userMapper.deleteAccessToken(username);
		
		return rst;
	}

	public void initUserConectErrorCntByUsername(String username) {
		userMapper.initUserConectErrorCntByUsername(username);		
	}

}
