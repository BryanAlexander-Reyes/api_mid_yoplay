package helper

import (
	"errors"
	"fmt"
	"net/http"
	"io"
	"net/url"

	beego "github.com/beego/beego/v2/server/web"
)

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

	

	return false, nil
}