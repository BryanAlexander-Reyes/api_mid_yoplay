
package controllers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/beego/beego/v2/server/web"
)

type Regla struct {
	ID int `json:"id_regla"`
}

type Reglamento struct {
	Idioma   string  `json:"idioma"`
	Estandar string  `json:"estandar"`
	Regla    []Regla `json:"regla"`
}

type ReglamentoController struct {
	web.Controller
}

func (c *ReglamentoController) CrearReglamento() {
	var reglamento Reglamento

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &reglamento); err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]interface{}{
			"ok":      false,
			"mensaje": "El cuerpo de la solicitud no tiene un formato válido",
		}
		c.ServeJSON()
		return
	}

	reglamento.Idioma = strings.TrimSpace(reglamento.Idioma)
	reglamento.Estandar = strings.TrimSpace(reglamento.Estandar)

	if reglamento.Idioma == "" {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]interface{}{
			"ok":      false,
			"mensaje": "El idioma es obligatorio",
		}
		c.ServeJSON()
		return
	}

	if reglamento.Estandar == "" {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]interface{}{
			"ok":      false,
			"mensaje": "Debe elegir un reglamento estándar",
		}
		c.ServeJSON()
		return
	}

	datos, err := json.Marshal(reglamento)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]interface{}{
			"ok":      false,
			"mensaje": "No fue posible preparar los datos del reglamento",
		}
		c.ServeJSON()
		return
	}

	urlReglamento, err := web.AppConfig.String("url_reglamentos")
	if err != nil || strings.TrimSpace(urlReglamento) == "" {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]interface{}{
			"ok":      false,
			"mensaje": "Falta configurar url_reglamentos en app.conf",
		}
		c.ServeJSON()
		return
	}

	respuesta, err := http.Post(
		urlReglamento,
		"application/json",
		bytes.NewBuffer(datos),
	)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadGateway)
		c.Data["json"] = map[string]interface{}{
			"ok":      false,
			"mensaje": "No fue posible comunicarse con el CRUD de reglamentos",
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
			"mensaje": "No fue posible leer la respuesta del CRUD",
		}
		c.ServeJSON()
		return
	}

	if respuesta.StatusCode < http.StatusOK ||
		respuesta.StatusCode >= http.StatusMultipleChoices {
		c.Ctx.Output.SetStatus(http.StatusBadGateway)
		c.Data["json"] = map[string]interface{}{
			"ok":      false,
			"mensaje": "El CRUD rechazó la creación del reglamento",
			"detalle": string(contenido),
		}
		c.ServeJSON()
		return
	}

	c.Ctx.Output.SetStatus(http.StatusCreated)
	c.Data["json"] = map[string]interface{}{
		"ok":      true,
		"mensaje": "Reglamento creado correctamente",
		"datos":   json.RawMessage(contenido),
	}
	c.ServeJSON()
}
