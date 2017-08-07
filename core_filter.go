package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var coresaveDir string = "/data/data/coresave/"

func coreDirExists() bool {
	finfo, err := os.Stat(coresaveDir)
	if err != nil {
		fmt.Println("dir not exists")
		return false
	}
	return finfo.IsDir()
}

func canCoredump(name string, interval int64) bool {
	cores, err := ioutil.ReadDir(coresaveDir)
	if err != nil {
		//unknow error happen, avoid coredump
		return false
	}
	var maxTimestamp int64 = 0
	for _, core := range cores {
		parts := strings.Split(core.Name(), ".")
		if len(parts) == 4 {
			coreName := parts[1]
			timestamp := parts[3]
			if name == coreName {
				timestamp, err := strconv.ParseInt(timestamp, 10, 64)
				if err != nil {
					continue
				}
				if timestamp > maxTimestamp {
					maxTimestamp = timestamp
				}
			}
		}
	}
	now := time.Now().Unix()
	return (now - maxTimestamp) >= (interval * 60) //coredump every hour
}

func main() {

	logfile, err := os.OpenFile("/var/log/core_filter.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	defer logfile.Close()
	if err != nil {
		fmt.Println("core filter log file not exists.")
		os.Exit(-1)
	}
	logger := log.New(logfile, "", log.Ldate|log.Ltime|log.Lshortfile)
	if !coreDirExists() {
		logger.Fatal("core save dir [", coresaveDir, "] doesn't exists")
	}

	pid := flag.String("p", "0", "pid of an excutable program.")
	timestamp := flag.Int64("t", 0, "time of dump, expressed as seconds since the Epoch")
	name := flag.String("e", "", "program name without path prefix")
	interval := flag.Int64("i", 60, "core dump interval. unit: min")
	//unit: mega byte
	maxCoreSize := flag.Int("s", 1024, "max core dump size. unit: mega byte")
	flag.Parse()

	corePath := coresaveDir + "core." + *name + "." + *pid + "." + strconv.FormatInt(*timestamp, 10)

	if canCoredump(*name, *interval) {
		coreFile, err := os.Create(corePath)
		if err != nil {
			logger.Fatal("unable to create core dump file")
			os.Exit(-1)
		}
		defer coreFile.Close()
		//write core dump
		buf := make([]byte, 1024)
		coreSize := 0

		bio := bufio.NewReader(os.Stdin)
		for {
			size, err := bio.Read(buf)
			if err != nil {
				break
			}
			coreSize += size
			if coreSize <= (*maxCoreSize * 1024 * 1024) {
				coreFile.Write(buf[:size])
			}
		}
		coreFile.Sync()
	} else {
		//already coredump with 1 hour, just touch core file
		coreFile, err := os.Create(corePath)
		if err != nil {
			logger.Fatal("unable to create core dump file")
			os.Exit(-1)
		}
		defer coreFile.Close()
	}
}
