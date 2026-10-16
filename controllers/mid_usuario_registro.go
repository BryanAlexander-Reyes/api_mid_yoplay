package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"golang.org/x/crypto/bcrypt"
	"time"
	

	beego "github.com/beego/beego/v2/server/web"
)

type Documento struct {
	IdDocumento    int    `json:"id"`
}

type Usuario struct {
	Id               int	`json:"Id"` 
	Nombre            string     `json:"Nombre"`
	Apellido          string     `json:"Apellido"`
	NumeroDocumento   string     `json:"NumeroDocumento"`
	Email             string     `json:"Email"`
	Telefono          string     `json:"Telefono"`
	Terminos bool `json:"Terminos"`
	FechaNacimiento   *time.Time  `json:"FechaNacimiento"`
	IdDocumento       *Documento `json:"IdDocumento"`
}
type Contrasena struct {
	Id       int  `json:"id"`
	HashContrasena    string    `json:"HashContrasena"`
	IdUsuario        *Usuario  `json:"IdUsuario"`
	
}

type  Registro struct{
	Nombre            string     `json:"Nombre"`
	Apellido          string     `json:"apellido"`
	NumeroDocumento   string     `json:"numero_documento"`
	Email             string     `json:"email"`
	Telefono          string     `json:"telefono"`
	Terminos bool `json:"terminos"`
	FechaNacimiento   string  `json:"fecha_nacimiento"`
	FechaRegistro     string  `json:"FechaRegistro"`
	TipoDocumento     *Documento    `json:"tipo_documento"`
	Contrasena        string    `json:"contrasena"`
}
type RegistroResponse struct {
    Mensaje   string `json:"mensaje"`
    IdUsuario int    `json:"id_usuario"`
    IdContrasena int `json:"id_contrasena"`
}
type UsuarioCRUDResponse struct {
    Succes bool    `json:"succes"`
    Status int     `json:"status"`
    Mensaje string `json:"Mensaje"`
    Data   Usuario `json:"Data"`
	ID		Contrasena `json:"id"`
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
	urlUsuarios, _ := beego.AppConfig.String("url_usuario")
	urlContrasena,_ := beego.AppConfig.String("url_contrasena")
	fmt.Println("URL de usuarios:", urlUsuarios)
	fmt.Println("URL de contrasenas:", urlContrasena)
	
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
	
	fechaNacimiento, err := time.Parse("2006-01-02", solicitud.FechaNacimiento)
	if err != nil {
    	c.Ctx.Output.SetStatus(http.StatusBadRequest)
    	c.Data["json"] = map[string]string{
        "mensaje": "La fecha de nacimiento no tiene un formato válido",
    }
    c.ServeJSON()
    return
}
	usuario:= Usuario{
		Nombre: solicitud.Nombre,
		Apellido: solicitud.Apellido,
		NumeroDocumento: solicitud.NumeroDocumento,
		Email: solicitud.Email,
		Telefono: solicitud.Telefono,
		FechaNacimiento: &fechaNacimiento,
		Terminos: solicitud.Terminos,
		IdDocumento: &Documento{
			IdDocumento: solicitud.TipoDocumento.IdDocumento,
	},
}
	
	UsuarioJSON, err:= json.Marshal(usuario)
	if err !=nil{
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"]= map[string]string{
			"error":"no fue posible convertir el json de usuario",
		}
		c.ServeJSON()
		return
	}
	fmt.Println("Respuesta del POST:", string(UsuarioJSON))
	respuestaUsuario, err:= http.Post(urlUsuarios,"application/json",bytes.NewBuffer(UsuarioJSON))

	if err!=nil{
		fmt.Println("Error al enviar POST al CRUD:", err)
    fmt.Println("URL del CRUD:", urlUsuarios)
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


	
fmt.Println("Respuesta completa del CRUD:", string(cuerpoUsuarioCreado))

var respuestaCRUD UsuarioCRUDResponse

err = json.Unmarshal(cuerpoUsuarioCreado, &respuestaCRUD)
if err != nil {
    fmt.Println("Error al interpretar respuesta:", err)
    fmt.Println("Respuesta recibida:", string(cuerpoUsuarioCreado))
    return
}

if !respuestaCRUD.Succes {
    fmt.Println("El CRUD rechazó la creación del usuario.")
    fmt.Println("Mensaje del CRUD:", respuestaCRUD.Mensaje)

    c.Ctx.Output.SetStatus(http.StatusBadGateway)
    c.Data["json"] = map[string]interface{}{
        "mensaje": "El CRUD no pudo crear el usuario",
        "detalle": respuestaCRUD.Mensaje,
    }
    c.ServeJSON()
    return
}

UsuarioCreado := &respuestaCRUD.Data
idUsuario := UsuarioCreado.Id

fmt.Println("ID del usuario creado:", idUsuario)

if idUsuario <= 0 {
    fmt.Println("Error: el CRUD no devolvió un ID válido.")
    c.Ctx.Output.SetStatus(http.StatusBadGateway)
    c.Data["json"] = map[string]interface{}{
        "mensaje": "El CRUD no devolvió un ID válido",
    }
    c.ServeJSON()
    return
}

	
	hash, err:= bcrypt.GenerateFromPassword(
		[]byte(solicitud.Contrasena),
		bcrypt.DefaultCost,
	)
	if err!=nil{
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"]=map[string]string{
			"ERROR":"No fue posible proteger la contraseña",
		}	
		c.ServeJSON()
		return
	}
	hashTexto := string(hash)

	contrasenaRequest := Contrasena{
		IdUsuario: &Usuario{
			Id: idUsuario,
		},
		HashContrasena: hashTexto,
	}

	jsoncontrasena, err := json.Marshal(contrasenaRequest)
	if err!=nil{
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"]= map[string]string{
			"Error": "No se pudo generar el json de la contrasena",
		}
		c.ServeJSON()
		return
	}
	respuestaContrasena,err:= http.Post(urlContrasena,"application/json",bytes.NewBuffer(jsoncontrasena),)
	
	if err!=nil{
		
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"]= map[string]string{
			"Error": "El usuario fue creado pero no fue posible no fue posible crear la contrasena",
		}
		c.ServeJSON()
		return
	}
	defer respuestaContrasena.Body.Close()
	if respuestaContrasena.StatusCode<http.StatusOK|| respuestaContrasena.StatusCode>=http.StatusMultipleChoices{
		fmt.Println("ID del usuario que se eliminará:", UsuarioCreado.Id)
		urlEliminarUsuario := urlUsuarios +"/"+ strconv.Itoa(UsuarioCreado.Id)

		request, err :=http.NewRequest(http.MethodDelete,urlEliminarUsuario,nil)
		if err!=nil{
			return  
		}
		client:= &http.Client{}
		response, err:= client.Do(request)

		if err!=nil{
			fmt.Println("falló el metodo delete")
			return 
		}
		defer response.Body.Close()

		if response.StatusCode< http.StatusOK || response.StatusCode>= http.StatusMultipleChoices{
			c.Ctx.Output.SetStatus(response.StatusCode)
		c.Data["json"]=map[string]string{
			"Error": "El usuario no se pudo eliminar",
		}
		c.ServeJSON()
		return
		}
		c.Ctx.Output.SetStatus(respuestaContrasena.StatusCode)
		c.Data["json"]=map[string]string{
			"Error": "Fallo el cracion de la contraseña, se realizo el rollback ",
		}
		c.ServeJSON()
		return
	}
	cuerpoContrasenaCreado, err := io.ReadAll(respuestaContrasena.Body)
	if err!=nil{
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"]=map[string]string{
			"Error": "El usuario fue creado pero no fue posible leer el rol",
		}
		c.ServeJSON()
		return
	}


	var contrasenaCreada Contrasena 

	err= json.Unmarshal(cuerpoContrasenaCreado, &contrasenaCreada)

	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["JSON"] = map[string]string {
			"Error" : "El usuario fué creado, pero no fué posible procesar la respuesta.",
		}
		c.ServeJSON()
		return
	}
	

	resultado :=RegistroResponse{
		Mensaje: "Usuario y contraseña creados exitosamente.",
		IdUsuario: UsuarioCreado.Id,
		IdContrasena: contrasenaCreada.Id,
	}
	
	
	fmt.Printf("Registro recivido: %+v\n",resultado)

	c.Data["json"]=resultado
	c.ServeJSON()

}