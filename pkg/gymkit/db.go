package gymkit

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
	"github.com/spf13/cast"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const MySQLType = "MYSQL"
const PostgresType = "POSTGRES"

type DatabaseConfig struct {
	Type           string `json:"type"`
	Username       string `json:"username"`
	Password       string `json:"password"`
	DatabaseHost   string `json:"database_host"`
	DatabasePort   string `json:"database_port"`
	DatabaseName   string `json:"database_name"`
	DatabaseConfig string `json:"database_config"`
	MaxLifeTime    uint64 `json:"max_life_time"`
	MaxIdleConn    uint64 `json:"max_idle_conn"`
	MaxOpenConn    uint64 `json:"max_open_conn"`
}

// DatabaseBase builds a DatabaseConfig from environment variables.
func DatabaseBase(databaseType string) *DatabaseConfig {
	if databaseType != MySQLType && databaseType != PostgresType {
		panic(fmt.Sprintf("database type [%s] not recognized", databaseType))
	}

	cfg := &DatabaseConfig{Type: databaseType}

	typeKey := databaseType
	cfg.Username = GetEnv(fmt.Sprintf("USERNAME_DB_%s", typeKey), os.Getenv("USERNAME_DB"))
	cfg.Password = GetEnv(fmt.Sprintf("PASSWORD_DB_%s", typeKey), os.Getenv("PASSWORD_DB"))
	cfg.DatabaseHost = GetEnv(fmt.Sprintf("DATABASE_HOST_%s", typeKey), os.Getenv("DATABASE_HOST"))
	cfg.DatabasePort = GetEnv(fmt.Sprintf("DATABASE_PORT_%s", typeKey), os.Getenv("DATABASE_PORT"))
	cfg.DatabaseName = GetEnv(fmt.Sprintf("DATABASE_NAME_%s", typeKey), os.Getenv("DATABASE_NAME"))
	cfg.DatabaseConfig = GetEnv(fmt.Sprintf("DATABASE_CONFIG_%s", typeKey), os.Getenv("DATABASE_CONFIG"))
	cfg.MaxLifeTime = cast.ToUint64(os.Getenv("SET_CONN_MAX_LIFETIME"))
	cfg.MaxIdleConn = cast.ToUint64(os.Getenv("SET_MAX_IDLE_CONNS"))
	cfg.MaxOpenConn = cast.ToUint64(os.Getenv("SET_MAX_OPEN_CONNS"))

	return cfg
}

func (c *DatabaseConfig) Validate() error {
	return validation.ValidateStruct(c,
		validation.Field(&c.Type, validation.Required),
		validation.Field(&c.Username, validation.Required),
		validation.Field(&c.Password, validation.Required),
		validation.Field(&c.DatabaseHost, validation.Required),
		validation.Field(&c.DatabasePort, validation.Required),
		validation.Field(&c.DatabaseName, validation.Required),
		validation.Field(&c.MaxLifeTime, validation.Required),
		validation.Field(&c.MaxIdleConn, validation.Required),
		validation.Field(&c.MaxOpenConn, validation.Required),
	)
}

func (c *DatabaseConfig) GetConnection() string {
	if c.Type == MySQLType {
		return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s",
			c.Username, c.Password, c.DatabaseHost, c.DatabasePort, c.DatabaseName, c.DatabaseConfig)
	}
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s %s",
		c.DatabaseHost, c.Username, c.Password, c.DatabaseName, c.DatabasePort, c.DatabaseConfig)
}

func (c *DatabaseConfig) getLogger() logger.Interface {
	isDebug, _ := strconv.ParseBool(os.Getenv("DATABASE_DEBUG"))
	if !isDebug {
		return logger.Default
	}
	return logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: false,
			Colorful:                  true,
		},
	)
}

func (c *DatabaseConfig) SetConnection() (*gorm.DB, error) {
	conn := c.GetConnection()
	gormCfg := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{SingularTable: true},
		Logger:         c.getLogger(),
	}
	if c.Type == MySQLType {
		return gorm.Open(mysql.New(mysql.Config{DSN: conn}), gormCfg)
	}
	return gorm.Open(postgres.New(postgres.Config{DSN: conn, PreferSimpleProtocol: true}), gormCfg)
}

// GetMysqlConnection opens and configures a GORM MySQL connection.
func (c *DatabaseConfig) GetMysqlConnection() *gorm.DB {
	log.Println("Initializing MySQL database connection...")
	db, err := c.SetConnection()
	if err != nil {
		panic(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	sqlDB.SetConnMaxLifetime(time.Second * time.Duration(c.MaxLifeTime))
	sqlDB.SetMaxIdleConns(int(c.MaxIdleConn))
	sqlDB.SetMaxOpenConns(int(c.MaxOpenConn))
	log.Println("MySQL database connection initialized")
	return db
}

// GetPostgresConnection opens and configures a GORM PostgreSQL connection.
func (c *DatabaseConfig) GetPostgresConnection() *gorm.DB {
	log.Println("Initializing Postgres database connection...")
	db, err := c.SetConnection()
	if err != nil {
		panic(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	sqlDB.SetConnMaxLifetime(time.Second * time.Duration(c.MaxLifeTime))
	sqlDB.SetMaxIdleConns(int(c.MaxIdleConn))
	sqlDB.SetMaxOpenConns(int(c.MaxOpenConn))
	log.Println("Postgres database connection initialized")
	return db
}
