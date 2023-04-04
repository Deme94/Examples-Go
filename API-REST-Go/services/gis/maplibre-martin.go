package gis

import (
	"API-REST/services/conf"
	"errors"
	"fmt"
	"net/http"
)

func Setup() error {
	// Read conf
	host := conf.Env.GetString("MARTIN_HOST")
	port := conf.Env.GetString("MARTIN_PORT")

	url := host + ":" + port + "/health"

	// Hacer la solicitud HTTP
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Verificar el código de estado
	if resp.StatusCode != http.StatusOK {
		return errors.New("error " + fmt.Sprint(resp.StatusCode))
	}
	return nil
}
