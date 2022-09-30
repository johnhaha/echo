package echo

type Storage interface {
	Save(key string, val string) (err error)
	ListAppend(key string, data string) (err error)
	ListRem(key string, data string) (err error)
	Find(key string) ([]string, error)
	Rem(key string) error
	Get(key string) (string, error)
	UUID() string
	AppName() string
}
