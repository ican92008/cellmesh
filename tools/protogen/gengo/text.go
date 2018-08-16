package gengo

// 报错行号+3
const goCodeTemplate = `// Auto generated by github.com/davyxu/cellmesh/protogen
// DO NOT EDIT!

package {{.PackageName}}

import (
	"fmt"
	"reflect"	
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/codec"
	"github.com/davyxu/cellmesh/service"
	{{if HasJsonCodec $}}_ "github.com/davyxu/cellnet/codec/json"{{end}}
)

// Make compiler import happy
var (
	_ service.Service
	_ cellnet.Peer
	_ cellnet.Codec
	_ reflect.Type
	_ fmt.Formatter
)


{{range $a, $enumobj := .Enums}}
type {{.Name}} int32
const (	{{range .Fields}}
	{{$enumobj.Name}}_{{.Name}} {{$enumobj.Name}} = {{TagNumber $enumobj .}} {{end}}
)

var (
{{$enumobj.Name}}MapperValueByName = map[string]int32{ {{range .Fields}}
	"{{.Name}}": {{TagNumber $enumobj .}}, {{end}}
}

{{$enumobj.Name}}MapperNameByValue = map[int32]string{ {{range .Fields}}
	{{TagNumber $enumobj .}}: "{{.Name}}" , {{end}}
}

{{$enumobj.Name}}MapperTrailingCommentByValue = map[int32]string{ {{range .Fields}}
	{{TagNumber $enumobj .}}: "{{.Trailing}}" , {{end}}
}
)

func (self {{$enumobj.Name}}) String() string {
	return {{$enumobj.Name}}MapperNameByValue[int32(self)]
}
{{end}}

{{range .Structs}}
{{ObjectLeadingComment .}}
type {{.Name}} struct{	{{range .Fields}}
	{{GoFieldName .}} {{GoTypeName .}} {{GoStructTag .}}{{FieldTrailingComment .}} {{end}}
}
{{end}}
{{range .Structs}}
func (self *{{.Name}}) String() string { return fmt.Sprintf("%+v",*self) } {{end}}

func GetRPCPair(req interface{}) reflect.Type {

	switch req.(type) { {{range RPCPair $}}
	case *{{.REQ.Name}}:
		return reflect.TypeOf((*{{.ACK.Name}})(nil)).Elem() {{end}}
	}

	return nil
}

{{range ServiceGroup $}}
// {{.Key}}
var ( {{range .Group}}
	Handler_{{.Name}} = func(ev service.Event, req *{{.Name}}){ panic("'{{.Name}}' not handled") } {{end}}
)
{{end}}

func GetDispatcher(svcName string) service.DispatcherFunc {

	switch svcName { {{range ServiceGroup $}}
	case "{{.Key}}":
		return func(ev service.Event) {
			switch req := ev.Message().(type) { {{range .Group}}
			case *{{.Name}}:
				Handler_{{.Name}}(ev, req) {{end}}
			}
		} {{end}}
	} 

	return nil
}

func init() {
	{{range .Structs}} {{ if IsMessage . }}
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("{{StructCodec .}}"),
		Type:  reflect.TypeOf((*{{.Name}})(nil)).Elem(),
		ID:    {{StructMsgID .}},
	}).SetContext("service", "{{StructService .}}")
	{{end}} {{end}}
}

`