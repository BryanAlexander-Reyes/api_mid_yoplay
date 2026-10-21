package services

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)


type GmailService struct{
	service *gmail.Service
}

// Crea y configura un cliente autenticado para utilizar Gmail


func NewGmailService()(*GmailService,error){

	ctx:= context.Background()

	// definicion de credenciales

	credentialsPath:="credentials/credentials.json"

	// leer archivo credenciales

	credentials, err:= os.ReadFile(credentialsPath)

	if err != nil{
		return nil, fmt.Errorf(
			"No fue posible leer las credenciales.json: %w",err,
		)
	}

	// crear configuracion de oauth utilizando credenciales


	config ,err:= google.ConfigFromJSON(
		credentials,
		gmail.GmailSendScope,
	)

	if err != nil{
		return nil, fmt.Errorf(
			"No fue posible configurar OAuth 2.0: %w",err,
		)
	}

	// obtener el cliente http autenticado

	client, err := getGmailClient(config)

	if err != nil{
		return nil, err
	}

	// crear el servicio utilizando el cliente http

	service, err := gmail.NewService(
		ctx,
		option.WithHTTPClient(client),
	)

	if err != nil{
		return nil, fmt.Errorf(
			"no fue posible crear el servicio Gmail: %w", err,
		)
	}

	return &GmailService{
		service: service,
	},nil

}


// obtener un client http autenticado mediante token.json

func getGmailClient(config *oauth2.Config)(*http.Client,error){
	// Definicion path del token 

	tokenPath:= "token/token.json"


	token, err:= tokenFromFile(tokenPath)
	if err!= nil{
		token, err= getTokenFromWeb(config)
		if err!=nil{
			return nil, err
		}
		
		err = saveToken(tokenPath, token)
		if err!=nil{
			return nil, err
		}
	}
	return config.Client(
		context.Background(),
		token,
	),nil

}


// funcion para trear el token desde la web

func getTokenFromWeb(config *oauth2.Config)(*oauth2.Token, error){
	// generar el url donde se solicita el token

	authURL := config.AuthCodeURL(
		"state-token",
		oauth2.AccessTypeOffline,
	)

	fmt.Println("====================")
	fmt.Println("Autorizacion de GMAIL.")
	fmt.Println("====================")
	fmt.Println("Abra la siguiente URL en su navegador")
	fmt.Println()
	fmt.Println(authURL)
	fmt.Println()
	fmt.Println("======================================")

	fmt.Print("Ingrese el codigo de autorizacion")

	var autotizationCode string
	if _, err:= fmt.Scan(&autotizationCode); err !=nil{
		return nil, fmt.Errorf(
			"No fue posible leer el codigo: %w", err,
		)
	}

	token, err := config.Exchange(
		context.Background(),
		autotizationCode,
	)
	if err!=nil{
		return nil, fmt.Errorf(
			"No fue posible obtener el token AUTH: %w",err,
		)
	}
	return token, nil
	
}


// funcion de lectura de token desde el json

func tokenFromFile(path string)(*oauth2.Token, error){
	file , err:= os.Open(path)

	if err!=nil{
		return nil, err
	}

	defer file.Close()

	token := &oauth2.Token{}
	err= json.NewDecoder(file).Decode(token)

	if err!=nil{
		return nil, fmt.Errorf(
			"No fue posible leer el token.json: %w",err,
		)
	}
	return token, nil
}


// metodo para guardar el token 

func saveToken(path string, Token*oauth2.Token) error{
	err:= os.MkdirAll("token",0700)

	if err!=nil{
		return fmt.Errorf(
			"No fue posible crear la carpeta token:%w",err,
		)
	}
	
	file, err:= os.OpenFile(
		path,
		os.O_RDWR | os.O_CREATE | os.O_TRUNC,
		0600,
	)

	if err!=nil{
		return fmt.Errorf(
			"No fue posible crear token.json%w",err,
		)
	}

	defer file.Close()

	err= json.NewEncoder(file).Encode(Token)

	if err!=nil{
		return  fmt.Errorf(
			"No fue posible guardar el token.json: %w",err,
		)
	}
	return nil



}


// sendEmail metodo que envia correo 

// parametros
// destinatario: correo
// asunto
// cuerpo del correo
// retorna:
// id del mensaje generado en gmail
// error en caso de que falle

func (g *GmailService) SendEmail(
	destinatario string,
	asunto string,
	cuerpoHTML string,

	)(string, error){
		if strings.TrimSpace(destinatario)==""{
			return "", fmt.Errorf(
				"El destinatario es obligatorio",
			)
		}
		if strings.TrimSpace(asunto)==""{
			return "", fmt.Errorf(
				"El asunto es obligatorio",
			)
		}
		if strings.TrimSpace(cuerpoHTML)==""{
			return "",fmt.Errorf("el cuerpo html es obligatorio")
		}

		// crear el mensaje mime
		// estructura de mime si o si 
		

		menesajeMIME:= strings.Join([]string{
			"To: " +destinatario,
			"Subject: " +asunto,
			"Content-Type: text/html; charset=\"UTF-8\"",
			"MIME-Version: 1.0",
			"Content-Transfer-Encoding: 8bit",
			"",
			cuerpoHTML,
		},"\r\n")

		mensajeCodificado:= base64.RawURLEncoding.EncodeToString([]byte(menesajeMIME))

		mensaje:= &gmail.Message{
			Raw: mensajeCodificado,
		}

		respuesta, err:= g.service.Users.Messages.Send(
			"me",
			mensaje,
		).Do()

		if err!= nil{
			return "", fmt.Errorf(
				"No fue posible enviar el correo mediante GMAIL: %w",err,
			)
		}

		if respuesta ==nil|| respuesta.Id==""{
			return"", fmt.Errorf(
				"Gmail API no retorno el id del mensaje enviado",
			)
		}
		return respuesta.Id, nil
}

func (g*GmailService)SendEmailConAdjunto(
	destinatario string,
	asunto string,
	cuerpo string,
	nombreArchivo string,
	archivo []byte,

)(string, error){
	var mensaje bytes.Buffer

	boundary:= "FacturaBoundary1234567"

	fmt.Fprintf(
		&mensaje,
		"From: me\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: multipart/mixed; boundary=\"%s\"\r\n"+
		"\r\n",
		destinatario,
		asunto,
		boundary,
	)

	// 
	fmt.Fprintf(
		&mensaje,
		"--%s\r\n"+
			"Content-Type: text/plain; charset=\"UTF-8\"\r\n"+
			"Content-Transfer-Encoding: 8bit\r\n"+
			"\r\n"+
			"%s\r\n",
		boundary,
		cuerpo,
	)
	// parte del pdf
	fmt.Fprintf(
		&mensaje,
		"--%s\r\n"+
			"Content-Type: application/pdf; name=\"%s\"\r\n"+
			"Content-Disposition: attachment; filename=\"%s\"\r\n"+
			"Content-Transfer-Encoding: base64\r\n"+
			"\r\n",
		boundary,
		nombreArchivo,
		nombreArchivo,
	)
	pdfBase64:= base64.StdEncoding.EncodeToString(archivo)

	// Escribir base64 en lineas de 76 caracteres
	for i:=0; i<len(pdfBase64); i+=76{
		fin:=i+76
		if fin>len(pdfBase64){
			fin=len(pdfBase64)
		}
		mensaje.WriteString(pdfBase64[i:fin])
		mensaje.WriteString("\r\n")
	}
	fmt.Fprintf(&mensaje,"--%s--\r\n",boundary)
	raw:= base64.URLEncoding.EncodeToString(mensaje.Bytes())
	msg:= &gmail.Message{Raw:raw}

	resultado, err:=g.service.Users.Messages.Send("me",msg).Do()
	if err!=nil{
		return "", err
	}
	return resultado.Id, nil
}