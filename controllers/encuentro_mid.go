package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"api_mid_yoplay/services"

	beego "github.com/beego/beego/v2/server/web"
)

type Encuentro struct {
	Id                 int       `json:"Id"`
	Fecha              time.Time `json:"Fecha"`
	IdTorneo           int       `json:"IdTorneo"`
	IdTipoDistribucion int       `json:"IdTipoDistribucion"`
	FaseTorneo         string    `json:"FaseTorneo"`
	IdEquipo1          int       `json:"IdEquipo1"`
	ResultadoEquipo1   string    `json:"ResultadoEquipo1"`
	IdEquipo2          int       `json:"IdEquipo2"`
	ResultadoEquipo2   string    `json:"ResultadoEquipo2"`
	Activo             bool      `json:"Activo"`
	FechaCreacion      time.Time `json:"FechaCreacion"`
	FechaModificacion  time.Time `json:"FechaModificacion"`
}

type RespuestaEncuentro struct {
	Data    *Encuentro `json:"Data"`
	Messaje string     `json:"Messaje"`
	Status  int        `json:"status"`
	Succes  bool       `json:"succes"`
}

type RespuestaEncuentros struct {
	Data    []Encuentro `json:"Data"`
	Messaje string      `json:"Messaje"`
	Status  int         `json:"status"`
	Succes  bool        `json:"succes"`
}

type EncuentroFrontend struct {
	NombreEquipo1 string `json:"NombreEquipo1"`
	NombreEquipo2 string `json:"NombreEquipo2"`
	Fecha         string `json:"Fecha"`
	Hora          string `json:"Hora"`
}

// Encuentro_midController operations for Encuentro_mid.
type Encuentro_midController struct {
	beego.Controller
}

// URLMapping ...
func (c *Encuentro_midController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
func (c *Encuentro_midController) Post() {
}

// GetOne ...
// @Title GetOne
// @Description Get encounter by ID and return frontend-friendly Data.
// @Param id path string true "ID del encuentro"
// @Success 200 {object} EncuentroFrontend
// @Failure 400 ID inválido
// @Failure 502 Error al consultar las APIs
// @router /:id [get]
func (c *Encuentro_midController) GetOne() {
	// Obtener y validar el ID.
	idString := c.Ctx.Input.Param(":id")

	id, err := strconv.Atoi(idString)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		c.Data["json"] = map[string]interface{}{
			"Error": "El ID del encuentro debe ser numérico.",
		}
		c.ServeJSON()
		return
	}

	// Consultar el encuentro en la CRUD de encuentros.
	urlEncuentro := fmt.Sprintf(
		"http://localhost:8082/v1/encuentro/%d",
		id,
	)

	responseEncuentro, err := http.Get(urlEncuentro)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusBadGateway)
		c.Data["json"] = map[string]interface{}{
			"Error": "No fue posible comunicarse con la API de encuentros.",
		}
		c.ServeJSON()
		return
	}
	defer responseEncuentro.Body.Close()

	if responseEncuentro.StatusCode != http.StatusOK {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusBadGateway)
		c.Data["json"] = map[string]interface{}{
			"Error":  "El API de encuentros respondió con un estado erróneo.",
			"Status": responseEncuentro.StatusCode,
		}
		c.ServeJSON()
		return
	}

	bodyEncuentro, err := io.ReadAll(responseEncuentro.Body)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusBadGateway)
		c.Data["json"] = map[string]interface{}{
			"Error": "Error al leer la respuesta del API de encuentros.",
		}
		c.ServeJSON()
		return
	}

	// La CRUD devuelve el encuentro como un objeto JSON directo.

	var encuentro Encuentro

	var envoltura struct {
		Data json.RawMessage `json:"Data"`
	}

	if err := json.Unmarshal(bodyEncuentro, &envoltura); err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusBadGateway)
		c.Data["json"] = map[string]interface{}{
			"Error":   "Respuesta JSON inválida del API de encuentros.",
			"Detalle": err.Error(),
		}
		c.ServeJSON()
		return
	}

	if len(envoltura.Data) > 0 &&
		string(envoltura.Data) != "null" {

		err = json.Unmarshal(envoltura.Data, &encuentro)
	} else {
		err = json.Unmarshal(bodyEncuentro, &encuentro)
	}

	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusBadGateway)
		c.Data["json"] = map[string]interface{}{
			"Error":   "No fue posible decodificar el encuentro.",
			"Detalle": err.Error(),
		}
		c.ServeJSON()
		return
	}

	fmt.Printf(
		"Encuentro recibido: ID=%d, Equipo1=%d, Equipo2=%d, Fecha=%s\n",
		encuentro.Id,
		encuentro.IdEquipo1,
		encuentro.IdEquipo2,
		encuentro.Fecha.Format(time.RFC3339),
	)

	if encuentro.IdEquipo1 == 0 || encuentro.IdEquipo2 == 0 {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusBadGateway)
		c.Data["json"] = map[string]interface{}{
			"Error":         "El encuentro no contiene los identificadores de ambos equipos.",
			"RespuestaCRUD": string(bodyEncuentro),
		}
		c.ServeJSON()
		return
	}

	err = json.Unmarshal(bodyEncuentro, &encuentro)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusBadGateway)
		c.Data["json"] = map[string]interface{}{
			"Error": "Respuesta JSON inválida por parte del API de encuentros.",
		}
		c.ServeJSON()
		return
	}

	// Crear el servicio de equipos.
	equipoService := services.NewEquipoService(
		"http://localhost:8081/v1/equipo/",
	)

	// Consultar el primer equipo.
	equipo1, err := equipoService.ConsultarPorID(encuentro.IdEquipo1)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusBadGateway)
		c.Data["json"] = map[string]interface{}{
			"Error": "Error al consultar el equipo 1.",
		}
		c.ServeJSON()
		return
	}

	// Consultar el segundo equipo.
	equipo2, err := equipoService.ConsultarPorID(encuentro.IdEquipo2)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusBadGateway)
		c.Data["json"] = map[string]interface{}{
			"Error": "Error al consultar el equipo 2.",
		}
		c.ServeJSON()
		return
	}

	// Preparar únicamente los Datas necesarios para el frontend.
	encuentroFrontend := EncuentroFrontend{
		NombreEquipo1: equipo1.NombreEquipo,
		NombreEquipo2: equipo2.NombreEquipo,
		Fecha:         encuentro.Fecha.Format("02/01/2006"),
		Hora:          encuentro.Fecha.Format("15:04"),
	}

	// Devolver la respuesta.
	c.Data["json"] = encuentroFrontend
	c.ServeJSON()
}

// GetAll ...
// @Title GetAll
// @Description Get all encounters.
// @Success 200 {object} RespuestaEncuentros
// @Failure 502 Error al consultar el API
// @router / [get]

func (c *Encuentro_midController) GetAll() {
	urlEncuentros := "http://localhost:8082/v1/encuentro"

	// Consultar todos los encuentros.
	responseEncuentro, err := http.Get(urlEncuentros)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusBadGateway)
		c.Data["json"] = map[string]interface{}{
			"Error": "No fue posible comunicarse con la API de encuentros.",
		}
		c.ServeJSON()
		return
	}
	defer responseEncuentro.Body.Close()

	if responseEncuentro.StatusCode != http.StatusOK {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusBadGateway)
		c.Data["json"] = map[string]interface{}{
			"Error":  "El API de encuentros respondió con un estado erróneo.",
			"Status": responseEncuentro.StatusCode,
		}
		c.ServeJSON()
		return
	}

	// Leer la respuesta de la CRUD de encuentros.
	bodyEncuentro, err := io.ReadAll(responseEncuentro.Body)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusBadGateway)
		c.Data["json"] = map[string]interface{}{
			"Error": "Error al leer la respuesta del API de encuentros.",
		}
		c.ServeJSON()
		return
	}

	// Decodificar la lista de encuentros.
	var respuesta RespuestaEncuentros

	err = json.Unmarshal(bodyEncuentro, &respuesta)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusBadGateway)
		c.Data["json"] = map[string]interface{}{
			"Error":   "Respuesta JSON inválida del API de encuentros.",
			"Detalle": err.Error(),
		}
		c.ServeJSON()
		return
	}

	// Crear el servicio para consultar equipos.
	equipoService := services.NewEquipoService(
		"http://localhost:8081/v1/equipo",
	)

	// Preparar la lista que recibirá el frontend.
	var encuentrosFrontend []EncuentroFrontend

	// Recorrer todos los encuentros.
	for _, encuentro := range respuesta.Data {

		// Consultar el equipo 1.
		equipo1, err := equipoService.ConsultarPorID(
			encuentro.IdEquipo1,
		)
		if err != nil {
			c.Ctx.ResponseWriter.WriteHeader(http.StatusBadGateway)
			c.Data["json"] = map[string]interface{}{
				"Error":       "Error al consultar el equipo 1.",
				"IdEncuentro": encuentro.Id,
				"IdEquipo":    encuentro.IdEquipo1,
				"Detalle":     err.Error(),
			}
			c.ServeJSON()
			return
		}

		// Consultar el equipo 2.
		equipo2, err := equipoService.ConsultarPorID(
			encuentro.IdEquipo2,
		)
		if err != nil {
			c.Ctx.ResponseWriter.WriteHeader(http.StatusBadGateway)
			c.Data["json"] = map[string]interface{}{
				"Error":       "Error al consultar el equipo 2.",
				"IdEncuentro": encuentro.Id,
				"IdEquipo":    encuentro.IdEquipo2,
				"Detalle":     err.Error(),
			}
			c.ServeJSON()
			return
		}

		// Preparar los Datas para el frontend.
		encuentroFrontend := EncuentroFrontend{
			NombreEquipo1: equipo1.NombreEquipo,
			NombreEquipo2: equipo2.NombreEquipo,
			Fecha:         encuentro.Fecha.Format("02/01/2006"),
			Hora:          encuentro.Fecha.Format("15:04"),
		}

		// Agregar el encuentro a la lista.
		encuentrosFrontend = append(
			encuentrosFrontend,
			encuentroFrontend,
		)
	}

	// Devolver todos los encuentros procesados.
	c.Data["json"] = encuentrosFrontend
	c.ServeJSON()
}

// Put ...
// @Title Put
// @Description update the Encuentro_mid
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Encuentro_mid	true		"body for Encuentro_mid content"
// @Success 200 {object} models.Encuentro_mid
// @Failure 403 :id is not int
// @router /:id [put]
func (c *Encuentro_midController) Put() {

}

// Delete ...
// @Title Delete
// @Description delete the Encuentro_mid
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *Encuentro_midController) Delete() {

}
