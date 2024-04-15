package com.wxy.config;

import com.github.xiaoymin.knife4j.spring.annotations.EnableKnife4j;
import io.swagger.annotations.Api;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.context.annotation.Import;
import springfox.bean.validators.configuration.BeanValidatorPluginsConfiguration;
import springfox.documentation.builders.ApiInfoBuilder;
import springfox.documentation.builders.PathSelectors;
import springfox.documentation.builders.RequestHandlerSelectors;
import springfox.documentation.service.ApiInfo;
import springfox.documentation.spi.DocumentationType;
import springfox.documentation.spring.web.plugins.Docket;

import java.sql.Timestamp;
import java.util.Date;

/**
 * Swagger2配置信息
 * 这里分了两组显示
 * 第一组是api，当作用户端接口
 * 第二组是admin，当作后台管理接口
 * 也可以根据实际情况来减少或者增加组
 */

@Configuration
@EnableKnife4j
// 对JSR303提供支持
@Import(BeanValidatorPluginsConfiguration.class)
public class Swagger2Config {
    @Value("${spring.application.name}")
    private String applicationName;
    private ApiInfo apiInfo() {
        return new ApiInfoBuilder()
                .title("SpringBoot整合Knife4j-API文档")
                .description("本文档描述了SpringBoot如何整合Knife4j")
                .version("1.0")
                .build();
    }

    @Bean
    public Docket api() {
        Docket docket = new Docket(DocumentationType.SWAGGER_2)
                .enable(true)
                .apiInfo(apiInfo())
                //分组名称
                .groupName(applicationName)
                .select()
                //这里指定Controller扫描包路径
                .apis(RequestHandlerSelectors.withClassAnnotation(Api.class))
                .paths(PathSelectors.any())
                .build()
                .directModelSubstitute(Timestamp.class, String.class)
                .directModelSubstitute(Date.class, String.class)
                .useDefaultResponseMessages(false);
        return docket;
    }
}


