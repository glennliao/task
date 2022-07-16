package ops

import (
	"github.com/dop251/goja"
	"github.com/fatih/color"
	"github.com/fsnotify/fsnotify"
	"github.com/glennliao/task/tasker/op"
)

func init() {
	op.AddOp(op.Op{
		Name: "watch",
		Handler: func(ctx op.Context, args []string) {
			watcher, err := fsnotify.NewWatcher()
			if err != nil {
				panic(err)
			}
			defer watcher.Close()

			watcher.Add(args[0])

			color.Blue("begin watch " + args[0])

			for {
				select {
				case event, ok := <-watcher.Events:
					if ok {
						callback, _ := goja.AssertFunction(ctx.GetOriCall().Argument(1))
						e := map[string]string{
							"name": event.Name,
							"op":   event.Op.String(),
						}
						_, err := callback(ctx.GetOriCall().This, ctx.GetVM().ToValue(e))
						if err != nil {
							panic(err)
						}
					}
				}
			}
		},
	})
}
