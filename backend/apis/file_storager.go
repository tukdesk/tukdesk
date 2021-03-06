package apis

import (
	"crypto/md5"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"

	"github.com/tukdesk/tukdesk/backend/models"
)

type FileStorager interface {
	Store(*multipart.FileHeader) (*models.File, error)
}

type sizer interface {
	Size() int64
}

type stater interface {
	Stat() (os.FileInfo, error)
}

var (
	ErrExpectingOSFile = fmt.Errorf("expected an *os.File")
)

const (
	defaultDirPerm  = 0755
	defaultFilePerm = 0666

	pathSeparater = "/"
)

type InternalLocalStorager struct {
	dir string
}

func newInternalLocalStorager(dir string) (FileStorager, error) {
	if err := os.MkdirAll(dir, defaultDirPerm); err != nil {
		return nil, err
	}

	return &InternalLocalStorager{
		dir: dir,
	}, nil
}

func (this *InternalLocalStorager) Store(header *multipart.FileHeader) (*models.File, error) {
	reader, err := header.Open()
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	// 生成 fileDoc
	fileDoc := models.NewFile()
	fileDoc.IsInternal = true
	fileDoc.FileName = header.Filename

	// 获取文件大小
	if rSizer, ok := reader.(sizer); ok {
		fileDoc.FileSize = rSizer.Size()
	} else if rStater, ok := reader.(stater); ok {
		stat, err := rStater.Stat()
		if err != nil {
			return nil, err
		}
		fileDoc.FileSize = stat.Size()
	}

	// TODO: 长度为0的文件?

	// 获取文件类型
	fileDoc.MimeType = detectFileType(reader)

	// 生成 sub dir
	fileKey := fileDoc.Id.Hex()
	subDir := generateSubDir([]byte(fileKey))
	dir := path.Join(this.dir, subDir)
	if err := os.MkdirAll(dir, defaultDirPerm); err != nil {
		return nil, err
	}

	// 完整文件名
	ext := path.Ext(header.Filename)
	if ext != "" {
		fileKey += ext
	}

	filePath := path.Join(dir, fileKey)

	// 存储到本地
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, defaultFilePerm)
	if err != nil {
		return nil, err
	}

	defer f.Close()
	if _, err := io.Copy(f, reader); err != nil {
		return nil, err
	}

	fileDoc.FileKey = path.Join(subDir, fileKey)

	return fileDoc, nil
}

func generateSubDir(fileKey []byte) string {
	b := md5.Sum(fileKey)
	s := fmt.Sprintf("%x", b)
	return s[0:3] + pathSeparater + s[3:6] + pathSeparater + s[6:9]
}

func detectFileType(reader multipart.File) string {
	head := make([]byte, 512)
	reader.Read(head)
	reader.Seek(0, 0)
	return http.DetectContentType(head)
}
