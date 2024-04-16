package source

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSource_BuildURL(t *testing.T) {
	s := &Source{
		Url: "https://api.freshop.com/1/products?app_key=%s&fields=%s&limit=%d&q=%s&relevance_sort=%s&sort=%s&store_id=%d&token=%s",
		RequestFields: map[string]interface{}{
			"app_key":        "fest_foods",
			"fields":         "id,name,department_id,status,shopper_location",
			"limit":          24.0,
			"relevance_sort": "relevance_sort",
			"sort":           "sort",
			"store_id":       1955,
			"token":          "123456",
		},
		QueryField: "q",
	}
	expected := "https://api.freshop.com/1/products?app_key=fest_foods&fields=id,name,department_id,status,shopper_location&limit=24&q=oranges&relevance_sort=relevance_sort&sort=sort&store_id=1955&token=123456"
	url, err := s.BuildURL("oranges")
	require.Nil(t, err)
	require.Equal(t, expected, url)
}
