package clients

import (
	"fmt"
	"net/http"
)

func CheckProductExists(productCode string) error {
	url := fmt.Sprintf("http://localhost:8081/products/%s", productCode)

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to check product existence: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("product %s does not exist in inventory", productCode)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected error from Inventory Service: status %d", resp.StatusCode)
	}

	return nil
}
