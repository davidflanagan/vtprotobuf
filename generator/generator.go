// Copyright (c) 2021 PlanetScale Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package generator

import (
	"runtime/debug"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/runtime/protoimpl"
)

type featureHelpers struct {
	path    protogen.GoImportPath
	feature int
}

type Extensions struct {
	Poolable                map[protogen.GoIdent]bool
	InternPackageImportPath string
	InternFunctionName      string
}

type Generator struct {
	seen     map[featureHelpers]bool
	ext      *Extensions
	features []Feature
	local    map[string]bool
}

func NewGenerator(allFiles []*protogen.File, featureNames []string, ext *Extensions) (*Generator, error) {
	features, err := findFeatures(featureNames)
	if err != nil {
		return nil, err
	}

	local := make(map[string]bool)
	for _, f := range allFiles {
		if f.Generate {
			local[string(f.Desc.Package())] = true
		}
	}

	return &Generator{
		seen:     make(map[featureHelpers]bool),
		ext:      ext,
		features: features,
		local:    local,
	}, nil
}

func (gen *Generator) GenerateFile(gf *protogen.GeneratedFile, file *protogen.File) bool {
	p := &GeneratedFile{
		GeneratedFile: gf,
		Ext:           gen.ext,
		LocalPackages: gen.local,
	}

	p.P("// Code generated by protoc-gen-go-vtproto. DO NOT EDIT.")
	if bi, ok := debug.ReadBuildInfo(); ok {
		p.P("// protoc-gen-go-vtproto version: ", bi.Main.Version)
	}
	p.P("// source: ", file.Desc.Path())
	p.P()
	p.P("package ", file.GoPackageName)
	p.P()

	protoimplPackage := protogen.GoImportPath("google.golang.org/protobuf/runtime/protoimpl")
	p.P("const (")
	p.P("// Verify that this generated code is sufficiently up-to-date.")
	p.P("_ = ", protoimplPackage.Ident("EnforceVersion"), "(", protoimpl.GenVersion, " - ", protoimplPackage.Ident("MinVersion"), ")")
	p.P("// Verify that runtime/protoimpl is sufficiently up-to-date.")
	p.P("_ = ", protoimplPackage.Ident("EnforceVersion"), "(", protoimplPackage.Ident("MaxVersion"), " - ", protoimpl.GenVersion, ")")
	p.P(")")
	p.P()

	var generated bool
	for fidx, feat := range gen.features {
		featGenerator := feat(p)
		if featGenerator.GenerateFile(file) {
			generated = true

			helpersForPlugin := featureHelpers{
				path:    file.GoImportPath,
				feature: fidx,
			}
			if !gen.seen[helpersForPlugin] {
				featGenerator.GenerateHelpers()
				gen.seen[helpersForPlugin] = true
			}
		}
	}

	return generated
}
