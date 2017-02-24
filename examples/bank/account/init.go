package account

import "github.com/yehohanan7/cqrs"

var store cqrs.EventStore

func InitAccounts(s cqrs.EventStore) {
	store = s
}
