package aurarath

type AurArath struct {
	config *Config
}

func New(config *Config) (aurarath *AurArath) {
	aurarath.config = config
	return
}

func (*AurArath) AddImport(appkey *AppKey) (imp *Import)

func (*AurArath) AddExport(appkey *AppKey) (exp *Export)
