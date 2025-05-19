package repositories

import "github.com/danysoftdev/p-go-update/models"

type PersonaRepository interface {

	ObtenerPersonaPorDocumento(documento string) (models.Persona, error)
	ActualizarPersona(documento string, persona models.Persona) error
	
}
