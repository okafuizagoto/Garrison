module members

go 1.24

require (
	github.com/go-sql-driver/mysql v1.9.3
	github.com/gold-gym/gymkit v0.0.0
	github.com/golang-jwt/jwt/v5 v5.2.2
	github.com/joho/godotenv v1.5.1
	github.com/kataras/iris/v12 v12.2.11
	github.com/rabbitmq/amqp091-go v1.10.0
	github.com/spf13/cast v1.9.2
	gorm.io/driver/mysql v1.6.0
	gorm.io/gorm v1.30.0
)

replace github.com/gold-gym/gymkit => /pkg/gymkit
