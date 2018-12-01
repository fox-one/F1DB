package storage

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/fox-one/broker"
	"github.com/fox-one/f1db/account"
	config "github.com/fox-one/f1db/config"
	fxAccount "github.com/fox-one/foxgo/account"
	fxWallet "github.com/fox-one/foxgo/wallet"
	uuid "github.com/satori/go.uuid"
	"github.com/vmihailenco/msgpack"
)

const ContentTypePlainText = "text/plain"
const ContentTypeImagePng = "image/png"

type ItemMeta struct {
	Creator     account.User
	Version     uint
	Type        string
	SnapshotID  string
	QuotaID     string
	QuotaAmount string
	Brief       string
	CID         string
}
type Item struct {
	ItemMeta
	Content string
}
type ItemHead struct {
	C string
	B string
}

func WriteItem(user *account.User, itemType string, brief string, content string) (*Item, error) {
	var err error
	var packed []byte
	var cid string
	meta := ItemMeta{
		Creator: *user,
		Version: 1,
		Type:    itemType,
		Brief:   brief,
		QuotaID: config.GetConfig().General.QuotaID,
		CID:     "",
	}
	item := Item{
		ItemMeta: meta,
		Content:  content,
	}

	if packed, err = msgpack.Marshal(&item); err != nil {
		return nil, err
	}
	if cid, err = WriteToIpfs(packed, false); err != nil {
		return nil, err
	}
	item.CID = cid
	return &item, err
}

// ReadItem read a specified content from DAG Network
func ReadItem(ctx context.Context, cid string) (*Item, error) {
	var err error
	var decoded []byte
	item := Item{}
	var contentItem *Item
	// Read Snapshot
	head := ItemHead{}
	var snapshot *broker.Snapshot
	snapshot, err = account.GetBroker().GetSnapshot(ctx, cid)
	if err != nil {
		return nil, err
	}
	if decoded, err = base64.StdEncoding.DecodeString(snapshot.Memo); err != nil {
		fmt.Printf("Error: %s\n", err)
		return nil, err
	}
	if err = msgpack.Unmarshal(decoded, &head); err != nil {
		fmt.Printf("Error: %s\n", err)
		return nil, err
	}
	// Read Record
	contentItem, err = ReadRecord(head.C)
	if err != nil {
		return nil, err
	}

	item.SnapshotID = snapshot.SnapshotId
	item.Brief = head.B
	item.CID = head.C
	item.QuotaID = snapshot.Asset.AssetId
	item.QuotaAmount = snapshot.Amount
	item.Content = contentItem.Content
	item.Type = contentItem.Type
	item.Creator.ID = contentItem.Creator.ID
	item.Version = contentItem.Version
	return &item, err
}

func KeepItem(ctx context.Context, item Item, token string, pin fxAccount.Pin, quota string) (*Item, error) {
	var err error
	var uid uuid.UUID
	var packed []byte
	payload := fxWallet.TransferRequest{}
	uid, err = uuid.NewV4()
	if err != nil {
		return nil, err
	}
	head := ItemHead{
		C: item.CID,
		B: item.Brief,
	}
	if packed, err = msgpack.Marshal(&head); err != nil {
		return nil, err
	}
	encoded := base64.StdEncoding.EncodeToString(packed)

	traceID := uid.String()
	payload.CounterUserId = config.GetConfig().General.CollectorUserID
	payload.TraceId = traceID
	payload.Amount = quota
	payload.AssetId = item.QuotaID
	payload.Memo = encoded
	snapshot, err := fxWallet.Transfer(ctx, token, pin, payload)
	if err != nil {
		return nil, err
	}
	item.SnapshotID = snapshot.SnapshotId
	item.QuotaID = snapshot.AssetId
	item.QuotaAmount = fmt.Sprintf("%f", snapshot.Amount)
	return &item, nil
}

// ReadRecord read a specified content from storage
func ReadRecord(cid string) (*Item, error) {
	var err error
	var packed []byte
	item := Item{}
	if packed, err = ReadFromIpfs(cid); err != nil {
		return nil, err
	}
	if err = msgpack.Unmarshal(packed, &item); err != nil {
		return nil, err
	}
	item.CID = cid
	return &item, err
}

func (item *Item) Response() interface{} {
	return map[string]interface{}{
		"user_id":      item.Creator.ID,
		"version":      item.Version,
		"type":         item.Type,
		"snapshot_id":  item.SnapshotID,
		"quota_id":     item.QuotaID,
		"quota_amount": item.QuotaAmount,
		"brief":        item.Brief,
		"cid":          item.CID,
		"content":      item.Content,
		"content_url":  fmt.Sprintf("https://ipfs.io/ipfs/%s", item.CID),
	}
}
