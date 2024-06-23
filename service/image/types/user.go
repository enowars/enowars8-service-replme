package types

type UserPasswdData struct {
	Username string
	Password string
	Uid      uint32
	Gid      uint32
	Gecos    string
	Home     string
	Shell    string
}
