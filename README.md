## Usage

```go
func main() {
    e := mua.GetEngine().Mapping("html", "/usr/local/wwwroot")

    e.Use(mua.ExceptionHandler).GET("/", func(ctx *mua.Context) {
        ctx.EchoString(ctx.Path)
    })
    e.GET("/test1/:name/test2", func(ctx *mua.Context) {
        ctx.EchoString(ctx.Params["name"])
    })
    e.GET("/test1/*filepath", func(ctx, *mua.Context) {
        ctx.EchoString(ctx.Params["filepath"])
    })
    e.Group("/api").Use(func1, func2).GET("/test", func(ctx, *muaConetxt) {})

    e.LoadTmpls("template/*", nil).GET("test", func(ctx *mua.Context) {
        ctx.EchoTMPL("test.tpl", nil)
    })
}
```