// Code generated by vfsgen; DO NOT EDIT.

package rbac

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	pathpkg "path"
	"time"
)

// Root statically implements the virtual filesystem provided to vfsgen.
var Root = func() http.FileSystem {
	fs := vfsgen۰FS{
		"/": &vfsgen۰DirInfo{
			name:    "/",
			modTime: time.Time{},
		},
		"/auth_proxy_role.yaml": &vfsgen۰CompressedFileInfo{
			name:             "auth_proxy_role.yaml",
			modTime:          time.Time{},
			uncompressedSize: 280,

			compressedContent: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x7c\xce\xbd\x4a\xc5\x40\x10\xc5\xf1\x7e\x9f\x62\x48\x9f\x88\x9d\x6c\x6b\x61\x6f\x61\x23\xb7\x98\xec\x3d\xe0\x98\x64\x67\x99\x99\x8d\x1f\x4f\x2f\xb1\x12\x09\xb7\x3e\x9c\x1f\x7f\x6e\xf2\x02\x73\xd1\x9a\xc9\x66\x2e\x13\xf7\x78\x53\x93\x6f\x0e\xd1\x3a\x2d\x0f\x3e\x89\xde\xed\xf7\x69\x91\x7a\xcd\xf4\xb8\x76\x0f\xd8\xb3\xae\x48\x1b\x82\xaf\x1c\x9c\x13\x51\xe5\x0d\x99\x9a\xe9\xe7\xd7\x68\xc7\x68\x7d\x85\xe7\x34\x12\x37\x79\x32\xed\xcd\x33\xbd\x0e\x07\x8e\x1a\x52\xfe\xea\xc3\x25\x11\x19\x5c\xbb\x95\xe3\x43\x34\x52\xe8\x82\x6a\xd8\x05\x1f\x9e\x88\x76\xd8\xfc\x0b\x14\x03\x07\x86\xcb\x19\xfc\xbf\xfa\xcc\xf5\x3e\xbf\xa3\x04\x97\x02\xf7\x5b\xfe\x4f\x00\x00\x00\xff\xff\x4b\x41\x4e\x63\x18\x01\x00\x00"),
		},
		"/auth_proxy_role_binding.yaml": &vfsgen۰CompressedFileInfo{
			name:             "auth_proxy_role_binding.yaml",
			modTime:          time.Time{},
			uncompressedSize: 257,

			compressedContent: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x7c\x8d\x31\x4e\xc0\x30\x0c\x45\xf7\x9c\xc2\x17\x68\x11\x1b\xca\x06\x0c\xec\x45\x62\x77\x13\x17\x4c\xd3\x38\xb2\x9d\x8a\x72\x7a\x54\xa9\xb0\x80\xd8\x6c\xe9\xbd\xff\xb0\xf1\x0b\xa9\xb1\xd4\x08\x3a\x63\x1a\xb1\xfb\x9b\x28\x7f\xa2\xb3\xd4\x71\xbd\xb3\x91\xe5\x66\xbf\x0d\x2b\xd7\x1c\xe1\xb1\x74\x73\xd2\x49\x0a\x3d\x70\xcd\x5c\x5f\xc3\x46\x8e\x19\x1d\x63\x00\xa8\xb8\x51\x84\xa6\xf2\x71\x0c\x2a\x85\xe6\x8b\x39\xef\x89\x96\x13\xc1\xc6\x4f\x2a\xbd\xfd\x93\x0b\x00\xbf\x6a\x7f\x8c\x07\xeb\xf3\x3b\x25\xb7\x18\x86\x4b\x78\x26\xdd\x39\xd1\x7d\x4a\xd2\xab\xff\x38\x99\x16\xec\xe5\xfb\xb7\x86\x89\x22\xd8\x61\x4e\x5b\xf8\x0a\x00\x00\xff\xff\x38\xc6\x1a\x54\x01\x01\x00\x00"),
		},
		"/auth_proxy_service.yaml": &vfsgen۰CompressedFileInfo{
			name:             "auth_proxy_service.yaml",
			modTime:          time.Time{},
			uncompressedSize: 379,

			compressedContent: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x8c\x8f\xb1\x6a\x04\x31\x0c\x44\x7b\x7f\x85\xd8\xde\x09\x21\x57\x04\x7f\x45\x20\x90\x5e\xf1\x0d\xb7\x26\xb6\x65\x24\xdd\x41\xfe\x3e\x78\xb3\x07\x29\xb6\xb8\x72\x34\x62\x66\x1e\x8f\xf2\x09\xb5\x22\x3d\xd1\xed\x25\x7c\x97\x7e\x4e\xf4\x01\xbd\x95\x8c\xd0\xe0\x7c\x66\xe7\x14\x88\xb8\x77\x71\xf6\x22\xdd\xa6\x24\x1a\x2a\x0d\xbe\xe2\x6a\x4f\x45\x9e\x87\xa8\x27\x5a\xde\x4e\xa7\xd7\xe5\xc0\xb6\xbc\xa2\x21\xd1\xea\x3e\xec\xd0\x57\x1e\x48\xb4\xb8\x5e\x31\x03\x2a\x7f\xa1\xee\x4d\x59\xba\xab\xd4\x38\x2a\x77\xa4\xbb\xac\xd0\xd8\xb8\xf3\x05\x1a\x88\x3a\xb7\x43\x2b\x36\xb8\x96\x6c\xd1\x76\xa6\xbf\x57\x1b\x9c\x91\xc8\x7e\xcc\xd1\x82\x0d\xe4\x59\x35\x21\xb6\xce\xb8\xe7\xfd\x9b\xbb\xf1\x4d\xbc\x4d\x3a\xeb\x05\xfe\xbe\x1d\xef\x4f\x86\x8a\xec\xa2\x8f\x8e\xfe\x0d\x00\x00\xff\xff\xe8\x4f\xf6\xaf\x7b\x01\x00\x00"),
		},
		"/kustomization.yaml": &vfsgen۰CompressedFileInfo{
			name:             "kustomization.yaml",
			modTime:          time.Time{},
			uncompressedSize: 344,

			compressedContent: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x6c\x8f\xb1\x4e\xc5\x30\x0c\x45\xf7\x7e\x85\xa5\x2e\x30\xb4\x19\xd8\xde\xca\x87\x44\x49\xea\xd7\x58\x24\x71\xe5\x38\x94\xf2\xf5\xe8\x45\x80\x5a\xc4\x66\xf9\x9e\xab\xab\x23\x58\xb9\x49\xc0\x7a\x1b\x26\x10\x4e\x38\x1f\x2e\xa7\xef\xdb\x7a\x2a\x0b\x95\xf5\xe7\x97\xd0\x2d\x28\x16\x13\x06\x25\x2e\xf6\xcc\xff\x97\x5d\xfb\x23\xbc\x72\xce\x58\x14\x34\x22\xdc\x39\x25\xde\xa9\xac\xf0\x02\x89\x0a\x56\xa0\x3b\x1c\xdc\x60\x77\x0f\x82\x61\xa1\xea\x7c\xc2\x61\xec\xb8\x6b\x1a\x61\x13\xfe\x38\xe0\x29\xaa\x6e\xf5\x66\xcc\x4a\x1a\x9b\x9f\x03\x67\xe3\xc5\x95\xf0\x69\xde\x9a\xc7\x49\xbc\x0b\x53\x47\x9f\x87\x11\xf6\x48\xa1\x37\x15\x83\xd6\xc7\x82\x80\xc9\xa8\x42\xa1\x02\x96\x65\x63\x2a\x3a\x0f\x53\x5f\xb0\xbd\x66\x2b\xca\x3b\x85\x5f\xb5\x53\x72\x36\xfe\xf3\xbe\xca\x7e\x05\x00\x00\xff\xff\xac\x90\x8c\xf8\x58\x01\x00\x00"),
		},
		"/leader_election_role.yaml": &vfsgen۰CompressedFileInfo{
			name:             "leader_election_role.yaml",
			modTime:          time.Time{},
			uncompressedSize: 355,

			compressedContent: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x94\x4e\xbd\x4e\xc3\x40\x0c\xde\xef\x29\xac\x32\x27\x15\x1b\xca\x0b\xb0\x33\xb0\xbb\xb9\x8f\xd6\xea\xe5\x7c\xb2\x7d\x45\xe2\xe9\x51\xd2\x94\x05\x31\x30\xf9\xd3\xf7\xeb\x27\x6a\xb0\x45\xdc\x45\xab\x53\x28\x65\xa5\x02\xce\x30\x42\xc1\x1c\xa2\x75\x4c\xdc\xe4\x1d\xb6\x5a\x26\xb2\x13\xcf\x23\xf7\xb8\xa8\xc9\x17\x6f\xfa\xf5\xc5\x47\xd1\xe3\xed\x39\x5d\xa5\xe6\x89\xde\xb4\x20\x2d\x08\xce\x1c\x3c\x25\xa2\xca\x0b\xa6\xbd\x76\x78\xd4\x0e\xb6\xda\xac\x17\xf8\x94\x06\xe2\x26\xaf\xa6\xbd\xf9\x1a\x18\xe8\x70\x48\x44\x06\xd7\x6e\x33\x76\x6e\xd6\xfa\x21\xe7\x85\x9b\x27\xa2\x1b\xec\xb4\xf3\x67\xc4\x76\x8b\xf8\x1d\x7c\x72\xcc\x97\x7b\xc4\xc0\x81\x0d\xf6\x96\x1f\xb0\xfd\xe8\x19\x05\x81\xff\xce\x1f\x3d\x38\xfa\x1f\x5f\xfc\xda\xf9\x0e\x00\x00\xff\xff\xcb\xa9\x3f\x0f\x63\x01\x00\x00"),
		},
		"/leader_election_role_binding.yaml": &vfsgen۰CompressedFileInfo{
			name:             "leader_election_role_binding.yaml",
			modTime:          time.Time{},
			uncompressedSize: 263,

			compressedContent: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x7c\x8d\x3d\x8a\xc3\x30\x10\x46\x7b\x9d\x62\x2e\x60\x2f\xdb\x2d\xea\x76\x9b\xed\x1d\x48\x3f\x96\x3e\x27\x13\xcb\x1a\xa1\x1f\x43\x72\xfa\x60\x70\x92\xce\xdd\x0c\xbc\xf7\x3d\x4e\x72\x46\x2e\xa2\xd1\x52\x1e\xd9\xf5\xdc\xea\x55\xb3\x3c\xb8\x8a\xc6\x7e\xfe\x29\xbd\xe8\xd7\xfa\x6d\x66\x89\xde\xd2\xa0\x01\x7f\x12\xbd\xc4\x8b\x59\x50\xd9\x73\x65\x6b\x88\x22\x2f\xb0\x14\xc0\x1e\xb9\x43\x80\xdb\xec\x2e\x6b\xc0\xb8\xd3\xdb\x3d\x60\xda\x60\x4e\xf2\x9f\xb5\xa5\x83\xa2\x21\xfa\x04\x0f\xf7\x4d\x69\xe3\x0d\xae\x16\x6b\xba\xdd\x39\x21\xaf\xe2\xf0\xeb\x9c\xb6\x58\xdf\xb6\xc7\xc4\x2d\xbc\xfe\x92\xd8\xc1\x52\xb9\x97\x8a\xc5\x3c\x03\x00\x00\xff\xff\xc0\xbe\x81\x75\x07\x01\x00\x00"),
		},
		"/role.yaml": &vfsgen۰CompressedFileInfo{
			name:             "role.yaml",
			modTime:          time.Time{},
			uncompressedSize: 1876,

			compressedContent: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xcc\x55\x3d\x6f\xdc\x30\x0c\xdd\xfd\x2b\x84\xec\xbe\xa0\x5b\xe1\xb5\x43\xf7\xa2\xe8\x4e\x4b\x3c\x87\x38\x49\x14\x48\xc9\xd7\xcb\xaf\x2f\xe4\x8f\x24\x3d\x9f\x81\x14\x48\x0e\x9d\x4c\x53\xa4\xde\x23\x9f\x44\x35\x6d\xdb\x36\x90\xe8\x17\x8a\x12\xc7\xce\x48\x0f\xf6\x00\x25\x3f\xb1\xd0\x33\x64\xe2\x78\x38\x7d\xd5\x03\xf1\xe3\xf8\xa5\x39\x51\x74\x9d\xf9\xe6\x8b\x66\x94\x1f\xec\xb1\x09\x98\xc1\x41\x86\xae\x31\xc6\x0a\x4e\x09\x3f\x29\xa0\x66\x08\xa9\x33\xb1\x78\xdf\x18\x13\x21\x60\x67\x02\x44\x18\x50\x5a\xa9\x89\x52\x3c\x6a\xd7\xb4\x06\x12\x7d\x17\x2e\x49\xeb\x16\xad\x79\x78\x68\x8c\x11\x54\x2e\x62\x71\xf1\x59\x8e\x47\x1a\x02\x24\x9d\x7e\x15\xad\x60\xae\xf6\x88\xd2\xaf\x31\x15\x1c\x27\xd3\xa1\xc7\xc5\x1c\x30\x4f\x5f\x4f\x3a\x1b\x09\xb2\x7d\x9a\xac\x92\xdc\x9a\x70\x9e\x9c\xef\xa2\x82\x23\xc6\x1d\xe8\x77\xe5\xd7\x4e\x68\x02\x8b\x73\x29\x91\x1d\x5e\xed\xb6\xa1\xfc\x0f\xec\x52\x55\x51\x33\xc6\x3c\xb2\x2f\x01\xad\x07\x0a\x33\x52\x62\xb7\x76\x4f\x46\xb2\x08\xd6\x72\x99\x4b\x79\xf1\x7d\x7e\x4b\x21\x4d\x1a\x5e\xd1\x76\x80\x81\xa3\xe2\xca\x26\x43\xc6\x63\xf1\x7a\x0f\x91\x17\x46\xad\xc1\xdf\x19\x63\xbd\x03\xb7\x08\x62\xf2\x7c\x09\xbb\xd2\x7f\x24\xa1\xbf\x78\xb4\x26\x62\x3e\xb3\x9c\x28\x0e\xcb\x3d\xdc\xb2\xa3\x38\x08\xaa\xde\x41\xbe\x2b\x6e\x89\x3d\xd9\xcb\x8d\x63\xc8\x4e\xd1\x16\xa1\x7c\x99\x42\xe8\x0e\xd4\x3c\x0f\x43\x6d\x52\x0f\xf1\x19\xc8\x7a\x2e\xee\x66\xb3\xec\x3c\xbc\x8e\x9e\xcf\xfa\xd6\xc1\x25\xa7\xb2\x1c\xc0\xd7\xc5\x65\xd7\xf9\xe7\x35\xe4\xff\x28\x65\x25\xf7\x58\x2f\x4c\xd9\x19\x23\x1b\xe0\x0d\x5c\xe0\x48\x99\xa5\x22\x5a\x16\x64\x3d\x58\x0e\x5b\xb4\x65\x46\x2c\xd1\x9f\xdf\x85\xdd\x87\x68\x57\xd3\xfa\xae\xf4\x14\xdd\x8b\x60\x6f\xfc\xb3\x63\x13\xb1\x2e\x7d\x68\x29\x7f\x02\x00\x00\xff\xff\x09\x96\xc5\x9e\x54\x07\x00\x00"),
		},
		"/role_binding.yaml": &vfsgen۰CompressedFileInfo{
			name:             "role_binding.yaml",
			modTime:          time.Time{},
			uncompressedSize: 261,

			compressedContent: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x7c\x8d\xb1\x4e\x03\x31\x0c\x86\xf7\x3c\x85\x5f\xa0\x87\xd8\x50\x36\x60\x60\x2f\x12\xbb\x2f\x71\x8b\x69\x62\x47\xb6\x53\x09\x9e\x1e\x9d\x74\xb0\x80\xba\xd9\xd2\xf7\xfd\x1f\x0e\x7e\x23\x73\x56\xc9\x60\x2b\x96\x05\x67\xbc\xab\xf1\x17\x06\xab\x2c\x97\x07\x5f\x58\xef\xae\xf7\xe9\xc2\x52\x33\x3c\xb7\xe9\x41\x76\xd4\x46\x4f\x2c\x95\xe5\x9c\x3a\x05\x56\x0c\xcc\x09\x40\xb0\x53\x86\x8e\x82\x67\xb2\x83\x69\xa3\x75\xa7\xb6\xfb\x48\xa7\x0d\xc2\xc1\x2f\xa6\x73\xdc\x08\x26\x80\x3f\xbd\x7f\xe7\x93\xcf\xf5\x83\x4a\x78\x4e\x87\x5d\x79\x25\xbb\x72\xa1\xc7\x52\x74\x4a\xfc\x5a\x95\x4e\x38\xdb\xcf\xef\x03\x0b\x65\xf0\x4f\x0f\xea\xe9\x3b\x00\x00\xff\xff\x3d\x57\xd3\xb6\x05\x01\x00\x00"),
		},
	}
	fs["/"].(*vfsgen۰DirInfo).entries = []os.FileInfo{
		fs["/auth_proxy_role.yaml"].(os.FileInfo),
		fs["/auth_proxy_role_binding.yaml"].(os.FileInfo),
		fs["/auth_proxy_service.yaml"].(os.FileInfo),
		fs["/kustomization.yaml"].(os.FileInfo),
		fs["/leader_election_role.yaml"].(os.FileInfo),
		fs["/leader_election_role_binding.yaml"].(os.FileInfo),
		fs["/role.yaml"].(os.FileInfo),
		fs["/role_binding.yaml"].(os.FileInfo),
	}

	return fs
}()

type vfsgen۰FS map[string]interface{}

func (fs vfsgen۰FS) Open(path string) (http.File, error) {
	path = pathpkg.Clean("/" + path)
	f, ok := fs[path]
	if !ok {
		return nil, &os.PathError{Op: "open", Path: path, Err: os.ErrNotExist}
	}

	switch f := f.(type) {
	case *vfsgen۰CompressedFileInfo:
		gr, err := gzip.NewReader(bytes.NewReader(f.compressedContent))
		if err != nil {
			// This should never happen because we generate the gzip bytes such that they are always valid.
			panic("unexpected error reading own gzip compressed bytes: " + err.Error())
		}
		return &vfsgen۰CompressedFile{
			vfsgen۰CompressedFileInfo: f,
			gr:                        gr,
		}, nil
	case *vfsgen۰DirInfo:
		return &vfsgen۰Dir{
			vfsgen۰DirInfo: f,
		}, nil
	default:
		// This should never happen because we generate only the above types.
		panic(fmt.Sprintf("unexpected type %T", f))
	}
}

// vfsgen۰CompressedFileInfo is a static definition of a gzip compressed file.
type vfsgen۰CompressedFileInfo struct {
	name              string
	modTime           time.Time
	compressedContent []byte
	uncompressedSize  int64
}

func (f *vfsgen۰CompressedFileInfo) Readdir(count int) ([]os.FileInfo, error) {
	return nil, fmt.Errorf("cannot Readdir from file %s", f.name)
}
func (f *vfsgen۰CompressedFileInfo) Stat() (os.FileInfo, error) { return f, nil }

func (f *vfsgen۰CompressedFileInfo) GzipBytes() []byte {
	return f.compressedContent
}

func (f *vfsgen۰CompressedFileInfo) Name() string       { return f.name }
func (f *vfsgen۰CompressedFileInfo) Size() int64        { return f.uncompressedSize }
func (f *vfsgen۰CompressedFileInfo) Mode() os.FileMode  { return 0444 }
func (f *vfsgen۰CompressedFileInfo) ModTime() time.Time { return f.modTime }
func (f *vfsgen۰CompressedFileInfo) IsDir() bool        { return false }
func (f *vfsgen۰CompressedFileInfo) Sys() interface{}   { return nil }

// vfsgen۰CompressedFile is an opened compressedFile instance.
type vfsgen۰CompressedFile struct {
	*vfsgen۰CompressedFileInfo
	gr      *gzip.Reader
	grPos   int64 // Actual gr uncompressed position.
	seekPos int64 // Seek uncompressed position.
}

func (f *vfsgen۰CompressedFile) Read(p []byte) (n int, err error) {
	if f.grPos > f.seekPos {
		// Rewind to beginning.
		err = f.gr.Reset(bytes.NewReader(f.compressedContent))
		if err != nil {
			return 0, err
		}
		f.grPos = 0
	}
	if f.grPos < f.seekPos {
		// Fast-forward.
		_, err = io.CopyN(ioutil.Discard, f.gr, f.seekPos-f.grPos)
		if err != nil {
			return 0, err
		}
		f.grPos = f.seekPos
	}
	n, err = f.gr.Read(p)
	f.grPos += int64(n)
	f.seekPos = f.grPos
	return n, err
}
func (f *vfsgen۰CompressedFile) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case io.SeekStart:
		f.seekPos = 0 + offset
	case io.SeekCurrent:
		f.seekPos += offset
	case io.SeekEnd:
		f.seekPos = f.uncompressedSize + offset
	default:
		panic(fmt.Errorf("invalid whence value: %v", whence))
	}
	return f.seekPos, nil
}
func (f *vfsgen۰CompressedFile) Close() error {
	return f.gr.Close()
}

// vfsgen۰DirInfo is a static definition of a directory.
type vfsgen۰DirInfo struct {
	name    string
	modTime time.Time
	entries []os.FileInfo
}

func (d *vfsgen۰DirInfo) Read([]byte) (int, error) {
	return 0, fmt.Errorf("cannot Read from directory %s", d.name)
}
func (d *vfsgen۰DirInfo) Close() error               { return nil }
func (d *vfsgen۰DirInfo) Stat() (os.FileInfo, error) { return d, nil }

func (d *vfsgen۰DirInfo) Name() string       { return d.name }
func (d *vfsgen۰DirInfo) Size() int64        { return 0 }
func (d *vfsgen۰DirInfo) Mode() os.FileMode  { return 0755 | os.ModeDir }
func (d *vfsgen۰DirInfo) ModTime() time.Time { return d.modTime }
func (d *vfsgen۰DirInfo) IsDir() bool        { return true }
func (d *vfsgen۰DirInfo) Sys() interface{}   { return nil }

// vfsgen۰Dir is an opened dir instance.
type vfsgen۰Dir struct {
	*vfsgen۰DirInfo
	pos int // Position within entries for Seek and Readdir.
}

func (d *vfsgen۰Dir) Seek(offset int64, whence int) (int64, error) {
	if offset == 0 && whence == io.SeekStart {
		d.pos = 0
		return 0, nil
	}
	return 0, fmt.Errorf("unsupported Seek in directory %s", d.name)
}

func (d *vfsgen۰Dir) Readdir(count int) ([]os.FileInfo, error) {
	if d.pos >= len(d.entries) && count > 0 {
		return nil, io.EOF
	}
	if count <= 0 || count > len(d.entries)-d.pos {
		count = len(d.entries) - d.pos
	}
	e := d.entries[d.pos : d.pos+count]
	d.pos += count
	return e, nil
}
