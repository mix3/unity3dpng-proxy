package main

import (
	"bytes"
	"html/template"
	"mime/multipart"
	"net/http"
	"net/http/httputil"
	"os"
	"os/exec"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/k0kubun/pp"
	"github.com/kelseyhightower/envconfig"
	"github.com/mholt/binding"
)

type Conf struct {
	Git  string `envconfig:"GIT"  default:"git"`
	Addr string `envconfig:"ADDR" default:":19300"`

	WorkDir        string `envconfig:"WORKDIR"           required:"true"`
	Repository     string `envconfig:"REPOSITORY"        required:"true"`
	Unity3d2PngURL string `envconfig:"UNITY3D2PNG_URL"   required:"true"`
	PathTemplate   string `envconfig:"PATH_TEMPLATE"     required:"true"`
}

var (
	conf         Conf
	git          *Git
	pathTemplate *template.Template

	logger = logrus.New()
)

func init() {
	if err := envconfig.Process("unity3d2png_proxy", &conf); err != nil {
		logger.Warn(err)
		os.Exit(1)
	}

	git = NewGit(conf.Git, conf.Repository, conf.WorkDir, logger)

	if !which(conf.Git) {
		logger.Warn("git not found:", conf.Git)
		os.Exit(1)
	}
	if _, dirErr := os.Stat(conf.WorkDir); dirErr != nil {
		if cloneErr := clone(); cloneErr != nil {
			logger.Warn(dirErr)
			logger.Warn(cloneErr)
			os.Exit(1)
		}
	}

	pathTemplate = template.Must(template.New("path_template").Parse(conf.PathTemplate))

	logger.Info(pp.Sprint(conf))
}

func which(cmd string) bool {
	return exec.Command("which", cmd).Run() == nil
}

type Form struct {
	Branch string
	File   string
	Fetch  bool
}

func (f *Form) FieldMap(r *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&f.Branch: binding.Field{
			Form:     "branch",
			Required: true,
		},
		&f.File: binding.Field{
			Form:     "file",
			Required: true,
		},
		&f.Fetch: binding.Field{
			Form: "fetch",
		},
	}
}

func (f *Form) Validate(r *http.Request, errs binding.Errors) binding.Errors {
	if !strings.HasSuffix(f.File, ".unity3d") {
		errs = append(errs, binding.Error{
			FieldNames:     []string{"file"},
			Classification: "SuffixError",
			Message:        "require suffix `.unity3d`",
		})
	}
	return errs
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		form := new(Form)
		errs := binding.Bind(r, form)
		if errs.Handle(w) {
			logger.Warn(errs)
			return
		}

		var pb bytes.Buffer
		if err := pathTemplate.Execute(&pb, r.URL.Query()); err != nil {
			logger.Warn(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		path := string(pb.Bytes())

		if form.Fetch {
			if err := fetch(); err != nil {
				logger.Warn(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		if !search(form.Branch, path) {
			logger.Warn("not found")
			http.NotFound(w, r)
			return
		}

		out, err := output(form.Branch, path)
		if err != nil {
			logger.Warn(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var b bytes.Buffer
		mw := multipart.NewWriter(&b)
		fw, err := mw.CreateFormFile("file", path)
		if err != nil {
			logger.Warn(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if _, err = fw.Write(out); err != nil {
			logger.Warn(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		mw.Close()

		proxy := httputil.ReverseProxy{Director: func(req *http.Request) {
			pr, _ := http.NewRequest("POST", conf.Unity3d2PngURL, &b)
			req.Method = pr.Method
			req.URL.Scheme = pr.URL.Scheme
			req.URL.Host = pr.URL.Host
			req.Header.Set("Content-Type", mw.FormDataContentType())
			req.ContentLength = pr.ContentLength
			req.Body = pr.Body
		}}
		proxy.ServeHTTP(w, r)
	})
	http.ListenAndServe(conf.Addr, nil)
}

func clone() error {
	stdout, stderr, err := git.Clone()
	if 0 < len(stdout) {
		logger.Info(string(stdout))
	}
	if 0 < len(stderr) {
		logger.Info(string(stderr))
	}
	return err
}

func fetch() error {
	stdout, stderr, err := git.Run("fetch")
	if 0 < len(stdout) {
		logger.Info(string(stdout))
	}
	if 0 < len(stderr) {
		logger.Info(string(stderr))
	}
	return err
}

func search(branch, path string) bool {
	stdout, stderr, err := git.Run(
		"cat-file",
		"-e",
		"remotes/origin/"+branch+":"+path,
	)
	if 0 < len(stdout) {
		logger.Info(string(stdout))
	}
	if 0 < len(stderr) {
		logger.Info(string(stderr))
	}
	return err == nil
}

func output(branch, path string) ([]byte, error) {
	stdout, stderr, err := git.Run(
		"cat-file",
		"-p",
		"remotes/origin/"+branch+":"+path,
	)
	if 0 < len(stderr) {
		logger.Info(string(stderr))
	}
	return stdout, err
}
