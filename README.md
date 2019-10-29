# 终端版QQ(based on redis)


# 功能
 - 大量用户同时在线聊天
 - 点对点聊天
 - 用户登录&注册

## 服务端功能

### 用户管理(数据存储: redis hash<users>)
 - 用户 id : 数字
 - 用户密码 : 字母数字组合
 - 用户昵称 : 显示
 - 用户性别 : 字符串
 - 用户头像 : url
 - 用户上线登录时间 : 字符串
 - 用户是否在线 ： online

### 用户行为
 - 发送信息
 - 接受信息

### 用户注册&登录

### 用户消息离线存储

## 客户端功能
 - 用户注册
 - 用户登录
 - 发送信息
 - 获取用户列表

## 通信协议
[0:4] 表示长度
[] json
