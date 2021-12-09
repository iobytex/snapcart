package postgres

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"snapcart/config"
)

func InitDBConnection(config *config.Config) (*gorm.DB,error) {
	return gorm.Open(postgres.New(postgres.Config{
		DSN: fmt.Sprintf("host= %s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Shanghai",
			config.Postgresql.PostgresqlHost,
			config.Postgresql.PostgresqlUser,
			config.Postgresql.PostgresqlPassword,
			config.Postgresql.PostgresqlDbname,
			config.Postgresql.PostgresqlPort,
			config.Postgresql.PostgresqlSslMode),
	}),&gorm.Config{})
}
