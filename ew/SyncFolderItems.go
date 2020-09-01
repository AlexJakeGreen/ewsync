package ew

import (
	"bytes"
	"encoding/xml"
	"text/template"
)

var syncFolderItemsRequestTemplate string = `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
               xmlns:m="http://schemas.microsoft.com/exchange/services/2006/messages"
               xmlns:t="http://schemas.microsoft.com/exchange/services/2006/types"
               xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <m:SyncFolderItems>
      <m:ItemShape>
        <t:BaseShape>IdOnly</t:BaseShape>
      </m:ItemShape>
      <m:SyncFolderId>
       {{ if eq .FolderId "inbox" }}
        <t:DistinguishedFolderId Id="{{ .FolderId }}"/>
       {{ else }}
        <t:FolderId Id="{{ .FolderId }}" ChangeKey="{{ .ChangeKey }}"/>
       {{ end }}
      </m:SyncFolderId>
      <m:MaxChangesReturned>50</m:MaxChangesReturned>
	  {{ if .SyncState }}<m:SyncState>{{ .SyncState }}</m:SyncState>{{ end }}
      <m:SyncScope>NormalItems</m:SyncScope>
    </m:SyncFolderItems>
  </soap:Body>
</soap:Envelope>`

type SyncFolderItemsResponse struct {
	XMLName xml.Name
	Body    struct {
		XMLName                 xml.Name
		SyncFolderItemsResponse struct {
			XMLName xml.Name `xml:"http://schemas.microsoft.com/exchange/services/2006/messages SyncFolderItemsResponse"`
			*BaseResponseMessage
		}
	} `xml:"Body"`
}

func (c *EW) SyncFolderItems(folderId, changeKey, syncState string) (*SyncFolderItemsResponse, error) {
	res, err := c.syncFolderItems(folderId, changeKey, syncState)
	if err != nil {
		return nil, err
	}

	var v *SyncFolderItemsResponse
	err = xml.Unmarshal(res, &v)
	if err != nil {
		return nil, err
	}

	return v, nil
}

func (c *EW) syncFolderItems(folderId, changeKey, syncState string) ([]byte, error) {
	t, err := template.New("SyncFolderItems").Parse(syncFolderItemsRequestTemplate)
	if err != nil {
		return nil, err
	}

	doc := &bytes.Buffer{}
	err = t.Execute(doc, &struct{ FolderId, ChangeKey, SyncState string }{FolderId: folderId, ChangeKey: changeKey, SyncState: syncState})
	if err != nil {
		return nil, err
	}

	res, err := c.DoCall(doc)
	if err != nil {
		return nil, err
	}

	return res, nil
}
