# Go技术文档大纲-礼品码

## 1.整体框架

​	结合gin框架和redis使用，开发了一个服务完成基本需求：创建和验证礼品码，实现三个接口1）管理后台调用-创建礼品码。2）管理后台调用-查询礼品码信息。3）客户调用-验证礼品码	

## 2.目录结构

```
.
├── Readme.md
├── controller
│   ├── StatusErr.go #状态码
│   └── controller.go
├── gift.png
├── giftmodify.png
├── go.mod
├── go.sum
├── http
│   ├── main
│   └── main.go
├── model
│   └── model.go
├── redis
│   └── redis.go
├── router
│   └── router.go  #路由 
├── service  #业务逻辑层
│   ├── rand.go
│   ├── service.go
│   └── service_test.go
└── test
    ├── __pycache__
    │   └── locust_test.cpython-39.pyc
    ├── locust_gift.html
    └── locust_test.py

```



## 3.代码逻辑分层

​	

| 层        | 文件夹          | 主要职责               | 调用关系               | 其他说明     |
| --------- | --------------- | ---------------------- | ---------------------- | ------------ |
| 应用层    | /http/main.go   | 程序启动               | 调用控制层，和redis层  |              |
| redis层   | /redis/redis.go | 初始化连接redis        | 被服务层和应用层调用   |              |
| 控制层    | /controller     | 处理请求和构建回复消息 | 被路由调用，调用服务层 |              |
| service层 | /service/       | 业务逻辑实现           | 被控制层调用           | 个层互不调用 |
| model层   | /model          | 数据模型               | 被服务层所调用         |              |
| router层  | /router         | 路由器                 | 被应用层所调用         |              |



## 4.存储设计

数据库存储信息

| 内容             | 数据库 | key  | 类型   | 说明       |
| ---------------- | ------ | ---- | ------ | ---------- |
| 礼品码及相关信息 | redis  | code | string | "HVALJMX9" |

保存礼品码所需信息

| 内容                            | field          | 类型              |
| ------------------------------- | -------------- | ----------------- |
| 礼品描述信息                    | Description    | string            |
| 礼品码类型                      | Type           | int               |
| 可领取用户                      | ReceivingUser  | string            |
| 可领取次数                      | AvailableTimes | string            |
| 有效期                          | ValidPeriod    | string            |
| 礼品码被创建的时间              | CreatTime      | string            |
| 创建这个礼品码的用户            | CreateUser     | string            |
| 存储礼品包内容                  | GiftPackages   | []GiftPackage     |
| 存储已经领取过该礼品码的用户    | ReceivedUsers  | Map[string]string |
| 礼品码已经被领取过的次数        | ReceivedCount  | int               |
| 存储内部生成的礼品码            | Code           | string            |
| 礼品码过期时间，内部 redis 使用 | Expiration     | Time.Duration     |

GiftPackage存储信息

| 内容     | field | 类型   |
| -------- | ----- | ------ |
| 礼品名字 | Name  | string |
| 礼品数量 | Num   | int    |

验证礼品码的请求实体

| 内容   | field | 类型   |
| ------ | ----- | ------ |
| 礼品码 | Code  | string |
| 用户   | User  | string |



## 5.接口设计

### 1.创建礼品码

### 	请求方式

​		http post

### 	接口地址

​		http://localhost:8080/create_gift_code

### 	请求参数

```json
{
  "type": 1,
  "receiving_user": "alms",
  "valid_period": 1625628980,
  "create_user": "admin",
  "available_times": 2,
  "description": "测试type1",
  "gift_packages": [
    {
      "name": "金币",
      "num": 10
    },
    {
      "name": "钻石",
      "num": 20
    }
  ]
}
```

### 请求响应

```json
{
    "Code": 201,
    "Message": "成功创建",
    "Data": "V5KX5NKC"
}
```

### 	响应状态吗	

| 状态码 | 说明                           |
| ------ | ------------------------------ |
| 201    | 创建成功                       |
| 400    | 没绑定参数成功                 |
| 401    | 请指定领取用户                 |
| 402    | 请指定可兑换次数               |
| 404    | 礼品码类型不合法               |
| 406    | 请输入有效期                   |
| 403    | 输入的礼品码有效期小于当前时间 |
| 411    | 请输入礼品码内容               |
| 409    | 五次创建礼品码失败，请重新创建 |
| 413    | 插入礼品码失败                 |

### 2.查询礼品码

### 	请求方式

​		http get

### 	接口地址

http://localhost:8080/query_gift_code

### 	请求参数

```json
code=0D45LM5R
```

### 请求响应

```json
{
    "Code": 200,
    "Message": "成功",
    "Data": {
        "description": "测试type1",
        "type": 1,
        "receiving_user": "alms",
        "available_times": 1,
        "valid_period": 1725628980,
        "create_time": 1625725957,
        "create_user": "admin",
        "gift_packages": [
            {
                "name": "金币",
                "num": 10
            },
            {
                "name": "钻石",
                "num": 20
            }
        ],
        "received_users": null,
        "received_count": 0,
        "code": "V5KX5NKC"
    }
}
```

### 	响应状态吗	

| 状态码 | 说明                         |
| ------ | ---------------------------- |
| 200    | 成功                         |
| 411    | 请输入礼品码                 |
| 410    | 礼品码输入错误或礼品码已过期 |

### 3.验证礼品码

### 1.创建礼品码

### 	请求方式

​		http post

### 	接口地址

​		http://localhost:8080/verify_gift_code

### 	请求参数

```json
{
  "user": "alms",
  "code": "1ORE35SJ"
}
```

### 请求响应

```json
{		
  
  	"Code": 200,
    "Message": "成功",
  "Data":{
    "GiftPackages": [
        {
            "name": "金币",
            "num": 10
        },
        {
            "name": "钻石",
            "num": 20
        }
    ]
  }
}
```

### 	响应状态吗	

| 状态码 | 说明           |
| ------ | -------------- |
| 200    | 成功           |
| 400    | 没绑定成功     |
| 410    | 请输入礼品码   |
| 401    | 请输入领取用户 |
| 412    | 验证失败       |

### 	

## 6.第三方库

```
github.com/go-redis/redis v6.15.9+incompatible
```

## 7.如何编译执行

### 代码格式化

```
make fmt
```

### 代码静态检测

```
make vet
```

### 执行可执行文件

```
./main
```

程序详细处理

![gift](gift.png)

程序逻辑思想

![giftmodify](giftmodify.png)

