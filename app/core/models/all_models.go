package models

//go:generate go run ../../../../scripts/generate_all_models.go

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*            - DB Models -            */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// All DB Models should implement this.
type Model interface {
	TableName() string
}

// —
// DO NOT EDIT this slice manually, just run go generate ./...
// and any model defined in this package should be added automatically.
var AllModels = []any{
	&GPTChat{},
	&GPTMessage{},
	&Group{},
	&User{},
	&UsersInGroup{},
}
