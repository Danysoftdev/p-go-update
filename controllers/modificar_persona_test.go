package controllers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/danysoftdev/p-go-update/controllers"
	"github.com/danysoftdev/p-go-update/models"
	"github.com/danysoftdev/p-go-update/services"
	"github.com/danysoftdev/p-go-update/tests/mocks"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestActualizarPersonaController_Success(t *testing.T) {
	mockRepo := new(mocks.MockPersonaRepo)
	services.SetPersonaRepository(mockRepo)

	persona := models.Persona{
		Documento: "123",
		Nombre:    "Juan",
		Apellido:  "Pérez",
		Edad:      30,
		Correo:    "juan@example.com",
		Telefono:  "1234567890",
		Direccion: "Calle Falsa 123",
	}

	// Mock de verificación de existencia y luego actualización
	mockRepo.On("ObtenerPersonaPorDocumento", "123").Return(persona, nil)
	mockRepo.On("ActualizarPersona", "123", persona).Return(nil)

	body, _ := json.Marshal(persona)
	req := httptest.NewRequest("PUT", "/personas/123", bytes.NewBuffer(body))
	req = mux.SetURLVars(req, map[string]string{"documento": "123"})
	rr := httptest.NewRecorder()

	controllers.ActualizarPersona(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), "actualizada exitosamente")

	mockRepo.AssertExpectations(t)
}

func TestActualizarPersonaController_ErrorFormato(t *testing.T) {
	req := httptest.NewRequest("PUT", "/personas/123", bytes.NewBuffer([]byte("invalido")))
	req = mux.SetURLVars(req, map[string]string{"documento": "123"})
	rr := httptest.NewRecorder()

	controllers.ActualizarPersona(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "formato del cuerpo es inválido")
}

func TestActualizarPersonaController_ErrorServicio(t *testing.T) {
	mockRepo := new(mocks.MockPersonaRepo)
	services.SetPersonaRepository(mockRepo)

	persona := models.Persona{
		Documento: "123",
		Nombre:    "Juan",
		Apellido:  "Pérez",
		Edad:      30,
		Correo:    "juan@example.com",
		Telefono:  "1234567890",
		Direccion: "Calle Falsa 123",
	}

	mockRepo.On("ObtenerPersonaPorDocumento", "123").Return(persona, nil)
	mockRepo.On("ActualizarPersona", "123", persona).Return(errors.New("fallo actualización"))

	body, _ := json.Marshal(persona)
	req := httptest.NewRequest("PUT", "/personas/123", bytes.NewBuffer(body))
	req = mux.SetURLVars(req, map[string]string{"documento": "123"})
	rr := httptest.NewRecorder()

	controllers.ActualizarPersona(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "fallo actualización")

	mockRepo.AssertExpectations(t)
}

func TestActualizarPersonaController_ErrorDocumentoDistinto(t *testing.T) {
	mockRepo := new(mocks.MockPersonaRepo)
	services.SetPersonaRepository(mockRepo)

	persona := models.Persona{
		Documento: "456", // diferente al de la ruta
		Nombre:    "María",
		Apellido:  "López",
		Edad:      35,
		Correo:    "maria@example.com",
		Telefono:  "3214567890",
		Direccion: "Carrera Falsa 456",
	}

	// Esto no se llega a ejecutar, pero por orden es buena práctica
	mockRepo.On("ObtenerPersonaPorDocumento", "123").Return(models.Persona{}, nil)

	body, _ := json.Marshal(persona)
	req := httptest.NewRequest("PUT", "/personas/123", bytes.NewBuffer(body))
	req = mux.SetURLVars(req, map[string]string{"documento": "123"})
	rr := httptest.NewRecorder()

	controllers.ActualizarPersona(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "no se puede modificar el documento")
}
