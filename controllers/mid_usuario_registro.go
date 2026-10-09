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
	IdDocumento    int    `json:"id"`
}

type Usuario struct {
	Id                string	`json:"id_usuario"` 
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
	Usuario:= Usuario{
		Nombre: solicitud.Nombre,
		Apellido: solicitud.Apellido,
		NumeroDocumento: solicitud.NumeroDocumento,
		Email: solicitud.Email,
		Telefono: solicitud.Telefono,
		FechaNacimiento: solicitud.FechaNacimiento,
		Terminos: solicitud.Terminos,
		IdDocumento: solicitud.TipoDocumento.IdDocumento,
	}
	
	usuarioCreado, err:= json.Marshal((Usuario))
	if err !=nil{
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"]= map[string]string{
			"error":"no fue posible convertir el json de usuario",
		}
		c.ServeJSON()
		return
	}
	respuestaUsuario, err:= http.Post("http://localhost:8080/v1/usuario","aplication/json",bytes.NewBuffer(usuarioCreado))
	if err!=nil{
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"]= map[string]string{
			"error":"No se pudo realizar el post de usuario",
		}
		c.ServeJSON()
		return
	}
	defer respuestaUsuario.Body.Close()
	if respuestaUsuario.StatusCode< http.StatusOK || respuestaUsuario.StatusCode>= http.StatusMultipleChoices{
		c.Ctx.Output.SetStatus(respuestaUsuario.StatusCode)
		c.Data["json"] =map[string]string{
			"error":"api CRUD de usuario no creo el usuario",
		}
		c.ServeJSON()
		return
	}
	cuerpoUsuarioCreado, err:= io.ReadAll(respuestaUsuario.Body)
	if err!=nil{
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"]=map[string]string{
			"error":"no fue posible leer la respuesta del usuario creado",
		}

		c.ServeJSON()
		return
	}

	
	
	fmt.Printf("Registro recivido: %+v\n",solicitud)

	c.Data["json"]= solicitud
	c.ServeJSON()

}
