package cmd

type exitMemento struct {
	code int
}

func (e *exitMemento) Exit(i int) {
	e.code = i
}
