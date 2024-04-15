package com.wxy.util;

import org.springframework.lang.Nullable;
import org.springframework.util.ObjectUtils;

public class ObjectUtil extends ObjectUtils {
    public ObjectUtil() {
    }

    public static boolean isNotArray(@Nullable Object obj) {
        return !isEmpty(obj);
    }

    public static boolean isNotEmpty(@Nullable Object[] array) {
        return !isEmpty(array);
    }

    public static boolean isNotEmpty(@Nullable Object obj) {
        return !isEmpty(obj);
    }
}
