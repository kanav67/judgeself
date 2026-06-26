package worker

import (
	job "execution-engine/Job"
	"execution-engine/sandbox"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

type Worker struct {
	BoxID    int
	Sandbox  *sandbox.Sandbox
	LocalDir string
	Job      *job.Job
}

const (
	StdinFile       = "stdin.txt"
	StdoutFile      = "stdout.txt"
	StderrFile      = "stderr.txt"
	MetadataFile    = "metadata.txt"
	LocalTemp       = "/tmp/execution-engine"
	CheckerBinFile  = "checkerBin"
	UserCodeBinFile = "userCodeBin"
	OutputDir       = "output"
)

func NewWorker(boxID int, job *job.Job) (*Worker, error) {
	w := &Worker{
		BoxID:    boxID,
		LocalDir: filepath.Join(LocalTemp, job.SubmissionData.SubmissionID),
		Job:      job,
	}

	sandbox, err := sandbox.NewSandbox(boxID)
	if err != nil {
		//todo
		panic("Failed to initialize sandbox: " + err.Error())
	}
	w.Sandbox = sandbox

	err = os.MkdirAll(w.LocalDir, 0755)
	if err != nil {
		return nil, fmt.Errorf("Failed to create destination directory: %w", err)
	}

	return w, nil
}

func (w *Worker) ExecuteWorker() error {
	userCodeSource := []byte(w.Job.SubmissionData.SourceCode)
	userCodeLanguage := w.Job.SubmissionData.SourceLanguage

	userCodeBin, _ := w.Sandbox.CompileCode(userCodeSource, userCodeLanguage, "", nil)

	//todo handle err

	checkerLanguage := w.Job.ProblemData.CheckerLanguage
	checkerSource, _ := os.ReadFile(w.Job.ProblemData.GetCheckerSourcePath())

	checkerBin, _ := w.Sandbox.CompileCode(checkerSource, checkerLanguage, w.Job.ProblemData.GetResourcesFolder(), nil)

	config := sandbox.NewIsolateConfig(func(ic *sandbox.IsolateConfig) {
		ic.TimeLimit = float64(w.Job.ProblemData.TimeLimit) / 1000
		ic.MemoryLimit = w.Job.ProblemData.MemoryLimit
		ic.WallTimeLimit = ic.TimeLimit * 2
	})

	for test := 1; test <= w.Job.ProblemData.TestCount; test++ {
		testInputPath := w.Job.ProblemData.GetTestFilePath(test)
		testAnswerPath := w.Job.ProblemData.GetAnswerFilePath(test)
		testOutputPath := w.GetOutputFilePath(test)

		// userMeta :=
		w.Sandbox.ExecuteCode(userCodeBin, userCodeLanguage, testInputPath, testOutputPath, nil, config)
		// checkerMeta :=
		w.Sandbox.ExecuteCode(checkerBin, checkerLanguage, "", "", []string{testInputPath, testOutputPath, testAnswerPath}, config)
	}

	return nil
}

func (w *Worker) GetOutputDir() string {
	return filepath.Join(w.LocalDir, OutputDir)
}

func (w *Worker) GetOutputFilePath(test int) string {
	length := len(strconv.Itoa(w.Job.ProblemData.TestCount))
	return filepath.Join(w.GetOutputDir(), fmt.Sprintf("%0*d.out", length, test))
}
