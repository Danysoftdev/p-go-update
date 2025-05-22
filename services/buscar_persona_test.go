package services_test

import (
	"errors"
	"testing"

	"github.com/danysoftdev/p-go-update/models"
	"github.com/danysoftdev/p-go-update/services"
	"github.com/danysoftdev/p-go-update/tests/mocks"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestBuscarPersonaPorDocumento_Exito(t *testing.T) {
	mockRepo := new(mocks.MockPersonaRepo)
	services.SetPersonaRepository(mockRepo)

	personaMock := models.Persona{
		Documento: "123",
		Nombre:    "Juan",
		Apellido:  "Pérez",
		Edad:      30,
		Correo:    "juan@example.com",
		Telefono:  "1234567890",
		Direccion: "Calle 123",
	}

	mockRepo.On("ObtenerPersonaPorDocumento", "123").Return(personaMock, nil)

	persona, err := services.BuscarPersonaPorDocumento("123")

	assert.Nil(t, err)
	assert.Equal(t, "Juan", persona.Nombre)
	mockRepo.AssertExpectations(t)
}

func TestBuscarPersonaPorDocumento_Vacio(t *testing.T) {
	mockRepo := new(mocks.MockPersonaRepo)
	services.SetPersonaRepository(mockRepo)

	_, err := services.BuscarPersonaPorDocumento("")

	assert.NotNil(t, err)
	assert.Equal(t, "el documento no puede estar vacío", err.Error())
}

func TestBuscarPersonaPorDocumento_NoEncontrado(t *testing.T) {
	mockRepo := new(mocks.MockPersonaRepo)
	services.SetPersonaRepository(mockRepo)

	mockRepo.On("ObtenerPersonaPorDocumento", "999").Return(models.Persona{}, mongo.ErrNoDocuments)

	_, err := services.BuscarPersonaPorDocumento("999")

	assert.NotNil(t, err)
	assert.Equal(t, "persona no encontrada", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestBuscarPersonaPorDocumento_ErrorBaseDeDatos(t *testing.T) {
	mockRepo := new(mocks.MockPersonaRepo)
	services.SetPersonaRepository(mockRepo)

	mockRepo.On("ObtenerPersonaPorDocumento", "123").Return(models.Persona{}, errors.New("error de base de datos"))

	_, err := services.BuscarPersonaPorDocumento("123")

	assert.NotNil(t, err)
	assert.Equal(t, "error de base de datos", err.Error())
	mockRepo.AssertExpectations(t)
}
