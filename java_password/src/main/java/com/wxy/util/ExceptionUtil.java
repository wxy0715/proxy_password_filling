package com.wxy.util;

public class ExceptionUtil {
    public static void isNull(Object value,String msg) {
        if (ObjectUtil.isEmpty(value)) {
            throw new RuntimeException(msg);
        }
    }

    public static void isTrue(boolean flag, String msg) {
        if (flag) {
            throw new RuntimeException(msg);
        }
    }

    public static void error(String msg) {
        throw new RuntimeException(msg);
    }
}
