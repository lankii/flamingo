package pugjs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path"
	"strings"

	"bytes"

	"github.com/pkg/errors"
)

type (
	// Attr is a simple key-value pair
	Attr struct {
		Name       string
		Val        interface{}
		MustEscape bool
	}

	// Fileref is used by include/extends
	Fileref struct {
		Type, Path string
		Line       int
	}

	// Token defines the basic token read by the tokenizer
	// Tokens form a tree, where the beginning root node starts the document
	Token struct {
		// default
		Type, Name string
		Mode, Val  string
		Line       int

		// subblock
		Block *Token
		// subblock childs
		Nodes []*Token

		// specific information
		AttributeBlocks []string
		Attrs           []*Attr
		MustEscape      bool
		File            *Fileref
		Filename        string
		SelfClosing     bool
		IsInline        *bool
		Obj             string
		Key             string

		// mixin
		Call bool   // mixin call?
		Args string // call args

		// if
		Test                  string // if
		Consequent, Alternate *Token // if result

		// Interpolated
		Expr string
	}
)

// Parse parses a filename into a Token-tree
func (p *renderState) Parse(file string) (*Token, error) {
	bytes, err := ioutil.ReadFile(path.Join(p.path, file) + ".ast.json")

	if err != nil {
		return nil, errors.Errorf("Cannot read %q", file)
	}

	return p.ParseJSON(bytes, file)
}

// ParseJSON parses a json into a Token-tree
func (p *renderState) ParseJSON(bytes []byte, file string) (*Token, error) {
	token := new(Token)

	err := json.Unmarshal(bytes, token)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return token, nil
}

// TokenToTemplate gets named Template from Token
func (p *renderState) TokenToTemplate(name string, t *Token) (*Template, string, error) {
	template := New(name).
		Funcs(funcmap).
		Funcs(p.funcs).
		Option("missingkey=error")

	nodes := p.build(t)
	wr := new(bytes.Buffer)

	for _, b := range nodes {
		b.Render(p, wr, 0)
	}

	for _, b := range p.mixinblocks {
		wr.WriteString("\n" + b)
	}

	for _, b := range p.mixin {
		wr.WriteString("\n" + b)
	}

	template, err := template.Parse(wr.String())

	if err != nil {
		e := err.Error() + "\n"
		for i, l := range strings.Split(wr.String(), "\n") {
			e += fmt.Sprintf("%03d: %s\n", i+1, l)
		}
		return nil, "", errors.New(e)
	}

	return template, wr.String(), nil
}

func (p *renderState) build(parent *Token) (res []Node) {
	if parent == nil {
		return
	}
	for _, t := range parent.Nodes {
		node := p.buildNode(t)
		if node != nil {
			res = append(res, node)
		}
	}
	return
}

func (p *renderState) buildNode(t *Token) (res Node) {
	switch t.Type {
	case "Tag":
		tag := new(Tag)
		tag.Name = t.Name
		tag.IsInline = t.IsInline
		tag.SelfClosing = t.SelfClosing
		for _, a := range t.AttributeBlocks {
			tag.AttributeBlocks = append(tag.AttributeBlocks, JavaScriptExpression(a))
		}
		tag.Block = Block{Nodes: p.build(t.Block)}
		for _, a := range t.Attrs {
			tag.Attrs = append(tag.Attrs, Attribute{Name: a.Name, Val: JavaScriptExpression(fmt.Sprintf("%v", a.Val)), MustEscape: a.MustEscape})
		}

		return tag

	case "Mixin":
		mixin := new(Mixin)
		mixin.Block = Block{Nodes: p.build(t.Block)}
		for _, a := range t.AttributeBlocks {
			mixin.AttributeBlocks = append(mixin.AttributeBlocks, JavaScriptExpression(a))
		}
		mixin.Name = JavaScriptIdentifier(t.Name)
		mixin.Args = t.Args
		for _, a := range t.Attrs {
			mixin.Attrs = append(mixin.Attrs, Attribute{Name: a.Name, Val: JavaScriptExpression(fmt.Sprintf("%v", a.Val)), MustEscape: a.MustEscape})
		}
		mixin.Call = t.Call

		return mixin

	case "Text":
		text := new(Text)
		text.Val = t.Val
		return text

	case "Code":
		code := new(Code)
		code.Val = t.Val
		code.Block = Block{Nodes: p.build(t.Block)}
		code.IsInline = t.IsInline
		code.MustEscape = t.MustEscape
		return code

	case "Conditional":
		cond := new(Conditional)
		cond.Test = JavaScriptExpression(t.Test)
		cond.Consequent = p.buildNode(t.Consequent)
		if t.Alternate != nil {
			cond.Alternate = p.buildNode(t.Alternate)
		}
		return cond

	case "Each":
		each := new(Each)
		each.Val = JavaScriptIdentifier(t.Val)
		each.Key = JavaScriptIdentifier(t.Key)
		each.Obj = JavaScriptExpression(t.Obj)
		each.Block = Block{Nodes: p.build(t.Block)}

		return each

	case "Doctype":
		doctype := new(Doctype)
		doctype.Val = t.Val

		return doctype

	case "NamedBlock", "Block":
		return &Block{Nodes: p.build(t)}

	case "Comment":
		return nil

	case "Case":
		cas := new(Case)
		cas.Expr = JavaScriptExpression(t.Expr)
		cas.Block = Block{Nodes: p.build(t.Block)}

		return cas

	case "When":
		when := new(When)
		when.Expr = JavaScriptExpression(t.Expr)
		when.Block = Block{Nodes: p.build(t.Block)}

		return when

	case "MixinBlock":
		return new(MixinBlock)

	default:
		log.Printf("%#v\n", t)
		panic(errors.Errorf("Cannot parse Pug block %#v", t))
	}
}