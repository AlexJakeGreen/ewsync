package ew

import (
	"bytes"
	"encoding/xml"
	"text/template"
)

var findFolderRequestTemplate string = `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
               xmlns:m="http://schemas.microsoft.com/exchange/services/2006/messages"
               xmlns:t="http://schemas.microsoft.com/exchange/services/2006/types"
               xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <m:FindFolder Traversal="Shallow">
      <m:FolderShape>
        <t:BaseShape>Default</t:BaseShape>
      </m:FolderShape>
      {{ if .FolderId }}
      <m:ParentFolderIds>
        <t:{{ if eq .FolderId "inbox" }}Distinguished{{ end }}FolderId Id="{{ .FolderId }}" />
      </m:ParentFolderIds>
      {{ end }}
    </m:FindFolder>
  </soap:Body>
</soap:Envelope>`

type FindFolderResponse struct {
	XMLName xml.Name
	Body    struct {
		XMLName            xml.Name
		FindFolderResponse struct {
			XMLName xml.Name
			*BaseResponseMessage
		} `xml:"FindFolderResponse"`
	} `xml:"Body"`
}

func (c *EW) FindFolder(folderId string) (*FindFolderResponse, error) {
	res, err := c.findFolder(folderId)
	if err != nil {
		return nil, err
	}
	// fmt.Println(string(res))
	v := &FindFolderResponse{}
	err = xml.Unmarshal(res, &v)
	if err != nil {
		return nil, err
	}

	return v, nil
}

func (c *EW) findFolder(rootFolder string) ([]byte, error) {
	t, err := template.New("FindFolder").Parse(findFolderRequestTemplate)
	if err != nil {
		return nil, err
	}

	doc := &bytes.Buffer{}
	err = t.Execute(doc, &struct{ FolderId string }{FolderId: rootFolder})
	if err != nil {
		return nil, err
	}

	return c.DoCall(doc)
}
