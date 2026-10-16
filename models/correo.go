package models

type DestinatarioCorreo struct{
	Correo string `json:"correo"`
	Nombre string `json:"nombre"`
}

type SolicitudEnvioCorreo struct{
	Destinatario []DestinatarioCorreo `json:"destinatario"`
	Asunto string `json:"asunto"`
	Cuerpo string `json:"cuerpo"`

}