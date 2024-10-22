# unsafepool

unsafepool提供了一种不安全的可复用内存池。一般用在框架代码中（框架代码中常用`reflect.New()`生成对象），用于复用未知结构内存。



**慎用！！！**



## 用法

基本等同`sync.Pool`：

```go
package main

import (
	"reflect"

	"github.com/eachain/unsafepool"
)

type YourType struct {
	A bool           // 清0后为false
	B byte           // 清0后为0
	C int            // 清0后为0
	D float64        // 清0后为
	E []byte         // 清0后为nil
	F map[string]any // 清0后为nil
}

func main() {
	pool := unsafepool.New(reflect.TypeOf(YourType{}))
	value := pool.Get().(*YourType)
	// Put将value所有字段清0。
  // 如果*YourType有Reset()方法，将调用(*YourType).Reset()。
  defer pool.Put(value)

	// ... do something with value
}
```

