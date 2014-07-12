package nodefinder

import "testing"

func Test_TagParser_1(t *testing.T) {
	sample := `body/div.good[id="a b",code="go"]/#gogogo/p.download[os=linux,arch=x86_64]/a[href="http://meow.com/",target=_blank]`
	tags := TagParser(sample)
	if tags[0].Tag != "body" {
		t.Logf("%#V\n", tags[0])
		t.Error("Failed")
	}
	if tags[1].Tag != "div" {
		t.Logf("%#V\n", tags[1])
		t.Error("Failed")
	}
	if _, ok := tags[1].Attr["class"]["good"]; !ok {
		t.Logf("%#V\n", tags[1])
		t.Error("Failed")
	}
	if _, ok := tags[1].Attr["id"]["a"]; !ok {
		t.Logf("%#V\n", tags[1])
		t.Error("Failed")
	}
	if _, ok := tags[1].Attr["id"]["b"]; !ok {
		t.Logf("%#V\n", tags[1])
		t.Error("Failed")
	}
	if tags[2].Tag != "" {
		t.Logf("%#V\n", tags[2])
		t.Error("Failed")
	}
	if tags[3].Tag != "p" {
		t.Logf("%#V\n", tags[3])
		t.Error("Failed")
	}
	if _, ok := tags[3].Attr["class"]["download"]; !ok {
		t.Logf("%#V\n", tags[3])
		t.Error("Failed")
	}
	if _, ok := tags[3].Attr["os"]["linux"]; !ok {
		t.Logf("%#V\n", tags[3])
		t.Error("Failed")
	}
	if tags[4].Tag != "a" {
		t.Logf("%#V\n", tags[3])
		t.Error("Failed")
	}
	if _, ok := tags[4].Attr["href"]["http://meow.com/"]; !ok {
		t.Logf("%#V\n", tags[4])
		t.Error("Failed")
	}
}
