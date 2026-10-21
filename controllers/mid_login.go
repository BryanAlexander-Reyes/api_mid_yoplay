package controllers

import (
	"api_mid_yoplay/helper"
	"encoding/json"

	"net/http"

	beego "github.com/beego/beego/v2/server/web"
	"golang.org/x/crypto/bcrypt"
)

type Login struct {
	Correo     string `json:"correo"`
	Contrasena string `json:"contrasena"`
}

type MidLoginController struct {
	beego.Controller
}

func (c *MidLoginController) IniciarSeSion() {
	var request Login
	err := json.NewDecoder(c.Ctx.Request.Body).Decode(&request)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]string{
			"mensaje": "No se pudo leer la peticion",
		}
		c.ServeJSON()
		return
	}
	if request.Contrasena == "" || request.Correo == "" {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]string{
			"mensaje": "es obligatorio el correo o la contraseña",
		}
		c.ServeJSON()
		return
	}
	idUsuario, err := helper.ObtenerIDUsuarioCorreo(request.Correo)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusUnauthorized)
		c.Data["json"] = map[string]string{
			"mensaje": "No se encontró el usuario",
		}
		c.ServeJSON()
		return
	}
	contrasena, err := helper.ObtenerIdContrasena(idUsuario)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{
			"mensaje": "No fue posible obtener la contraseña del usuario",
		}
		c.ServeJSON()
		return
	}
	hashContrasena, err := helper.ObtenerHashContrasena(contrasena)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{
			"mensaje": "No fue posible obtener el hash de la contrasena",
		}
		c.ServeJSON()
		return
	}
	err = bcrypt.CompareHashAndPassword(
		[]byte(hashContrasena),
		[]byte(request.Contrasena),
	)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusUnauthorized)
		c.Data["json"] = map[string]string{
			"mensaje": "Correo o contraseña incorrectos ingrese un correo valido o una contraseña valida",
		}
		c.ServeJSON()
		return
	}
	c.Ctx.Output.SetStatus(http.StatusOK)
	c.Data["json"]=map[string]interface{}{
		"mensaje":"Inicio de sesión exitosa",
		"id_usuario":idUsuario,
	}
	c.ServeJSON()
}
