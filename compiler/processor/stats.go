package processor

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Stats struct {
	CPU int64
	Mem int64
}

// containerStats should be run in a go routine. This function
// waits for the required files to be created and starts to collect metrics for them.
func (job *Job) containerStats(containerID string) {
	// TODO: What can i do with the job.logAndError to pass back the error back
	// to the job?
	defer job.wg.Done()

	var stats = &Stats{}

	var limit = 250
	var count = 0

	cpuFile := fmt.Sprintf("/sys/fs/cgroup/system.slice/docker-%s.scope/cpu.stat", containerID)
	memFile := fmt.Sprintf("/sys/fs/cgroup/system.slice/docker-%s.scope/memory.current", containerID)
	// Wait for the required files to be created.
	for {
		if _, err := os.Stat(cpuFile); err == nil ||
			count > limit {
			break
		}

		time.Sleep(1 * time.Millisecond)
		count++
	}

	// Iterate until file is no longer available
	for {
		cpu, err := getFieldFromFile(cpuFile, "usage_usec")
		if err != nil {
			job.Stats = stats
			return
		}
		if stats.CPU < cpu {
			stats.CPU = cpu
		}

		mem, err := getMemoryFieldFromFile(memFile)
		if err != nil {
			job.Stats = stats
			return
		}
		if stats.Mem < mem {
			stats.Mem = mem
		}
		time.Sleep(1 * time.Millisecond)
	}
}

func getFieldFromFile(filePath, field string) (int64, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) == 2 && fields[0] == field {
			return strconv.ParseInt(fields[1], 10, 64)
		}
	}

	return 0, fmt.Errorf("field not found")
}

func getMemoryFieldFromFile(filePath string) (int64, error) {
	bs, err := os.ReadFile(filePath)
	if err != nil {
		return 0, err
	}

	return strconv.ParseInt(strings.TrimSpace(string(bs)), 10, 64)
}
