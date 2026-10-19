
package controllers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/beego/beego/v2/server/web"
)

type Integrante struct {
	Nombre   string `json:"nombre"`
	Posicion string `json:"posicion"`
	Numero   int    `json:"numero"`
	Correo   string `json:"correo"`
	Telefono string `json:"telefono"`
}

type Equipos struct {
	IdTorneo    int          `json:"id_torneo"`
	NombreEquipo string      `json:"nombre_equipo"`
	Categoria   string       `json:"categoria"`
	Region      string       `json:"region"`
	Integrantes []Integrante `json:"integrantes"`
}

type EquipoCreado struct {
	IdEquipo int `json:"id_equipo"`
}

type EquipoMIDController struct {
	web.Controller
}

func (c *EquipoMIDController) Post() {
	var equipo Equipos

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &equipo); err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]interface{}{
			"ok":      false,
			"mensaje": "El formato de los datos no es válido",
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
			"mensaje": "Debe seleccionar un torneo válido",
		}
		c.ServeJSON()
		return
	}

	if equipo.NombreEquipo == "" || equipo.Region == "" {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]interface{}{
			"ok":      false,
			"mensaje": "El nombre del equipo y la región son obligatorios",
		}
		c.ServeJSON()
		return
	}

	if len(equipo.Integrantes) == 0 {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]interface{}{
			"ok":      false,
			"mensaje": "Debe registrar al menos un integrante",
		}
		c.ServeJSON()
		return
	}

	for i := range equipo.Integrantes {
		equipo.Integrantes[i].Nombre = strings.TrimSpace(equipo.Integrantes[i].Nombre)
		equipo.Integrantes[i].Posicion = strings.TrimSpace(equipo.Integrantes[i].Posicion)
		equipo.Integrantes[i].Correo = strings.TrimSpace(equipo.Integrantes[i].Correo)
		equipo.Integrantes[i].Telefono = strings.TrimSpace(equipo.Integrantes[i].Telefono)

		if equipo.Integrantes[i].Nombre == "" || equipo.Integrantes[i].Posicion == "" {
			c.Ctx.Output.SetStatus(http.StatusBadRequest)
			c.Data["json"] = map[string]interface{}{
				"ok":      false,
				"mensaje": "Cada integrante debe tener nombre y posición",
			}
			c.ServeJSON()
			return
		}

		if equipo.Integrantes[i].Numero < 0 {
			c.Ctx.Output.SetStatus(http.StatusBadRequest)
			c.Data["json"] = map[string]interface{}{
				"ok":      false,
				"mensaje": "El número del integrante no puede ser negativo",
			}
			c.ServeJSON()
			return
		}
	}

	urlEquipos, err := web.AppConfig.String("url_equipos")
	if err != nil || strings.TrimSpace(urlEquipos) == "" {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]interface{}{
			"ok":      false,
			"mensaje": "No está configurada la URL del CRUD de equipos",
		}
		c.ServeJSON()
		return
	}

	urlIntegrantes, err := web.AppConfig.String("url_integrantes")
	if err != nil || strings.TrimSpace(urlIntegrantes) == "" {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]interface{}{
			"ok":      false,
			"mensaje": "No está configurada la URL del CRUD de integrantes",
		}
		c.ServeJSON()
		return
	}

	datosEquipo := map[string]interface{}{
		"id_torneo":          equipo.IdTorneo,
		"nombre_equipo":      equipo.NombreEquipo,
		"categoria":          equipo.Categoria,
		"region":             equipo.Region,
		"cantidad_jugadores": len(equipo.Integrantes),
	}

	contenidoEquipo, err := json.Marshal(datosEquipo)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]interface{}{
			"ok":      false,
			"mensaje": "No fue posible preparar los datos del equipo",
		}
		c.ServeJSON()
		return
	}

	respuestaEquipo, err := http.Post(urlEquipos, "application/json", bytes.NewBuffer(contenidoEquipo))
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadGateway)
		c.Data["json"] = map[string]interface{}{
			"ok":      false,
			"mensaje": "No fue posible conectar con el CRUD de equipos",
		}
		c.ServeJSON()
		return
	}
	defer respuestaEquipo.Body.Close()

	cuerpoEquipo, err := io.ReadAll(respuestaEquipo.Body)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadGateway)
		c.Data["json"] = map[string]interface{}{
			"ok":      false,
			"mensaje": "No fue posible leer la respuesta del CRUD de equipos",
		}
		c.ServeJSON()
		return
	}

	if respuestaEquipo.StatusCode < 200 || respuestaEquipo.StatusCode >= 300 {
		c.Ctx.Output.SetStatus(http.StatusBadGateway)
		c.Data["json"] = map[string]interface{}{
			"ok":      false,
			"mensaje": "El CRUD rechazó el registro del equipo",
			"detalle": string(cuerpoEquipo),
		}
		c.ServeJSON()
		return
	}

	var equipoCreado EquipoCreado
	if err := json.Unmarshal(cuerpoEquipo, &equipoCreado); err != nil || equipoCreado.IdEquipo <= 0 {
		c.Ctx.Output.SetStatus(http.StatusBadGateway)
		c.Data["json"] = map[string]interface{}{
			"ok":      false,
			"mensaje": "El CRUD no devolvió un id_equipo válido",
		}
		c.ServeJSON()
		return
	}

	integrantesCreados := make([]map[string]interface{}, 0, len(equipo.Integrantes))

	for _, integrante := range equipo.Integrantes {
		datosIntegrante := map[string]interface{}{
			"id_equipo": equipoCreado.IdEquipo,
			"nombre":    integrante.Nombre,
			"posicion":  integrante.Posicion,
			"numero":    integrante.Numero,
			"correo":    integrante.Correo,
			"telefono":  integrante.Telefono,
		}

		contenidoIntegrante, err := json.Marshal(datosIntegrante)
		if err != nil {
			c.Ctx.Output.SetStatus(http.StatusInternalServerError)
			c.Data["json"] = map[string]interface{}{
				"ok":      false,
				"mensaje": "No fue posible preparar los datos de los integrantes",
			}
			c.ServeJSON()
			return
		}

		respuestaIntegrante, err := http.Post(urlIntegrantes, "application/json", bytes.NewBuffer(contenidoIntegrante))
		if err != nil {
			c.Ctx.Output.SetStatus(http.StatusBadGateway)
			c.Data["json"] = map[string]interface{}{
				"ok":      false,
				"mensaje": "No fue posible conectar con el CRUD de integrantes",
				"id_equipo": equipoCreado.IdEquipo,
			}
			c.ServeJSON()
			return
		}

		cuerpoIntegrante, err := io.ReadAll(respuestaIntegrante.Body)
		respuestaIntegrante.Body.Close()

		if err != nil || respuestaIntegrante.StatusCode < 200 || respuestaIntegrante.StatusCode >= 300 {
			c.Ctx.Output.SetStatus(http.StatusBadGateway)
			c.Data["json"] = map[string]interface{}{
				"ok":        false,
				"mensaje":   "El equipo fue creado, pero falló el registro de un integrante",
				"id_equipo": equipoCreado.IdEquipo,
				"detalle":   string(cuerpoIntegrante),
			}
			c.ServeJSON()
			return
		}

		integrantesCreados = append(integrantesCreados, map[string]interface{}{
			"nombre":   integrante.Nombre,
			"posicion": integrante.Posicion,
			"numero":   integrante.Numero,
		})
	}

	c.Ctx.Output.SetStatus(http.StatusCreated)
	c.Data["json"] = map[string]interface{}{
		"ok":      true,
		"mensaje": "Equipo e integrantes registrados correctamente",
		"equipo": map[string]interface{}{
			"id_equipo":          equipoCreado.IdEquipo,
			"id_torneo":          equipo.IdTorneo,
			"nombre_equipo":      equipo.NombreEquipo,
			"categoria":          equipo.Categoria,
			"region":             equipo.Region,
			"cantidad_jugadores": len(equipo.Integrantes),
		},
		"integrantes": integrantesCreados,
	}
	c.ServeJSON()
}
