package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)


const (
	timeFormat = "[15:04:05.000]"
)

var Reset = "\033[0m"
var Red = "\033[31m"
var Green = "\033[32m"
var Yellow = "\033[33m"
var Blue = "\033[34m"
var Magenta = "\033[35m"
var Cyan = "\033[36m"
var Gray = "\033[37m"
var White = "\033[97m"

var Path string

func CreateFile(path string) error {
	err := os.MkdirAll(filepath.Dir(path), os.ModePerm)
	if err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}

	defer file.Close()
	Path = path
	return nil
}

func ErrorFromErr(err error, a ...any) {
	printedMessage := fmt.Sprintf(Red+ "[Error] " + err.Error(), a)
	writedMessage :=  fmt.Sprintf("[Error] " + err.Error(), a)
	writeToFileAndPrint(printedMessage, writedMessage)
}

func ErrorFromString(err string, a ...any) {
	printedMessage := fmt.Sprintf(Red+ "[Error] " + err, a)
	writedMessage :=  fmt.Sprintf("[Error] " + err, a)
	writeToFileAndPrint(printedMessage, writedMessage)
}

func WarningFromString(warning string, a ...any) {
	printedMessage := fmt.Sprintf(Yellow+ "[Warning] " + warning, a)
	writedMessage :=  fmt.Sprintf("[Warning] " + warning, a)
	writeToFileAndPrint(printedMessage, writedMessage)
}

func LogInfo(message string) {
	printedMessage := "[Info] " + message
	writeToFileAndPrint(printedMessage, printedMessage)
}

func LogInfoF(message string, a ...any) {
	printedMessage := fmt.Sprintf("[Info] " + message, a...)
	writeToFileAndPrint(printedMessage, printedMessage)
}

func WarningFromErr(err error, a ...any) {
	printedMessage := fmt.Sprintf(Yellow+ "[Warning] " + err.Error(), a)
	writedMessage :=  fmt.Sprintf("[Warning] " + err.Error(), a)
	writeToFileAndPrint(printedMessage, writedMessage)
}

func writeToFileAndPrint(printedMessage string, writedMessage string) {
	timeInFormat := time.Now().Format(timeFormat) + " "
	fmt.Println(timeInFormat + printedMessage)
	file, err := os.OpenFile(Path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		fmt.Printf("[Warning] could not log: %v\n", err)
	}
	defer file.Close()

	_, err = file.WriteString(timeInFormat + writedMessage + "\n")
	if err != nil {
		fmt.Printf("[Warning] could not log: %v\n", err)
	}

}