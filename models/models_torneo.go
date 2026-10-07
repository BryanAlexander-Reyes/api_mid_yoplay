package models

type Premiaciones struct{
	Puesto string `json:"puesto"`
	Premio string `json:"premio"`
}

type CrearTorneo struct{
	IDTorneo string `json:"idTorneo"`
	Lugar string `json:"lugar"`
	Objetivo string `json:"objetivo"`
	CantidadEquipos int `json:"cantidadequipos"`
	FechaInicio string `json:"fechaInicio"`
	FechaFin string `json:"fechaFin"`
	IDDeporte *TipoDeporte `json:"idDeporte"`
	IDDivision *TipoDivision `json:"idDivision"`
	PremiacionPuesto []Premiaciones `json:"premiacion"`
}

type TipoDeporte struct{
	IDDeporte string `json:"idDeporte"`
	Nombre_Deporte string `json:"deporte"`

}

type TipoDivision struct{
	IDDivision string `json:"idDivision"`
	Division string`json:"division"`
}
