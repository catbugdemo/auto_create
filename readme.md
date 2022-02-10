# 自动生成 不再烦恼

## 自动生成 crud
```go


```


## 自动生成 models

```go

func TestModels(t *testing.T) {
	normal := Normal{
		DataSource: fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s password=%s", "127.0.0.1", "5432", "postgres", "mydb", "disable", "123456"),
		TableName:  "tmp_user",
		Driver:     "postgres",
	}

	generate, err := AutoGenerateModel(&normal)
	assert.Nil(t, err)
	fmt.Printf(generate)
}
```

## 自动生成 contrller 层
```go
// 生成 Controller 层
func TestGenerateController(t *testing.T) {
	c := Control{
		ControlName: "GetPages", // 输入创建名称
		Describe:    "获取落地页",    // 输入描述 -- 可不填
		
		// 是否需要绑定参数
		Req: Req{
			ReqBool: true, // 是否需要手动填写绑定参数，推荐3个以内为true
			Req:     "",   // 请求名称,如果 ReqBool == true 不填
		},
		DbConfig: "c.MustGet(DB_CONFIG).(*gorm.DB)",

		ServiceStr:     "",   // 一般不填 service层名称 一般只有 2 个返回 data,err
		ReturnDataBool: true, // 是否需要返回数据

		LogOrSave: "", // 默认不填
	}
	fmt.Println(GenerateController(c))
}
```

- 将会自动生成
```go
// 获取落地页
func GetPages(c *gin.Context) {
	type Req struct {
	// TODO 请填写请求参数
	}
	var req Req
	if err := c.Bind(&req); err != nil {
		log.Printf("%+v",errors.WithStack(err))
		c.JSON(200, gin.H{"code": 1, "msg": "request binding failed", "debug": err.Error()})
		return
	}

	db := c.MustGet(DB_CONFIG).(*gorm.DB)
	data, err := service.GetPages(req, db)
	if err != nil {
		log.Printf("%+v",errors.WithStack(err))
		c.JSON(200, gin.H{"code": 2, "msg": "service operate failed", "debug": err.Error()})
		return
	}
	
	log.Printf("way:%v ; req:%v ; data:%v ;", "GetPages", req, data)
	c.JSON(200, gin.H{"code": 0, "msg": "success", "data": data})
}
```