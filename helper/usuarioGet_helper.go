package helper

import (
	"errors"
	"fmt"
	"net/http"
	"io"
	"net/url"
	"encoding/json"
	"strconv"

	beego "github.com/beego/beego/v2/server/web"
)

type responseUsuario struct {
	Data []IdUsuario `json:"Data"`}

type IdUsuario struct {
	ID    int    `json:"Id"`
	Email string `json:"Email"`
}
type responseContrasena struct {
    Data []IdContrasena `json:"Data"`
}

type IdContrasena struct {
    ID int `json:"Id"`
}

func BuscarUsuarioCorreo(correo string) (bool, error){

	urlUsuarios, _ := beego.AppConfig.String("url_usuario")
	urlconsulta := urlUsuarios + "?query=" + url.QueryEscape("Email:"+correo)

	if correo == ""{
		return  false, errors.New("el correo es obligatorio")
	}

	peticion, err := http.NewRequest(http.MethodGet,urlconsulta,nil,)
	if err!=nil{
		return false, err
	}

	cliente := &http.Client{}
	respuesta, err:= cliente.Do(peticion)
	if err!=nil{
		return false, err
	}
	defer respuesta.Body.Close()
	if respuesta.StatusCode != http.StatusOK{
		return false, fmt.Errorf(
			"el CRUD respondio con el estado:%d",
			respuesta.StatusCode,
		)
	}
	cuerpo, err:=io.ReadAll(respuesta.Body)
	if err!=nil{
		return false, err
	}
	fmt.Println("Respuesta Crud", string(cuerpo))

	var resultado responseUsuario
	err= json.Unmarshal(cuerpo,&resultado)
	if err!=nil{
		return false, fmt.Errorf("no se pudo convertir la peticion")
	}
	if len(resultado.Data)>0{
		return true, nil
	}

	return false, nil
}


// funcion para obtener el id del usuario a partir del correo
func ObtenerIDUsuarioCorreo(correo string)(int, error){
	urlUsuarios, _ := beego.AppConfig.String("url_usuario")
	urlconsulta := urlUsuarios + "?query=" + url.QueryEscape("Email:"+correo)

	peticion, err:= http.NewRequest(http.MethodGet, urlconsulta,nil)
	if err!=nil{
		return 0, err
	}
	cliente:= &http.Client{}

	respuesta, err:= cliente.Do(peticion)
	if err!=nil{
		return 0, err
	}
	defer respuesta.Body.Close()
	if respuesta.StatusCode != http.StatusOK{
		return 0, fmt.Errorf("El crud respondio con el estado: %d", respuesta.StatusCode,)
	}

	cuerpo, err:=io.ReadAll(respuesta.Body)
	if err!=nil{
		return 0, err
	}
	fmt.Printf("Respuesta Crud: %s", string(cuerpo))


	var resultado responseUsuario
	err = json.Unmarshal(cuerpo,&resultado)
	if err!=nil{
		return 0, fmt.Errorf("no se pudo convertir la peticion ")
	}

	if len(resultado.Data)==0{
		return 0, errors.New("no se encontró el usuario con ese correo")
	}

	return resultado.Data[0].ID, nil

}

func ObtenerIdContrasena(IdUSUARIO int)(int, error){
	urlContrasena,_ := beego.AppConfig.String("url_contrasena")
	urlConsulta := urlContrasena +"?query=" + url.QueryEscape("IdUsuario.id: "+ strconv.Itoa(IdUSUARIO),)

	peticion, err:= http.NewRequest(http.MethodGet, urlConsulta, nil)
	if err!=nil{
		return 0, errors.New("No se pudo leer la peticion")
	}

	cliente:= &http.Client{}

	respuesta, err:= cliente.Do(peticion)
	if err!=nil{
		return 0, err
	}
	defer respuesta.Body.Close()
	if respuesta.StatusCode != http.StatusOK{
		return 0, fmt.Errorf("El crud respondio con el estado: %d", respuesta.StatusCode,)
	}

	cuerpo, err:=io.ReadAll(respuesta.Body)
	if err!=nil{
		return 0, err
	}
	fmt.Printf("Respuesta Crud: %s", string(cuerpo))


	var resultado responseContrasena
	err = json.Unmarshal(cuerpo,&resultado)
	if err!=nil{
		return 0, fmt.Errorf("no se pudo convertir la peticion ")
	}

	if len(resultado.Data)==0{
		return 0, errors.New("no se encontro una contraseña")
	}

	return resultado.Data[0].ID, nil
}