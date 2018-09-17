package main

// Don't edit this file. 

var tmpl = `// Code generated by bindgen, DO NOT EDIT.

package {{.package}}

import "github.com/newalchemylimited/seth"

{{range $c := .contracts}}

	type {{$c.Name}} struct {
		addr  *seth.Address
		s     *seth.Sender
	}

	func New{{$c.Name}}(addr *seth.Address, sender *seth.Sender) *{{$c.Name}} {
		return &{{$c.Name}}{addr: addr, s: sender}
	}

	{{range $d := $c.ABI}}
		
		{{if eq $d.Type "function" }}

			{{if $d.Constant}}

				// {{FuncName $d.Name}} calls the solidity view {{$c.Name}}.{{$d.Signature}}
				func (c *{{$c.Name}}) {{FuncName $d.Name}}({{range $i, $input := $d.Inputs}}{{if gt $i 0}}, {{end}}{{ArgName $input.Name}} {{ArgType $input.Type}}{{end}}) ({{range $i, $output := $d.Outputs}}{{ArgName $output.Name}} {{RetType $output.Type}}, {{end}}err error) {
					d := seth.NewABIDecoder({{range $i, $output := $d.Outputs}}{{if gt $i 0}}, {{end}}&{{ArgName $output.Name}}{{end}})
					err = c.s.ConstCall(c.addr, "{{$d.Signature}}", d{{range $i, $input := $d.Inputs}}, {{ArgName $input.Name}}{{end}})
					return
				}
			
			{{else}}

				// {{FuncName $d.Name}} calls the solidity function {{$c.Name}}.{{$d.Signature}}
				func (c *{{$c.Name}}) {{FuncName $d.Name}}({{range $i, $input := $d.Inputs}}{{if gt $i 0}}, {{end}}{{ArgName $input.Name}} {{ArgType $input.Type}}{{end}}) (res seth.Hash, err error) {
					return c.s.Send(c.addr, "{{$d.Signature}}"{{range $i, $input := $d.Inputs}}, {{ArgName $input.Name}}{{end}})
				}

			{{end}}

		{{end}}

		{{if eq $d.Type "event" }}

			type {{FuncName $d.Name}}Event struct {
				Log *seth.Log{{range $i, $input := $d.Inputs}}
					{{ArgNameUpper $input.Name}} {{ArgType $input.Type}}{{end}}
			}

			func (e *{{FuncName $d.Name}}Event) FromABI(data []byte) error {
				return seth.DecodeABI(data{{range $i, $input := $d.Inputs}}, &e.{{ArgNameUpper $input.Name}}{{end}})
			}

			type {{FuncName $d.Name}}EventIterator struct {
				Event *{{FuncName $d.Name}}Event
				Error error
				Close func()

				errors chan error
				events chan *{{FuncName $d.Name}}Event
			}

			func (i *{{FuncName $d.Name}}EventIterator) Next() bool {

				select {
					case i.Error = <-i.errors:
						return false
					case i.Event = <-i.events:
						return i.Event != nil
				}

			}
			//outChan chan *{{FuncName $d.Name}}Event, close func(), errChan chan error

			func (c *{{$c.Name}}) Filter{{FuncName $d.Name}}Event(ctx context.Context, start, end int64) (*{{FuncName $d.Name}}EventIterator, error) {
				

				topic := seth.HashString("{{$d.Signature}}")
				filter, err := c.s.FilterTopics([]*seth.Hash{&topic}, c.addr, start, end)
				if err != nil {
					return nil, err
				}

				i := &{{FuncName $d.Name}}EventIterator{
					errors: make(chan error, 1),
					events: make(chan *{{FuncName $d.Name}}Event),
					Close: filter.Close,
				}

				go func() {
					defer filter.Close()
					for {
						if filter.Err() != nil {
							i.errors <- err
							return
						}
						select {
						case <-ctx.Done():
							i.errors <- ctx.Err()
							return
						case msg := <-filter.Out():
							if msg == nil {
								i.events <- nil
								return
							}

							x := &{{FuncName $d.Name}}Event{
								Log: msg,
							}
							if err := x.FromABI(msg.Data); err != nil {
								i.errors <- err
								return
							}
							i.events <- x
						}
					}
				}()

				return i, nil
			}
		{{end}}

	{{end}}
	

{{end}}


{{range $c := .contracts}}

	var {{$c.Name}}Code = {{CodeVar $c.Code}}

{{end}}`
