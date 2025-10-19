package core

const (
	AppName        = "jebi"
	AppVersion     = "0.1.0"
	KeyLengthBytes = 32

	KeyFileName        = "enc"
	SecretFileName     = "sec"
	ProjectConfigFile  = "pro"
	DefaultEnvironment = "dev"
	CommitFileName     = "commits"
	CurrentFileName    = "current"
)

const (
	KdfAlgo      = "argon2id"
	CipherAlgo   = "aes-gcm"
	SaltLen      = 16 // bytes
	KeyLen       = 32 // 256-bit AES key
	NonceLen     = 12 // AES-GCM nonce size
	ArgonTime    = 3
	ArgonMemory  = 64 * 1024 // 64 MB
	ArgonThreads = 4
)
