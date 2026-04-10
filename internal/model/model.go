package model

type AnalyzeRequest struct {
	JobDescription string `json:"job_description"`
	CVText         string `json:"cv_text"`
}

type AnalysisResult struct {
	MatchScore      int      `json:"match_score"`
	Summary         string   `json:"summary"`
	MatchedSkills   []string `json:"matched_skills"`
	MissingSkills   []string `json:"missing_skills"`
	Recommendations []string `json:"recommendations"`
}
