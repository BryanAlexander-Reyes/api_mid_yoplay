package main

import (
	_ "api_mid_yoplay/routers"
	beego "github.com/beego/beego/v2/server/web"
    cors "github.com/beego/beego/v2/server/web/filter/cors"

)

func main() {
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowOrigins: []string {"http://localhost:4200"},
		AllowMethods: []string {"PUT", "PATCH", "GET", "OPTIONS", "DELETE"},
		AllowHeaders: []string{
			"Origin",
			"x-requested-with",
			"content-type",
			"accept",
			"origin",
			"authorization",
			"x-csrftoken"},
		ExposeHeaders: []string{"Content-Length"},
		AllowCredentials:true,
	}))
	beego.Run()
}

