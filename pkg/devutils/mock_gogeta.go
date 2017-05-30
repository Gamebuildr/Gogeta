package devutils

// GetMessage will create a fake message, clone a repo, and upload it to google cloud
func GetMessage() string {
	return `{
			"id":"58dc12e993179a0012a592dc",
			"project":"RepoSizeTest",
			"enginename":"Godot",
			"engineversion":"2.1",
			"engineplatform":"PC",
			"repotype":"Git",
			"repourl":
			"https://github.com/dirty-casuals/Calamity.git",
			"buildowner":"herman.rogers@gmail.com"
		}`
}
