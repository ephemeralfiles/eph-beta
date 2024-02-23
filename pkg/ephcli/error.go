package ephcli

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// parseError is a helper function to parse the error from the response
func parseError(resp *http.Response) error {
	var jsonResponse ResponseError
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(respBody, &jsonResponse)
	if err != nil {
		return err
	}
	return fmt.Errorf("status not ok %d: %s", resp.StatusCode, jsonResponse.Message)
}
