package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	beego "github.com/beego/beego/v2/server/web"
)

type usuarioCRUD struct {
	Nombre          string `json:"Nombre"`
	Apellido        string `json:"Apellido"`
	Email           string `json:"Email"`
	NumeroDocumento string `json:"NumeroDocumento"`
	Activo          bool   `json:"Activo"`
	IdDocumento     struct {
		TipoDocumento string `json:"TipoDocumento"`
	} `json:"IdDocumento"`
}

type integranteCRUD struct {
	Correo   string `json:"Correo"`
	IdEquipo *struct {
		NombreEquipo string `json:"NombreEquipo"`
	} `json:"IdEquipo"`
}

type usuarioDetalle struct {
	Nombre          string `json:"Nombre"`
	TipoDocumento   string `json:"TipoDocumento"`
	NumeroDocumento string `json:"NumeroDocumento"`
	Equipo          string `json:"Equipo"`
	Rol             string `json:"Rol"`
	Estado          string `json:"Estado"`
}

// Info_usuario_midController provides read-only access to a user's details.
type Info_usuario_midController struct {
	beego.Controller
}

func (c *Info_usuario_midController) GetOne() {
	baseURL, err := beego.AppConfig.String("url_usuario")
	if err != nil || strings.TrimSpace(baseURL) == "" {
		c.writeError(http.StatusInternalServerError, "No está configurada la URL del servicio de usuarios.")
		return
	}

	id := strings.TrimSpace(c.Ctx.Input.Param(":id"))
	if id == "" {
		c.writeError(http.StatusBadRequest, "El ID del usuario es obligatorio.")
		return
	}

	client := &http.Client{Timeout: 10 * time.Second}
	userURL := strings.TrimRight(baseURL, "/") + "/" + url.PathEscape(id)
	userBody, status, contentType, err := getServiceResponse(c.Ctx.Request.Context(), client, userURL, c.Ctx.Request)
	if err != nil {
		c.writeError(http.StatusBadGateway, "No fue posible comunicarse con el servicio de usuarios.")
		return
	}
	if status < http.StatusOK || status >= http.StatusMultipleChoices {
		if contentType != "" {
			c.Ctx.ResponseWriter.Header().Set("Content-Type", contentType)
		}
		c.Ctx.ResponseWriter.WriteHeader(status)
		_, _ = c.Ctx.ResponseWriter.Write(userBody)
		return
	}

	var usuario usuarioCRUD
	if err := decodeCRUDData(userBody, &usuario); err != nil {
		c.writeError(http.StatusBadGateway, "El servicio de usuarios devolvió una respuesta inválida.")
		return
	}

	nombreEquipo := "No definido"
	if strings.TrimSpace(usuario.Email) != "" {
		integrantesURL, err := beego.AppConfig.String("url_integrantes")
		if err != nil || strings.TrimSpace(integrantesURL) == "" {
			c.writeError(http.StatusInternalServerError, "No está configurada la URL del servicio de integrantes.")
			return
		}

		equipo, err := consultarEquipoUsuario(
			c.Ctx.Request.Context(),
			client,
			integrantesURL,
			usuario.Email,
			c.Ctx.Request,
		)
		if err != nil {
			c.writeError(http.StatusBadGateway, "No fue posible consultar el equipo vinculado al usuario.")
			return
		}
		if equipo != "" {
			nombreEquipo = equipo
		}
	}

	estado := "Inactivo"
	if usuario.Activo {
		estado = "Activo"
	}

	c.Data["json"] = usuarioDetalle{
		Nombre:          strings.TrimSpace(usuario.Nombre + " " + usuario.Apellido),
		TipoDocumento:   usuario.IdDocumento.TipoDocumento,
		NumeroDocumento: usuario.NumeroDocumento,
		Equipo:          nombreEquipo,
		Rol:             "Usuario",
		Estado:          estado,
	}
	c.ServeJSON()
}

func consultarEquipoUsuario(
	ctx context.Context,
	client *http.Client,
	baseURL string,
	correo string,
	originalRequest *http.Request,
) (string, error) {
	query := url.Values{}
	query.Set("query", "Correo:"+correo)
	targetURL := strings.TrimRight(baseURL, "/") + "?" + query.Encode()

	body, status, _, err := getServiceResponse(ctx, client, targetURL, originalRequest)
	if err != nil {
		return "", err
	}
	if status < http.StatusOK || status >= http.StatusMultipleChoices {
		return "", fmt.Errorf("el servicio de integrantes respondió con estado HTTP %d", status)
	}

	var resultado struct {
		Data []integranteCRUD `json:"Data"`
	}
	if err := json.Unmarshal(body, &resultado); err != nil {
		return "", fmt.Errorf("respuesta inválida del servicio de integrantes: %w", err)
	}

	for _, integrante := range resultado.Data {
		if !strings.EqualFold(strings.TrimSpace(integrante.Correo), strings.TrimSpace(correo)) {
			continue
		}
		if integrante.IdEquipo != nil {
			return strings.TrimSpace(integrante.IdEquipo.NombreEquipo), nil
		}
	}
	return "", nil
}

func getServiceResponse(
	ctx context.Context,
	client *http.Client,
	targetURL string,
	originalRequest *http.Request,
) ([]byte, int, string, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, targetURL, nil)
	if err != nil {
		return nil, 0, "", err
	}
	for _, header := range []string{"Accept", "Authorization"} {
		if value := originalRequest.Header.Get(header); value != "" {
			request.Header.Set(header, value)
		}
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, 0, "", err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, 0, "", err
	}
	return body, response.StatusCode, response.Header.Get("Content-Type"), nil
}

func decodeCRUDData(body []byte, destination interface{}) error {
	var wrapper struct {
		Data json.RawMessage `json:"Data"`
	}
	if err := json.Unmarshal(body, &wrapper); err != nil {
		return err
	}
	if len(wrapper.Data) > 0 && string(wrapper.Data) != "null" {
		return json.Unmarshal(wrapper.Data, destination)
	}
	return json.Unmarshal(body, destination)
}

func (c *Info_usuario_midController) writeError(status int, message string) {
	c.Ctx.Output.SetStatus(status)
	c.Data["json"] = map[string]string{"error": message}
	c.ServeJSON()
}
