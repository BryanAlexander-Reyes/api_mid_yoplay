
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

type TipoDeporte struct {
	ID            int    `json:"id"`
	NombreDeporte string `json:"nombre_deporte"`
}

type TipoDistribucion struct {
	ID                 int    `json:"id"`
	NombreDistribucion string `json:"nombre_distribucion"`
}

type Imagen struct {
	URLImagen  string `json:"url_imagen"`
	TipoImagen string `json:"tipo_imagen"`
}

type Premiacion struct {
	Puesto string `json:"puesto"`
	Premio int    `json:"premio"`
}

type Torneo struct {
	NombreTorneo   string            `json:"nombre_torneo"`
	Lugar          string            `json:"lugar"`
	Ubicacion      string            `json:"ubicacion"`
	Objetivo       string            `json:"objetivo"`
	CantidadEquipo int               `json:"cantidad_equipo"`
	FechaInicio    string            `json:"fechaInicio"`
	FechaFin       string            `json:"fechaFin"`
	TipoDeporte    *TipoDeporte      `json:"tipo_deporte"`
	Distribucion   *TipoDistribucion `json:"tipo_distribucion"`
	Division       *TipoDistribucion `json:"divicion"`
	Imagen         *Imagen           `json:"imagen"`
	Premiaciones   []Premiacion      `json:"premiaciones"`
}

func enviarSolicitud(metodo, direccion string, datos interface{}) ([]byte, int, error) {
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

func obtenerID(respuesta []byte) (int, error) {
	var resultado struct {
		Datos struct {
			ID int `json:"Id"`
		} `json:"Data"`
		ID int `json:"Id"`
	}

	if err := json.Unmarshal(respuesta, &resultado); err != nil {
		return 0, err
	}

	if resultado.Datos.ID > 0 {
		return resultado.Datos.ID, nil
	}

	if resultado.ID > 0 {
		return resultado.ID, nil
	}

	return 0, fmt.Errorf("no se encontró el ID en la respuesta del CRUD")
}

func convertirFecha(fecha string) (string, error) {
	fecha = strings.TrimSpace(fecha)

	if fecha == "" {
		return "", fmt.Errorf("la fecha está vacía")
	}

	if fechaValidada, err := time.Parse("2006-01-02", fecha); err == nil {
		return fechaValidada.Format(time.RFC3339), nil
	}

	if fechaValidada, err := time.Parse(time.RFC3339, fecha); err == nil {
		return fechaValidada.Format(time.RFC3339), nil
	}

	return "", fmt.Errorf("la fecha debe tener formato AAAA-MM-DD o RFC3339")
}

func (c *TorneoController) CrearTorneos() {
	var torneo Torneo

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &torneo); err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "El cuerpo JSON no es válido",
		}
		c.ServeJSON()
		return
	}

	if strings.TrimSpace(torneo.NombreTorneo) == "" ||
		strings.TrimSpace(torneo.Objetivo) == "" ||
		torneo.TipoDeporte == nil ||
		torneo.TipoDeporte.ID <= 0 ||
		torneo.Distribucion == nil ||
		torneo.Distribucion.ID <= 0 ||
		torneo.Imagen == nil ||
		strings.TrimSpace(torneo.Imagen.URLImagen) == "" ||
		len(torneo.Premiaciones) == 0 {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "Faltan datos obligatorios o hay valores no válidos",
		}
		c.ServeJSON()
		return
	}

	cantidad := torneo.CantidadEquipo
	if cantidad <= 0 {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "La cantidad de equipos debe ser mayor que cero",
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
			"success": false,
			"message": "La ubicación del torneo es obligatoria",
		}
		c.ServeJSON()
		return
	}

	fechaInicio, err := convertirFecha(torneo.FechaInicio)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]interface{}{"message": "La fecha de inicio no es válida"}
		c.ServeJSON()
		return
	}

	fechaFin, err := convertirFecha(torneo.FechaFin)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]interface{}{"message": "La fecha de finalización no es válida"}
		c.ServeJSON()
		return
	}

	inicio, _ := time.Parse(time.RFC3339, fechaInicio)
	fin, _ := time.Parse(time.RFC3339, fechaFin)
	if fin.Before(inicio) {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "La fecha de finalización no puede ser anterior a la fecha de inicio",
		}
		c.ServeJSON()
		return
	}

	urlImagenes, err := web.AppConfig.String("url_imagenes")
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]interface{}{"message": "Falta configurar url_imagenes"}
		c.ServeJSON()
		return
	}

	urlTorneos, err := web.AppConfig.String("url_torneos")
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]interface{}{"message": "Falta configurar url_torneos"}
		c.ServeJSON()
		return
	}

	urlPremiaciones, err := web.AppConfig.String("url_premiaciones")
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]interface{}{"message": "Falta configurar url_premiaciones"}
		c.ServeJSON()
		return
	}

	// Crear la imagen inicialmente sin relacionarla con el torneo.
	datosImagen := map[string]interface{}{
		"UrlImagen":  strings.TrimSpace(torneo.Imagen.URLImagen),
		"TipoImagen": torneo.Imagen.TipoImagen,
		"Activo":     true,
	}

	respuestaImagen, estado, err := enviarSolicitud(http.MethodPost, urlImagenes, datosImagen)
	if err != nil || estado < 200 || estado >= 300 {
		c.Ctx.Output.SetStatus(http.StatusBadGateway)
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "No se pudo crear la imagen en el CRUD",
			"detalle": string(respuestaImagen),
		}
		c.ServeJSON()
		return
	}

	idImagen, err := obtenerID(respuestaImagen)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadGateway)
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "No se pudo obtener el ID de la imagen creada",
		}
		c.ServeJSON()
		return
	}

	// Crear el torneo usando el ID de la imagen.
	datosTorneo := map[string]interface{}{
		"IdTipoDeporte":      torneo.TipoDeporte.ID,
		"NombreTorneo":       strings.TrimSpace(torneo.NombreTorneo),
		"FechaInicio":        fechaInicio,
		"FechaFin":           fechaFin,
		"Ubicacion":          ubicacion,
		"Objetivo":           strings.TrimSpace(torneo.Objetivo),
		"IdImagen":           idImagen,
		"CantidadEquipo":     cantidad,
		"IdTipoDistribucion": torneo.Distribucion.ID,
		"Activo":             true,
	}

	respuestaTorneo, estado, err := enviarSolicitud(http.MethodPost, urlTorneos, datosTorneo)
	if err != nil || estado < 200 || estado >= 300 {
		c.Ctx.Output.SetStatus(http.StatusBadGateway)
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "No se pudo crear el torneo en el CRUD",
			"detalle": string(respuestaTorneo),
		}
		c.ServeJSON()
		return
	}

	idTorneo, err := obtenerID(respuestaTorneo)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadGateway)
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "El torneo se creó, pero no se pudo obtener su ID",
		}
		c.ServeJSON()
		return
	}

	// Relacionar la imagen con el torneo.
	datosImagenActualizados := map[string]interface{}{
		"IdTorneo": map[string]interface{}{
			"Id": idTorneo,
		},
		"UrlImagen":  strings.TrimSpace(torneo.Imagen.URLImagen),
		"TipoImagen": torneo.Imagen.TipoImagen,
		"Activo":     true,
	}

	_, estado, err = enviarSolicitud(
		http.MethodPut,
		fmt.Sprintf("%s/%d", strings.TrimRight(urlImagenes, "/"), idImagen),
		datosImagenActualizados,
	)

	if err != nil || estado < 200 || estado >= 300 {
		c.Ctx.Output.SetStatus(http.StatusMultiStatus)
		c.Data["json"] = map[string]interface{}{
			"success":   false,
			"message":   "El torneo se creó, pero no se pudo relacionar la imagen",
			"id_torneo": idTorneo,
			"id_imagen": idImagen,
		}
		c.ServeJSON()
		return
	}

	// Crear las premiaciones relacionadas con el torneo.
	var errores []string

	for posicion, premiacion := range torneo.Premiaciones {
		if strings.TrimSpace(premiacion.Puesto) == "" || premiacion.Premio < 0 {
			errores = append(errores, fmt.Sprintf("La premiación %d tiene datos no válidos", posicion+1))
			continue
		}

		datosPremiacion := map[string]interface{}{
			"IdTorneo": map[string]interface{}{
				"Id": idTorneo,
			},
			"CantidadPremiacion": premiacion.Premio,
			"Descripcion":        strings.TrimSpace(premiacion.Puesto),
			"Activo":             true,
		}

		_, estado, err := enviarSolicitud(http.MethodPost, urlPremiaciones, datosPremiacion)
		if err != nil || estado < 200 || estado >= 300 {
			errores = append(errores, fmt.Sprintf("No se pudo guardar la premiación %d", posicion+1))
		}
	}

	if len(errores) > 0 {
		c.Ctx.Output.SetStatus(http.StatusMultiStatus)
		c.Data["json"] = map[string]interface{}{
			"success":   false,
			"message":   "El torneo se creó, pero algunas premiaciones fallaron",
			"id_torneo": idTorneo,
			"errores":   errores,
		}
		c.ServeJSON()
		return
	}

	c.Ctx.Output.SetStatus(http.StatusCreated)
	c.Data["json"] = map[string]interface{}{
		"success":   true,
		"status":    http.StatusCreated,
		"message":   "Torneo, imagen y premiaciones creados correctamente",
		"id_torneo": idTorneo,
		"id_imagen": idImagen,
	}
	c.ServeJSON()
}

func (c *TorneoController) ObtenerTorneo() {
	idTexto := c.Ctx.Input.Param(":id")
	id, err := strconv.Atoi(idTexto)

	if err != nil || id <= 0 {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "El ID del torneo no es válido",
		}
		c.ServeJSON()
		return
	}

	urlTorneos, err := web.AppConfig.String("url_torneos")
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]interface{}{"message": "Falta configurar url_torneos"}
		c.ServeJSON()
		return
	}

	respuesta, estado, err := enviarSolicitud(
		http.MethodGet,
		fmt.Sprintf("%s/%d", strings.TrimRight(urlTorneos, "/"), id),
		nil,
	)

	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadGateway)
		c.Data["json"] = map[string]interface{}{"message": "No se pudo consultar el torneo"}
		c.ServeJSON()
		return
	}

	c.Ctx.Output.SetStatus(estado)
	c.Ctx.Output.Header("Content-Type", "application/json")
	c.Ctx.WriteString(string(respuesta))
}
