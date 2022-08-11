package transactions

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

type Transaction struct {
	Category string
	Amount   int
}

func NewTransactionClient(url string) TransactionClient {
	return TransactionClient{
		url:        url,
		httpClient: http.Client{},
	}
}

type TransactionClient struct {
	httpClient http.Client
	url        string
}

func (t TransactionClient) GetTransactionsForUser(userId string, start time.Time, end time.Time) ([]Transaction, error) {
	req, err := http.NewRequest(http.MethodGet, t.url, nil)
	if err != nil {
		return []Transaction{}, err
	}

	resp, err := t.httpClient.Do(req)
	if err != nil {
		return []Transaction{}, err
	}

	//need to handle other status codes than just this
	if resp.StatusCode != http.StatusOK {
		return []Transaction{}, errors.New("returned a non 200")
	}

	defer resp.Body.Close()

	var txns []Transaction
	err = json.NewDecoder(resp.Body).Decode(&txns)
	if err != nil {
		return []Transaction{}, err
	}

	return txns, nil
}
