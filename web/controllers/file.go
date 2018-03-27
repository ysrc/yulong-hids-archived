package controllers

import (
	"yulong-hids/web/models"
	"yulong-hids/web/settings"
	"yulong-hids/web/utils"

	"fmt"
	"path"

	"github.com/astaxie/beego"
)

// FileController /file
type FileController struct {
	BaseController
}

// Upload HTTP method POST
func (c *FileController) Upload() {

	// see DownloadController
	system := c.GetString("system")
	platform := c.GetString("platform")
	c.Data["json"] = FileUpload(system, platform, c)
	c.ServeJSON()
	return
}

// FileUpload save file and file info after upload file
func FileUpload(system string, platform string, c *FileController) *models.CodeInfo {
	if !utils.StringInSlice(system, settings.SystemArray) ||
		!utils.StringInSlice(platform, settings.PlatformArray) {
		// check all param in white list
		beego.Info("文件上传参数错误")
		return ErrorReturn()
	}

	// maybe web should add a fix to filename with random string
	// but we delete this feature, it is not suitable for this
	filename := fmt.Sprintf("%s-%s-%s", system, platform, "agent")
	filename = path.Join(settings.FilePath, filename)

	err := c.SaveToFile("file", filename)
	if err != nil {
		beego.Error("SaveToFile Error:", err)
		return ErrorReturn()
	}
	md5 := utils.GetFileMD5Hash(filename)
	if md5 == "" {
		beego.Error("GetFileMD5Hash null error.")
		return ErrorReturn()
	}
	filemodel := &models.File{
		Platform: platform,
		System:   system,
		Hash:     md5,
		Type:     "agent",
	}
	if res := filemodel.Update(); !res {
		beego.Error("Filemodel Update Error!!!")
		return ErrorReturn()
	}
	return models.NewNormalInfo(settings.Succeed)
}

// ErrorReturn return the error struct
func ErrorReturn() *models.CodeInfo {
	return models.NewErrorInfo(settings.Failure)
}
