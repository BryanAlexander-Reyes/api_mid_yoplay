package helper

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/beego/beego/v2"
)

func BuscarUsuarioCorreo(correo string) (bool, error){

	urlUsuarios, _ := beego.AppConfig.String("url_usuario")
	urlconsulta := urlUsuarios + "?email=" + correo

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


	return false, nil
}