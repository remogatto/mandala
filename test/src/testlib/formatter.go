// +build android

package testlib

import (
	"github.com/remogatto/gorgasm"
	"github.com/remogatto/prettytest"
	"log"
	"path/filepath"
)

const formatTag = "%s\t"

var (
	labelFAIL         = red("F")
	labelMUSTFAIL     = green("EF")
	labelPASS         = green("OK")
	labelPENDING      = yellow("PE")
	labelNOASSERTIONS = yellow("NA")
)

func green(text string) string {
	return "\033[32m" + text + "\033[0m"
}

func red(text string) string {
	return "\033[31m" + text + "\033[0m"
}

func yellow(text string) string {
	return "\033[33m" + text + "\033[0m"
}

/*TDDFormatter is a very simple TDD-like formatter.

Legend:

* F  - Test Failed

* OK - Test Passed

* EF - An Expected Failure occured

* NA - Not Assertions found

* PE - Pending test
*/
type TDDFormatter struct {
	log *log.Logger
}

func NewTDDFormatter() *TDDFormatter {
	return &TDDFormatter{
		log: log.New(gorgasm.AndroidWriter{}, "[gorgasm-test] ", 0),
	}
}

func (formatter *TDDFormatter) PrintSuiteInfo(suite *prettytest.Suite) {
	formatter.log.Printf("%s.%s:\n\n", suite.Package, suite.Name)
}

func (formatter *TDDFormatter) PrintStatus(testFunc *prettytest.TestFunc) {
	callerName := testFunc.Name
	switch testFunc.Status {
	case prettytest.STATUS_FAIL:
		formatter.log.Printf(formatTag+"%-30s(%d assertion(s))\n", labelFAIL, callerName, len(testFunc.Assertions))
	case prettytest.STATUS_MUST_FAIL:
		formatter.log.Printf(formatTag+"%-30s(%d assertion(s))\n", labelMUSTFAIL, callerName, len(testFunc.Assertions))
	case prettytest.STATUS_PASS:
		formatter.log.Printf(formatTag+"%-30s(%d assertion(s))\n", labelPASS, callerName, len(testFunc.Assertions))
	case prettytest.STATUS_PENDING:
		formatter.log.Printf(formatTag+"%-30s(%d assertion(s))\n", labelPENDING, callerName, len(testFunc.Assertions))
	case prettytest.STATUS_NO_ASSERTIONS:
		formatter.log.Printf(formatTag+"%-30s(%d assertion(s))\n", labelNOASSERTIONS, callerName, len(testFunc.Assertions))

	}
}

func (formatter *TDDFormatter) PrintErrorLog(logs []*prettytest.Error) {
	if len(logs) > 0 {
		currentTestFuncHeader := ""
		for _, error := range logs {
			if currentTestFuncHeader != error.TestFunc.Name {
				formatter.log.Printf("\n%s:\n", error.TestFunc.Name)
			}
			filename := filepath.Base(error.Assertion.Filename)
			formatter.log.Printf("\t(%s:%d) %s\n", filename, error.Assertion.Line, error.Assertion.ErrorMessage)
			currentTestFuncHeader = error.TestFunc.Name
		}
	}
}

func (formatter *TDDFormatter) PrintFinalReport(report *prettytest.FinalReport) {
	formatter.log.Printf("%d tests, %d passed, %d failed, %d expected failures, %d pending, %d with no assertions\n",
		report.Total(), report.Passed, report.Failed, report.ExpectedFailures, report.Pending, report.NoAssertions)
}

func (formatter *TDDFormatter) AllowedMethodsPattern() string {
	return "^Test.*"
}
