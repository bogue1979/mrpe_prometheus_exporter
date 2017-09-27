package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
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
			if len(line) > 3 {
				c.Comment.Type = line[2]
			}
		case "INTERVAL":
			if len(line) > 3 {
				i, err := strconv.Atoi(line[2])
				if err != nil {
					return false, fmt.Errorf("could not read INTERVAL: %s", err)
				}
				c.Interval = time.Second * time.Duration(i)
			}
		default:
			c.Comment.Comment += strings.Join(line[1:], " ") + "\n"
		}
	}
	// no comment
	if line[len(line)-1] == "\\" {
		if len(line) > 2 {
			if c.Name == "" {
				c.Name, c.Command = line[0], strings.Join(line[1:len(line)-1], " ")
			} else {
				c.Command += " " + strings.Join(line[:len(line)-1], " ")
			}
		} else {
			if c.Name == "" {
				c.Name = line[0]
			} else {
				c.Command += line[0]
			}
		}
	}
	if line[0] != "#" && line[len(line)-1] != "\\" {
		if c.Name == "" && c.Command == "" {
			c.Name, c.Command = line[0], strings.Join(line[1:], " ")
		} else {
			if c.Command == "" {
				c.Command = strings.Join(line, " ")
			} else {
				c.Command += " " + strings.Join(line, " ")
			}
		}
		if c.Valid() {
			return true, nil
		}
	}
	return false, nil
}

func loadCfgDir(d string) (c Checks, e error) {

	files, err := ioutil.ReadDir(d)
	if err != nil {
		log.Fatal("could not read in conf dir:", err)
	}

	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".cfg") {
			fmt.Printf("read Configfile: %s\n", f.Name())
			check, err := loadCfg(d + "/" + f.Name())
			if err != nil {
				return c, err
			}
			c = append(c, check...)
		}
	}
	return c, nil
}

func loadCfg(f string) (c Checks, err error) {
	check := NewCheck()

	file, err := os.Open(f)
	if err != nil {
		return c, fmt.Errorf("could not open file %s: %s", f, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.Fields(scanner.Text())

		// empty
		if len(line) == 0 {
			continue
		}
		complete, err := check.parseLine(line)
		if err != nil {
			return c, fmt.Errorf("could not parse config %s", file.Name())
		}
		if complete {
			c = append(c, check)
			check = NewCheck()
		}
	}

	if err := scanner.Err(); err != nil {
		return c, fmt.Errorf("Scanner error: %s", err)
	}

	return c, nil
}
