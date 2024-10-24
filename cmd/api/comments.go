package main

import (
	"fmt"
	"net/http"

	"github.com/luigiacunaUB/comments/internal/data"
	// import the data package which contains the definition for Comment
	//"github.com/luigiacunaUB/comments"
	"github.com/luigiacunaUB/comments/internal/validator"
)

func (a *applicationDependencies) createCommentHandler(w http.ResponseWriter, r *http.Request) {
	// create a struct to hold a comment
	// we use struct tags[“] to make the names display in lowercase
	var incomingData struct {
		Content string `json:"content"`
		Author  string `json:"author"`
	}
	// perform the decoding
	//err := json.NewDecoder(r.Body).Decode(&incomingData)
	err := a.readJSON(w, r, &incomingData)
	if err != nil {
		//a.errorResponseJSON(w, r, http.StatusBadRequest, err.Error())
		a.badRequestResponse(w, r, err)
		return
	}
	// Copy the values from incomingData to a new Comment struct
	// At this point in our code the JSON is well-formed JSON so now
	// we will validate it using the Validator which expects a Comment
	comment := &data.Comment{
		Content: incomingData.Content,
		Author:  incomingData.Author,
	}

	//Initialize a Validator instance
	v := validator.New()

	//Do the Validation
	data.ValidateComment(v, comment)
	if !v.IsEmpty() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	//Add the comment to the database table
	err = a.commentModel.Insert(comment)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/comments/%d", comment.ID))

	//Send a JSON response with 201 (new resource created) status code
	data := envelope{
		"comment": comment,
	}
	err = a.writeJSON(w, http.StatusCreated, data, headers)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	// for now display the result
	//fmt.Fprintf(w, "%+v\n", incomingData)
}
