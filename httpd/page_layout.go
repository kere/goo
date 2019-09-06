package httpd

import (
	"io"

	"github.com/kere/gno/libs/util"
)

var (
	bHeadBegin  = []byte("<head>\n")
	bHeadEnd    = []byte("</head>\n")
	bTitleBegin = []byte("<title>")
	bTitleEnd   = []byte("</title>\n")

	metaCharset     = []byte("<meta charset=\"UTF-8\">\n")
	bytesHTMLBegin  = []byte("<!DOCTYPE HTML>\n<html lang=\"")
	bytesHTMLBegin2 = []byte("\">\n")
	bytesHTMLEnd    = []byte("</html>\n")
	// BytesHTMLBodyBegin bytes
	bytesHTMLBodyBegin = []byte("\n<body>\n")
	// BytesHTMLBodyEnd bytes
	bytesHTMLBodyEnd = []byte("\n</body>\n")

	bRenderS1 = []byte("\n<script type=\"text/javascript\">var MYENV='")
	bRenderS2 = []byte("'," + PageAccessTokenField + "='")
	bRenderS3 = []byte("';</script>")

	contentTypePage = []byte("text/html; charset=utf-8")
)

// renderPage func
func renderPage(w io.Writer, pd *PageData, bPath []byte) error {
	// <html>
	w.Write(bytesHTMLBegin)
	w.Write(util.Str2Bytes(pd.SiteData.Lang))
	w.Write(bytesHTMLBegin2)

	// head -------------------------
	w.Write(bHeadBegin)
	w.Write(metaCharset)

	w.Write(bTitleBegin)
	w.Write(util.Str2Bytes(pd.Title))
	w.Write(bTitleEnd)

	w.Write(bRenderS1)
	w.Write([]byte(RunMode))
	w.Write(bRenderS2)

	token := buildToken(bPath, pd.SiteData.Secret, pd.SiteData.Nonce)

	w.Write(util.Str2Bytes(token))

	// opt := render.Opt{AssetsURL: pd.AssetsURL, JSVersion: pd.JSVersion, CSSVersion: pd.CSSVersion}

	w.Write(bRenderS3)
	for _, r := range pd.Head {
		if err := r.Render(w); err != nil {
			return err
		}
	}

	for _, r := range pd.CSS {
		if err := r.RenderWith(w, pd); err != nil {
			return err
		}
	}

	if pd.JSPosition == JSPositionHead {
		for _, r := range pd.JS {
			if err := r.RenderWith(w, pd); err != nil {
				return err
			}
		}
	}
	w.Write(bHeadEnd)

	// <body>
	w.Write(bytesHTMLBodyBegin)

	var err error
	for _, r := range pd.Top {
		if err = r.Render(w); err != nil {
			return err
		}
	}

	// if len(pd.Body) == 0 {
	// 	r := NewTemplate(filepath.Join(pd.Dir, pd.Name+defaultTemplateSubfix))
	// 	if err = r.Render(w); err != nil {
	// 		return err
	// 	}
	// } else {
	for _, r := range pd.Body {
		if err = r.Render(w); err != nil {
			return err
		}
	}
	// }

	if pd.JSPosition == JSPositionBottom {
		for _, r := range pd.JS {
			if err := r.RenderWith(w, pd); err != nil {
				return err
			}
		}
	}

	for _, r := range pd.Bottom {
		if err = r.Render(w); err != nil {
			return err
		}
	}

	w.Write(bytesHTMLBodyEnd)
	w.Write(bytesHTMLEnd)

	return nil
}
