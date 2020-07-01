package log

import (
	"fmt"
	"os"
)

func Println(format string, a ...interface{}) {
	if len(a) == 0 {
		fmt.Println(format)
		return
	}
	fmt.Printf(format+"\n", a)
}

func Errorln(format string, a ...interface{}) {
	if len(a) == 0 {
		fmt.Fprintln(os.Stderr, "ERROR: "+format)
	}
	fmt.Fprintf(os.Stderr, "ERROR: "+format+"\n", a)
}

func Warnln(format string, a ...interface{}) {
	if len(a) == 0 {
		fmt.Fprintln(os.Stderr, "WARN: "+format)
	}
	fmt.Fprintf(os.Stderr, "WARN: "+format+"\n", a)
}

func PrintCmdOutput(id string, output []byte) {
	fmt.Printf("====== Output start for event: %s ======\n%s\n====== Output end ======\n", id, string(output))
}
