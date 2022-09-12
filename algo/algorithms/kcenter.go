package algorithms

import (
	"fmt"
	"io/ioutil"
	"os/exec"

	"os"

	"regexp"
	"strings"
)

//Kcenter package struct
type Kcenter struct{}

//Method parses input and executes the the requested algorithm
func (k Kcenter) run(request Parameters) (string, string, error) {
	if len(request.Parameters) < 1 {
		return "Format not specified, please refer to documentation", "0s", nil
	}
	//Input parsing
	tfile, err := ioutil.TempFile("", "graph")
	if err != nil {
		return "", "", err
	}
	defer os.Remove(tfile.Name())
	if request.Parameters[0] == "pmed" {
		stringLines := strings.Replace(request.Input, "\n", " ", -1)
		reg := regexp.MustCompile(`(\w+) (\w+) (\w+)`)
		stringLines = reg.ReplaceAllString(stringLines, "$0\n")
		tfile.WriteString(stringLines)
		err = os.Chmod(tfile.Name(), 0777)
		if err != nil {
			return "", "", err
		}
	} else if request.Parameters[0] == "graph" {
		if len(request.Parameters) < 2 {
			return "Number of centers missing, please refer to documentation", "", nil
		}
		reg := regexp.MustCompile(`(.*?)(\[[ Nodes ^\]]*\].*)(\[[ Edges ^\]]*\].*)`)
		stringLines := reg.FindAllStringSubmatch(request.Input, -1)
		tfile.WriteString(stringLines[0][1])
		reg = regexp.MustCompile(`(\(.*?\))`)
		stringLine := reg.ReplaceAllString(stringLines[0][2], "\n$1")
		tfile.WriteString("\n" + stringLine)
		stringLine = reg.ReplaceAllString(stringLines[0][3], "\n$1")
		tfile.WriteString("\n" + stringLine)
		err = os.Chmod(tfile.Name(), 0777)
		if err != nil {
			return "", "", err
		}
	} else if request.Parameters[0] == "pajek" {
		if len(request.Parameters) < 2 {
			return "Number of centers missing, please refer to documentation", "", nil
		}
		reg := regexp.MustCompile(`(\*vertices\s\d*)\s*(\*arcs)\s*(.*$)`)
		stringLines := reg.FindAllStringSubmatch(request.Input, -1)
		tfile.WriteString(stringLines[0][1] + "\n")
		tfile.WriteString(stringLines[0][2] + "\n")
		reg = regexp.MustCompile(`(\w+) (\w+) (\w+)`)
		stringLine := reg.ReplaceAllString(stringLines[0][3], "$0\n")
		tfile.WriteString(stringLine)
	} else {
		return "Format not available, please refer to documentation", "", nil
	}
	//Algorithm execution
	var b []byte
	err = nil
	if request.Parameters[0] == "pmed" {
		cmd := exec.Command("./kcenter", "-i", tfile.Name(), "-m", request.Algorithm, "-f", request.Parameters[0])
		b, err = cmd.CombinedOutput()
	} else {
		cmd := exec.Command("./kcenter", "-i", tfile.Name(), "-m", request.Algorithm, "-f", request.Parameters[0], "-c", request.Parameters[1])
		b, err = cmd.CombinedOutput()
	}
	if err != nil {
		return "Error executing algorithm: " + err.Error(), "", nil
	}
	st := fmt.Sprintf("%s", b)
	re := regexp.MustCompile(`\r?\n`)
	st = re.ReplaceAllString(st, "")
	var outs []string = strings.Split(st, ",")
	return outs[0], outs[1] + "ms", nil
}
