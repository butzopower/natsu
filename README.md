# 夏 – natsu – summer.

sum type generator for go 1.18+

## Features

What benefit does a sum type give over a type switch?

Consider the below example of function constrained by a type union that uses a type switch: 

```go
type Cat struct {
	Name string 
	SharpClaws bool
}

type Dog struct {
	Name string
	Trained bool
}

type Pet interface {
	Cat | Dog
}

func Cuddle[T Pet](pet T) {
	var switchablePet interface{}
	switchablePet = pet
	switch p := switchablePet.(type) {
	case Cat:
		if p.SharpClaws {
			print("ow, it scratched me")
		}
	case Dog:
		if !p.Trained {
			print("ah, it slobbered me")
		}
	case string:
		// shouldn't be matchable
		print("uh wut")
	default:
		// shouldn't be required
		print("there is no pet?")
	}
}

func main() {
	Cuddle(Cat{Name: "Tex", SharpClaws: true})
	Cuddle(Dog{Name: "Fifi", Trained: false})
	
	// does not compile: string does not implement Pet 
	//Cuddle("strings are what cats play with")
}

```

## Usage