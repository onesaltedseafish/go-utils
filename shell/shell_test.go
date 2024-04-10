package shell

import (
	"context"
	"errors"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// get err string
// return "" for nil err
func getErrString(err error) string{
	if err == nil{
		return ""
	}
	return err.Error()
}

func TestValidateEnv(t *testing.T) {
	testcases := []struct {
		Target     string
		WantResult bool
	}{
		{"a=b", true},
		{"a = b", true},
		{"a==b", false},
		{"a", false},
	}

	for _, testcase := range testcases {
		assert.Equal(t, testcase.WantResult, validateEnv(testcase.Target))
	}
}

func TestValidateEnvs(t *testing.T) {
	testcases := []struct {
		Envs    []string
		MetaErr error
	}{
		{nil, nil},
		{[]string{}, nil},
		{[]string{"a=b", "c=d"}, nil},
		{[]string{"a", "b"}, CommandEnvInvalidErr},
		{[]string{"a=b", "c"}, CommandEnvInvalidErr},
	}

	for _, testcase := range testcases {
		if testcase.MetaErr == nil {
			assert.Empty(t, validateEnvs(testcase.Envs))
		} else {
			err := validateEnvs(testcase.Envs)
			if !errors.Is(err, testcase.MetaErr) {
				t.Errorf("err is not derived from CommandEnvInvalidErr")
			}
		}
	}
}

func TestGetExecAndArgsFromCommand(t *testing.T) {
	testcases := []struct {
		Command      string
		WantExecName string
		WantArgs     []string
	}{
		{"ls ", "ls", []string{}},
		{"ls  -l -h", "ls", []string{"-l", "-h"}},
	}

	for _, testcase := range testcases {
		tmpExecName, tmpArgs := getExecAndArgsFromCommand(testcase.Command)
		assert.Equal(t, testcase.WantExecName, tmpExecName)
		assert.Equal(t, testcase.WantArgs, tmpArgs)
	}
}

func TestRun(t *testing.T) {
	if runtime.GOOS != "linux" { // tests run only under linux
		return
	}
	testcases := []struct {
		Command       string
		WantOutput    string
		WantErrString string
	}{
		{"echo test", "test\n", ""},
		{"cat /tmp", "", "exit status 1\ncat: /tmp: Is a directory\n"},
		{"ls /t1p", "", "exit status 2\nls: cannot access '/t1p': No such file or directory\n"},
	}

	for _, testcase := range testcases {
		output, err := Run(testcase.Command, nil)
		assert.Equal(t, testcase.WantOutput, string(output))
		assert.Equal(t, testcase.WantErrString, getErrString(err))
	}
}

// finds shell.Run will return an empty error message error
// so we need a specified test case
func TestRunErr(t *testing.T) {
	if runtime.GOOS != "linux" { // tests run only under linux
		return
	}
	_, err := Run("echo 2", nil)
	if err != nil {
		t.Errorf("the command will always success, but got err: %v", err)
	}
}

func TestRunWithCtx(t *testing.T) {
	if runtime.GOOS != "linux" {
		return
	}
	timeoutCtx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	testcases := []struct {
		Ctx           context.Context
		Command       string
		Envs          []string
		WantOutput    string
		WantErrString string
	}{
		{timeoutCtx, "sleep 2", []string{}, "", "context deadline exceeded"},
		{context.Background(), "echo 2", []string{}, "2\n", ""},
	}

	for _, testcase := range testcases {
		output, err := RunWithCtx(testcase.Ctx, testcase.Command, testcase.Envs)
		assert.Equal(t, testcase.WantOutput, string(output))
		assert.Equal(t, testcase.WantErrString, getErrString(err))
	}
}
