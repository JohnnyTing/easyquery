package main

import (
	"easyquery/examples/pkg/db"
	"easyquery/examples/pkg/user"
)

func main() {
	// AutoMigrate
	db.InitDB()
	defer db.CloseDB()
	current := db.Postgres
	current.Migrator().DropTable(&user.User{})
	current.AutoMigrate(&user.User{})
	current.Migrator().DropTable(&user.Role{})
	current.AutoMigrate(&user.Role{})
}
