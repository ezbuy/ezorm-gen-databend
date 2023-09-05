package main

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/ezbuy/ezorm-gen-databend/internal/databend"
	"github.com/ezbuy/ezorm-gen-databend/internal/handler"
	"github.com/ezbuy/ezorm/v2/pkg/plugin"
)

var (
	genHandlers = []handler.SchemaHandler{}
	printers    = []handler.Printer{}
)

func main() {
	r := bufio.NewScanner(os.Stdin)
	ctx := context.Background()
	for r.Scan() {
		req, err := plugin.Decode(r.Bytes())
		if err != nil {
			fmt.Fprintf(os.Stdout, "decode error: %s\n", err.Error())
			return
		}
		if err := req.Each(func(_ plugin.TemplateName, s plugin.Schema) error {
			d, err := s.GetDriver()
			if err != nil {
				return err
			}
			if d != "mysqlr" {
				fmt.Fprintln(os.Stdout, "driver is not mysqlr , skip")
				return nil
			}
			c := &databend.CreateTable{}
			genHandlers = append(genHandlers, c)
			for _, h := range genHandlers {
				if err := h.Handle(ctx, s); err != nil {
					fmt.Fprintf(os.Stdout, "handle error: %q\n", err)
					return err
				}
			}
			printers = append(printers, c)
			for _, p := range printers {
				if err := p.Print(ctx, req.GetOutputPath()); err != nil {
					fmt.Fprintf(os.Stdout, "print error: %q\n", err)
					return err
				}
			}
			return nil
		}); err != nil {
			fmt.Fprintf(os.Stdout, "each error: %s\n", err.Error())
		}
	}
}
