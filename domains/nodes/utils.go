package nodes

import "github.com/v4lproik/simple-blockchain-quickstart/common/models"

func txsMapToArr(txsMap map[models.TransactionId]models.Transaction) []models.Transaction {
	arr := make([]models.Transaction, len(txsMap))

	i := 0
	for _, tx := range txsMap {
		arr[i] = tx
		i++
	}

	return arr
}
