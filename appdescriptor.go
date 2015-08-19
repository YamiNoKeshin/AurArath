package connection

import (
	"encoding/json"
	"fmt"
)

type AppDescriptor struct {
	Functions          []Function
	Tags               map[string]string
}

type Function struct {
	Name   string
	Input  []Parameter
	Output []Parameter
}

func (f Function) String(){
	inpar := ""
	for _,p := range f.Input {
		inpar = fmt.Sprintf("%s,%s",inpar)
	}
	outpar := ""
	for _,p := range f.Output{
		outpar = fmt.Sprintf("%s,%s",outpar)
	}
	return fmt.Sprintf("%s(%s)%s",f.Name,inpar,outpar)
}

type Parameter struct {
	Name string
	Type string
}

//name(par:type,...)par:type

func (p Parameter) String(){
	return fmt.Sprintf("%s:%s",p.Name,p.Type)
}

func AppDescriptorFromJson(JSON string) (appkey AppDescriptor) {
	err := json.Unmarshal([]byte(JSON), &appkey)
	if err != nil {
		panic(fmt.Sprint("Insane Appkey", err))
	}
	return
}

func (a AppDescriptor) AsTagSet() (tagset map[string]string){
	tagset = a.Tags

	for _, fn := range a.Functions {
		tagset[fmt.Sprintf("function_%s",fn)] = ""
	}
	return
}

/*
func AppKeyFromYaml(YAML string) (appkey AppKey) {
	if err := yaml.Unmarshal([]byte(YAML), &appkey); err != nil {
		panic(fmt.Sprint("Insane Appkey", err))
	}
	return
}
*/
