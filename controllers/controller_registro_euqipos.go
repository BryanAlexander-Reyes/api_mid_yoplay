
package controllers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/beego/beego/v2/server/web"
)

type EquipoMID struct {
	IdTorneo          int    `json:"id_torneo"`
	NombreEquipo      string `json:"nombre_equipo"`
	Categoria         string `json:"categoria"`
	CantidadJugadores int    `json:"cantidad_jugadores"`
	Region            string `json:"region"`
}

type EquipoController struct {
	web.Controller
}

func (c *EquipoController) Post() {
	var equipo EquipoMID

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &equipo); err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]interface{}{
			"ok":      false,
			"Mensaje": "Los datos del equipo no tienen un formato válido",
		}
		c.ServeJSON()
		return
	}

	equipo.NombreEquipo = strings.TrimSpace(equipo.NombreEquipo)
	equipo.Categoria = strings.TrimSpace(equipo.Categoria)
	equipo.Region = strings.TrimSpace(equipo.Region)

	if equipo.IdTorneo <= 0 {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]interface{}{
			"ok":      false,
			"Mensaje": "Debe seleccionar un torneo válido",
		}
		c.ServeJSON()
		return
	}

	if equipo.NombreEquipo == "" || equipo.Region == "" {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]interface{}{
			"ok":      false,
			"Mensaje": "El nombre del equipo y la región son obligatorios",
		}
		c.ServeJSON()
		return
	}

	if equipo.CantidadJugadores < 1 {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]interface{}{
			"ok":      false,
			"Mensaje": "La cantidad de jugadores debe ser mayor que cero",
		}
		c.ServeJSON()
		return
	}

	datos, err := json.Marshal(equipo)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]interface{}{
			"ok":      false,
			"Mensaje": "No fue posible preparar los datos del equipo",
		}
		c.ServeJSON()
		return
	}

	respuesta, err := http.Post(
		"http://localhost:8080/v1/equipo",
		"application/json",
		bytes.NewBuffer(datos),
	)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadGateway)
		c.Data["json"] = map[string]interface{}{
			"ok":      false,
			"Mensaje": "No fue posible comunicarse con el CRUD de equipos",
		}
		c.ServeJSON()
		return
	}
	defer respuesta.Body.Close()

	contenido, err := io.ReadAll(respuesta.Body)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadGateway)
		c.Data["json"] = map[string]interface{}{
			"ok":      false,
			"Mensaje": "No fue posible leer la respuesta del CRUD",
		}
		c.ServeJSON()
		return
	}

	if respuesta.StatusCode < 200 || respuesta.StatusCode >= 300 {
		c.Ctx.Output.SetStatus(http.StatusBadGateway)
		c.Data["json"] = map[string]interface{}{
			"ok":      false,
			"Mensaje": "El CRUD rechazó el registro del equipo",
			"Detalle": string(contenido),
		}
		c.ServeJSON()
		return
	}

	c.Ctx.Output.SetStatus(http.StatusCreated)
	c.Data["json"] = map[string]interface{}{
		"ok":      true,
		"Mensaje": "Equipo registrado correctamente",
		"Datos":   json.RawMessage(contenido),
	}
	c.ServeJSON()
}
