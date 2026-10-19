
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

type TorneoController struct {
	web.Controller
}

type TorneoTipoDeporte struct {
	ID int `json:"id"`
}

type TorneoTipoDistribucion struct {
	ID int `json:"id"`
}

type TorneoImagen struct {
	URLImagen  string `json:"url_imagen"`
	ImagenUrl  string `json:"imagenurl"`
	TipoImagen string `json:"tipo_imagen"`
}

type TorneoPremiacion struct {
	Puesto string `json:"puesto"`
	Premio int    `json:"premio"`
}

type TorneoEntrada struct {
	NombreTorneo   string                   `json:"nombre_torneo"`
	Lugar          string                   `json:"lugar"`
	Ubicacion      string                   `json:"ubicacion"`
	Objetivo       string                   `json:"objetivo"`
	CantidadEquipo int                      `json:"cantidad_equipo"`
	CantidadEquipos int                     `json:"cantidadequipos"`
	FechaInicio    string                   `json:"fechaInicio"`
	FechaFin       string                   `json:"fechaFin"`
	TipoDeporte    *TorneoTipoDeporte       `json:"tipo_deporte"`
	Distribucion   *TorneoTipoDistribucion  `json:"tipo_distribucion"`
	Division       *TorneoTipoDistribucion  `json:"divicion"`
	Imagen         *TorneoImagen            `json:"imagen"`
	Premiaciones   []TorneoPremiacion       `json:"premiaciones"`
}

func torneoSolicitud(metodo, direccion string, datos interface{}) ([]byte, int, error) {
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

	cliente := &http.Client{
		Timeout: 15 * time.Second,
	}

	respuesta, err := cliente.Do(solicitud)
	if err != nil {
		return nil, 0, err
	}
	defer respuesta.Body.Close()

	respuestaCuerpo, err := io.ReadAll(respuesta.Body)
	if err != nil {
		return nil, respuesta.StatusCode, err
	}

	return respuestaCuerpo, respuesta.StatusCode, nil
}

func torneoObtenerID(respuesta []byte) (int, error) {
	var datos interface{}

	if err := json.Unmarshal(respuesta, &datos); err != nil {
		return 0, err
	}

	return torneoBuscarID(datos)
}

func torneoBuscarID(datos interface{}) (int, error) {
	switch valor := datos.(type) {
	case map[string]interface{}:
		for _, clave := range []string{
			"Id", "ID", "id",
			"id_imagen", "id_torneo",
		} {
			if dato, existe := valor[clave]; existe {
				switch numero := dato.(type) {
				case float64:
					if numero > 0 {
						return int(numero), nil
					}
				case string:
					id, err := strconv.Atoi(numero)
					if err == nil && id > 0 {
						return id, nil
					}
				}
			}
		}

		for _, clave := range []string{
			"Data", "data", "Datos", "datos",
			"resultado", "result", "imagen", "torneo",
		} {
			if dato, existe := valor[clave]; existe {
				id, err := torneoBuscarID(dato)
				if err == nil {
					return id, nil
				}
			}
		}

	case []interface{}:
		for _, elemento := range valor {
			id, err := torneoBuscarID(elemento)
			if err == nil {
				return id, nil
			}
		}
	}

	return 0, fmt.Errorf("no se encontró el identificador en la respuesta del CRUD")
}

func torneoFecha(fecha string) (string, error) {
	fecha = strings.TrimSpace(fecha)

	if fecha == "" {
		return "", fmt.Errorf("la fecha es obligatoria")
	}

	if valor, err := time.Parse("2006-01-02", fecha); err == nil {
		return valor.Format(time.RFC3339), nil
	}

	if valor, err := time.Parse(time.RFC3339, fecha); err == nil {
		return valor.Format(time.RFC3339), nil
	}

	return "", fmt.Errorf("la fecha debe tener formato AAAA-MM-DD")
}

func (c *TorneoController) CrearTorneos() {
	var torneo TorneoEntrada

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &torneo); err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]interface{}{
			"ok":      false,
			"Mensaje": "El JSON enviado no es válido",
		}
		c.ServeJSON()
		return
	}

	if strings.TrimSpace(torneo.NombreTorneo) == "" {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]interface{}{
			"ok":      false,
			"Mensaje": "El nombre del torneo es obligatorio",
		}
		c.ServeJSON()
		return
	}

	if strings.TrimSpace(torneo.Objetivo) == "" {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]interface{}{
			"ok":      false,
			"Mensaje": "El objetivo es obligatorio",
		}
		c.ServeJSON()
		return
	}

	ubicacion := strings.TrimSpace(torneo.Ubicacion)
	if ubicacion == "" {
		ubicacion = strings.TrimSpace(torneo.Lugar)
	}

	if ubicacion == "" {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]interface{}{
			"ok":      false,
			"Mensaje": "El lugar del torneo es obligatorio",
		}
		c.ServeJSON()
		return
	}

	cantidad := torneo.CantidadEquipo
	if cantidad == 0 {
		cantidad = torneo.CantidadEquipos
	}

	if cantidad < 1 {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]interface{}{
			"ok":      false,
			"Mensaje": "La cantidad de equipos debe ser mayor que cero",
		}
		c.ServeJSON()
		return
	}

	if torneo.TipoDeporte == nil || torneo.TipoDeporte.ID < 1 {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]interface{}{
			"ok":      false,
			"Mensaje": "Debe seleccionar un tipo de deporte válido",
		}
		c.ServeJSON()
		return
	}

	distribucion := torneo.Distribucion
	if distribucion == nil {
		distribucion = torneo.Division
	}

	if distribucion == nil || distribucion.ID < 1 {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]interface{}{
			"ok":      false,
			"Mensaje": "Debe seleccionar una división para el torneo",
		}
		c.ServeJSON()
		return
	}

	if torneo.Imagen == nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]interface{}{
			"ok":      false,
			"Mensaje": "La imagen del torneo es obligatoria",
		}
		c.ServeJSON()
		return
	}

	urlImagen := strings.TrimSpace(torneo.Imagen.URLImagen)
	if urlImagen == "" {
		urlImagen = strings.TrimSpace(torneo.Imagen.ImagenUrl)
	}

	if urlImagen == "" {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]interface{}{
			"ok":      false,
			"Mensaje": "La URL de la imagen no puede estar vacía",
		}
		c.ServeJSON()
		return
	}

	if len(torneo.Premiaciones) == 0 {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]interface{}{
			"ok":      false,
			"Mensaje": "Debe registrar al menos una premiación",
		}
		c.ServeJSON()
		return
	}

	for i, premiacion := range torneo.Premiaciones {
		if strings.TrimSpace(premiacion.Puesto) == "" || premiacion.Premio < 0 {
			c.Ctx.Output.SetStatus(http.StatusBadRequest)
			c.Data["json"] = map[string]interface{}{
				"ok":      false,
				"Mensaje": fmt.Sprintf("La premiación %d tiene datos no válidos", i+1),
			}
			c.ServeJSON()
			return
		}
	}

	fechaInicio, err := torneoFecha(torneo.FechaInicio)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]interface{}{
			"ok":      false,
			"Mensaje": "La fecha de inicio no es válida",
		}
		c.ServeJSON()
		return
	}

	fechaFin, err := torneoFecha(torneo.FechaFin)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]interface{}{
			"ok":      false,
			"Mensaje": "La fecha de finalización no es válida",
		}
		c.ServeJSON()
		return
	}

	inicio, _ := time.Parse(time.RFC3339, fechaInicio)
	fin, _ := time.Parse(time.RFC3339, fechaFin)

	if fin.Before(inicio) {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]interface{}{
			"ok":      false,
			"Mensaje": "La fecha de finalización no puede ser anterior a la fecha de inicio",
		}
		c.ServeJSON()
		return
	}

	urlImagenes, errImagen := web.AppConfig.String("url_imagenes")
	urlTorneos, errTorneo := web.AppConfig.String("url_torneos")
	urlPremiaciones, errPremiacion := web.AppConfig.String("url_premiaciones")

	if errImagen != nil || errTorneo != nil || errPremiacion != nil ||
		strings.TrimSpace(urlImagenes) == "" ||
		strings.TrimSpace(urlTorneos) == "" ||
		strings.TrimSpace(urlPremiaciones) == "" {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]interface{}{
			"ok":      false,
			"Mensaje": "Falta configurar url_imagenes, url_torneos o url_premiaciones en app.conf",
		}
		c.ServeJSON()
		return
	}



	tipoImagen := strings.TrimSpace(torneo.Imagen.TipoImagen)
	if tipoImagen == "" {
		tipoImagen = "afiche"
	}

	datosImagen := map[string]interface{}{
		"UrlImagen":  urlImagen,
		"TipoImagen": tipoImagen,
		"Activo":     true,
	}

	respuesta, estado, err := torneoSolicitud(
		http.MethodPost,
		urlImagenes,
		datosImagen,
	)

	if err != nil || estado < 200 || estado >= 300 {
		c.Ctx.Output.SetStatus(http.StatusBadGateway)
		c.Data["json"] = map[string]interface{}{
			"ok":      false,
			"Mensaje": "No se pudo crear la imagen en el CRUD",
			"detalle": string(respuesta),
		}
		c.ServeJSON()
		return
	}

	idImagen, err := torneoObtenerID(respuesta)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadGateway)
		c.Data["json"] = map[string]interface{}{
			"ok":      false,
			"Mensaje": "La imagen se creó, pero no se pudo obtener su identificador",
			"detalle": string(respuesta),
		}
		c.ServeJSON()
		return
	}


	datosTorneo := map[string]interface{}{
		"IdTipoDeporte": map[string]interface{}{
			"Id": torneo.TipoDeporte.ID,
		},
		"NombreTorneo": strings.TrimSpace(torneo.NombreTorneo),
		"FechaInicio":  fechaInicio,
		"FechaFin":     fechaFin,
		"Ubicacion":    ubicacion,
		"Objetivo":     strings.TrimSpace(torneo.Objetivo),
		"IdImagen": map[string]interface{}{
			"Id": idImagen,
		},
		"CantidadEquipo": cantidad,
		"IdTipoDistribucion": map[string]interface{}{
			"Id": distribucion.ID,
		},
		"Activo": true,
	}

	respuesta, estado, err = torneoSolicitud(
		http.MethodPost,
		urlTorneos,
		datosTorneo,
	)

	if err != nil || estado < 200 || estado >= 300 {
		c.Ctx.Output.SetStatus(http.StatusBadGateway)
		c.Data["json"] = map[string]interface{}{
			"ok":       false,
			"Mensaje":  "No se pudo crear el torneo en el CRUD",
			"id_imagen": idImagen,
			"detalle":  string(respuesta),
		}
		c.ServeJSON()
		return
	}

	idTorneo, err := torneoObtenerID(respuesta)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadGateway)
		c.Data["json"] = map[string]interface{}{
			"ok":       false,
			"Mensaje":  "El torneo se creó, pero no se pudo obtener su identificador",
			"id_imagen": idImagen,
			"detalle":  string(respuesta),
		}
		c.ServeJSON()
		return
	}

	datosImagenActualizados := map[string]interface{}{
		"IdTorneo": map[string]interface{}{
			"Id": idTorneo,
		},
		"UrlImagen":  urlImagen,
		"TipoImagen": tipoImagen,
		"Activo":     true,
	}

	respuesta, estado, err = torneoSolicitud(
		http.MethodPut,
		fmt.Sprintf("%s/%d", strings.TrimRight(urlImagenes, "/"), idImagen),
		datosImagenActualizados,
	)

	if err != nil || estado < 200 || estado >= 300 {
		c.Ctx.Output.SetStatus(http.StatusBadGateway)
		c.Data["json"] = map[string]interface{}{
			"ok":        false,
			"Mensaje":   "El torneo se creó, pero no se pudo relacionar la imagen",
			"id_torneo": idTorneo,
			"id_imagen": idImagen,
			"detalle":   string(respuesta),
		}
		c.ServeJSON()
		return
	}
	var errores []string

	for i, premiacion := range torneo.Premiaciones {
		datosPremiacion := map[string]interface{}{
			"IdTorneo": map[string]interface{}{
				"Id": idTorneo,
			},
			"CantidadPremiacion": premiacion.Premio,
			"Descripcion":        strings.TrimSpace(premiacion.Puesto),
			"Activo":             true,
		}

		respuesta, estado, err = torneoSolicitud(
			http.MethodPost,
			urlPremiaciones,
			datosPremiacion,
		)

		if err != nil || estado < 200 || estado >= 300 {
			errores = append(
				errores,
				fmt.Sprintf("No se pudo guardar la premiación %d: %s", i+1, string(respuesta)),
			)
		}
	}

	if len(errores) > 0 {
		c.Ctx.Output.SetStatus(http.StatusMultiStatus)
		c.Data["json"] = map[string]interface{}{
			"ok":        false,
			"Mensaje":   "El torneo se creó, pero algunas premiaciones fallaron",
			"id_torneo": idTorneo,
			"id_imagen": idImagen,
			"errores":   errores,
		}
		c.ServeJSON()
		return
	}

	c.Ctx.Output.SetStatus(http.StatusCreated)
	c.Data["json"] = map[string]interface{}{
		"ok":        true,
		"Mensaje":   "Torneo, imagen y premiaciones creados correctamente",
		"id_torneo": idTorneo,
		"id_imagen": idImagen,
	}
	c.ServeJSON()
}
