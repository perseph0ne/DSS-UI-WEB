package model
type Result struct {
	Code int
	Message string
	Json string
}

type RequestBase struct {
	Action string
}

type RequestDeleteDocument struct {
	Base RequestBase
	ID string
}

type RequestGetDocument struct {
	Base RequestBase
	ID string
}


type RequestDownloadDocument struct {
	Base RequestBase
	ID string
}


type RequestListDocument struct {
	Base RequestBase
}


type RequestCreateDocument struct {
	Base RequestBase
	Name string
	Content []byte
}
type RequestSendMail struct {
	Base RequestBase
	From string
	To []string
	Message string
}
