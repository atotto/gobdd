package gobdd

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func TestGen(t *testing.T) {
	buf, err := Gen("_testdata/features/addition.feature")
	fmt.Println(string(buf))
	if err != nil {
		t.Fatal(err)
	}
	ioutil.WriteFile("_testdata/glue/addition_test.go", buf, 0644)
}
