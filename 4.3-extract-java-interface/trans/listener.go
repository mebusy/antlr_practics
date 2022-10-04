package trans

import (
	"javaparser/parser" // module/package
	// "github.com/antlr/antlr4/runtime/Go/antlr/v4"
	"fmt"
    // "strings"
)

// your listener, say you want convert C declaration to C#
type MyListener struct {
	*parser.BaseJava8ParserListener
}

func NewListener( _ *parser.Java8Parser  ) *MyListener {
    // parser not used
    return &MyListener{  } 
}

func (l *MyListener) EnterClassDeclaration(ctx *parser.ClassDeclarationContext) {
    // println("enter class")
    // get parser.NormalClassDeclarationContext
    nc := ctx.NormalClassDeclaration().(*parser.NormalClassDeclarationContext )
    id := nc.Identifier()

    fmt.Println( "interface I"+id.GetText()+" {" )
}

func (l *MyListener) ExitClassDeclaration(ctx *parser.ClassDeclarationContext) {
    // println("exit class")
    fmt.Println( "}" )
}

func (l *MyListener) EnterMethodDeclaration(ctx *parser.MethodDeclarationContext) {
    // need parser to get tokens, to get source file text
    tokens := ctx.GetParser().GetTokenStream() 
    // NOTE: 1. tokens.GetTokenSource().GetInputStream() 
    //           is CharStream of original source file

    method_modifiers := ""
    if len(ctx.AllMethodModifier()) > 0 {
        for _, m := range ctx.AllMethodModifier() {
            method_modifiers += m.GetText() + " "
        }
    }


    mh_context := ctx.MethodHeader().(*parser.MethodHeaderContext)

    // fmt.Printf( "%+v %v\n", mh_context.GetSourceInterval().Start , mh_context.GetStart().() )
    is :=tokens.GetTokenSource().GetInputStream()
    // NOTE: call GetStart() on a token will return the start index in source file
    src_start := mh_context.GetStart().GetStart() // index in source file
    src_stop := mh_context.GetStop().GetStop() // index in source file

    method_header := is.GetText( src_start, src_stop )

    fmt.Printf( "\t%s %+v;\n", method_modifiers , method_header )
}


func (l *MyListener) EnterImportDeclaration(ctx *parser.ImportDeclarationContext) {
    // println("enter import")
    // get parser.NormalClassDeclarationContext
    // nc := ctx.NormalClassDeclaration().(*parser.NormalClassDeclarationContext )
    // id := nc.Identifier()

    // fmt.Println( "interface I"+id.GetText()+" {" )
    is := ctx.GetParser().GetTokenStream().GetTokenSource().GetInputStream()
    src_start := ctx.GetStart().GetStart() // index in source file
    src_stop := ctx.GetStop().GetStop() // index in source file

    fmt.Printf( "%s\n", is.GetText( src_start, src_stop ) )
}




















