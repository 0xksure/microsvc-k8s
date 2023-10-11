package shared

type Args struct{}
type Reply string
type MessageServer string

func (t *MessageServer) GetMessage(args *Args, reply *Reply) error {
	*reply = "hello from your server"
	return nil
}
