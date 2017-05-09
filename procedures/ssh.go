package procedures

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

type SSHCommander struct {
	User         string
	IP           string
	Port         int
	IdentityFile string
	Arguments    []string
}

func (ssh *SSHCommander) Command(cmd ...string) *exec.Cmd {
	var args []string = make([]string, 0)
	if ssh.IdentityFile != "" {
		args = append(args, "-i", ssh.IdentityFile)
	}

	if len(ssh.Arguments) > 0 {
		args = append(args, ssh.Arguments...)
	}

	args = append(args, fmt.Sprintf("%s@%s:%d", ssh.User, ssh.IP, ssh.Port))

	args = append(args, cmd...)

	return exec.Command("ssh", args...)
}

func (ssh *SSHCommander) Procedure(commands ...string) ([]string, error) {
	var err error
	var out []string = make([]string, 0)
	var args []string = make([]string, 0)
	if ssh.IdentityFile != "" {
		args = append(args, "-i", ssh.IdentityFile)
	}

	if len(ssh.Arguments) > 0 {
		args = append(args, ssh.Arguments...)
	}

	args = append(args, fmt.Sprintf("%s@%s:%d", ssh.User, ssh.IP, ssh.Port))

	cmd := exec.Command("ssh", args...)
	reader, err := cmd.StdoutPipe()
	scanner := bufio.NewScanner(reader)
	defer reader.Close()
	err = scanner.Err()
	if err == nil {
		for _, command := range commands {
			cmd.Stdin = strings.NewReader(command)
			for scanner.Scan() {
				err = scanner.Err()
				if err != nil {
					break
				}
				out = append(out, string(scanner.Bytes()))
			}
		}
	}

	// Graceful ssh exit
	cmd.Stdin = strings.NewReader("exit")

	time.Sleep(2 * time.Second)

	if cmd.Process.Pid > 0 {
		// Forced ssh exit
		cmd.Process.Kill()
	}

	return out, err
}

func SShGoLangInstall(user string, ip string, port int, identityFilePath string) {
	if port == 0 {
		port = 22
	}
	commander := SSHCommander{
		User:         user,
		IP:           ip,
		Port:         port,
		IdentityFile: identityFilePath,
	}

	command := []string{
		"apt-get",
		"install",
		"-y",
		"golang-go",
	}

	cmd := commander.Command(command...)
	reader, err := cmd.StdoutPipe()
	reader.Close()
	if err != nil {
		fmt.Fprintln(os.Stderr, "There was an error running SSH command: ", err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		err = scanner.Err()
		if err == nil {
			fmt.Printf("%s\n", string(scanner.Bytes()))
		} else {
			fmt.Fprintln(os.Stderr, "Unable to read SSH output: ", err)
		}
	}
	fmt.Fprintln(os.Stdout, fmt.Sprintf("SSH Connection to %s:%d with user %s closed!!", ip, port, user))
}
