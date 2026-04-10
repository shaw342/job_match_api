package llm_test

import (
	"testing"

	"job_match_api/internal/llm"
	"job_match_api/internal/model"
)

func TestAnalyze(t *testing.T) {
	job := model.AnalyzeRequest{
		JobDescription: "Backend Go developer",
		CVText:         "Backend developer with Go and Docker experience",
	}
	result, err := llm.Analyze(job)
	if err != nil {
		t.Log(err)
	}

	t.Log(result)
}
