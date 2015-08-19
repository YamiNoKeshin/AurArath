package connection

import "github.com/joernweissenborn/aurarath/config"

type AurArath struct {
	config *config.Config
}

func New(config *config.Config) (aurarath *AurArath) {
	aurarath.config = config
	return
}

func (*AurArath) AddImport(appkey *AppKey) (imp *Import){
	return nil
}

func (*AurArath) AddExport(appkey *AppKey) (exp *Export){
	return nil
}
