package main

import (
	"javaparser/parser" // module/package
	"javaparser/trans" // 
	"log"
	"os"
    "io/ioutil"
    "fmt"
	"github.com/antlr/antlr4/runtime/Go/antlr/v4"
)


func main() {
    var lexer *parser.Java8Lexer
    var err error

    if len(os.Args) < 2 {
        data, err := ioutil.ReadAll(os.Stdin)
        if err != nil {
            log.Fatal( "read from stdin error" )
        }
        // Setup the input
        input := antlr.NewInputStream( string(data) )
        // Create the Lexer
        lexer = parser.NewJava8Lexer (input)
    } else {
        // Setup the input
        var input *antlr.FileStream
        input, err = antlr.NewFileStream(os.Args[1])
        if err != nil {
            log.Fatal( "input error:", err  )
        }
        // Create the Lexer
        lexer = parser.NewJava8Lexer (input)
    }

    
    // https://blog.gopheracademy.com/advent-2017/parsing-with-antlr4-and-go/

	stream := antlr.NewCommonTokenStream(lexer, 
                antlr.TokenDefaultChannel /*0*/ ) 
    // Create the Parser
	p := parser.NewJava8Parser(stream)
	p.AddErrorListener(antlr.NewDiagnosticErrorListener(true))
	p.BuildParseTrees = true
    // Finally parse the expression
	gramma := p.CompilationUnit()  // start gramma
    listener := trans.NewListener( p ) 
	antlr.ParseTreeWalkerDefault.Walk( listener , gramma )

    // print rewriter
    listener.PrintRewritedSrc()

    fmt.Println("done")
}

