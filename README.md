# beego-KindEditor
KindEditor use under beego framework

The Beego WebSite https://beego.me/
The KindEditor WebSite http://kindeditor.net

These is no version for beego, So I create this.

How to use Beego-KindEditor

1. Copy static/kindeditor to your beego project static directory and Create 'attached' directory under this directory

2. Copy controllers/upload.go to your beego project controllers directory 

3. Insert follow code into routers/router.go
    beego.Router("/upload", &controllers.UploadController{})
	beego.Router("/uploadfilemgr", &controllers.UploadFileMgrController{})

4. Add the script and html like the code in views/index.html

'''html
<link rel="stylesheet" href="/static/Kindeditor/themes/default/default.css" />
  <script charset="utf-8" src="/static/kindeditor/kindeditor-all.js"></script>
  <script charset="utf-8" src="/static/Kindeditor/lang/zh-CN.js"></script>
  <script>
    var editor;
    KindEditor.ready(function(K) {
      editor = K.create('textarea[name="content"]', {
        allowFileManager : true,
                  /*items : [
          'fontname', 'fontsize', '|', 'forecolor', 'hilitecolor', 'bold', 'italic', 'underline',
          'removeformat', '|', 'justifyleft', 'justifycenter', 'justifyright', 'insertorderedlist',
          'insertunorderedlist', '|', 'emoticons', 'image', 'link'],*/
          height:500
      });
    });
  </script>

   <form>
      <textarea name="content" style="width:800px;height:400px;visibility:hidden;">KindEditor</textarea>
    </form>
'''



