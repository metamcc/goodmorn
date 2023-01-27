package io.mcc.auth.mapper;

import org.apache.ibatis.annotations.Delete;
import org.apache.ibatis.annotations.Mapper;
import org.apache.ibatis.annotations.Param;
import org.apache.ibatis.annotations.Select;
import org.apache.ibatis.annotations.Update;

import io.mcc.auth.entity.CommonVO;

@Mapper
public interface UserMapper {

	@Select("SELECT user_seq, username, password, enabled, IFNULL(conect_error_cnt, 0) AS conect_error_cnt \r\n" + 
			"  FROM t_user \r\n" + 
			" WHERE username = #{username}")
	CommonVO findOneByUsername(@Param("username")String username);
	
	@Select("SELECT U.user_seq, \r\n" +
			"		U.password, \r\n" + 
			"		U.user_status_cd, \r\n" +
			"		IFNULL(U.conect_error_cnt, 0) AS conect_error_cnt, \r\n" + 
			"       CASE WHEN user_type_cd='03' THEN 'N' ELSE 'Y' END AS snsUserYn, \r\n" +			
			"		D.uuid, D.regist_dt \r\n" + 
			"  FROM t_user U \r\n" + 
			"  LEFT JOIN t_user_device D \r\n" + 
			"    ON D.user_seq = U.user_seq \r\n" + 
			" WHERE U.username = #{username} \r\n" + 
			" ORDER BY D.regist_dt desc \r\n" + 
			" LIMIT 1")
	CommonVO selectUserByUsername(@Param("username")String username);

	@Select("SELECT * FROM t_user WHERE user_seq = #{user_seq}")
	CommonVO findOne(@Param("user_seq")String user_seq);
	
	//@Select("SELECT user_seq,username,password,enabled FROM t_user WHERE sns_id = #{snsId}")
	@Select("SELECT U.user_seq, U.username, U.password, U.user_status_cd, U.enabled , D.uuid\r\n" + 
			"  FROM t_user U \r\n" + 
			"  LEFT JOIN t_user_device D \r\n" + 
			"    ON D.user_seq = U.user_seq \r\n" + 
			" WHERE U.sns_id = #{snsId} \r\n" + 
			" ORDER BY D.regist_dt DESC \r\n" + 
			" LIMIT 1")
	CommonVO selectUserBySnsId(@Param("snsId")String snsId);
	
	@Update("UPDATE t_user SET last_conect_dt = CURRENT_TIMESTAMP WHERE username = #{username}")
	void updateLastConectTime(String username);

	@Update("UPDATE t_user SET conect_error_cnt = conect_error_cnt+1, last_update_dt = CURRENT_TIMESTAMP WHERE username = #{username}")
	void increaseConectErrorCnt(@Param("username")String username);
	
	@Update("UPDATE t_user SET conect_error_cnt = 0 WHERE username = #{username}")
	void initUserConectErrorCntByUsername(@Param("username")String username);

	@Delete("DELETE FROM oauth_access_token WHERE user_name = #{username}")
	int deleteAccessToken(@Param("username")String username);

	@Delete("DELETE FROM oauth_refresh_token WHERE token_id = (SELECT refresh_token FROM oauth_access_token WHERE user_name = #{username})")
	int deleteRefreshToken(@Param("username")String username);

	


}
