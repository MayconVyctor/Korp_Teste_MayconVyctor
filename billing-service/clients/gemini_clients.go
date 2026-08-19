package clients

import (
	"context"
	"fmt"
	"os"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"

	"github.com/MayconVyctor/Korp_Teste_MayconVyctor/billing-service/models"
)

func AnalyzeInvoice(invoice models.Invoice) (string, error) {

	promptText := fmt.Sprintf("You are a senior sales analyst. The customer has made the following purchase (Invoice ID: %d, Status: %s):\n", invoice.ID, invoice.Status)
	for _, item := range invoice.Items {
		promptText += fmt.Sprintf("- Product Code: %s | Quantity: %d\n", item.ProductCode, item.Quantity)
	}

	promptText += "\nBased on these items, write a short paragraph (maximum 3 lines) suggesting which would be the next ideal product to offer to this customer for an upsell, Please, write your response in English"

	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("the generated API key is empty, please set the GEMINI_API_KEY environment variable with your valid API key")
	}

	ctx := context.Background()

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return "", fmt.Errorf("failed to create Gemini client: %v", err)
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-3.6-flash")

	resp, err := model.GenerateContent(ctx, genai.Text(promptText))
	if err != nil {
		return "", fmt.Errorf("failed to generate analysis with Gemini: %v", err)
	}

	if len(resp.Candidates) > 0 && len(resp.Candidates[0].Content.Parts) > 0 {
		aiText := fmt.Sprintf("%v", resp.Candidates[0].Content.Parts[0])
		return aiText, nil
	}

	return "It was not possible to generate an analysis for this invoice", nil
}
