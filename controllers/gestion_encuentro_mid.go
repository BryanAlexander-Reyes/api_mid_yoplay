package controllers

import (
	beego "github.com/beego/beego/v2/server/web"
)

// Gestion_encuentro_midController operations for Gestion_encuentro_mid
type Gestion_encuentro_midController struct {
	beego.Controller
}

// URLMapping ...
func (c *Gestion_encuentro_midController) URLMapping() {
}

// Post ...
// @Title Create
// @Description create Gestion_encuentro_mid
// @Param	body		body 	models.Gestion_encuentro_mid	true		"body for Gestion_encuentro_mid content"
// @Success 201 {object} models.Gestion_encuentro_mid
// @Failure 403 body is empty
// @router / [post]
func (c *Gestion_encuentro_midController) Post() {

}

// GetOne ...
// @Title GetOne
// @Description get Gestion_encuentro_mid by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Gestion_encuentro_mid
// @Failure 403 :id is empty
// @router /:id [get]
func (c *Gestion_encuentro_midController) GetOne() {

}

// GetAll ...
// @Title GetAll
// @Description get Gestion_encuentro_mid
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Gestion_encuentro_mid
// @Failure 403
// @router / [get]
func (c *Gestion_encuentro_midController) GetAll() {

}

// Put ...
// @Title Put
// @Description update the Gestion_encuentro_mid
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Gestion_encuentro_mid	true		"body for Gestion_encuentro_mid content"
// @Success 200 {object} models.Gestion_encuentro_mid
// @Failure 403 :id is not int
// @router /:id [put]
func (c *Gestion_encuentro_midController) Put() {

}

// Delete ...
// @Title Delete
// @Description delete the Gestion_encuentro_mid
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *Gestion_encuentro_midController) Delete() {

}
