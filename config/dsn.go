package config

import (
	"fmt"
	"os"
)

func DSN() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s",
		os.Getenv("MYSQL_USER"),     // mysql
		os.Getenv("MYSQL_PASSWORD"), // mysqlpass
		os.Getenv("MYSQL_HOST"),     // db
		os.Getenv("MYSQL_PORT"),     // 3306
		os.Getenv("MYSQL_DATABASE"), // ddd_go
	) + "?parseTime=true&collation=utf8mb4_bin"
}
