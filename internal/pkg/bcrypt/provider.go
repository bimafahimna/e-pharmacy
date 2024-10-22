package bcrypt

type Provider interface {
	Hash(password string) (string, error)
	CompareHashAndPassword(hash, password string) error
}
