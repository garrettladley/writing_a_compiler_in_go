package parser

import (
	"fmt"
	"strings"
)

//nolint:unused
var traceLevel int = 0

//nolint:unused
const traceIdentPlaceholder string = "\t"

//nolint:unused
func identLevel() string {
	return strings.Repeat(traceIdentPlaceholder, traceLevel-1)
}

//nolint:unused
func tracePrint(fs string) {
	fmt.Printf("%s%s\n", identLevel(), fs)
}

//nolint:unused
func incIdent() { traceLevel = traceLevel + 1 }

//nolint:unused
func decIdent() { traceLevel = traceLevel - 1 }

//nolint:unused
func trace(msg string) string {
	incIdent()
	tracePrint("BEGIN " + msg)
	return msg
}

//nolint:unused
func untrace(msg string) {
	tracePrint("END " + msg)
	decIdent()
}
