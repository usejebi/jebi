package core

const (
	AppName        = "jebi"
	AppVersion     = "0.1.0"
	KeyLengthBytes = 32

	KeyFilePath       = "keys/enc.key"
	EnvDirPath        = "envs"
	SecretFileName    = "sec"
	ProjectConfigFile = "pro"
	CommitFileName    = "commits"
	CurrentFileName   = "current"

	DefaultEnvironment = "dev"
	DefaultProjectName = "my-jebi-project"
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
