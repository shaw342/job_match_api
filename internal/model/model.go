package model

type AnalyzeRequest struct {
	JobDescription string `json:"job_description"`
	CVText         string `json:"cv_text"`
}
