package jwt

type Provider interface {
	Sign(userID int64, role string, isVerified bool) (string, error)
	Parse(tokenString string) (*jwtClaims, error)
}
