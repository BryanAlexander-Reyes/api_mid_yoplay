package main

import (
	_ "api_mid_yoplay/routers"
	beego "github.com/beego/beego/v2/server/web"
)

func main() {
	beego.Run()
}

