package main

import (
	"easyquery"
	"easyquery/examples/pkg/db"
	"easyquery/examples/pkg/user"
	"fmt"
	"strings"
)

func main() {
	db.InitDB()
	defer db.CloseDB()
	current := db.Postgres
	var vos []easyquery.GroupVO
	var models []user.User
	current.Config.DryRun = true
	current = current.Model(&user.User{}).Joins("Company")
	current.Find(&models)
	sql := current.Statement.SQL.String()
	start := strings.Index(sql, "LEFT")
	end := strings.Index(sql, "WHERE")
	join := sql[start:end]
	fmt.Println("join: ", join)
	fmt.Println("sql: ", sql)

	current.DryRun = false
	current.Statement.SQL.Reset()
	current.
		Select(fmt.Sprintf("%s as label, count(1) as count", `"Company"."name"`)).
		Group(`"Company"."name"`).
		Find(&vos)
	//current.Model(&user.User{}).Joins("Company").Group(`"Company"."name"`).
	//	Select(fmt.Sprintf("%s as label, count(1) as count", `"Company"."name"`)).
	//	Find(&vos)
	fmt.Println(vos)
	db.Postgres.Model(&user.User{}).Find(&models)
	fmt.Println(models)
}
