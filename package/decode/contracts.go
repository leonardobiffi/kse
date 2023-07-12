package decode

type Service interface {
	Execute(content []byte) (output []byte, err error)
}
