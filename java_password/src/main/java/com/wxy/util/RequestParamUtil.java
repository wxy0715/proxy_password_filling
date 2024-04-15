package com.wxy.util;

import java.util.HashMap;
import java.util.Map;

public class RequestParamUtil {
    public RequestParamUtil() {
    }

    public static Map<String, String> urlSplit(String URL) {
        Map<String, String> mapRequest = new HashMap();
        String[] arrSplit = null;
        String strUrlParam = truncateUrlPage(URL);
        if (strUrlParam == null) {
            return mapRequest;
        } else {
            arrSplit = strUrlParam.split("[&]");
            String[] var4 = arrSplit;
            int var5 = arrSplit.length;

            for(int var6 = 0; var6 < var5; ++var6) {
                String strSplit = var4[var6];
                String[] arrSplitEqual = null;
                arrSplitEqual = strSplit.split("[=]");
                if (arrSplitEqual.length > 1) {
                    mapRequest.put(arrSplitEqual[0], arrSplitEqual[1]);
                } else if (arrSplitEqual[0] != "") {
                    mapRequest.put(arrSplitEqual[0], "");
                }
            }

            return mapRequest;
        }
    }

    private static String truncateUrlPage(String strURL) {
        String strAllParam = null;
        String[] arrSplit = null;
        strURL = strURL.trim();
        arrSplit = strURL.split("[?]");
        if (strURL.length() > 1 && arrSplit.length > 1) {
            for(int i = 1; i < arrSplit.length; ++i) {
                strAllParam = arrSplit[i];
            }
        }

        return strAllParam;
    }
}
