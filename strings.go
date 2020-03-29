package errors

//String is a string that implements the error and stringer interfaces
type String string

func (s String) String() string {
	return string(s)
}

func (s String) Error() string {
	return string(s)
}
