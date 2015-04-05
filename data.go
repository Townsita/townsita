package townsita

type DataAdapter interface {
	Init()
	MustGetMessageTypes() []*MessageType
	MustGetMessageSubTypes(id string) []*MessageType
	GetMessageTypeById(id string) *MessageType
	SaveMessage(message *Message, user *User) (string, error)
	GetMessageById(id string) (*Message, error)
	LoginUser(user *User) (string, error)
	RegisterUser(user *User) (string, error)
	LoadUserByID(userID string) (*User, error)
	GetOwnMessages(userID string, limit, offset int) []*Message
	GetReceivedMessages(userID string, limit, offset int) []*Message
}
