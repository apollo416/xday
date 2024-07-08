package pbuilder

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestDataDirectoryExists(t *testing.T) {
	t.Log("DefaultData: ", DefaultData)
	if _, err := os.Stat(DefaultData.DataDir); os.IsNotExist(err) {
		t.Error(err)
		t.Error("DefaultData Data directory does not exist")
	}

	if _, err := os.Stat(DefaultData.ServicesDir); os.IsNotExist(err) {
		t.Error("DefaultData Services data directory does not exist")
	}

	if _, err := os.Stat(DefaultData.FunctionDataDir); os.IsNotExist(err) {
		t.Error("DefaultData Function data directory does not exist")
	}

	dir := getTestDataDirectory()
	testData := NewData(dir)
	t.Log("testData: ", testData)

	if _, err := os.Stat(testData.DataDir); os.IsNotExist(err) {
		t.Error(err)
		t.Error("testData Data directory does not exist")
	}

	if _, err := os.Stat(testData.ServicesDir); os.IsNotExist(err) {
		t.Error(err)
		t.Error("testData Services directory does not exist")
	}

	if _, err := os.Stat(DefaultData.FunctionDataDir); os.IsNotExist(err) {
		t.Error("testData Function data directory does not exist")
	}
}

func getTestDataDirectory() string {
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	dir = filepath.Join(dir, "testdata")
	return dir
}
