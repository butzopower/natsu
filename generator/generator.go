package generator

import (
	"fmt"
	. "github.com/dave/jennifer/jen"
	"natsu/core"
)

func Generate(r core.Result) string {
	file := NewFilePath(r.Path)

	containerInterface, containerInterfaceFn := generateContainerInterface(file, r)
	generateContainers(file, r, containerInterfaceFn)
	generateStruct(file, r, containerInterface)
	generateConstructor(file, r)

	return fmt.Sprintf("%#v", file)
}

func generateContainerInterface(file *File, r core.Result) (string, string) {
	interfaceId := fmt.Sprintf("container%s", r.Union.Local)
	interfaceFn := fmt.Sprintf("%sFn", interfaceId)

	file.Type().Id(interfaceId).Interface(
		Id(interfaceFn).Params(),
	)

	return interfaceId, interfaceFn
}

func generateContainers(file *File, r core.Result, interfaceFn string) {
	for _, term := range r.Terms {
		memberContainerId := fmt.Sprintf("container%s", term.Local)
		file.Type().Id(memberContainerId).Struct(
			Id("v").Qual(term.Package, term.Local),
		)

		file.Func().Params(Id("c").Id(memberContainerId)).Id(interfaceFn).Params().Block()
	}
}

func generateConstructor(file *File, r core.Result) {
	typeId := "T"
	constructorId := fmt.Sprintf("%sOf", r.Union.Local)
	file.Func().Id(constructorId).
		Types(Id(typeId).Qual(r.Union.Package, r.Union.Local)).
		Params(Id("value").Id(typeId)).
		Block()
}

func generateStruct(file *File, r core.Result, containerInterfaceId string) {
	structId := fmt.Sprintf("Tagged%s", r.Union.Local)

	file.Type().Id(structId).Struct(
		Id("v").Id(containerInterfaceId),
	)
}
