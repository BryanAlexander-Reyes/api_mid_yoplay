
package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	beego "github.com/beego/beego/v2/server/web"
)

type ReglamentoEntrada struct {
	IDTorneo         int    `json:"id_torneo"`
	IDTipoReglamento int    `json:"id_tipo_reglamento"`
	IDTipoDeporte    int    `json:"id_tipo_deporte"`
	IDReglas         int    `json:"id_reglas"`
	Idioma           string `json:"idioma"`
}

type ReglamentoMIDController struct {
	beego.Controller
}

func solicitudReglamento(
	metodo string,
	direccion string,
	datos interface{},
) ([]byte, int, error) {
	var cuerpo io.Reader

	if datos != nil {
		contenido, err := json.Marshal(datos)
		if err != nil {
			return nil, 0, err
		}
		cuerpo = bytes.NewBuffer(contenido)
	}

	solicitud, err := http.NewRequest(metodo, direccion, cuerpo)
	if err != nil {
		return nil, 0, err
	}

	solicitud.Header.Set("Content-Type", "application/json")

	cliente := &http.Client{
		Timeout: 10 * time.Second,
	}

	respuesta, err := cliente.Do(solicitud)
	if err != nil {
		return nil, 0, err
	}
	defer respuesta.Body.Close()

	contenido, err := io.ReadAll(respuesta.Body)
	if err != nil {
		return nil, respuesta.StatusCode, err
	}

	if respuesta.StatusCode < 200 || respuesta.StatusCode >= 300 {
		return contenido, respuesta.StatusCode,
			fmt.Errorf("el CRUD respondió con estado %d: %s",
				respuesta.StatusCode, string(contenido))
	}

	return contenido, respuesta.StatusCode, nil
}

func buscarIDReglamento(datos interface{}) int {
	switch valor := datos.(type) {
	case map[string]interface{}:
		for _, campo := range []string{
			"id_reglamento",
			"IdReglamento",
			"IDReglamento",
			"id",
			"Id",
			"ID",
		} {
			if dato, existe := valor[campo]; existe {
				switch numero := dato.(type) {
				case float64:
					if numero > 0 {
						return int(numero)
					}
				case string:
					var id int
					if _, err := fmt.Sscanf(numero, "%d", &id); err == nil && id > 0 {
						return id
					}
				}
			}
		}

		for _, dato := range valor {
			if id := buscarIDReglamento(dato); id > 0 {
				return id
			}
		}

	case []interface{}:
		for _, dato := range valor {
			if id := buscarIDReglamento(dato); id > 0 {
				return id
			}
		}
	}

	return 0
}

func (c *ReglamentoMIDController) CrearReglamento() {
	var request ReglamentoEntrada

	err := json.NewDecoder(c.Ctx.Request.Body).Decode(&request)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]string{
			"Mensaje": "No se pudo leer la solicitud",
			"error":   err.Error(),
		}
		c.ServeJSON()
		return
	}

	if request.IDTorneo <= 0 {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]string{
			"Mensaje": "El torneo es obligatorio",
		}
		c.ServeJSON()
		return
	}

	if request.IDTipoReglamento <= 0 {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]string{
			"Mensaje": "El tipo de reglamento es obligatorio",
		}
		c.ServeJSON()
		return
	}

	if request.IDTipoDeporte <= 0 {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]string{
			"Mensaje": "El tipo de deporte es obligatorio",
		}
		c.ServeJSON()
		return
	}

	if request.IDReglas <= 0 {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]string{
			"Mensaje": "Debe seleccionar las reglas del reglamento",
		}
		c.ServeJSON()
		return
	}

	request.Idioma = strings.TrimSpace(request.Idioma)

	if request.Idioma == "" {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]string{
			"Mensaje": "El idioma es obligatorio",
		}
		c.ServeJSON()
		return
	}

	urlReglamentos, err := beego.AppConfig.String("url_reglamentos")
	if err != nil || strings.TrimSpace(urlReglamentos) == "" {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{
			"Mensaje": "No está configurada la dirección del CRUD de reglamentos",
		}
		c.ServeJSON()
		return
	}

	urlTorneos, err := beego.AppConfig.String("url_torneos")
	if err != nil || strings.TrimSpace(urlTorneos) == "" {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{
			"Mensaje": "No está configurada la dirección del CRUD de torneos",
		}
		c.ServeJSON()
		return
	}

	urlReglas, err := beego.AppConfig.String("url_reglas")
	if err != nil || strings.TrimSpace(urlReglas) == "" {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{
			"Mensaje": "No está configurada la dirección del CRUD de reglas",
		}
		c.ServeJSON()
		return
	}

	_, _, err = solicitudReglamento(
		http.MethodGet,
		fmt.Sprintf("%s/%d", strings.TrimRight(urlTorneos, "/"), request.IDTorneo),
		nil,
	)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadGateway)
		c.Data["json"] = map[string]string{
			"Mensaje": "No se pudo verificar el torneo en el CRUD",
			"error":   err.Error(),
		}
		c.ServeJSON()
		return
	}

	_, _, err = solicitudReglamento(
		http.MethodGet,
		fmt.Sprintf("%s/%d", strings.TrimRight(urlReglas, "/"), request.IDReglas),
		nil,
	)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadGateway)
		c.Data["json"] = map[string]string{
			"Mensaje": "No se pudieron verificar las reglas en el CRUD",
			"error":   err.Error(),
		}
		c.ServeJSON()
		return
	}

datosReglamento := map[string]interface{}{
    "IdTipoReglamento": request.IDTipoReglamento,
	"IdTorneo": map[string]interface{}{
   		 "Id": request.IDTorneo,
	},
    "IdTipoDeporte":    request.IDTipoDeporte,
    "IdReglas":         request.IDReglas,
    "Idioma":           request.Idioma,
    "Activo":           true,
}

	contenido, _, err := solicitudReglamento(
		http.MethodPost,
		urlReglamentos,
		datosReglamento,
	)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadGateway)
		c.Data["json"] = map[string]string{
			"Mensaje": "No se pudo registrar el reglamento en el CRUD",
			"error":   err.Error(),
		}

		c.ServeJSON()
		return

	}
	fmt.Println("Respuesta del CRUD de reglamentos:", string(contenido))

	var resultado interface{}

	if err := json.Unmarshal(contenido, &resultado); err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadGateway)
		c.Data["json"] = map[string]string{
			"Mensaje": "La respuesta del CRUD de reglamentos no es válida",
		}
		c.ServeJSON()
		return
	}

	idReglamento := buscarIDReglamento(resultado)

	if idReglamento <= 0 {
		c.Ctx.Output.SetStatus(http.StatusBadGateway)
		c.Data["json"] = map[string]string{
			"Mensaje": "El CRUD no devolvió el identificador del reglamento",
		}
		c.ServeJSON()
		return
	}

	datosTorneo := map[string]interface{}{
		"IdReglamento": map[string]interface{}{
			"Id": idReglamento,
		},
	}

	_, _, err = solicitudReglamento(
		http.MethodPut,
		fmt.Sprintf("%s/%d", strings.TrimRight(urlTorneos, "/"), request.IDTorneo),
		datosTorneo,
	)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadGateway)
		c.Data["json"] = map[string]interface{}{
			"Mensaje":       "El reglamento se creó, pero no se pudo vincular con el torneo",
			"id_reglamento": idReglamento,
			"error":         err.Error(),
		}
		c.ServeJSON()
		return
	}

	c.Ctx.Output.SetStatus(http.StatusCreated)
	c.Data["json"] = map[string]interface{}{
		"ok":            true,
		"Mensaje":       "Reglamento registrado y vinculado con el torneo correctamente",
		"id_reglamento": idReglamento,
		"id_torneo":     request.IDTorneo,
		"id_reglas":     request.IDReglas,
	}
	c.ServeJSON()
}
