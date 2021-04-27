package application

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/smarest/smarest-account/domain/service"
	"github.com/smarest/smarest-account/infrastructure/persistence"
	"github.com/smarest/smarest-common/util"
	"gopkg.in/gorp.v3"
)

type Bean struct {
	LoginService      *LoginService
	UserService       *UserService
	RestaurantService *RestaurantService
	DbMap             *gorp.DbMap
}

func (bean *Bean) DestroyBean() {

	// Turn off tracing
	bean.DbMap.TraceOff()
	bean.DbMap.Db.Close()
}

func InitBean() (*Bean, error) {
	user := util.GetEnvDefault("DB_USER", "root")
	password := util.GetEnvDefault("DB_PASSWORD", "")
	//	host := getEnvWithDefault("DB_HOST", "127.0.0.1")
	//	port := getEnvWithDefault("DB_PORT", "3306")
	dbName := util.GetEnvDefault("DB_NAME", "smarest")
	//	dsn := fmt.Sprintf("%s:%s@unix(%s:%s)/%s?parseTime=true", user, password, host, port,dbName)
	dsn := fmt.Sprintf("%s:%s@unix(/Applications/XAMPP/xamppfiles/var/mysql/mysql.sock)/%s?parseTime=true", user, password, dbName)
	fmt.Printf("dns: %s", dsn)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	jwtKey := util.GetEnvDefault("POS_JWT_KEY", "b3BlbnNzaC1rZXktdjEAAAAACm")
	dbMap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}
	userRepository := persistence.NewUserRepository(dbMap)
	restaurantRepository := persistence.NewRestaurantRepository(dbMap)
	domainUserService := service.NewUserService(jwtKey, userRepository)
	userService := NewUserService(domainUserService)

	restaurantService := NewRestaurantService(service.NewRestaurantService("b3BlbnNzaC1rZXktdjEAAAAACmFlczI1Ni1jdHIAAAAGYmNyeXB0AAAAGAAAABD1WGaxt2", restaurantRepository))

	loginService := NewLoginService(
		util.GetEnvDefault("POS_HOME_PAGE", "http://pos.server.vn/pos/pos-portal"),
		strings.Split(util.GetEnvDefault("POS_DOMAINS", "localhost:2020"), ","),
		util.GetEnvDefault("smarest-account_TOKEN", "pos_access_token"),
		domainUserService)
	// Will log all SQL statements + args as they are run
	// The first arg is a string prefix to prepend to all log messages
	dbMap.TraceOn("[gorp]", log.New(os.Stdout, "myapp:", log.Lmicroseconds))

	return &Bean{LoginService: loginService, UserService: userService, DbMap: dbMap, RestaurantService: restaurantService}, nil
}
