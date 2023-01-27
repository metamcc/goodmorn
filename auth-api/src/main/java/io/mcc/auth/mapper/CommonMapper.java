package io.mcc.auth.mapper;

import org.apache.ibatis.annotations.Insert;
import org.apache.ibatis.annotations.Mapper;
import org.apache.ibatis.annotations.Param;
import org.apache.ibatis.annotations.Select;
import org.apache.ibatis.annotations.Update;

@Mapper
public interface CommonMapper {

	@Select("SELECT value FROM t_system_config WHERE name = #{name}")
	String getSystemConfigValue(String name);

	@Update("UPDATE t_user_device SET regist_dt = CURRENT_TIMESTAMP WHERE user_seq=#{userSeq} AND UUID=#{deviceId}")
	int updateUserDevice(@Param("userSeq")long userSeq, @Param("deviceId") String deviceId);

	@Insert("INSERT IGNORE INTO t_user_device (user_seq, uuid) VALUES (#{userSeq}, #{deviceId})")
	int insertUserDevice(@Param("userSeq")long userSeq, @Param("deviceId") String deviceId);

	
    @Insert("insert into t_device_mgmt(uuid, fcm_token, os_ver, device_type, device_name) values(#{deviceId}, '', #{osVersion},  #{deviceType}, #{device_name})"
            + " on duplicate key update device_type = #{deviceType}, device_name = #{device_name}, os_ver = #{osVersion}")
	int updateUserDevicePlatformMgmt(@Param("deviceId") String deviceId, @Param("deviceType") String deviceType, @Param("osVersion") String osVersion, @Param("device_name") String device_name);
}
