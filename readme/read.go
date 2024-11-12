package read

import "github.com/gin-gonic/gin"

func ReadMe(c *gin.Context) {
	c.File("./html/index.html")
}

func Documentation(c *gin.Context) {
	c.FileAttachment("./README.md", "JUGGERNAUT_Documentation.md")
}
