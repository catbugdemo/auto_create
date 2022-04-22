package test

import (
	"fmt"
	"github.com/catbugdemo/auto_create/test/model"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"testing"
)

func InitSqlx() *sqlx.DB {
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s password=%s", "1.117.233.151", "5432", "ytf", "smb", "disable", "ytf@2021"))
	if err != nil {
		panic(err)
	}
	return db
}

func TestCRUD(t *testing.T) {
	db := InitSqlx()
	svdl := model.SmbVocherDetailLog{
		SmbVocherDetailId: 2,
		SmbVocherId:       2,
	}
	t.Run("insert", func(t *testing.T) {
		err := svdl.Insert(db)
		assert.Nil(t, err)
	})
	t.Run("find", func(t *testing.T) {
		find, err := svdl.Find(db)
		assert.Nil(t, err)
		fmt.Println(find)
	})

	t.Run("find", func(t *testing.T) {
		err := svdl.First(db)
		assert.Nil(t, err)
		fmt.Println(svdl)
	})

	t.Run("update", func(t *testing.T) {
		err := svdl.Update(db, "where id=2")
		assert.Nil(t, err)
	})

	t.Run("delete", func(t *testing.T) {
		err := svdl.Delete(db, "where id=2")
		assert.Nil(t, err)
	})

}
