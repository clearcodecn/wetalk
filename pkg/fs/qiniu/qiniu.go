package qiniu

import (
	"fmt"
	"github.com/clearcodecn/wetalk/pkg/fs"
	"github.com/clearcodecn/wetalk/pkg/util"
	"github.com/gabriel-vasile/mimetype"
	"github.com/pkg/errors"
	"github.com/qiniu/api.v7/auth/qbox"
	"github.com/qiniu/api.v7/storage"
	"golang.org/x/net/context"
	"gopkg.in/yaml.v2"
	"io"
	"mime/multipart"
	"os"
	"strings"
)

var (
	ErrMimeTypeNotAllowed = errors.New("mime-type not to allowed to upload")
	zones                 = map[string]*storage.Zone{
		"huadong": &storage.ZoneHuadong,
		"huanan":  &storage.ZoneHuanan,
		"huabei":  &storage.ZoneHuabei,
	}
)

type Config struct {
	Ak              string   `yaml:"ak" json:"ak"`
	Sk              string   `yaml:"sk" json:"sk"`
	Bucket          string   `yaml:"bucket" json:"bucket"`
	Domain          string   `yaml:"domain" json:"domain"`
	AllowExtensions []string `yaml:"allow_extensions" json:"allow_extensions"`
	Zone            string   `yaml:"zone" json:"zone"`
	Https           bool     `yaml:"https" json:"https"`
	EnableCdnUpload bool     `yaml:"enable_cdn_upload" json:"enable_cdn_upload"`
}

type Uploader struct {
	cfg             *Config
	allowExtensions map[string]struct{}
	uploadToken     string
	mac             *qbox.Mac
	qcfg            storage.Config
}

func (u *Uploader) Init(b []byte) error {
	cfg := new(Config)
	err := yaml.Unmarshal(b, cfg)
	if err != nil {
		return err
	}
	if err := checkConfig(cfg); err != nil {
		return err
	}
	u.allowExtensions = make(map[string]struct{})
	u.cfg = cfg
	for _, ae := range cfg.AllowExtensions {
		u.allowExtensions[ae] = struct{}{}
	}

	qcfg := storage.Config{}
	// 空间对应的机房
	if zone, ok := zones[cfg.Zone]; ok {
		qcfg.Zone = zone
	}
	// 是否使用https域名
	qcfg.UseHTTPS = cfg.Https
	// 上传是否使用CDN上传加速
	qcfg.UseCdnDomains = cfg.EnableCdnUpload
	u.mac = qbox.NewMac(u.cfg.Ak, u.cfg.Sk)
	u.qcfg = qcfg
	return nil
}

func (u *Uploader) UploadHTTP(mh *multipart.FileHeader) (*fs.FileInfo, error) {
	f, err := mh.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()

	_, ext, err := mimetype.DetectReader(f)
	if err != nil {
		return nil, err
	}

	if _, ok := u.allowExtensions[ext]; !ok {
		return nil, ErrMimeTypeNotAllowed
	}

	_, err = f.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}
	key := fmt.Sprintf("%s.%s", util.UUID(), ext)
	ret, err := u.uploadReader(f, key, mh.Size)
	if err != nil {
		return nil, err
	}
	fi := &fs.FileInfo{
		FullURL:  fmt.Sprintf("%s/%s", u.cfg.Domain, key),
		BaseURL:  key,
		Domain:   u.cfg.Domain,
		Hash:     ret.Hash,
		FileSize: mh.Size,
		Filename: key,
	}
	return fi, nil

}

func (u *Uploader) UploadLocal(dst string) (*fs.FileInfo, error) {
	fi, err := os.Stat(dst)
	if err != nil {
		return nil, err
	}
	f, err := os.OpenFile(dst, os.O_RDONLY, 0755)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var ext string
	extArr := strings.Split(fi.Name(), ".")
	if len(extArr) > 1 {
		ext = extArr[len(extArr)-1]
	}
	key := fmt.Sprintf("%s.%s", util.UUID(), ext)
	ret, err := u.uploadReader(f, key, fi.Size())
	if err != nil {
		return nil, err
	}
	retFi := &fs.FileInfo{
		FullURL:  fmt.Sprintf("%s/%s", u.cfg.Domain, key),
		BaseURL:  key,
		Domain:   u.cfg.Domain,
		Hash:     ret.Hash,
		FileSize: fi.Size(),
		Filename: key,
	}
	return retFi, nil
}

func (u *Uploader) uploadReader(r io.Reader, name string, size int64) (*storage.PutRet, error) {
	putPolicy := storage.PutPolicy{
		Scope: u.cfg.Bucket,
	}
	upToken := putPolicy.UploadToken(u.mac)

	formUploader := storage.NewFormUploader(&u.qcfg)
	ret := storage.PutRet{}
	err := formUploader.Put(context.Background(), &ret, upToken, name, r, size, nil)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

func checkConfig(cfg *Config) error {
	if cfg.Ak == "" {
		return errors.New("ak can not be empty")
	}
	if cfg.Sk == "" {
		return errors.New("sk can not be empty")
	}
	if cfg.Domain == "" {
		return errors.New("domain can not be empty")
	}
	if cfg.Https && !strings.HasPrefix(cfg.Domain, "https") {
		return errors.New("invalid https config, domain without prefix of https")
	}
	return nil
}
