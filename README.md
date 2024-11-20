
![Logo](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/th5xamgrr6se0x5ro4g6.png)


# JUGGERNAUT

一个处理转发的多路由服务端


## 构建

```bash
  git clone https://github.com/wa3721/juggernaut.git
  cd juggernaut
  CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o juggernaut 
  docker build -t juggernaut:{$tag} .
```

## 部署

使用docker部署

```bash
cp /juggernaut/config.yaml {$hostpath}
docker run --name juggernaut -v {$hostpath}:/JUGGERNAUT/config/ --restart=always juggernaut:{$tag}
```


## 技术栈

**Client:** UDESK触发器，手机短信转发器

**Server:** golang，gin，html，css，JavaScript


## API 文档

#### 服务说明主页

```http
  GET /juggernaut/readme
```

| Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `null`    | `null`   | 浏览器访问 |

#### 服务说明文档

```http
  GET /juggernaut/documentation
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `null`      | `null` | 直接下载本说明文档 |

#### 贡献者主页

```http
  GET /juggernaut/contributors
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `null`      | `null` | 请求贡献者主页 |

## 以下路由经由手机短信转发器触发
客户端详见：https://github.com/pppscn/SmsForwarder

#### 验证码转发路由
```http
  POST /juggernaut/captcha
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `null`      | `null` | `null` |

### 请求体示例

```string
from=1069119987553739999&content=1069119987553739999%0A%E3%80%90%E6%B7%B1%E5%9C%B3%E6%94%BF%E5%8A%A1%E7%9F%AD%E4%BF%A1%E5%B9%B3%E5%8F%B0%E3%80%91%E6%82%A8%E7%9A%84%E5%8F%A3%E4%BB%A4:%20073184%20%5B%E6%98%8E%E5%BE%A1%E8%BF%90%E7%BB%B4%E5%AE%A1%E8%AE%A1%E4%B8%8E%E9%A3%8E%E9%99%A9%E6%8E%A7%E5%88%B6%E7%B3%BB%E7%BB%9F%5D%0ASIM1_%0ASubId%EF%BC%9A1%0A2024-11-04%2011:17:51&timestamp=1730690271813
```

## 以下路由经由udesk触发器触发


#### 新工单提醒路由
```http
  POST /juggernaut/udesk/newticket
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `null`      | `null` | `null` |

### 请求体示例

```json
{
"标题":"{{ticket.subject}}",
"级别":"{{ticket.priority}}",
"环境":[
"{{ticket.SelectField_123041}}",
"{{ticket.SelectField_1500534}}",
"{{ticket.SelectField_1654224}}"
],
"提单人":"{{ticket.creator}}",
"提单时间":"{{ticket.created_at}}",
"工单地址":"{{ticket.web_url}}"
}
```


#### 新工单分配提醒路由
```http
  POST /juggernaut/udesk/remind
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `null`      | `null` | `null` |

### 请求体示例

```json
{
"工单id":"{{ticket.TextField_219254}}",
"客户":"{{ticket.user}}",
"创建时间":"{{ticket.created_at}}",
"主题":"{{ticket.subject}}",
"优先级":"{{ticket.priority}}",
"受理客服":"{{ticket.assignee}}",
"首次受理客服id":"{{ticket.first_assignee.id}}",
"工单链接":"{{ticket.web_url}}",
"客户手机号":"{{customer.cellphone}}"
}
```

#### 新客户回复提醒路由（循环提醒）
```http
  POST /juggernaut/udesk/reply
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `null`      | `null` | `null` |

### 请求体示例
### 接受两种格式json 分别处理回复提醒和取消消息
#### 回复提醒
```json
{
  "cloudId":"{{ticket.TextField_219254}}",
  "subject":"{{ticket.subject}}",
  "assignee":"{{ticket.assignee}}",
  "ticketUser":"{{ticket.user}}",
  "webUrl":"{{ticket.web_url}}",
  "latest_comment":"{{ticket.latest_customer_comment}}",
  "udeskId":"{{ticket.id}}"
}
```
#### 取消提醒
```json
{
  "cloudId":"{{ticket.TextField_219254}}",
  "assignee":"{{ticket.assignee}}",
  "silence":"{{ticket.SelectField_1666114}}"
}
```

## FAQ

#### reply:我重新点击取消发送，将他置成空，为什么不能重新发送

目前reply处的取消是一个一次性的动作，直接取消了对应的发送进程，在客户有新的回复之前，不能重新发送。

#### reply:发送进程有几种退出的方式？

发送进程有三种退出的方式
1.udesk页面点击取消发送
2.正常回复客户
3.工单交接，会先取消当前客服的发送进程，然后重新启动新进程到新客服


## 贡献者
感谢 `王奥` `王培伦` `吴舒汀` 等人对于此项目的大力支持！