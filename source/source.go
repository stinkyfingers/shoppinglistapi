package source

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

type Source struct {
	Name          string                 `json:"name"`
	Url           string                 `json:"url"`
	Authorization string                 `json:"authorization"`
	RequestFields map[string]interface{} `json:"request_fields"`
	QueryField    string                 `json:"query_field"`
	ResponseMap   map[string]string      `json:"response_map"`
}

type Response []Item

type Item struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Status       string   `json:"status"`
	Location     string   `json:"shopper_location"`
	DeparmentIDs []string `json:"department_id"`
}

type Results interface {
	ParseResponse(response map[string]interface{}) Response
}

// GetSource reads the source json file and returns the Source struct.
func GetSource(store string) (*Source, error) {
	filename := fmt.Sprintf("config/%s.json", store)
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var source Source
	if err = json.NewDecoder(file).Decode(&source); err != nil {
		return nil, err
	}
	return &source, nil
}

// BuildURL builds the url for the source. It matches the request fields with the url params and returns the url.
func (s *Source) BuildURL(query string) (string, error) {
	urlArray := strings.Split(s.Url, "?")
	if len(urlArray) != 2 {
		return "", fmt.Errorf("url is not in the correct format")
	}
	paramsFields := strings.Split(urlArray[1], "&")
	var paramValues []interface{}
	for _, field := range paramsFields {
		fieldValue := strings.Split(field, "=")
		if fieldValue[0] == s.QueryField {
			paramValues = append(paramValues, query)
			continue
		}
		paramValue, ok := s.RequestFields[fieldValue[0]]
		if !ok {
			return "", fmt.Errorf("field %s not found in request fields", fieldValue[0])
		}
		paramType := fieldValue[1]
		switch paramValue.(type) {
		case string:
			if paramType != "%s" {
				return "", fmt.Errorf("field %s is not of type string", fieldValue[0])
			}
		case int, int64, int32:
			if paramType == "%f" { // convert int to float
				paramValue = float64(paramValue.(int))
			}
			if paramType != "%d" && paramType != "%f" {
				return "", fmt.Errorf("field %s is not of type number", fieldValue[0])
			}
		case float64, float32:
			if paramType == "%d" { // convert float to int
				paramValue = int(paramValue.(float64))
			}
			if paramType != "%f" && paramType != "%d" {
				return "", fmt.Errorf("field %s is not of type float", fieldValue[0])
			}
		default:
			return "", fmt.Errorf("field %s is not of type string or int", fieldValue[0])
		}

		paramValues = append(paramValues, paramValue)
	}
	return fmt.Sprintf(s.Url, paramValues...), nil
}

func (s *Source) Search(url string) (map[string]interface{}, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	cli := &http.Client{}
	req.Header.Set("Authorization", s.Authorization)
	resp, err := cli.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var response map[string]interface{}
	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}
	return response, nil
}

func (s *Source) GetItems(response map[string]interface{}) ([]Item, error) {
	switch s.Name {
	case "festival":
		f := &Festival{}
		return f.ParseResponse(response)
	case "schnucks":
		s := &Schnucks{}
		return s.ParseResponse(response)
	default:
		return nil, fmt.Errorf("source not found")
	}
}
