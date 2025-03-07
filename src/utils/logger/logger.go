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

var reset = "\033[0m"
var red = "\033[31m"
var green = "\033[32m"
var yellow = "\033[33m"
var blue = "\033[34m"
var magenta = "\033[35m"
var cyan = "\033[36m"
var gray = "\033[37m"
var white = "\033[97m"

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
	printedMessage := fmt.Sprintf(red+ "[Error] " + err.Error(), a...)
	writedMessage :=  fmt.Sprintf("[Error] " + err.Error(), a...)
	writeToFileAndPrint(printedMessage, writedMessage)
}

func ErrorFromString(err string, a ...any) {
	printedMessage := fmt.Sprintf(red+ "[Error] " + err, a...)
	writedMessage :=  fmt.Sprintf("[Error] " + err, a...)
	writeToFileAndPrint(printedMessage, writedMessage)
}

func WarningFromStringF(warning string, a ...any) {
	printedMessage := fmt.Sprintf(yellow+ "[Warning] " + warning, a...)
	writedMessage :=  fmt.Sprintf("[Warning] " + warning, a...)
	writeToFileAndPrint(printedMessage, writedMessage)
}

func WarningFromString(warning string) {
	printedMessage := yellow+ "[Warning] " + warning
	writedMessage :=  "[Warning] " + warning
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

func writeToFileAndPrint(printedMessage string, writedMessage string) {
	timeInFormat := time.Now().Format(timeFormat) + " "
	fmt.Println(timeInFormat + printedMessage + reset)
	file, err := os.OpenFile(Path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		WarningFromStringF(timeInFormat + yellow + "[Warning] Could not log: %v\n" + reset, err)
	}
	defer file.Close()

	_, err = file.WriteString(timeInFormat + writedMessage + "\n")
	if err != nil {
		WarningFromStringF(timeInFormat + yellow + "[Warning] Could not log: %v\n" + reset, err)
	}
}