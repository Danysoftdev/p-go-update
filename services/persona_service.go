package services

import (
	"errors"
	"strings"

	"github.com/danysoftdev/p-go-update/models"
	"github.com/danysoftdev/p-go-update/repositories"
	"go.mongodb.org/mongo-driver/mongo"
)

var Repo repositories.PersonaRepository

func SetPersonaRepository(r repositories.PersonaRepository) {
	Repo = r
}

func ValidarPersona(p models.Persona) error {
	if strings.TrimSpace(p.Documento) == "" {
		return errors.New("el documento no puede estar vacío")
	}
	if strings.TrimSpace(p.Nombre) == "" {
		return errors.New("el nombre no puede estar vacío")
	}
	if strings.TrimSpace(p.Apellido) == "" {
		return errors.New("el apellido no puede estar vacío")
	}
	if p.Edad <= 0 {
		return errors.New("la edad debe ser un número entero mayor a 0")
	}
	if strings.TrimSpace(p.Correo) == "" || !strings.Contains(p.Correo, "@") {
		return errors.New("el correo es inválido")
	}
	if strings.TrimSpace(p.Telefono) == "" {
		return errors.New("el teléfono no puede estar vacío")
	}
	if strings.TrimSpace(p.Direccion) == "" {
		return errors.New("la dirección no puede estar vacía")
	}
	return nil
}


func BuscarPersonaPorDocumento(doc string) (models.Persona, error) {
	if strings.TrimSpace(doc) == "" {
		return models.Persona{}, errors.New("el documento no puede estar vacío")
	}

	persona, err := Repo.ObtenerPersonaPorDocumento(doc)
	if err == mongo.ErrNoDocuments {
		return models.Persona{}, errors.New("persona no encontrada")
	}

	return persona, err
}

func ModificarPersona(documento string, p models.Persona) error {
	if strings.TrimSpace(documento) == "" {
		return errors.New("el documento no puede estar vacío")
	}

	if err := ValidarPersona(p); err != nil {
		return err
	}

	_, err := Repo.ObtenerPersonaPorDocumento(documento)
	if err == mongo.ErrNoDocuments {
		return errors.New("persona no encontrada")
	}

	if p.Documento != documento {
		return errors.New("no se puede modificar el documento de una persona")
	}

	return Repo.ActualizarPersona(documento, p)
}
