package models

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*         - Database Models -         */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// -> These models translate to tables in the database

// Used to migrate all models at once
var AllDBModels = []any{
	User{},
	Group{},
	UsersInGroup{},
}
