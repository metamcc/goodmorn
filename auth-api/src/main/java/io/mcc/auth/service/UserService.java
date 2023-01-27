package io.mcc.auth.service;

import java.util.ArrayList;
import java.util.HashSet;
import java.util.List;
import java.util.Set;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.security.core.GrantedAuthority;
import org.springframework.security.core.userdetails.UserDetails;
import org.springframework.security.core.userdetails.UserDetailsService;
import org.springframework.security.core.userdetails.UsernameNotFoundException;
import org.springframework.stereotype.Component;

import io.mcc.auth.entity.CommonVO;
import io.mcc.auth.mapper.UserMapper;

import lombok.extern.slf4j.Slf4j;

@Slf4j
@Component("userDetailsService")
public class UserService implements UserDetailsService{
	
	@Autowired
	private UserMapper userMapper;

	@Override
	public UserDetails loadUserByUsername(String username) throws UsernameNotFoundException {
		log.info("Authenticating {}", username);
		String lowercaseLogin = username.toLowerCase();
		
		CommonVO user =  userMapper.findOneByUsername(lowercaseLogin);
		if (user == null) {
            //throw new UsernameNotFoundException("User " + lowercaseLogin + " was not found in the database");
			return null;
        } 
		
		Set<GrantedAuthority> dbAuthSet = new HashSet<GrantedAuthority>();
		List<GrantedAuthority> dbAuths = new ArrayList<GrantedAuthority>(dbAuthSet);
		
		return new org.springframework.security.core.userdetails.User(user.getInfo("username"), user.getInfo("password"), dbAuths);
	}
	
	public CommonVO selectUserBySnsId(String snsId)throws UsernameNotFoundException {		
		log.info("Authenticating by sns_id : {}", snsId);

		CommonVO user = userMapper.selectUserBySnsId(snsId);
		if (user == null) {
            //throw new UsernameNotFoundException("User " + lowercaseLogin + " was not found in the database");
			return null;
        } 
		
		return user ;
	}

}
