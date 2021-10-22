package auto

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestModels(t *testing.T) {
	t.Run("models", func(t *testing.T) {
		normal := Normal{
			DataSource: fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s password=%s", "10.0.10.209", "5432", "ytf", "ytf", "disable", "ytf@2021"),
			TableName:  "ytf_adv_industry_people",
			Driver:     "postgres",
		}

		generate, err := AutoGenerateModel(&normal)
		assert.Nil(t, err)

		fmt.Printf(generate)
	})
}

func TestCRUD(t *testing.T) {
	type YtfAdvIndustryPeople struct {
		Id         int       `gorm:"column:id;default:" json:"id" form:"id"`
		CreatedAt  time.Time `gorm:"column:created_at;default:" json:"created_at" form:"created_at"`
		UpdatedAt  time.Time `gorm:"column:updated_at;default:" json:"updated_at" form:"updated_at"`
		UserId     int       `gorm:"column:user_id;default:" json:"user_id" form:"user_id"`
		IndustryId int       `gorm:"column:industry_id;default:" json:"industry_id" form:"industry_id"`
		PeopleId   int       `gorm:"column:people_id;default:" json:"people_id" form:"people_id"`
	}
	st := St{
		Stru:        YtfAdvIndustryPeople{},
		DbConfig:    `c.MustGet(DB_CONFIG).(*gorm.DB)`,
		ModelsName:  "models.YtfAdvIndustryPeople",
		RedisConfig: `c.MustGet(REDIS_TOKEN).(*redis.Pool).Get()`,
		Handlers:    "handlers",
	}

	crud := AutoGenerateCRUD(st)
	fmt.Printf(crud)
}
