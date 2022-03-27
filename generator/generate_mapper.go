package generator

import (
    "fmt"
    "github.com/butzopower/natsu/core"
    . "github.com/dave/jennifer/jen"
    "sort"
)

const genericKey = "T"

type mapFnType struct {
    Term       core.TermPath
    FnTypeName string
}

type mapFn struct {
    Term   core.TermPath
    Type   mapFnType
    FnName string
}

type mapFnContainerStruct struct {
    Id  string
    Fns containersMapFn
}

type mapChainStruct struct {
    Id            string
    FnContainerId string
}

type mapChainFn struct {
    Id                string
    ChainStructId     string
    NextChainStructId string
    FnContainerId     string
    FnName            string
    FnTypeId          string
    ThisId            string
    ParamId           string
}

type mapConstructorFn struct {
    Id                 string
    FirstChainStructId string
    ContainerStructId  string
}

type containersMapFnType map[string]mapFnType
type containersMapFn map[string]mapFn
type containersMapChainStruct map[string]mapChainStruct

func generateMapper(
    file *File,
    c containers,
    sumType sumTypeStruct,
) {
    namespace := sumType.Id
    mapFnTypesMetadata := buildMapFnTypesMetadata(namespace, c)
    mapFnContainerStructMetadata := buildMapFnContainerMetadata(namespace, mapFnTypesMetadata)
    mapperStructMetadata := buildMapperStructMetadata(namespace)
    mapChainStructsMetadata := buildMapperChainStructMetadata(namespace, c)
    mapChainFnMetadata := buildMapperChainFnMetadata(mapFnContainerStructMetadata, mapChainStructsMetadata, mapperStructMetadata)
    mapConstructorFnMetadata := buildMapperConstructorFnMetadata(namespace, mapChainFnMetadata[0], mapFnContainerStructMetadata)

    file.Line()
    file.Comment("Mapper")

    generateMapFnTypes(file, mapFnTypesMetadata)
    generateMapFnContainerStruct(file, mapFnContainerStructMetadata)
    generateMapperStruct(file, mapperStructMetadata, mapFnContainerStructMetadata.Id)
    generateMapperChainStructs(file, mapChainStructsMetadata, mapFnContainerStructMetadata.Id)
    generateMapperChainFns(file, mapChainFnMetadata)
    generateMapperMap(file, mapFnContainerStructMetadata, mapperStructMetadata, sumType)
    generateMapperConstructorFn(file, mapConstructorFnMetadata)
}

func buildMapperChainStructMetadata(
    namespace string,
    c containers,
) containersMapChainStruct {
    containersToReturn := make(containersMapChainStruct)

    for containerId, term := range c {
        id := fmt.Sprintf("mapper%sChain%s", namespace, term.Local)
        fnContainerId := "fns"
        containersToReturn[containerId] = mapChainStruct{
            Id:            id,
            FnContainerId: fnContainerId,
        }
    }

    return containersToReturn
}

func buildMapFnTypesMetadata(namespace string, c containers) containersMapFnType {
    containersToReturn := make(containersMapFnType)

    for containerId, term := range c {
        functionTypeName := fmt.Sprintf("mapFn%s%sFn", namespace, term.Local)

        containersToReturn[containerId] = mapFnType{
            Term:       term,
            FnTypeName: functionTypeName,
        }
    }

    return containersToReturn
}

func buildMapFnContainerMetadata(
    namespace string,
    metadata containersMapFnType,
) mapFnContainerStruct {
    structId := fmt.Sprintf("map%sFns", namespace)
    fns := make(containersMapFn)

    for containerId, fnType := range metadata {
        functionName := fmt.Sprintf("mapFn%s", fnType.Term.Local)

        fns[containerId] = mapFn{
            Term:   fnType.Term,
            Type:   fnType,
            FnName: functionName,
        }
    }

    return mapFnContainerStruct{
        Id:  structId,
        Fns: fns,
    }
}

func buildMapperStructMetadata(
    namespace string,
) mapChainStruct {
    id := fmt.Sprintf("mapper%s", namespace)
    containerId := "fns"

    return mapChainStruct{
        Id:            id,
        FnContainerId: containerId,
    }
}

func buildMapperChainFnMetadata(
    fnContainerStruct mapFnContainerStruct,
    execChainStructs containersMapChainStruct,
    mapperStruct mapChainStruct,
) []mapChainFn {

    type chainFn struct {
        Id            string
        ChainStructId string
        FnContainerId string
        FnName        string
        FnTypeId      string
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

    var mapChainFns []mapChainFn

    for i := 0; i < len(chainFns)-1; i++ {
        currentChain := chainFns[i]
        nextChain := chainFns[i+1]

        mapChainFn := mapChainFn{
            Id:                currentChain.Id,
            ChainStructId:     currentChain.ChainStructId,
            NextChainStructId: nextChain.ChainStructId,
            FnContainerId:     currentChain.FnContainerId,
            FnName:            currentChain.FnName,
            FnTypeId:          currentChain.FnTypeId,
            ThisId:            thisId,
            ParamId:           paramId,
        }

        mapChainFns = append(mapChainFns, mapChainFn)
    }

    lastChain := chainFns[len(chainFns)-1]

    lastMapChainFn := mapChainFn{
        Id:                lastChain.Id,
        ChainStructId:     lastChain.ChainStructId,
        NextChainStructId: mapperStruct.Id,
        FnContainerId:     lastChain.FnContainerId,
        FnName:            lastChain.FnName,
        FnTypeId:          lastChain.FnTypeId,
        ThisId:            thisId,
        ParamId:           paramId,
    }

    mapChainFns = append(mapChainFns, lastMapChainFn)

    return mapChainFns
}

func buildMapperConstructorFnMetadata(
    namespace string,
    firstChainFn mapChainFn,
    fnContainerStruct mapFnContainerStruct,
) mapConstructorFn {
    id := fmt.Sprintf("%sMapper", namespace)

    return mapConstructorFn{
        Id:                 id,
        FirstChainStructId: firstChainFn.ChainStructId,
        ContainerStructId:  fnContainerStruct.Id,
    }
}

func generateMapFnTypes(file *File, c containersMapFnType) {
    for _, fn := range c {
        file.Type().Id(fn.FnTypeName).Types(Id(genericKey).Any()).Func().Params(Id("x").Qual(fn.Term.Package, fn.Term.Local)).Id(genericKey)
    }
}

func generateMapFnContainerStruct(file *File, s mapFnContainerStruct) {
    var fields []Code

    for _, execFn := range s.Fns {
        field := Id(execFn.FnName).Id(execFn.Type.FnTypeName).Types(Id(genericKey))
        fields = append(fields, field)
    }

    file.Type().Id(s.Id).Types(Id(genericKey).Any()).Struct(fields...)
}

func generateMapperStruct(
    file *File,
    mapperStruct mapChainStruct,
    mapFnContainerStructId string,
) {
    generateMapperChainStruct(file, mapperStruct, mapFnContainerStructId)
}

func generateMapperChainStructs(
    file *File,
    containerChainStructs containersMapChainStruct,
    mapFnContainerStructId string,
) {
    for _, chainStruct := range containerChainStructs {
        generateMapperChainStruct(file, chainStruct, mapFnContainerStructId)
    }
}

func generateMapperChainFns(file *File, mapChainFns []mapChainFn) {
    for _, chainFn := range mapChainFns {
        file.Func().
            Params(Id(chainFn.ThisId).Op("*").Id(chainFn.ChainStructId).Types(Id(genericKey))).
            Id(chainFn.Id).Params(Id(chainFn.ParamId).Id(chainFn.FnTypeId).Types(Id(genericKey))).
            Op("*").Id(chainFn.NextChainStructId).Types(Id(genericKey)).
            Block(
                Id(chainFn.ThisId).Dot(chainFn.FnContainerId).Dot(chainFn.FnName).Op("=").Id(chainFn.ParamId),
                Return(Op("&").Id(chainFn.NextChainStructId).Types(Id(genericKey)).Values(Id(chainFn.ThisId).Dot(chainFn.FnContainerId))),
            )
    }
}

func generateMapperChainStruct(
    file *File,
    mapperStruct mapChainStruct,
    mapFnContainerStructId string,
) {
    file.Type().Id(mapperStruct.Id).Types(Id(genericKey).Any()).Struct(
        Id(mapperStruct.FnContainerId).Id(mapFnContainerStructId).Types(Id(genericKey)),
    )
}

func generateMapperMap(
    file *File,
    mapFnContainer mapFnContainerStruct,
    mapperStruct mapChainStruct,
    sumType sumTypeStruct,
) {
    mapMethodName := "Map"
    thisId := "e"
    sumTypeInstanceId := "sum"
    narrowedInstanceId := "value"
    file.Func().
        Params(Id(thisId).Op("*").Id(mapperStruct.Id).Types(Id(genericKey))).
        Id(mapMethodName).Params(Id(sumTypeInstanceId).Id(sumType.Id)).
        Id(genericKey).
        Block(
            Switch(Id(narrowedInstanceId).Op(":=").Id(sumTypeInstanceId).Dot(sumType.ValueId).Dot("").Parens(Type())).
                Block(mapSwitchOptions(mapFnContainer.Fns, mapperStruct, thisId, narrowedInstanceId, mapMethodName)...),
    )
}

func mapSwitchOptions(
    mapFnContainer containersMapFn,
    mapperStruct mapChainStruct,
    mapperId string,
    narrowedInstanceId string,
    mapMethodName string,
) []Code {
    var options []Code

    for containerId, execFn := range mapFnContainer {
        var caseStatement = Case(Id(containerId)).Block(
            Return(Id(mapperId).Dot(mapperStruct.FnContainerId).Dot(execFn.FnName).Call(Id(narrowedInstanceId).Dot("v"))),
        )

        options = append(options, caseStatement)
    }

    options = append(
        options,
        Default().Block(Panic(Lit("called "+mapMethodName+" with invalid option"))),
    )

    return options
}

func generateMapperConstructorFn(
    file *File,
    fn mapConstructorFn,
) {
    file.Func().Id(fn.Id).Types(Id(genericKey).Any()).Params().Op("*").Id(fn.FirstChainStructId).Types(Id(genericKey)).Block(
        Return().Op("&").Id(fn.FirstChainStructId).Types(Id(genericKey)).Values(Id(fn.ContainerStructId).Types(Id(genericKey)).Values()),
    )
}
