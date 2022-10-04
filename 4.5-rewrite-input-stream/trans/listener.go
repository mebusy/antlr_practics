package trans

import (
	"javaparser/parser" // module/package
	"github.com/antlr/antlr4/runtime/Go/antlr/v4"
	"fmt"
    // "strings"
)

type MyListener struct {
	*parser.BaseJava8ParserListener
    Rewriter *antlr.TokenStreamRewriter
    p *parser.Java8Parser
}

func NewListener( p *parser.Java8Parser  ) *MyListener {
    // parser not used
    r := antlr.NewTokenStreamRewriter( p.GetTokenStream() )
    _ = r
    return &MyListener{ Rewriter: r, p:p } 
}


func (l *MyListener) EnterClassBody(ctx *parser.ClassBodyContext) {
    field := "\n\tpublic static final long serialVersionUID = 1L;";
    l.Rewriter.InsertAfterDefault( ctx.GetStart().GetTokenIndex(), field )
}   

func (l *MyListener) PrintRewritedSrc() {
    fmt.Printf( "%s \n", l.Rewriter.GetTextDefault()  )
}


func (l *MyListener) EnterClassDeclaration(ctx *parser.ClassDeclarationContext) {
    // println("enter class")
    // get parser.NormalClassDeclarationContext
    nc := ctx.NormalClassDeclaration().(*parser.NormalClassDeclarationContext )
    id := nc.Identifier()
    
    id.GetSymbol().SetText( "changed_id" )  // change class identifier
}
