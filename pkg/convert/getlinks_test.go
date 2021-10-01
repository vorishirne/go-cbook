package convert

import (
	"fmt"
	"testing"
)

func TestMod_GetFilePath(t *testing.T) {
	m := Mod{
		CssFile: "",
		// here it is compulsary for scheme(https/http) to be present
		// otherwise net package wouldn't work.
		// hence, no need for exclusive normalization for scheme
		BaseUrl:    "https://www.medium.com/.*",
		BaseDir:    "envoy-blogs",
		dirVisited: map[string]*DirVisited{},
	}
	filePath, fileName, err := m.GetRawFilePath("https://medium.com/.*/galgodas/@abhbose6/bazel-101-2b0272b15da8/")
	if err != nil {
		t.Fatal(err)
	}
	if filePath != "galgodas/@abhbose6" {
		t.Fatalf("path %v not expected", filePath)
	}
	if fileName != "bazel-101-2b0272b15da8" {
		t.Fatalf("name %v not expected", fileName)
	}
	dir, index, err := m.EnsureNumberedDir(filePath)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(dir, index)
}
