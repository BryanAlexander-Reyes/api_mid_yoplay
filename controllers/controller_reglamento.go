package controllers

import (
	"api_mid_yoplay/models"
	"encoding/json"
	"net/http"

	"github.com/beego/beego/v2/server/web"
)


type ReglamentoController struct{
	web.Controller
}


func(c *ReglamentoController)CrearReglamento(){


	var reglamentos models.Reglamento


	err:=json.Unmarshal(
		c.Ctx.Input.RequestBody, &reglamentos,
	)

	if err !=nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		c.Data["json"]=map[string]interface{}{
			"ok":false,
			"mensaje":"el cuerpo de la solicitud no tiene formato valido",
			"error": err.Error(),
		}
		c.ServeJSON()
		return
	}

	if reglamentos.Idioma ==""{
		c.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		c.Data["json"] = map[string]interface{}{
			"ok":      false,
			"mensaje": "el idioma es obligatorio",
		}
		c.ServeJSON()
		return
	}

	if reglamentos.Estandar == ""{
		c.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		c.Data["json"] = map[string]interface{}{
			"ok":      false,
			"mensaje": "debe elejir un reglamento estandar",
		}
		c.ServeJSON()
		return
	}

	if len(reglamentos.Regla)>=1{
		c.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		c.Data["json"] = map[string]interface{}{
			"ok":      false,
			"mensaje": "debe a ver minimo una regla",
		}
		c.ServeJSON()
		return

	}



}
