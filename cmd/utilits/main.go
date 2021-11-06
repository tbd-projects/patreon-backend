package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/sirupsen/logrus"
	"log"
	"net/url"
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
	Url      url.URL   `json:"urls,omitempty"`
	Time     time.Time `json:"time,omitempty"`
	WorkTime int64     `json:"work_time,omitempty"`
	ReqID    string    `json:"req_id,omitempty"`
}

var (
	configPath    string
	logLevel      string
	needFile      string
	allFiles      bool
	InitMigration int64
	MigrationUp   bool
	MigrationDown bool
	MigrationAdd   bool
	MigrationDel bool
	SearchURL     string
	useServerRepository bool
)

func init() {
	flag.StringVar(&configPath, "config-path", "configs/utilits.toml", "path to config file")
	flag.StringVar(&logLevel, "level", "trace", "skip levels")
	flag.StringVar(&needFile, "name-file", "", "concrate files to print")
	flag.BoolVar(&allFiles, "all", false, "print all logs")
	flag.Int64Var(&InitMigration, "init-migration", -1, "init-migration")
	flag.BoolVar(&MigrationUp, "mig-up", false, "migration-up")
	flag.BoolVar(&MigrationDown, "mig-down", false, "migration-down")
	flag.BoolVar(&MigrationAdd, "mig-add", false, "migration-add-one-change")
	flag.BoolVar(&MigrationDel, "mig-del", false, "migration-delete-one-change")
	flag.StringVar(&SearchURL, "search-url", "", "search url")
	flag.BoolVar(&useServerRepository, "server-run", false, "true if it server run, false if it local run")
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

	tmp := time.Now().In(time.UTC)
	diff := tmp.Sub(fileTime)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		bytes := scanner.Bytes()
		lg := Log{}
		err = json.Unmarshal(bytes, &lg)
		if err != nil {
			return err
		}

		if SearchURL != "" && SearchURL != lg.Url.String() {
			continue
		}

		level, err := logrus.ParseLevel(lg.Level)
		if err != nil {
			return err
		}

		logger.WithTime(lg.Time.In(time.Now().Location()).Add(diff)).WithFields(logrus.Fields{
			"urls":        lg.Url.String(),
			"method":      lg.Method,
			"remote_addr": lg.Adr,
			"work_time":   lg.WorkTime,
			"req_id":      lg.ReqID,
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

	mig := Migrations{log: logger, dbConect: config.LocalRepository.DataBaseUrl}
	if useServerRepository {
		mig = Migrations{log: logger, dbConect: config.ServerRepository.DataBaseUrl}
	}

	switch {
	case InitMigration >= 0:
		mig.MigrationInit(InitMigration)
		return
	case MigrationUp:
		mig.MigrationUp()
		return
	case MigrationDown:
		mig.MigrationDown()
		return
	case MigrationAdd:
		mig.MigrationAddOne()
		return
	case MigrationDel:
		mig.MigrationDelOne()
		return
	}

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

	files, err := os.ReadDir(config.LogAddr)
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
	first := true
	for _, file := range files {
		tmp, err := parseTimeFromFileName(file.Name())
		if err == nil && (lastestTime.Before(tmp) || first) {
			lastestTime = tmp
			lastestFile = file.Name()
			first = false
		}
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
