package mysql

import (
	"app/models"
)

//just stubbing to show, this would probably be done with a sql reader and writer
type Store struct {
}

func (s Store) GetUserByUID(id string) (models.User, error) {
	//join on userID on the account table and budget table to return the user....
	return models.User{}, nil
}

func (s Store) UpdateUser(models.User) error {
	//update those tables with a transaction to avoid partial failures
	return nil
}

func (s Store) CreateBudget(userID string, buget models.Budget) (models.Budget, error) {
	//update those tables with a transaction to avoid partial failures
	return models.Budget{}, nil
}

func (s Store) UpdateBudgetCategory(id string, category string) (models.Budget, error) {
	//update those tables with a transaction to avoid partial failures
	return models.Budget{}, nil
}

func (s Store) DeleteBudgetById(id string) error {
	//update those tables with a transaction to avoid partial failures
	return nil
}
