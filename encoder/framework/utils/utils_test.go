package utils_test

import (
	"github.com/diasYuri/encoder-go/framework/utils"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestIsJson(t *testing.T) {
	json := `{
				"_id": "6091713b6a00f546c5f59141",
				"createdOn": "2021-05-04T16:07:23.146Z",
				"name": "Home"
			  }`

	err := utils.IsJson(json)
	require.Nil(t, err)


	json = `wes`
	err = utils.IsJson(json)
	require.Error(t, err)

}
