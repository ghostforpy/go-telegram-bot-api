package tgbotapi

func NoCommand(update Update) bool {
	return update.Message == nil || !update.Message.IsCommand()
}

func Private(update Update) bool {
	if update.EffectiveChat() == nil {
		return false
	}
	return update.EffectiveChat().Type == "private"
}

func FromGroup(update Update) bool {
	if update.EffectiveChat() == nil {
		return false
	}
	chatType := update.EffectiveChat().Type
	return chatType == "group" || chatType == "supergroup"
}

func FromChannel(update Update) bool {
	if update.EffectiveChat() == nil {
		return false
	}
	return update.EffectiveChat().Type == "channel"
}

func TextMessage(update Update) bool {
	return update.Message.Text != ""
}

func HandleDocument(update Update) bool {
	return update.Message.Document != nil
}

func HandlePhoto(update Update) bool {
	return update.Message.Photo != nil
}

func HandleAudio(update Update) bool {
	return update.Message.Audio != nil
}

func HandleVideo(update Update) bool {
	return update.Message.Video != nil
}

func HandleVideoNote(update Update) bool {
	return update.Message.VideoNote != nil
}

func HandleSticker(update Update) bool {
	return update.Message.Sticker != nil
}

func HandleVoice(update Update) bool {
	return update.Message.Voice != nil
}

func HandleVenue(update Update) bool {
	return update.Message.Venue != nil
}

func HandleLocation(update Update) bool {
	return update.Message.Location != nil
}

func HandleContact(update Update) bool {
	return update.Message.Contact != nil
}

func HandleAnimation(update Update) bool {
	return update.Message.Animation != nil
}

func HandlePremiumAnimation(update Update) bool {
	return update.Message.PremiumAnimation != nil
}
