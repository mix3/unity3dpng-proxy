package main

import (
	"bytes"
	"os/exec"
	"strings"

	"github.com/Sirupsen/logrus"
)

type Git struct {
	GitCmd  string
	Repo    string
	WorkDir string
	Logger  *logrus.Logger
}

func NewGit(gitCmd, repo, workDir string, logger *logrus.Logger) *Git {
	return &Git{
		GitCmd:  gitCmd,
		Repo:    repo,
		WorkDir: workDir,
		Logger:  logger,
	}
}

func (g *Git) Run(arg string, more ...string) ([]byte, []byte, error) {
	args := append([]string{
		g.GitCmd,
		"--git-dir",
		g.WorkDir + "/.git",
		"--work-tree",
		g.WorkDir,
		arg,
	}, more...)
	logger.Info(strings.Join(args, " "))
	cmd := exec.Command(args[0], args[1:]...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return nil, nil, err
	}
	return stdout.Bytes(), stderr.Bytes(), nil
}

func (g *Git) Clone() ([]byte, []byte, error) {
	args := append([]string{
		g.GitCmd,
		"clone",
		g.Repo,
		g.WorkDir,
	})
	logger.Info(strings.Join(args, " "))
	cmd := exec.Command(args[0], args[1:]...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return nil, nil, err
	}
	return stdout.Bytes(), stderr.Bytes(), nil
}
