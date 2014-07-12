package nodefinder

import "testing"
import "os"

var testfile = "test.html"

func Test_Find_1(t *testing.T) {
	hf, _ := os.Open(testfile)
	defer hf.Close()
	path := "table.downloads/tbody/tr/td/a"
	tags := TagParser(path)
	nodes, _ := Find(tags, hf)
	if len(nodes) != 11 {
		t.Logf("%#V\n", nodes)
		t.Error("Failed")
	}
}

func Test_Find_2(t *testing.T) {
	hf, _ := os.Open(testfile)
	defer hf.Close()
	path := "table.downloads/tbody/tr/td/code"
	tags := TagParser(path)
	nodes, _ := Find(tags, hf)
	if len(nodes) != 11 {
		t.Logf("%#V\n", nodes)
		t.Error("Failed")
	}
}
