package main

import (
	"encoding/json"
	"net/http"
)

func (a *applicationDependencies) writeJSON(w http.ResponseWriter,
	status int,
	data any,
	headers http.Header) error {
	jsResponse, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	jsResponse = append(jsResponse, '\n')
	//additonal headers to be set
	for key, value := range headers {
		w.Header()[key] = value
		//w.Header().Set(key, value[0])
	}

	//set content type header
	w.Header().Set("Content-Type", "application/json")
	//explicity set the response status code
	w.WriteHeader(status)
	_, err = w.Write(jsResponse)
	if err != nil {
		return err
	}
	return nil
	//w.Write(jsResponse)

	//return nil
}
