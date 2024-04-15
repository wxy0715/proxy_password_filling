package com.wxy.util;


import lombok.Data;
import lombok.ToString;

import java.util.Date;
import java.util.Map;

@Data
@ToString
public class WsMessage {
    /** 消息id */
    private String id;
    /** 消息发送类型 */
    private Integer code;
    /** 发送人用户id */
    private String sendUserId;
    /** 发送人用户名 */
    private String username;
    /** 接收人用户id，多个逗号分隔 */
    private String receiverUserId;
    /** 发送时间 */
    private Date sendTime;
    /** 消息类型 */
    private Integer type;
    /** 消息内容 */
    private String msg;
    /** 消息扩展内容 */
    private Map<String, Object> ext;
}
