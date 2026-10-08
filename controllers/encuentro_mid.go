package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	beego "github.com/beego/beego/v2/server/web"
)

type Encuentro struct {
	Id                 int       `json:"Id"`
	Fecha              time.Time `json:"Fecha"`
	IdTorneo           int       `json:"IdTorneo"`
	IdTipoDistribucion int       `json:"IdTipoDistribucion"`
	FaseTorneo         string    `json:"FaseTorneo"`
	IdEquipo1          int       `json:"IdEquipo1"`
	ResultadoEquipo1   string    `json:"ResultadoEquipo1"`
	IdEquipo2          int       `json:"IdEquipo2"`
	ResultadoEquipo2   string    `json:"ResultadoEquipo2"`
	Activo             bool      `json:"Activo"`
	FechaCreacion      time.Time `json:"FechaCreacion"`
	FechaModificacion  time.Time `json:"FechaModificacion"`
}

type RespuestaEncuentro struct {
	Data    Encuentro `json:"Data"`
	Messaje string    `json:"Messaje"`
	Status  int       `json:"status"`
	Succes  bool      `json:"succes"`
}

type RespuestaEncuentros struct {
	Data    []Encuentro `json:"Data"`
	Messaje string      `json:"Messaje"`
	Status  int         `json:"status"`
	Succes  bool        `json:"succes"`
}

// Encuentro_midController operations for Encuentro_mid
type Encuentro_midController struct {
	beego.Controller
}

// URLMapping ...
func (c *Encuentro_midController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Create
// @Description create Encuentro_mid
// @Param	body		body 	models.Encuentro_mid	true		"body for Encuentro_mid content"
// @Success 201 {object} models.Encuentro_mid
// @Failure 403 body is empty
// @router / [post]
func (c *Encuentro_midController) Post() {

}

// GetOne ...
// @Title GetOne
// @Description get Encuentro_mid by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Encuentro_mid
// @Failure 403 :id is empty
// @router /:id [get]
func (c *Encuentro_midController) GetOne() {

	idString := c.Ctx.Input.Param(":id")

	id, err := strconv.Atoi(idString)

	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest,)
		c.Data["json"] = map[string]interface{}{
			"Error" : "El ID del encuentro debe ser numérico.",
		}
		c.ServeJSON()
		return
	}

	urlEncuentro := fmt.Sprintf(
		"http://localhost:8082/v1/encuentro/%d", id,
	)

	responseEncuentro, err := http.Get(urlEncuentro)

	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusBadGateway,)

		c.Data["json"] = map[string]interface{}{
			"error": "No fué posible comunicarse con la API de encuentros",
		}
		c.ServeJSON()
		return
	}
	
	defer responseEncuentro.Body.Close()

	if responseEncuentro.StatusCode != http.StatusOK {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusBadGateway,)
		c.Data["json"] = map[string]interface{}{
			"Error" : "El API de encuentros respondió erróneamente.",
			"Status" : responseEncuentro.StatusCode,
		}
		c.ServeJSON()
		return
	}


	bodyEncuentro, err := io.ReadAll(responseEncuentro.Body,)

	fmt.Println(bodyEncuentro)

	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusBadGateway,)
		c.Data["json"] = map[string]interface{}{
			"Error" : "Error leyendo el API de encuentros.",
		}
		c.ServeJSON()
		return
	}

	var respuesta RespuestaEncuentro

	err = json.Unmarshal(bodyEncuentro, &respuesta)

	fmt.Println(err)

	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusBadGateway,)
		c.Data["json"] = map[string]interface{}{
			"Error" : "Respuesta inválida por parte del API de encuentros.",
		}
		c.ServeJSON()
		return
	}

	fmt.Println("========== GET ONE ==========")
	fmt.Println("URL:", urlEncuentro)
	fmt.Println("BODY:", string(bodyEncuentro))
	fmt.Println("TIPO: RespuestaEncuentro")

	encuentro := respuesta.Data

	encuentroCmp := Encuentro{
		Id:                 encuentro.Id,
		Fecha:              encuentro.Fecha,
		IdTorneo:           encuentro.IdTorneo,
		IdTipoDistribucion: encuentro.IdTipoDistribucion,
		FaseTorneo:         encuentro.FaseTorneo,
		IdEquipo1:          encuentro.IdEquipo1,
		ResultadoEquipo1:   encuentro.ResultadoEquipo1,
		IdEquipo2:         encuentro.IdEquipo2,
		ResultadoEquipo2:   encuentro.ResultadoEquipo2,
		Activo:             encuentro.Activo,
		FechaCreacion:      encuentro.FechaCreacion,
		FechaModificacion:  encuentro.FechaModificacion,
	}


	c.Data["json"] = encuentroCmp
	c.ServeJSON()
}

// GetAll ...
// @Title GetAll
// @Description get Encuentro_mid
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Encuentro_mid
// @Failure 403
// @router / [get]
func (c *Encuentro_midController) GetAll() {

	urlEncuentros := "http://localhost:8082/v1/encuentro"

	fmt.Println(urlEncuentros)

	responseEncuentro, err := http.Get(urlEncuentros)

	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusBadGateway)
		c.Data["json"] = map[string]interface{}{
			"Error" : "No fué posible comunicarse con la API de encuentros",
		}
		c.ServeJSON()
		return
	}

	defer responseEncuentro.Body.Close()

	if responseEncuentro.StatusCode != http.StatusOK {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusBadGateway,)
		c.Data["json"] = map[string]interface{}{
			"Error" : "El API de encuentros respondió erróneamente.",
			"Status" : responseEncuentro.StatusCode,
		}
		c.ServeJSON()
		return
	}

	bodyEncuentro, err := io.ReadAll(responseEncuentro.Body,)

	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusBadGateway,)
		c.Data["json"] = map[string]interface{}{
			"Error" : "Error en lectura del API de encuentros.",
			"Status" : responseEncuentro.StatusCode,
		}
		c.ServeJSON()
		return
	}

	var respuesta RespuestaEncuentros

	err = json.Unmarshal(bodyEncuentro, &respuesta)

	fmt.Println(err)

	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusBadGateway)
		c.Data["json"] = map[string]interface{}{
			"Error":  "Respuesta inválidaa del API de encuentros.",
			"Status": responseEncuentro.StatusCode,
			"Data" : err,
		}
		c.ServeJSON()
		return
	}

	var encuentrosCmp []Encuentro

	for _, encuentro := range respuesta.Data {

		encuentroCompleto := Encuentro{
			Id:                 encuentro.Id,
			Fecha:              encuentro.Fecha,
			IdTorneo:           encuentro.IdTorneo,
			IdTipoDistribucion: encuentro.IdTipoDistribucion,
			FaseTorneo:         encuentro.FaseTorneo,
			IdEquipo1:          encuentro.IdEquipo1,
			ResultadoEquipo1:   encuentro.ResultadoEquipo1,
			IdEquipo2:          encuentro.IdEquipo2,
			ResultadoEquipo2:   encuentro.ResultadoEquipo2,
			Activo:              encuentro.Activo,
			FechaCreacion:       encuentro.FechaCreacion,
			FechaModificacion:   encuentro.FechaModificacion,
		}

		encuentrosCmp = append(encuentrosCmp, encuentroCompleto)
	}

	c.Data["json"] = encuentrosCmp
	c.ServeJSON()
}

// Put ...
// @Title Put
// @Description update the Encuentro_mid
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Encuentro_mid	true		"body for Encuentro_mid content"
// @Success 200 {object} models.Encuentro_mid
// @Failure 403 :id is not int
// @router /:id [put]
func (c *Encuentro_midController) Put() {

}

// Delete ...
// @Title Delete
// @Description delete the Encuentro_mid
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *Encuentro_midController) Delete() {

}
