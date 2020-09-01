package ew

import (
	"bytes"
	"encoding/xml"
	"text/template"
)

var getFolder string = `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
      xmlns:m="http://schemas.microsoft.com/exchange/services/2006/messages"
      xmlns:t="http://schemas.microsoft.com/exchange/services/2006/types"
      xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <m:GetFolder>
      <m:FolderShape>
        <t:BaseShape>AllProperties</t:BaseShape>
      </m:FolderShape>
      <m:FolderIds>
        {{ if .ChangeKey }}
        <t:FolderId Id="{{ .FolderId }}" ChangeKey="{{ .ChangeKey }}"/>
        {{ else }}
        <t:DistinguishedFolderId Id="{{ .FolderId }}"/>
        {{ end }}
      </m:FolderIds>
    </m:GetFolder>
  </soap:Body>
</soap:Envelope>`

type GetFolderResponse struct {
	XMLName xml.Name
	Body    struct {
		XMLName           xml.Name
		GetFolderResponse struct {
			XMLName xml.Name
			*BaseResponseMessage
		} `xml:"GetFolderResponse"`
	} `xml:"Body"`
}

func (c *EW) GetFolder(folderId string, changeKey string) (*GetFolderResponse, error) {
	res, err := c.getFolder(folderId, changeKey)
	if err != nil {
		return nil, err
	}

	v := &GetFolderResponse{}
	err = xml.Unmarshal(res, &v)
	if err != nil {
		return nil, err
	}

	return v, nil
}

func (c *EW) getFolder(folderId string, changeKey string) ([]byte, error) {
	t, err := template.New("GetFolder").Parse(getFolder)
	if err != nil {
		return nil, err
	}

	doc := &bytes.Buffer{}
	err = t.Execute(doc, &struct {
		FolderId  string
		ChangeKey string
	}{FolderId: folderId, ChangeKey: changeKey})
	if err != nil {
		return nil, err
	}

	return c.DoCall(doc)
}
