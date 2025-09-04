package models

const (
	TypeSimpleUtterance = "SimpleUtterance"
)

// Request описывает запрос пользователя.
// см. https://yandex.ru/dev/dialogs/alice/doc/request.html
type Request struct {
	Request  SimpleUtterance `json:"request"`
	Version  string          `json:"version"`
	Timezone string          `json:"timezone"`
	Session  Session         `json:"session"`
}

type User struct {
	UserID string `json:"user_id"`
}

type Session struct {
	New  bool `json:"new"`
	User User
}

// SimpleUtterance описывает команду, полученную в запросе типа SimpleUtterance.
type SimpleUtterance struct {
	Type    string `json:"type"`
	Command string `json:"command"`
}

// Response описывает ответ сервера.
// см. https://yandex.ru/dev/dialogs/alice/doc/response.html
type Response struct {
	Response ResponsePayload `json:"response"`
	Version  string          `json:"version"`
}

// ResponsePayload описывает ответ, который нужно озвучить.
type ResponsePayload struct {
	Text string `json:"text"`
}
