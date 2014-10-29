package nodefinder

import "testing"
import "os"
import "fmt"

var testfile = "test.html"

func Test_Find_1(t *testing.T) {
	hf, _ := os.Open(testfile)
	defer hf.Close()
	path := "ul[class='list3 list2']//li:1"
	tags := NewPath(path)
	fmt.Println(tags)
	//fmt.Printf("%#V\n", tags)
	nodes, _ := Find(tags, hf)
	if len(nodes) != 1 {
		t.Error("Failed")
	}
}
