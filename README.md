# btcp
可扩展的TCP服务器

- 支持头解析器注册
- 支持body解析器注册
- 支持注册日期解析器


# server

- 连接管理
- 拦截器
- 中间件
- 过滤器
- 路由
- 路由组

# 连接
- 元数据
- 心跳管理

# 数据包流图
- 读包
- 完整包
- 解析
- 中间件

# 经典问题
- 粘包拆包
- 连接管理
- 鉴权
- 鉴权回调


# todo
- 心跳管理
- 鉴权
- 

```go
func main(){
	s := btcp.NewServer(
	btcp.WithPort(6666)
		)
	// 使用中间件
	s.Use(btcp.Recover)
	S.Use(func(req Request){
		// 前置中间件
		req.Next()
    })

    S.Use(func(req Request){
        req.Next()
        // 后置中间件
    })
	s.Handler(1, func(req Request){
		
		
    })
	if err := s.Run();err != nil {
	    panic(err)	
    }
	
}

```