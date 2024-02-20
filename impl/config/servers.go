package config

import "net/http"

type (
	IFileServer interface {
		UploadFile(w http.ResponseWriter, r *http.Request)
		ListFile(w http.ResponseWriter, r *http.Request)
		DeleteFile(w http.ResponseWriter, r *http.Request)
	}
	IProcessServer interface {
		StartProc(request *StartProcRequest, response *StartProcResponse) error
		GetProcLog(request *GetProcLogRequest, response *GetProcLogResponse) error
		KillProc(request *KillProcLogRequest, response *KillProcLogResponse) error
	}
)
