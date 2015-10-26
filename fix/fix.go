package main

import (
	"bufio"
	"compress/gzip"
	"flag"
	// "fmt"
	filepath "github.com/MichaelTJones/walk"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"sync"
	"time"
)

type Logline struct {
	hashedID  string
	unixtime  uint64
	fileorder uint16
	timestamp time.Time
	PID       uint32
	TID       uint32
	level     string
	rest      string
}

var formatters = [...]*regexp.Regexp{
	regexp.MustCompile(`^(?P<hashedID>[0-9a-f]{40})\s+(?:\d+)\s+(?P<unixtime>[\d.]+)\s+(?P<timestamp>\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}\.\d+)\s+(?P<PID>\d+)\s+(?P<TID>\d+)\s+(?P<level>\w+)(?P<rest>.*)$`),
}

func parseFile(dirname string) error {
	defer waitGroup.Done()
	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		return err
	}

	var currentFormatter *regexp.Regexp

	for _, file := range files {
		file, err := os.Open(path.Join(dirname, file.Name()))
		if err != nil {
			continue
		}
		defer file.Close()

		gunzip, err := gzip.NewReader(file)
		if err != nil {
			continue
		}
		defer gunzip.Close()

		scanner := bufio.NewScanner(gunzip)
		for scanner.Scan() {
			line := scanner.Text()
			if currentFormatter == nil {
				for _, formatter := range formatters {
					if formatter.MatchString(line) {
						currentFormatter = formatter
					}
				}
			}
			if currentFormatter == nil {
				continue
			}

			match := currentFormatter.FindStringSubmatch(line)
			matches := make(map[string]string)
			for i, name := range currentFormatter.SubexpNames() {
				if i == 0 || name == "" {
					continue
				}
				matches[name] = match[i]
			}
			logline := Logline{
				hashedID: matches["hashedID"],
			}
			_ = logline
		}
	}

	return nil
}

var oldDirPattern = regexp.MustCompile(`.*time/\d{4}/\d{2}/\d{2}$`)
var tagDirPattern = regexp.MustCompile(`.*tag$`)
var waitGroup sync.WaitGroup

func check(path string, f os.FileInfo, err error) error {
	if tagDirPattern.MatchString(path) {
		return filepath.SkipDir
	}
	if oldDirPattern.MatchString(path) {
		waitGroup.Add(1)
		go parseFile(path)
	}
	return nil
}

func main() {
	flag.Parse()
	root := flag.Arg(0)
	filepath.Walk(root, check)
	waitGroup.Wait()
}
