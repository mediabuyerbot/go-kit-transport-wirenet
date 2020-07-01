package test

type UpdateBalanceRequest struct {
	A int
	B int
}

type UpdateBalanceResponse struct {
	Sum int
	Err error
}
