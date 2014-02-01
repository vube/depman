package result

var err bool

func Error() {
	err = true
}

func ExitWithError() bool {
	return err
}
