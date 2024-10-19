package fd

import (
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/fhke/loadsheding-go/usage"
)

type fileDescriptorTracker struct {
	procFdDir string
}

func New() usage.Tracker {
	return &fileDescriptorTracker{
		procFdDir: getFdDir(),
	}
}

func (f *fileDescriptorTracker) Utilization() float64 {
	fds, err := os.ReadDir(f.procFdDir)
	if err != nil {
		log.Printf("[ERROR] Error reading file descriptors: %s", err.Error())
		return 0
	}
	return float64(len(fds))
}

func getFdDir() string {
	pid := os.Getpid()
	pidStr := strconv.Itoa(pid)
	return filepath.Join("/proc", pidStr, "fd")
}
