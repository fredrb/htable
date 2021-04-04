# Go HTable

This is a Hash Table implementation in Go. It was developed for learning and experimenting purposes only, and shouldn't be used for anything serious. 

## Run the example:

```
make countwords
```

## Usage

```go
import (
	"github.com/fredrb/htable"
)

type stringKey = htable.StringKey

func main() {
    ht := htable.New()
    ht.Set(stringKey("key_1"), "Some value")
    ht.Set(stringKey("key_2"), 32)

    value, ok := ht.Get(S("key_1"))
    // ...
}
```