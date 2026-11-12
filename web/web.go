// Package web 把构建好的控制台静态资源嵌入二进制，实现单文件分发。
package web

import (
	"embed"
	"io/fs"
	"net/http"
	"strings"
)

//go:embed all:dist
var distFS embed.FS

// Handler 返回控制台静态资源处理器；调用方需先用 StripPrefix 去掉挂载前缀。
// 对不存在的路径回退到 index.html，以支持前端单页路由。
func Handler() http.Handler {
	sub, err := fs.Sub(distFS, "dist")
	if err != nil {
		panic(err)
	}
	indexHTML, err := fs.ReadFile(sub, "index.html")
	if err != nil {
		panic(err)
	}
	fileServer := http.FileServer(http.FS(sub))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rel := strings.TrimPrefix(r.URL.Path, "/")

		// 根目录：交给 FileServer，它会自动输出该目录的 index.html。
		if rel == "" {
			r.URL.Path = "/"
			fileServer.ServeHTTP(w, r)
			return
		}

		// 资源不存在：直接写回 index.html（不能用 /index.html 路径，
		// 标准库会把它 301 到目录）。
		if _, statErr := fs.Stat(sub, rel); statErr != nil {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			_, _ = w.Write(indexHTML)
			return
		}

		r.URL.Path = "/" + rel
		fileServer.ServeHTTP(w, r)
	})
}
