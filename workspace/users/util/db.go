package util

import (
	"log"
	"os"
	"strconv"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gold-gym/gymkit"
	"gorm.io/gorm"
)

var (
	dbInstance *gorm.DB
	dbOnce     sync.Once
)

func GetDBConnection() *gorm.DB {
	dbOnce.Do(func() {
		log.Println("Initializing users MySQL connection...")

		maxLifeTime, _ := strconv.Atoi(os.Getenv("SET_CONN_MAX_LIFETIME_MYSQL"))
		if maxLifeTime == 0 {
			maxLifeTime = 60
		}
		maxIdleConns, _ := strconv.Atoi(os.Getenv("SET_MAX_IDLE_CONNS_MYSQL"))
		if maxIdleConns == 0 {
			maxIdleConns = 20
		}
		maxOpenConns, _ := strconv.Atoi(os.Getenv("SET_MAX_OPEN_CONNS_MYSQL"))
		if maxOpenConns == 0 {
			maxOpenConns = 100
		}

		cfg := gymkit.DatabaseBase(gymkit.MySQLType)
		cfg.MaxLifeTime = uint64(maxLifeTime)
		cfg.MaxIdleConn = uint64(maxIdleConns)
		cfg.MaxOpenConn = uint64(maxOpenConns)
		cfg.DatabaseConfig = "charset=utf8mb4&parseTime=True&loc=Local"

		dbInstance = cfg.GetMysqlConnection()
		log.Println("users MySQL connection ready")
	})
	return dbInstance
}
