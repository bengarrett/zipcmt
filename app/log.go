// © Ben Garrett https://github.com/bengarrett/zipcmt

package app

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strings"
	"text/tabwriter"
	"time"

	humanize "github.com/dustin/go-humanize"
	"github.com/gookit/color"
	gap "github.com/muesli/go-app-paths"
)

var ErrLog = errors.New("log name not set")

// Error saves the error to either a new or append an existing log file.
func (c *Config) Error(err error) {
	const format = "ERROR: %s"
	if err == nil {
		return
	}

	color.Error.Tips(fmt.Sprint(err))
	if err := c.WriteLog(fmt.Sprintf(format, err)); err != nil {
		log.Fatal(err)
	}
}

// WriteLog saves the string to an appended or new log file.
func (c *Config) WriteLog(s string) error {
	if !c.Log || s == "" {
		return nil
	}

	const format = "%s log file: %w"
	if c.LogName() == "" {
		c.SetLog()
	}
	logPath := c.LogName()
	if logPath == "" {
		return fmt.Errorf(format, "", ErrLog)
	}

	dir := filepath.Dir(logPath)
	const dPerm = 0o755
	if err := os.MkdirAll(dir, dPerm); err != nil {
		return fmt.Errorf(format, "make directory", err)
	}

	flags := os.O_APPEND | os.O_CREATE | os.O_WRONLY
	const fPerm = 0o644
	out, err := os.OpenFile(logPath, flags, fPerm)
	if err != nil {
		return fmt.Errorf(format, "open", err)
	}
	defer out.Close()

	st, err := out.Stat()
	if err != nil {
		return fmt.Errorf(format, "stat", err)
	}
	logger := log.New(out, "zipcmt|", log.LstdFlags)
	if st.Size() == 0 {
		c.logHeader(logger)
	}

	var buf strings.Builder
	fmt.Fprintf(&buf, "zip#: %07d; cmmt#: %07d; ", c.Zips, c.Cmmts)
	if !c.Dupes {
		const hashLen = 32
		bytes := uint64(len(c.hashes)) * uint64(hashLen)
		fmt.Fprintf(&buf, "hashes: %s; ", humanize.Bytes(bytes))
	}

	l := fmt.Sprintf("zip#: %07d; cmmt#: %07d; ", c.Zips, c.Cmmts)
	if !c.Dupes {
		const hashLen = 32
		x := uint64(len(c.hashes)) * uint64(hashLen)
		l += fmt.Sprintf("hashes: %s; ", humanize.Bytes(x))
	}
	if c.SaveName != "" {
		l += fmt.Sprintf("names: %s; ", humanize.Bytes(uint64(c.names)))
	}
	l += s + "\n"
	logger.Print(l)

	return nil
}

// logHeader creates a header for new log files that lists all the values of Config.
func (c *Config) logHeader(logger *log.Logger) {
	w := new(tabwriter.Writer)
	const tabWidth = 8
	w.Init(logger.Writer(), 0, tabWidth, 0, '\t', 0)
	fmt.Fprintln(w, "Zip Comment Log - Configurations and arguments")
	fmt.Fprintln(w, "")
	// see: https://scene-si.org/2017/12/21/introduction-to-reflection/
	v := reflect.ValueOf(c).Elem()
	t := v.Type()
	const format = "%02d. %s:\t\t%v\n"
	for i := range v.NumField() {
		fmt.Fprintf(w, format, i+1, t.Field(i).Name, v.Field(i))
		if t.Field(i).Name == "test" {
			break
		}
	}
	fmt.Fprintln(w)
	w.Flush()
}

func logName() string {
	const yyyymmddTime = "20060102150405"
	filename := time.Now().Format(yyyymmddTime) + ".log"
	name, err := gap.NewScope(gap.User, "zipcmt").LogPath(filename)
	if err != nil {
		dir, err2 := os.UserHomeDir()
		if err2 != nil {
			log.Fatalln(fmt.Errorf("log name user home dir: %w", err2))
		}
		name = path.Join(dir, filename)
	}

	return name
}
