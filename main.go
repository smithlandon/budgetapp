package main

import (
	"app/models"
	"app/mysql"
	"app/transactions"
	"app/users"
	"encoding/json"
	"log"
	"net/http"
)

func main() {
	store := mysql.Store{}
	transactionClient := transactions.NewTransactionClient("some.url")

	userService := users.NewUserService(store, transactionClient)

	handler := handler{
		userService: userService,
	}

	handleRequests(handler)
}

type handler struct {
	userService userService
}

type userService interface {
	GetUserByUID(uid string) (models.User, error)

	CreateBudget(userUid string, targetAmount int, category string) (models.Budget, error)
	UpdateBudget(budgetUid string, category string) (models.Budget, error)
	DeleteBudget(budgetUId string) error

	//Dont have a GETBudget nor PutUser, but would be easy to have those endpoints. I'd rather bake it into
	//the get user endpoint for now to be quick and manage requirements
}

//moving fast here just to show roughly routes
func handleRequests(hanlder handler) {
	http.HandleFunc("/users/{uid}", hanlder.handleUsers)
	http.HandleFunc("/users/{uid}/budgets", hanlder.handleBudgets)
	http.HandleFunc("/users/{uid}/budgets/id", hanlder.handleBudgets)

	//random port to run
	log.Fatal(http.ListenAndServe(":10000", nil))
}

//using default http client instead of mux to move fast. With mux we could have split this up
//at the router level by request method
func (h handler) handleUsers(w http.ResponseWriter, r *http.Request) {
	//for now just have GET since that's the only requirement
	switch r.Method {
	case http.MethodGet:
		h.GetUser(w, r)
	default:
		http.Error(w, "", 405)
	}
}

func (h handler) handleBudgets(w http.ResponseWriter, r *http.Request) {
	//for now just have GET since that's the only requirement
	switch r.Method {
	case http.MethodPost:
		h.CreateBudget(w, r)
	case http.MethodPut:
		h.UpdateBudget(w, r)
	case http.MethodDelete:
		h.DeleteBudget(w, r)
	default:
		http.Error(w, "", 405)
	}
}

//I like strong structs for my requests and responses
type UserResponse struct {
	User *models.User `json:"user,omitempty"`
}

func (h handler) GetUser(w http.ResponseWriter, r *http.Request) {
	//usuaully I use gorilla/mux -> very easy to get id string from address, we'll just hardcode for speed
	uid := "uid_1"

	user, err := h.userService.GetUserByUID(uid)
	if err != nil {
		Respond(http.StatusInternalServerError, err, w)
	}

	response := UserResponse{
		User: &user,
	}

	Respond(http.StatusOK, response, w)
}

type NewBudgetRequest struct {
	UserUid  string
	Category string
	Target   int
}

type NewBudgetResponse struct {
	Budget *models.Budget `json:"budget,omitempty"`
}

func (h handler) CreateBudget(w http.ResponseWriter, r *http.Request) {
	var newBudgetRequest NewBudgetRequest

	defer r.Body.Close()

	err := json.NewDecoder(r.Body).Decode(&newBudgetRequest)
	if err != nil {
		Respond(http.StatusInternalServerError, err, w)
	}

	budget, err := h.userService.CreateBudget(newBudgetRequest.UserUid, newBudgetRequest.Target, newBudgetRequest.Category)
	if err != nil {
		Respond(http.StatusInternalServerError, err, w)
	}

	response := NewBudgetResponse{
		Budget: &budget,
	}

	Respond(http.StatusOK, response, w)
}

type UpdateBudgetRequest struct {
	BudgetUid string
	Category  string
}

func (h handler) UpdateBudget(w http.ResponseWriter, r *http.Request) {
	//same deal as Create, get request and thenc all the service layer to do the update
	//skipping to stay under 45 mins

}

func (h handler) DeleteBudget(w http.ResponseWriter, r *http.Request) {
	//same program, grab the id from the request and then call the service layer to do the delete
	//skipping to stay under 45 mins

}

func Respond(statusCode int, payload interface{}, w http.ResponseWriter) {
	bytes, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}
