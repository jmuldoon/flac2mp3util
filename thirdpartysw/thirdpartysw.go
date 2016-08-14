package thirdpartysw

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/golang/glog"
	"github.com/jmuldoon/flac2mp3util/util"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const (
	thirdpartyswdir  string        = "./thirdpartysw"
	sourceconfigname string        = "sources.json"
	dependencyswpath string        = "./deps"
	clienttimeout    time.Duration = 300
)

type ThirdPartyer interface {
	InitializeClient()
	ReadURLs() (err error)
	Download() (err error)
}

type ThirdPartyType struct {
	Client       *http.Client
	Dependencies []*Url
}

type Url struct {
	URL string `json:"url"`
}

// InitializeClient sets up the client to be used for the http requests since the default
// one is garbage due to not having a timeout.
func (tp *ThirdPartyType) InitializeClient() {
	tp.Client = &http.Client{Timeout: time.Second * clienttimeout}
}

// verifyContentIsCompressed returns nil if the content header specifies gzip
// compression, else returns an error
// TODO: Ensure that the server if it doesn't send a gzip encoding header value
// since it may not send it or it is munged by some formatting of antivirus
// garbage. Handle it gracefully. atm this won't be called so I can proceed.
func verifyContentIsCompressed(resp *http.Response) (err error) {
	fmt.Printf("%+v\n", resp.Header.Get("Content-Encoding"))
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		return nil
	default:
		return fmt.Errorf(
			"verifyContentIsCompressed: failed due to incorrect header response")
	}
}

// ungzip decompresses the specified source to the target destination directory.
// if it doesn't exist, it will be created.
func ungzip(source *gzip.Reader, target string) error {
	target = filepath.Join(target, source.Name)
	writer, err := os.Create(target)
	if err != nil {
		return err
	}
	defer writer.Close()

	_, err = io.Copy(writer, source)
	return err
}

// Download retrieves the tarballs from the url list given as a parameter.
func (tp *ThirdPartyType) Download() (err error) {
	if err := util.MakeDirectory(dependencyswpath); err != nil {
		return err
	}
	for _, el := range tp.Dependencies {
		// setup request string and headers
		req, err := http.NewRequest("GET", el.URL, nil)
		if err != nil {
			return err
		}
		req.Header.Add("Accept-Encoding", "gzip")
		glog.V(2).Infof("%+v\n", req)

		// submit the request through the client
		resp, err := tp.Client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		// TODO: re-add this back in see comment on function
		// // validate that response is gzip compression type
		// if err := verifyContentIsCompressed(resp); err != nil {
		// 	return err
		// }

		// assign the io.ReadCloser
		gzipReader, err := gzip.NewReader(resp.Body)
		if err != nil {
			return err
		}
		defer gzipReader.Close()

		// todo: likely needs to be untar instead, which may be a bit more work.
		// either way, figure that out and ensure that the string placeholde "lame"
		// is dynamic. will need to update the struct for that though as well.
		if err := ungzip(gzipReader, filepath.Join(dependencyswpath, "Lame")); err != nil {
			return err
		}

	}
	return nil
}

// decode takes the ThirdPartyType as a receiver to decode the io.Reader json
// structure and store it in the ThirdPartyType's Dependencies field.
func (tp *ThirdPartyType) decode(file io.Reader) (err error) {
	if err = json.NewDecoder(file).Decode(&tp.Dependencies); err != nil {
		return err
	}
	return nil
}

// read takes the ThirdPartyType as a receiver to readin the json information
// from the file specified by the path
func (tp *ThirdPartyType) read(path string) (err error) {
	file, err := os.OpenFile(path, os.O_RDONLY, 0777)
	if err != nil {
		return err
	}
	if err := tp.decode(file); err != nil {
		return err
	}

	return nil
}

// ReadURLs reads in the URLs from the sources.json file and stores them in the
// receiver ThirdPartyType's field Dependencies.
func (tp *ThirdPartyType) ReadURLs() (err error) {
	absPath, err := util.AbsolutePathHelper(thirdpartyswdir, sourceconfigname)
	glog.V(2).Infoln(absPath)
	if err != nil {
		return err
	}
	if err := tp.read(absPath); err != nil {
		return err
	}
	for i, el := range tp.Dependencies {
		glog.V(2).Infof("%d: %s\n", i, el.URL)
	}
	return nil
}
