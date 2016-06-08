package main

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"sync"
	"time"
)

type _escLocalFS struct{}

var _escLocal _escLocalFS

type _escStaticFS struct{}

var _escStatic _escStaticFS

type _escDirectory struct {
	fs   http.FileSystem
	name string
}

type _escFile struct {
	compressed string
	size       int64
	modtime    int64
	local      string
	isDir      bool

	once sync.Once
	data []byte
	name string
}

func (_escLocalFS) Open(name string) (http.File, error) {
	f, present := _escData[path.Clean(name)]
	if !present {
		return nil, os.ErrNotExist
	}
	return os.Open(f.local)
}

func (_escStaticFS) prepare(name string) (*_escFile, error) {
	f, present := _escData[path.Clean(name)]
	if !present {
		return nil, os.ErrNotExist
	}
	var err error
	f.once.Do(func() {
		f.name = path.Base(name)
		if f.size == 0 {
			return
		}
		var gr *gzip.Reader
		b64 := base64.NewDecoder(base64.StdEncoding, bytes.NewBufferString(f.compressed))
		gr, err = gzip.NewReader(b64)
		if err != nil {
			return
		}
		f.data, err = ioutil.ReadAll(gr)
	})
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (fs _escStaticFS) Open(name string) (http.File, error) {
	f, err := fs.prepare(name)
	if err != nil {
		return nil, err
	}
	return f.File()
}

func (dir _escDirectory) Open(name string) (http.File, error) {
	return dir.fs.Open(dir.name + name)
}

func (f *_escFile) File() (http.File, error) {
	type httpFile struct {
		*bytes.Reader
		*_escFile
	}
	return &httpFile{
		Reader:   bytes.NewReader(f.data),
		_escFile: f,
	}, nil
}

func (f *_escFile) Close() error {
	return nil
}

func (f *_escFile) Readdir(count int) ([]os.FileInfo, error) {
	return nil, nil
}

func (f *_escFile) Stat() (os.FileInfo, error) {
	return f, nil
}

func (f *_escFile) Name() string {
	return f.name
}

func (f *_escFile) Size() int64 {
	return f.size
}

func (f *_escFile) Mode() os.FileMode {
	return 0
}

func (f *_escFile) ModTime() time.Time {
	return time.Unix(f.modtime, 0)
}

func (f *_escFile) IsDir() bool {
	return f.isDir
}

func (f *_escFile) Sys() interface{} {
	return f
}

// FS returns a http.Filesystem for the embedded assets. If useLocal is true,
// the filesystem's contents are instead used.
func FS(useLocal bool) http.FileSystem {
	if useLocal {
		return _escLocal
	}
	return _escStatic
}

// Dir returns a http.Filesystem for the embedded assets on a given prefix dir.
// If useLocal is true, the filesystem's contents are instead used.
func Dir(useLocal bool, name string) http.FileSystem {
	if useLocal {
		return _escDirectory{fs: _escLocal, name: name}
	}
	return _escDirectory{fs: _escStatic, name: name}
}

// FSByte returns the named file from the embedded assets. If useLocal is
// true, the filesystem's contents are instead used.
func FSByte(useLocal bool, name string) ([]byte, error) {
	if useLocal {
		f, err := _escLocal.Open(name)
		if err != nil {
			return nil, err
		}
		return ioutil.ReadAll(f)
	}
	f, err := _escStatic.prepare(name)
	if err != nil {
		return nil, err
	}
	return f.data, nil
}

// FSMustByte is the same as FSByte, but panics if name is not present.
func FSMustByte(useLocal bool, name string) []byte {
	b, err := FSByte(useLocal, name)
	if err != nil {
		panic(err)
	}
	return b
}

// FSString is the string version of FSByte.
func FSString(useLocal bool, name string) (string, error) {
	b, err := FSByte(useLocal, name)
	return string(b), err
}

// FSMustString is the string version of FSMustByte.
func FSMustString(useLocal bool, name string) string {
	return string(FSMustByte(useLocal, name))
}

var _escData = map[string]*_escFile{

	"/static/test.css": {
		local:   "static/test.css",
		size:    59,
		modtime: 1464367588,
		compressed: `
H4sIAAAJbogA/0rKT6lUqOblUgCCpMTk7PSi/NK8FN3k/Jz8IiuFlMSi7KLUFGuIPFQwuTIxDyhSCwgA
AP//R+/zwjsAAAA=
`,
	},

	"/static/test.html": {
		local:   "static/test.html",
		size:    817,
		modtime: 1464367325,
		compressed: `
H4sIAAAJbogA/3yTvW7jMAzH9wPuHQgNt+WE5MZTPHfo0CEvwNhspJS2BVEJkrcv/ZFW+aphQKL449/U
X7LzueXq9y/Qx3nCZp6PcQ6ZqdqQZHjDHTk7LRQEh+4DfKL3tcmK/a1FDCTitZF8ZhJPlA3kc6QBOGU7
AraUkDqFmEFSPWvs5apij0ecGFM5O80uDduiY7ftm3Mp7JcQmklzEbV9U+7EL6/Q1Yh6Yu5N9TIMiqxK
JFYbHwT03R9UBSF/i8UrMBHUjCJrQ03IfVrUPfdJgLHbLYYtmQIfS5pwvJSooXSbn2yK2F2gQQQiY+ge
oXd4SxkhYsJdwujhUQt3tWpCoj/dVuL/INMofTuv0AnbyHMwyOm5DDVPenmWfJhwVt0o/bRqaBkf+LaC
Q/Ua9DCWzur0SXL1U/LfTdLZr684O98rvRDDr/IZAAD//3bVt40xAwAA
`,
	},

	"/static/test.js": {
		local:   "static/test.js",
		size:    16,
		modtime: 1464366989,
		compressed: `
H4sIAAAJbogA/0rMSS0q0VDKSM3JyVdU0rQGBAAA//9uPMx/EAAAAA==
`,
	},

	"/": {
		isDir: true,
		local: "/",
	},

	"/static": {
		isDir: true,
		local: "/static",
	},
}
