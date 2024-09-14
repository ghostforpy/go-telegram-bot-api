package tgbotapi

type InlineKeyboard struct {
	HeaderBtns []InlineKeyboardButton
	MainBtns   []InlineKeyboardButton
	FooterBtns []InlineKeyboardButton
	Columns    int
}

func (kb *InlineKeyboard) Make() InlineKeyboardMarkup {
	var keyboard [][]InlineKeyboardButton
	if len(kb.HeaderBtns) > 0 {
		keyboard = append(keyboard, kb.HeaderBtns)
	}
	if kb.Columns == 0 {
		kb.Columns = 1
	}
	if len(kb.MainBtns) > 0 {
		var j int
		for i := 0; i < len(kb.MainBtns); i += kb.Columns {
			j += kb.Columns
			if j > len(kb.MainBtns) {
				j = len(kb.MainBtns)
			}
			keyboard = append(keyboard, kb.MainBtns[i:j])
		}
	}
	if len(kb.FooterBtns) > 0 {
		keyboard = append(keyboard, kb.FooterBtns)

	}
	return InlineKeyboardMarkup{
		InlineKeyboard: keyboard,
	}
}

func MakeInlineKeyboard(headerBtns []InlineKeyboardButton, mainBtns []InlineKeyboardButton, footerBtns []InlineKeyboardButton, columns int) InlineKeyboardMarkup {
	kb := &InlineKeyboard{
		HeaderBtns: headerBtns,
		MainBtns:   mainBtns,
		FooterBtns: footerBtns,
		Columns:    columns,
	}
	return kb.Make()
}
