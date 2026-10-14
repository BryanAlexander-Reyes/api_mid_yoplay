
package controllers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	web "github.com/beego/beego/v2/server/web"
)

type Premios struct {
	ID     int    `json:"id,omitempty"`
	Puesto string `json:"puesto"`
	Premio string `json:"premio"`
}

type TipoDeporte struct {
	ID             int    `json:"id"`
	Nombre_Deporte string `json:"nombre_deporte"`
}

type TipoDivision struct {
	ID       int    `json:"id"`
	Division string `json:"Nombre_distribucion"`
}

type Imagen struct {
	ID        int    `json:"id,omitempty"`
	ImagenUrl string `json:"imagenurl"`
}

type CrearTorneo struct {
	NombreTorneo    string        `json:"nombre_torneo"`
	Lugar           string        `json:"lugar"`
	Objetivo        string        `json:"objetivo"`
	CantidadEquipos int           `json:"cantidadequipos"`
	FechaInicio     string        `json:"fechaInicio"`
	FechaFin        string        `json:"fechaFin"`
	TipoDeporte     *TipoDeporte  `json:"tipo_deporte"`
	Division        *TipoDivision `json:"divicion"`
	Imagen          *Imagen       `json:"imagen"`
	Premiaciones    []Premios     `json:"premiaciones"`
}

type TorneoCompleto struct {
	NombreTorneo    string        `json:"nombre_torneo"`
	Lugar           string        `json:"lugar"`
	Objetivo        string        `json:"objetivo"`
	CantidadEquipos int           `json:"cantidadequipos"`
	FechaInicio     string        `json:"fechaInicio"`
	FechaFin        string        `json:"fechaFin"`
	TipoDeporte     *TipoDeporte  `json:"tipo_deporte"`
	Division        *TipoDivision `json:"divicion"`
	Imagen          *Imagen       `json:"imagen"`
	Premiaciones    []Premios     `json:"premiaciones"`
}

type ImagenCreadaResponse struct {
	ID int `json:"id"`
}

type PremiacionCreadaResponse struct {
	ID int `json:"id"`
}

type TorneoController struct {
	web.Controller
}

// @Title Enviar Torneo
// @Description Crear imagen, premiaciones y torneo
// @Param body body controllers.CrearTorneo true "Datos para crear el torneo"
// @Success 201 {object} controllers.TorneoCompleto
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @router /api/torneo [post]

func (c *TorneoController) CrearTorneos() {
	body := c.Ctx.Input.RequestBody

	if len(body) == 0 {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]string{
			"error": "El registro del torneo está vacío",
		}
		c.ServeJSON()
		return
	}

	var torneo CrearTorneo

	err := json.Unmarshal(body, &torneo)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]string{
			"error": "El JSON enviado no tiene un formato válido",
		}
		c.ServeJSON()
		return
	}

	if strings.TrimSpace(torneo.NombreTorneo) == "" {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]string{
			"error": "El nombre del torneo es obligatorio",
		}
		c.ServeJSON()
		return
	}

	if strings.TrimSpace(torneo.Objetivo) == "" {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]string{
			"error": "El objetivo es obligatorio",
		}
		c.ServeJSON()
		return
	}

	if strings.TrimSpace(torneo.Lugar) == "" {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]string{
			"error": "El lugar es obligatorio",
		}
		c.ServeJSON()
		return
	}

	if torneo.CantidadEquipos < 1 {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]string{
			"error": "La cantidad de equipos debe ser mayor que cero",
		}
		c.ServeJSON()
		return
	}

	if len(torneo.Premiaciones) < 1 {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]string{
			"error": "Debe registrar al menos una premiación",
		}
		c.ServeJSON()
		return
	}

	if torneo.TipoDeporte == nil || torneo.TipoDeporte.ID < 1 {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]string{
			"error": "Debe seleccionar un deporte válido",
		}
		c.ServeJSON()
		return
	}

	if torneo.Division == nil || torneo.Division.ID < 1 {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]string{
			"error": "Debe seleccionar una división válida",
		}
		c.ServeJSON()
		return
	}

	if torneo.Imagen == nil ||
		strings.TrimSpace(torneo.Imagen.ImagenUrl) == "" {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]string{
			"error": "La URL de la imagen es obligatoria",
		}
		c.ServeJSON()
		return
	}

	urlImagenes, _ := web.AppConfig.String("url_imagenes")
	urlPremiaciones, _ := web.AppConfig.String("url_premiaciones")
	urlTorneos, _ := web.AppConfig.String("url_torneos")

	if urlImagenes == "" || urlPremiaciones == "" || urlTorneos == "" {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{
			"error": "Falta configurar url_imagenes, url_premiaciones o url_torneos en app.conf",
		}
		c.ServeJSON()
		return
	}


	solicitudImagen := Imagen{
		ImagenUrl: torneo.Imagen.ImagenUrl,
	}

	jsonImagen, err := json.Marshal(solicitudImagen)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{
			"error": "No fue posible convertir el JSON de imagen",
		}
		c.ServeJSON()
		return
	}

	respuestaImagen, err := http.Post(
		urlImagenes,
		"application/json",
		bytes.NewBuffer(jsonImagen),
	)

	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{
			"error": "No se pudo realizar el POST de imagen",
		}
		c.ServeJSON()
		return
	}

	defer respuestaImagen.Body.Close()

	if respuestaImagen.StatusCode < http.StatusOK ||
		respuestaImagen.StatusCode >= http.StatusMultipleChoices {
		c.Ctx.Output.SetStatus(respuestaImagen.StatusCode)
		c.Data["json"] = map[string]string{
			"error": "El CRUD no creó la imagen",
		}
		c.ServeJSON()
		return
	}

	cuerpoImagenCreada, err := io.ReadAll(respuestaImagen.Body)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{
			"error": "No fue posible leer la respuesta de imagen",
		}
		c.ServeJSON()
		return
	}

	var imagenCreada ImagenCreadaResponse

	err = json.Unmarshal(cuerpoImagenCreada, &imagenCreada)
	if err != nil || imagenCreada.ID < 1 {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{
			"error": "No fue posible obtener el ID de la imagen creada",
		}
		c.ServeJSON()
		return
	}

	jsonPremiaciones, err := json.Marshal(torneo.Premiaciones)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{
			"error": "No fue posible convertir el JSON de premiaciones",
		}
		c.ServeJSON()
		return
	}

	respuestaPremiaciones, err := http.Post(
		urlPremiaciones,
		"application/json",
		bytes.NewBuffer(jsonPremiaciones),
	)

	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{
			"error": "No se pudo realizar el POST de premiaciones",
		}
		c.ServeJSON()
		return
	}

	defer respuestaPremiaciones.Body.Close()

	if respuestaPremiaciones.StatusCode < http.StatusOK ||
		respuestaPremiaciones.StatusCode >= http.StatusMultipleChoices {
		c.Ctx.Output.SetStatus(respuestaPremiaciones.StatusCode)
		c.Data["json"] = map[string]string{
			"error": "El CRUD no creó las premiaciones",
		}
		c.ServeJSON()
		return
	}

	cuerpoPremiacionesCreadas, err := io.ReadAll(respuestaPremiaciones.Body)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{
			"error": "No fue posible leer la respuesta de premiaciones",
		}
		c.ServeJSON()
		return
	}

	var premiacionesCreadas []PremiacionCreadaResponse

	err = json.Unmarshal(cuerpoPremiacionesCreadas, &premiacionesCreadas)
	if err != nil || len(premiacionesCreadas) != len(torneo.Premiaciones) {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{
			"error": "No fue posible procesar los IDs de las premiaciones",
		}
		c.ServeJSON()
		return
	}

	for i := range premiacionesCreadas {
		if premiacionesCreadas[i].ID < 1 {
			c.Ctx.Output.SetStatus(http.StatusInternalServerError)
			c.Data["json"] = map[string]string{
				"error": "El CRUD devolvió un ID de premiación inválido",
			}
			c.ServeJSON()
			return
		}

		torneo.Premiaciones[i].ID = premiacionesCreadas[i].ID
	}

	// POST 3: Crear torneo completo

	torneo.Imagen.ID = imagenCreada.ID

	jsonTorneo, err := json.Marshal(torneo)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{
			"error": "No fue posible convertir el JSON del torneo",
		}
		c.ServeJSON()
		return
	}

	respuestaTorneo, err := http.Post(
		urlTorneos,
		"application/json",
		bytes.NewBuffer(jsonTorneo),
	)

	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{
			"error": "No se pudo realizar el POST del torneo",
		}
		c.ServeJSON()
		return
	}

	defer respuestaTorneo.Body.Close()

	if respuestaTorneo.StatusCode < http.StatusOK ||
		respuestaTorneo.StatusCode >= http.StatusMultipleChoices {
		c.Ctx.Output.SetStatus(respuestaTorneo.StatusCode)
		c.Data["json"] = map[string]string{
			"error": "El CRUD no pudo crear el torneo",
		}
		c.ServeJSON()
		return
	}

	cuerpoTorneoCreado, err := io.ReadAll(respuestaTorneo.Body)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{
			"error": "No fue posible leer la respuesta del torneo creado",
		}
		c.ServeJSON()
		return
	}

	var torneoCreado map[string]interface{}

	err = json.Unmarshal(cuerpoTorneoCreado, &torneoCreado)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{
			"error": "No fue posible procesar la respuesta del torneo creado",
		}
		c.ServeJSON()
		return
	}

	c.Ctx.Output.SetStatus(http.StatusCreated)
	c.Data["json"] = map[string]interface{}{
		"ok":                  true,
		"mensaje":             "Torneo creado exitosamente",
		"id_imagen":           imagenCreada.ID,
		"premiaciones":        premiacionesCreadas,
		"torneo":              torneoCreado,
	}
	c.ServeJSON()
}
