package sandbox

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const CompileTimeLimit = 5.0
const CompileMemoryLimit = 256 * 1024 // 256 MB


func (s *Sandbox) CompileCode(sourceCode []byte, language, additionalFilesDir string, additionalArgs []string) ([]byte, error) {
	s.ReInitialize()

	sourcePath := filepath.Join(s.BoxDir, "source")
	_ = os.WriteFile(sourcePath, sourceCode, 0644)

	compileCmd, exists := GetCompileCommand(language, "source", "outputBin");
	if !exists {
		return nil, fmt.Errorf("Unsupported checker language: %s", language)
	}

	if additionalFilesDir != "" {
		CopyDirFiles(additionalFilesDir, s.BoxDir)
	}

	compileConfig:= NewIsolateConfig()

	args := s.baseIsolateArgs()
	args = append(args, compileConfig.toArgs()...)
	args = append(args,
		"-M", s.GetMetadataPath(),
		"-E", "HOME=/tmp",
		"-E", "PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
		"-d", "/etc:noexec",
		"--run", "--", compileCmd, strings.Join(additionalArgs, " "),
	)

	cmd := exec.Command("isolate", args...)
	output, err := cmd.CombinedOutput()

	if err != nil {
		return nil, fmt.Errorf("checker compile failed: %v, output: %s", err, string(output))
	}

	//todo check metadata for errors, tle etc
	// metaData := s.GetParsedMetadata()
	
	outputBin, _ := os.ReadFile(filepath.Join(s.BoxDir, "outputBin"))
	
	return outputBin, nil
}