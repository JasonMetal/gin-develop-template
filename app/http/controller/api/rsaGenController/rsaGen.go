/*
*

	@author:
	@since:
	@desc: //生成rsa 公私钥对

*
*/
package rsaGenController

import (
	"archive/zip"
	"crypto/md5"
	baseController "develop-template/app/http/controller"
	"develop-template/app/logic/rsaGenLogic"
	"fmt"
	"github.com/JasonMetal/submodule-support-go.git/bootstrap"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type controller struct {
	baseController.BaseController
}

// NewController
//
//	@Description:
//	@since: time
//	@param ctx
//	@return *controller
func NewController(ctx *gin.Context) *controller {
	return &controller{baseController.NewBaseController(ctx)}
}

// ZipRSAKey
//
//	@Description: /Rsa/gen
//	@receiver c
//	@since: time
func (c *controller) ZipRSAKey() {
	defer func() {
		r := recover()
		fmt.Println("============panic ZipRSAKey============", r)
	}()
	logic := rsaGenLogic.NewLogic(c.GCtx)
	prvKey, pubKey := logic.GenRsaKey()

	//
	filesDir := fmt.Sprintf("%sconfig/%s", bootstrap.ProjectPath(), "certs")
	err := os.MkdirAll(filesDir, 0644)
	if err != nil {
		fmt.Println("创建配置文件夹", err)
	}
	savePath := fmt.Sprintf("%sconfig/%s/%s.pem", bootstrap.ProjectPath(), "certs", "privateKey")
	savePathPublicKey := fmt.Sprintf("%sconfig/%s/%s.pem", bootstrap.ProjectPath(), "certs", "publicKey")

	filesZipPath := fmt.Sprintf("%sconfig/%s/%s.zip", bootstrap.ProjectPath(), "certs", "certs")

	// 打开文件
	filePrvKey, err := os.OpenFile(savePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		panic(err)
	}
	defer filePrvKey.Close()

	// 清空内容
	if err := emptyFile(filePrvKey); err != nil {
		fmt.Println("清空内容", err)
	}

	fmt.Println("文件filePrvKey清空成功!")

	// 打开文件
	filePubKey, err := os.OpenFile(savePathPublicKey, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		panic(err)
	}
	defer filePubKey.Close()

	// 清空内容
	if err := emptyFile(filePubKey); err != nil {
		panic(err)
	}
	fmt.Println("文件filePubKey清空成功!")

	c.writeContent(savePath, prvKey)
	c.writeContent(savePathPublicKey, pubKey)

	// 定义文件列表和zip文件路径
	files := []string{savePath, savePathPublicKey}
	dist := filesZipPath

	// 创建zip文件
	zipFile, err := os.Create(dist)
	if err != nil {
		panic(err)
	}
	defer zipFile.Close()

	// 创建zip归档
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// 将文件添加到zip
	for _, file := range files {
		c.addFileToZip(zipWriter, file)
	}
	//newFilesZipPath := fmt.Sprintf("%sconfig/%s/%s.zip", bootstrap.ProjectPath(), "certs", "certs_"+getFileMd5(filesZipPath))
	//newFilesZipPath := fmt.Sprintf("%sconfig/%s/%s.zip", bootstrap.ProjectPath(), "certs", "certs")
	//c.renameFileZip(filesZipPath, newFilesZipPath)

	fmt.Println("打包完成!")
	c.GCtx.JSON(0, "打包完成!!!")
}

func (c *controller) addFileToZip(zipWriter *zip.Writer, filename string) {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// 获取文件信息
	info, err := file.Stat()
	if err != nil {
		panic(err)
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		panic(err)
	}

	// 设置zip内文件名
	header.Name = filepath.Base(filename)

	// 添加文件到zip
	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		panic(err)
	}

	_, err = io.Copy(writer, file)
	if err != nil {
		panic(err)
	}
}

func (c *controller) writeContent(savePath string, content []byte) {
	// 打开abc.pem文件
	file, err := os.OpenFile(savePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	// 写入内容
	_, err = file.Write(content)
	if err != nil {
		panic(err)
	}

	// 关闭文件
	if err = file.Close(); err != nil {
		panic(err)
	}

	fmt.Println("内容写入成功!")
}

func (c *controller) renameFileZip(oldPath, newPath string) {
	// zip文件路径
	// 重命名为
	// 打开旧zip文件
	oldZip, err := zip.OpenReader(oldPath)
	if err != nil {
		panic(err)
	}
	defer oldZip.Close()

	// 创建新zip文件
	newZip, err := os.Create(newPath)
	if err != nil {
		panic(err)
	}
	defer newZip.Close()

	// 创建zip.Writer
	w := zip.NewWriter(newZip)
	defer w.Close()

	// 遍历旧zip文件中的文件
	for _, file := range oldZip.File {
		// 获取旧zip里的文件名
		fname := file.Name

		// 设置新zip内文件名
		newfile := zip.FileHeader{
			Name:   fname,
			Method: file.Method,
		}

		// 将旧zip内容复制到新zip
		if file.FileInfo().IsDir() {
			w.CreateHeader(&newfile)
		} else {
			in, err := file.Open()
			if err != nil {
				panic(err)
			}
			out, err := w.CreateHeader(&newfile)
			if err != nil {
				panic(err)
			}
			_, err = io.Copy(out, in)
			if err != nil {
				panic(err)
			}
			in.Close()
		}
	}
	// 删除旧zip文件
	os.Remove(oldPath)

}

func (c *controller) DownloadCerts() {
	defer func() {
		r := recover()
		fmt.Println("============panic DownloadCerts============", r)
	}()
	//D:\projects\one\metal\git-init-main\config\certs\certs.zip
	filesZipPath := fmt.Sprintf("%sconfig/%s/", bootstrap.ProjectPath(), "certs")
	// 打开文件
	fmt.Println("filesZipPath, ", filesZipPath)
	files, err := os.ReadDir(filesZipPath)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		name := file.Name()
		// 检查后缀名是否为.zip
		if filepath.Ext(name) == ".zip" {
			// 打印zip包文件名
			//println(name)

			newFilesZipPath := fmt.Sprintf("%sconfig/%s/%s", bootstrap.ProjectPath(), "certs", name)
			file, err := os.Open(newFilesZipPath)
			if err != nil {
				c.GCtx.String(http.StatusNotFound, fmt.Sprintf("file %s not found", "certs.zip"))
				return
			}

			fmt.Println("download file", *file)
			defer file.Close()

			// 获取文件信息
			fileInfo, _ := file.Stat()
			size := fileInfo.Size()
			buffer := make([]byte, size)

			// 读取文件内容到缓冲区
			file.Read(buffer)

			// 设置响应头
			c.GCtx.Writer.Header()["Content-Disposition"] = []string{fmt.Sprintf(`attachment; filename="%s"`, "certs.zip")}
			c.GCtx.Data(http.StatusOK, "application/octet-stream", buffer)

		}
	}

}

// 获取文件MD5方法
func getFileMd5(path string) string {
	// 打开文件
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// 创建MD5哈希计算实例
	md5Hash := md5.New()

	// 读取文件内容并计算MD5
	if _, err := io.Copy(md5Hash, file); err != nil {
		panic(err)
	}

	// 获取MD5字符串
	md5Str := fmt.Sprintf("%x", md5Hash.Sum(nil))

	return md5Str
}

// 清空文件内容方法
func emptyFile(file *os.File) error {
	// 获取当前文件的属性
	stat, err := file.Stat()
	if err != nil {
		return err
	}

	// 获取文件大小
	size := stat.Size()
	// 创建临时文件
	tmpFile, err := os.CreateTemp("", "tmp")
	if err != nil {
		return err
	}
	defer tmpFile.Close()

	// 将原文件内容复制到临时文件
	if _, err = file.Seek(0, 0); err != nil {
		return err
	}
	if _, err = io.CopyN(tmpFile, file, size); err != nil {
		return err
	}

	// 截断原文件
	if err = file.Truncate(0); err != nil {
		return err
	}

	// 重置原文件读取位置
	if _, err = file.Seek(0, 0); err != nil {
		return err
	}

	return nil
}
