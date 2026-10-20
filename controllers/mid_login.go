package controllers

import (
	"api_mid_yoplay/helper"
	"encoding/json"
	
	"net/http"

	beego "github.com/beego/beego/v2/server/web"
)

type Login struct{
	Correo string `json:"correo"`
	Contrasena string `json:"contrasena"`
}

type midLoginController struct{
	beego.Controller
}

func (c *midLoginController) IniciarSecion(){
	var request Login
	err:= json.NewDecoder(c.Ctx.Request.Body).Decode(&request)
	if err!=nil{
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"]=map[string]string{
			"mensaje":"No se pudo leer la peticion",
		}
		c.ServeJSON()
		return
	}
	if request.Contrasena== "" || request.Correo==""{
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"]=map[string]string{
			"mensaje":"es obligatorio el correo o la contraseña",
		}
	}
	
}
