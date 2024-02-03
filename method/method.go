package method

type Method int

const (
	GET Method = iota
	POST
	PUT
	DELETE
)

func (m Method) String() string {
	return [...]string{"GET", "POST", "PUT", "DELETE"}[m]
}
