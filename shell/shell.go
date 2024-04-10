// Package shell wraps exec.Command
package shell

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

var (
	// CommandEnvInvalidErr means environment passed have errors
	CommandEnvInvalidErr = errors.New("command environment invalid")
)

// Run Shell command
// run this command will block until task finished
func Run(command string, envs []string) ([]byte, error) {
	if err := validateEnvs(envs); err != nil {
		return nil, err
	}
	var stdout, stderr bytes.Buffer
	name, args := getExecAndArgsFromCommand(command)
	cmd := exec.Command(name, args...)
	cmd.Env = append(cmd.Env, envs...) // adding environments
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	// run command
	err := cmd.Run()

	stderrContent := string(stderr.Bytes())
	if stderrContent != "" {
		err = errors.Join(err, errors.New(stderrContent))
	}

	out := stdout.Bytes()

	return out, err
}

// RunWithCtx pass ctx to Run Func
func RunWithCtx(ctx context.Context, command string, envs []string) ([]byte, error) {
	var err error
	var output []byte
	c := make(chan struct{}, 1)

	go func() {
		output, err = Run(command, envs)
		c <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-c:
		// done the Run func
		return output, err
	}
}

func validateEnv(env string) bool {
	s := strings.Split(env, "=")
	if len(s) == 2 {
		return true
	}
	return false
}

func validateEnvs(envs []string) error {
	var err error

	for _, env := range envs {
		if !validateEnv(env) {
			err = errors.Join(fmt.Errorf("%w, %s", CommandEnvInvalidErr, env))
		}
	}
	return err
}

func getExecAndArgsFromCommand(command string) (execName string, args []string) {
	r := strings.Split(command, " ")
	r2 := make([]string, 0)
	// delete empty element
	for _, v := range r {
		if v != "" {
			r2 = append(r2, v)
		}
	}
	return r2[0], r2[1:]
}
