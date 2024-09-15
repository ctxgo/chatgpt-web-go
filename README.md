<div align="center">
<img src="./docs/openai.svg" style="width:64px;height:64px;margin:0 32px" alt="icon"/>

<h1 align="center">ChatGPT Web Go</h1>

A commercially-viable ChatGpt web application built with Go.

可部署商业化的 ChatGpt 网页应用。


💡 本项目是后端服务，前端对应的项目是：[ctxgo/ChatGpt-Web](https://github.com/ctxgo/ChatGpt-Web)

</div>
<br>

> [!IMPORTANT]
> **此项目 Fork 自 [neoean/chatgpt-web-go](https://github.com/neoean/chatgpt-web-go)**
>
> 添加 gemini 支持


## 主要功能
#### 包括但不限于：
- 后台管理系统,可对用户,Token,商品,卡密等进行管理
- 精心设计的 UI，响应式设计
- 极快的首屏加载速度（~100kb）
- 支持Midjourney绘画和DALL·E模型绘画,GPT4等应用
- 海量的内置 prompt 列表，来自[中文](https://github.com/PlexPt/awesome-chatgpt-prompts-zh)和[英文](https://github.com/f/awesome-chatgpt-prompts)
- 一键导出聊天记录，完整的 Markdown 支持
- 支持自定义API地址（如：[openAI](https://api.openai.com) / [API2D](https://api2d.com/r/192767)）


#### TODO：
- [x] API Key 功能实现
- [x] API Proxy 代理
- [ ] 绘画功能
- [ ] 思维导图功能
- [ ] 支付功能完善
- [ ] server端渲染模式支持
- [ ] Docker 支持 


## 本地启动
**0.环境要求准备**
- golang1.18
- mysql 5.7+
- redis
- goland

**1.先 `Fork` 本项目，然后克隆到本地。**
```
建议目录 ~/go/src/github.com/heimeropen/
git clone https://github.com/ctxgo/chatgpt-web-go.git
```

**2.导入sql**
```
# sql文件
./model/sql/chatgpt.sql
```

**3.配置文件**
在 ./config 目录下新建文件 dev.yml 内容如下：
（配置内容需要更具自己环境更改）
```
port: 8899

db:
  type: mysql
  host: 127.0.0.1:3306
  hostR1: 127.0.0.1:3306
  user: root
  password: yourpassword
  name: chatgpt_web_new_go

redis:
  addr: 127.0.0.1:6379
  password: yourpassword

aiClient:
  proxy: xxxxxxxxxxxxx

emailServer:
  host: "xxxxxxxxxx"
  port: "587"
  SenderName: xxxxxx
  User: xxxxxxxxx
  Password: xxxxxx
```

**4.运行**
```
用goland打开项目
启动main函数：
./cmd/server/main.go
```

**前端服务**

前端服务安装参考：
[ctxgo/ChatGpt-Web](https://github.com/ctxgo/ChatGpt-Web)
```
前端项目需要修改配置文件 .env.development, 指向本地服务端：
VITE_APP_REQUEST_HOST=http://127.0.0.1:8899
```

#### 登录信息
```
账号: admin@c.om
密码: admin123
```

### 页面截图

![cover](https://files.catbox.moe/tp963e.png)
![cover](https://files.catbox.moe/y5avbx.png)
![cover](https://files.catbox.moe/k16jsz.png)
![cover](https://files.catbox.moe/8o5oja.png)


## 📋 开源协议

[![License MIT](https://img.shields.io/badge/License-MIT-brightgreen.svg)](https://github.com/79E/ChatGpt-Web/blob/master/license)
