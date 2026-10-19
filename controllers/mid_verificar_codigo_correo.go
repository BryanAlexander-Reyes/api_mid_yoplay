package controllers

import(
	"api_mid_yoplay/helper"
	"encoding/json"
	beego "github.com/beego/beego/v2/server/web"
	"net/http"
	
)

type VerificarCodigo struct{
	Correo string `json:"correo"`
	Codigo string `json:"codigo"`
	Proposito string `json:"proposito"`
}

type VerificarCodigoCorreo struct{
	beego.Controller
}

func (c* VerificarCodigoCorreo) VerificarCodigo(){
	var request VerificarCodigo
	err:= json.NewDecoder(c.Ctx.Request.Body).Decode(&request)
	if err!=nil{
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"]=map[string]string{
			"Mensaje":"No se pudo leer la solicitud",
			"error":err.Error(),
		}
		c.ServeJSON()
		return
	}

	if request.Correo==""{
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"]=map[string]string{
			"Mensaje":"El correo es obligatorio",
		}
		c.ServeJSON()
		return
	}
	if request.Codigo==""{
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"]=map[string]string{
			"Mensaje":"El codigo es obligatorio",
		}
		c.ServeJSON()
		return
	}
	if request.Proposito==""{
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"]=map[string]string{
			"Mensaje":"El proposito es obligatorio",
		}
		c.ServeJSON()
		return
	}
	if request.Proposito!="registro"&& request.Proposito!="recuperacion"{
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"]=map[string]string{
			"mensaje":"El proposito no es valido",
		}
		c.ServeJSON()
		return
	}

	
	valido := helper.VerificarCodigo(
		request.Correo,
		request.Codigo,
		request.Proposito,
	)

	if valido {
    respuesta := map[string]string{
        "mensaje": "El codigo de verificacion es valido",
    }

    if request.Proposito == "recuperacion" {
        idUsuario, err := helper.ObtenerIDUsuarioCorreo(request.Correo)
        if err != nil {
            c.Ctx.Output.SetStatus(http.StatusInternalServerError)
            c.Data["json"] = map[string]string{
                "mensaje": "Error al obtener usuario",
            }
            c.ServeJSON()
            return
        }

        token, err := helper.GenerarTokenRecuperacion()
        if err != nil {
            c.Ctx.Output.SetStatus(http.StatusInternalServerError)
            c.Data["json"] = map[string]string{
                "mensaje": "No fue posible generar token de recuperacion",
            }
            c.ServeJSON()
            return
        }

        helper.GuardarTokenRecuperacion(token, idUsuario)
        respuesta["token"] = token
    }

    c.Data["json"] = respuesta
    c.ServeJSON()
    return
}

c.Ctx.Output.SetStatus(http.StatusBadRequest)
c.Data["json"] = map[string]string{
    "mensaje": "El codigo de verificacion es invalido",
}
c.ServeJSON()
return}