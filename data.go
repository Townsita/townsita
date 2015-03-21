package townsita

type DataAdapter interface {
	Init()
	MustGetMessageTypes() []*MessageType
	MustGetMessageSubTypes(id string) []*MessageType
	GetMessageTypeById(id string) *MessageType
}
