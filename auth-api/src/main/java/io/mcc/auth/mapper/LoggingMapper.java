package io.mcc.auth.mapper;

import org.apache.ibatis.annotations.Insert;
import org.apache.ibatis.annotations.Mapper;

import io.mcc.auth.entity.CommonVO;

@Mapper
public interface LoggingMapper {

	@Insert("INSERT INTO t_user_conect_hist "
			+ "(user_seq, hist_type_cd, conect_dt, conect_ip_addr, conect_trmnl, conect_os, conect_brwsr) "
		  + "VALUES (#{user_seq},#{hist_type_cd},CURRENT_TIMESTAMP,#{conect_ip_addr},#{conect_trmnl},#{conect_os},#{conect_brwsr})")
	int insertConecHistory(CommonVO param);
}
