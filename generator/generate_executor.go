package generator

import (
	"fmt"
	"github.com/butzopower/natsu/core"
	. "github.com/dave/jennifer/jen"
	"sort"
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

type execChainStruct struct {
	Id            string
	FnContainerId string
}

// func (chain *execChainA) WithTypeA(fn execFnTypeAFn) *execChainB {
//	 chain.fns.execFnTypeA = fn
//	 return &execChainB{chain.fns}
// }

type execChainFn struct {
	Id                string // WithTypeA
	ChainStructId     string // execChainA
	NextChainStructId string // execChainB
	FnContainerId     string // fns
	FnName            string // execFnTypeA
	FnTypeId          string // execFnTypeAFn
	ThisId            string // chain
	ParamId           string // fn
}

type execConstructorFn struct {
	Id                 string
	FirstChainStructId string
	ContainerStructId  string
}

type containersExecFnType map[string]execFnType
type containersExecFn map[string]execFn
type containersExecChainStruct map[string]execChainStruct

func generateExecutor(
	file *File,
	c containers,
	sumType sumTypeStruct,
) {
	namespace := sumType.Id
	execFnTypesMetadata := buildExecFnTypesMetadata(c)
	execFnContainerStructMetadata := buildExecFnContainerMetadata(namespace, execFnTypesMetadata)
	executorStructMetadata := buildExecutorStructMetadata(namespace)
	execChainStructsMetadata := buildExecutorChainStructMetadata(namespace, c)
	execChainFnMetadata := buildExecutorChainFnMetadata(execFnContainerStructMetadata, execChainStructsMetadata, executorStructMetadata)
	execConstructorFnMetadata := buildExecutorConstructorFnMetadata(namespace, execChainFnMetadata[0], execFnContainerStructMetadata)

	generateExecFnTypes(file, execFnTypesMetadata)
	generateExecFnContainerStruct(file, execFnContainerStructMetadata)
	generateExecutorStruct(file, executorStructMetadata, execFnContainerStructMetadata.Id)
	generateExecutorChainStructs(file, execChainStructsMetadata, execFnContainerStructMetadata.Id)
	generateExecutorChainFns(file, execChainFnMetadata)
	generateExecutorExec(file, execFnContainerStructMetadata, executorStructMetadata, sumType)
	generateExecutorConstructorFn(file, execConstructorFnMetadata)
}

func buildExecutorChainStructMetadata(
	namespace string,
	c containers,
) containersExecChainStruct {
	containersToReturn := make(containersExecChainStruct)

	for containerId, term := range c {
		id := fmt.Sprintf("executor%sChain%s", namespace, term.Local)
		fnContainerId := "fns"
		containersToReturn[containerId] = execChainStruct{
			Id:            id,
			FnContainerId: fnContainerId,
		}
	}

	return containersToReturn
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
) execChainStruct {
	id := fmt.Sprintf("executor%s", namespace)
	containerId := "fns"

	return execChainStruct{
		Id:            id,
		FnContainerId: containerId,
	}
}

func buildExecutorChainFnMetadata(
	fnContainerStruct execFnContainerStruct,
	execChainStructs containersExecChainStruct,
	executorStruct execChainStruct,
) []execChainFn {

	type chainFn struct {
		Id            string // WithTypeA
		ChainStructId string // doChainA
		FnContainerId string // fns
		FnName        string // doFnTypeA
		FnTypeId      string // doFnTypeAFn
	}

	var chainFns []chainFn

	for containerName, fn := range fnContainerStruct.Fns {
		id := fmt.Sprintf("With%s", fn.Term.Local)
		chainStruct := execChainStructs[containerName]

		chainFn := chainFn{
			Id:            id,
			ChainStructId: chainStruct.Id,
			FnContainerId: chainStruct.FnContainerId,
			FnName:        fn.FnName,
			FnTypeId:      fn.Type.FnTypeName,
		}

		chainFns = append(chainFns, chainFn)
	}

	sort.Slice(chainFns, func(i, j int) bool {
		return chainFns[i].Id < chainFns[j].Id
	})

	thisId := "chain"
	paramId := "fn"

	var execChainFns []execChainFn

	for i := 0; i < len(chainFns)-1; i++ {
		currentChain := chainFns[i]
		nextChain := chainFns[i+1]

		execChainFn := execChainFn{
			Id:                currentChain.Id,
			ChainStructId:     currentChain.ChainStructId,
			NextChainStructId: nextChain.ChainStructId,
			FnContainerId:     currentChain.FnContainerId,
			FnName:            currentChain.FnName,
			FnTypeId:          currentChain.FnTypeId,
			ThisId:            thisId,
			ParamId:           paramId,
		}

		execChainFns = append(execChainFns, execChainFn)
	}

	lastChain := chainFns[len(chainFns)-1]

	lastExecChainFn := execChainFn{
		Id:                lastChain.Id,
		ChainStructId:     lastChain.ChainStructId,
		NextChainStructId: executorStruct.Id,
		FnContainerId:     lastChain.FnContainerId,
		FnName:            lastChain.FnName,
		FnTypeId:          lastChain.FnTypeId,
		ThisId:            thisId,
		ParamId:           paramId,
	}

	execChainFns = append(execChainFns, lastExecChainFn)

	return execChainFns
}

func buildExecutorConstructorFnMetadata(
	namespace string,
	firstChainFn execChainFn,
	fnContainerStruct execFnContainerStruct,
) execConstructorFn {
	id := fmt.Sprintf("%sExecutor", namespace)

	return execConstructorFn{
		Id:                 id,
		FirstChainStructId: firstChainFn.ChainStructId,
		ContainerStructId:  fnContainerStruct.Id,
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
	executorStruct execChainStruct,
	execFnContainerStructId string,
) {
	generateExecutorChainStruct(file, executorStruct, execFnContainerStructId)
}

func generateExecutorChainStructs(
	file *File,
	containerChainStructs containersExecChainStruct,
	execFnContainerStructId string,
) {
	for _, chainStruct := range containerChainStructs {
		generateExecutorChainStruct(file, chainStruct, execFnContainerStructId)
	}
}

func generateExecutorChainFns(file *File, execChainFns []execChainFn) {
	for _, chainFn := range execChainFns {
		file.Func().
			Params(Id(chainFn.ThisId).Op("*").Id(chainFn.ChainStructId)).
			Id(chainFn.Id).Params(Id(chainFn.ParamId).Id(chainFn.FnTypeId)).
			Op("*").Id(chainFn.NextChainStructId).
			Block(
				Id(chainFn.ThisId).Dot(chainFn.FnContainerId).Dot(chainFn.FnName).Op("=").Id(chainFn.ParamId),
				Return(Op("&").Id(chainFn.NextChainStructId).Values(Id(chainFn.ThisId).Dot(chainFn.FnContainerId))),
			)
	}
}

func generateExecutorChainStruct(
	file *File,
	executorStruct execChainStruct,
	execFnContainerStructId string,
) {
	file.Type().Id(executorStruct.Id).Struct(Id(executorStruct.FnContainerId).Id(execFnContainerStructId))
}

func generateExecutorExec(
	file *File,
	execFnContainer execFnContainerStruct,
	executorStruct execChainStruct,
	sumType sumTypeStruct,
) {
	execMethodName := "Exec"
	thisId := "e"
	sumTypeInstanceId := "sum"
	narrowedInstanceId := "value"
	file.Func().Params(Id(thisId).Op("*").Id(executorStruct.Id)).Id(execMethodName).Params(Id(sumTypeInstanceId).Id(sumType.Id)).Block(
		Switch(Id(narrowedInstanceId).Op(":=").Id(sumTypeInstanceId).Dot(sumType.ValueId).Dot("").Parens(Type())).Block(execSwitchOptions(execFnContainer.Fns, executorStruct, thisId, narrowedInstanceId)...),
	)
}

func execSwitchOptions(
	execFnContainer containersExecFn,
	executorStruct execChainStruct,
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

func generateExecutorConstructorFn(
	file *File,
	fn execConstructorFn,
) {
	file.Func().Id(fn.Id).Params().Op("*").Id(fn.FirstChainStructId).Block(
		Return().Op("&").Id(fn.FirstChainStructId).Values(Id(fn.ContainerStructId).Values()),
	)
}
