package ew

import (
	"bytes"
	"encoding/xml"
	"text/template"
	"time"
)

var getItemRequestTemplate string = `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope
  xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
  xmlns:xsd="http://www.w3.org/2001/XMLSchema"
  xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/"
  xmlns:t="http://schemas.microsoft.com/exchange/services/2006/types">
  <soap:Body>
    <GetItem
      xmlns="http://schemas.microsoft.com/exchange/services/2006/messages"
      xmlns:t="http://schemas.microsoft.com/exchange/services/2006/types">
      <ItemShape>
        <t:BaseShape>IdOnly</t:BaseShape>
        <t:IncludeMimeContent>true</t:IncludeMimeContent>
      </ItemShape>
      <ItemIds>
        {{ range .Items }}
        <t:ItemId Id="{{ .Id }}" ChangeKey="{{ .ChangeKey}}" />
        {{ end }}
      </ItemIds>
    </GetItem>
  </soap:Body>
</soap:Envelope>`

type ItemId struct {
	Id        string `xml:"Id,attr"`
	ChangeKey string `xml:"ChangeKey,attr"`
}

type Change struct {
	Create Create `xml:"Create"`
}

type Create struct {
	Message struct {
		ItemId ItemId `xml:"ItemId"`
	}
}

type Message struct {
	ItemId      ItemId `xml:"ItemId"`
	MimeContent string
	Subject     string
	Sensitivity string
	Body        string
	Attachments struct {
		FileAttachment struct {
			AttachmentId struct {
				Id string `xml:"Id,attr"`
			}
			Name        string
			ContentType string
			ContentId   string
		}
	}
	Size            int
	DateTimeSent    time.Time
	DateTimeCreated time.Time
	HasAttachments  bool
	From            struct {
		Mailbox struct {
			Name         string
			EmailAddress string
			RoutingType  string
		}
	}
	IsRead bool
}

type MessageResponse struct {
	XMLName xml.Name
	Body    struct {
		XMLName         xml.Name
		GetItemResponse struct {
			XMLName          xml.Name
			ResponseMessages struct {
				XMLName                xml.Name
				GetItemResponseMessage struct {
					XMLName       xml.Name
					ResponseClass string `xml:"ResponseClass,attr"`
					ResponseCode  string `xml:"ResponseCode"`
					Items         struct {
						XMLName xml.Name
						Message []*Message `xml:"Message"`
					} `xml:"Items"`
				} `xml:"GetItemResponseMessage"`
			} `xml:"ResponseMessages"`
		} `xml:"GetItemResponse"`
	} `xml:"Body"`
}

func (c *EW) GetItem(items []*ItemId) (*MessageResponse, error) {
	res, err := c.getItem(items)
	if err != nil {
		return nil, err
	}

	v := &MessageResponse{}
	err = xml.Unmarshal(res, &v)
	if err != nil {
		return nil, err
	}

	return v, nil
}

func (c *EW) getItem(items []*ItemId) ([]byte, error) {
	t, err := template.New("GetItem").Parse(getItemRequestTemplate)
	if err != nil {
		return nil, err
	}

	doc := &bytes.Buffer{}
	err = t.Execute(doc, &struct{ Items []*ItemId }{Items: items})
	if err != nil {
		return nil, err
	}

	res, err := c.DoCall(doc)
	if err != nil {
		return nil, err
	}

	return res, nil
}
