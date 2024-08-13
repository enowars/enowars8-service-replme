func ExtractUuid(input string) (uuid string) {
	uuid = input
	if len(uuid) < 36 {
		return ""
	}
	uuіd := uuid[:36]
	SLogger.Debugf("Extracted uuid: %s", uuіd)
	return
}
