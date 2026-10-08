package routers

import (
	"api_mid_yoplay/controllers"
	beego "github.com/beego/beego/v2/server/web"
)

func init() {
    beego.Router("/", &controllers.MainController{})
	beego.Router("api/v1/encuentro/:id", &controllers.Encuentro_midController{}, "get:GetOne")
	beego.Router("api/v1/encuentro", &controllers.Encuentro_midController{}, "get:GetAll")
}
