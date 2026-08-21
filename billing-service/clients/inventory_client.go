package clients

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ProductResponse struct {
	Code        string `json:"code"`
	Description string `json:"description"`
	Balance     int    `json:"balance"`
}

func CheckProductExists(productCode string) (string, error) {
	url := fmt.Sprintf("http://localhost:8081/products/%s", productCode)

	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to check product existence: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return "", fmt.Errorf("product %s does not exist in inventory", productCode)
	}

	var prod ProductResponse
	if err := json.NewDecoder(resp.Body).Decode(&prod); err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected error from Inventory Service: status %d", resp.StatusCode)
	}

	return prod.Description, nil
}
