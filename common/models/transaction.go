package models

type Transaction struct {
	From   Account `json:"from"`
	To     Account `json:"to"`
	Value  uint    `json:"value"`
	Reason string  `json:"reason"`
}

func NewTransaction(from Account, to Account, value uint, reason string) *Transaction {
	return &Transaction{
		From:   from,
		To:     to,
		Value:  value,
		Reason: string(getReason(reason)),
	}
}

type Reason string

const (
	OTHER       Reason = ""
	SELF_REWARD        = "self-reward"
	BIRTHDAY           = "birthday"
	LOAN               = "loan"
)

func getReason(reason string) Reason {
	switch reason {
	case "self-reward":
		return SELF_REWARD
	case "birthday":
		return BIRTHDAY
	case "loan":
		return LOAN
	}
	return OTHER
}

func (s Reason) IsValid() bool {
	switch s {
	case OTHER, SELF_REWARD, BIRTHDAY, LOAN:
		return true
	}

	return false
}
