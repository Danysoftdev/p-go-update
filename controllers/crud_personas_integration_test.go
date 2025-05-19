//go:build integration
// +build integration

package controllers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	"github.com/danysoftdev/p-go-update/config"
	"github.com/danysoftdev/p-go-update/controllers"
	"github.com/danysoftdev/p-go-update/models"
	"github.com/danysoftdev/p-go-update/repositories"
	"github.com/danysoftdev/p-go-update/services"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.mongodb.org/mongo-driver/bson"
)

func TestEndpointsControllerIntegration(t *testing.T) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "mongo:6.0",
		ExposedPorts: []string{"27017/tcp"},
		WaitingFor:   wait.ForListeningPort("27017/tcp").WithStartupTimeout(20 * time.Second),
	}
	mongoC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	assert.NoError(t, err)
	defer mongoC.Terminate(ctx)

	endpoint, err := mongoC.Endpoint(ctx, "")
	assert.NoError(t, err)

	os.Setenv("MONGO_URI", "mongodb://"+endpoint)
	os.Setenv("MONGO_DB", "testdb")
	os.Setenv("COLLECTION_NAME", "personas_test")

	err = config.ConectarMongo()
	assert.NoError(t, err)
	defer config.CerrarMongo()

	repositories.SetCollection(config.Collection)
	services.SetPersonaRepository(repositories.RealPersonaRepository{})

	// 3. Limpiar colección
	_, err = config.Collection.DeleteMany(ctx, bson.M{})
	assert.NoError(t, err)

	// 4. Insertar persona
	persona := models.Persona{
		Documento: "999",
		Nombre:    "Lucía",
		Apellido:  "Pérez",
		Edad:      30,
		Correo:    "lucia@example.com",
		Telefono:  "3111234567",
		Direccion: "Calle 123",
	}
	_, err = config.Collection.InsertOne(ctx, persona)
	assert.NoError(t, err)

	// Setup router
	router := mux.NewRouter()
	router.HandleFunc("/personas/{documento}", controllers.ActualizarPersona).Methods("PUT")

	// 4. Modificar persona
	persona.Nombre = "Actualizado"
	bodyUpdate, _ := json.Marshal(persona)
	reqUpdate := httptest.NewRequest("PUT", "/personas/999", bytes.NewReader(bodyUpdate))
	reqUpdate = mux.SetURLVars(reqUpdate, map[string]string{"documento": "999"})
	resUpdate := httptest.NewRecorder()
	router.ServeHTTP(resUpdate, reqUpdate)

	assert.Equal(t, http.StatusOK, resUpdate.Code)

}
