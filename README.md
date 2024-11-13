
![Logo](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/th5xamgrr6se0x5ro4g6.png)


# JUGGERNAUT

一个处理转发的多路由服务端


## Installation

Install my-project with npm

```bash
  npm install my-project
  cd my-project
```
    
## Deployment

To deploy this project run

```bash
  npm run deploy
```


## Tech Stack

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
  POST /juggernaut/udesk/newreply
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `null`      | `null` | `null` |

### 请求体示例

```json
{
"标题":"{{ticket.subject}}",
"客户":"{{ticket.user}}",
"回复时间":"{{ticket.replied_at}}",
"回复内容":"{{ticket.latest_customer_comment}}",
"工单地址":"{{ticket.web_url}}",
"工单id":"{{ticket.id}}",
"工单受理人":"{{ticket.assignee}}"
}
```

## FAQ

#### Question 1

Answer 1

#### Question 2

Answer 2


## 贡献者  
感谢 `王奥` `王培伦` `吴舒汀` 对于此项目的大力支持！

