
package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/beego/beego/v2/server/web"
)

type Integrante struct {
	Nombre   string `json:"nombre"`
	Posicion string `json:"posicion"`
	Correo   string `json:"correo"`
	Telefono string `json:"telefono"`
}

type Equipos struct {
	IdTorneo     int          `json:"id_torneo"`
	NombreEquipo string       `json:"nombre_equipo"`
	Categoria    string       `json:"categoria"`
	Region       string       `json:"region"`
	Integrantes  []Integrante `json:"integrantes"`
}

type EquipoMIDController struct {
	web.Controller
}

func solicitudEquipo(metodo, direccion string, datos interface{}) ([]byte, int, error) {
	var cuerpo io.Reader

	if datos != nil {
		contenido, err := json.Marshal(datos)
		if err != nil {
			return nil, 0, err
		}
		cuerpo = bytes.NewReader(contenido)
	}

	solicitud, err := http.NewRequest(metodo, direccion, cuerpo)
	if err != nil {
		return nil, 0, err
	}

	solicitud.Header.Set("Content-Type", "application/json")

	cliente := &http.Client{Timeout: 15 * time.Second}

	respuesta, err := cliente.Do(solicitud)
	if err != nil {
		return nil, 0, err
	}
	defer respuesta.Body.Close()

	contenido, err := io.ReadAll(respuesta.Body)
	if err != nil {
		return nil, respuesta.StatusCode, err
	}

	return contenido, respuesta.StatusCode, nil
}

func obtenerIDEquipo(contenido []byte) (int, error) {
	var datos map[string]interface{}

	if err := json.Unmarshal(contenido, &datos); err != nil {
		return 0, err
	}

	for _, clave := range []string{"id_equipo", "IdEquipo", "IDEquipo", "id", "Id", "ID"} {
		if valor, existe := datos[clave]; existe {
			switch numero := valor.(type) {
			case float64:
				if numero > 0 {
					return int(numero), nil
				}
			case string:
				var id int
				if _, err := fmt.Sscanf(numero, "%d", &id); err == nil && id > 0 {
					return id, nil
				}
			}
		}
	}

	for _, clave := range []string{"data", "Data", "datos", "resultado"} {
		if valor, existe := datos[clave]; existe {
			contenidoInterno, err := json.Marshal(valor)
			if err == nil {
				if id, err := obtenerIDEquipo(contenidoInterno); err == nil {
					return id, nil
				}
			}
		}
	}

	return 0, fmt.Errorf("el CRUD no devolvió un identificador válido para el equipo")
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
				"mensaje": fmt.Sprintf("El integrante %d debe tener nombre y posición", i+1),
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
		"IdTorneo": map[string]interface{}{
			"Id": equipo.IdTorneo,
		},
		"NombreEquipo":     equipo.NombreEquipo,
		"Categoria":        equipo.Categoria,
		"Region":           equipo.Region,
		"CantidadJugadores": len(equipo.Integrantes),
	}

	respuesta, estado, err := solicitudEquipo(http.MethodPost, urlEquipos, datosEquipo)
	if err != nil || estado < 200 || estado >= 300 {
		c.Ctx.Output.SetStatus(http.StatusBadGateway)
		c.Data["json"] = map[string]interface{}{
			"ok":      false,
			"mensaje": "No se pudo registrar el equipo en el CRUD",
			"detalle": string(respuesta),
		}
		c.ServeJSON()
		return
	}

	idEquipo, err := obtenerIDEquipo(respuesta)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadGateway)
		c.Data["json"] = map[string]interface{}{
			"ok":      false,
			"mensaje": "El equipo se creó, pero no se pudo obtener su identificador",
			"detalle": string(respuesta),
		}
		c.ServeJSON()
		return
	}

	integrantesCreados := make([]map[string]interface{}, 0, len(equipo.Integrantes))

	for i, integrante := range equipo.Integrantes {
		datosIntegrante := map[string]interface{}{
			"IdEquipo": map[string]interface{}{
				"Id": idEquipo,
			},
			"Nombre":   integrante.Nombre,
			"Posicion": integrante.Posicion,
			"Correo":   integrante.Correo,
			"Telefono": integrante.Telefono,
			"Activo":   true,
		}

		respuestaIntegrante, estado, err := solicitudEquipo(
			http.MethodPost,
			urlIntegrantes,
			datosIntegrante,
		)

		if err != nil || estado < 200 || estado >= 300 {
			c.Ctx.Output.SetStatus(http.StatusBadGateway)
			c.Data["json"] = map[string]interface{}{
				"ok":        false,
				"mensaje":   fmt.Sprintf("El equipo se creó, pero falló el registro del integrante %d", i+1),
				"id_equipo": idEquipo,
				"detalle":   string(respuestaIntegrante),
			}
			c.ServeJSON()
			return
		}

		integrantesCreados = append(integrantesCreados, map[string]interface{}{
			"nombre":   integrante.Nombre,
			"posicion": integrante.Posicion,
			"correo":   integrante.Correo,
			"telefono": integrante.Telefono,
		})
	}

	c.Ctx.Output.SetStatus(http.StatusCreated)
	c.Data["json"] = map[string]interface{}{
		"ok":      true,
		"mensaje": "Equipo e integrantes registrados correctamente",
		"equipo": map[string]interface{}{
			"id_equipo":          idEquipo,
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
