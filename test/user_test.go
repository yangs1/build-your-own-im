package test

import (
	"gim/global"
	"gim/initialize"
	"gim/models"
	"testing"
)

func TestUserDB(t *testing.T) {
	initialize.InitDB()
	err := global.DB.AutoMigrate(&models.UserBasic{})

	if err != nil {
		panic(err)
	}
}
