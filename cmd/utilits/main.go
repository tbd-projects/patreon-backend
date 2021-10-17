package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/sirupsen/logrus"
	"golang.org/x/tools/container/intsets"
	"io/ioutil"
	"log"
	"os"
	"patreon/internal/app"
	"strings"
	"time"
)

type Log struct {
	Level    string    `json:"level,omitempty"`
	Method   string    `json:"method,omitempty"`
	Msg      string    `json:"msg,omitempty"`
	Adr      string    `json:"remote_addr,omitempty"`
	Url      string    `json:"urls,omitempty"`
	Time     time.Time `json:"time,omitempty"`
	WorkTime int64     `json:"work_time,omitempty"`
}

var (
	configPath string
	logLevel   string
	needFile   string
	allFiles   bool
)

func init() {
	flag.StringVar(&configPath, "config-path", "configs/server.toml", "path to config file")
	flag.StringVar(&logLevel, "level", "trace", "skip levels")
	flag.StringVar(&needFile, "name-file", "", "concrate files to print")
	flag.BoolVar(&allFiles, "all", false, "print all logs")
}

func printLogFromFile(logger *logrus.Logger, fileName string, fileTime time.Time) error {
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}

	defer func() {
		if err = file.Close(); err != nil {
			logger.Fatal(err)
		}
	}()

	diff := time.Now().Sub(fileTime)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		bytes := scanner.Bytes()
		lg := Log{}
		err = json.Unmarshal(bytes, &lg)
		if err != nil {
			return err
		}

		level, err := logrus.ParseLevel(lg.Level)
		if err != nil {
			return err
		}

		logger.WithTime(lg.Time.Add(diff)).WithFields(logrus.Fields{
			"urls":        lg.Url,
			"method":      lg.Method,
			"remote_addr": lg.Adr,
			"work_time":   lg.WorkTime,
		}).Log(level, lg.Msg)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return nil
}

func fileNameWithoutExtension(fileName string) string {
	if pos := strings.LastIndexByte(fileName, '.'); pos != -1 {
		return fileName[:pos]
	}
	return fileName
}

func parseTimeFromFileName(fileName string) (time.Time, error) {
	formatTime := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d", 2006, 1, 2, 15, 04, 05)
	tmp, err := time.Parse(formatTime, fileNameWithoutExtension(fileName))
	if err != nil {
		return time.Now(), err
	}
	return tmp, err
}

func main() {
	flag.Parse()

	config := app.NewConfig()
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Fatal(err)
	}

	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		log.Fatal(err)
	}

	logger := logrus.New()
	logger.SetLevel(level)

	if needFile != "" {
		tmp, err := parseTimeFromFileName(needFile)
		if err != nil {
			log.Printf("error in file %v", err)
		}
		err = printLogFromFile(logger, config.LogAddr+needFile, tmp)
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	files, err := ioutil.ReadDir(config.LogAddr)
	if err != nil {
		log.Fatal(err)
	}

	if allFiles {
		for _, file := range files {
			fmt.Printf("Log from : %s\n", file.Name())
			tmp, err := parseTimeFromFileName(file.Name())
			if err != nil {
				log.Printf("error in file %v", err)
			}

			err = printLogFromFile(logger, config.LogAddr+file.Name(), tmp)
			if err != nil {
				log.Printf("error in file %v", err)
			}
		}
		return
	}

	var lastestFile string
	var lastestTime time.Time
	lastestTime.AddDate(intsets.MaxInt, intsets.MaxInt, intsets.MaxInt)
	for _, file := range files {
		if lastestTime.Second() < file.ModTime().Second() {
			lastestTime = file.ModTime()
			lastestFile = file.Name()
		}
		file.ModTime()
	}

	fmt.Printf("Log from : %s\n", lastestFile)
	tmp, err := parseTimeFromFileName(lastestFile)
	if err != nil {
		log.Printf("error in file %v", err)
	}
	err = printLogFromFile(logger, config.LogAddr+lastestFile, tmp)
	if err != nil {
		log.Printf("error in file %v", err)
	}
}
