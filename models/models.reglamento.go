package models


type Reglamento struct {
	Estandar string `json:"estandar"`
	Regla []reglas `json:"id_reglas"`
	Idioma string `json:"idioma"`
}

type reglas struct {
	Numero_Regla string `json:"numero_reglas"`
	Descripcion_Reglas string `json:"descripcion"`
}