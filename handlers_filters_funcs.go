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
