package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/danysoftdev/p-go-update/models"
	"github.com/danysoftdev/p-go-update/services"

	"github.com/gorilla/mux"
)

func ActualizarPersona(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	documento := params["documento"]

	var persona models.Persona
	err := json.NewDecoder(r.Body).Decode(&persona)
	if err != nil {
		http.Error(w, "El formato del cuerpo es inválido", http.StatusBadRequest)
		return
	}

	err = services.ModificarPersona(documento, persona)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"mensaje": "Persona actualizada exitosamente"})
}
