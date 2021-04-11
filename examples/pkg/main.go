package main

import (
	"easyquery/examples/pkg/db"
	"easyquery/examples/pkg/router"
)

func main() {
	db.InitDB()
	defer db.CloseDB()
	router.Register()
}
