package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ifilespace/ifilesdk-demo/data"
	"github.com/ifilespace/ifilesdk-demo/model"
	"github.com/ifilespace/ifilesdk-go"
	"github.com/ifilespace/ifilesdk-go/config"
	ifilemodel "github.com/ifilespace/ifilesdk-go/model"
	ifilesdkutil "github.com/ifilespace/ifilesdk-go/util"
)

func main() {
	data.Init()
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.GET("/", Index)
	r.POST("/getinfo", Getinfo)
	r.POST("/saveconfig", Saveconfig)
	r.POST("/createuser", Createuser)
	r.POST("/createproject", Createproject)
	r.POST("/getproject", Getproject)
	r.POST("/getiframeurl", Getiframeurl)
	r.POST("/deleteproject", DeleteProject)
	r.POST("/joinproject", JoinProject)
	r.Run(":8080")
}
func Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{"title": "iFileSDK DEMO演示"})
}
func Getinfo(c *gin.Context) {
	var config model.Config
	data.SQLDB.Get(&config, "select * from config")
	userlist := make([]model.Users, 0)
	data.SQLDB.Select(&userlist, "select * from users")
	projectlist := make([]model.Project, 0)
	if len(userlist) > 0 {
		data.SQLDB.Select(&projectlist, "select * from project ")
	}
	tasklist := make([]model.Task, 0)
	data.SQLDB.Select(&tasklist, "select * from task")
	c.JSON(http.StatusOK, gin.H{"status": 1, "config": config, "userlist": userlist, "projectlist": projectlist, "tasklist": tasklist})
}
func Getproject(c *gin.Context) {
	var post map[string]int
	err := c.BindJSON(&post)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": -1, "msg": err.Error()})
		return
	}
	if post["uid"] <= 0 {
		c.JSON(http.StatusOK, gin.H{"status": -1, "msg": "参数错误"})
		return
	}
	projectlist := make([]model.Project, 0)

	data.SQLDB.Select(&projectlist, "select * from project where uid=?", post["uid"])
	c.JSON(http.StatusOK, gin.H{"status": 1, "data": projectlist})
}
func Saveconfig(c *gin.Context) {
	var config model.Config
	err := c.BindJSON(&config)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": -1, "msg": err.Error()})
		return
	}
	if config.ID > 0 {
		_, err = data.SQLDB.Exec("update config set keyid=?,keysecret=?,ifileurl=? where id=1", config.Keyid, config.Keysecret, config.Ifileurl)
	} else {
		_, err = data.SQLDB.Exec("insert into  config (keyid,keysecret,ifileurl) values(?,?,?)", config.Keyid, config.Keysecret, config.Ifileurl)
	}
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": -1, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": 1, "msg": "保存成功"})
}
func Md5s(s string) string {
	r := md5.Sum([]byte(s))
	return hex.EncodeToString(r[:])
}
func Createuser(c *gin.Context) {
	var post model.Users
	err := c.BindJSON(&post)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": -1, "msg": err.Error()})
		return
	}
	if post.Username == "" || post.Password == "" {
		c.JSON(http.StatusOK, gin.H{"status": -1, "msg": "请输入信息"})
		return
	}
	var conf model.Config
	err = data.SQLDB.Get(&conf, "select * from config where id=1")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": -1, "msg": "请配置config"})
		return
	}
	ifile := ifilesdk.NewIfile(&config.Config{
		Keyid:     conf.Keyid,
		Keysecret: conf.Keysecret,
		Url:       conf.Ifileurl,
	})
	res, err := ifile.CreateUser(&ifilemodel.CreateUserReq{Username: post.Username, Email: post.Email, Password: post.Password, Mobile: post.Mobile})
	fmt.Println(res)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusOK, gin.H{"status": -1, "msg": "创建iFile用户失败"})
		return
	}
	if res.Status == 1 {
		_, err = data.SQLDB.Exec("insert into users (username,password,ifileuid,email,mobile) values(?,?,?,?,?)", post.Username, Md5s(post.Password), res.Data.ID, post.Email, post.Mobile)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"status": -1, "msg": "绑定ifile用户失败"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": 1, "msg": "添加成功"})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": -1, "msg": res.Msg})
		return
	}

}
func Createproject(c *gin.Context) {
	var post model.Project
	err := c.BindJSON(&post)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": -1, "msg": err.Error()})
		return
	}
	if post.Title == "" || post.UID <= 0 {
		c.JSON(http.StatusOK, gin.H{"status": -1, "msg": "请输入信息"})
		return
	}
	var ifileuid int
	err = data.SQLDB.Get(&ifileuid, "select ifileuid from users where id=?", post.UID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": -1, "msg": "获取绑定id失败"})
		return
	}
	var conf model.Config
	err = data.SQLDB.Get(&conf, "select * from config where id=1")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": -1, "msg": "请配置config"})
		return
	}
	ifile := ifilesdk.NewIfile(&config.Config{
		Keyid:     conf.Keyid,
		Keysecret: conf.Keysecret,
		Url:       conf.Ifileurl,
	})
	res, err := ifile.CreateFolder(&ifilemodel.CreateFolderReq{Name: post.Title + "-" + strconv.Itoa(ifileuid), Special: 1, PFileID: "root", Uid: ifileuid})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": -1, "msg": "新建文件夹失败"})
		return
	}
	fmt.Println(res)
	if res.Status == 1 {
		_, err = data.SQLDB.Exec("insert into project (title,ifile_root,uid) values(?,?,?)", post.Title, res.Data, post.UID)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"status": -1, "msg": "添加项目失败"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": 1, "msg": "添加成功"})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": -1, "msg": "添加项目失败1"})
		return
	}
}
func DeleteProject(c *gin.Context) {
	var post map[string]string
	err := c.BindJSON(&post)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": -1, "msg": err.Error()})
		return
	}
	var projectid = post["projectid"]
	var leixing = post["leixing"]
	if projectid == "" || leixing == "" {
		c.JSON(http.StatusOK, gin.H{"status": -1, "msg": "参数错误"})
		return
	}
	var ifileuid int
	var fileid string
	err = data.SQLDB.QueryRow("select a.ifile_root,b.ifileuid from project a left join users b on a.uid=b.id where a.id=?", projectid).Scan(&fileid, &ifileuid)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": -1, "msg": "获取信息失败"})
		return
	}
	var conf model.Config
	err = data.SQLDB.Get(&conf, "select * from config where id=1")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": -1, "msg": "请配置config"})
		return
	}
	ifile := ifilesdk.NewIfile(&config.Config{
		Keyid:     conf.Keyid,
		Keysecret: conf.Keysecret,
		Url:       conf.Ifileurl,
	})
	if leixing == "delete" {

		res, err := ifile.DeleteFolder(&ifilemodel.DeleteFolderReq{FileID: fileid, Uid: ifileuid})
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"status": -1, "msg": "删除项目失败"})
			return
		}
		if res.Status == 1 {
			data.SQLDB.Exec("delete from project where id=?", projectid)
			c.JSON(http.StatusOK, gin.H{"status": 1, "msg": "删除项目成功"})
		} else {
			c.JSON(http.StatusOK, gin.H{"status": -1, "msg": "删除文件夹失败"})
			return
		}
	} else if leixing == "cancelbind" {
		res, err := ifile.CancelBindFolder(&ifilemodel.CancelBindFolderReq{FileID: fileid, Uid: ifileuid})
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"status": -1, "msg": "删除项目失败"})
			return
		}
		if res.Status == 1 || res.Msg == "设置成功" {
			data.SQLDB.Exec("delete from project where id=?", projectid)
			c.JSON(http.StatusOK, gin.H{"status": 1, "msg": "删除项目成功"})
		} else {
			c.JSON(http.StatusOK, gin.H{"status": -1, "msg": "取消绑定失败"})
			return
		}
	} else {
		c.JSON(http.StatusOK, gin.H{"status": -1, "msg": "类型错误"})
	}
}
func Getiframeurl(c *gin.Context) {
	var post map[string]int
	err := c.BindJSON(&post)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": -1, "msg": err.Error()})
		return
	}
	if post["uid"] <= 0 {
		c.JSON(http.StatusOK, gin.H{"status": -1, "msg": "请选择用户"})
		return
	}
	var ifileuid int
	err = data.SQLDB.Get(&ifileuid, "select ifileuid from users where id=?", post["uid"])
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": -1, "msg": "获取绑定id失败"})
		return
	}
	rootpath := "root"
	if post["projectid"] > 0 {
		err = data.SQLDB.Get(&rootpath, "select ifile_root from project where id=?", post["projectid"])
		if err != nil {
			rootpath = "root"
		}
	}
	var conf model.Config
	err = data.SQLDB.Get(&conf, "select * from config where id=1")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": -1, "msg": "请配置config"})
		return
	}
	ifile := ifilesdk.NewIfile(&config.Config{
		Keyid:     conf.Keyid,
		Keysecret: conf.Keysecret,
		Url:       conf.Ifileurl,
	})
	var res ifilemodel.GetUserTokenRet
	if rootpath == "root" {
		res, err = ifile.GetUserToken(ifileuid)
	} else {
		res, err = ifile.GetProjectToken(ifileuid, rootpath)
	}

	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": -1, "msg": "获取用户Token失败"})
		return
	}
	if res.Status != 1 {
		c.JSON(http.StatusOK, gin.H{"status": -1, "msg": res.Msg})
		return
	}
	now := time.Now().Unix()
	sign := ifilesdkutil.GetSignToken(conf.Keyid, conf.Keysecret, now)
	//参数
	// token ifile用户token
	// driveid 用户空间id
	// keyid ifile 绑定的keyid
	// sign 访问ifile的签名授权
	// now 当前时间戳
	// theme 页面样式 dark暗黑模式其他为 普通模式
	// hideside 是否隐藏侧边栏 1隐藏，不传或0为显示侧边栏
	// root 可访问文件的根目录
	iurl := conf.Ifileurl + "/thirdapp/?token=" + res.Data + "&driveid=" + res.Driveid + "&keyid=" + conf.Keyid + "&sign=" + sign + "&now=" + strconv.FormatInt(now, 10) + "&theme=dark&hideside=0&root=" + rootpath
	// iurl := "http://127.0.0.1:3000/thirdapp/?token=" + res.Data + "&driveid=" + res.Driveid + "&keyid=" + conf.Keyid + "&sign=" + sign + "&now=" + strconv.FormatInt(now, 10) + "&theme=dark&hideside=0&root=" + rootpath
	c.JSON(http.StatusOK, gin.H{"status": 1, "msg": "获取成功", "data": iurl})

}
func JoinProject(c *gin.Context) {
	var post map[string]interface{}
	err := c.BindJSON(&post)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": -1, "msg": err.Error()})
		return
	}
	if post["uid"].(float64) <= 0 {
		c.JSON(http.StatusOK, gin.H{"status": -1, "msg": "请选择用户"})
		return
	}
	if post["projectid"].(float64) <= 0 {
		c.JSON(http.StatusOK, gin.H{"status": -1, "msg": "请选择项目"})
		return
	}
	var ifileuid int
	err = data.SQLDB.Get(&ifileuid, "select ifileuid from users where id=?", post["uid"])
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": -1, "msg": "获取ifileuid错误"})
		return
	}
	var fileid string
	err = data.SQLDB.Get(&fileid, "select ifile_root from project where id=?", post["projectid"])
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": -1, "msg": "获取fileid错误"})
		return
	}
	auth := post["auth"].(string)
	var conf model.Config
	err = data.SQLDB.Get(&conf, "select * from config where id=1")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": -1, "msg": "请配置config"})
		return
	}
	ifile := ifilesdk.NewIfile(&config.Config{
		Keyid:     conf.Keyid,
		Keysecret: conf.Keysecret,
		Url:       conf.Ifileurl,
	})
	res, err := ifile.JoinProject(ifileuid, fileid, auth)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": -1, "msg": "加入项目失败"})
		return
	}
	if res.Status != 1 {
		c.JSON(http.StatusOK, gin.H{"status": -1, "msg": res.Msg})
		return
	}
	data.SQLDB.Exec("update project set userid=?,auth=? where id=?", post["uid"], auth, post["projectid"])
	c.JSON(http.StatusOK, gin.H{"status": 1, "msg": "加入成功"})
}
