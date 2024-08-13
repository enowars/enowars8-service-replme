devenv := devenvs.Group("/:uuid", func(ctx *gin.Context) {
	...
	id := ctx.Query("uuid")
	if id == "" {
		id = ctx.Param("uuid")
	}
	...
	uuid := util.ExtractUuid(id)
	...
	var devenvs []model.Devenv
	err := database.DB.Model(&user).Where("id = ?", uuid[:36]).Association("Devenvs").Find(&devenvs)
	...
	ctx.Set("uuid", uuid)
	ctx.Set("current_devenv", devenvs[0])
	...
}

