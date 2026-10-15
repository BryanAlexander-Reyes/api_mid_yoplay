package controllers

import (
	"encoding/json"
	"net/http"
	"net/mail"
	"strings"

	"github.com/beego/beego/v2/server/web"

	"api_mid_yoplay/models"
	"api_mid_yoplay/services"
)

type CorreoController struct{
	web.Controller
}

func(c*CorreoController)EnciarCorreo(){
	body:= c.Ctx.Input.RequestBody
	if len(body)==0{
		c.Data["json"]=map[string]interface{}{
			"OK":false,
			"Mensaje":"El cuerpo de la solicitud esta vacio",
		}
		c.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		c.ServeJSON()
		return
	}

	var solicitud models.SolicitudEnvioCorreo

	err:= json.Unmarshal(body,&solicitud)
	if err!=nil{
		c.Data["json"]=map[string]interface{}{
			"OK":false,
			"Mensaje":"El cuerpo de la solicitud no tiene formato valido",
			"error":err.Error(),
		}
		c.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		c.ServeJSON()
		return
	}

	solicitud.Asunto=strings.TrimSpace(solicitud.Asunto)
	if solicitud.Asunto==""{
		c.Data["json"]=map[string]interface{}{
			"OK":false,
			"Mensaje":"El Asunto es obligatorio",
			
		}
		c.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		c.ServeJSON()
		return
	}
	if strings.TrimSpace(solicitud.Cuerpo)==""{
		c.Data["json"]=map[string]interface{}{
			"OK":false,
			"Mensaje":"El Cuerpo del correo  es obligatorio",
			
		}
		c.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		c.ServeJSON()
		return
	}
	if len(solicitud.Destinatario)==0{
		c.Data["json"]=map[string]interface{}{
			"OK":false,
			"Mensaje":"Debe de existir un destinatario",
			
		}
		c.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		c.ServeJSON()
		return
	}

	gmailService, err:=services.NewGmailService()

	if err!=nil{
		c.Data["json"]=map[string]interface{}{
			"OK":false,
			"Mensaje":"No fue posible inicializar el servicio de gmail ",
			"error":err.Error(),
			
		}
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.ServeJSON()
		return
	}

	var enviados []map[string]interface{}
	var errores []map[string]interface{}

	for _, destinatario:= range solicitud.Destinatario{
		correo:= strings.TrimSpace(strings.ToLower(destinatario.Correo),)
		nombre:=strings.TrimSpace(destinatario.Nombre)
	
		if correo==""{
		errores =append(errores, map[string]interface{}{
			"correo":correo,
			"nombre":nombre,
			"Mensaje":"el correo no esta",
		})
		continue
	}
	_, err := mail.ParseAddress(correo)

	if err!=nil{
		errores =append(errores, map[string]interface{}{
			"correo":correo,
			"nombre":nombre,
			"Mensaje":"el correo no es valido",
		})
		continue
	}
	cuerpoPersonalizado:= solicitud.Cuerpo

	cuerpoPersonalizado=strings.ReplaceAll(
		cuerpoPersonalizado,
		"{{nombre}}",
		nombre,
	)

	asuntoPersonalizado := strings.ReplaceAll(
		solicitud.Asunto,
		"{{nombre}}",
		nombre,
	)
	mensajeID, err := gmailService.SendEmail(
			correo, asuntoPersonalizado, cuerpoPersonalizado,
		)

		if err != nil {
			errores = append(errores, map[string]interface{}{
				"Corrreo" : correo,
				"Nombre" : nombre,
				"Mensaje" : "No fué posible enviar el correo.",
				"Error" : err.Error(),
			})
			continue
		}

		enviados = append(enviados, map[string]interface{}{
			"Correo" : correo,
			"Nombre" : nombre,
			"MensajeID" : mensajeID,
		})
	}
	
	if len(enviados) == 0 {
		c.Data["json"] = map[string]interface{} {
			"OK" : false,
			"Mensaje" : "No fué posible enviar ningún correo.",
			"Enviado" : len(enviados),
			"Errores" : errores,
		}
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.ServeJSON()
		return
	}
	c.Data["json"] = map[string]interface{} {
		"OK" : true,
		"Mensaje" : "Proceso de envíos de correos finalizado.",
		"Enviados" : enviados,
		"Errores" : errores,
	}

	c.Ctx.ResponseWriter.WriteHeader(http.StatusOK)
	c.ServeJSON()
}