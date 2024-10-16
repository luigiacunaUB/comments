package main

import (
	"net/http"
)

func (a *applicationDependencies) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	/*fmt.Fprintln(w, "status: available")
	fmt.Fprintf(w, "environment: %s\n", a.config.environment)
	fmt.Fprintf(w, "version: %s\n", appVersion)*/

	/*jsResponse := `{"status": "avaliable", "enviroment": %q,"version": %q}`
	jsResponse = fmt.Sprintf(jsResponse, a.config.environment, appVersion)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(jsResponse))*/

	//panic("Apples & Oranges")
	data := envelope{
		"status": "avalibale",
		"system_info": map[string]string{
			"enviroment": a.config.environment,
			"version":    appVersion,
		},
	}

	//jsResponse, err := json.Marshal(data)
	err := a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		//a.logger.Error(err.Error())
		a.serverErrorResponse(w, r, err)
		http.Error(w, "The server encounted a problem and could not process your request", http.StatusInternalServerError)
		return
	}

	/*jsResponse = append(jsResponse, '\n')
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsResponse)*/
}
