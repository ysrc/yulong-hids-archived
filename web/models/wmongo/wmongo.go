package wmongo

import (
	"github.com/astaxie/beego"
	"gopkg.in/mgo.v2"
)

var session *mgo.Session

// Conn return mongodb session.
func Conn() *mgo.Session {
	return session.Copy()
}

/*
func Close() {
	session.Close()
}
*/

func init() {
	url := beego.AppConfig.String("mongodb::url")

	sess, err := mgo.Dial(url)
	if err != nil {
		beego.Error("mongodb url:", url)
		beego.Error("mongodb session connect", err)
	}

	session = sess
	session.SetMode(mgo.Monotonic, true)
}
