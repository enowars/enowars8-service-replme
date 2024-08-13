func (devenv *DevenvController) GetFileContent(ctx *gin.Context) {
	_uuid, _ := ctx.Get("uuid")
	uuid := _uuid.(string)
	name := ctx.Param("name")
	path := filepath.Join(devenv.DevenvFilesPath, uuid, name)

	if !strings.HasPrefix(path, devenv.DevenvFilesPath) {
		...
		return
	}

	content, err := util.GetFileContent(path)
	...
}
