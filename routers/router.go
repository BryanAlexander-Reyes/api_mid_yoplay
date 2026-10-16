package routers

import (
	"api_mid_yoplay/controllers"
	beego "github.com/beego/beego/v2/server/web"
)

// ADVERTENCIA PARA TODOS LOS DESARROLLADORES:
// NO eliminar rutas de otros desarrolladores. Si hay rutas que generan conflicto en la ejecución del programa, por favor comentarear la ruta, y al momento de hacer push, descomentarear la misma.

func init() {
    beego.Router("/", &controllers.MainController{})
	beego.Router("api/v1/encuentro/:id", &controllers.Encuentro_midController{}, "get:GetOne")
	beego.Router("api/v1/encuentro", &controllers.Encuentro_midController{}, "get:GetAll")
	beego.Router("api/v1/torneo/:id",&controllers.TorneoController{}, "get:CrearTorneos")
	beego.Router("api/v1/info-torneo/", &controllers.Info_torneo_midController{}, "get:GetAll")
	beego.Router("api/v1/info-torneo/:id", &controllers.Info_torneo_midController{}, "get:GetOne")
	beego.Router("/api/creacion-usuario",&controllers.RegistroController{}, "post:CrearUsuario")
	beego.Router("/api/enviarcorreoregistro",&controllers.MidCorreoRegistroController{}, "post:EnviarCorreoRegistro")
}
