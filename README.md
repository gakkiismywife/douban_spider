# 监测豆瓣小组最新讨论

### 项目介绍

恰逢双11，女朋友为了省钱整天在各大豆瓣小组里抄作业，为了方便她抄作业开发了这个小工具。
本工具只爬取指定豆瓣小组的首页，取第一页的帖子发布时间与当前时间进行对比，在10分钟之内的就微信通知。
方便女朋友工作之余还能摸鱼抄作业。

### 进度

- [x] ip代理
- [x] 帖子时间计算
- [x] 企业微信通知
- [x] 自定义帖子标题过滤

### 项目说明


### 配置说明

```yaml
proxy:
  id: "" #快代理的id
  secret: "" #快代理secret
  tunnel: "" #代理ip 我自己用的隧道代理 每次请求更换ip
wechat:
  key: "" #企业微信key
  secret: "" #企业微信应用secret
  agentId: "" #企业微信应用id
task:
  home: "https://www.douban.com/group/698716" #豆瓣小组主页 主要是为了标识
  urls:
    - "https://www.douban.com/group/698716/discussion?start=0&type=new" #讨论列表第一页
    - "https://www.douban.com/group/698716/discussion?start=25&type=new"#讨论列表第二页
    - "https://www.douban.com/group/698716/discussion?start=50&type=new"#讨论列表第三页
  interval: 300 #监测频率
  seconds: 1800 #发布时间在30分钟内的帖子 算新帖子
database: # 数据库配置 主要记录 消息发送 避免重复消息
  host: ""
  port:
  username: ""
  password: ""
  dbName: ""
  charset: "utf8mb4"
  parseTime: True
  maxOpen: 5
  maxIdle: 2
  maxIdleTime: 1
service:
  maxRetryTimes: 15 #重试次数
```

### 消息表

```sql
-- auto-generated definition
create table messages
(
    id         int auto_increment
        primary key,
    title      varchar(255) default '' not null,
    url        varchar(255) default '' not null,
    created_at timestamp null,
    updated_at timestamp null,
    constraint id
        unique (id)
);
```