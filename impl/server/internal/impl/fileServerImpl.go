package impl

import (
	"ProcessIsolator/impl/config"
	"ProcessIsolator/impl/server/internal"
	"fmt"
	"mime/multipart"
	"net/http"
)

type FileServerImpl struct {
	config.IFileServer
	errorChan chan<- error
	msgChan   chan<- string
}

func (impl *FileServerImpl) Init(errorChan chan<- error, msgChan chan<- string) {
	impl.errorChan = errorChan
	impl.msgChan = msgChan
}

func (impl *FileServerImpl) UploadFile(w http.ResponseWriter, r *http.Request) {
	name := r.Header.Get("name")
	file, _, err := r.FormFile("file")
	if name == "" {
		http.Error(w, "param name is \"\" ", http.StatusBadRequest)
		impl.errorChan <- fmt.Errorf("http %s error %s", "UploadFile", "param name is \"\" ")
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		impl.errorChan <- fmt.Errorf("http %s error %s", "UploadFile", err.Error())
		return
	}
	defer func(file multipart.File) {
		_ = file.Close()
	}(file)

	impl.msgChan <- fmt.Sprintf("http %s %s start", "UploadFile", name)
	err = internal.WriteFile(name, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		impl.errorChan <- fmt.Errorf("http %s error %s", "UploadFile", err.Error())
		return
	}

	_, _ = fmt.Fprintf(w, "文件上传成功: %s\n", name)
	impl.msgChan <- fmt.Sprintf("http %s %s success", "UploadFile", name)
}

func (impl *FileServerImpl) ListFile(w http.ResponseWriter, r *http.Request) {
	// todo: listFile
}

func (impl *FileServerImpl) DeleteFile(w http.ResponseWriter, r *http.Request) {
	name := r.Header.Get("name")
	if name == "" {
		http.Error(w, "param name is \"\" ", http.StatusBadRequest)
		impl.errorChan <- fmt.Errorf("http %s error %s", "DeleteFile", "param name is \"\" ")
		return
	}

	impl.msgChan <- fmt.Sprintf("http %s %s start", "DeleteFile", name)
	err := internal.DeleteFile(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		impl.errorChan <- fmt.Errorf("http %s error %s", "DeleteFile", err.Error())
		return
	}

	_, _ = fmt.Fprintf(w, "文件删除成功: %s\n", name)
	impl.msgChan <- fmt.Sprintf("http %s %s success", "DeleteFile", name)
}
