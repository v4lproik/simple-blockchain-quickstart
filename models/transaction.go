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
		Reason: getReason(reason),
	}
}

var (
	UNKNOWN     = ""
	SELF_REWARD = "self-reward"
	BIRTHDAY    = "birthday"
	LOAN        = "loan"
)

func getReason(reason string) string {
	switch reason {
	case "self-reward":
		return SELF_REWARD
	case "birthday":
		return BIRTHDAY
	case "loan":
		return LOAN
	}
	return UNKNOWN
}
