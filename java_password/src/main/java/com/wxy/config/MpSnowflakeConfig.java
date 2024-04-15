package com.wxy.config;

import com.baomidou.mybatisplus.core.incrementer.DefaultIdentifierGenerator;
import com.baomidou.mybatisplus.core.incrementer.IdentifierGenerator;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

@Configuration
public class MpSnowflakeConfig {
    @Bean
    public IdentifierGenerator identifierGenerator(){
        return new DefaultIdentifierGenerator();
    }
}
