package aurarath

type AurArath struct {
}

func New(cfg *Config) (aurarath *AurArath)

func (*AurArath) AddImport(appkey *AppKey) (imp *Import)

func (*AurArath) AddExport(appkey *AppKey) (exp *Export)
