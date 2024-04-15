package com.wxy.aop;

import com.alibaba.fastjson.JSON;
import lombok.extern.slf4j.Slf4j;
import org.aspectj.lang.JoinPoint;
import org.aspectj.lang.annotation.AfterReturning;
import org.aspectj.lang.annotation.Aspect;
import org.aspectj.lang.annotation.Before;
import org.aspectj.lang.annotation.Pointcut;
import org.springframework.stereotype.Component;
import org.springframework.util.ObjectUtils;
import org.springframework.web.context.request.RequestContextHolder;
import org.springframework.web.context.request.ServletRequestAttributes;
import org.springframework.web.multipart.MultipartFile;

import javax.servlet.ServletRequest;
import javax.servlet.ServletResponse;
import java.lang.reflect.Field;

@Slf4j
@Aspect
@Component
public class LogAop {

    @Pointcut("execution(* com.wxy.controller.*.*(..))")
    public void point() {
    }
    @Before("point()")
    public void before(JoinPoint joinPoint) {
    }

    @AfterReturning(value = "point()", returning = "result")
    public void afterReturn(JoinPoint joinPoint, Object result) {
        Object[] args = joinPoint.getArgs();
        String info = "[ZJH]接口:".concat(((ServletRequestAttributes) RequestContextHolder.currentRequestAttributes()).getRequest().getRequestURI());
        if (!ObjectUtils.isEmpty(args)) {
            for (int i = 0; i < args.length; i++) {
                boolean flag = !ObjectUtils.isEmpty(args[i])
                        && !(args[0] instanceof ServletResponse)
                        && !(args[0] instanceof ServletRequest)
                        && !(args[0] instanceof MultipartFile[])
                        && !(args[0] instanceof MultipartFile)
                        && !hasByteArray(args[0]);
                if (!flag) {
                    args[i] = null;
                }
            }
            String jsonString = JSON.toJSONString(args);
            info = info.concat(" 参数:".concat(jsonString));
        }
        info = info.concat("\n[ZJH]响应数据:".concat(JSON.toJSONString(result)));
        if (info.length() > 1000) {
            try {
                log.info(info.substring(0, 1000));
            } catch (Exception e) {
                log.error("截断日志报错");
            }
        } else {
            log.info(info);
        }
    }

    private boolean hasByteArray(Object args) {
        try {
            boolean flag = false;
            for (Field field : args.getClass().getDeclaredFields()) {
                if (field.getType() == byte[].class) {
                    flag = true;
                }
            }
            return flag;
        } catch (Exception e) {
            return false;
        }
    }
}
