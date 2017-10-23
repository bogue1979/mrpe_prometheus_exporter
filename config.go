package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

//TODO feels not good but it works ...
func (c *Check) parseLine(line []string) (cpl bool, err error) {

	// comment
	if len(line) > 1 && line[0] == "#" {
		switch line[1] {
		case "HELP":
			if len(line) > 3 {
				c.Comment.Help = strings.Join(line[2:], " ")
			}
		case "TYPE":
			if len(line) == 3 {
				c.Comment.Type = line[2]
			}
		case "INTERVAL":
			if len(line) == 3 {
				i, err := strconv.Atoi(line[2])
				if err != nil {
					return false, fmt.Errorf("could not read INTERVAL: %s", err)
				}
				c.Interval = int64(i)
			}
		default:
			c.Comment.Comment += strings.Join(line[1:], " ") + "\n"
		}
		return false, nil
	}

	if len(line) >= 2 {
		c.Name, c.Command = line[0], strings.Join(line[1:], " ")
	}
	if c.Valid() {
		return true, nil
	}

	return false, nil
}

func loadCfgDir(d string) (c Checks, e error) {

	files, err := ioutil.ReadDir(d)
	if err != nil {
		return c, fmt.Errorf("could not read in conf dir: %s", err)
	}

	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".cfg") {
			log.Infof("read Configfile: %s", f.Name())
			check, err := loadCfg(d + "/" + f.Name())
			if err != nil {
				return c, err
			}
			c = append(c, check...)
		}
	}
	return c, nil
}

//ldCfg merge escaped lines into []string
func ldCfg(f string) (lines []string, err error) {

	file, err := os.Open(f)
	if err != nil {
		return lines, fmt.Errorf("could not open file %s: %s", f, err)
	}
	defer file.Close()

	var appendLine string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := string(scanner.Text())
		if len(line) > 0 && string(line[len(line)-1]) == "\\" {
			appendLine = strings.Join([]string{appendLine, string(line[:len(line)-1])}, " ")
			continue
		}
		appendLine = strings.Join([]string{appendLine, line}, " ")
		lines = append(lines, appendLine)
		appendLine = ""
	}
	if err := scanner.Err(); err != nil {
		return lines, fmt.Errorf("Scanner error: %s", err)
	}
	return lines, nil
}

func loadCfg(f string) (c Checks, err error) {
	check := NewCheck()

	lines, err := ldCfg(f)
	if err != nil {
		return c, fmt.Errorf("could not sanitize file %s: %s", f, err)
	}

	for _, line := range lines {

		l := strings.Fields(line)
		// empty
		if len(l) == 0 {
			continue
		}
		complete, err := check.parseLine(l)
		if err != nil {
			return c, fmt.Errorf("could not parse config %s", f)
		}
		if complete {
			c = append(c, check)
			check = NewCheck()
		}
	}

	return c, nil
}
