package controllers

import (
	"api_mid_yoplay/helper"
	"encoding/json"
	"net/http"
	beego "github.com/beego/beego/v2/server/web"
)

type RecuperarContrasena struct{
	Correo string `json:"correo"`
	ContrasenaNueva string `json:"contrasena_nueva"`
}

type MidRecuperarContrasenaController struct{
	beego.Controller}

func (c* MidRecuperarContrasenaController) RecuperarContrasena(){
	var request RecuperarContrasena
	err:= json.NewDecoder(c.Ctx.Request.Body).Decode(&request)
	if err!=nil{
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"]=map[string]string{
			"mensaje":"no se pudo leer la solicitud",
			"error":err.Error(),
		}
		c.ServeJSON()
		return
	}
	existe, err:= helper.BuscarUsuarioCorreo(request.Correo)
	if err!=nil{
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"]=map[string]string{
			"mensaje":"no fue posible consultar el correo",
			"error":err.Error(),
		}
		c.ServeJSON()
		return
	}
	if !existe{
		c.Ctx.Output.SetStatus(http.StatusNotFound)
		c.Data["json"]=map[string]string{
			"mensaje":"el correo no existe",
		}
		c.ServeJSON()
		return
	}
	err= helper.EnviarCodigoVerificacion(request.Correo, "recuperacion",)
	if err!=nil{
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"]=map[string]string{
			"mensaje":"no fue posible enviar el codigo de verificacion",
			"error":err.Error(),
		}
		c.ServeJSON()
		return
	}
	c.Data["json"]=map[string]string{
		"mensaje":"el codigo de verificacion fue enviado correctamente",
	}
	c.ServeJSON()
}