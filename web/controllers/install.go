package controllers

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"
	"yulong-hids/web/models"
	"yulong-hids/web/models/wmongo"
	"yulong-hids/web/settings"
	"yulong-hids/web/utils"

	"github.com/astaxie/beego"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// InstallController /install
type InstallController struct {
	BaseController
}

// Get return the install page
func (c *InstallController) Get() {

	currentStep := HasInstall()
	if currentStep > settings.InstallStep {
		c.Ctx.Redirect(302, beego.URLFor("MainController.Get"))
		return
	}

	step, err := c.GetInt("step")
	if err != nil {
		step = 1
	}

	// API for key file download
	dlfile := c.GetString("download")
	if dlfile != "" {
		if dlfile == "public" {
			c.Ctx.Output.Download(path.Join(settings.FilePath, settings.PublicKeyName))
			return
		}
		if dlfile == "private" {
			c.Ctx.Output.Download(path.Join(settings.FilePath, settings.PrivateKeyName))
			return
		}
	}

	// csrf token
	token := utils.RandStringBytesMaskImprSrc(64)
	c.Data["token"] = token
	c.Ctx.SetCookie("request_token", token, 1073741824, "/")

	c.Data["currentstep"] = currentStep
	c.Data["step"] = step
	c.TplName = "install.tpl"

}

// Post HTTP method POST
func (c *InstallController) Post() {

	currentStep := HasInstall()
	if currentStep > settings.InstallStep {
		c.Data["json"] = bson.M{"status": false, "msg": "already Install"}
		c.ServeJSON()
		return
	}

	step, err := c.GetInt("step")
	if err != nil {
		step = 1
	}

	mConn := wmongo.Conn()
	defer mConn.Close()

	// you will nevet allowed to dump or retry any step
	if step == 1 && currentStep == 1 {
		// init mongodb
		c.Data["json"] = initDB()
		c.ServeJSON()
		return
	}
	if step == 2 && currentStep == 2 {
		// init rules
		ruleModel := models.NewRule()
		var rulelist []interface{}

		json.Unmarshal(c.Ctx.Input.RequestBody, &rulelist)
		if err := ruleModel.InsertMany(rulelist); err != nil {
			beego.Error("Rule InsertMany error:", err)
			c.Data["json"] = bson.M{"status": false, "msg": "初始化规则失败"}
		} else {
			c.Data["json"] = bson.M{"status": true, "msg": "初始化规则成功"}
		}
		c.ServeJSON()
		return
	}

	if step == 3 && currentStep > 2 {
		// upload data, agent and daemon follow install doc
		system := c.GetString("system")
		platform := c.GetString("platform")
		c.Data["json"] = c.saveFile(system, platform)
		c.ServeJSON()
		return
	}

	if step == 4 && currentStep == 4 {
		// init config collection and create new key file
		action := c.GetString("action")

		if action == "createkey" {
			err := utils.GenRsaKey(1024)
			if err != nil {
				beego.Error("Create key(utils.GenRsaKey) Error", err)
				c.Data["json"] = bson.M{"status": false, "msg": "Create key Error"}
			} else {
				// run command `openssl req -new -x509 -key private.pem -out cert.pem -days 3650 -subj "/CN=domain-sec-project.com"`
				prikey := path.Join(settings.FilePath, settings.PrivateKeyName)
				certkey := path.Join(settings.FilePath, settings.CertKeyName)
				cmdLst := []string{
					"req", "-new", "-x509", "-key", prikey,
					"-out", certkey, "-days", "3650", "-subj",
					"/CN=domain-sec-project.com",
				}
				_, err := exec.Command("openssl", cmdLst...).Output()
				if err != nil {
					beego.Error("exec Command(exec.Command)", err)
					c.Data["json"] = bson.M{"status": true, "cert": "", "msg": err.Error()}
				} else {
					cert, _ := ioutil.ReadFile(certkey)
					c.Data["json"] = bson.M{"status": true, "cert": string(cert), "msg": "yes"}
				}
			}
			c.ServeJSON()
			return
		}
		if action == "addconfig" {
			var configJOSN bson.M
			json.Unmarshal(c.Ctx.Input.RequestBody, &configJOSN)
			iplst := configJOSN["ip"].([]interface{})
			prolist := configJOSN["process"].([]interface{})
			cert := configJOSN["cert"].(string)
			toConfig(iplst, prolist, cert)
			if !updateHTTPSKeys(cert) {
				c.Data["json"] = bson.M{"status": false, "msg": "更新https文件时可能发生错误"}
				c.ServeJSON()
				return
			}
			c.Data["json"] = bson.M{"status": true, "msg": "HTTPS证书及基础配置已经更新，请重启web服务"}
			c.ServeJSON()
			return
		}
	}

	beego.Error("Real step:", step)
	beego.Error("Current step params:", currentStep)
	c.Data["json"] = bson.M{"status": false, "msg": "error"}
	c.ServeJSON()
	return
}

func updateHTTPSKeys(cert string) bool {
	newPrivate := path.Join(settings.FilePath, settings.PrivateKeyName)
	if _, err := os.Stat(newPrivate); os.IsNotExist(err) {
		beego.Error("OS Stat, file not exist error", err)
		return false
	}
	if err := os.Rename(newPrivate, beego.AppConfig.String("HTTPSKeyFile")); err != nil {
		beego.Error("OS Rename", err)
		return false
	}
	if err := ioutil.WriteFile(beego.AppConfig.String("HTTPSCertFile"), []byte(cert), 0644); err != nil {
		beego.Error("Ioutil WriteFile", err)
		return false
	}
	config := models.NewConfig()
	serverID := config.FindOne(bson.M{"type": "server"}).Id.Hex()
	private, _ := ioutil.ReadFile(beego.AppConfig.String("HTTPSKeyFile"))
	config.EditByID(serverID, "privatekey", string(private))
	public, _ := ioutil.ReadFile(path.Join(settings.FilePath, settings.PublicKeyName))
	config.EditByID(serverID, "publickey", string(public))
	return true
}

func toConfig(iplst []interface{}, prolist []interface{}, certText interface{}) {
	config := models.NewConfig()
	white := config.FindOne(bson.M{"type": "whitelist"})
	idstr := white.Id.Hex()
	for _, ipaddr := range iplst {
		if strings.Trim(ipaddr.(string), " ") != "" {
			config.AddOne(idstr, "ip", ipaddr.(string))
		}
	}
	for _, program := range prolist {
		if strings.Trim(program.(string), " ") != "" {
			config.AddOne(idstr, "process", program.(string))
		}
	}
	serverID := config.FindOne(bson.M{"type": "server"}).Id.Hex()
	config.EditByID(serverID, "cert", certText.(string))
}

// saveFile much the same with FileContorller.FileUpload
func (c *InstallController) saveFile(system string, platform string) *models.CodeInfo {
	if !utils.StringInSlice(system, settings.SystemArray) ||
		!utils.StringInSlice(platform, settings.PlatformArray) {
		beego.Error("文件上传参数错误")
		return ErrorReturn()
	}

	filename := fmt.Sprintf("%s-%s", system, platform)
	filepath := path.Join(settings.FilePath, filename)
	file, _, err := c.GetFile("file")

	if err == nil && file != nil {
		err := c.SaveToFile("file", filepath)
		if err != nil {
			beego.Error("SaveToFile Error:", err)
			return ErrorReturn()
		}
		err = deCompressZip(filepath, settings.FilePath, filename, system, platform)
		if err != nil {
			beego.Error("decompress zip file error:", err)
			return ErrorReturn()
		}

		return models.NewNormalInfo(settings.Succeed)
	}
	beego.Error("Beego GetFile error: ", err)
	return ErrorReturn()
}

// HasInstall check current step with install
func HasInstall() int {
	config := models.NewConfig()
	ser := config.FindOne(bson.M{"type": "server"})

	if ser.Dic == nil {
		return 1
	}

	publickey, is := ser.Dic.(bson.M)["publickey"]
	rule := models.NewRule()
	filemgo := models.NewFile()
	filecount := len(filemgo.FindAll(nil))
	firstRule := rule.FindOne(nil)

	if is && firstRule == nil {
		return 2
	}

	if firstRule != nil && filecount == 0 {
		return 3
	}

	if filecount > 0 && (publickey == "" || publickey == nil) {
		return 4
	}

	if publickey != "" {
		return 42 // why 42?
	}

	return 0
}

func initDB() bson.M {
	mConn := wmongo.Conn()
	defer mConn.Close()
	db := mConn.DB("")

	index := mgo.Index{
		Key:    []string{"ip", "info", "type", "status", "uptime"},
		Unique: true,
	}
	err := db.C("notice").EnsureIndex(index)
	if err != nil {
		beego.Error("Collection EnsureIndex", err)
		return bson.M{"status": false, "msg": "create index error"}
	}
	var defualtConfig []interface{}
	json.Unmarshal(settings.DefualtConfig, &defualtConfig)
	err = db.C("config").Insert(defualtConfig...)

	if err != nil {
		beego.Error("Rule Insert", err)
		return bson.M{"status": false, "msg": "init config error"}
	}
	return bson.M{"status": true, "msg": "init mongodb doned"}
}

// unzip zipfile
func deCompressZip(zipFile string, unpath string, fileprefix string, system string, platform string) error {
	reader, err := zip.OpenReader(zipFile)
	var filelst []string
	if err != nil {
		return err
	}
	defer reader.Close()

	// filecol := models.NewFile()
	toDB := func(platform string, system string, rtype string, filepath string) error {
		md5 := utils.GetFileMD5Hash(filepath)
		if md5 == "" {
			return errors.New("md5 hash string is null")
		}
		filecol := models.File{
			Platform: platform,
			System:   system,
			Type:     rtype,
			Hash:     md5,
			Uptime:   time.Now(),
		}
		if res := filecol.Update(); !res {
			return errors.New("file model insert error")
		}
		return nil
	}

	deCompress2dest := func(f *zip.File) error {

		rc, err := f.Open()
		if err != nil {
			return nil
		}
		defer rc.Close()

		filelst = append(filelst, f.Name)
		rtype := settings.FileName2Type[f.Name].(string)

		filename := path.Join(unpath, fmt.Sprintf("%s-%s-%s", system, platform, rtype))
		// only allow agent, daemon, data
		if !utils.StringInSlice(f.Name, settings.TypeArray) {
			return errors.New("type not allow")
		}
		w, err := os.Create(filename)
		if err != nil {
			return err
		}
		defer w.Close()
		_, err = io.Copy(w, rc)
		if err != nil {
			return err
		}
		err = toDB(platform, system, rtype, filename)
		return err
	}

	for _, file := range reader.File {
		err := deCompress2dest(file)
		if err != nil {
			return err
		}
	}
	return nil
}
