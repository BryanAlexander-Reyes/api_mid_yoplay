package controllers

import (
	beego "github.com/beego/beego/v2/server/web"
)

// Info_torneo_midController operations for Info_torneo_mid
type Info_torneo_midController struct {
	beego.Controller
}

type TorneoInfo struct {

}

type GrupoEquipo struct {

}

type GrupoEncuentro struct {

}

type Equipo struct {

}


// URLMapping ...
func (c *Info_torneo_midController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Create
// @Description create Info_torneo_mid
// @Param	body		body 	models.Info_torneo_mid	true		"body for Info_torneo_mid content"
// @Success 201 {object} models.Info_torneo_mid
// @Failure 403 body is empty
// @router / [post]
func (c *Info_torneo_midController) Post() {
	// Es informativo.
}

// GetOne ...
// @Title GetOne
// @Description get Info_torneo_mid by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Info_torneo_mid
// @Failure 403 :id is empty
// @router /:id [get]
func (c *Info_torneo_midController) GetOne() {
	// NEcesito generar una funcionalidad en la cual, se tenga una fecha de inscripción (No visible en CLI).
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
	


}

// GetAll ...
// @Title GetAll
// @Description get Info_torneo_mid
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Info_torneo_mid
// @Failure 403
// @router / [get]
func (c *Info_torneo_midController) GetAll() {
	// Solo se traerá la información básica, como la imagen del torneo, y así.
	// Requiere de información básica del torneo.
	// Imagen


}

// Put ...
// @Title Put
// @Description update the Info_torneo_mid
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Info_torneo_mid	true		"body for Info_torneo_mid content"
// @Success 200 {object} models.Info_torneo_mid
// @Failure 403 :id is not int
// @router /:id [put]
func (c *Info_torneo_midController) Put() {
	// Es un modal informativo.
}

// Delete ...
// @Title Delete
// @Description delete the Info_torneo_mid
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *Info_torneo_midController) Delete() {
	// Es un modal informativo.
}
