package whatsmeow

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"
	"unicode"

	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"

	"github.com/nocodeleaks/quepasa/library"
	whatsapp "github.com/nocodeleaks/quepasa/whatsapp"
	whatsmeow "go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	types "go.mau.fi/whatsmeow/types"
)

// Must Implement IWhatsappConnection
type WhatsmeowConnection struct {
	library.LogStruct // logging
	Client            *whatsmeow.Client
	Handlers          *WhatsmeowHandlers

	failedToken  bool
	paired       func(string)
	IsConnecting bool `json:"isconnecting"` // used to avoid multiple connection attempts
}

//#region IMPLEMENT WHATSAPP CONNECTION OPTIONS INTERFACE

func (conn *WhatsmeowConnection) GetWid() string {
	if conn != nil {
		wid, err := conn.GetWidInternal()
		if err != nil {
			return wid
		}
	}

	return ""
}

// get default log entry, never nil
func (source *WhatsmeowConnection) GetLogger() *log.Entry {
	if source != nil && source.LogEntry != nil {
		return source.LogEntry
	}

	logentry := library.NewLogEntry(source)
	if source != nil {
		wid, _ := source.GetWidInternal()
		if len(wid) > 0 {
			logentry = logentry.WithField(LogFields.WId, wid)
		}
		source.LogEntry = logentry
	}

	logentry.Level = log.ErrorLevel
	logentry.Infof("generating new log entry for %s, with level: %s", reflect.TypeOf(source), logentry.Level)
	return logentry
}

func (conn *WhatsmeowConnection) SetReconnect(value bool) {
	if conn != nil {
		if conn.Client != nil {
			conn.Client.EnableAutoReconnect = value
		}
	}
}

func (conn *WhatsmeowConnection) GetReconnect() bool {
	if conn != nil {
		if conn.Client != nil {
			return conn.Client.EnableAutoReconnect
		}
	}

	return false
}

//#endregion

//region IMPLEMENT INTERFACE WHATSAPP CONNECTION

func (conn *WhatsmeowConnection) GetVersion() string { return "multi" }

func (conn *WhatsmeowConnection) GetWidInternal() (string, error) {
	if conn.Client == nil {
		err := fmt.Errorf("client not defined on trying to get wid")
		return "", err
	}

	if conn.Client.Store == nil {
		err := fmt.Errorf("device store not defined on trying to get wid")
		return "", err
	}

	if conn.Client.Store.ID == nil {
		err := fmt.Errorf("device id not defined on trying to get wid")
		return "", err
	}

	wid := conn.Client.Store.ID.User
	return wid, nil
}

func (conn *WhatsmeowConnection) IsValid() bool {
	if conn != nil {
		if conn.Client != nil {
			if conn.Client.IsConnected() {
				if conn.Client.IsLoggedIn() {
					return true
				}
			}
		}
	}
	return false
}

func (source *WhatsmeowConnection) IsConnected() bool {
	if source != nil {

		// manual checks for avoid thread locking
		if source.IsConnecting {
			return false
		}

		if source.Client != nil {
			if source.Client.IsConnected() {
				return true
			}
		}
	}
	return false
}

func (source *WhatsmeowConnection) GetStatus() whatsapp.WhatsappConnectionState {
	if source != nil {
		if source.Client == nil {
			return whatsapp.UnVerified
		} else {

			// manual checks for avoid thread locking
			if source.IsConnecting {
				return whatsapp.Connecting
			}

			// this is connected method locks the socket thread, so, if its in connecting state, it will be blocked here
			if source.Client.IsConnected() {
				if source.Client.IsLoggedIn() {
					return whatsapp.Ready
				} else {
					return whatsapp.Connected
				}
			} else {
				if source.failedToken {
					return whatsapp.Failed
				} else {
					return whatsapp.Disconnected
				}
			}
		}
	} else {
		return whatsapp.UnPrepared
	}
}

// returns a valid chat title from local memory store
func (conn *WhatsmeowConnection) GetChatTitle(wid string) string {
	jid, err := types.ParseJID(wid)
	if err == nil {
		return GetChatTitle(conn.Client, jid)
	}

	return ""
}

// Connect to websocket only, dot not authenticate yet, errors come after
func (source *WhatsmeowConnection) Connect() (err error) {
	source.GetLogger().Info("starting whatsmeow connection")

	if source.IsConnecting {
		return
	}

	source.IsConnecting = true

	err = source.Client.Connect()
	source.IsConnecting = false

	if err != nil {
		source.failedToken = true
		return
	}

	// waits 2 seconds for loggedin
	// not required
	_ = source.Client.WaitForConnection(time.Millisecond * 2000)

	source.failedToken = false
	return
}

func (source *WhatsmeowConnection) GetContacts() (chats []whatsapp.WhatsappChat, err error) {
	if source.Client == nil {
		err = errors.New("invalid client")
		return chats, err
	}

	if source.Client.Store == nil {
		err = errors.New("invalid store")
		return chats, err
	}

	contacts, err := source.Client.Store.Contacts.GetAllContacts(context.TODO())
	if err != nil {
		return chats, err
	}

	for jid, info := range contacts {

		title := info.FullName
		if len(title) == 0 {
			title = info.BusinessName
			if len(title) == 0 {
				title = info.PushName
			}
		}

		chats = append(chats, whatsapp.WhatsappChat{
			Id:    jid.String(),
			Title: title,
		})
	}

	return chats, nil
}

// func (cli *Client) Download(msg DownloadableMessage) (data []byte, err error)
func (source *WhatsmeowConnection) DownloadData(imsg whatsapp.IWhatsappMessage) (data []byte, err error) {
	msg := imsg.GetSource()
	downloadable, ok := msg.(whatsmeow.DownloadableMessage)
	if ok {
		return source.Client.Download(context.TODO(), downloadable)
	}

	logentry := source.GetLogger()
	logentry.Debug("not downloadable type, trying default message")

	waMsg, ok := msg.(*waE2E.Message)
	if !ok {
		attach := imsg.GetAttachment()
		if attach != nil {
			data := attach.GetContent()
			if data != nil {
				return *data, err
			}
		}

		err = fmt.Errorf("parameter msg cannot be converted to an original message")
		return
	}
	switch {
	case waMsg.ImageMessage != nil:
		replaced := strings.Replace(*waMsg.ImageMessage.URL, "mmg", "web", 1)
		waMsg.ImageMessage.URL = &replaced
	case waMsg.VideoMessage != nil:
		replaced := strings.Replace(*waMsg.VideoMessage.URL, "mmg", "web", 1)
		waMsg.VideoMessage.URL = &replaced
	case waMsg.AudioMessage != nil:
		replaced := strings.Replace(*waMsg.AudioMessage.URL, "mmg", "web", 1)
		waMsg.AudioMessage.URL = &replaced
	case waMsg.DocumentMessage != nil:
		replaced := strings.Replace(*waMsg.DocumentMessage.URL, "mmg", "web", 1)
		waMsg.DocumentMessage.URL = &replaced
	case waMsg.StickerMessage != nil:
		replaced := strings.Replace(*waMsg.StickerMessage.URL, "mmg", "web", 1)
		waMsg.StickerMessage.URL = &replaced
	}

	data, err = source.Client.DownloadAny(context.TODO(), waMsg)
	if err != nil {
		if strings.Contains(err.Error(), whatsmeow.ErrFileLengthMismatch.Error()) {
			logentry.Infof("ignoring (%s) whatsmeow error for msg id: %s", whatsmeow.ErrFileLengthMismatch.Error(), imsg.GetId())
			err = nil
		}
	}

	return
}

func (conn *WhatsmeowConnection) Download(imsg whatsapp.IWhatsappMessage, cache bool) (att *whatsapp.WhatsappAttachment, err error) {
	att = imsg.GetAttachment()
	if att == nil {
		err = fmt.Errorf("message (%s) does not contains attachment info", imsg.GetId())
		return
	}

	if !att.HasContent() && !att.CanDownload {
		err = fmt.Errorf("message (%s) attachment with invalid content and not available to download", imsg.GetId())
		return
	}

	if !att.HasContent() || (att.CanDownload && !cache) {
		data, err := conn.DownloadData(imsg)
		if err != nil {
			return att, err
		}

		if !cache {
			newAtt := *att
			att = &newAtt
		}

		att.SetContent(&data)
	}

	return
}

func (source *WhatsmeowConnection) Revoke(msg whatsapp.IWhatsappMessage) error {
	logentry := source.GetLogger()

	jid, err := types.ParseJID(msg.GetChatId())
	if err != nil {
		logentry.Infof("revoke error on get jid: %s", err)
		return err
	}

	participantJid, err := types.ParseJID(msg.GetParticipantId())
	if err != nil {
		logentry.Infof("revoke error on get jid: %s", err)
		return err
	}

	newMessage := source.Client.BuildRevoke(jid, participantJid, msg.GetId())
	_, err = source.Client.SendMessage(context.Background(), jid, newMessage)
	if err != nil {
		logentry.Infof("revoke error: %s", err)
		return err
	}

	return nil
}

func (conn *WhatsmeowConnection) IsOnWhatsApp(phones ...string) (registered []string, err error) {
	results, err := conn.Client.IsOnWhatsApp(phones)
	if err != nil {
		return
	}

	for _, result := range results {
		if result.IsIn {
			registered = append(registered, result.JID.String())
		}
	}

	return
}

func (conn *WhatsmeowConnection) GetProfilePicture(wid string, knowingId string) (picture *whatsapp.WhatsappProfilePicture, err error) {
	jid, err := types.ParseJID(wid)
	if err != nil {
		return
	}

	params := &whatsmeow.GetProfilePictureParams{}
	params.ExistingID = knowingId
	params.Preview = false

	pictureInfo, err := conn.Client.GetProfilePictureInfo(jid, params)
	if err != nil {
		return
	}

	if pictureInfo != nil {
		picture = &whatsapp.WhatsappProfilePicture{
			Id:   pictureInfo.ID,
			Type: pictureInfo.Type,
			Url:  pictureInfo.URL,
		}
	}
	return
}

func isASCII(s string) bool {
	for _, c := range s {
		if c > unicode.MaxASCII {
			return false
		}
	}
	return true
}

func (source *WhatsmeowConnection) GetContextInfo(msg whatsapp.WhatsappMessage) *waE2E.ContextInfo {

	var contextInfo *waE2E.ContextInfo
	if len(msg.InReply) > 0 {
		contextInfo = source.GetInReplyContextInfo(msg)
	}

	// mentions ---------------------------------------
	if msg.FromGroup() {
		messageText := msg.GetText()
		mentions := GetMentions(messageText)
		if len(mentions) > 0 {

			if contextInfo == nil {
				contextInfo = &waE2E.ContextInfo{}
			}

			contextInfo.MentionedJID = mentions
		}
	}

	// disapering messages, not implemented yet
	if contextInfo != nil {
		contextInfo.Expiration = proto.Uint32(0)
		contextInfo.EphemeralSettingTimestamp = proto.Int64(0)
		contextInfo.DisappearingMode = &waE2E.DisappearingMode{Initiator: waE2E.DisappearingMode_CHANGED_IN_CHAT.Enum()}
	}

	return contextInfo
}

func (source *WhatsmeowConnection) GetInReplyContextInfo(msg whatsapp.WhatsappMessage) *waE2E.ContextInfo {
	logentry := source.GetLogger()

	// default information for cached messages
	var info types.MessageInfo

	// getting quoted message if available on cache
	// (optional) another devices will process anyway, but our devices will show quoted only if it exists on cache
	var quoted *waE2E.Message

	cached, _ := source.Handlers.WAHandlers.GetById(msg.InReply)
	if cached != nil {

		// update cached info
		info, _ = cached.InfoForHistory.(types.MessageInfo)

		if cached.Content != nil {
			if content, ok := cached.Content.(*waE2E.Message); ok {

				// update quoted message content
				quoted = content

			} else {
				logentry.Warnf("content has an invalid type (%s), on reply to msg id: %s", reflect.TypeOf(cached.Content), msg.InReply)
			}
		} else {
			logentry.Warnf("message content not cached, on reply to msg id: %s", msg.InReply)
		}
	} else {
		logentry.Warnf("message not cached, on reply to msg id: %s", msg.InReply)
	}

	var participant *string
	if (types.MessageInfo{}) != info {
		var sender string
		if msg.FromGroup() {
			sender = fmt.Sprint(info.Sender.User, "@", info.Sender.Server)
		} else {
			sender = fmt.Sprint(info.Chat.User, "@", info.Chat.Server)
		}
		participant = proto.String(sender)
	}

	return &waE2E.ContextInfo{
		StanzaID:      proto.String(msg.InReply),
		Participant:   participant,
		QuotedMessage: quoted,
	}
}

// Default SEND method using WhatsappMessage Interface
func (source *WhatsmeowConnection) Send(msg *whatsapp.WhatsappMessage) (whatsapp.IWhatsappSendResponse, error) {
	logentry := source.GetLogger()
	loglevel := logentry.Level
	logentry = logentry.WithField(LogFields.MessageId, msg.Id)
	logentry.Level = loglevel

	var err error

	// Formatting destination accordingly
	formattedDestination, _ := whatsapp.FormatEndpoint(msg.GetChatId())

	// avoid common issue with incorrect non ascii chat id
	if !isASCII(formattedDestination) {
		err = fmt.Errorf("not an ASCII formatted chat id")
		return msg, err
	}

	// validating jid before remote commands as upload or send
	jid, err := types.ParseJID(formattedDestination)
	if err != nil {
		logentry.Infof("send error on get jid: %s", err)
		return msg, err
	}

	// request message text
	messageText := msg.GetText()

	var newMessage *waE2E.Message
	if !msg.HasAttachment() {
		if IsValidForButtons(messageText) {
			internal := GenerateButtonsMessage(messageText)
			internal.ContextInfo = source.GetContextInfo(*msg)
			newMessage = &waE2E.Message{
				ButtonsMessage: internal,
			}
		} else {
			if msg.Poll != nil {
				newMessage, err = GeneratePollMessage(msg)
				if err != nil {
					return msg, err
				}
			} else {
				internal := &waE2E.ExtendedTextMessage{Text: &messageText}
				internal.ContextInfo = source.GetContextInfo(*msg)
				newMessage = &waE2E.Message{ExtendedTextMessage: internal}
			}
		}
	} else {
		newMessage, err = source.UploadAttachment(*msg)
		if err != nil {
			return msg, err
		}
	}

	// Generating a new unique MessageID
	if len(msg.Id) == 0 {
		msg.Id = source.Client.GenerateMessageID()
	}

	extra := whatsmeow.SendRequestExtra{
		ID: msg.Id,
	}

	// saving cached content for instance of future reply
	if msg.Content == nil {
		msg.Content = newMessage
	}

	resp, err := source.Client.SendMessage(context.Background(), jid, newMessage, extra)
	if err != nil {
		logentry.Errorf("whatsmeow connection send error: %s", err)
		return msg, err
	}

	// updating timestamp
	msg.Timestamp = resp.Timestamp

	if msg.Id != resp.ID {
		logentry.Warnf("send success but msg id differs from response id: %s, type: %v, on: %s", resp.ID, msg.Type, msg.Timestamp)
	} else {
		logentry.Infof("send success, type: %v, on: %s", msg.Type, msg.Timestamp)
	}

	// testing, mark read function
	if source.Handlers.ReadUpdate {
		go source.Handlers.MarkRead(msg, types.ReceiptTypeRead)
	}

	return msg, err
}

// useful to check if is a member of a group before send a msg.
// fails on recently added groups.
// pending a more efficient code !!!!!!!!!!!!!!
func (source *WhatsmeowConnection) HasChat(chat string) bool {
	jid, err := types.ParseJID(chat)
	if err != nil {
		return false
	}

	info, err := source.Client.Store.ChatSettings.GetChatSettings(context.TODO(), jid)
	if err != nil {
		return false
	}

	return info.Found
}

// func (cli *Client) Upload(ctx context.Context, plaintext []byte, appInfo MediaType) (resp UploadResponse, err error)
func (source *WhatsmeowConnection) UploadAttachment(msg whatsapp.WhatsappMessage) (result *waE2E.Message, err error) {

	content := *msg.Attachment.GetContent()
	if len(content) == 0 {
		err = fmt.Errorf("null or empty content")
		return
	}

	mediaType := GetMediaTypeFromWAMsgType(msg.Type)
	response, err := source.Client.Upload(context.Background(), content, mediaType)
	if err != nil {
		return
	}

	inreplycontext := source.GetInReplyContextInfo(msg)
	result = NewWhatsmeowMessageAttachment(response, msg, mediaType, inreplycontext)
	return
}

func (conn *WhatsmeowConnection) Disconnect() (err error) {
	if conn.Client != nil {
		if conn.Client.IsConnected() {
			conn.Client.Disconnect()
		}
	}
	return
}

func (source *WhatsmeowConnection) GetInvite(groupId string) (link string, err error) {
	jid, err := types.ParseJID(groupId)
	if err != nil {
		source.GetLogger().Infof("getting invite error on parse jid: %s", err)
	}

	link, err = source.Client.GetGroupInviteLink(jid, false)
	return
}

//region PAIRING

func (source *WhatsmeowConnection) PairPhone(phone string) (string, error) {

	if !source.Client.IsConnected() {
		err := source.Client.Connect()
		if err != nil {
			log.Errorf("error on connecting for getting whatsapp qrcode: %s", err.Error())
			return "", err
		}
	}

	return source.Client.PairPhone(context.TODO(), phone, true, whatsmeow.PairClientChrome, "Chrome (Linux)")
}

func (conn *WhatsmeowConnection) GetWhatsAppQRCode() string {

	var result string

	// No ID stored, new login
	qrChan, err := conn.Client.GetQRChannel(context.Background())
	if err != nil {
		log.Errorf("error on getting whatsapp qrcode channel: %s", err.Error())
		return ""
	}

	if !conn.Client.IsConnected() {
		err = conn.Client.Connect()
		if err != nil {
			log.Errorf("error on connecting for getting whatsapp qrcode: %s", err.Error())
			return ""
		}
	}

	var wg sync.WaitGroup
	wg.Add(1)

	for evt := range qrChan {
		if evt.Event == "code" {
			result = evt.Code
		}

		wg.Done()
		break
	}

	wg.Wait()
	return result
}

//endregion

func TryUpdateChannel(ch chan<- string, value string) (closed bool) {
	defer func() {
		if recover() != nil {
			// the return result can be altered
			// in a defer function call
			closed = false
		}
	}()

	ch <- value // panic if ch is closed
	return true // <=> closed = false; return
}

func (source *WhatsmeowConnection) GetWhatsAppQRChannel(ctx context.Context, out chan<- string) error {
	logger := source.GetLogger()

	// No ID stored, new login
	qrChan, err := source.Client.GetQRChannel(ctx)
	if err != nil {
		logger.Errorf("error on getting whatsapp qrcode channel: %s", err.Error())
		return err
	}

	if !source.Client.IsConnected() {
		err = source.Client.Connect()
		if err != nil {
			logger.Errorf("error on connecting for getting whatsapp qrcode: %s", err.Error())
			return err
		}
	}

	var wg sync.WaitGroup
	wg.Add(1)

	for evt := range qrChan {
		if evt.Event == "code" {
			if !TryUpdateChannel(out, evt.Code) {
				// expected error, means that websocket was closed
				// probably user has gone out page
				return fmt.Errorf("cant write to output")
			}
		} else {
			if evt.Event == "timeout" {
				return errors.New("timeout")
			}
			wg.Done()
			break
		}
	}

	wg.Wait()
	return nil
}

func (source *WhatsmeowConnection) HistorySync(timestamp time.Time) (err error) {
	logentry := source.GetLogger()

	leading := source.Handlers.WAHandlers.GetLeading()
	if leading == nil {
		err = fmt.Errorf("no valid msg in cache for retrieve parents")
		return err
	}

	// Convert interface to struct using type assertion
	info, ok := leading.InfoForHistory.(types.MessageInfo)
	if !ok {
		logentry.Error("error converting leading for history")
	}

	logentry.Infof("getting history from: %s", timestamp)
	extra := whatsmeow.SendRequestExtra{Peer: true}

	//info := &types.MessageInfo{ }
	msg := source.Client.BuildHistorySyncRequest(&info, 50)
	response, err := source.Client.SendMessage(context.Background(), source.Client.Store.ID.ToNonAD(), msg, extra)
	if err != nil {
		logentry.Errorf("getting history error: %s", err.Error())
	}

	logentry.Infof("history: %v", response)
	return
}

func (conn *WhatsmeowConnection) UpdateHandler(handlers whatsapp.IWhatsappHandlers) {
	conn.Handlers.WAHandlers = handlers
}

func (conn *WhatsmeowConnection) UpdatePairedCallBack(callback func(string)) {
	conn.paired = callback
}

func (conn *WhatsmeowConnection) PairedCallBack(jid types.JID, platform, businessName string) bool {
	if conn.paired != nil {
		go conn.paired(jid.String())
	}
	return true
}

//endregion

/*
<summary>

	Disconnect if connected
	Cleanup Handlers
	Dispose resources
	Does not erase permanent data !

</summary>
*/
func (source *WhatsmeowConnection) Dispose(reason string) {

	source.GetLogger().Infof("disposing connection: %s", reason)

	if source.Handlers != nil {
		go source.Handlers.UnRegister()
		source.Handlers = nil
	}

	if source.Client != nil {
		if source.Client.IsConnected() {
			go source.Client.Disconnect()
		}
		source.Client = nil
	}

	source = nil
}

/*
<summary>

	Erase permanent data + Dispose !

</summary>
*/
func (source *WhatsmeowConnection) Delete() (err error) {
	if source != nil {
		if source.Client != nil {
			if source.Client.IsLoggedIn() {
				err = source.Client.Logout(context.TODO())
				if err != nil {
					err = fmt.Errorf("whatsmeow connection, delete logout error: %s", err.Error())
					return
				}
				source.GetLogger().Infof("logged out for delete")
			}

			if source.Client.Store != nil {
				err = source.Client.Store.Delete(context.TODO())
				if err != nil {
					// ignoring error about JID, just checked and the delete process was succeed
					if strings.Contains(err.Error(), "device JID must be known before accessing database") {
						err = nil
					} else {
						err = fmt.Errorf("whatsmeow connection, delete store error: %s", err.Error())
						return
					}
				}
			}
		}
	}

	source.Dispose("Delete")
	return
}

func (conn *WhatsmeowConnection) IsInterfaceNil() bool {
	return nil == conn
}
func (conn *WhatsmeowConnection) GetJoinedGroups() ([]interface{}, error) {
	if conn.Client == nil {
		return nil, fmt.Errorf("client not defined")
	}

	// Get the group info slice
	groupInfos, err := conn.Client.GetJoinedGroups()
	if err != nil {
		return nil, err
	}

	// Convert []*types.GroupInfo to []interface{}
	groups := make([]interface{}, len(groupInfos))
	for i, group := range groupInfos {
		groups[i] = group
	}

	return groups, nil
}

func (conn *WhatsmeowConnection) GetGroupInfo(groupId string) (interface{}, error) {
	if conn.Client == nil {
		return nil, fmt.Errorf("client not defined")
	}

	jid, err := types.ParseJID(groupId)
	if err != nil {
		return nil, err
	}

	return conn.Client.GetGroupInfo(jid)
}

func (conn *WhatsmeowConnection) CreateGroup(name string, participants []string) (interface{}, error) {
	if conn.Client == nil {
		return nil, fmt.Errorf("client not defined")
	}

	// Convert participants to JID format
	var participantsJID []types.JID
	for _, participant := range participants {
		// Check if it's already in JID format
		if strings.Contains(participant, "@") {
			jid, err := types.ParseJID(participant)
			if err != nil {
				return nil, fmt.Errorf("invalid JID format for participant %s: %v", participant, err)
			}
			participantsJID = append(participantsJID, jid)
		} else {
			// Assume it's a phone number and convert to JID
			jid := types.JID{
				User:   participant,
				Server: "s.whatsapp.net", // Use the standard WhatsApp server
			}
			participantsJID = append(participantsJID, jid)
		}
	}

	// Create the request struct
	groupConfig := whatsmeow.ReqCreateGroup{
		Name:         name,
		Participants: participantsJID,
	}

	// Call the existing method with the constructed request
	return conn.Client.CreateGroup(groupConfig)
}

func (conn *WhatsmeowConnection) UpdateGroupSubject(groupID string, name string) (interface{}, error) {
	if conn.Client == nil {
		return nil, fmt.Errorf("client not defined")
	}

	// Parse the group ID to JID format
	jid, err := types.ParseJID(groupID)
	if err != nil {
		return nil, fmt.Errorf("invalid group JID format: %v", err)
	}

	// Update the group subject
	err = conn.Client.SetGroupName(jid, name)
	if err != nil {
		return nil, fmt.Errorf("failed to update group subject: %v", err)
	}

	// Return the updated group info
	return conn.Client.GetGroupInfo(jid)
}

func (conn *WhatsmeowConnection) UpdateGroupPhoto(groupID string, imageData []byte) (string, error) {
	if conn.Client == nil {
		return "", fmt.Errorf("client not defined")
	}

	// Parse the group ID to JID format
	jid, err := types.ParseJID(groupID)
	if err != nil {
		return "", fmt.Errorf("invalid group JID format: %v", err)
	}

	// Update the group photo
	pictureID, err := conn.Client.SetGroupPhoto(jid, imageData)
	if err != nil {
		return "", fmt.Errorf("failed to update group photo: %v", err)
	}

	return pictureID, nil
}

func (conn *WhatsmeowConnection) UpdateGroupParticipants(groupJID string, participants []string, action string) ([]interface{}, error) {
	if conn.Client == nil {
		return nil, fmt.Errorf("client not defined")
	}

	// Parse the group JID
	jid, err := types.ParseJID(groupJID)
	if err != nil {
		return nil, fmt.Errorf("invalid group JID format: %v", err)
	}

	// Convert participant strings to JIDs
	participantJIDs := make([]types.JID, len(participants))
	for i, participant := range participants {
		participantJIDs[i], err = types.ParseJID(participant)
		if err != nil {
			return nil, fmt.Errorf("invalid participant JID format for %s: %v", participant, err)
		}
	}

	// Map the action string to the ParticipantChange type
	var participantAction whatsmeow.ParticipantChange
	switch action {
	case "add":
		participantAction = whatsmeow.ParticipantChangeAdd
	case "remove":
		participantAction = whatsmeow.ParticipantChangeRemove
	case "promote":
		participantAction = whatsmeow.ParticipantChangePromote
	case "demote":
		participantAction = whatsmeow.ParticipantChangeDemote
	default:
		return nil, fmt.Errorf("invalid action %s", action)
	}

	// Call the whatsmeow method
	result, err := conn.Client.UpdateGroupParticipants(jid, participantJIDs, participantAction)
	if err != nil {
		return nil, fmt.Errorf("failed to update group participants: %v", err)
	}

	// Convert to interface array for the generic return type
	interfaceResults := make([]interface{}, len(result))
	for i, r := range result {
		interfaceResults[i] = r
	}

	return interfaceResults, nil
}

func (conn *WhatsmeowConnection) GetGroupJoinRequests(groupJID string) ([]interface{}, error) {
	if conn.Client == nil {
		return nil, fmt.Errorf("client not defined")
	}

	// Parse the group JID
	jid, err := types.ParseJID(groupJID)
	if err != nil {
		return nil, fmt.Errorf("invalid group JID format: %v", err)
	}

	// Call the whatsmeow method
	requests, err := conn.Client.GetGroupRequestParticipants(jid)
	if err != nil {
		return nil, fmt.Errorf("failed to get group join requests: %v", err)
	}

	// Convert to interface array for the generic return type
	interfaceResults := make([]interface{}, len(requests))
	for i, r := range requests {
		interfaceResults[i] = r
	}

	return interfaceResults, nil
}

func (conn *WhatsmeowConnection) HandleGroupJoinRequests(groupJID string, participants []string, action string) ([]interface{}, error) {
	if conn.Client == nil {
		return nil, fmt.Errorf("client not defined")
	}

	// Parse the group JID
	jid, err := types.ParseJID(groupJID)
	if err != nil {
		return nil, fmt.Errorf("invalid group JID format: %v", err)
	}

	// Convert participant strings to JIDs
	participantJIDs := make([]types.JID, len(participants))
	for i, participant := range participants {
		participantJIDs[i], err = types.ParseJID(participant)
		if err != nil {
			return nil, fmt.Errorf("invalid participant JID format for %s: %v", participant, err)
		}
	}

	// Map the action string to the ParticipantRequestChange type
	var requestAction whatsmeow.ParticipantRequestChange
	switch action {
	case "approve":
		requestAction = whatsmeow.ParticipantChangeApprove
	case "reject":
		requestAction = whatsmeow.ParticipantChangeReject
	default:
		return nil, fmt.Errorf("invalid action %s", action)
	}

	// Call the correct WhatsApp method which returns participant results
	result, err := conn.Client.UpdateGroupRequestParticipants(jid, participantJIDs, requestAction)
	if err != nil {
		return nil, fmt.Errorf("failed to handle group join requests: %v", err)
	}

	// Convert the typed results to interface array
	interfaceResults := make([]interface{}, len(result))
	for i, r := range result {
		interfaceResults[i] = r
	}

	return interfaceResults, nil
}

func (conn *WhatsmeowConnection) CreateGroupExtended(title string, participants []string) (interface{}, error) {
	if conn.Client == nil {
		return nil, fmt.Errorf("client not defined")
	}

	// Convert participants to JIDs
	participantJIDs := make([]types.JID, len(participants))
	for i, participant := range participants {
		jid, err := types.ParseJID(participant)
		if err != nil {
			return nil, fmt.Errorf("invalid participant JID: %v", err)
		}
		participantJIDs[i] = jid
	}

	// Create request structure
	req := whatsmeow.ReqCreateGroup{
		Name:         title,
		Participants: participantJIDs,
	}

	// Call the WhatsApp method
	return conn.Client.CreateGroup(req)
}

// SendChatPresence updates typing status in a chat
func (conn *WhatsmeowConnection) SendChatPresence(chatId string, presenceType uint) error {
	if conn.Client == nil {
		return fmt.Errorf("client not defined")
	}

	jid, err := types.ParseJID(chatId)
	if err != nil {
		return fmt.Errorf("invalid chat id format: %v", err)
	}

	var state types.ChatPresence
	var media types.ChatPresenceMedia

	// Map our custom presence type to whatsmeow types
	switch whatsapp.WhatsappChatPresenceType(presenceType) {
	case whatsapp.WhatsappChatPresenceTypeText:
		state = types.ChatPresenceComposing // typing
		media = types.ChatPresenceMediaText
	case whatsapp.WhatsappChatPresenceTypeAudio:
		state = types.ChatPresenceComposing // typing
		media = types.ChatPresenceMediaAudio
	default:
		// Default is paused (stop typing)
		state = types.ChatPresencePaused
		media = types.ChatPresenceMediaText
	}

	return conn.Client.SendChatPresence(jid, state, media)
}
