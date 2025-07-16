package server

import "github.com/gin-gonic/gin"

func (s *Server) handleRoot(c *gin.Context) {
	c.String(200, `
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣀⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣀⣠⣤⣴⠶⠾⠛⠋⠛⠷⢦⣤⣀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⢀⣀⣤⣤⣶⣶⣿⣿⣿⣿⣿⣿⣷⣦⣤⣀⢀⣀⣠⣬⡿⠿⠖⠀⠀⠀
⠀⢠⣶⣤⣄⣉⠛⠻⢿⣿⣿⣿⣿⣿⣿⣿⠿⠿⠛⠛⢛⣉⣉⣤⠤⢶⣶⣿⠀⠀
⠀⠀⠛⠿⢿⣿⣿⣶⣤⣈⡙⠛⠋⣉⣀⣤⣴⣶⣾⣿⣿⡇⢠⣤⣤⠀⠛⣉⠀⠀
⠀⠰⣿⣶⣤⣈⠙⠻⢿⣿⣿⣿⣿⣿⣿⣿⠿⠟⠛⠋⣉⡀⢸⣿⣿⠀⣿⣿⠀⠀
⠀⠀⠙⠻⢿⣿⣿⣶⣤⣄⡉⠛⢉⣉⣠⣤⣴⣶⣿⣿⣿⡇⢸⣿⣿⠀⠋⣉⠀⠀
⠀⠰⣿⣶⣦⣤⣉⠙⠻⢿⣿⣿⣿⣿⠿⠿⠛⠛⢉⣉⣩⡄⠘⢉⡉⠀⣿⣿⠀⠀
⠀⠀⠉⠛⠻⢿⣿⣿⣶⣤⣌⣉⣁⣤⣤⣶⣾⣿⣿⣿⣿⠇⠴⠛⠛⠀⠉⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠈⠙⠻⠿⣿⣿⣿⣿⠿⠿⠛⠋⠉⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠉⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
This is Nook, my AT Protocol Personal Data Server (ATProto PDS).

Code: https://github.com/hbjydev/nook
Bluesky: https://bsky.app/profile/hayden.moe
`)
}

func (s *Server) handleRobotsTxt(c *gin.Context) {
	c.String(200, `# Hello there! Feel free to scrape my PDS at your leisure.
# Crawling the public API is allowed
User-agent: *
Allow: /
`)
}
