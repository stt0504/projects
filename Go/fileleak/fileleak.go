//go:build !solution

package fileleak

import (
	"os"
	"path/filepath"
)

type testingT interface {
	Errorf(msg string, args ...interface{})
	Cleanup(func())
}

func VerifyNone(t testingT) {
	startFiles, err := getOpenFiles()
	if err != nil {
		t.Errorf("Failed to get open files: %v", err)
		return
	}

	t.Cleanup(func() {
		finalFiles, err := getOpenFiles()
		if err != nil {
			t.Errorf("Failed to get open files: %v", err)
			return
		}

		leaked := 0

		for link := range finalFiles {
			if _, exists := startFiles[link]; !exists {
				leaked++
			}
		}

		if leaked > 0 {
			t.Errorf("File leak detected")
		}
	})
}

func getOpenFiles() (map[string]struct{}, error) {
	fileSet := make(map[string]struct{})

	fds, err := os.ReadDir("/proc/self/fd")
	if err != nil {
		return fileSet, err
	}

	for _, fd := range fds {
		fdPath, err := os.Readlink(filepath.Join("/proc/self/fd", fd.Name()))
		if err == nil {
			fileSet[fd.Name()+fdPath] = struct{}{}
		}
	}
	return fileSet, nil
}
