package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func checkDir(dir string) (bool, error) {
	_, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func copyFile(sourceFile, destFile string) error {
	input, err := ioutil.ReadFile(sourceFile)
	if err != nil {
		return err
	}

	if mkdirIfNotExist(filepath.Dir(destFile)) != nil {
		return err
	}

	err = ioutil.WriteFile(destFile, input, 0644)
	if err != nil {
		return err
	}
	return nil
}

func mkdirIfNotExist(path string) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		if err = os.MkdirAll(path, 0700); err != nil {
			return err
		}
	}
	return nil
}

func gitCheck(cloneToDir, filePath string) error {
	stringToCheck := fmt.Sprintf("A\t%s\n", filePath)

	output, err := exec.Command("git", "-C", cloneToDir, "log", "--name-status", "HEAD^..HEAD").CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to \"git -C cloneToDir log --name-status HEAD^..HEAD\": %s", err)
	}
	if !strings.Contains(string(output), stringToCheck) {
		return errors.New("mismatch git check")
	}
	return nil
}

func gitPush(cloneToDir, branchName string) error {
	output, err := exec.Command("git", "-C", cloneToDir, "push", "--set-upstream", "origin", branchName).CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to \"git -C %s push\": %s", cloneToDir, output)
	}
	return nil
}
