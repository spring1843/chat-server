package bootstrap

import (
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/spring1843/chat-server/src/config"
	"github.com/spring1843/chat-server/src/shared/logs"
)

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func serveStaticWeb(mux *http.ServeMux, config config.Config) {
	if config.StaticWeb == "" {
		logs.Infof("Not serving static web files")
		return
	}

	absolutePath, err := filepath.Abs(filepath.Join(config.CWD, config.StaticWeb))
	if err != nil {
		logs.Errf("Error finding absolute path of %q + %q", config.CWD, config.StaticWeb)
	}

	_, err = os.Stat(absolutePath)
	if os.IsNotExist(err) {
		logs.Errf("Directory for StaticWeb defined in config does not exist. CWD %s Absolute Path %s", config.CWD, absolutePath)
		return
	}
	logs.Infof("Serving static web files from %s", absolutePath)
	fs := http.FileServer(http.Dir(absolutePath))
	mux.Handle("/", fs)
}
