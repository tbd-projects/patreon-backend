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
	"os/exec"
	"path/filepath"
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
	configPath string
	logLevel   string
	needFile   string
	allFiles   bool
	GenMock    string
	SearchURL  string
)

func init() {
	flag.StringVar(&configPath, "config-path", "configs/server.toml", "path to config file")
	flag.StringVar(&logLevel, "level", "trace", "skip levels")
	flag.StringVar(&needFile, "name-file", "", "concrate files to print")
	flag.BoolVar(&allFiles, "all", false, "print all logs")
	flag.StringVar(&GenMock, "gen-mock", "", "genmock")
	flag.StringVar(&SearchURL, "search-url", "*", "search url")
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

		if SearchURL != "*" && SearchURL != lg.Url.String() {
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

//mockgen  -destination=mocks/mock_awards_usecase.go -package=mock_usecase -mock_names=Usecase=AwardsUsecase . Usecase

func generateMock(dir string) {
	files, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if file.IsDir() {
			filesInFile, err := os.ReadDir(dir + "/" + file.Name())
			if err != nil {
				log.Fatal(err)
			}

			haveMockDir := false
			for _, checkedFiles := range filesInFile {
				if checkedFiles.Name() == "mocks" {
					haveMockDir = true
					break
				}
			}

			if !haveMockDir {
				continue
			}

			baseDir := filepath.Base(dir)
			interfaceName := strings.Title(strings.ToLower(baseDir))
			cmd := exec.Command("mockgen", fmt.Sprintf("-destination=mocks/mock_%s_%s.go", file.Name(), baseDir),
				fmt.Sprintf("-package=mock_%s", baseDir),
				fmt.Sprintf("-mock_names=%s=%s%s", interfaceName,
					strings.Title(strings.ToLower(file.Name())), interfaceName), ".", interfaceName)
			cmd.Dir = dir + file.Name() + "/"
			cmd.Stdout = log.Writer()
			cmd.Stderr = log.Writer()
			err = cmd.Run()
			if err != nil {
				log.Fatal(err)
				return
			}
		}
	}
}

func main() {
	flag.Parse()

	if GenMock != "" {
		generateMock(GenMock)
		return
	}

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
