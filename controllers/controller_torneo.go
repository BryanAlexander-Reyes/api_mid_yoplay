package controllers

import (
	"api_mid_yoplay/models"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/astaxie/beego/orm"
	"github.com/beego/beego/v2/server/web"
)

type ReglamentoController struct{
	web.Controller 
}

func (c *ReglamentoController)CrearTorneos(){
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

	if torneo.CantidadEquipos==0{
		c.Data["json"]=map[string]interface{}{
			"ok":false,
			"Mensaje":"La cantidad de equipos no puede ser 0",
		}
		c.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		c.ServeJSON()
		return
	}
	
	if len(torneo.PremiacionPuesto)<1{
		c.Data["json"]=map[string]interface{}{
			"ok":false,
			"Mensaje":"La premiacion tiene que tener minimo una premiacion",
		}
		c.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		c.ServeJSON()
		return
	}

	if torneo.IDDeporte == nil{
				c.Data["json"]=map[string]interface{}{
			"ok":false,
			"Mensaje":"el tipo de deporte es obligatorio",
		}
		c.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		c.ServeJSON()
		return
	}

	o:=orm.NewOrm()

	Tipodeporte:= models.TipoDeporte{}


	err = o.QueryTable("TipoDeporte").
		Filter("IDdeporte",torneo.IDDeporte.IDDeporte).
		One(&Tipodeporte)

	if err !=nil{
		c.Data["json"]=map[string]interface{}{
			"ok":false,
			"Mensaje":"el tipo de deporte no existe",
		}
		c.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		c.ServeJSON()
		return
	}

	
	
}

