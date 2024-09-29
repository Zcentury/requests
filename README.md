### 介绍
一个方便好用的 **Go Requests** 库

### 安装

```shell
go get -u github.com/Zcentury/requests
```

### 使用

> GET

```go
resp := requests.Get("https://www.baidu.com/")
```

> POST

```go
requests.Post("https://example.com", params.BodyString("user=xxx&pass=xxx"))

requests.Post("https://example.com", params.BodyMap2String{
    "user": "xxx",
    "pass": "xxx",
})

requests.Post("https://example.com", params.BodyJsonString("{\"user\":\"xxx\",\"pass\":\"xxx\"}"))

requests.Post("https://example.com", params.BodyMap2Json{
    "user": "xxx",
    "pass": "xxx",
})
```



> 自动保存 **Cookie**

```go
session := requests.Session()

get := session.Get(method.GET, params.Url("http://www.baidu.com/"), headers)
if get != nil {
    fmt.Println(get.Response.Status)
    fmt.Println(get.Response.Cookies())
}

get1 := session.Post("http://www.baidu.com/", headers, params.BodyMap2Json{
    "user": "xxx",
    "pass": "xxx",
})
if get1 != nil {
    fmt.Println(get1.Response.Status)
    fmt.Println(get1.Response.Cookies())
}
```

