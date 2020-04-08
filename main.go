package main

import (
	"fmt"
	"github.com/IdentitySystem/database"
	"github.com/IdentitySystem/web"
	"github.com/IdentitySystem/web/controller"
	"github.com/IdentitySystem/web/stateverifier"
)

func main()  {
	//fabricserivc:=service.ServicStup()
	db,err:=database.Connesql()
	if err!=nil{
		fmt.Println(err.Error())
	}
	dblogic:=&database.DBlogic{DB:db,Name:"identity"}
	verifier:=stateverifier.NewVerifier()
	verifier.Register(controller.VERLOGINED,web.Verlogined)
	verifier.Register(controller.VERISADMIN,web.VerIsAdmin)
	app:=&controller.Application{nil,dblogic,verifier}
	web.WebStart(app)
}
