package errors

func JudgeError(err error) bool {
	if err != nil {
		panic(CustomError(err.Error()))
	}
	return false
}

func ThrowError(msg string) {
	panic(CustomError(msg))
}
