package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Equipo representa la información recibida desde la CRUD de equipos.
type Equipo struct {
	Id                int       `json:"Id"`
	NombreEquipo      string    `json:"NombreEquipo"`
	Categoria         string    `json:"Categoria"`
	CantidadJugadores int       `json:"CantidadJugadores"`
	Region            string    `json:"Region"`
	Activo            bool      `json:"Activo"`
	FechaCreacion     time.Time `json:"FechaCreacion"`
	FechaModificacion time.Time `json:"FechaModificacion"`
}

// EquipoService contiene la configuración para consultar la CRUD.
type EquipoService struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewEquipoService crea una instancia del servicio.
func NewEquipoService(baseURL string) *EquipoService {
	return &EquipoService{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// ConsultarPorID consulta un equipo por su identificador.
func (s *EquipoService) ConsultarPorID(id int) (Equipo, error) {

	url := fmt.Sprintf("%s/%d", s.BaseURL, id)

	response, err := s.HTTPClient.Get(url)

	if err != nil {
		return Equipo{}, fmt.Errorf(
			"error al consultar el equipo %d: %w", id, err,
		)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return Equipo{}, fmt.Errorf(
			"la CRUD de equipos respondió con estado %d para el ID %d",
			response.StatusCode,
			id,
		)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return Equipo{}, fmt.Errorf(
			"error al leer la respuesta del equipo %d: %w", id, err,
		)
	}

	// Intentar leer una respuesta que contenga el campo Data.
	var envoltura struct {
		Data json.RawMessage `json:"Data"`
	}

	if err := json.Unmarshal(body, &envoltura); err != nil {
		return Equipo{}, fmt.Errorf(
			"respuesta JSON inválida para el equipo %d: %w", id, err,
		)
	}

	var equipo Equipo

	// Aceptar tanto {"Data": {...}} como un objeto directo.
	if len(envoltura.Data) > 0 &&
		string(envoltura.Data) != "null" {
		err = json.Unmarshal(envoltura.Data, &equipo)
	} else {
		err = json.Unmarshal(body, &equipo)
	}

	if err != nil {
		return Equipo{}, fmt.Errorf(
			"error al decodificar el equipo %d: %w", id, err,
		)
	}

	return equipo, nil
}
