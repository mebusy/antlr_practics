

## FAQ

### Vistor vs Listener

1. **listener** methods are called by the ANTLR-provided walker object, 
    - whereas **visitor** methods must walk their children with explicit visit calls. Forgetting to invoke visit() on a node’s children means those subtrees don’t get visited.
2.  In **visitor** pattern you have the ability to direct tree walking 
    - while in **listener** you are only reacting to the tree walker.
3. a visitor uses the **call stack** to manage tree traversals, 
    - whereas the listener uses an **explicit stack** allocated on the heap, managed by a walker. 
4. Visitors work very well if we need application-specific return values because we get to use the built-in Java return value mechanism. 
    - If we prefer not having to explicitly invoke visitor methods to visit children, we can switch to the listener mechanism. Unfortunately, that means giving up the cleanliness of using Java method return values.



## Quick Start

[ANTLR 4简明教程](https://wizardforcel.gitbooks.io/antlr4-short-course/content/calculator-visitor.html)

[Parsing with ANTLR 4 and Go](https://blog.gopheracademy.com/advent-2017/parsing-with-antlr4-and-go/)

### Install

```bash
$ python3 -m pip install antlr4-tools 
$ antlr4
Downloading antlr4-4.11.1-complete.jar
...
```

### Try parsing

```bash
$ antlr4-parse C.g4 compilationUnit  -tree
int *a[] = NULL ;
^D
(compilationUnit:1 (translationUnit:1 (externalDeclaration:2 (declaration:1 (declarationSpecifiers:1 (declarationSpecifier:2 (typeSpecifier:1 int))) (initDeclaratorList:1 (initDeclarator:1 (declarator:1 (pointer:1 *) (directDeclarator:3 (directDeclarator:1 a) [ ])) = (initializer:1 (assignmentExpression:1 (conditionalExpression:1 (logicalOrExpression:1 (logicalAndExpression:1 (inclusiveOrExpression:1 (exclusiveOrExpression:1 (andExpression:1 (equalityExpression:1 (relationalExpression:1 (shiftExpression:1 (additiveExpression:1 (multiplicativeExpression:1 (castExpression:2 (unaryExpression:1 (postfixExpression:1 (primaryExpression:1 NULL))))))))))))))))))) ;))) <EOF>)
```

To get the tokens and trace through the parse:


```bash
$ antlr4-parse C.g4 compilationUnit -tokens -trace
int *a[] = NULL ; 
[@0,0:2='int',<'int'>,1:0]  # tokens start
[@1,4:4='*',<'*'>,1:4]
[@2,5:5='a',<Identifier>,1:5]
[@3,6:6='[',<'['>,1:6]
[@4,7:7=']',<']'>,1:7]
[@5,9:9='=',<'='>,1:9]
[@6,11:14='NULL',<Identifier>,1:11]
[@7,16:16=';',<';'>,1:16]
[@8,19:18='<EOF>',<EOF>,2:0]
enter   compilationUnit, LT(1)=int  # trace start
enter   translationUnit, LT(1)=int
enter   externalDeclaration, LT(1)=int
enter   declaration, LT(1)=int
enter   declarationSpecifiers, LT(1)=int
enter   declarationSpecifier, LT(1)=int
enter   typeSpecifier, LT(1)=int
...
```

Here's how to get a visual tree view:

```bash
$ antlr4-parse C.g4 compilationUnit -gui     
```


### Generating parser code

```bash
$ antlr4 C.g4
$ ls *.java
CBaseListener.java	CLexer.java		CListener.java		CParser.java
```

To generate C++ code from the same grammar:

```bash
$ antlr4 -Dlanguage=Cpp C.g4 
$ ls *.cpp *.h
CBaseListener.cpp	CLexer.cpp		CListener.cpp		CParser.cpp
CBaseListener.h		CLexer.h		CListener.h		CParser.h
```


## A First Example

[antlr target doc](https://github.com/antlr/antlr4/tree/master/doc) 有大量的目标语言说明。

Here we use [golang target](https://github.com/antlr/antlr4/blob/master/doc/go-target.md) for example.


go version: 1.18


```bash
$ mkdir proj_cparser
$ cd proj_cparser
$ go mod init cparser
```

```bash
# default package is `parser`
$ antlr4 -o parser -Dlanguage=Go C.g4
qibinyi@Qis-Mac-mini-6 antlr (master) $ ls parser
C.interp		CLexer.interp		c_base_listener.go	c_listener.go
C.tokens		CLexer.tokens		c_lexer.go		c_parser.go

$ go mod tidy
```

### Using the Lexer

<details>
<summary>
main.go
</summary>

```go
package main

import (
	"cparser/parser" // module/package
	"fmt"
	"log"
	"os"
    "io/ioutil"

	"github.com/antlr/antlr4/runtime/Go/antlr/v4"
)


func main() {
    var lexer *parser.CLexer
    var err error

    if len(os.Args) < 2 {
        data, err := ioutil.ReadAll(os.Stdin)
        if err != nil {
            log.Fatal( "read from stdin error" )
        }
        // Setup the input
        input := antlr.NewInputStream( string(data) )
        // Create the Lexer
        lexer = parser.NewCLexer (input)
    } else {
        // Setup the input
        var input *antlr.FileStream
        input, err = antlr.NewFileStream(os.Args[1])
        if err != nil {
            log.Fatal( "input error:", err  )
        }
        // Create the Lexer
        lexer = parser.NewCLexer (input)
    }

	// Read all tokens
	for {
		t := lexer.NextToken()
		if t.GetTokenType() == antlr.TokenEOF {
			break
		}
		fmt.Printf("%s (%q)\n",
			lexer.SymbolicNames[t.GetTokenType()], t.GetText())
	}
    // end read tokens
}


```

</details>


```bash
$ go run main.go
int a = 3;
^D
Int ("int")
Identifier ("a")
Assign ("=")
Constant ("3")
Semi (";")
```


### Using the Parser

antlr的listener默认对语法树进行前序遍历, antlr go runtime中的ParseTreeListener接口包含EnterEveryRule和ExitEveryRule两个方法：

```go
// VisitTerminal is called when a terminal node is visited.
func (s *BaseCListener) VisitTerminal(node antlr.TerminalNode) {}

// VisitErrorNode is called when an error node is visited.
func (s *BaseCListener) VisitErrorNode(node antlr.ErrorNode) {}

// EnterEveryRule is called when any rule is entered.
func (s *BaseCListener) EnterEveryRule(ctx antlr.ParserRuleContext) {}

// ExitEveryRule is called when any rule is exited.
func (s *BaseCListener) ExitEveryRule(ctx antlr.ParserRuleContext) {}
```

<details>
<summary>
main.go
</summary>

```go
package main

import (
	"cparser/parser" // module/package
	"fmt"
	"log"
	"os"
    "io/ioutil"

	"github.com/antlr/antlr4/runtime/Go/antlr/v4"
    "strings"
)

type MyListener struct {
	*parser.BaseCListener
    p *parser.CParser
    nLvl int  // record tree level
}


func (this *MyListener) EnterEveryRule(ctx antlr.ParserRuleContext) {
    i := ctx.GetRuleIndex()
    ruleName := this.p.RuleNames[i]
    fmt.Printf( "%s==> %s : %s \n", strings.Repeat(" ", this.nLvl), ruleName, ctx.GetText() )

    this.nLvl+=1
}
func (this *MyListener) ExitEveryRule(ctx antlr.ParserRuleContext) {
    this.nLvl-=1

    i := ctx.GetRuleIndex()
    ruleName := this.p.RuleNames[i]
    fmt.Printf( "%s<== %s\n", strings.Repeat(" ", this.nLvl), ruleName )
}

func main() {
    var lexer *parser.CLexer
    var err error

    if len(os.Args) < 2 {
        data, err := ioutil.ReadAll(os.Stdin)
        if err != nil {
            log.Fatal( "read from stdin error" )
        }
        // Setup the input
        input := antlr.NewInputStream( string(data) )
        // Create the Lexer
        lexer = parser.NewCLexer (input)
    } else {
        // Setup the input
        var input *antlr.FileStream
        input, err = antlr.NewFileStream(os.Args[1])
        if err != nil {
            log.Fatal( "input error:", err  )
        }
        // Create the Lexer
        lexer = parser.NewCLexer (input)
    }

    

	stream := antlr.NewCommonTokenStream(lexer, 
                antlr.TokenDefaultChannel /*0*/ ) 
    // Create the Parser
	p := parser.NewCParser(stream)
	p.AddErrorListener(antlr.NewDiagnosticErrorListener(true))
	p.BuildParseTrees = true
    // Finally parse the expression
	gramma := p.CompilationUnit()  // start gramma
	antlr.ParseTreeWalkerDefault.Walk(&MyListener{ p:p } , gramma )
}



```

</details>


```bash
$ go run main.go
int a = 3;
^D
==> compilationUnit : inta=3;<EOF> 
 ==> translationUnit : inta=3; 
  ==> externalDeclaration : inta=3; 
   ==> declaration : inta=3; 
    ==> declarationSpecifiers : int 
     ==> declarationSpecifier : int 
      ==> typeSpecifier : int 
      <== typeSpecifier
     <== declarationSpecifier
    <== declarationSpecifiers
    ==> initDeclaratorList : a=3 
     ==> initDeclarator : a=3 
      ==> declarator : a 
...
```



