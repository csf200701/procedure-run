package asset

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"strings"
)

func bindata_read(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	return buf.Bytes(), nil
}

var _file_mouse_txt = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x3c\xcd\xd1\xa9\xc3\x30\x0c\x85\xe1\xf7\x0b\x77\x87\x9f\xbc\xd8\x82\x4a\x5a\xa0\xe9\x24\x26\x27\x1b\x64\x80\xa6\x9d\xbd\xb8\xb5\x73\x1e\x8c\xf5\x49\x70\x20\xdc\x6f\xfc\xd2\xbf\xff\x7f\x50\xa9\xd0\xc2\xdd\x23\xc1\xb0\x8e\x94\x90\xf2\x00\x8e\x26\x45\xf9\x1a\xf0\x5c\x61\x83\xf5\x35\x81\x07\x38\xdc\xe7\x9c\x83\xdb\x80\xbc\x60\x48\xe6\xd9\xcb\x81\x73\xd2\x52\xda\x38\xca\xb2\x28\xfc\xbd\xed\xc5\xe3\x6a\x68\x20\x48\xf7\xc2\xcc\x3e\x96\x92\xac\x62\x55\xd2\x90\x5a\xab\x64\xf4\xc7\x0c\xf8\x04\x00\x00\xff\xff\xe5\xfd\x15\xa8\xf3\x00\x00\x00")

func file_mouse_txt() ([]byte, error) {
	return bindata_read(
		_file_mouse_txt,
		"file/mouse.txt",
	)
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		return f()
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() ([]byte, error){
	"file/mouse.txt": file_mouse_txt,
}
// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for name := range node.Children {
		rv = append(rv, name)
	}
	return rv, nil
}

type _bintree_t struct {
	Func func() ([]byte, error)
	Children map[string]*_bintree_t
}
var _bintree = &_bintree_t{nil, map[string]*_bintree_t{
	"file/mouse.txt": &_bintree_t{file_mouse_txt, map[string]*_bintree_t{
	}},
}}
