package ew

import (
	"time"
)

type FolderId struct {
	Id        string `xml:"Id,attr,omitempty" json:"Id,omitempty"`
	ChangeKey string `xml:"ChangeKey,attr,omitempty" json:"ChangeKey,omitempty"`
}

type BaseFolderType struct {
	FolderId         *FolderId               `xml:"FolderId,omitempty" json:"FolderId,omitempty"`
	ParentFolderId   *FolderId               `xml:"ParentFolderId,omitempty" json:"ParentFolderId,omitempty"`
	FolderClass      string                  `xml:"FolderClass,omitempty" json:"FolderClass,omitempty"`
	DisplayName      string                  `xml:"DisplayName,omitempty" json:"DisplayName,omitempty"`
	TotalCount       int                     `xml:"TotalCount,omitempty" json:"TotalCount,omitempty"`
	ChildFolderCount int                     `xml:"ChildFolderCount,omitempty" json:"ChildFolderCount,omitempty"`
	ExtendedProperty []*ExtendedPropertyType `xml:"ExtendedProperty,omitempty" json:"ExtendedProperty,omitempty"`
}

type Folder struct {
	*BaseFolderType
	Parent   *Folder   `xml:"-"`
	Children []*Folder `xml:"-"`
	// PermissionSet *PermissionSetType `xml:"PermissionSet,omitempty"`
	UnreadCount int `xml:"UnreadCount,omitempty"`
}

type ExtendedPropertyType struct {
	// ExtendedFieldURI *PathToExtendedFieldType `xml:"ExtendedFieldURI,omitempty" json:"ExtendedFieldURI,omitempty"`
	Value string `xml:"Value,omitempty" json:"Value,omitempty"`
	// Values *NonEmptyArrayOfPropertyValuesType `xml:"Values,omitempty" json:"Values,omitempty"`
}

type BaseResponseMessage struct {
	ResponseMessages *ArrayOfResponseMessages `xml:"ResponseMessages,omitempty" json:"ResponseMessages,omitempty"`
}

type ArrayOfResponseMessages struct {
	GetFolderResponseMessage       *FolderInfoResponseMessage      `xml:"GetFolderResponseMessage,omitempty"`
	SyncFolderItemsResponseMessage *SyncFolderItemsResponseMessage `xml:"SyncFolderItemsResponseMessage,omitempty"`
	FindFolderResponseMessage      *FindFolderResponseMessage      `xml:"FindFolderResponseMessage,omitempty"`
}

type FindFolderResponseMessage struct {
	*ResponseMessage
	RootFolder *FindFolderParent `xml:"RootFolder,omitempty"`
}

type FindFolderParent struct {
	Folders *ArrayOfFolders `xml:"Folders,omitempty"`
}

type SyncFolderItemsResponseMessage struct {
	*ResponseMessage
	SyncState               string                  `xml:"SyncState,omitempty"`
	IncludesLastItemInRange bool                    `xml:"IncludesLastItemInRange,omitempty"`
	Changes                 *SyncFolderItemsChanges `xml:"Changes,omitempty"`
}

type SyncFolderItemsChanges struct {
	Create []*SyncFolderItemsCreateOrUpdate `xml:"Create,omitempty" json:"Create,omitempty"`
	// Update *SyncFolderItemsCreateOrUpdateType `xml:"Update,omitempty" json:"Update,omitempty"`
	// Delete *SyncFolderItemsDeleteType `xml:"Delete,omitempty" json:"Delete,omitempty"`
	// ReadFlagChange *SyncFolderItemsReadFlagType `xml:"ReadFlagChange,omitempty" json:"ReadFlagChange,omitempty"`
}

type SyncFolderItemsCreateOrUpdate struct {
	Item    *Item    `xml:"Item,omitempty"`
	Message *Message `xml:"Message,omitempty" json:"Message,omitempty"`
}

type Item struct {
	// MimeContent *MimeContentType `xml:"MimeContent,omitempty" json:"MimeContent,omitempty"`
	// ItemId *ItemIdType `xml:"ItemId,omitempty" json:"ItemId,omitempty"`
	ParentFolderId *FolderId `xml:"ParentFolderId,omitempty" json:"ParentFolderId,omitempty"`
	// ItemClass *ItemClassType `xml:"ItemClass,omitempty" json:"ItemClass,omitempty"`
	Subject string `xml:"Subject,omitempty" json:"Subject,omitempty"`
	// Sensitivity *SensitivityChoicesType `xml:"Sensitivity,omitempty" json:"Sensitivity,omitempty"`
	// Body *BodyType `xml:"Body,omitempty" json:"Body,omitempty"`
	// Attachments *NonEmptyArrayOfAttachmentsType `xml:"Attachments,omitempty" json:"Attachments,omitempty"`
	DateTimeReceived time.Time `xml:"DateTimeReceived,omitempty" json:"DateTimeReceived,omitempty"`
	Size             int32     `xml:"Size,omitempty" json:"Size,omitempty"`
	// Categories *ArrayOfStringsType `xml:"Categories,omitempty" json:"Categories,omitempty"`
	// Importance *ImportanceChoicesType `xml:"Importance,omitempty" json:"Importance,omitempty"`
	InReplyTo       string    `xml:"InReplyTo,omitempty" json:"InReplyTo,omitempty"`
	IsSubmitted     bool      `xml:"IsSubmitted,omitempty" json:"IsSubmitted,omitempty"`
	IsDraft         bool      `xml:"IsDraft,omitempty" json:"IsDraft,omitempty"`
	IsFromMe        bool      `xml:"IsFromMe,omitempty" json:"IsFromMe,omitempty"`
	IsResend        bool      `xml:"IsResend,omitempty" json:"IsResend,omitempty"`
	IsUnmodified    bool      `xml:"IsUnmodified,omitempty" json:"IsUnmodified,omitempty"`
	DateTimeSent    time.Time `xml:"DateTimeSent,omitempty" json:"DateTimeSent,omitempty"`
	DateTimeCreated time.Time `xml:"DateTimeCreated,omitempty" json:"DateTimeCreated,omitempty"`
	DisplayCc       string    `xml:"DisplayCc,omitempty" json:"DisplayCc,omitempty"`
	DisplayTo       string    `xml:"DisplayTo,omitempty" json:"DisplayTo,omitempty"`
	DisplayBcc      string    `xml:"DisplayBcc,omitempty" json:"DisplayBcc,omitempty"`
	HasAttachments  bool      `xml:"HasAttachments,omitempty" json:"HasAttachments,omitempty"`
}

type FolderInfoResponseMessage struct {
	*ResponseMessage
	Folders *ArrayOfFolders `xml:"Folders,omitempty" json:"Folders,omitempty"`
}

type ResponseMessage struct {
	MessageText string `xml:"MessageText,omitempty" json:"MessageText,omitempty"`
	// ResponseCode *ResponseCodeType `xml:"ResponseCode,omitempty" json:"ResponseCode,omitempty"`
	// DescriptiveLinkKey int32 `xml:"DescriptiveLinkKey,omitempty" json:"DescriptiveLinkKey,omitempty"`
	MessageXml struct {
	} `xml:"MessageXml,omitempty" json:"MessageXml,omitempty"`
	// ResponseClass *ResponseClassType `xml:"ResponseClass,attr,omitempty" json:"ResponseClass,omitempty"`
}

type ArrayOfFolders struct {
	Folder []*Folder `xml:"Folder,omitempty" json:"Folder,omitempty"`
	// CalendarFolder *CalendarFolderType `xml:"CalendarFolder,omitempty" json:"CalendarFolder,omitempty"`
	// ContactsFolder *ContactsFolderType `xml:"ContactsFolder,omitempty" json:"ContactsFolder,omitempty"`
	// SearchFolder *SearchFolderType `xml:"SearchFolder,omitempty" json:"SearchFolder,omitempty"`
	// TasksFolder *TasksFolderType `xml:"TasksFolder,omitempty" json:"TasksFolder,omitempty"`
}
