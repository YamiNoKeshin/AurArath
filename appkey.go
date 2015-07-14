package aurarath

type AppKey struct {
	ApplicationKeyName string

	Functions []Function

	Tags []string

	Key []string
}

type Function struct {
	Name string

	Input []Parameter

	Output []Parameter
}

type Parameter struct {
	Name string

	Type string
}

func AppKeyFromJson(JSON string) (appkey AppKey) {
	err := json.Unmarshal([]byte(JSON), &appkey)
	if err != nil {
		panic(fmt.Sprint("Insane Appkey", err))
	}
	return
}

func AppKeyFromYaml(YAML string) (appkey AppKey) {
	if err := yaml.Unmarshal([]byte(YAML), &appkey); err != nil {
		panic(fmt.Sprint("Insane Appkey", err))
	}
	return
}
