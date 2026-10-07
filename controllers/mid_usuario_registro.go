package controllers

import (
	beego "github.com/beego/beego/v2/server/web"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

type Documento struct {
	TipoDocumento     string    `json:"(tipo_documento)"`
}

type Usuario struct {
	Nombre            string     `json:"(nombre)"`
	Apellido          string     `json:"(apellido)"`
	NumeroDocumento   string     `json:"(numero_documento)"`
	Email             string     `json:"(email)"`
	Telefono          string     `json:"(telefono)"`
	FechaNacimiento   string  `json:"(fecha_nacimiento);;type(timestamp without time zone)"`
	IdDocumento       *Documento `json:"(id_documento);rel(fk);on_delete(cascade)"`
	FechaRegistro     string  `json:"(fecha_registro);;type(timestamp without time zone);null"`
}
type Contrasena struct {
	Contrasena        string    `json:"(contrasena)"`
	HashContrasena    string    `json:"(hash_contrasena)"`
	IdUsuario         *Usuario  `json:"(id_usuario);rel(fk); on_delete(cascade)"`
	
}

type  FormularioRegistro struct{
	
}

// MidController operations for Mid
type MidController struct {
	beego.Controller
}
