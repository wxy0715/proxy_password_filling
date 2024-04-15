package com.wxy.config;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.env.EnvironmentPostProcessor;
import org.springframework.core.Ordered;
import org.springframework.core.env.ConfigurableEnvironment;

/**
 * springfox3.0与springboot2.7.6版本适配
 */
public class WebEnvironmentPostProcessor implements EnvironmentPostProcessor, Ordered {

	@Override
	public void postProcessEnvironment(ConfigurableEnvironment environment, SpringApplication application) {
        System.setProperty("spring.mvc.pathmatch.matching-strategy", "ant_path_matcher");
		// 开启容器中的默认servlet
		System.setProperty("server.servlet.register-default-servlet", "true");
	}

	@Override
	public int getOrder() {
		return Ordered.HIGHEST_PRECEDENCE;
	}

}
