package main

import (
	"fmt"
	"os"
	"flag"
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"text/template"

	"github.com/stretchr/objx"
	"github.com/joho/godotenv"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/gomniauth/providers/google"
)

// アクティブなAvatar
// ユーザーのアバターのURLが必要な時に使用される。
var avatars Avatar = TryAvatars{
	UseFileSystemAvatar,
	UseAuthAvatar,
	UseGravatar,
}

type templateHandler struct {
	once     sync.Once
	filename string
	temp1    *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// once でtemplateを一度だけコンパイルする様にしている。
	t.once.Do(func() {
		t.temp1 =
			template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	data := map[string]interface{}{
		"Host": r.Host,
	}
	if authCookie, err := r.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}
	// 第二引数にrをしていする事で、httpを出力する際にhttp.Requestい含まれるデータを参照できる様にする。
	// ポート番号の情報もここに格納されている。
	t.temp1.Execute(w, data)
}
func main() {
	// envファイル読み込み
	err := godotenv.Load(fmt.Sprintf("../%s.env", os.Getenv("GO_ENV")))
	if err != nil {
		fmt.Println(err)
	}

	// フラグのデフォルト値に8080を指定
	var addr = flag.String("addr", ":8080", "アプリケーションのアドレス")
	// フラグを解釈
	// コマンドラインで指定された文字列から必要な情報を取り出し *addrにセット
	flag.Parse()
	// Gomniauthのセットアップ
	gomniauth.SetSecurityKey("tom1233908745aizu")
	gomniauth.WithProviders(
		google.New(os.Getenv("GOOGLECLIENTID"), os.Getenv("GOOGLECLIENTSECRET"), "http://localhost:8080/auth/callback/google"),
		github.New(os.Getenv("GITHUBCLIENTID"), os.Getenv("GITHUBCLIENTSECRET"), "http://localhost:8080/auth/callback/github"),
	)
	// ここで、様々な画像の適応の仕方が出来る。
	r := newRoom(avatars)
	// r.tracer = trace.New(os.Stdout) traceにOff() を定義してい無い場合、これを使うとターミなる上でログが表示される。

	//　各エンドポイントとそれに対応するファイルの対応付け
	// URLとハンドラーをDefaultSeveMuxに登録する。
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room", r)
	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:   "auth",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		})
		w.Header()["Location"] = []string{"/chat"}
		w.WriteHeader(http.StatusTemporaryRedirect)
	})
	http.Handle("/upload", &templateHandler{filename: "upload.gtpl"})
	http.HandleFunc("/uploader", uploaderHandler)
	http.Handle("/avatars/", http.StripPrefix("/avatars/", http.FileServer(http.Dir("./avatars"))))
	// チャットルームを開始
	go r.run()
	// Webサーバーを起動
	// port番号をターミナルに表示
	log.Println("Webサーバーを開始します。ポート:", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}

}
