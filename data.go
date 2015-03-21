package townsita

type DataAdapter interface {
	Init()
	MustGetMessageTypes() []MessageType
	MustGetMessageSubTypes(messageType int) ([]MessageType, error)
}
