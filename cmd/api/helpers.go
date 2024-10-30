package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/luigiacunaUB/comments/internal/validator"

	"github.com/julienschmidt/httprouter"
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
	err := json.NewDecoder(r.Body).Decode(destination)
	if err != nil {
		//check for the diffrent errors
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

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

		case errors.As(err, &invalidUnmarshalError):
			panic(err)
		default:
			return err
		}

	}
	return nil
}

func (a *applicationDependencies) readIDParam(r *http.Request) (int64, error) {
	// Get the URL parameters
	params := httprouter.ParamsFromContext(r.Context())
	// Convert the id from string to int
	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}

	return id, nil

}

func (a *applicationDependencies) getSingleQueryParameters(queryParameters url.Values, key string, defaultValue string) string {
	// url.Values is a key:value hash map of the query parameters
	result := queryParameters.Get(key)
	if result == "" {
		return defaultValue
	}
	return result

}

// call when we have multiple comma-separated values
func (a *applicationDependencies) getMultipleQueryParameters(queryParameters url.Values, key string, defaultValue []string) []string {
	result := queryParameters.Get(key)
	if result == "" {
		return defaultValue
	}
	return strings.Split(result, ",")

}

func (a *applicationDependencies) getSingleIntegerParameter(queryParameters url.Values, key string, defaultValue int, v *validator.Validator) int {
	result := queryParameters.Get(key)
	if result == "" {
		return defaultValue
	}
	// try to convert to an integer
	intValue, err := strconv.Atoi(result)
	if err != nil {
		v.AddError(key, "must be an integer value")
		return defaultValue
	}

	return intValue

}
