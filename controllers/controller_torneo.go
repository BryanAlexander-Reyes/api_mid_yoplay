package controllers

import (
	"api_mid_yoplay/models"
	"encoding/json"
	"net/http"
	"strings"
	"github.com/beego/beego/v2/server/web"
)
//@Title Enviar Torneo
//@Description Formulario para la creacion de torneo para el usuario
//@param boby body models.CrearTorneo true "Torneo guardado para la creacion"
//@Success 200 {object} map[string]interface{} "torneo creado"
//@Failure 400 {object} map[string]interface{} "solicitud invalidad"
//@Failure 500 {object} map[string]interface{} "Error Interno"
//@ROUTER /api/torneo [post]
type TorneoController struct{
	web.Controller 
}

func (c *TorneoController)CrearTorneos(){
	body:=c.Ctx.Input.RequestBody

	if len(body)==0{
		c.Data["json"]=map[string]interface{}{
			"ok":false,
			"Mensaje":"El registro del torneo esta vacio",
		}
		c.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		c.ServeJSON()
		return
	}
	var torneo models.CrearTorneo
	err := json.Unmarshal(body,&torneo)

	if err !=nil{
		c.Data["json"]=map[string]interface{}{
			"ok":false,
			"Mensaje":"El registro de la solicitud no tiene el formato valido",
			"error":err.Error(),
		}
		c.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		c.ServeJSON()
		return
	}


	if strings.TrimSpace(torneo.Objetivo)==""{
		c.Data["json"]=map[string]interface{}{
			"ok":false,
			"Mensaje":"El objectivo es obligatorio",
		}
		c.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		c.ServeJSON()
		return
	}
	if strings.TrimSpace(torneo.Lugar)==""{
		c.Data["json"]=map[string]interface{}{
			"ok":false,
			"Mensaje":"El Lugar es obligatorio",
		}
		c.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		c.ServeJSON()
		return
	}

	if torneo.CantidadEquipos<1{
		c.Data["json"]=map[string]interface{}{
			"ok":false,
			"Mensaje":"La cantidad de equipos no puede ser 0",
		}
		c.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		c.ServeJSON()
		return
	}
	
	if len(torneo.Premiaciones)<1{
		c.Data["json"]=map[string]interface{}{
			"ok":false,
			"Mensaje":"La premiacion tiene que tener minimo una premiacion",
		}
		c.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		c.ServeJSON()
		return
	}

	if len(torneo.Tipo_Deporte)!=1{
				c.Data["json"]=map[string]interface{}{
			"ok":false,
			"Mensaje":"el tipo de deporte es obligatorio",
		}
		c.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		c.ServeJSON()
		return
	}


	if len(torneo.Divicion)!=1{
	c.Data["json"]=map[string]interface{}{
			"ok":false,
			"Mensaje":"debe seleccionar una divicion para el torneo ",
		}
		c.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		c.ServeJSON()
		return
	}
	

	for _, imagen := range torneo.ImagenUrl {
    if strings.TrimSpace(imagen.Imagen) == "" {
        c.Data["json"] = map[string]interface{}{
            "ok": false,
            "Mensaje": "La URL de la imagen no puede estar vacía",
        }
        c.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
        c.ServeJSON()
        return
    }
}


}

