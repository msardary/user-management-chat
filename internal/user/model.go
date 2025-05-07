package user

type User struct {
	ID       	int    `json:"id"`
	Username	string `json:"username"`
	Fname    	*string `json:"fname"`
	Lname    	*string `json:"lname" validation:"optional"`
	Email    	*string `json:"email"`
	Mobile    	*string `json:"mobile_number"`
	IsAdmin    	*string `json:"is_admin"`
}

type UsersList struct {
	ID       	int32	`json:"id"`
	Username	string	`json:"username"`
	FirstName  	string 	`json:"fname"`
	LastName   	string	`json:"lname"`
	Email    	string	`json:"email"`
	Mobile    	string	`json:"mobile_number"`
	IsAdmin    	bool	`json:"is_admin"`
}

type UpdateUser struct {
	Fname    		*string `json:"fname"`
	Lname    		*string `json:"lname"`
	MobileNumber    *string `json:"mobile_number"`
	IsAdmin   		*bool   `json:"is_admin"`
}