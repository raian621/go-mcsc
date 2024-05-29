package minecraft

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
)

type Console struct {
	stdin  *bufio.Writer
	stdout *bufio.Reader
	stderr *bufio.Reader
}

func NewConsole(cmd *exec.Cmd) (*Console, error) {
	var (
		stdin  io.WriteCloser
		stdout io.ReadCloser
		stderr io.ReadCloser
		err    error
	)

	if stdin, err = cmd.StdinPipe(); err != nil {
		return nil, err
	}
	if stdout, err = cmd.StdoutPipe(); err != nil {
		return nil, err
	}
	if stderr, err = cmd.StdoutPipe(); err != nil {
		return nil, err
	}

	return &Console{
		stdin:  bufio.NewWriter(stdin),
		stdout: bufio.NewReader(stdout),
		stderr: bufio.NewReader(stderr),
	}, nil
}

func (c Console) SendCommand(cmd string) error {
	_, err := c.stdin.WriteString(
		fmt.Sprintf("%s\r\n", cmd),
	)

	if err != nil {
		return err
	}

	return c.stdin.Flush()
}

func (c Console) ReadLine() (string, error) {
	return c.stdout.ReadString('\n')
}

func (c Console) ReadError() (string, error) {
	return c.stderr.ReadString('\n')
}
