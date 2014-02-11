package result

var err bool

func RegisterError() {
	err = true
}

func ShouldExitWithError() bool {
	return err
}
