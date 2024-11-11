package delivery

import (
	"fmt"

	"github.com/kpango/glg"
)

const (
	PRINTOUT_TYPE_LOG = iota
	PRINTOUT_TYPE_WARN
	PRINTOUT_TYPE_ERR
)

var (
	printOutChan = make(chan printOutObj)
)

type printOutObj struct {
	Type    int
	Message []interface{}
}

func PrintWarnf(format string, message ...interface{}) {
	PrintWarn(fmt.Sprintf(format, message...))
}

func PrintWarn(message ...interface{}) {
	printOut(PRINTOUT_TYPE_ERR, message...)
}

func PrintErrorf(format string, message ...interface{}) {
	PrintError(fmt.Sprintf(format, message...))
}

func PrintError(message ...interface{}) {
	printOut(PRINTOUT_TYPE_ERR, message...)
}

func PrintLogf(format string, message ...interface{}) {
	PrintLog(fmt.Sprintf(format, message...))
}

func PrintLog(message ...interface{}) {
	printOut(PRINTOUT_TYPE_LOG, message...)
}

func printOut(printType int, message ...interface{}) {
	po := printOutObj{
		Type:    printType,
		Message: message,
	}
	printOutChan <- po
}

func StartPrintoutObserver() {
	for po := range printOutChan {
		if po.Type == PRINTOUT_TYPE_ERR {
			_ = glg.Error(po.Message...)
		} else if po.Type == PRINTOUT_TYPE_LOG {
			_ = glg.Log(po.Message...)
		}
	}
}
