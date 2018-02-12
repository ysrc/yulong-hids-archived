package utils

import (
	"fmt"
	"os/exec"
	"testing"

	"github.com/astaxie/beego"
)

func TestRandStringBytesMaskImprSrc(t *testing.T) {
	t.Logf(RandStringBytesMaskImprSrc(16))
}

func TestFunction(t *testing.T) {
	cmdLst := []string{"req", "-new", "-x509", "-key", "upload_files/private.pem", "-out", "/tmp/cert1.pem", "-days", "3650", "-subj", "/CN=domain-sec-project.com"}
	ex := exec.Command("openssl", cmdLst...)
	beego.Debug(cmdLst)
	out, err := ex.Output()
	if err != nil {
		beego.Error(err)
	}
	fmt.Printf("The date is %s\n", out)
}
