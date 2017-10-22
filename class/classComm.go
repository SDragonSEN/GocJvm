package class

type Class interface {
	ReadClass(classname string) ([]byte, error)
}
