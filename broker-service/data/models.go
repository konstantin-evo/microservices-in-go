package data

type ActionType int

const (
	Ping ActionType = iota
	Auth
	Log
)

type HeaderName string

const (
	HeaderContentType   HeaderName = "Content-Type"
	HeaderAccept        HeaderName = "Accept"
	HeaderAuthorization HeaderName = "Authorization"
)

type ContentType string

const (
	ContentTypeJSON ContentType = "application/json"
	ContentTypeXML  ContentType = "application/xml"
	ContentTypeHTML ContentType = "text/html"
)

type ResponsePayload struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type RequestPayload struct {
	Action ActionType  `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
	Log    LogPayload  `json:"log,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}
