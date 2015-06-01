package apis

type BrandInitArgs struct {
	BrandName string `json:"brandName"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type SignupArgs struct {
	Password string `json:"password"`
}
