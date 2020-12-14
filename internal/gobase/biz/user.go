package biz

type User struct {
	ID       uint   `json:"id"`
	Account  string `json:"account"`
	Password string `json:"password"`
}

func (b *SampleBiz) CreateUser(user *User) (*User, error) {
	// b.daoCtx.CreateUser()
	return nil, nil
}
