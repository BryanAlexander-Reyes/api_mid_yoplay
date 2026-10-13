package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"golang.org/x/crypto/bcrypt"

	beego "github.com/beego/beego/v2/server/web"
)

type Documento struct {
	IdDocumento    int    `json:"id"`
}

type Usuario struct {
	Id               int	`json:"id_usuario"` 
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
	Id       int  `json:"id_contrasena"`
	HashContrasena    string    `json:"hash_contrasena"`
	IdUsuario         int  `json:"id_usuario"`
	
}

type  Registro struct{
	Nombre            string     `json:"nombre"`
	Apellido          string     `json:"apellido"`
	NumeroDocumento   string     `json:"numero_documento"`
	Email             string     `json:"email"`
	Telefono          string     `json:"telefono"`
	Terminos bool `json:"terminos"`
	FechaNacimiento   string  `json:"fecha_nacimiento"`
	FechaRegistro     string  `json:"fecha_registro"`
	TipoDocumento     *Documento    `json:"tipo_documento"`
	Contrasena        string    `json:"contrasena"`
}
type RegistroResponse struct {
    Mensaje   string `json:"mensaje"`
    IdUsuario int    `json:"id_usuario"`
    IdContrasena int `json:"id_contrasena"`
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
	usuario:= Usuario{
		Nombre: solicitud.Nombre,
		Apellido: solicitud.Apellido,
		NumeroDocumento: solicitud.NumeroDocumento,
		Email: solicitud.Email,
		Telefono: solicitud.Telefono,
		FechaNacimiento: solicitud.FechaNacimiento,
		Terminos: solicitud.Terminos,
		IdDocumento: solicitud.TipoDocumento.IdDocumento,
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
	respuestaUsuario, err:= http.Post(urlUsuarios,"aplication/json",bytes.NewBuffer(UsuarioJSON))
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
	var UsuarioCreado Usuario
	err=json.Unmarshal(cuerpoUsuarioCreado, &UsuarioCreado)
	
	if err!=nil{
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"]=map[string]string{
			"Error":"No fue posible procesar la respuesta del crud",
		}
		c.ServeJSON()
		return
	}

	idUsuario:= UsuarioCreado.Id
	
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

	contrasenaRequest := Contrasena{
		IdUsuario: idUsuario,
		HashContrasena: string(hash),
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
		Mensaje: "Usuario y rol creados exitosamente.",
		IdUsuario: contrasenaCreada.IdUsuario,
		IdContrasena: contrasenaCreada.Id,
	}
	
	
	fmt.Printf("Registro recivido: %+v\n",resultado)

	c.Data["json"]=resultado
	c.ServeJSON()

}
