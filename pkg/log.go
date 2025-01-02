// © Ben Garrett https://github.com/bengarrett/zipcmt

package zipcmt

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"text/tabwriter"
	"time"

	humanize "github.com/dustin/go-humanize"
	"github.com/gookit/color"
	gap "github.com/muesli/go-app-paths"
)

// Error saves the error to either a new or append an existing log file.
func (c *Config) Error(err error) {
	if err == nil {
		return
	}
	color.Error.Tips(fmt.Sprint(err))
	c.WriteLog(fmt.Sprintf("ERROR: %s", err))
}

// WriteLog saves the string to an appended or new log file.
func (c *Config) WriteLog(s string) {
	if !c.Log || s == "" {
		return
	}

	if c.LogName() == "" {
		c.SetLog()
		d := filepath.Dir(c.LogName())
		_, err := os.Stat(d)
		if os.IsNotExist(err) {
			const perm = 0o755
			if err := os.MkdirAll(d, perm); err != nil {
				log.Fatalln(err)
			}
		}
	}

	const perm = 0o644
	f, err1 := os.OpenFile(c.LogName(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, perm)
	if err1 != nil {
		log.Fatalln(err1)
	}

	logger := log.New(f, "zipcmt|", log.LstdFlags)
	st, err2 := f.Stat()
	if err2 != nil {
		f.Close()
		log.Fatalln(err2)
	}
	defer f.Close()
	if st.Size() == 0 {
		c.logHeader(logger)
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
	for i := range v.NumField() {
		fmt.Fprintf(w, "%02d. %s:\t\t%v\n", i+1, t.Field(i).Name, v.Field(i))
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
		h, err2 := os.UserHomeDir()
		if err2 != nil {
			log.Fatalln(fmt.Errorf("logName UserHomeDir: %w", err2))
		}
		name = path.Join(h, filename)
	}
	return name
}
