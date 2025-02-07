// testRegen
// This programme helps port test cases from one output format to
// another.
//
// The porting consists of:
//  1. parsing the current test file for email test cases
//  2. running the parser to generate output for each test case
//     in a format suitable format for future tests
//
// No validation of the correctness of the output (other than each test
// file being parsed ok) is made. Manual checking is strongly advised.
//
// Some adjustments to the output need to be made, including
// reformatting by a text editor and the "TestHeaders" tests.
//
// RCL 07 February 2025
package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"text/template"
	"time"

	"github.com/mnako/letters"
	"github.com/mnako/letters/email"
	"github.com/sanity-io/litter"
)

var (
	testFile              string         = "../letters_test.go"
	extractFileNameRegexp *regexp.Regexp = regexp.MustCompile(`fp := "(.*txt)"`)
	extractFuncRegexp     *regexp.Regexp = regexp.MustCompile(`^func (Test.*)\(t \*testing.T\) {\s*`)
)

// results
var (
	fileNames []string = []string{}
	funcs     []string = []string{}
)

// grabTop gets the top lines from testFile, until the first `func Test...`
func grabTop() error {
	c, err := ioutil.ReadFile(testFile)
	if err != nil {
		return err
	}
	s := string(c)
	ss := strings.Split(s, "\n")
	for _, l := range ss {
		if extractFuncRegexp.MatchString(l) {
			break
		}
		fmt.Println(l)
	}
	return nil
}

// gatherTestCases grabs the test names and files in testFile
func gatherTestCases() error {
	f, err := os.Open(testFile)
	if err != nil {
		return err
	}
	buf := bufio.NewScanner(f)
	for buf.Scan() {
		if te := extractFileNameRegexp.FindAllStringSubmatch(buf.Text(), -1); te != nil {
			fileNames = append(fileNames, te[0][1])
			continue
		}
		if te := extractFuncRegexp.FindAllStringSubmatch(buf.Text(), -1); te != nil {
			funcs = append(funcs, te[0][1])
			continue
		}
	}
	if len(fileNames) != len(funcs) {
		return errors.New("number of fileNames != number of funcs")
	}
	return nil
}

// return email parsed from fp path
func parseEmail(fp string) (*email.Email, error) {
	p := letters.NewParser()
	o, err := os.Open(fp)
	if err != nil {
		return nil, err
	}
	parsedEmail, err := p.Parse(o)
	if err != nil {
		return nil, err
	}
	return parsedEmail, nil
}

// dump email
func dumpEmail(t any) string {
	// https://github.com/sanity-io/litter/issues/12#issuecomment-1144643251
	timeType := reflect.TypeOf(time.Time{})
	litter.Config.DumpFunc = func(v reflect.Value, w io.Writer) bool {
		if v.Type() != timeType {
			return false
		}
		t := v.Interface().(time.Time)
		t = t.In(time.UTC)
		t.Truncate(time.Second) // only need accuracy to second for email parsing
		fmt.Fprintf(
			w,
			`(time.Date(%d, %d, %d, %d, %d, %d, %d, time.UTC))`,
			t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(),
		)
		return true
>>>>>>> c96ad5a (utils: tweak testRegen to add header)
	}
	litter.Config.FieldExclusions = regexp.MustCompile("^(Reader|Encoding)$")
	litter.Config.DisablePointerReplacement = true
	return litter.Sdump(t)
}

var isTextRegexp *regexp.Regexp = regexp.MustCompile(`^\s+(Text|EnrichedText|HTML):`)
var byteLineLen int = 80

func formatDump(s string) string {
	outputLines := []string{}
	lines := strings.Split(s, "\n")
	inBytes := false
	storedString := ""
	for _, l := range lines {
		switch {
		case isTextRegexp.MatchString(l):
			l = strings.ReplaceAll(l, `\n`, `\n`+"\" +\n\"")
			outputLines = append(outputLines, l)
			continue
		case strings.Contains(l, "Data: []uint8{"):
			inBytes = true
			l = strings.ReplaceAll(l, "Data: []uint8{", "Data: []byte{")
			outputLines = append(outputLines, l)
			continue
		case strings.Contains(l, "},"):
			if inBytes {
				inBytes = false
				if storedString != "" {
					outputLines = append(outputLines, storedString)
					storedString = ""
				}
			}
			outputLines = append(outputLines, l)
			continue
		case inBytes:
			fragment := strings.TrimSpace(l)
			storedString += fragment + " "
			if len(storedString) > byteLineLen {
				outputLines = append(outputLines, storedString)
				storedString = ""
			}
			continue
		}
		outputLines = append(outputLines, l)
	}
	return strings.Join(outputLines, "\n")
>>>>>>> 2092f58 (utils: tweak testRegen.go)
}

var tpl string = `
func {{ .Func }}(t *testing.T) {
	fp := "{{ .Filename }}"
	expectedEmail := {{ .ParseResults }}
	testEmailFromFile(t, fp, expectedEmail)
}
`

var t = template.Must(template.New("tpl").Parse(tpl))

// print the output via a template
func printOutput(fn, fileName, parseResults string) error {
	return t.Execute(os.Stdout, struct {
		Func         string
		Filename     string
		ParseResults string
	}{
		Func:         fn,
		Filename:     fileName,
		ParseResults: parseResults,
	})
}

func main() {
	testFilePath := "../"
	err := gatherTestCases()
	if err != nil {
		fmt.Println("gather tests error:", err)
		os.Exit(1)
	}
	err = grabTop()
	if err != nil {
		fmt.Println("could not grab top of file error:", err)
		os.Exit(1)
	}
	for i := 0; i < len(funcs); i++ {
		email, err := parseEmail(filepath.Join(testFilePath, fileNames[i]))
		if err != nil {
			os.Stderr.WriteString(fmt.Sprintf("parseEmail error for %s: %s", fileNames[1], err))
			os.Exit(1)
		}
		s := dumpEmail(email)
		s = formatDump(s)
		err = printOutput(funcs[i], fileNames[i], s)
		if err != nil {
			fmt.Println("printOut err:", err)
			os.Exit(1)
		}
	}
}
