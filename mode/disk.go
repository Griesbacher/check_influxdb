package mode

import (
	"github.com/griesbacher/check_x/Units"

	"bytes"
	"fmt"
	"github.com/griesbacher/check_x"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
)

//based on: http://stackoverflow.com/a/32482941
func dirSize(path string) (int64, error) {
	var size int64
	adjSize := func(_ string, info os.FileInfo, err error) error {
		if info != nil && !info.IsDir() {
			size += info.Size()
		}
		return err
	}
	err := filepath.Walk(path, adjSize)

	return size, err
}

//Disk will be called on disk mode
func Disk(folder, warn, crit, filterRegex string) (err error) {
	regex, err := regexp.Compile(filterRegex)
	if err != nil {
		return
	}

	warning, err := check_x.NewThreshold(warn)
	if err != nil {
		return
	}

	critical, err := check_x.NewThreshold(crit)
	if err != nil {
		return
	}

	files, err := ioutil.ReadDir(folder)
	if err != nil {
		return
	}
	foundData := false
	for _, f := range files {
		if f.IsDir() && f.Name() == "data" {
			foundData = true
		}
		if f.IsDir() {
			size, err := dirSize(path.Join(folder, f.Name()))
			if err != nil {
				return err
			}
			check_x.NewPerformanceData(f.Name(), float64(size)).Unit("B").Min(0)
		}
	}

	if !foundData {
		return fmt.Errorf("Data folder not found within: %s", folder)
	}

	dataFolder := path.Join(folder, "data")
	files, err = ioutil.ReadDir(dataFolder)
	if err != nil {
		return
	}
	states := check_x.States{}
	var okPrint bytes.Buffer
	var warnPrint bytes.Buffer
	var critPrint bytes.Buffer
	printDB := func(buf *bytes.Buffer, name string, size int64) {
		buf.WriteString(fmt.Sprintf("%s: %s ", name, Units.ByteSize(size)))
	}
	for _, f := range files {
		if f.IsDir() && regex.MatchString(f.Name()) {
			size, err := dirSize(path.Join(dataFolder, f.Name()))
			if err != nil {
				return err
			}

			state := check_x.Evaluator{Warning: warning, Critical: critical}.Evaluate(float64(size))
			states = append(states, state)

			switch state {
			case check_x.OK:
				printDB(&okPrint, f.Name(), size)
				break
			case check_x.Warning:
				printDB(&warnPrint, f.Name(), size)
				break
			case check_x.Critical:
				printDB(&critPrint, f.Name(), size)
				break
			}
			check_x.NewPerformanceData(f.Name(), float64(size)).Unit("B").Min(0).Warn(warning).Crit(critical)
		}
	}
	var toPrint bytes.Buffer
	if critPrint.Len() > 0 {
		toPrint.WriteString("Critical: ")
		toPrint.WriteString(critPrint.String())
	}
	if warnPrint.Len() > 0 {
		toPrint.WriteString("Warning: ")
		toPrint.WriteString(warnPrint.String())
	}
	if okPrint.Len() > 0 {
		toPrint.WriteString("OK: ")
		toPrint.WriteString(okPrint.String())
	}

	worst, err := states.GetWorst()
	if err != nil {
		return
	}
	check_x.Exit(*worst, toPrint.String())

	return nil
}
