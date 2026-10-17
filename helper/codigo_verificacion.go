package helper

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"
	"sync"
	"strings"
	
	"api_mid_yoplay/services"
)
// estructura para generar el codgio 
type CodigoVerificacion struct{
	Codigo string
	FechaExpiracion time.Time
	Proposito string
}
// variables globales para almacennar los codigos y para las solicitudes 

var (
	codigosTemporales = make(map[string]CodigoVerificacion)
	mutexCodigo  sync.Mutex
)

func GenerarCodigoVerificacion()(string, error){
	numero, err := rand.Int(rand.Reader, big.NewInt(900000))
	if err != nil{
		return "", err
	}

	codigo:= numero.Int64() + 100000
	
	return fmt.Sprintf("%06d", codigo), nil
}


// funcion para guardar el codgio de verificacion 
func GuardarCodigoVerificacion( correo, codigo, proposito string){
	correo =strings.ToLower(strings.TrimSpace(correo))
	clave:= correo +":"+ proposito

	mutexCodigo.Lock()
	defer mutexCodigo.Unlock()

	codigosTemporales[clave] = CodigoVerificacion{
		Codigo: codigo,
		FechaExpiracion: time.Now().Add(3*time.Minute),
		Proposito: proposito,
	}
}


func VerificarCodigo(correo, codigo, proposito string)bool{
	correo = strings.ToLower(strings.TrimSpace(correo))

	clave := correo +":" + proposito

	mutexCodigo.Lock()
	defer mutexCodigo.Unlock()

	codigoGuardado, existe := codigosTemporales[clave]

	if !existe{
		return false
	}
	if time.Now().After(codigoGuardado.FechaExpiracion){
		delete(codigosTemporales, clave)
		return false
	}
	return false
}

func EnviarCodigoVerificacion(correo, proposito string)error{
	codigo, err:= GenerarCodigoVerificacion()
	if err!=nil{
		return fmt.Errorf("Error al generar el codigo de verificacion: %v", err)
	}
	
	servicio, err := services.NewGmailService()
	if err!=nil{
		return fmt.Errorf("Error al crear el servicio de Gmail: %v", err)
	}

	asunto := "Codigo de verificacion YoPLAY"
	cuerpo:="Su codigo de verificacion es: " + codigo + "\nEste codigo es valido por 3 minutos."

	_,err= servicio.SendEmail(correo, asunto,cuerpo)
	if err!=nil{
		return fmt.Errorf("Error al enviar el correo de verificacion: %v", err)
	}
	// Guardar el código de verificación si la peticion fue exitosa

	GuardarCodigoVerificacion(correo, codigo, proposito)
	return nil
}