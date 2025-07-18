package whatsapp

import library "github.com/nocodeleaks/quepasa/library"

//
// Methods to create messages to send
//

// Default method to generate message
func ToMessage(
	recipient string,
	text string,
	trackid string,
) (msg *WhatsappMessage, err error) {
	chatId, err := FormatEndpoint(recipient)
	if err != nil {
		return
	}

	chat := WhatsappChat{Id: chatId}

	// Try to populate phone field
	if phone, _ := library.ExtractPhoneIfValid(chatId); len(phone) > 0 {
		chat.Phone = phone
	}

	msg = &WhatsappMessage{
		TrackId:      trackid,
		Text:         text,
		Chat:         chat,
		FromMe:       true,
		FromInternal: true,
	}

	// setting default type
	if len(text) > 0 {
		msg.Type = TextMessageType
	}

	return
}

// (Extension) Send Text Message
func ToMessageText(recipient string, text string) (msg *WhatsappMessage, err error) {
	return ToMessageTextWTrack(recipient, text, "")
}

// (Extension) Send Text Message with optional track id
func ToMessageTextWTrack(recipient string, text string, trackid string) (msg *WhatsappMessage, err error) {
	return ToMessage(recipient, text, trackid)
}
