package interior_models

type SadException struct {
	Err error
}

func (e *SadException) Message() string {
	return e.Err.Error()
}

func (e *SadException) Error() string {
	return e.Err.Error()
}

type SuprisingException struct {
	Err error
}

func (e *SuprisingException) Message() string {
	return e.Err.Error()
}

func (e *SuprisingException) Error() string {
	return e.Err.Error()
}
