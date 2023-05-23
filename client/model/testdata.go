package model

var (
	TestCard = Card{
		Info: Info{
			TypeID:      SecretTypes["CARD"],
			Title:       "Tinkoff Bank",
			Description: "Tinka",
		},
		CardholderName:  "PETR ",
		ExpirationMonth: 2,
		ExpirationYear:  25,
		PAN:             "111 2222 5525 4544",
		ServiceCode:     546,
	}
	TestAuth = Auth{
		Info: Info{
			TypeID:      SecretTypes["AUTH"],
			Title:       "Login to kk.com",
			Description: "Just login",
		},
		Password: "passw",
		Login:    "login",
	}
	TestText = Text{
		Info: Info{
			TypeID:      SecretTypes["TEXT"],
			Title:       "chapter1",
			Description: "full text of chapter one",
		},
		Text: "Big long text",
	}
)
