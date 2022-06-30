## 自动生成 crud
```go


```


## 自动生成 models -- 支持 sqlx 

```go
func TestSqlx(t *testing.T) {
	normal := Normal{
		//DataSource: fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s password=%s", "1.117.233.151", "5432", "ytf", "smb", "disable", "ytf@2021"),
		DataSource: fmt.Sprint("new_retailers:mnHL63mzzNX2GyYd@tcp(1.15.221.224:3306)/new_retailers"),
		TableName:  "adv_shop_online_customer",
		Driver:     "mysql",
	}

	generate, err := AutoGenerateSqlx(&normal)
	assert.Nil(t, err)
	fmt.Println(generate)
}

```

## 自动生成 controller -- 支持 sqlx 
```go

func TestSqlxCrud(t *testing.T) {
	type AdvShopOnlineCustomer struct {
		Id          int       `gorm:"column:id;default:" json:"id" form:"id" db:"id"`                                             // 自增id
		CreateTime  time.Time `gorm:"column:create_time;default:" json:"create_time" form:"create_time" db:"create_time"`         // 创建时间
		Name        string    `gorm:"column:name;default:" json:"name" form:"name" db:"name"`                                     // 商户名称
		NickName    string    `gorm:"column:nick_name;default:" json:"nick_name" form:"nick_name" db:"nick_name"`                 // 商户昵称
		IconUrl     string    `gorm:"column:icon_url;default:" json:"icon_url" form:"icon_url" db:"icon_url"`                     // 头像url
		IconMediaId string    `gorm:"column:icon_media_id;default:" json:"icon_media_id" form:"icon_media_id" db:"icon_media_id"` // 微信头像id
		ShopId      int       `gorm:"column:shop_id;default:" json:"shop_id" form:"shop_id" db:"shop_id"`                         // 商户id
		BusinessId  int       `gorm:"column:business_id;default:" json:"business_id" form:"business_id" db:"business_id"`         // 商户在线id
		Status      int       `gorm:"column:status;default:" json:"status" form:"status" db:"status"`                             // 1:正常 2:删除
	}
	st := St{
		Stru:     AdvShopOnlineCustomer{},
		DbConfig: `c.MustGet(DB_CONFIG).(*sqlx.DB)`,
		Info: map[string]interface{}{
			"tag": "测试",
		},
	}
	fmt.Println(AutoGenerateSqlxControl(&st))
}

```

## 需要导入包
```go
/utils/model_utils.go
// 如果是 mysql -- 将 model_utils.go 中 mysql  注释释放

// 例 
// mysql 
// params = append(params, "?")
```