package nodefinder

import "testing"
import "os"

var testfile = "test.html"

func Test_Find_1(t *testing.T) {
	hf, _ := os.Open(testfile)
	defer hf.Close()
	path := "ul[class='list3 list2']//li:1"
	tags := NewPath(path)
	//fmt.Printf("%#V\n", tags)
	nodes, _ := Find(tags, hf)
	if len(nodes) != 1 {
		t.Error("Failed")
	}
}

func Test_Find_2(t *testing.T) {
	hf, _ := os.Open(testfile)
	defer hf.Close()
	path := "html:1//.table//.table2"
	tags := NewPath(path)
	t.Logf("%#V\n", tags)
	nodes, _ := Find(tags, hf)
	if len(nodes) != 2 {
		t.Error("Failed")
	}
}

func Test_Find_3(t *testing.T) {
	hf, _ := os.Open(testfile)
	defer hf.Close()
	path := "//:1/:2//a"
	tags := NewPath(path)
	t.Logf("%#V\n", tags)
	nodes, _ := Find(tags, hf)
	t.Logf("%#V\n", nodes)
	if len(nodes) != 2 {
		t.Error("Failed")
	}
}
