package townsita

type DataAdapter interface {
	Init()
	MustGetMessageTypes() []*MessageType
	MustGetMessageSubTypes(id string) []*MessageType
	GetMessageTypeById(id string) *MessageType
	SaveMessage(message *Message, user *User) (string, error)
	GetMessageById(id string) (*Message, error)
	LoginUser(*User) (string, error)
}
