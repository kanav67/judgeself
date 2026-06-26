package sandbox

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

type IsolateMetadata struct {
	Time     float64
	WallTime float64
	Memory   int
	ExitCode int
	ExitSig  int
	Status   string
	Message  string
}

func (s *Sandbox) GetParsedMetadata() *IsolateMetadata {
	meta := &IsolateMetadata{}

	file, err := os.Open(s.GetMetadataPath())
	if err != nil {
		return meta // Return empty struct if missing
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	//todo
	if scanner.Err() != nil {
		return meta // Return empty struct if error reading
	}

	for scanner.Scan() {
		parts := strings.SplitN(scanner.Text(), ":", 2)
		if len(parts) != 2 {
			continue
		}

		key, val := parts[0], parts[1]
		switch key {
		case "time":
			meta.Time, _ = strconv.ParseFloat(val, 64)
		case "time-wall":
			meta.WallTime, _ = strconv.ParseFloat(val, 64)
		case "cg-mem", "max-rss":
			meta.Memory, _ = strconv.Atoi(val)
		case "exitcode":
			meta.ExitCode, _ = strconv.Atoi(val)
		case "exitsig":
			meta.ExitSig, _ = strconv.Atoi(val)
		case "status":
			meta.Status = val
		case "message":
			meta.Message = val
		}
	}
	return meta
}

func (s *Sandbox) GetMetadataPath() string {
	return MetadataFile
}

//todo
func GetUserVerdict(meta *IsolateMetadata) string {
	switch meta.Status {
	case "TO":
		return "TLE"
	case "SG":
		return "RE"
	case "RE":
		return "NZEC"
	case "XX":
		return "IE"
	}
	return "OK"
}