package storer

type Storer struct {
}

func New() (Storer, error) {
	return Storer{}, nil
}
