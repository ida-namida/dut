package kopuro

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type kopuroClient struct {
	baseURL    string
	httpClient *http.Client
}

func newKopuroClient(baseURL string) kopuroClient {
	httpClient := &http.Client{}
	return kopuroClient{baseURL, httpClient}
}

func (kc kopuroClient) readJSONFile(filename string, response interface{}) error {
	request, err := http.NewRequest("GET", fmt.Sprintf("%s/json/read?filename=%s", kc.baseURL, filename), nil)
	if err != nil {
		return err
	}

	resp, err := kc.httpClient.Do(request)
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusOK {
		return json.NewDecoder(resp.Body).Decode(response)
	}

	errorMap := make(map[string]interface{})
	err = json.NewDecoder(resp.Body).Decode(&errorMap)
	if err != nil {
		return err
	}

	errorResponse, ok := errorMap["error"]
	if !ok {
		return fmt.Errorf("kopuro returned http status code %v with no error description", resp.StatusCode)
	}

	return fmt.Errorf("http %s from kopuro - %s", resp.Status, errorResponse)	
}