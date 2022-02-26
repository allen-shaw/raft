package log

func CHECK(condition bool, statement string) {
	if !condition {
		Fatal("Check failed: %s", statement)
	}
}
