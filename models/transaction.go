package models

type Transaction struct {
	From   Account `json:"from"`
	To     Account `json:"to"`
	Value  uint    `json:"value"`
	reason string  `json:"reason"`
}

func NewTransaction(from Account, to Account, value uint, data string) *Transaction {
	return &Transaction{
		From:   from,
		To:     to,
		Value:  value,
		reason: data,
	}
}

var (
	UNKNOWN     = ""
	SELF_REWARD = "self-reward"
	BIRTHDAY    = "birthday"
	LOAN        = "loan"
)

func (t Transaction) getReason() string {
	switch t.reason {
	case "self-reward":
		return SELF_REWARD
	case "birthday":
		return BIRTHDAY
	case "loan":
		return LOAN
	}
	return UNKNOWN
}
