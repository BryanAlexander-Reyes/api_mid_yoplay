package models

type premios struct{
	Puesto string `json:"puesto"`
	Premio string `json:"premio"`
}

type CrearTorneo struct{
	Lugar string `json:"lugar"`
	Premiaciones []premios `json:"premiaciones"`
	Tipo_Deporte []TipoDeporte `json:"tipo_deporte"`
	Divicion []TipoDivision `json:"divicion"`
	CantidadEquipos int `json:"cantidadequipos"`
	Objetivo string `json:"objetivo"`
	FechaInicio string `json:"fechaInicio"`
	FechaFin string `json:"fechaFin"`
	ImagenUrl []imagenUrl `json:"imagenurl"`
}

type TipoDeporte struct{
	IDDeporte string `json:"id_tipo_deporte"`
	Nombre_Deporte string `json:"nombre_deporte"`

}

type TipoDivision struct{
	IDDivision string `json:"id_tipo_distribucion"`
	Division string`json:"Nombre_distribucion"`
}

type imagenUrl struct {
	IDImange string `json:"Idimagen"`
	Imagen string `json:"imagen"`
}