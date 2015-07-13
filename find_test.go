package nodefinder

import "testing"
import "os"
import "golang.org/x/net/html"

var testfile = "test.html"

func printNode(nodes []*html.Node, t *testing.T) {
	for _, n := range nodes {
		t.Logf("%#V\n", n)
	}
}

func Test_Find_1(t *testing.T) {
	hf, _ := os.Open(testfile)
	defer hf.Close()
	path := "ul[class='list3 list2']//li:1"
	tags := NewPath(path)
	nodes, _ := Find(tags, hf)
	if len(nodes) != 1 {
		printNode(nodes, t)
		t.Error("Failed")
	}
}

func Test_Find_2(t *testing.T) {
	hf, _ := os.Open(testfile)
	defer hf.Close()
	path := "html:1//.table//.table2"
	tags := NewPath(path)
	nodes, _ := Find(tags, hf)
	if len(nodes) != 2 {
		for _, n := range nodes {
			t.Log(n.Data)
		}
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

func Test_Find_4(t *testing.T) {
	hf, _ := os.Open(testfile)
	defer hf.Close()
	path := "//div.table/.."
	tags := NewPath(path)
	t.Logf("%#V\n", tags)
	nodes, _ := Find(tags, hf)
	t.Logf("%#V\n", nodes)
	if len(nodes) != 1 {
		t.Error("Failed")
	}
}
