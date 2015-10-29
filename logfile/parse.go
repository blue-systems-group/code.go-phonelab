package phonelab

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
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
	message   string
}

var formatters = [...]*regexp.Regexp{
	regexp.MustCompile(`^(?P<hashedID>[0-9a-f]{40})\s+(?:\d+)\s+(?P<unixtime>[\d.]+)\s+(?P<timestamp>\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}\.\d+)\s+(?P<PID>\d+)\s+(?P<TID>\d+)\s+(?P<level>\w+)(?P<message>.*)$`),
}

func count(reader io.Reader) (int, error) {
	buf := make([]byte, 32768)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := reader.Read(buf)
		if err != nil && err != io.EOF {
			return count, err
		}
		count += bytes.Count(buf[:c], lineSep)
		if err == io.EOF {
			break
		}
	}

	return count, nil
}

func reopenFile(filename string, file *os.File) (io.Reader, error) {
	var reader io.Reader
	var err error
	extension := filepath.Ext(filename)

	if extension == ".out" {
		reader = bufio.NewReader(file)
	} else if extension == ".gz" {
		reader, err = gzip.NewReader(file)
	}
	if err != nil {
		return nil, err
	}
	return reader, nil
}

func Parse(filename string) ([]Logline, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader, err := reopenFile(filename, file)
	lineCount, err := count(reader)

	if err != nil {
		return nil, err
	}

	file.Seek(0, 0)
	reader, err = reopenFile(filename, file)

	loglines := make([]Logline, 0, lineCount)
	var currentFormatter *regexp.Regexp
	scanner := bufio.NewScanner(reader)

	i := 0
	for scanner.Scan() {
		i += 1
		line := scanner.Text()
		if currentFormatter == nil {
			for _, formatter := range formatters {
				if formatter.MatchString(line) {
					currentFormatter = formatter
				}
			}
		}
		if currentFormatter == nil {
			return nil, errors.New("Could not determine logfile format.")
		}

		match := currentFormatter.FindStringSubmatch(line)
		matches := make(map[string]string)
		for i, name := range currentFormatter.SubexpNames() {
			if i == 0 || name == "" {
				continue
			}
			matches[name] = match[i]
		}
		loglines = append(loglines, Logline{
			hashedID: matches["hashedID"],
		})
	}
	if len(loglines) != lineCount {
		return nil, errors.New(fmt.Sprintf("Loglines (%v) does not match line count (%v).\n", len(loglines), lineCount))
	}
	return loglines, nil
}