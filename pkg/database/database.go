package database

import (
	"fmt"
	"github.com/SoftclubIT/todo-service/pkg/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Init(host, user, password, name, port string) (db *gorm.DB, err error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Dushanbe ",
		host,
		user,
		password,
		name,
		port,
	)
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return
	}

	err = db.AutoMigrate(
		&models.Task{},
	)
	if err != nil {
		return nil, err
	}

	return db, err
}
