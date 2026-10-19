package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	beego "github.com/beego/beego/v2/server/web"

	// Ajusta esta ruta al nombre del módulo definido en tu go.mod.
	"api_mid_yoplay/services"
)

// Info_torneo_midController operations for Info_torneo_mid
type Info_torneo_midController struct {
	beego.Controller
}

// Estructuras para construir la respuesta del torneo.
// Se completarán al conectar las consultas reales del repositorio.
type TorneoInfo struct {
	ID                 int    `json:"Id"`
	Deporte            int    `json:"IdTipoDeporte"`
	Nombre             string `json:"NombreTorneo"`
	Objetivo           string `json:"Objetivo"`
	Descripcion        string `json:"Descripcion"`
	DescripcionTorneo  string `json:"DescripcionTorneo"`
	Lugar              string `json:"Ubicacion"`
	Fecha              string `json:"FechaInicio"`
	FechaFin           string `json:"FechaFin"`
	FechaCierre        string `json:"FechaCierre"`
	Premiacion         struct {
		ID int `json:"Id"`
	} `json:"IdPremiacion"`
	Estado             string `json:"estado"`
}

type PremioTorneo struct {
	ID                  int    `json:"Id"`
	CantidadPremiacion  int    `json:"CantidadPremiacion"`
	Descripcion         string `json:"Descripcion"`
}

type TorneoGrupoCRUD struct {
	ID       int    `json:"Id"`
	Nombre   string `json:"Grupo"`
	IDTorneo int    `json:"IdTorneo"`
}

type TorneoGrupoEquipoCRUD struct {
	Grupo    *TorneoGrupoCRUD `json:"IdGrupo"`
	IDEquipo int              `json:"IdEquipo"`
}

type TorneoEncuentroCRUD struct {
	ID                 int       `json:"Id"`
	Fecha              time.Time `json:"Fecha"`
	IDTorneo           int       `json:"IdTorneo"`
	Fase               string    `json:"FaseTorneo"`
	IDEquipo1          int       `json:"IdEquipo1"`
	ResultadoEquipo1   string    `json:"ResultadoEquipo1"`
	IDEquipo2          int       `json:"IdEquipo2"`
	ResultadoEquipo2   string    `json:"ResultadoEquipo2"`
}

type TorneoGrupoEncuentroCRUD struct {
	Grupo     *TorneoGrupoCRUD     `json:"IdGrupo"`
	Encuentro *TorneoEncuentroCRUD `json:"IdEncuentro"`
}

type TorneoEquipoDetalle struct {
	ID     int    `json:"Id"`
	Nombre string `json:"NombreEquipo"`
}

type TorneoEncuentroDetalle struct {
	ID               int    `json:"Id"`
	IDGrupo          int    `json:"IdGrupo,omitempty"`
	NombreGrupo      string `json:"Grupo,omitempty"`
	Fase             string `json:"FaseTorneo"`
	Fecha            string `json:"Fecha"`
	Hora             string `json:"Hora"`
	NombreEquipo1    string `json:"NombreEquipo1"`
	ResultadoEquipo1 string `json:"ResultadoEquipo1"`
	NombreEquipo2    string `json:"NombreEquipo2"`
	ResultadoEquipo2 string `json:"ResultadoEquipo2"`
}

type TorneoGrupoDetalle struct {
	ID         int                     `json:"Id"`
	Nombre     string                  `json:"Grupo"`
	Equipos    []TorneoEquipoDetalle   `json:"Equipos"`
	Encuentros []TorneoEncuentroDetalle `json:"Encuentros"`
}

type RespuestaDetalleTorneo struct {
	NombreTorneo string                   `json:"NombreTorneo"`
	Ubicacion    string                   `json:"Ubicacion"`
	Fecha        string                   `json:"Fecha"`
	Estado       string                   `json:"Estado"`
	Premios      []PremioTorneo           `json:"Premios"`
	Grupos       []TorneoGrupoDetalle     `json:"Grupos"`
	Encuentros   []TorneoEncuentroDetalle `json:"Encuentros,omitempty"`
}

type GrupoEquipo struct {
	IDGrupo int      `json:"id_grupo"`
	Nombre  string   `json:"nombre"`
	Equipos []Equipo `json:"equipos"`
}

type GrupoEncuentro struct {
	ID        int    `json:"id_encuentro"`
	IDGrupo   int    `json:"id_grupo,omitempty"`
	EquipoLocal string `json:"equipo_local"`
	EquipoVisitante string `json:"equipo_visitante"`
	Fecha     string `json:"fecha,omitempty"`
	Estado    string `json:"estado,omitempty"`
}

type Equipo struct {
	ID     int    `json:"id_equipo"`
	Nombre string `json:"nombre"`
}

type RespuestaInfoTorneo struct {
	Torneo        TorneoInfo      `json:"torneo"`
	Grupos        []GrupoEquipo   `json:"grupos"`
	Encuentros    []GrupoEncuentro `json:"encuentros"`
	Clasificacion []interface{}   `json:"clasificacion"`
	Fixture       []GrupoEncuentro `json:"fixture"`
}

type RespuestaTournaments struct {
	NombreDeporte     string `json:"NombreDeporte"`
	NombreTorneo      string `json:"NombreTorneo"`
	DescripcionTorneo string `json:"DescripcionTorneo"`
	Ubicacion         string `json:"Ubicacion"`
	FechaInicio       string `json:"FechaInicio"`
	Estado            string `json:"Estado"`
}

type Deporte struct {
	ID            int    `json:"Id"`
	NombreDeporte string `json:"NombreDeporte"`
}

// Servicio de distribución configurado durante el inicio del MID.
var distribucionService *services.DistribucionService

func ConfigurarDistribucionService(
	service *services.DistribucionService,
) {
	distribucionService = service
}

// URLMapping ...
func (c *Info_torneo_midController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
	c.Mapping("GenerarDistribucion", c.GenerarDistribucion)
}

// Post ...
// @router / [post]
func (c *Info_torneo_midController) Post() {
	c.Ctx.Output.SetStatus(http.StatusNotImplemented)
	c.Data["json"] = map[string]interface{}{
		"mensaje": "La creación de información del torneo aún no está implementada",
	}
	c.ServeJSON()
}

// GetOne ...
// @router /:id [get]
func (c *Info_torneo_midController) GetOne() {
	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil || id <= 0 {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]interface{}{
			"estado":  "ERROR",
			"mensaje": "El ID del torneo debe ser un entero positivo",
		}
		c.ServeJSON()
		return
	}
	ctx := c.Ctx.Request.Context()
	var torneo TorneoInfo
	if err := consultarCRUD(ctx, fmt.Sprintf("http://localhost:8081/v1/torneo/%d", id), &torneo); err != nil {
		responderErrorConsulta(c, "No fue posible consultar el torneo.", err)
		return
	}

	if torneo.Premiacion.ID <= 0 {
		c.Ctx.Output.SetStatus(http.StatusBadGateway)
		c.Data["json"] = map[string]interface{}{
			"Error":    "El torneo no contiene un ID de premiación válido.",
			"IdTorneo": id,
		}
		c.ServeJSON()
		return
	}

	var premio PremioTorneo
	if err := consultarCRUD(
		ctx,
		fmt.Sprintf("http://localhost:8081/v1/premiacion/%d", torneo.Premiacion.ID),
		&premio,
	); err != nil {
		responderErrorConsulta(c, "No fue posible consultar la premiación del torneo.", err)
		return
	}

	var gruposCRUD []TorneoGrupoCRUD
	if err := consultarCRUD(ctx, "http://localhost:8082/v1/grupo", &gruposCRUD); err != nil {
		responderErrorConsulta(c, "No fue posible consultar los grupos.", err)
		return
	}

	gruposPorID := make(map[int]*TorneoGrupoDetalle)
	grupos := make([]TorneoGrupoDetalle, 0)
	for _, grupo := range gruposCRUD {
		if grupo.IDTorneo != id || grupo.ID <= 0 {
			continue
		}
		if _, existe := gruposPorID[grupo.ID]; existe {
			continue
		}
		detalle := &TorneoGrupoDetalle{
			ID:         grupo.ID,
			Nombre:     grupo.Nombre,
			Equipos:    make([]TorneoEquipoDetalle, 0),
			Encuentros: make([]TorneoEncuentroDetalle, 0),
		}
		gruposPorID[grupo.ID] = detalle
		grupos = append(grupos, *detalle)
	}
	for i := range grupos {
		gruposPorID[grupos[i].ID] = &grupos[i]
	}

	var relacionesEquipo []TorneoGrupoEquipoCRUD
	if err := consultarCRUD(ctx, "http://localhost:8082/v1/grupo_equipo", &relacionesEquipo); err != nil {
		responderErrorConsulta(c, "No fue posible consultar los equipos de los grupos.", err)
		return
	}

	equipoService := services.NewEquipoService("http://localhost:8081/v1/equipo")
	equiposCache := make(map[int]TorneoEquipoDetalle)
	equiposPorGrupo := make(map[int]map[int]struct{})
	for _, relacion := range relacionesEquipo {
		if relacion.Grupo == nil || relacion.Grupo.IDTorneo != id || relacion.IDEquipo <= 0 {
			continue
		}
		grupo, existe := gruposPorID[relacion.Grupo.ID]
		if !existe {
			continue
		}
		if equiposPorGrupo[relacion.Grupo.ID] == nil {
			equiposPorGrupo[relacion.Grupo.ID] = make(map[int]struct{})
		}
		if _, existe := equiposPorGrupo[relacion.Grupo.ID][relacion.IDEquipo]; existe {
			continue
		}
		equipo, existe := equiposCache[relacion.IDEquipo]
		if !existe {
			equipoCRUD, err := equipoService.ConsultarPorID(relacion.IDEquipo)
			if err != nil {
				responderErrorConsulta(c, "No fue posible consultar un equipo del torneo.", err)
				return
			}
			equipo = TorneoEquipoDetalle{ID: equipoCRUD.Id, Nombre: equipoCRUD.NombreEquipo}
			equiposCache[relacion.IDEquipo] = equipo
		}
		grupo.Equipos = append(grupo.Equipos, equipo)
		equiposPorGrupo[relacion.Grupo.ID][relacion.IDEquipo] = struct{}{}
	}

	var relacionesEncuentro []TorneoGrupoEncuentroCRUD
	if err := consultarCRUD(ctx, "http://localhost:8082/v1/grupo_encuentro", &relacionesEncuentro); err != nil {
		responderErrorConsulta(c, "No fue posible consultar los encuentros del torneo.", err)
		return
	}

	encuentros := make([]TorneoEncuentroDetalle, 0)
	encuentrosVistos := make(map[int]struct{})
	for _, relacion := range relacionesEncuentro {
		if relacion.Encuentro == nil ||
			relacion.Encuentro.IDTorneo != id ||
			relacion.Encuentro.ID <= 0 {
			continue
		}
		if _, existe := encuentrosVistos[relacion.Encuentro.ID]; existe {
			continue
		}

		encuentro := relacion.Encuentro
		detalle := TorneoEncuentroDetalle{
			ID:               encuentro.ID,
			Fase:             encuentro.Fase,
			Fecha:            encuentro.Fecha.Format("02/01/2006"),
			Hora:             encuentro.Fecha.Format("3:04pm"),
			ResultadoEquipo1: encuentro.ResultadoEquipo1,
			ResultadoEquipo2: encuentro.ResultadoEquipo2,
		}
		if relacion.Grupo != nil {
			detalle.IDGrupo = relacion.Grupo.ID
			detalle.NombreGrupo = relacion.Grupo.Nombre
		}

		for _, idEquipo := range []int{encuentro.IDEquipo1, encuentro.IDEquipo2} {
			if _, existe := equiposCache[idEquipo]; existe {
				continue
			}
			equipoCRUD, err := equipoService.ConsultarPorID(idEquipo)
			if err != nil {
				responderErrorConsulta(c, "No fue posible consultar un equipo de un encuentro.", err)
				return
			}
			equiposCache[idEquipo] = TorneoEquipoDetalle{
				ID:     equipoCRUD.Id,
				Nombre: equipoCRUD.NombreEquipo,
			}
		}
		detalle.NombreEquipo1 = equiposCache[encuentro.IDEquipo1].Nombre
		detalle.NombreEquipo2 = equiposCache[encuentro.IDEquipo2].Nombre
		encuentros = append(encuentros, detalle)
		encuentrosVistos[encuentro.ID] = struct{}{}

		if grupo, existe := gruposPorID[detalle.IDGrupo]; existe {
			grupo.Encuentros = append(grupo.Encuentros, detalle)
		}
	}

	fechaInicio, err := parseFechaTorneo(torneo.Fecha)
	if err != nil {
		responderErrorConsulta(c, "La fecha de inicio del torneo no es válida.", err)
		return
	}
	fechaFinTexto := torneo.FechaFin
	if strings.TrimSpace(fechaFinTexto) == "" {
		fechaFinTexto = torneo.FechaCierre
	}
	fechaFin, err := parseFechaTorneo(fechaFinTexto)
	if err != nil {
		responderErrorConsulta(c, "La fecha de finalización del torneo no es válida.", err)
		return
	}
	estado := estadoTorneo(fechaInicio, fechaFin, time.Now())

	respuesta := RespuestaDetalleTorneo{
		NombreTorneo: torneo.Nombre,
		Ubicacion:    torneo.Lugar,
		Fecha:        formatearRangoFechas(fechaInicio, fechaFin),
		Estado:       estado,
		Premios:      []PremioTorneo{premio},
		Grupos:       grupos,
	}
	if len(grupos) == 0 {
		respuesta.Encuentros = encuentros
	}
	c.Data["json"] = respuesta
	c.ServeJSON()
}

func consultarCRUD(ctx context.Context, url string, destino interface{}) error {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("no fue posible preparar la consulta: %w", err)
	}

	client := &http.Client{Timeout: 5 * time.Second}
	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("falló el GET a %s: %w", url, err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("no fue posible leer la respuesta de %s: %w", url, err)
	}
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("el CRUD de %s respondió con estado %d: %s", url, response.StatusCode, strings.TrimSpace(string(body)))
	}

	var envoltura struct {
		Data json.RawMessage `json:"Data"`
	}
	if err := json.Unmarshal(body, &envoltura); err == nil &&
		len(envoltura.Data) > 0 && string(envoltura.Data) != "null" {
		if err := json.Unmarshal(envoltura.Data, destino); err != nil {
			return fmt.Errorf("respuesta JSON inválida de %s: %w", url, err)
		}
		return nil
	}
	if err := json.Unmarshal(body, destino); err != nil {
		return fmt.Errorf("respuesta JSON inválida de %s: %w", url, err)
	}
	return nil
}

func responderErrorConsulta(c *Info_torneo_midController, mensaje string, err error) {
	c.Ctx.Output.SetStatus(http.StatusBadGateway)
	c.Data["json"] = map[string]interface{}{
		"Error":   mensaje,
		"Detalle": err.Error(),
	}
	c.ServeJSON()
}

func parseFechaTorneo(valor string) (time.Time, error) {
	valor = strings.TrimSpace(valor)
	for _, formato := range []string{
		time.RFC3339Nano,
		time.RFC3339,
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		"2006-01-02",
	} {
		fecha, err := time.Parse(formato, valor)
		if err == nil {
			return fecha, nil
		}
	}
	return time.Time{}, fmt.Errorf("formato de fecha no válido: %q", valor)
}

func estadoTorneo(fechaInicio, fechaFin, ahora time.Time) string {
	hoy := time.Date(ahora.Year(), ahora.Month(), ahora.Day(), 0, 0, 0, 0, ahora.Location())
	diaInicio := time.Date(fechaInicio.Year(), fechaInicio.Month(), fechaInicio.Day(), 0, 0, 0, 0, ahora.Location())
	diaFin := time.Date(fechaFin.Year(), fechaFin.Month(), fechaFin.Day(), 0, 0, 0, 0, ahora.Location())

	if hoy.Before(diaInicio) {
		return "Por Confirmar"
	}
	if hoy.After(diaFin) {
		return "Finalizado"
	}
	return "En curso"
}

func formatearRangoFechas(inicio, fin time.Time) string {
	meses := [...]string{
		"enero", "febrero", "marzo", "abril", "mayo", "junio",
		"julio", "agosto", "septiembre", "octubre", "noviembre", "diciembre",
	}
	mesInicio := meses[int(inicio.Month())-1]
	mesFin := meses[int(fin.Month())-1]

	if inicio.Year() == fin.Year() && inicio.Month() == fin.Month() {
		return fmt.Sprintf("%d al %d de %s de %d", inicio.Day(), fin.Day(), mesInicio, inicio.Year())
	}
	if inicio.Year() == fin.Year() {
		return fmt.Sprintf(
			"%d de %s al %d de %s de %d",
			inicio.Day(), mesInicio, fin.Day(), mesFin, inicio.Year(),
		)
	}
	return fmt.Sprintf(
		"%d de %s de %d al %d de %s de %d",
		inicio.Day(), mesInicio, inicio.Year(), fin.Day(), mesFin, fin.Year(),
	)
}

// GenerarDistribucion solicita la generación manual de la distribución.
// @router /:id/distribucion [post]
func (c *Info_torneo_midController) GenerarDistribucion() {
	// 1. Validar el ID del torneo.
	idTorneo, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil || idTorneo <= 0 {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]interface{}{
			"estado_distribucion": string(services.ConError),
			"mensaje": "El ID del torneo debe ser un entero positivo",
		}
		c.ServeJSON()
		return
	}

	// 2. Verificar que el servicio esté configurado.
	if distribucionService == nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]interface{}{
			"id_torneo": idTorneo,
			"estado_distribucion": string(services.ConError),
			"mensaje": "El servicio de distribución no está configurado",
		}
		c.ServeJSON()
		return
	}

	// 3. Ejecutar el servicio de negocio.
	resultado, err := distribucionService.Generar(
		c.Ctx.Request.Context(),
		idTorneo,
	)

	// 4. Traducir errores a respuestas HTTP.
	if err != nil {
		status := http.StatusInternalServerError
		mensaje := err.Error()

		switch {
		case strings.Contains(mensaje, "debe ser mayor que cero"):
			status = http.StatusBadRequest

		case strings.Contains(mensaje, "no existe"):
			status = http.StatusNotFound

		case strings.Contains(mensaje, "se necesitan al menos"):
			status = http.StatusUnprocessableEntity
		}

		c.Ctx.Output.SetStatus(status)
		c.Data["json"] = map[string]interface{}{
			"id_torneo": idTorneo,
			"estado_distribucion": string(services.ConError),
			"mensaje": mensaje,
		}
		c.ServeJSON()
		return
	}

	// 5. Si ya está en proceso, responder conflicto.
	if resultado.Estado == services.EnProceso {
		c.Ctx.Output.SetStatus(http.StatusConflict)
	} else {
		c.Ctx.Output.SetStatus(http.StatusOK)
	}

	c.Data["json"] = resultado
	c.ServeJSON()
}

// GetAll ...
// @router / [get]
func (c *Info_torneo_midController) GetAll() {
	responseTournament, err := http.Get("http://localhost:8081/v1/torneo")
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadGateway)
		c.Data["json"] = map[string]interface{}{
			"Error": "No ha sido posible comunicarse con el API de torneos.",
			"Detalle": err.Error(),
		}
		c.ServeJSON()
		return
	}
	defer responseTournament.Body.Close()

	if responseTournament.StatusCode != http.StatusOK {
		c.Ctx.Output.SetStatus(http.StatusBadGateway)
		c.Data["json"] = map[string]interface{}{
			"Error":  "El API de torneos respondió con un estado erróneo.",
			"Status": responseTournament.StatusCode,
		}
		c.ServeJSON()
		return
	}

	bodyTorneos, err := io.ReadAll(responseTournament.Body)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadGateway)
		c.Data["json"] = map[string]interface{}{
			"Error":   "Error al leer la respuesta del API de torneos.",
			"Detalle": err.Error(),
		}
		c.ServeJSON()
		return
	}

	var torneos []TorneoInfo
	if err = json.Unmarshal(bodyTorneos, &torneos); err != nil {
		var respuesta struct {
			Data []TorneoInfo `json:"Data"`
		}
		if envelopeErr := json.Unmarshal(bodyTorneos, &respuesta); envelopeErr != nil {
			c.Ctx.Output.SetStatus(http.StatusBadGateway)
			c.Data["json"] = map[string]interface{}{
				"Error":   "Respuesta JSON inválida del API de torneos.",
				"Detalle": envelopeErr.Error(),
			}
			c.ServeJSON()
			return
		}
		torneos = respuesta.Data
	}

	responseDeportes, err := http.Get("http://localhost:8083/v1/tipo_deporte")
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadGateway)
		c.Data["json"] = map[string]interface{}{
			"Error":   "No ha sido posible comunicarse con el API de tipos de deporte.",
			"Detalle": err.Error(),
		}
		c.ServeJSON()
		return
	}
	defer responseDeportes.Body.Close()

	if responseDeportes.StatusCode != http.StatusOK {
		c.Ctx.Output.SetStatus(http.StatusBadGateway)
		c.Data["json"] = map[string]interface{}{
			"Error":  "El API de tipos de deporte respondió con un estado erróneo.",
			"Status": responseDeportes.StatusCode,
		}
		c.ServeJSON()
		return
	}

	bodyDeportes, err := io.ReadAll(responseDeportes.Body)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadGateway)
		c.Data["json"] = map[string]interface{}{
			"Error":   "Error al leer la respuesta del API de tipos de deporte.",
			"Detalle": err.Error(),
		}
		c.ServeJSON()
		return
	}

	var deportes []Deporte
	if err = json.Unmarshal(bodyDeportes, &deportes); err != nil {
		var respuesta struct {
			Data []Deporte `json:"Data"`
		}
		if envelopeErr := json.Unmarshal(bodyDeportes, &respuesta); envelopeErr != nil {
			c.Ctx.Output.SetStatus(http.StatusBadGateway)
			c.Data["json"] = map[string]interface{}{
				"Error":   "Respuesta JSON inválida del API de tipos de deporte.",
				"Detalle": envelopeErr.Error(),
			}
			c.ServeJSON()
			return
		}
		deportes = respuesta.Data
	}

	deportesPorID := make(map[int]string, len(deportes))
	for _, deporte := range deportes {
		deportesPorID[deporte.ID] = deporte.NombreDeporte
	}

	parseFecha := func(valor string) (time.Time, error) {
		valor = strings.TrimSpace(valor)
		for _, formato := range []string{
			time.RFC3339Nano,
			time.RFC3339,
			"2006-01-02T15:04:05",
			"2006-01-02 15:04:05",
			"2006-01-02",
		} {
			fecha, parseErr := time.Parse(formato, valor)
			if parseErr == nil {
				return fecha, nil
			}
		}
		return time.Time{}, fmt.Errorf("formato de fecha no válido: %q", valor)
	}

	ahora := time.Now()
	hoy := time.Date(ahora.Year(), ahora.Month(), ahora.Day(), 0, 0, 0, 0, ahora.Location())
	resultado := make([]RespuestaTournaments, 0, len(torneos))
	for _, torneo := range torneos {
		nombreDeporte, existe := deportesPorID[torneo.Deporte]
		if !existe {
			c.Ctx.Output.SetStatus(http.StatusBadGateway)
			c.Data["json"] = map[string]interface{}{
				"Error":         "No se encontró el tipo de deporte del torneo.",
				"IdTorneo":      torneo.ID,
				"IdTipoDeporte": torneo.Deporte,
			}
			c.ServeJSON()
			return
		}

		fechaInicio, err := parseFecha(torneo.Fecha)
		if err != nil {
			c.Ctx.Output.SetStatus(http.StatusBadGateway)
			c.Data["json"] = map[string]interface{}{
				"Error":       "La fecha de inicio de un torneo no es válida.",
				"IdTorneo":    torneo.ID,
				"FechaInicio": torneo.Fecha,
				"Detalle":     err.Error(),
			}
			c.ServeJSON()
			return
		}

		fechaCierreTexto := torneo.FechaFin
		if strings.TrimSpace(fechaCierreTexto) == "" {
			fechaCierreTexto = torneo.FechaCierre
		}
		fechaCierre, err := parseFecha(fechaCierreTexto)
		if err != nil {
			c.Ctx.Output.SetStatus(http.StatusBadGateway)
			c.Data["json"] = map[string]interface{}{
				"Error":       "La fecha de cierre de un torneo no es válida.",
				"IdTorneo":    torneo.ID,
				"FechaCierre": fechaCierreTexto,
				"Detalle":     err.Error(),
			}
			c.ServeJSON()
			return
		}

		diaInicio := time.Date(
			fechaInicio.Year(), fechaInicio.Month(), fechaInicio.Day(),
			0, 0, 0, 0, ahora.Location(),
		)
		diaCierre := time.Date(
			fechaCierre.Year(), fechaCierre.Month(), fechaCierre.Day(),
			0, 0, 0, 0, ahora.Location(),
		)
		estado := "En curso"
		if hoy.Before(diaInicio) {
			estado = "Por Confirmar"
		} else if hoy.After(diaCierre) {
			estado = "Finalizado"
		}

		meses := [...]string{
			"enero", "febrero", "marzo", "abril", "mayo", "junio",
			"julio", "agosto", "septiembre", "octubre", "noviembre", "diciembre",
		}
		fechaFormateada := fmt.Sprintf(
			"%d de %s de %d",
			fechaInicio.Day(),
			meses[int(fechaInicio.Month())-1],
			fechaInicio.Year(),
		)

		descripcion := torneo.Objetivo
		if descripcion == "" {
			descripcion = torneo.Descripcion
		}
		if descripcion == "" {
			descripcion = torneo.DescripcionTorneo
		}

		resultado = append(resultado, RespuestaTournaments{
			NombreDeporte:     nombreDeporte,
			NombreTorneo:      torneo.Nombre,
			DescripcionTorneo: descripcion,
			Ubicacion:         torneo.Lugar,
			FechaInicio:       fechaFormateada,
			Estado:            estado,
		})
	}

	c.Data["json"] = resultado
	c.ServeJSON()
}

// Put ...
// @router /:id [put]
func (c *Info_torneo_midController) Put() {
	c.Ctx.Output.SetStatus(http.StatusNotImplemented)
	c.Data["json"] = map[string]interface{}{
		"mensaje": "La actualización de información del torneo aún no está implementada",
	}
	c.ServeJSON()
}

// Delete ...
// @router /:id [delete]
func (c *Info_torneo_midController) Delete() {
	c.Ctx.Output.SetStatus(http.StatusNotImplemented)
	c.Data["json"] = map[string]interface{}{
		"mensaje": "La eliminación de información del torneo aún no está implementada",
	}
	c.ServeJSON()
}
