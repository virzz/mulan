package web

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/virzz/mulan/code"
	"github.com/virzz/mulan/rsp"
	"go.uber.org/zap"
)

func handleSystemUpgrade(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.AbortWithError(400, err)
		return
	}
	exePath, err := os.Executable()
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	newFile := exePath + "_upgrade"
	err = c.SaveUploadedFile(file, newFile)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	err = os.Chmod(newFile, 0o755)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	err = os.Rename(newFile, exePath)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	buf, err := exec.Command(exePath, "restart").CombinedOutput()
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	c.String(200, string(buf))
}

func handleSystemUpload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		zap.L().Error("Failed to get file", zap.Error(err))
		c.AbortWithStatusJSON(400, rsp.E(code.ParamInvalid, err))
		return
	}
	exePath, err := os.Executable()
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	dstPath := filepath.Join(filepath.Dir(exePath), "_"+filepath.Base(file.Filename))
	err = c.SaveUploadedFile(file, dstPath)
	if err != nil {
		zap.L().Error("Failed to save file", zap.String("file", file.Filename), zap.String("path", dstPath), zap.Error(err))
		c.AbortWithStatusJSON(400, rsp.E(code.UnknownErr, err))
		return
	}
	c.JSON(200, rsp.S(dstPath))
}
