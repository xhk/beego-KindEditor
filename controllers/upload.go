package controllers

import (
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
)

// 判断文件夹是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

type UploadController struct {
	beego.Controller
}

func (this *UploadController) Post() {
	this.TplName = "uploadresult.html"
	//文件保存目录路径
	var savePath = "static/attached/"

	//文件保存目录URL
	var saveUrl = "/static/attached/"

	//定义允许上传的文件扩展名
	var extTable = make(map[string]string)
	extTable["image"] = "gif,jpg,jpeg,png,bmp"
	extTable["flash"] = "swf,flv"
	extTable["media"] = "swf,flv,mp3,wav,wma,wmv,mid,avi,mpg,asf,rm,rmvb"
	extTable["file"] = "doc,docx,xls,xlsx,ppt,htm,html,txt,zip,rar,gz,bz2"

	//最大文件大小
	var maxSize = 1000000

	f, imgFile, err := this.GetFile("imgFile")

	if err != nil {
		this.showError("请选择文件。")
		return
	}
	defer f.Close()
	var dirPath = savePath
	exist, err := PathExists(dirPath)
	if err != nil || !exist {
		this.showError("上传目录不存在。")
		return
	}

	var dirName = this.GetString("dir")
	if dirName == "" {
		dirName = "image"
	}

	extStr, ok := extTable[dirName]
	if !ok {
		this.showError("目录名不正确。")
		return
	}

	var fileName = imgFile.Filename
	var fileExt = strings.ToLower(path.Ext(fileName))

	if imgFile.Size <= 0 || imgFile.Size > int64(maxSize) {
		this.showError("上传文件大小超过限制。")
		return
	}

	if fileExt == "" || !strings.ContainsAny(extStr, fileExt[1:]) {
		this.showError("上传文件扩展名是不允许的扩展名。\n只允许" + extTable[dirName] + "格式。")
		return
	}

	//创建文件夹
	dirPath += dirName + "/"
	saveUrl += dirName + "/"
	exist, err = PathExists(dirPath)
	if err != nil {
		this.showError("获取文件夹信息失败")
		return
	}

	if !exist {
		err = os.Mkdir(dirPath, os.ModePerm)
		if err != nil {
			this.showError("创建文件夹失败")
			return
		}
	}

	var ymd = time.Now().Format("20060102")
	dirPath += ymd + "/"
	saveUrl += ymd + "/"

	exist, err = PathExists(dirPath)
	if err != nil {
		this.showError("获取文件夹信息失败")
		return
	}

	if !exist {
		err = os.Mkdir(dirPath, os.ModePerm)
		if err != nil {
			this.showError("创建文件夹失败")
			return
		}
	}

	var newFileName = time.Now().Format("20060102150405") + "_" + strconv.Itoa(int(time.Now().Unix()%10000)) + fileExt
	var filePath = dirPath + newFileName

	err = this.SaveToFile("imgFile", filePath)
	if err != nil {
		this.showError("保存文件失败")
		return
	}

	var fileUrl = saveUrl + newFileName

	this.Data["error"] = 0
	this.Data["url"] = fileUrl
}

func (this *UploadController) showError(errMsg string) {
	this.Data["error"] = 1
	this.Data["url"] = errMsg

	// 这里不调用StopRun来终止执行逻辑，因为不清楚其机制，
	// 万一造成了文件资源没释放呢，所以采用有错误就return的方式
}

//---------------File Manager Controller------------------
type UploadFileMgrController struct {
	beego.Controller
}

type FileInfo struct {
	Datetime string `json:"datetime"`
	Dir_path string `json:"dir_path"`
	Filename string `json:"filename"`
	Filetype string `json:"filetype"`
	Filesize int64  `json:"filesize"`
	Is_dir   bool   `json:"is_dir"`
	Has_file bool   `json:"has_file"`
	Is_photo bool   `json:"is_photo"`
}

type FileList struct {
	Current_dir_path string     `json:"current_dir_path"`
	Current_url      string     `json:"current_url"`
	Moveup_dir_path  string     `json:"moveup_dir_path"`
	Total_count      int        `json:"total_count"`
	File_list        []FileInfo `json:"file_list"`
}

// 文件列表排序接口
type BySize []FileInfo

func (a BySize) Len() int           { return len(a) }
func (a BySize) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a BySize) Less(i, j int) bool { return a[i].Filesize < a[j].Filesize }

type ByType []FileInfo

func (a ByType) Len() int           { return len(a) }
func (a ByType) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByType) Less(i, j int) bool { return a[i].Filetype < a[j].Filetype }

type ByName []FileInfo

func (a ByName) Len() int           { return len(a) }
func (a ByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByName) Less(i, j int) bool { return a[i].Filename < a[j].Filename }

// 一次遍历，把目录和文件都获取到
func getDirFileList(dir string) ([]FileInfo, []FileInfo) {
	//图片扩展名
	var fileTypes = "gif,jpg,jpeg,png,bmp"
	dirList := make([]FileInfo, 0, 100)
	fileList := make([]FileInfo, 0, 100)
	files, _ := ioutil.ReadDir(dir)
	for _, f := range files {
		//fmt.Println(f.Name())
		var fi = FileInfo{f.ModTime().Format("2006-01-02 15:04:05"), "", f.Name(), "", f.Size(), false, false, true}
		if f.IsDir() {
			fi.Is_dir = true
			fi.Is_photo = false
			fi.Has_file = true
			dirList = append(dirList, fi)
		} else {
			fi.Filetype = strings.ToLower(path.Ext(fi.Filename))[1:]
			fileList = append(fileList, fi)
			if strings.ContainsAny(fileTypes, fi.Filetype) {
				fi.Is_photo = true
			} else {
				fi.Is_photo = false
			}
		}
	}

	return dirList, fileList
}

// 处理文件管理请求
func (this *UploadFileMgrController) Get() {
	//根目录路径，相对路径
	var rootPath = "static/attached/"
	//根目录URL，可以指定绝对路径，比如 http://www.yoursite.com/attached/
	var rootUrl = "/static/attached/"
	//图片扩展名
	//var fileTypes = "gif,jpg,jpeg,png,bmp"
	var currentUrl = ""
	var currentDirPath = ""
	var currentPath = ""
	var moveupDirPath = "" // 上级目录

	var dirPath = rootPath
	var dirName = this.GetString("dir")
	if !(dirName == "") {
		if !strings.ContainsAny("image,flash,media,file", dirName) {
			this.showError("Invalid Directory name.")
			return
		}
		dirPath += dirName + "/"
		rootUrl += dirName + "/"
		exist, err := PathExists(dirPath)
		if err != nil {
			this.showError("获取文件夹信息失败")
			return
		}

		if !exist {
			err = os.Mkdir(dirPath, os.ModePerm)
			if err != nil {
				this.showError("创建文件夹失败")
				return
			}
		}
	}

	//根据path参数，设置各路径和URL
	var path = this.GetString("path")
	if path == "" {
		currentPath = dirPath
		currentUrl = rootUrl
		currentDirPath = ""
		moveupDirPath = ""
	} else {
		currentPath = dirPath + path
		currentUrl = rootUrl + path
		currentDirPath = path
		var reg = regexp.MustCompile("(.*?)[^\\/]+\\/$")
		moveupDirPath = reg.ReplaceAllString(currentDirPath, "$1")
	}

	//排序形式，name or size or type
	var order = this.GetString("order")
	order = strings.ToLower(order)

	//不允许使用..移动到上一级目录
	if strings.ContainsAny(path, "..") {
		this.showError("Access is not allowed.")
		return
	}
	//最后一个字符不是/
	if path != "" && path[len(path)-1:] != "/" {
		this.showError("Parameter is not valid.")
		return
	}
	//目录不存在或不是目录
	exist, err := PathExists(dirPath)
	if err != nil {
		this.showError("获取文件夹信息失败")
		return
	}

	if !exist {
		this.showError("Directory does not exist.")
		return
	}

	//遍历目录取得文件信息
	dirList, fileList := getDirFileList(currentPath)

	switch order {
	case "size":
		sort.Sort(BySize(dirList))
		sort.Sort(BySize(fileList))
		break
	case "type":
		sort.Sort(ByType(dirList))
		sort.Sort(ByType(fileList))
		break
	case "name":
	default:
		sort.Sort(ByName(dirList))
		sort.Sort(ByName(fileList))
		break
	}

	var fl = &FileList{
		currentDirPath,
		currentUrl,
		moveupDirPath,
		len(dirList) + len(fileList),
		append(dirList, fileList...),
	}

	this.Data["json"] = fl
	this.ServeJSON()
}

func (this *UploadFileMgrController) showError(errMsg string) {
	this.Data["json"] = errMsg
	this.ServeJSON()
}
