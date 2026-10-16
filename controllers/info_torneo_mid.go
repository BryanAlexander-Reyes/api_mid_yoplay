package controllers

import (
	"net/http"
	"strconv"
	"strings"

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
	ID                 int    `json:"id_torneo"`
	Nombre             string `json:"nombre"`
	Lugar              string `json:"lugar"`
	Fecha              string `json:"fecha"`
	Estado             string `json:"estado"`
	Modalidad           string `json:"modalidad"`
	EstadoDistribucion string `json:"estado_distribucion"`
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
	// Necesito generar una funcionalidad en la cual, se tenga una fecha de inscripción (No visible en CLI).
	// Cuando se finalice la fecha de inscripción, se genere un servicio que cree grupos acorde
	// a la cantidad de equipos que hay inscritos en el torneo. Los grupos y equipos están definidos en
	// el cliente al momento de crear el torneo. Tengo que ver en qué parte se guardan esos campos.
	// Requiere que se busque en grupos, equipos,
	// y torneos que tengan vinculación con el torneo, por lo cual se crea el campo idTorneo.
	// La información que aparecerá en el FrontEnd es:
	// Información Nombre, Lugar, Fecha, Ganadores.
	// Grupos, o equipos confirmados para participación en el torneo.
	// Encuentros, los cuales serán por grupos, o por equipos, dependiendo de si son torneos con fase o no.
	// Tabla de posiciones y fixture.

	// Debo traer: Torneo, Grupo, equipos del grupo, Encuentros, tabla de posiciones (URG: hablar con Brayan).
	// También requiero de saber si es por grupo o por equipo la distribución del torneo.

	// Por otra parte, para la generación de la distribución, una vez finalizada, debería enviarse un
	// correo al creador del torneo, con el fin de notificar que la distribución fué generada correctamente.
	// Otra posibilidad sería que el creador del torneo cree de manera manual, la distribución.
	// Otra opción sería que el creador sea quien ejecute el trigger para la creación automática del torneo.
	// Sin necesidad de crear un campo más en el torneo, podría hacer una funcionalidad que cada cierto tiempo
	// pueda cambiar el estado del torneo como en curso, disponible, o finalizado. Requiere de una funcionalidad
	// que debe ser entregada por el mismo MID.

	// TODO:
	// Consultar PostgreSQL para obtener:
	// torneo, grupos, equipos, encuentros, clasificación y fixture.
	//
	// No se devuelve información ficticia mientras esas consultas
	// no estén implementadas.

	c.Ctx.Output.SetStatus(http.StatusNotImplemented)
	c.Data["json"] = map[string]interface{}{
		"id_torneo": id,
		"mensaje": "La consulta integral del torneo está pendiente de conectar con el repositorio",
	}
	c.ServeJSON()
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
	// TODO: consultar y devolver la información básica
	// de los torneos desde PostgreSQL.

	c.Ctx.Output.SetStatus(http.StatusNotImplemented)
	c.Data["json"] = map[string]interface{}{
		"mensaje": "La consulta de torneos aún no está implementada",
	}
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
