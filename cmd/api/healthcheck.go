package main

import (
	"net/http"
)

func healthcheckHandler(app *application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := envelope{
			"status": "available",
			"system_info": map[string]string{
				"environment": app.config.env,
				"version":     "1.0.0",
			},
		}

		err := app.writeJSON(w, http.StatusOK, data, nil)
		if err != nil {
			app.serverErrorResponse(w, r, err)
		}
	}

}
