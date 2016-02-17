package main

import (
	"bufio"
	"bytes"
	"flag"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"log"
	"os"
	"strings"
	"text/template"
)

//Struct defines the field set array for the struct
type Struct struct {
	Name   string
	Fields []Field
}

//Field defines the Fields pf
type Field struct {
	FieldName string
	FieldType string
}

var (
	fileName   = flag.String("f", "", "The source file path")
	structName = flag.String("t", "", "The name of the struct for which the files has to be generated")
)

func getTemplate() (temp *template.Template, err error) {

	// Create a new template and parse the letter into it.
	temp = template.
		Must(template.
		New("getter_setter.tpl").
		Funcs(template.FuncMap{
		"lower":    strings.ToLower,
		"receiver": func(s string) string { return strings.ToLower(s) }}).
		ParseFiles("getter_setter.tpl"))
	return
}

func generateCode(t *template.Template, str Struct) (code []byte, err error) {
	var buf bytes.Buffer
	w := bufio.NewWriter(&buf)
	err = t.Execute(w, str)
	if err != nil {
		return
	}

	err = w.Flush()
	if err != nil {
		return
	}

	//format the bytes from bufffer
	code, err = format.Source(buf.Bytes())
	return
}

func generateBuilder(str Struct) (err error) {
	var t *template.Template
	t, err = getTemplate()

	if err != nil {
		return
	}

	//create the template into a bytes Buffer
	var code []byte
	code, err = generateCode(t, str)
	if err != nil {
		return
	}

	//write to a file
	var f *os.File
	f, err = os.Create(strings.ToLower(str.Name) + "_builder.go")
	if err != nil {
		return
	}
	_, err = bytes.NewReader(code).WriteTo(f)
	if err != nil {
		return
	}
	return
}

func structFields(structName string, f *ast.File, src string) (fs []Field) {

	fields := f.Scope.Lookup(structName).Decl.(*ast.TypeSpec).Type.(*ast.StructType).Fields

	for _, field := range fields.List {
		if field.Names[0].String() != "XMLName" {
			fs = append(fs, Field{
				FieldName: field.Names[0].String(),
				FieldType: src[field.Type.Pos()-1 : field.Type.End()]})
		}
	}

	return
}

func getStructInfo(filename string, structName string) (str Struct, err error) {
	var buf bytes.Buffer
	fset := token.NewFileSet()
	var file *os.File
	file, err = os.Open(filename)
	if err != nil {
		return
	}

	buf.ReadFrom(file)
	b := buf.Bytes()
	var f *ast.File
	f, err = parser.ParseFile(fset, "", b, 0)
	if err != nil {
		return
	}

	str.Name = structName
	str.Fields = structFields(structName, f, string(b))
	return
}

func createStruct() (str Struct, err error) {
	flag.Parse()
	str, err = getStructInfo(*fileName, *structName)
	return
}

func main() {
	str, err := createStruct()
	if err != nil {
		log.Fatal("Error happened creating the resource information", err.Error())
	}
	if err = generateBuilder(str); err != nil {
		log.Fatalln("Error happened generating the builder ", err.Error())
	}
}
