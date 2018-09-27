package Errors

var errText = map[int]string{
	UserNameAlreadyExists: "User with same name already exists. Please choose anther one",
	UserEmailIsUsed:       "Email is used by another user",
	InvalidUsername:       "User is not found",
	InvalidPassword:       "Incorrect password",
	FailedToMarshalStruct: "Struct Marshaling failed",
	InvalidMail:           "Mail is invalid.",
	InvalidToken:          "Specified token is invalid",
}

func ErrorText(code int) string {
	return errText[code]
}
