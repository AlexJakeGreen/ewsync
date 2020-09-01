package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"strings"

	"./ew"
	"github.com/amalfra/maildir"
)

var svc *ew.EW
var cfg *config

type folderState struct {
	SyncState string     `json:"sync-state"`
	Folder    *ew.Folder `json:"folder"`
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	cfg = readConfig()
}

func main() {

	for _, account := range cfg.Accounts {
		svc = ew.NewEW(account.Login, account.Password)

		for _, f := range account.Folders {
			path := f.Remote
			state := readState(account, path)

			var folder *ew.Folder
			var err error
			if state.Folder == nil || state.Folder.FolderId.Id == "" {
				folder, err = findFolder(svc, path)
				if err != nil {
					panic(err)
				}
			} else {
				resp, err := svc.GetFolder(state.Folder.FolderId.Id, state.Folder.FolderId.ChangeKey)
				if err != nil {
					panic(err)
				}

				folder = resp.Body.GetFolderResponse.ResponseMessages.GetFolderResponseMessage.Folders.Folder[0]
			}
			state.Folder = folder
			saveState(account, path, state)

			syncFolder(folder, account, path, state)
		}
	}
}

func findFolder(svc *ew.EW, path string) (*ew.Folder, error) {
	parts := strings.Split(path, "/")
	root, err := svc.GetFolder(parts[0], "")
	if err != nil {
		return nil, err
	}

	rootFolder := root.Body.GetFolderResponse.ResponseMessages.GetFolderResponseMessage.Folders.Folder[0]
	rootId := rootFolder.FolderId.Id
	folderId := rootId
	result := rootFolder
	for _, part := range parts[1:] {
		fmt.Println(part)
		folder, err := svc.FindFolder(folderId)
		if err != nil {
			return nil, err
		}

		for _, i := range folder.Body.FindFolderResponse.ResponseMessages.FindFolderResponseMessage.RootFolder.Folders.Folder {
			if part == i.DisplayName {
				folderId = i.FolderId.Id
				result = i
			}
		}
	}

	return result, nil
}

func syncFolder(folder *ew.Folder, account *Account, path string, state *folderState) {
	log.Printf("Syncing folder %s\n", path)

	for {
		v, err := svc.SyncFolderItems(folder.FolderId.Id, folder.FolderId.ChangeKey, state.SyncState)
		if err != nil {
			panic(err)
		}

		d := maildir.NewMaildir(account.Maildir + "/" + path)

		resp := v.Body.SyncFolderItemsResponse.ResponseMessages.SyncFolderItemsResponseMessage
		state.SyncState = resp.SyncState
		items := make([]*ew.ItemId, 0)
		for _, item := range resp.Changes.Create {
			if item.Message == nil {
				continue
			}
			items = append(items, &item.Message.ItemId)
		}

		log.Printf("==> getting %d items \n", len(items))
		v2, err := svc.GetItem(items)
		if err != nil {
			panic(err)
		}

		for _, m := range v2.Body.GetItemResponse.ResponseMessages.GetItemResponseMessage.Items.Message {
			decoded, err := base64.StdEncoding.DecodeString(m.MimeContent)
			if err != nil {
				fmt.Println("decode error:", err)
				return
			}

			_, err = d.Add(string(decoded))
			if err != nil {
				panic(err)
			}

		}

		if resp.IncludesLastItemInRange {
			break
		}
	}

	saveState(account, path, state)
}
