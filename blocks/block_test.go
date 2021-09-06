package blocks

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/kataras/iris/v12/view"
	"github.com/stretchr/testify/assert"
)

func Test_Default_Layout(t *testing.T) {
	assert := assert.New(t)
	buf := new(bytes.Buffer)
	blockEngine := NewBlocks("../views", ".html")

	blockEngine.Layout("default")
	err := blockEngine.Load()
	assert.Nil(err)

	data := map[string]interface{}{
		"Title": "Vietnam",
		"user":  "Donald Trump",
		"email": "trump@whitehouse.gov.us",
	}

	err = blockEngine.ExecuteWriter(buf, "test", "", data)
	assert.Nil(err)

	body := buf.String()
	assert.Contains(body, "Donald Trump")
	assert.Contains(body, "Vietnam")
}

func Test_No_Default_Layout(t *testing.T) {
	assert := assert.New(t)
	buf := new(bytes.Buffer)
	blockEngine := NewBlocks("../views", ".html")

	err := blockEngine.Load()
	assert.Nil(err)

	data := map[string]interface{}{
		"user":  "Donald Trump",
		"email": "trump@whitehouse.gov.us",
	}

	err = blockEngine.ExecuteWriter(buf, "test", "test_layout", data)
	assert.Nil(err)

	body := buf.String()
	assert.Contains(body, "Test Layout") //views/layouts/test_layout.html contains "Test Layout"
	assert.Contains(body, "Donald Trump")
	assert.Contains(body, "trump@whitehouse.gov.us")
}

func Test_No_Layout_At_All(t *testing.T) {
	assert := assert.New(t)
	buf := new(bytes.Buffer)
	blockEngine := NewBlocks("../views", ".html")

	err := blockEngine.Load()
	assert.Nil(err)

	data := map[string]interface{}{
		"user":  "Donald Trump",
		"email": "trump@whitehouse.gov.us",
	}

	err = blockEngine.ExecuteWriter(buf, "test", view.NoLayout, data)
	assert.Nil(err)

	body := buf.String()
	fmt.Println(body)
	assert.Contains(body, "Donald Trump")
	assert.Contains(body, "trump@whitehouse.gov.us")

	//Vì không sử dụng layout nên sẽ không có những thẻ này
	assert.NotContains(body, "<body>")
	assert.NotContains(body, "</html>")
}
