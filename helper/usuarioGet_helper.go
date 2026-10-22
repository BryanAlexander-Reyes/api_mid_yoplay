package helper

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"golang.org/x/crypto/bcrypt"

	beego "github.com/beego/beego/v2/server/web"
)

type responseUsuario struct {
	Data []IdUsuario `json:"Data"`
}

type IdUsuario struct {
	ID    int    `json:"Id"`
	Email string `json:"Email"`
}
type responseContrasena struct {
	Data []IdContrasena `json:"Data"`
}

type IdContrasena struct {
	ID int `json:"Id"`
	HashContrasena string `json:"HashContrasena"`
}

type responseContrasenaID struct{
	Data *IdContrasena `json:"Data"`
}

func BuscarUsuarioCorreo(correo string) (bool, error) {

	urlUsuarios, _ := beego.AppConfig.String("url_usuario")
	urlconsulta := urlUsuarios + "?query=" + url.QueryEscape("Email:"+correo)

	if correo == "" {
		return false, errors.New("el correo es obligatorio")
	}

	peticion, err := http.NewRequest(http.MethodGet, urlconsulta, nil)
	if err != nil {
		return false, err
	}

	cliente := &http.Client{}
	respuesta, err := cliente.Do(peticion)
	if err != nil {
		return false, err
	}
	defer respuesta.Body.Close()
	if respuesta.StatusCode != http.StatusOK {
		return false, fmt.Errorf(
			"el CRUD respondio con el estado:%d",
			respuesta.StatusCode,
		)
	}
	cuerpo, err := io.ReadAll(respuesta.Body)
	if err != nil {
		return false, err
	}

	var resultado responseUsuario
	err = json.Unmarshal(cuerpo, &resultado)
	if err != nil {
		return false, fmt.Errorf("no se pudo convertir la peticion")
	}
	if len(resultado.Data) > 0 {
		return true, nil
	}

	return false, nil
}

// funcion para obtener el id del usuario a partir del correo
func ObtenerIDUsuarioCorreo(correo string) (int, error) {
	urlUsuarios, _ := beego.AppConfig.String("url_usuario")
	urlconsulta := urlUsuarios + "?query=" + url.QueryEscape("Email:"+correo)

	peticion, err := http.NewRequest(http.MethodGet, urlconsulta, nil)
	if err != nil {
		return 0, err
	}
	cliente := &http.Client{}

	respuesta, err := cliente.Do(peticion)
	if err != nil {
		return 0, err
	}
	defer respuesta.Body.Close()
	if respuesta.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("El crud respondio con el estado: %d", respuesta.StatusCode)
	}

	cuerpo, err := io.ReadAll(respuesta.Body)
	if err != nil {
		return 0, err
	}
	fmt.Printf("Respuesta Crud: %s", string(cuerpo))

	var resultado responseUsuario
	err = json.Unmarshal(cuerpo, &resultado)
	if err != nil {
		return 0, fmt.Errorf("no se pudo convertir la peticion ")
	}

	if len(resultado.Data) == 0 {
		return 0, errors.New("no se encontró el usuario con ese correo")
	}

	return resultado.Data[0].ID, nil

}

func ObtenerIdContrasena(IdUSUARIO int) (int, error) {
	urlContrasena, _ := beego.AppConfig.String("url_contrasena")
	urlConsulta := urlContrasena + "?query=" + url.QueryEscape("IdUsuario.id: "+strconv.Itoa(IdUSUARIO))

	peticion, err := http.NewRequest(http.MethodGet, urlConsulta, nil)
	if err != nil {
		return 0, errors.New("No se pudo leer la peticion")
	}

	cliente := &http.Client{}

	respuesta, err := cliente.Do(peticion)
	if err != nil {
		return 0, err
	}
	defer respuesta.Body.Close()
	if respuesta.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("El crud respondio con el estado: %d", respuesta.StatusCode)
	}

	cuerpo, err := io.ReadAll(respuesta.Body)
	if err != nil {
		return 0, err
	}
	fmt.Printf("Respuesta Crud: %s", string(cuerpo))

	var resultado responseContrasena
	err = json.Unmarshal(cuerpo, &resultado)
	if err != nil {
		return 0, fmt.Errorf("no se pudo convertir la peticion ")
	}

	if len(resultado.Data) == 0 {
		return 0, errors.New("no se encontro una contraseña")
	}

	return resultado.Data[0].ID, nil
}

func ObtenerHashContrasena(IdContrasena int)(string, error){
	urlContrasena, _ := beego.AppConfig.String("url_contrasena")
	urlConsulta := urlContrasena +"/"+strconv.Itoa(IdContrasena)

	peticion, err:=http.NewRequest(http.MethodGet,urlConsulta,nil)
	if err!=nil{
		return "", err
	}

	cliente := &http.Client{}
	respuesta, err:= cliente.Do(peticion)
	if err!=nil{
		return "", err
	}
	defer respuesta.Body.Close()

	if respuesta.StatusCode != http.StatusOK {
		return "", fmt.Errorf("El crud respondio con el estado: %d", respuesta.StatusCode)
	}
	cuerpo, err := io.ReadAll(respuesta.Body)
	if err != nil {
		return "", err
	}
	fmt.Printf("Respuesta Crud: %s", string(cuerpo))

	var resultado responseContrasenaID
	err = json.Unmarshal(cuerpo, &resultado)
	if err != nil {
		return "", fmt.Errorf("no se pudo convertir la peticion ")
	}
	if resultado.Data.HashContrasena==""{
		return "", errors.New("no se encontro el hash de la contraseña")
	}


	return resultado.Data.HashContrasena, nil
}
func GenerarHashContrasena(contrasena string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(
		[]byte(contrasena),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func ActualizarContrasena(IdContrasena int, hashTexto string) error {
	urlContrasena, err := beego.AppConfig.String("url_contrasena")
	if err != nil {
		return fmt.Errorf("no se podo obtener la url de contraseñas", err)
	}

	datos := map[string]interface{}{
		"HashContrasena": hashTexto,
	}
	cuerpoJSON, err := json.Marshal(datos)
	if err != nil {
		return fmt.Errorf("no se podo crear el json: %w", err)
	}
	fmt.Println("JSON enviado al CRUD:", string(cuerpoJSON))
	urlActualizacion := urlContrasena + "/" + strconv.Itoa(IdContrasena)

	peticion, err := http.NewRequest(http.MethodPut, urlActualizacion, bytes.NewBuffer(cuerpoJSON))
	if err != nil {
		return fmt.Errorf("no se pudo crear la peticion put: %w", err)
	}
	peticion.Header.Set("Content-Type", "application/json")

	cliente := &http.Client{}
	respuesta, err := cliente.Do(peticion)
	if err != nil {
		return fmt.Errorf("no se pudo actualizar la contraseña %w", err)
	}
	defer respuesta.Body.Close()
	cuerpoRespuesta, err := io.ReadAll(respuesta.Body)
	if err != nil {
		return fmt.Errorf("no se pudo leer la respuesta del CRUD: %w", err)
	}

	fmt.Println("Respuesta del CRUD:", string(cuerpoRespuesta))

	if respuesta.StatusCode < http.StatusOK || respuesta.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf(
			"el Crud respondio con el estado HTTP %d",
			respuesta.StatusCode,
		)

	}
	return nil
}
