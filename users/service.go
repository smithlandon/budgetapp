package users

import (
	"app/models"
	"app/transactions"
	"time"
)

type UserService struct {
	repo              repo
	transactionClient transactionClient
}

func NewUserService(repo repo, transtransactionClient transactionClient) UserService {
	return UserService{
		repo:              repo,
		transactionClient: transtransactionClient,
	}
}

//whatever our DB implementation is -> be it something like prepared statements or ORM like gorm
//in this case, I'd just go prepared statements
//this repo function will join on all of our important tables
type repo interface {

	//user
	GetUserByUID(id string) (models.User, error)
	UpdateUser(models.User) error

	//budget
	CreateBudget(userId string, buget models.Budget) (models.Budget, error)
	UpdateBudgetCategory(id string, category string) (models.Budget, error)
	DeleteBudgetById(id string) error
}

//this interface could just be an http client.
type transactionClient interface {
	GetTransactionsForUser(userId string, start time.Time, end time.Time) ([]transactions.Transaction, error)
}

func (u UserService) GetUserByUID(userUid string) (models.User, error) {
	//there was a requirement to get the total user monthly budget to display it in the front end....
	//I don't think we need to have that on the user struct exactly -> just sum up the target amounts from the
	//budgets in the user.budgets slice -> the front end can do that (or if we really want to, the back end can
	//easily do that and pass it up to the front end, but since the front end will already be wanting to have the budgets
	// to display, I think it makes the most sense to do it that way)

	user, err := u.repo.GetUserByUID(userUid)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (u UserService) CreateBudget(userUid string, targetAmount int, category string) (models.Budget, error) {
	budget := models.Budget{
		Target:   targetAmount,
		Category: category,
	}
	newBudget, err := u.repo.CreateBudget(userUid, budget)
	if err != nil {
		return models.Budget{}, err
	}

	return newBudget, nil
}

func (u UserService) UpdateBudget(budgetUid string, category string) (models.Budget, error) {
	budget, err := u.repo.UpdateBudgetCategory(budgetUid, category)
	if err != nil {
		return models.Budget{}, err
	}

	return budget, nil
}

func (u UserService) DeleteBudget(budgetUid string) error {
	err := u.repo.DeleteBudgetById(budgetUid)
	if err != nil {
		return err
	}

	return nil
}

// a few ways we can go about this
//expose this function as part of the userService interface at the application layer
//allow a cron job of some kind to run a standard 1x all the way up to n*x times a day
//depending on our requirements. We could also allow the GetUserId to run this if the delta
//from the last updated at is great enough -> with this of course now our GET user call will
//be slower and more likely to fail, so I'mn not necessarily in love with that option
func (u UserService) UpdateUserTransactions(userUid string) error {
	user, err := u.repo.GetUserByUID(userUid)
	if err != nil {
		return err
	}

	transactions, err := u.transactionClient.GetTransactionsForUser(user.Account.AccountID, user.Account.LastUpdatedAt, time.Now())
	if err != nil {
		return err
	}

	user = updateBudgetsFromTransactions(user, transactions)

	err = u.repo.UpdateUser(user)
	if err != nil {
		return err
	}

	return nil
}

func updateBudgetsFromTransactions(user models.User, txns []transactions.Transaction) models.User {
	//I was told that we could just blindly assume that the categories were 1:1 and matching from transactions
	//to the actual user budget category names. So just going to reconcile that

	// a few things, the transactions seemed to be negative for the spending transactions
	//we want to handle that here

	for _, txn := range txns {
		for _, budget := range user.Budgets {
			if txn.Category == budget.Category {
				budget.ActualSpending += abs(txn.Amount)
			}
		}
	}

	return user

}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
