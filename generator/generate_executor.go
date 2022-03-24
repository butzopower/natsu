package generator

import (
	"fmt"
	. "github.com/dave/jennifer/jen"
	"natsu/core"
)

type execFnType struct {
	Term       core.TermPath
	FnTypeName string
}

type execFn struct {
	Term   core.TermPath
	Type   execFnType
	FnName string
}

type execFnContainerStruct struct {
	Id  string
	Fns containersExecFn
}

type executorStruct struct {
	Id            string
	FnContainerId string
}

type containersExecFnType map[string]execFnType
type containersExecFn map[string]execFn

func generateExecutor(
	file *File,
	r core.Result,
	c containers,
	sumType sumTypeStruct,
) {
	execFnTypesMetadata := buildExecFnTypesMetadata(c)
	execFnContainerStructMetadata := buildExecFnContainerMetadata(r.Union.Local, execFnTypesMetadata)
	executorStructMetadata := buildExecutorStructMetadata(r.Union.Local)

	generateExecFnTypes(file, execFnTypesMetadata)
	generateExecFnContainerStruct(file, execFnContainerStructMetadata)
	generateExecutorStruct(file, executorStructMetadata, execFnContainerStructMetadata.Id)
	generateExecutorExec(file, execFnContainerStructMetadata, executorStructMetadata, sumType)
}

func buildExecFnTypesMetadata(c containers) containersExecFnType {
	containersToReturn := make(containersExecFnType)

	for containerId, term := range c {
		functionTypeName := fmt.Sprintf("execFn%sFn", term.Local)

		containersToReturn[containerId] = execFnType{
			Term:       term,
			FnTypeName: functionTypeName,
		}
	}

	return containersToReturn
}

func buildExecFnContainerMetadata(
	namespace string,
	metadata containersExecFnType,
) execFnContainerStruct {
	structId := fmt.Sprintf("exec%sFns", namespace)
	fns := make(containersExecFn)

	for containerId, fnType := range metadata {
		functionName := fmt.Sprintf("execFn%s", fnType.Term.Local)

		fns[containerId] = execFn{
			Term:   fnType.Term,
			Type:   fnType,
			FnName: functionName,
		}
	}

	return execFnContainerStruct{
		Id:  structId,
		Fns: fns,
	}
}

func buildExecutorStructMetadata(
	namespace string,
) executorStruct {
	id := fmt.Sprintf("executor%s", namespace)
	containerId := "fns"

	return executorStruct{
		Id:            id,
		FnContainerId: containerId,
	}
}

func generateExecFnTypes(file *File, c containersExecFnType) {
	for _, fn := range c {
		file.Type().Id(fn.FnTypeName).Func().Params(Id("x").Qual(fn.Term.Package, fn.Term.Local))
	}
}

func generateExecFnContainerStruct(file *File, s execFnContainerStruct) {
	var fields []Code

	for _, execFn := range s.Fns {
		field := Id(execFn.FnName).Id(execFn.Type.FnTypeName)
		fields = append(fields, field)
	}

	file.Type().Id(s.Id).Struct(fields...)
}

func generateExecutorStruct(
	file *File,
	executorStruct executorStruct,
	execFnContainerStructId string,
) {
	file.Type().Id(executorStruct.Id).Struct(Id(executorStruct.FnContainerId).Id(execFnContainerStructId))
}

func generateExecutorExec(
	file *File,
	execFnContainer execFnContainerStruct,
	executorStruct executorStruct,
	sumType sumTypeStruct,
) {
	execMethodName := "Exec"
	thisId := "e"
	sumTypeInstanceId := "sum"
	narrowedInstanceId := "value"
	file.Func().Params(Id(thisId).Id(executorStruct.Id)).Id(execMethodName).Params(Id(sumTypeInstanceId).Id(sumType.Id)).Block(
		Switch(Id(narrowedInstanceId).Op(":=").Id(sumTypeInstanceId).Dot(sumType.ValueId).Dot("").Parens(Type())).Block(execSwitchOptions(execFnContainer.Fns, executorStruct, thisId, narrowedInstanceId)...),
	)
}

func execSwitchOptions(
	execFnContainer containersExecFn,
	executorStruct executorStruct,
	executorId string,
	narrowedInstanceId string,
) []Code {
	var options []Code

	for containerId, execFn := range execFnContainer {
		var caseStatement = Case(Id(containerId)).Block(
			Id(executorId).Dot(executorStruct.FnContainerId).Dot(execFn.FnName).Call(Id(narrowedInstanceId).Dot("v")),
		)

		options = append(options, caseStatement)
	}

	return options
}
