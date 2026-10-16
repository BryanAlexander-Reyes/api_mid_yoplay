package controllers

import (
	"api_mid_yoplay/helper"
	"encoding/json"
	beego "github.com/beego/beego/v2/server/web"
	"net/http"
)

type correoRegistro struct {
	Correo string `json:"correo"`
}

type MidCorreoRegistroController struct {
	beego.Controller
}

func (c *MidCorreoRegistroController) EnviarCorreoRegistro() {
	var request correoRegistro
	err := json.NewDecoder(c.Ctx.Request.Body).Decode(&request)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]string{"error": "Error al decodificar la solicitud"}
		c.ServeJSON()
		return
	}
	if request.Correo == "" {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]string{"error": "El correo es obligatorio"}
		c.ServeJSON()
		return
	}
	existe, err := helper.BuscarUsuarioCorreo(request.Correo)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{"error": "No FUE POSIBLE CONSULTAR EL CORREO "}
		c.ServeJSON()
		return
	}

	if existe {
		c.Ctx.Output.SetStatus(http.StatusConflict)
		c.Data["json"] = map[string]string{"error": "El correo ya existe"}
		c.ServeJSON()
		return
	}
	c.Data["json"] = map[string]string{
		"mensaje": "El correo no existe enviando codigo",
	}
	c.ServeJSON()

}
