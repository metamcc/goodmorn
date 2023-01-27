package io.mcc.auth.service;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import io.mcc.auth.mapper.CommonMapper;


@Service
public class CommonService {
	
	@Autowired
	private CommonMapper commonMapper;

	public String getSystemConfigValue(String name) {
		
		
		return commonMapper.getSystemConfigValue(name);
	}

	public int updateUserDevice(long userSeq, String deviceId) {
		int rst = commonMapper.updateUserDevice(userSeq, deviceId);
		
		if(rst == 0) {
			rst = commonMapper.insertUserDevice(userSeq, deviceId);
		}
		
		return rst;
	}
	
	public int updateUserDevicePlatform(long userSeq, String OsVer, String device_name, String deviceId, String deviceType) {
		int rst = 0;
		commonMapper.updateUserDevice(userSeq, deviceId);
		if(deviceId != null && (deviceId.trim()).length() > 1) {
			rst = commonMapper.updateUserDevicePlatformMgmt(deviceId, deviceType, OsVer, device_name);//유저 디바이스 정보 등록
		}
		return rst;
	}
}
