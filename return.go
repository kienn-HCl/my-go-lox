package myloxgo

type ReturnValue struct {
	value any
}

func NewReturnValue(value any) *ReturnValue {
	return &ReturnValue{
		value: value,
	}
}
