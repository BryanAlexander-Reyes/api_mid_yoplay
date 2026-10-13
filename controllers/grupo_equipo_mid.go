package controllers

import (
	beego "github.com/beego/beego/v2/server/web"
)

// Grupo_equipo_midController operations for Grupo_equipo_mid
type Grupo_equipo_midController struct {
	beego.Controller
}

// URLMapping ...
func (c *Grupo_equipo_midController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Create
// @Description create Grupo_equipo_mid
// @Param	body		body 	models.Grupo_equipo_mid	true		"body for Grupo_equipo_mid content"
// @Success 201 {object} models.Grupo_equipo_mid
// @Failure 403 body is empty
// @router / [post]
func (c *Grupo_equipo_midController) Post() {

}

// GetOne ...
// @Title GetOne
// @Description get Grupo_equipo_mid by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Grupo_equipo_mid
// @Failure 403 :id is empty
// @router /:id [get]
func (c *Grupo_equipo_midController) GetOne() {

}

// GetAll ...
// @Title GetAll
// @Description get Grupo_equipo_mid
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Grupo_equipo_mid
// @Failure 403
// @router / [get]
func (c *Grupo_equipo_midController) GetAll() {

}

// Put ...
// @Title Put
// @Description update the Grupo_equipo_mid
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Grupo_equipo_mid	true		"body for Grupo_equipo_mid content"
// @Success 200 {object} models.Grupo_equipo_mid
// @Failure 403 :id is not int
// @router /:id [put]
func (c *Grupo_equipo_midController) Put() {

}

// Delete ...
// @Title Delete
// @Description delete the Grupo_equipo_mid
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *Grupo_equipo_midController) Delete() {

}
