package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	beego "github.com/beego/beego/v2/server/web"
)

type Documento struct {
	Id     string    `json:"id"`
}

type Usuario struct {
	Nombre            string     `json:"nombre"`
	Apellido          string     `json:"apellido"`
	NumeroDocumento   string     `json:"numero_documento"`
	Email             string     `json:"email"`
	Telefono          string     `json:"telefono"`
	Terminos bool `json:"terminos"`
	FechaNacimiento   string  `json:"fecha_nacimiento"`
	IdDocumento       int `json:"id_documento"`
	FechaRegistro     string  `json:"fecha_registro"`
}
type Contrasena struct {
	Contrasena        string    `json:"contrasena"`
	HashContrasena    string    `json:"hash_contrasena"`
	IdUsuario         int  `json:"id_usuario"`
	
}

type  Registro struct{
	Nombre            string     `json:"nombre"`
	Apellido          string     `json:"apellido"`
	NumeroDocumento   string     `json:"numero_documento"`
	Email             string     `json:"email"`
	Telefono          string     `json:"telefono"`
	Terminos bool `json:"terminos;default:false"`
	FechaNacimiento   string  `json:"fecha_nacimiento"`
	FechaRegistro     string  `json:"fecha_registro"`
	TipoDocumento     *Documento    `json:"tipo_documento"`
	Contrasena        string    `json:"contrasena"`
}



// MidController operations for Mid
type RegistroController struct {
	beego.Controller
}


// @Title Resgistro usuario 
// @Description recibe json completo
// @Param body body models.ResgistroUsuarioRequest true "Informacion de usuario y rol"
// @Success 201 {object} models.ResgistroUsuariosResponse
// @Failure 400 {object} map[string]String
// @failure 500 {object} map[string]String
// @router /Resgistro [post]

func (c* RegistroController) CrearUsuario(){

	var solicitud Registro
	err := json.Unmarshal(c.Ctx.Input.RequestBody,&solicitud)
	if err!=nil{
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] =map [string]string{
			"error": "El JSON enviado no contiene formulario valido",
		}
		c.ServeJSON()
		return
	}
	
	if solicitud.Nombre ==""{
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"]= map[string]string{
			"error":"El campo nombre es obligatorio",
		}
		c.ServeJSON()
		return
	}
	if solicitud.Apellido ==""{
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"]= map[string]string{
			"error":"El campo apellido es obligatorio",
		}
		c.ServeJSON()
		return
	}
	if solicitud.Email ==""{
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"]= map[string]string{
			"error":"El campo email es obligatorio",
		}
		c.ServeJSON()
		return
	}
	if solicitud.NumeroDocumento ==""{
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"]= map[string]string{
			"error":"El campo numero de documento es obligatorio",
		}
		c.ServeJSON()
		return
	}
	if solicitud.Telefono ==""{
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"]= map[string]string{
			"error":"El campo telefono es obligatorio",
		}
		c.ServeJSON()
		return
	}
	if solicitud.FechaNacimiento ==""{
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"]= map[string]string{
			"error":"El campo fecha nacimiento es obligatorio",
		}
		c.ServeJSON()
		return
	}
	
	if solicitud.Contrasena== ""{
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"]= map[string]string{
			"error":"El campo contraseña es obligatorio",
		}
		c.ServeJSON()
		return
	}
	if !solicitud.Terminos {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"]= map[string]string{
			"error":"El campo Terminos es obligatorio",
		}
		c.ServeJSON()
		return
	
	}
	documento:= Documento{
		TipoDocumento: solicitud.TipoDocumento,
	}
	datosDocumento, err:= json.Marshal(documento)
	if err !=nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"]= map[string]string{
		"error":"no fue posible de convertir el json de documento",
		}
		c.ServeJSON()
		return
	}
	respuestaDocumento, err:= http.Post("http://localhost:8080/v1/documento","aplication/json",bytes.NewBuffer(datosDocumento),)

	if err!=nil{
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"]= map[string]string{
		"error":"no fue posible hacer post del documento",
		}
		c.ServeJSON()
		return
	}
	defer respuestaDocumento.Body.Close()
	if respuestaDocumento.StatusCode<http.StatusOK || respuestaDocumento.StatusCode>=http.StatusMultipleChoices{
		c.Ctx.Output.SetStatus(respuestaDocumento.StatusCode)
		c.Data["json"]= map[string]string{
			"error":"api Crud de documento no creo el documento",
		}
		c.ServeJSON()
		return
	}

	cuerpoDocumento, err:=io.ReadAll(respuestaDocumento.Body)
	if err!=nil{
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"]= map[string]string{
			"error":"No fue posible leer la respuesta del documento creado",
		}
		c.ServeJSON()
		return
	}

	var DocumentoCreado Documento
	err= json.Unmarshal(cuerpoDocumento,&DocumentoCreado)

	if err!=nil{
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"]=map[string]string{
			"Error":"No fue posible procesar la respuesta",
		}
		c.ServeJSON()
		return
	}
	idDocumento, err:=strconv.Atoi(DocumentoCreado.TipoDocumento)
	
	fmt.Printf("Registro recivido: %+v\n",solicitud)

	c.Data["json"]= solicitud
	c.ServeJSON()

}
