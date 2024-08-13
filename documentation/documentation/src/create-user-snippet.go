func createUser(username string, password string) *types.ResponseError {
	...
	cmd = exec.Command(
		"sh",
		"-c",
		fmt.Sprintf("echo %s:%s | chpasswd", username, password),
	)
	...
}
