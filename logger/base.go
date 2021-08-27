package logger

import "os"

type LogConfig struct {
	LogFolder     string // thư mục chứa log file. Nếu rỗng có nghĩa là không ghi log ra file
	ErrorTemplate string // tên view template sẽ render error page
	InfoTemplate  string // tên view template sẽ render info page
	Top           int    // số dòng đỉnh stack trace sẽ được in ra
}

var logConfig LogConfig
var logFile *os.File

func Init(_logConfig ...LogConfig) *os.File {
	if len(_logConfig) > 0 {
		logConfig = _logConfig[0]
	} else { //Truyền cấu hình nil thì tạo cấu hình mặc định
		logConfig = LogConfig{
			LogFolder:     "logs/", // thư mục chứa log file. Nếu rỗng có nghĩa là không ghi log ra file
			ErrorTemplate: "error", // tên view template sẽ render error page
			InfoTemplate:  "info",  // tên view template sẽ render info page
			Top:           3,       // số dòng đầu tiên trong stack trace sẽ được giữ lại
		}
	}

	if logConfig.LogFolder != "" {
		deleteZeroByteLogFiles(logConfig.LogFolder)
		logFile = newLogFile(logConfig.LogFolder)
		return logFile
	} else {
		return nil
	}
}
