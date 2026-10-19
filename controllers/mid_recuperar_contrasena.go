package controllers

import (
	"api_mid_yoplay/helper"
	"encoding/json"
	
	"net/http"

	beego "github.com/beego/beego/v2/server/web"
)

type RecuperarContrasena struct {
	Correo string `json:"correo"`
}

type CambiarContrasena struct {
	Token           string `json:"token"`
	ContrasenaNueva string `json:"contrasenaNueva"`
}

type MidRecuperarContrasenaController struct {
	beego.Controller
}

func (c *MidRecuperarContrasenaController) RecuperarContrasena() {
	var request RecuperarContrasena
	err := json.NewDecoder(c.Ctx.Request.Body).Decode(&request)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]string{
			"mensaje": "no se pudo leer la solicitud",
			"error":   err.Error(),
		}
		c.ServeJSON()
		return
	}
	existe, err := helper.BuscarUsuarioCorreo(request.Correo)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{
			"mensaje": "no fue posible consultar el correo",
			"error":   err.Error(),
		}
		c.ServeJSON()
		return
	}
	if !existe {
		c.Ctx.Output.SetStatus(http.StatusNotFound)
		c.Data["json"] = map[string]string{
			"mensaje": "el correo no existe",
		}
		c.ServeJSON()
		return
	}
	err = helper.EnviarCodigoVerificacion(request.Correo, "recuperacion")
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{
			"mensaje": "no fue posible enviar el codigo de verificacion",
			"error":   err.Error(),
		}
		c.ServeJSON()
		return
	}
	c.Data["json"] = map[string]string{
		"mensaje": "el codigo de verificacion fue enviado correctamente",
	}
	c.ServeJSON()
}

func (c *MidRecuperarContrasenaController) CambiarContrasena() {
	var request CambiarContrasena
	err := json.NewDecoder(c.Ctx.Request.Body).Decode(&request)

	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]string{
			"mensaje": "no se pudo leer la solicitud",
		}
		c.ServeJSON()
		return
	}
	if request.Token == "" || request.ContrasenaNueva == "" {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]string{
			"mensaje": "el token y la nueva contraseña son obligatorios",
		}
		c.ServeJSON()
		return
	}
	idUsuario, err := helper.ValidarTokenRecuperacion((request.Token))
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]string{
			"mensaje": "el token ha expirado",
		}
		c.ServeJSON()
		return
	}
	idContrasena, err := helper.ObtenerIdContrasena((idUsuario))
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{
			"mensaje": "No fue posible obtener el ID de la contraseña",
			"error":   err.Error(),
		}
		c.ServeJSON()
		return
	}
	hashTexto, err := helper.GenerarHashContrasena(request.ContrasenaNueva)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{
			"mensaje": "No fue posible procesar la nueva contraseña",
			"error":   err.Error(),
		}
		c.ServeJSON()
		return
	}
	err = helper.ActualizarContrasena(idContrasena, hashTexto)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{
			"mensaje": "No fue posible actualizar la contraseña nueva",
			"error":   err.Error(),
		}
		c.ServeJSON()
		return
	}

	c.Ctx.Output.SetStatus(http.StatusOK)
	c.Data["json"] = map[string]string{
		"mensaje": "la contraseña se actualizó correctamente",
		"hash":hashTexto,
	}
	c.ServeJSON()
}
