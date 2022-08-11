package models

import "time"

type User struct {
	Uid     string
	Account Account
	Budgets []Budget
}

//can have this data on user itself, or on an account or something

type Account struct {
	Uid string
	//when talking on the API, I remember it took an ID, time, time as inputs. I assumed that we were high level and
	//that for the sake of the interview, I could pass in the userId. If this were real we'd need to pass in the actual
	//id and get the correct auth tied to the bank account.... We'll have accountID be that
	AccountID     string
	LastUpdatedAt time.Time
}

type Budget struct {
	Uid string
	//Like I said yesterday, can store these actual categories as enums on the db, constants in code, or as defintions in
	//their own db table -> I think enums in db is a good way to go,the front end could also have valudation on these inputs
	Category       string
	Target         int
	ActualSpending int
}
