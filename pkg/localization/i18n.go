package localization

import (
	"strings"
)


const (
	LangEN      = "en"
	LangKA      = "ka"
	LangDefault = LangEN
)


var Messages = map[string]map[string]string{
	LangEN: {
		"user.created":                 "User created successfully",
		"user.login.success":           "Login successful",
		"user.invalid_credentials":     "Invalid email or password",
		"user.already_exists":          "User with this email already exists",
		"team.created":                 "Team created successfully",
		"team.updated":                 "Team updated successfully",
		"team.not_found":               "Team not found",
		"player.updated":               "Player updated successfully",
		"player.not_found":             "Player not found",
		"player.not_owned":             "Player does not belong to your team",
		"player.listed":                "Player listed for transfer",
		"player.already_listed":        "Player is already on transfer list",
		"player.removed_from_list":     "Player removed from transfer list",
		"player.not_on_list":           "Player is not on transfer list",
		"transfer.purchased":           "Player purchased successfully",
		"transfer.insufficient_budget": "Insufficient budget",
		"transfer.team_full":           "Team already has maximum number of players",
		"transfer.cannot_buy_own":      "Cannot buy your own player",
		"transfer.listing_not_found":   "Transfer listing not found",
		"error.internal":               "Internal server error",
		"error.validation":             "Validation error",
		"error.unauthorized":           "Unauthorized",
	},
	LangKA: {
		"user.created":                 "მომხმარებელი წარმატებით შეიქმნა",
		"user.login.success":           "შესვლა წარმატებულია",
		"user.invalid_credentials":     "არასწორი ელფოსტა ან პაროლი",
		"user.already_exists":          "ამ ელფოსტით მომხმარებელი უკვე არსებობს",
		"team.created":                 "გუნდი წარმატებით შეიქმნა",
		"team.updated":                 "გუნდი განახლდა",
		"team.not_found":               "გუნდი ვერ მოიძებნა",
		"player.updated":               "მოთამაშე განახლდა",
		"player.not_found":             "მოთამაშე ვერ მოიძებნა",
		"player.not_owned":             "მოთამაშე არ ეკუთვნის თქვენს გუნდს",
		"player.listed":                "მოთამაშე განთავსდა გადაცემის სიაში",
		"player.already_listed":        "მოთამაშე უკვე არის გადაცემის სიაში",
		"player.removed_from_list":     "მოთამაშე წაიშალა გადაცემის სიიდან",
		"player.not_on_list":           "მოთამაშე არ არის გადაცემის სიაში",
		"transfer.purchased":           "მოთამაშე წარმატებით შეიძინა",
		"transfer.insufficient_budget": "არასაკმარისი ბიუჯეტი",
		"transfer.team_full":           "გუნდს უკვე აქვს მაქსიმალური რაოდენობის მოთამაშე",
		"transfer.cannot_buy_own":      "ვერ შეიძენთ საკუთარ მოთამაშეს",
		"transfer.listing_not_found":   "გადაცემის სია ვერ მოიძებნა",
		"error.internal":               "შიდა სერვერის შეცდომა",
		"error.validation":             "ვალიდაციის შეცდომა",
		"error.unauthorized":           "არაავტორიზებული",
	},
}


func GetMessage(lang, key string) string {
	lang = normalizeLang(lang)
	if messages, ok := Messages[lang]; ok {
		if msg, ok := messages[key]; ok {
			return msg
		}
	}

	if messages, ok := Messages[LangDefault]; ok {
		if msg, ok := messages[key]; ok {
			return msg
		}
	}
	return key
}


func normalizeLang(lang string) string {
	lang = strings.ToLower(strings.TrimSpace(lang))
	if lang == "" {
		return LangDefault
	}

	if len(lang) > 2 {
		lang = lang[:2]
	}
	if lang != LangEN && lang != LangKA {
		return LangDefault
	}
	return lang
}


func GetLanguageFromHeader(acceptLang string) string {
	if acceptLang == "" {
		return LangDefault
	}

	langs := strings.Split(acceptLang, ",")
	if len(langs) > 0 {
		lang := strings.Split(strings.TrimSpace(langs[0]), ";")[0]
		return normalizeLang(lang)
	}
	return LangDefault
}
