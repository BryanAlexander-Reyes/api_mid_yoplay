
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

type TorneoEntrada struct {
	IdTipoDeporte      int                 `json:"IdTipoDeporte"`
	NombreTorneo       string              `json:"NombreTorneo"`
	FechaInicio        string              `json:"FechaInicio"`
	FechaFin           string              `json:"FechaFin"`
	Ubicacion          string              `json:"Ubicacion"`
	Objetivo           string              `json:"Objetivo"`
	IdImagen           int                 `json:"IdImagen"`
	FechaFase          string              `json:"FechaFase"`
	CantidadEquipo     int                 `json:"CantidadEquipo"`
	IdTipoDistribucion int                 `json:"IdTipoDistribucion"`
	CantidadGrupos	   int 				    `json:"CantidadGrupos"`
	Activo             bool                `json:"Activo"`
	Imagen             *ImagenEntrada      `json:"Imagen"`
	Premiaciones       []PremiacionEntrada `json:"Premiaciones"`
}

type ImagenEntrada struct {
	UrlImagen  string `json:"UrlImagen"`
	TipoImagen string `json:"TipoImagen"`
	Activo     bool   `json:"Activo"`
}

type RespuestaTipoDistribucion struct {
    Success bool `json:"success"`
    Status  int  `json:"status"`
    Data    struct {
        Id       int    `json:"Id"`
        Divicion string `json:"Divicion"`
    } `json:"Data"`
}

type PremiacionEntrada struct {
	IdTorneo           int    `json:"IdTorneo"`
	CantidadPremiacion int    `json:"CantidadPremiacion"`
	Descripcion        string `json:"Descripcion"`
	Activo             bool   `json:"Activo"`
}

func torneoSolicitud(metodo, direccion string, datos interface{}) ([]byte, int, error) {
	var cuerpo io.Reader

	if datos != nil {
		jsonDatos, err := json.Marshal(datos)
		if err != nil {
			return nil, 0, err
		}
		cuerpo = bytes.NewBuffer(jsonDatos)
	}

	req, err := http.NewRequest(metodo, direccion, cuerpo)
	if err != nil {
		return nil, 0, err
	}

	req.Header.Set("Content-Type", "application/json")

	cliente := &http.Client{
		Timeout: 15 * time.Second,
	}

	resp, err := cliente.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	respuesta, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}

	return respuesta, resp.StatusCode, nil
}

func torneoBuscarID(valor interface{}) int {
	switch dato := valor.(type) {
	case map[string]interface{}:
		for _, clave := range []string{"Id", "ID", "id", "id_torneo", "id_imagen"} {
			if valorID, existe := dato[clave]; existe {
				switch id := valorID.(type) {
				case float64:
					if id > 0 {
						return int(id)
					}
				case int:
					if id > 0 {
						return id
					}
				case json.Number:
					numero, err := id.Int64()
					if err == nil && numero > 0 {
						return int(numero)
					}
				case string:
					numero, err := strconv.Atoi(id)
					if err == nil && numero > 0 {
						return numero
					}
				}
			}
		}

		for _, clave := range []string{
			"Data", "data", "Datos", "datos",
			"resultado", "result", "imagen", "torneo",
		} {
			if interno, existe := dato[clave]; existe {
				id := torneoBuscarID(interno)
				if id > 0 {
					return id
				}
			}
		}

	case []interface{}:
		for _, elemento := range dato {
			id := torneoBuscarID(elemento)
			if id > 0 {
				return id
			}
		}
	}

	return 0
}

func torneoObtenerID(respuesta []byte) int {
	decodificador := json.NewDecoder(bytes.NewReader(respuesta))
	decodificador.UseNumber()

	var datos interface{}
	if err := decodificador.Decode(&datos); err != nil {
		return 0
	}

	return torneoBuscarID(datos)
}

func torneoFecha(fecha string) (string, error) {
	fecha = strings.TrimSpace(fecha)

	if fecha == "" {
		return "", fmt.Errorf("la fecha es obligatoria")
	}

	formatos := []string{
		time.RFC3339,
		"2006-01-02",
		"2006-01-02 15:04:05",
	}

	for _, formato := range formatos {
		fechaParseada, err := time.Parse(formato, fecha)
		if err == nil {
			return fechaParseada.Format(time.RFC3339), nil
		}
	}

	return "", fmt.Errorf("formato de fecha inválido: %s", fecha)
}

func (c *TorneoController) CrearTorneos() {
	responder := func(status int, datos map[string]interface{}) {
		c.Ctx.Output.SetStatus(status)
		c.Data["json"] = datos
		c.ServeJSON()
	}

	var torneo TorneoEntrada

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &torneo); err != nil {
		responder(http.StatusBadRequest, map[string]interface{}{
			"ok":      false,
			"Mensaje": "El cuerpo JSON del torneo no es válido",
			"detalle": err.Error(),
		})
		return
	}

	torneo.NombreTorneo = strings.TrimSpace(torneo.NombreTorneo)
	torneo.Ubicacion = strings.TrimSpace(torneo.Ubicacion)
	torneo.Objetivo = strings.TrimSpace(torneo.Objetivo)

	if torneo.NombreTorneo == "" {
		responder(http.StatusBadRequest, map[string]interface{}{
			"ok":      false,
			"Mensaje": "El nombre del torneo es obligatorio",
		})
		return
	}

	if torneo.IdTipoDeporte < 1 {
		responder(http.StatusBadRequest, map[string]interface{}{
			"ok":      false,
			"Mensaje": "El tipo de deporte es obligatorio",
		})
		return
	}

	if torneo.IdTipoDistribucion < 1 {
		responder(http.StatusBadRequest, map[string]interface{}{
			"ok":      false,
			"Mensaje": "Debe seleccionar un tipo de distribución para el torneo",
		})
		return
	}

	if torneo.Ubicacion == "" {
		responder(http.StatusBadRequest, map[string]interface{}{
			"ok":      false,
			"Mensaje": "La ubicación del torneo es obligatoria",
		})
		return
	}

	if torneo.Objetivo == "" {
		responder(http.StatusBadRequest, map[string]interface{}{
			"ok":      false,
			"Mensaje": "El objetivo del torneo es obligatorio",
		})
		return
	}

	if torneo.CantidadEquipo < 1 {
		responder(http.StatusBadRequest, map[string]interface{}{
			"ok":      false,
			"Mensaje": "La cantidad de equipos debe ser mayor que cero",
		})
		return
	}

	if torneo.Imagen == nil || strings.TrimSpace(torneo.Imagen.UrlImagen) == "" {
		responder(http.StatusBadRequest, map[string]interface{}{
			"ok":      false,
			"Mensaje": "La URL de la imagen no puede estar vacía",
		})
		return
	}

	if len(torneo.Premiaciones) < 1 {
		responder(http.StatusBadRequest, map[string]interface{}{
			"ok":      false,
			"Mensaje": "La premiación tiene que tener mínimo una premiación",
		})
		return
	}

	for i, premiacion := range torneo.Premiaciones {
		if strings.TrimSpace(premiacion.Descripcion) == "" {
			responder(http.StatusBadRequest, map[string]interface{}{
				"ok":      false,
				"Mensaje": fmt.Sprintf("La descripción de la premiación %d es obligatoria", i+1),
			})
			return
		}

		if premiacion.CantidadPremiacion < 1 {
			responder(http.StatusBadRequest, map[string]interface{}{
				"ok":      false,
				"Mensaje": fmt.Sprintf("La cantidad de la premiación %d no puede ser negativa", i+1),
			})
			return
		}
	}

	fechaInicio, err := torneoFecha(torneo.FechaInicio)
	if err != nil {
		responder(http.StatusBadRequest, map[string]interface{}{
			"ok":      false,
			"Mensaje": "La fecha de inicio no es válida",
			"detalle": err.Error(),
		})
		return
	}

	fechaFin, err := torneoFecha(torneo.FechaFin)
	if err != nil {
		responder(http.StatusBadRequest, map[string]interface{}{
			"ok":      false,
			"Mensaje": "La fecha de finalización no es válida",
			"detalle": err.Error(),
		})
		return
	}

	inicio, _ := time.Parse(time.RFC3339, fechaInicio)
	fin, _ := time.Parse(time.RFC3339, fechaFin)

	if fin.Before(inicio) {
		responder(http.StatusBadRequest, map[string]interface{}{
			"ok":      false,
			"Mensaje": "La fecha de finalización no puede ser anterior a la fecha de inicio",
		})
		return
	}
	urlTiposDistribucion, err := web.AppConfig.String("url_tipos_distribucion")
if err != nil || strings.TrimSpace(urlTiposDistribucion) == "" {
    responder(http.StatusInternalServerError, map[string]interface{}{
        "ok":      false,
        "Mensaje": "No está configurada la URL del CRUD de tipos de distribución",
    })
    return
}

urlTipoDistribucion := strings.TrimRight(
    strings.TrimSpace(urlTiposDistribucion), "/",
) + "/" + strconv.Itoa(torneo.IdTipoDistribucion)

respuestaDistribucion, statusDistribucion, err := torneoSolicitud(
    http.MethodGet,
    urlTipoDistribucion,
    nil,
)

if err != nil || statusDistribucion < 200 || statusDistribucion >= 300 {
    detalle := string(respuestaDistribucion)
    if err != nil {
        detalle = err.Error()
    }

    responder(http.StatusBadGateway, map[string]interface{}{
        "ok":      false,
        "Mensaje": "No se pudo consultar el tipo de distribución en el CRUD",
        "detalle": detalle,
    })
    return
}

var distribucion RespuestaTipoDistribucion

if err := json.Unmarshal(respuestaDistribucion, &distribucion); err != nil {
    responder(http.StatusBadGateway, map[string]interface{}{
        "ok":      false,
        "Mensaje": "La respuesta del tipo de distribución no es válida",
        "detalle": err.Error(),
    })
    return
}

if !distribucion.Success || distribucion.Data.Id != torneo.IdTipoDistribucion {
    responder(http.StatusBadRequest, map[string]interface{}{
        "ok":      false,
        "Mensaje": "El tipo de distribución seleccionado no existe o no es válido",
    })
    return
}

if strings.EqualFold(
    strings.TrimSpace(distribucion.Data.Divicion), "grupos",
) {
    if torneo.CantidadGrupos <= 2 {
        responder(http.StatusBadRequest, map[string]interface{}{
            "ok":      false,
            "Mensaje": "La cantidad de grupos debe ser mayor a 2",
        })
        return
    }

    if torneo.CantidadGrupos > torneo.CantidadEquipo {
        responder(http.StatusBadRequest, map[string]interface{}{
            "ok":      false,
            "Mensaje": "La cantidad de grupos no puede superar la cantidad de equipos",
        })
        return
    }
}

	urlImagenes, err := web.AppConfig.String("url_imagenes")
	if err != nil || strings.TrimSpace(urlImagenes) == "" {
		responder(http.StatusInternalServerError, map[string]interface{}{
			"ok":      false,
			"Mensaje": "No está configurada la URL del CRUD de imágenes",
		})
		return
	}

	urlTorneos, err := web.AppConfig.String("url_torneos")
	if err != nil || strings.TrimSpace(urlTorneos) == "" {
		responder(http.StatusInternalServerError, map[string]interface{}{
			"ok":      false,
			"Mensaje": "No está configurada la URL del CRUD de torneos",
		})
		return
	}

	urlPremiaciones, err := web.AppConfig.String("url_premiaciones")
	if err != nil || strings.TrimSpace(urlPremiaciones) == "" {
		responder(http.StatusInternalServerError, map[string]interface{}{
			"ok":      false,
			"Mensaje": "No está configurada la URL del CRUD de premiaciones",
		})
		return
	}

	urlImagen := strings.TrimRight(strings.TrimSpace(urlImagenes), "/")
	urlTorneo := strings.TrimRight(strings.TrimSpace(urlTorneos), "/")
	urlPremiacion := strings.TrimRight(strings.TrimSpace(urlPremiaciones), "/")

	tipoImagen := strings.TrimSpace(torneo.Imagen.TipoImagen)
	if tipoImagen == "" {
		tipoImagen = "flyer"
	}

	// 1. Crear la imagen.
	datosImagen := map[string]interface{}{
		"UrlImagen":  strings.TrimSpace(torneo.Imagen.UrlImagen),
		"TipoImagen": tipoImagen,
		"Activo":     true,
	}

	respuestaImagen, statusImagen, err := torneoSolicitud(
		http.MethodPost,
		urlImagen,
		datosImagen,
	)

	

	if err != nil || statusImagen < 200 || statusImagen >= 300 {
		detalle := string(respuestaImagen)
		if err != nil {
			detalle = err.Error()
		}

		responder(http.StatusBadGateway, map[string]interface{}{
			"ok":      false,
			"Mensaje": "Falló la creación de la imagen en el CRUD",
			"detalle": detalle,
		})
		return
	}

	idImagen := torneoObtenerID(respuestaImagen)
	if idImagen < 1 {
		responder(http.StatusBadGateway, map[string]interface{}{
			"ok":      false,
			"Mensaje": "La imagen se creó, pero no se pudo obtener su ID",
			"detalle": string(respuestaImagen),
		})
		return
	}

	// 2. Crear el torneo con los IDs de las relaciones.
datosTorneo := map[string]interface{}{
    "IdTipoDeporte":      torneo.IdTipoDeporte,
    "NombreTorneo":       torneo.NombreTorneo,
    "FechaInicio":        fechaInicio,
    "FechaFin":           fechaFin,
    "Ubicacion":          torneo.Ubicacion,
    "Objetivo":            torneo.Objetivo,
    "IdImagen": map[string]interface{}{
        "Id": idImagen,
    },
    "CantidadEquipo":      torneo.CantidadEquipo,
    "IdTipoDistribucion":  torneo.IdTipoDistribucion,
    "CantidadGrupos":      torneo.CantidadGrupos,
    "Activo":              true,
}

	if strings.TrimSpace(torneo.FechaFase) != "" {
		fechaFase, err := torneoFecha(torneo.FechaFase)
		if err != nil {
			responder(http.StatusBadRequest, map[string]interface{}{
				"ok":      false,
				"Mensaje": "La fecha de fase no es válida",
				"detalle": err.Error(),
			})
			return
		}
		datosTorneo["FechaFase"] = fechaFase
	}

	respuestaTorneo, statusTorneo, err := torneoSolicitud(
		http.MethodPost,
		urlTorneo,
		datosTorneo,
	)

	if err != nil || statusTorneo < 200 || statusTorneo >= 300 {
		detalle := string(respuestaTorneo)
		if err != nil {
			detalle = err.Error()
		}

		responder(http.StatusBadGateway, map[string]interface{}{
			"ok":        false,
			"Mensaje":   "La imagen se creó, pero falló la creación del torneo",
			"id_imagen": idImagen,
			"detalle":   detalle,
		})
		return
	}

	idTorneo := torneoObtenerID(respuestaTorneo)
	if idTorneo < 1 {
		responder(http.StatusBadGateway, map[string]interface{}{
			"ok":        false,
			"Mensaje":   "El torneo se creó, pero no se pudo obtener su ID",
			"id_imagen": idImagen,
			"detalle":   string(respuestaTorneo),
		})
		return
	}

	// 4. Crear las premiaciones asociadas al torneo.
	var idsPremiaciones []int
	var erroresPremiaciones []map[string]interface{}

	for i, premiacion := range torneo.Premiaciones {
		datosPremiacion := map[string]interface{}{
			"IdTorneo": map[string]interface{}{
				"Id": idTorneo,
			},
			"CantidadPremiacion": premiacion.CantidadPremiacion,
			"Descripcion":        strings.TrimSpace(premiacion.Descripcion),
			"Activo":             true,
		}

		respuestaPremiacion, statusPremiacion, err := torneoSolicitud(
			http.MethodPost,
			urlPremiacion,
			datosPremiacion,
		)

		if err != nil || statusPremiacion < 200 || statusPremiacion >= 300 {
			detalle := string(respuestaPremiacion)
			if err != nil {
				detalle = err.Error()
			}

			erroresPremiaciones = append(erroresPremiaciones, map[string]interface{}{
				"premiacion": i + 1,
				"descripcion": premiacion.Descripcion,
				"detalle":    detalle,
			})
			continue
		}

		if idPremiacion := torneoObtenerID(respuestaPremiacion); idPremiacion > 0 {
			idsPremiaciones = append(idsPremiaciones, idPremiacion)
		}
	}

		// 5. Devolver todos los datos creados.
	if len(erroresPremiaciones) > 0 {
		responder(http.StatusMultiStatus, map[string]interface{}{
			"ok":                   false,
			"Mensaje":              "El torneo y la imagen se crearon, pero algunas premiaciones fallaron",
			"torneo":               json.RawMessage(respuestaTorneo),
			"imagen":               json.RawMessage(respuestaImagen),
			"ids_premiaciones":     idsPremiaciones,
			"errores_premiaciones": erroresPremiaciones,
		})
		return
	}

	var datosTorneoCompleto interface{}
	var datosImagenCompleta interface{}

	if err := json.Unmarshal(respuestaTorneo, &datosTorneoCompleto); err != nil {
		datosTorneoCompleto = string(respuestaTorneo)
	}

	if err := json.Unmarshal(respuestaImagen, &datosImagenCompleta); err != nil {
		datosImagenCompleta = string(respuestaImagen)
	}

	responder(http.StatusCreated, map[string]interface{}{
		"ok":      true,
		"Mensaje": "Torneo, imagen y premiaciones creados correctamente",
		"torneo":  datosTorneoCompleto,
		"imagen":  datosImagenCompleta,
		"premiaciones": idsPremiaciones,
	})
}
