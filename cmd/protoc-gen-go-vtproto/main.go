package main

import (
	"flag"
	"fmt"
	"strings"

	_ "github.com/davidflanagan/vtprotobuf/features/equal"
	_ "github.com/davidflanagan/vtprotobuf/features/grpc"
	_ "github.com/davidflanagan/vtprotobuf/features/marshal"
	_ "github.com/davidflanagan/vtprotobuf/features/pool"
	_ "github.com/davidflanagan/vtprotobuf/features/size"
	_ "github.com/davidflanagan/vtprotobuf/features/unmarshal"
	"github.com/davidflanagan/vtprotobuf/generator"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"
)

type ObjectSet map[protogen.GoIdent]bool

func (o ObjectSet) String() string {
	return fmt.Sprintf("%#v", o)
}

func (o ObjectSet) Set(s string) error {
	idx := strings.LastIndexByte(s, '.')
	if idx < 0 {
		return fmt.Errorf("invalid object name: %q", s)
	}

	ident := protogen.GoIdent{
		GoImportPath: protogen.GoImportPath(s[0:idx]),
		GoName:       s[idx+1:],
	}
	o[ident] = true
	return nil
}

func main() {
	var (
		features                string
		poolable                = make(ObjectSet)
		internPackageImportPath = "github.com/davidflanagan/vtprotobuf/vtproto"
		internFunctionName      = "Intern"
	)

	var f flag.FlagSet
	f.Var(poolable, "pool", "use memory pooling for this object")
	f.StringVar(&features, "features", "all", "list of features to generate (separated by '+')")
	f.StringVar(&internPackageImportPath, "internPkg", internPackageImportPath, "the package to import for string interning")
	f.StringVar(&internFunctionName, "internFunc", internFunctionName, "the name of the string interning function")

	protogen.Options{ParamFunc: f.Set}.Run(func(plugin *protogen.Plugin) error {
		return generateAllFiles(plugin, strings.Split(features, "+"), poolable, internPackageImportPath, internFunctionName)
	})
}

var SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)

func generateAllFiles(
	plugin *protogen.Plugin,
	featureNames []string,
	poolable ObjectSet,
	internPackageImportPath string,
	internFunctionName string,
) error {
	ext := &generator.Extensions{
		Poolable:                poolable,
		InternPackageImportPath: internPackageImportPath,
		InternFunctionName:      internFunctionName,
	}
	gen, err := generator.NewGenerator(plugin.Files, featureNames, ext)
	if err != nil {
		return err
	}

	for _, file := range plugin.Files {
		if !file.Generate {
			continue
		}

		gf := plugin.NewGeneratedFile(file.GeneratedFilenamePrefix+"_vtproto.pb.go", file.GoImportPath)
		if !gen.GenerateFile(gf, file) {
			gf.Skip()
		}
	}

	plugin.SupportedFeatures = SupportedFeatures
	return nil
}
