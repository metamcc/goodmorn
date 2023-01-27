package io.mcc.auth;

import java.util.TimeZone;

import javax.annotation.PostConstruct;

import org.apache.ibatis.annotations.Mapper;
import org.mybatis.spring.annotation.MapperScan;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RestController;
import org.springframework.web.servlet.config.annotation.WebMvcConfigurationSupport;

import lombok.extern.slf4j.Slf4j;

@Slf4j
@RestController
@SpringBootApplication
@MapperScan(annotationClass=Mapper.class, basePackages="io.mcc.auth")
public class AuthApi2Application extends WebMvcConfigurationSupport{
//	public static final String CONFIG_LOCATION_PATH = "classpath:mybatis-config.xml";
	
	@PostConstruct
	void init() {
		TimeZone.setDefault(TimeZone.getTimeZone("UTC"));
	}
	
	public static void main(String[] args) {
		SpringApplication.run(AuthApi2Application.class, args);
	}
	
	@GetMapping("/health")
	public ResponseEntity<String> health() {
		log.info("AuthApi2Application ===>>> start");
		return new ResponseEntity<String>("Health OK", HttpStatus.OK);
	}
	
		
//	@Bean
//    public SqlSessionFactory sqlSessionFactory(DataSource dataSource) throws Exception{
//		PathMatchingResourcePatternResolver pathResolver = new PathMatchingResourcePatternResolver();
//		
//        SqlSessionFactoryBean sqlSessionFactoryBean = new SqlSessionFactoryBean();
//        sqlSessionFactoryBean.setDataSource(dataSource);
//        sqlSessionFactoryBean.setConfigLocation(pathResolver.getResource(CONFIG_LOCATION_PATH));
//                
//        return sqlSessionFactoryBean.getObject();
//    }
}
