
package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/beego/beego/v2/server/web"
)

type ReglamentoController struct {
	web.Controller
}

type Reglamento struct {
	IDTorneo         int    `json:"id_torneo"`
	IDTipoReglamento int    `json:"id_tipo_reglamento"`
	IDReglas          int    `json:"id_reglas"`
	IDTipoDeporte     int    `json:"id_tipo_deporte"`
	Idioma            string `json:"idioma"`
}

func enviarSolicitudReglamento(metodo, direccion string, datos interface{}) ([]byte, int, error) {
	var contenido io.Reader

	if datos != nil {
		cuerpo, err := json.Marshal(datos)
		if err != nil {
			return nil, 0, err
		}
		contenido = bytes.NewBuffer(cuerpo)
	}

	solicitud, err := http.NewRequest(metodo, direccion, contenido)
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

	cuerpoRespuesta, err := io.ReadAll(respuesta.Body)
	if err != nil {
		return nil, respuesta.StatusCode, err
	}

	return cuerpoRespuesta, respuesta.StatusCode, nil
}

func (c *ReglamentoController) CrearReglamento() {
	var reglamento Reglamento

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &reglamento); err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "El cuerpo JSON no es válido",
		}
		c.ServeJSON()
		return
	}

	if reglamento.IDTorneo <= 0 ||
		reglamento.IDTipoReglamento <= 0 ||
		reglamento.IDReglas <= 0 ||
		reglamento.IDTipoDeporte <= 0 ||
		strings.TrimSpace(reglamento.Idioma) == "" {

		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "Todos los campos son obligatorios y los ID deben ser mayores que cero",
		}
		c.ServeJSON()
		return
	}

	urlReglamentos, err := web.AppConfig.String("url_reglamentos")
	if err != nil || strings.TrimSpace(urlReglamentos) == "" {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "Falta configurar url_reglamentos",
		}
		c.ServeJSON()
		return
	}

	datos := map[string]interface{}{
		"IdTorneo": map[string]interface{}{
			"Id": reglamento.IDTorneo,
		},
		"IdTipoReglamento": reglamento.IDTipoReglamento,
		"IdReglas":         reglamento.IDReglas,
		"IdTipoDeporte":    reglamento.IDTipoDeporte,
		"Idioma":           strings.TrimSpace(reglamento.Idioma),
		"Activo":           true,
	}

	respuesta, estado, err := enviarSolicitudReglamento(
		http.MethodPost,
		strings.TrimRight(urlReglamentos, "/"),
		datos,
	)

	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadGateway)
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "No fue posible conectar con el CRUD de reglamentos",
			"error":   err.Error(),
		}
		c.ServeJSON()
		return
	}

	if estado < 200 || estado >= 300 {
		c.Ctx.Output.SetStatus(http.StatusBadGateway)
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "El CRUD no pudo crear el reglamento",
			"estado":  estado,
			"detalle": string(respuesta),
		}
		c.ServeJSON()
		return
	}

	var resultado interface{}
	if err := json.Unmarshal(respuesta, &resultado); err != nil {
		resultado = string(respuesta)
	}

	c.Ctx.Output.SetStatus(http.StatusCreated)
	c.Data["json"] = map[string]interface{}{
		"success": true,
		"status":  http.StatusCreated,
		"message": "Reglamento creado correctamente",
		"Data":    resultado,
	}
	c.ServeJSON()
}

func (c *ReglamentoController) ObtenerReglamento() {
	idTexto := c.Ctx.Input.Param(":id")
	id, err := strconv.Atoi(idTexto)

	if err != nil || id <= 0 {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "El ID del reglamento no es válido",
		}
		c.ServeJSON()
		return
	}

	urlReglamentos, err := web.AppConfig.String("url_reglamentos")
	if err != nil || strings.TrimSpace(urlReglamentos) == "" {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "Falta configurar url_reglamentos",
		}
		c.ServeJSON()
		return
	}

	direccion := fmt.Sprintf("%s/%d", strings.TrimRight(urlReglamentos, "/"), id)

	respuesta, estado, err := enviarSolicitudReglamento(
		http.MethodGet,
		direccion,
		nil,
	)

	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadGateway)
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "No se pudo consultar el reglamento",
		}
		c.ServeJSON()
		return
	}

	c.Ctx.Output.SetStatus(estado)
	c.Ctx.Output.Header("Content-Type", "application/json")
	c.Ctx.WriteString(string(respuesta))
}
