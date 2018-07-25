package goossselect

import(
	"testing"
	"fmt"
)

func TestSelectQuery(t *testing.T) {

	result, err := SelectQuery("objselect", "test.csv", "Select _1 from ossobject where _3 > 1")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%v", string(result))
	t.Logf("%v", string(result))
}
