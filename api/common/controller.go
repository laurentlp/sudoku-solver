package common

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/laurentlp/sudoku-solver/api/errors"
)

// Controller handle all base methods
type Controller struct {
}

// SendJSON marshals v to a json struct and sends appropriate headers to w
func (c *Controller) SendJSON(w http.ResponseWriter, r *http.Request, v interface{}, code int) {
	w.Header().Add("Content-Type", "application/json")

	j, err := json.Marshal(v)

	if err != nil {
		log.Print(fmt.Sprintf("Error while encoding JSON: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, `{"error": "Internal server error"}`)
	} else {
		w.WriteHeader(code)
		io.WriteString(w, string(j))
	}
}

// MapJSON marshals v to a json struct
// Return true if succesful, false otherwise
func (c *Controller) MapJSON(w http.ResponseWriter, r *http.Request, v interface{}) *errors.APIError {
	bodyBuffer, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Println(fmt.Sprintf("Error while reading response body: %v", err))

		errors.BadRequest("The informations sent to the server contains errors.")
	} else if len(bodyBuffer) == 0 {
		log.Println("Body is empty:")

		return errors.BadRequest("No information was sent to the server. Please send a valid sudoku.")
	}

	if err := json.Unmarshal(bodyBuffer, &v); err != nil {
		log.Println(fmt.Sprintf("Error while encoding JSON: %v", err))

		return errors.BadRequest("The informations sent to the server contains errors.")
	}
	return nil
}
