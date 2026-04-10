package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"job_match_api/internal/model"

	"github.com/joho/godotenv"
	"github.com/revrost/go-openrouter"
)

func buildPrompt(data model.AnalyzeRequest) string {
	return fmt.Sprintf(`Tu es un expert en recrutement. Analyse le CV et la description de poste ci-dessous.
Réponds UNIQUEMENT avec un objet JSON valide contenant exactement ces champs :
- "match_score": entier entre 0 et 100
- "summary": une phrase résumant le profil par rapport au poste
- "matched_skills": liste des compétences présentes dans le CV et demandées par l'offre
- "missing_skills": liste des compétences demandées mais absentes du CV
- "recommendations": liste de conseils concrets pour améliorer le CV

Ne mets aucun texte avant ou après le JSON.

--- DESCRIPTION DU POSTE ---
%s

--- CV ---
%s`, data.JobDescription, data.CVText)
}

func Analyze(data model.AnalyzeRequest) (model.AnalysisResult, error) {
	err := godotenv.Load("../../.env")
	if err != nil {
		return model.AnalysisResult{}, fmt.Errorf("error to load .env %v", err)
	}

	apiKey := os.Getenv("OPENROUTER_API_KEY")
	if apiKey == "" {
		return model.AnalysisResult{}, fmt.Errorf("OPENROUTER_API_KEY is not set")
	}

	client := openrouter.NewClient(apiKey)

	prompt := buildPrompt(data)

	resp, err := client.CreateChatCompletion(context.Background(),
		openrouter.ChatCompletionRequest{
			Model: "openai/gpt-oss-120b:free",
			Messages: []openrouter.ChatCompletionMessage{
				{
					Role:    openrouter.ChatMessageRoleUser,
					Content: openrouter.Content{Text: prompt},
				},
			},
		},
	)
	if err != nil {
		return model.AnalysisResult{}, fmt.Errorf("openrouter API call failed: %w", err)
	}

	if len(resp.Choices) == 0 {
		return model.AnalysisResult{}, fmt.Errorf("no response from model")
	}

	rawContent := resp.Choices[0].Message.Content.Text

	var result model.AnalysisResult
	if err := json.Unmarshal([]byte(rawContent), &result); err != nil {
		return model.AnalysisResult{}, fmt.Errorf("failed to parse model response as JSON: %w", err)
	}

	return result, nil
}
