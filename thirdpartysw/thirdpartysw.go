package thirdpartysw

import (
	"encoding/json"
	"fmt"
	"github.com/golang/glog"
	"github.com/jmuldoon/flac2mp3util/util"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

const (
	thirdpartyswdir  = "./thirdpartysw"
	sourceconfigname = "sources.json"
	dependencyswpath = "./deps"
)

type ThirdPartyer interface {
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

var thirdParty *ThirdPartyType

func init() {
	thirdParty = &ThirdPartyType{
		Client: &http.Client{},
	}
}

// InitDepsFolder creates the deps (dependencies) folder if it doesn't already
// exist
func initDepsFolder() (err error) {
	abs, err := filepath.Abs(dependencyswpath)
	if err != nil {
		return err
	}
	// if it exists we ignore the error, by not creating the dir
	if _, err := os.Stat(abs); os.IsNotExist(err) {
		os.Mkdir(abs, 0777)
	}
	return nil
}

// Download retrieves the tarballs from the url list given as a parameter.
func (tp *ThirdPartyType) Download() (err error) {
	if err := initDepsFolder(); err != nil {
		return err
	}
	for _, el := range tp.Dependencies {
		resp, err := tp.Client.Get(el.URL)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		fmt.Printf("%+v\n", resp.Body)
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
