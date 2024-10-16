package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// creating an envelope type
type envelope map[string]any

func (a *applicationDependencies) writeJSON(w http.ResponseWriter,
	status int,
	data envelope,
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

func (a *applicationDependencies) readJSON(w http.ResponseWriter, r *http.Request, destination any) error {
	//what is the max size of the request body(250K is reasonable)
	maxBytes := 256_000
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))
	//our decoder will check for unknown fields
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	//let start the decoding
	err := dec.Decode(destination)
	//err := json.NewDecoder(r.Body).Decode(destination)
	if err != nil {
		//check for the diffrent errors
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError
		var MaxBytesError *http.MaxBytesError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("the body contains badly-formed JSON (at charater %d)", syntaxError.Offset)

			//Decode can also sed back an io error message
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("the body contains badly-formed JSON")

		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("the body contains the incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("the body contains the incorrect JSON type (at charater %d)", unmarshalTypeError.Offset)

		case errors.Is(err, io.EOF):
			return errors.New("the body must not be empty")

		case strings.HasPrefix(err.Error(), "json:unknown field"):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field")
			return fmt.Errorf("body contains unknown key %s", fieldName)

		//does the body exceed our limit of 250KB
		case errors.As(err, &MaxBytesError):
			return fmt.Errorf("The body must not be larger that %d bytes", MaxBytesError.Limit)

		case errors.As(err, &invalidUnmarshalError):
			panic(err)
		default:
			return err
		}

	}
	// almost done. Let's lastly check if there is any data after
	// the valid JSON data. Maybe the person is trying to send
	// multiple request bodies during one request
	// We call decode once more to see if it gives us by anything
	// we use a throw away struct 'struct{}{}' to hold the result

	err = dec.Decode(&struct{}{})
	if !errors.Is(err, io.EOF) {
		return errors.New("the body must only contain a single JSON value")
	}
	return nil
}
