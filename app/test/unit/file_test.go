package unit

import (
	"educ-gpt/config/dic"
	"educ-gpt/services"
)

var (
	fileSrv services.FileService
)

func TestInitFileService() {
	fileSrv = dic.FileService()
}
