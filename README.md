# Introduce

To use `go/ast` easily.

This lib may not completed, but work well in my project, I'll keep improve it with my project.

## Install

```sh
go get -u -v github.com/szyhf/go-aster
```

## Example

```go
import (
	"fmt"

	aster "github.com/szyhf/go-aster"
)

func main(){
	// an go path directory
	goDir := "./data"

	pkgsTyp, err := aster.ParseDir(goDir, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	// simple display the result
	fmt.Println(pkgsTyp[0].String())
}
```

## TODO

More details please see the `test` directory.