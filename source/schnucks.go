package source

import "fmt"

type Schnucks struct{}

func (f *Schnucks) ParseResponse(response map[string]interface{}) (Response, error) {
	var items []Item
	if response["data"] == nil {
		return items, nil
	}
	if _, ok := response["data"].([]interface{}); !ok {
		return nil, fmt.Errorf("data is not an array")
	}
	for _, item := range response["data"].([]interface{}) {
		// Parse item
		itemMap := item.(map[string]interface{})
		parsedItem := Item{
			Name: itemMap["name"].(string),
		}
		items = append(items, parsedItem)
	}
	return items, nil
}
