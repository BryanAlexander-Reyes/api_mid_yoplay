package helper

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"html"
	"math/big"
	"os"
	"strings"
	"sync"
	"time"

	"api_mid_yoplay/services"
)

// estructura para generar el codgio
type CodigoVerificacion struct{
	Codigo string
	FechaExpiracion time.Time
	Proposito string
}
// estrutura para generar el token de recuperacion de contrasena
type TokenRecuperacion struct {
    Token            string
    IdUsuario        int
    FechaExpiracion  time.Time
}
// variables globales para almacennar los codigos y para las solicitudes 

var (
	codigosTemporales = make(map[string]CodigoVerificacion)
	mutexCodigo  sync.Mutex
)
var tokensRecuperacion = make(map[string]TokenRecuperacion)
var mutexToken sync.Mutex

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

	if codigoGuardado.Codigo == codigo{
		return false
	}
	delete(codigosTemporales, clave)
	return true
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
	// lee la plantilla
	plantilla, err:= os.ReadFile("plantillaCodigo/codigo_verificacion.html")
	if err!=nil{
		return fmt.Errorf("Error al leer la plantilla html:%w", err)
	} 
	// convertir la plantilla a texto para remplazar los marcadores

	cuerpo:= string(plantilla)

	cuerpo= strings.ReplaceAll(
		cuerpo,
		"{{proposito}}",
		html.EscapeString(proposito),
	)

	cuerpo= strings.ReplaceAll(
		cuerpo,
		"{{codigo}}",
		codigo,
	)
	_,err= servicio.SendEmail(correo, asunto,cuerpo)
	if err!=nil{
		return fmt.Errorf("Error al enviar el correo de verificacion: %v", err)
	}
	// Guardar el código de verificación si la peticion fue exitosa

	GuardarCodigoVerificacion(correo, codigo, proposito)
	return nil
}

func GenerarTokenRecuperacion()(string, error){
	tokenBytes := make([]byte, 32)
	_, err := rand.Read(tokenBytes)
	if err!=nil{
		return "", fmt.Errorf("Error al generar el token de recuperacion: %v", err)
	}
	token := hex.EncodeToString(tokenBytes)
	return token, nil
}

func GuardarTokenRecuperacion(token string, IdUsuario int){
	mutexToken.Lock()
	defer mutexToken.Unlock()


	tokensRecuperacion[token] = TokenRecuperacion{
		Token: token,
		IdUsuario: IdUsuario,
		FechaExpiracion: time.Now().Add(3*time.Minute),
	}
}

func ValidarTokenRecuperacion(token string) (int, error){
	mutexToken.Lock()
	defer mutexToken.Unlock()

	tokenGuardado, existe := tokensRecuperacion[token]
	if !existe{
		return 0, errors.New("el token no existe o no es valido")
	}
	if time.Now().After(tokenGuardado.FechaExpiracion){
		delete(tokensRecuperacion, token)
		return 0, errors.New("el token ha expirado")
	}
	delete(tokensRecuperacion, token)
	return tokenGuardado.IdUsuario, nil
}