package redis

type CommandName string

const (
	Get CommandName = "Get"
)

type Command struct {
}

func NewCommand(name string, args ...[]interface{}) {

}
