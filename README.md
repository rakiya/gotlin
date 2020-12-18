# gotlin

Gotlin make you be able to write Go lang like Kotlin.

```
package main

import(
  "fmt"
  "gotlin/stream"
)

func main() {
  v := gotlin.NewStream([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}).
    Filter(func(it int) bool { return it%2 == 0 }).
    Map(func(it int) int { return 3 * it }).
    Inject(0, func(res int, it int) int { return res + it }).
    Let(func(it int) float64 { return float64(it) }).
    End().(float64)
   
  fmt.Println(v) // => 90
}
```

- You can use `apply` and `let` function like them in Kotlin.

- You can use `map` / `filter` / `inject` function like them of a functional programming.
