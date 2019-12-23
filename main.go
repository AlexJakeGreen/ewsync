package main

import (
	"bytes"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/amalfra/maildir"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func setSyncState(state string) {
	log.Printf("setSyncState: %s\n", state)
	b := []byte(state)
	err := ioutil.WriteFile(".maildir_state", b, 0644)
	if err != nil {
		panic(err)
	}
}

func getSyncState() string {
	if _, err := os.Stat(".maildir_state"); os.IsNotExist(err) {
		return ""
	}

	data, err := ioutil.ReadFile(".maildir_state")
	if err != nil {
		panic(err)
	}

	return string(data)
}

func post(message string) ([]byte, error) {
	client := &http.Client{}

	body := bytes.NewBuffer([]byte(message))
	req, err := http.NewRequest("POST", "https://outlook.office365.com/EWS/Exchange.asmx", body)
	req.SetBasicAuth("user@domain.com", "password")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	res, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return res, nil
}

type ItemId struct {
	Id        string `xml:"Id,attr"`
	ChangeKey string `xml:"ChangeKey,attr"`
}
type Create struct {
	Message struct {
		ItemId ItemId `xml:"ItemId"`
	}
}

type SyncFolderItemsResponse struct {
	XMLName xml.Name
	Body    struct {
		SyncFolderItemsResponse struct {
			ResponseMessages struct {
				SyncFolderItemsResponseMessage struct {
					ResponseCode            string
					SyncState               string
					IncludesLastItemInRange bool
					Changes                 struct {
						Create []Create
					}
				}
			}
		}
	}
}

func syncFolder() int {
	syncStateElement := ""
	syncState := getSyncState()
	if syncState != "" {
		syncStateElement = fmt.Sprintf("<m:SyncState>%s</m:SyncState>", syncState)
	}

	rawRequest := `<?xml version="1.0" encoding="utf-8"?>
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
        <t:DistinguishedFolderId Id="inbox" />
      </m:SyncFolderId>
      <m:MaxChangesReturned>10</m:MaxChangesReturned>` +
		syncStateElement +
		`<m:SyncScope>NormalItems</m:SyncScope>
    </m:SyncFolderItems>
  </soap:Body>
</soap:Envelope>`

	rawResponse, err := post(rawRequest)
	if err != nil {
		panic(err)
	}

	v := SyncFolderItemsResponse{}
	err = xml.Unmarshal(rawResponse, &v)
	if err != nil {
		fmt.Println(err.Error())
	}

	//fmt.Printf("%+v\n", v)
	for _, item := range v.Body.SyncFolderItemsResponse.ResponseMessages.SyncFolderItemsResponseMessage.Changes.Create {
		getItem(item.Message.ItemId.Id, item.Message.ItemId.ChangeKey)
	}

	setSyncState(v.Body.SyncFolderItemsResponse.ResponseMessages.SyncFolderItemsResponseMessage.SyncState)

	return len(v.Body.SyncFolderItemsResponse.ResponseMessages.SyncFolderItemsResponseMessage.Changes.Create)
}

type Message struct {
	ItemId struct {
		Id        string `xml:"Id,attr"`
		ChangeKey string `xml:"ChangeKey,attr"`
	} `xml:"ItemId"`
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
						Message Message `xml:"Message"`
					} `xml:"Items"`
				} `xml:"GetItemResponseMessage"`
			} `xml:"ResponseMessages"`
		} `xml:"GetItemResponse"`
	} `xml:"Body"`
}

func getItem(id, key string) {
	fmt.Printf("Fetching %s %s\n", id, key)
	rawRequest := fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
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
        <t:ItemId Id="%s" ChangeKey="%s" />
      </ItemIds>
    </GetItem>
  </soap:Body>
</soap:Envelope>`, id, key)

	rawResponse, err := post(rawRequest)
	if err != nil {
		panic(err)
	}

	//fmt.Println(string(rawResponse))
	v := MessageResponse{}
	err = xml.Unmarshal(rawResponse, &v)
	if err != nil {
		fmt.Println(err.Error())
	}

	//fmt.Printf("%+v\n", v)
	m := v.Body.GetItemResponse.ResponseMessages.GetItemResponseMessage.Items.Message.MimeContent
	decoded, err := base64.StdEncoding.DecodeString(m)
	if err != nil {
		fmt.Println("decode error:", err)
		return
	}
	//fmt.Println(string(decoded))

	d := maildir.NewMaildir(".maildir")
	_, err = d.Add(string(decoded))
	if err != nil {
		panic(err)
	}
}

func main() {
	for {
		n := syncFolder()
		fmt.Printf("Sync %d items\n", n)
		if n < 1 {
			break
		}
	}
}
