package routers

import (
	"api_mid_yoplay/controllers"
	beego "github.com/beego/beego/v2/server/web"
)

func init() {
    beego.Router("/", &controllers.MainController{})
	beego.Router("/api/creacion-usuario",&controllers.RegistroController{}, "post:CrearUsuario")
}
