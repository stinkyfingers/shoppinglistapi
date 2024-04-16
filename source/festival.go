package source

import "fmt"

type Festival struct{}

func (f *Festival) ParseResponse(response map[string]interface{}) (Response, error) {
	var items []Item
	if response["items"] == nil {
		return items, nil
	}
	if _, ok := response["items"].([]interface{}); !ok {
		return nil, fmt.Errorf("items is not an array")
	}
	for _, item := range response["items"].([]interface{}) {
		// Parse item
		itemMap := item.(map[string]interface{})
		var parsedItem Item
		if val, ok := itemMap["id"].(string); ok {
			parsedItem.ID = val
		}
		if val, ok := itemMap["name"].(string); ok {
			parsedItem.Name = val
		}
		if val, ok := itemMap["status"].(string); ok {
			parsedItem.Status = val
		}
		if val, ok := itemMap["shopper_location"].(string); ok {
			parsedItem.Location = val
		}
		if val, ok := itemMap["department_id"].([]string); ok {
			parsedItem.DeparmentIDs = val
		}
		items = append(items, parsedItem)
	}
	return items, nil
}
