package repl

import (
	"bufio"
	"fmt"
	"io"

	"monkey/compiler"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"monkey/vm"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	constants := []object.Object{}
	globals := make([]object.Object, vm.GlobalsSize)

	symbolTable := compiler.NewSymbolTable()
	for i, v := range object.Builtins {
		symbolTable.DefineBuiltin(i, v.Name)
	}

	for {
		fmt.Print(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		_, tokens := lexer.New(line)
		p := parser.New(&tokens)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		comp := compiler.NewWithState(symbolTable, constants)
		err := comp.Compile(program)
		if err != nil {
			fmt.Fprintf(out, "Woops! Compilation failed:\n %s\n", err)
			continue
		}

		code := comp.Bytecode()
		constants = code.Constants

		machine := vm.NewWithGlobalsStore(code, globals)
		err = machine.Run()
		if err != nil {
			fmt.Fprintf(out, "Woops! Executing bytecode failed:\n %s\n", err)
			continue
		}

		lastPopped := machine.LastPoppedStackElem()
		if _, err := io.WriteString(out, lastPopped.Inspect()); err != nil {
			fmt.Fprintf(out, "Woops! Writing to output failed:\n %s\n", err)
		}
		if _, err := io.WriteString(out, "\n"); err != nil {
			fmt.Fprintf(out, "Woops! Writing to output failed:\n %s\n", err)
		}
	}
}

const MONKEY_FACE = `            __,__
   .--.  .-"     "-.  .--.
  / .. \/  .-. .-.  \/ .. \
 | |  '|  /   Y   \  |'  | |
 | \   \  \ 0 | 0 /  /   / |
  \ '- ,\.-"""""""-./, -' /
   ''-' /_   ^ ^   _\ '-''
       |  \._   _./  |
       \   \ '~' /   /
        '._ '-=-' _.'
           '-----'
`

func printParserErrors(out io.Writer, errors []string) {
	if _, err := io.WriteString(out, MONKEY_FACE); err != nil {
		fmt.Fprintf(out, "Woops! Writing to output failed:\n %s\n", err)
	}
	if _, err := io.WriteString(out, "Woops! We ran into some monkey business here!\n"); err != nil {
		fmt.Fprintf(out, "Woops! Writing to output failed:\n %s\n", err)
	}
	if _, err := io.WriteString(out, " parser errors:\n"); err != nil {
		fmt.Fprintf(out, "Woops! Writing to output failed:\n %s\n", err)
	}
	for _, msg := range errors {
		if _, err := io.WriteString(out, "\t"+msg+"\n"); err != nil {
			fmt.Fprintf(out, "Woops! Writing to output failed:\n %s\n", err)
		}
	}
}
