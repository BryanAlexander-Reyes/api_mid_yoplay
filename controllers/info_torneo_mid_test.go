package controllers

import (
	"encoding/json"
	"testing"
)

func TestTorneoInfoAcceptsImageID(t *testing.T) {
	var torneo TorneoInfo
	if err := json.Unmarshal([]byte(`{"IdImagen":12}`), &torneo); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}
	if torneo.Imagen == nil {
		t.Fatal("expected image reference")
	}
	if torneo.Imagen.ID != 12 {
		t.Fatalf("image ID = %d, want 12", torneo.Imagen.ID)
	}
}

func TestTorneoInfoAcceptsImageObject(t *testing.T) {
	const expectedURL = "https://cdn.example.com/torneo.png"
	data := []byte(`{"IdImagen":{"Id":12,"UrlImagen":"` + expectedURL + `"}}`)

	var torneo TorneoInfo
	if err := json.Unmarshal(data, &torneo); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}
	if torneo.Imagen == nil || torneo.Imagen.ID != 12 {
		t.Fatalf("image reference = %#v, want ID 12", torneo.Imagen)
	}
	if torneo.Imagen.URLImagen != expectedURL {
		t.Fatalf("image URL = %q, want %q", torneo.Imagen.URLImagen, expectedURL)
	}
}

func TestRespuestaTournamentsIncludesImageURL(t *testing.T) {
	const expectedURL = "https://www.pexels.com/es-es/foto/tecnologia-musica-negro-diseno-27507165/"
	body, err := json.Marshal(RespuestaTournaments{Imagen: expectedURL})
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}
	if response["Imagen"] != expectedURL {
		t.Fatalf("Imagen = %v, want %q", response["Imagen"], expectedURL)
	}
}
