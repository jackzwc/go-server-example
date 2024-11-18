package config

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	logger "github.com/sirupsen/logrus"
	"go-server-example/model"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"time"
)

var (
	AppConfig *Config
	Db        *gorm.DB
)

type Config struct {
	Server `yaml:"server"`
	DB     `yaml:"db"`
}

type Server struct {
	Port int64 `yaml:"port"`
}

type DB struct {
	Host         string `yaml:"host"`
	Port         string `yaml:"port"`
	UserName     string `yaml:"user_name"`
	Password     string `yaml:"password"`
	Database     string `yaml:"database"`
	Timeout      int64  `yaml:"timeout"`
	MaxOpenConns int    `yaml:"max_open_conns"`
	MaxIdleConns int    `yaml:"max_idle_conns"`
}

func InitConfig() {
	AppConfig = initAppConfig()
	Db = initSQLConfig()
}

// initAppConfig init app config
func initAppConfig() *Config {
	// Read config file
	configData, err := ioutil.ReadFile("./config.yml")
	if err != nil {
		logger.Errorf("Read config file err: %v", err)
		panic(err)
	}

	var config = &Config{}
	err = yaml.Unmarshal(configData, config)
	if err != nil {
		logger.Errorf("Unmarshal config file err: %v", err)
	}

	return config
}

// initSQLConfig init postgreSQL config
func initSQLConfig() *gorm.DB {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		AppConfig.DB.Host, AppConfig.DB.Port, AppConfig.DB.UserName, AppConfig.DB.Password, AppConfig.Database,
	)

	// connecting to a database
	db, err := gorm.Open("postgres", connStr)
	if err != nil {
		panic(fmt.Sprintf("connect database err: %v", err))
	}
	err = db.DB().Ping()
	if err != nil {
		logger.Error("DB ping fail:", err)
		return nil
	}

	logger.Info("Connect db success!")
	db.LogMode(true)
	db.DB().SetMaxIdleConns(AppConfig.MaxIdleConns)
	db.DB().SetMaxOpenConns(AppConfig.MaxOpenConns)

	db.AutoMigrate(model.User{})
	db.Model(&model.User{}).AddIndex("idx_number", "number")

	// Failed to reconnect
	go func() {
		timer := time.NewTicker(time.Duration(AppConfig.Timeout * int64(time.Second)))
		for {
			select {
			case <-timer.C:
				err = db.DB().Ping()
				if err != nil {
					logger.Error("DB connect fail:", err)
					logger.Info("Reconnect beginning...")
					// Connecting to a database
					db, err = gorm.Open("postgres", connStr)
					if err != nil {
						logger.Error("Connect database err:", err)
					}
					err = db.DB().Ping()
					if err != nil {
						logger.Error("DB ping fail:", err)
					}
					logger.Info("Reconnect db success!")
					db.LogMode(true)
					db.DB().SetMaxIdleConns(AppConfig.MaxIdleConns)
					db.DB().SetMaxOpenConns(AppConfig.MaxOpenConns)
				}
			}
		}
	}()

	return db
}
