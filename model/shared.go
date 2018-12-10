package model
type App struct {
	UserLogged *User
	Users [] *User
}
type AppResult struct {
	UserLogged User
	Admin bool
}
type AppResultEmail struct {
	UserLogged User
	Admin bool
	MsgResult string
	Email  Email
}
type AppResultDocument struct {
	UserLogged User
	Admin bool
	MsgResult string
	Docs [] Document
}
type AppResultUser struct {
	UserLogged User
	Admin bool
	MsgResult string
	Users [] User
}