package handlers

func isAdminPasswordValid(provided, adminPassword string) bool {
	return provided == adminPassword
}
