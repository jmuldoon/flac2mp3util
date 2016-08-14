/*
	$ go run flac2mp3util.go -stderrthreshold=Info -log_dir=./log -v=2 -src "test/dir"
*/

package main

import (
	"flag"
	"fmt"
	"github.com/golang/glog"
	"github.com/jmuldoon/flac2mp3util/thirdpartysw"
	"os"
)

const (
	ERR_STANDARDERR = 1 << iota
	ERR_ARGUMENTINPUT
	ERR_THIRDPARTYSWDEP
)

var args *Arguments

// Arguments is the cli storage struct for the parameters used to run the script
type Arguments struct {
	SourceDir string
	OutputDir string
	BitRate   int
}

// usage redefines the flag.Usage() function.
func usage() {
	usageStr := fmt.Sprintf(`usage: %s -stderrthreshold=[INFO|WARN|FATAL]
		-log_dir=[string] -src "/path/to/src"`+"\n", os.Args[0])
	fmt.Fprintf(os.Stderr, usageStr)
	glog.Errorln(usageStr)
	flag.PrintDefaults()
	os.Exit(ERR_ARGUMENTINPUT)
}

// getPWD gets the present working directory
func getPWD() (pwd string) {
	pwd, err := os.Getwd()
	if err != nil {
		glog.Error(err)
	}
	return pwd
}

// init declares and sets up the flags that are being utilized as well as the
// Arguments struct for usability purposes.
func init() {
	args = &Arguments{}

	flag.StringVar(&args.SourceDir, "src", "",
		"Mandatorily specify the source directory to be converted.")
	flag.StringVar(&args.OutputDir, "out", getPWD(),
		"Optionally specify the output directory for the converted files.")
	flag.IntVar(&args.BitRate, "br", 320,
		"Optionally specify the sample bitrate in kB.")
	flag.Usage = usage
	flag.Parse()
}

// validateArguments checks if the mandatory arguments were passed as parameters
func (a *Arguments) assertArgumentsValid() {
	glog.V(1).Infof("%+v\n", args)
	if a.SourceDir != "" {
		flag.Usage()
	}
}

// main is where the majority of the business logic will reside to run this pkg.
func main() {
	args.assertArgumentsValid()

	thirdpartyer := thirdpartysw.ThirdPartyer(&thirdpartysw.ThirdPartyType{})
	thirdpartyer.InitializeClient() // additional setup for the default client
	if err := thirdpartyer.ReadURLs(); err != nil {
		glog.Errorf("error in %+v . Exited with: %d", err, ERR_THIRDPARTYSWDEP)
	}
	if err := thirdpartyer.Download(); err != nil {
		glog.Errorf("error in %+v . Exited with: %d", err, ERR_THIRDPARTYSWDEP)
	}
}
