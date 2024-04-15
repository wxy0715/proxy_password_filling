package com.wxy.util;

import com.baomidou.mybatisplus.core.incrementer.IdentifierGenerator;

public class SnowflakeMpUtil {

    public static Long nextId(){
        IdentifierGenerator identifierGenerator = SpringUtils.getBean(IdentifierGenerator.class);
        return identifierGenerator.nextId(null).longValue();
    }
}
