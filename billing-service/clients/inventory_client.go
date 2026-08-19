package clients

import (
	"fmt"
	"net/http"
)

func CheckProductExists(productCode string) error {
	url := fmt.Sprintf("http://localhost:8081/products/%s", productCode)

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("falha de comunicação com o Inventory Service: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("produto %s não existe no estoque", productCode)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("erro inesperado do Inventory Service: status %d", resp.StatusCode)
	}

	return nil
}
