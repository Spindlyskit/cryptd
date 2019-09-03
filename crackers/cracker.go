package crackers

// ReplyData is used to send the result of decrypt via rpc
type ReplyData struct {
	Key       string
	PlainText []byte
}

// SolveOptions defines options for solving a cipher
type SolveOptions struct {
	CT  []byte
	Key string
}

// Cracker defines methods used to decipher ciphertext
type Cracker interface {
	Crack(ct []byte) ([]byte, string, error)
	Decrypt(options SolveOptions, reply *ReplyData) error
	Solve(ct []byte, key string) ([]byte, error)
}
