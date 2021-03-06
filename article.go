package html2article

import (
	"net/url"
	"path"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type Article struct {
	// Basic
	Html        string `json:"content_html"`
	Content     string `json:"content"`
	Title       string `json:"title"`
	Publishtime int64  `json:"publish_time"`

	// Others
	Images      []string `json:"images"`
	ReadContent string   `json:"read_content"`
	contentNode *html.Node
}

func (a *Article) Readable(urlStr string) {
	a.ParseReadContent()
	a.ParseImage(urlStr)
}

// ParseReadContent parse the ReadContent to be readability
func (a *Article) ParseReadContent() {
	a.cleanStyle(a.contentNode, "class", "id", "style", "width", "height", "onclick", "onmouseover", "border")
	a.clean(a.contentNode, atom.Object, atom.H1)
	a.ReadContent, _ = getHtml(a.contentNode)
}

// ParseImage parse the image src to the absolute path
func (a *Article) ParseImage(urlStr string) {
	_url, err := url.Parse(urlStr)
	if err != nil {
		return
	}
	mp := make(map[string]string)
	for i, _ := range a.Images {
		if strings.Index(a.Images[i], "http") != 0 {
			var newImg string
			if strings.Index(a.Images[i], "//") == 0 {
				newImg = _url.Scheme + ":" + a.Images[i]
			} else if strings.Index(a.Images[i], "/") == 0 {
				newImg = _url.Scheme + "://" + _url.Host + a.Images[i]
			} else {
				newImg = _url.Scheme + "://" + _url.Host + path.Join(_url.Path, "../", a.Images[i])
			}
			mp[a.Images[i]] = newImg
			a.Images[i] = newImg
		}
	}
	for k, v := range mp {
		a.Html = strings.Replace(a.Html, k, v, -1)
		a.ReadContent = strings.Replace(a.ReadContent, k, v, -1)
	}
}

func (a *Article) clean(sel *html.Node, tags ...atom.Atom) {
	for c := sel.FirstChild; c != nil; c = c.NextSibling {
		a.clean(c, tags...)
		for _, tag := range tags {
			if isTag(tag)(c) {
				pre := c.PrevSibling
				sel.RemoveChild(c)
				c = pre
				break
			}
		}
		if c == nil {
			break
		}
	}
}

func (a *Article) cleanStyle(sel *html.Node, attrs ...string) {
	for _, attr := range attrs {
		removeAttr(sel, attr)
	}
	for c := sel.FirstChild; c != nil; c = c.NextSibling {
		a.cleanStyle(c, attrs...)
	}
}
