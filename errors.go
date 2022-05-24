package jenga

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type JengaError struct {
	Status  bool   `json:"status"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// ErrorHandler takes a http response as an input and returns the appropriate error
//  given the response returned from Jenga
func ErrorHandler(resp *http.Response) error {
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error response body could not be read: %w", err)
	}

	var je JengaError
	err = json.Unmarshal(body, &je)
	if err != nil {
		return fmt.Errorf("error response body could not be unmarshalled: %w", err)
	}

	return fmt.Errorf(
		"jenga api call failed [status-code: %d] - %s (%d)", je.Code, je.Message, resp.StatusCode)
}
