package logs

import (
	"github.com/ziliscite/messaging-app/pkg/must"
	"io"
	"os"
)

func Set() io.Writer {
	must.MustServe(os.MkdirAll("./logs", 0755))

	logFile := must.Must(os.OpenFile("./logs/messaging_log.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666))
	return io.MultiWriter(os.Stdout, logFile)
}
