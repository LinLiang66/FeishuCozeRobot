# FeishuCozeRobot

#### 介绍

基于技术栈Go+Gin+Redis 结合飞书机器人和字节扣子API实现的扣子智能体DIY飞书机器人，更多可能等你发现~

## 飞书机器人

https://open.feishu.cn/app

### 1.创建企业自建应用

![image](https://github.com/LinLiang66/FeishuCozeRobot/assets/93206426/9a77a642-80a3-4c2d-b597-57b3d298617d)

### 2.应用凭证

App ID、 App Secret

![image](https://github.com/LinLiang66/FeishuCozeRobot/assets/93206426/b986a6a5-66e2-440f-8ac3-306d6b6f0fcf)


### 3.添加飞书机器人能力
![image](https://github.com/LinLiang66/FeishuCozeRobot/assets/93206426/2957575a-a9da-4a91-b536-9d4a56cfddc0)

### 4.应用配置信息入库Redis缓存
``json
{
    "appid": "cli_a51e57ce4179900b",
    "app_secret": "gRUSCqIIeJKSdPwCG8uq0fW2c2myGy3g",
    "verification_token": "Mj3kyPKUCVC9ChmfqiB2U5ZMrArY6hJr",
    "encrypt_key": "Lin927919732Liang",
    "robot_appid": "您的扣子BotId",
    "robot_api_key": "您的扣子ApiPersonal_Access_Token"
}
``
### 4.事件订阅安全验证

Encrypt Key和 Verification Token 用于验证请求是否合法

![image](https://github.com/LinLiang66/FeishuCozeRobot/assets/93206426/32efa452-5a16-4616-ab61-6a7df0cb7a90)

配置消息事件接收地址

![image](https://github.com/LinLiang66/FeishuCozeRobot/assets/93206426/f117ca5b-3197-41d0-ba1e-206feb46a9bc)

配置卡片事件接收地址

![image](https://github.com/LinLiang66/FeishuCozeRobot/assets/93206426/5838c078-7911-41ca-9c8e-471a1b05cb05)

### 5.事件订阅

订阅：接收消息即可，其他事件随意

 1.消息事件订阅 im.message.receive_v1【接收消息v2.0】
 
![image](https://github.com/LinLiang66/FeishuCozeRobot/assets/93206426/3d2bce86-d8de-4041-93ec-59081b61c8b8)

  2.卡片事件订阅 card.action.trigger【卡片回传交互】、card.action.trigger_v1【消息卡片回传交互（旧）】
  
![image](https://github.com/LinLiang66/FeishuCozeRobot/assets/93206426/9924a1fd-4814-4368-80f3-d17b54ba589b)

### 4.权限管理

接收群聊中@机器人消息事件
读取用户发给机器人的单聊消息
获取用户发给机器人的单聊消息
获取与发送单聊、群组消息
以应用的身份发消息

## 扣子创建和搭建

### 

扣子 [https://www.coze.cn](https://www.coze.cn/docs/developer_guides/coze_api_overview)

飞书 [https://feishu.cn](https://open.feishu.cn)

