package adapter

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/ikmski/git-lfs3/usecase"
)

type transferController struct {
	transferService usecase.TransferService
}

// TransferController is ...
type TransferController interface {
	Download(ctx Context)
	Upload(ctx Context)
}

// NewTransferController is ...
func NewTransferController(s usecase.TransferService) TransferController {
	return &transferController{
		transferService: s,
	}
}

func (c *transferController) Download(ctx Context) {

	or := parseObjectRequest(ctx)

	exists := c.transferService.Exists(or)
	if !exists {
		ctx.SetStatus(404)
		return
	}

	rangeHeader := ctx.GetHeader("Range")
	if rangeHeader != "" {

		size := c.transferService.GetSize(or)
		var fromByte int64
		var toByte int64 = size
		regex := regexp.MustCompile(`bytes=(.*)\-(.*)`)
		match := regex.FindStringSubmatch(rangeHeader)
		if match != nil && len(match) >= 3 {
			if len(match[1]) > 0 {
				fromByte, _ = strconv.ParseInt(match[1], 10, 64)
			}
			if len(match[2]) > 0 {
				toByte, _ = strconv.ParseInt(match[2], 10, 64)
			}
		}
		or.From = fromByte
		or.To = toByte

		ctx.SetHeader("Content-Range", fmt.Sprintf("bytes %d-%d/%d", fromByte, toByte-1, int64(toByte-fromByte)))
		ctx.SetStatus(206)
	}

	_, err := c.transferService.Download(or, ctx.GetResponseWriter())
	if err != nil {
		ctx.SetStatus(404)
		return
	}
}

func (c *transferController) Upload(ctx Context) {

	o := parseObjectRequest(ctx)
	exists := c.transferService.Exists(o)
	if !exists {
		ctx.SetStatus(404)
		return
	}

	err := c.transferService.Upload(o, ctx.GetRequestReader())
	if err != nil {
		ctx.SetStatus(500)
		//fmt.Fprintf(c.Writer, `{"message":"%s"}`, err)
		return
	}
}

func parseObjectRequest(ctx Context) *usecase.ObjectRequest {

	oid := ctx.GetParam("oid")

	or := &usecase.ObjectRequest{
		Oid: oid,
	}

	return or
}

/*
func (o *ObjectRequest) DownloadLink() string {

	return o.internalLink("objects")
}

func (o *ObjectRequest) UploadLink() string {

	return o.internalLink("objects")
}

func (o *ObjectRequest) internalLink(subpath string) string {

	path := ""

	if len(o.User) > 0 {
		path += fmt.Sprintf("/%s", o.User)
	}

	if len(o.Repo) > 0 {
		path += fmt.Sprintf("/%s", o.Repo)
	}

	path += fmt.Sprintf("/%s/%s", subpath, o.Oid)

	if config.Server.Tls {
		return fmt.Sprintf("https://%s%s", config.Server.Host, path)
	}

	return fmt.Sprintf("http://%s%s", config.Server.Host, path)
}

func (o *ObjectRequest) VerifyLink() string {

	path := fmt.Sprintf("/verify/%s", o.Oid)

	if config.Server.Tls {
		return fmt.Sprintf("https://%s%s", config.Server.Host, path)
	}

	return fmt.Sprintf("http://%s%s", config.Server.Host, path)
}
*/
