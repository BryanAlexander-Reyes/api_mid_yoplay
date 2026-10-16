package controllers

import (
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
	err:= json.NewDecoder(c.Ctx.Request.Body).Decode(&request)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]string{"error": "Error al decodificar la solicitud"}
		c.ServeJSON()
		return
	}
	if request.correo ==""{
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]string{"error": "El correo es obligatorio"}
		c.ServeJSON()
		return
	}

	
}