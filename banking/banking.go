package banking

import (
	"errors"
	"fmt"
)

type Account struct {
	owner string
	balance int
}

var errNoDeposit = errors.New("Not enough deposit")

func NewAccount(owner string) *Account{
	account := Account{owner: owner, balance: 0}
	return &account
}

func (a *Account) Deposit(amount int) {
	a.balance += amount
}

func (a *Account) Withdraw(amount int) error {
	switch {
	case a.balance < amount:
		return errNoDeposit
		
	default:
		a.balance -= amount
		return nil
		
	}
}

func (a Account) String() string {
	return fmt.Sprint(a.owner, "'s account\nHas: ", a.balance)
}