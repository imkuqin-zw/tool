package jwt_auth

type Header struct {
	Typ string `json:"typ"` // the media type of this complete JWT
	Alg string `json:"alg"` // algorithm to produce the JWE Encrypted Key
}

type Payload struct {
	Iss string `json:"iss"` // the principal that issued the JWT
	Sub string `json:"sub"` // the principal that is the subject of the JWT
	Aud string `json:"aud"` // the recipients that the JWT is intended for
	Exp int64  `json:"exp"` // the expiration time on or after which the JWT MUST NOT be accepted for processing
	Nbf int64  `json:"nbf"` // the time before which the JWT MUST NOT be accepted for processing
	Iat int64  `json:"iat"` // the time at which the JWT was issued
	//Jti string  `json:"jti"` // a unique identifier for the JWT
}
