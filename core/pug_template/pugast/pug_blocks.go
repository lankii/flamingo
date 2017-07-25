package pugast

// Render blocks
// see https://github.com/pugjs/pug-ast-spec/blob/master/parser.md
// Not complete, and some minor things have been stripped

type (
	// Node objects

	// Node is something renderable
	Node interface {
		Render(p *PugAst, depth int) (result string, isinline bool)
	}

	// Blocks

	// Block is a list of nodes
	Block struct {
		Nodes []Node
	}

	// Abstract Node Types

	// AttributedNode extends a node with attributes
	AttributedNode struct {
		Attrs           []Attribute
		AttributeBlocks []JavaScriptExpression
	}

	// BlockNode is a block (which is renderable)
	BlockNode struct {
		Block Block
	}

	// ExpressionNode is a javascript expression
	ExpressionNode struct {
		Expr JavaScriptExpression
	}

	// PlaceholderNode does nothing
	PlaceholderNode struct{}

	// ValueNode contains a value
	ValueNode struct {
		Val string
	}

	// Attribute is a key-value attribute (might have to be escaped)
	Attribute struct {
		Name       string
		Val        JavaScriptExpression
		MustEscape bool
	}

	// JavaScriptExpression is a string
	JavaScriptExpression string
	// JavaScriptIdentifier is a string
	JavaScriptIdentifier string

	// Doctypes

	// Doctype is a `<!DOCTYPE...` node
	Doctype struct {
		ValueNode
	}

	// Comments

	// CommonComment is the basic value for comments
	CommonComment struct {
		ValueNode
	}

	// Comment is single line comment
	Comment struct {
		CommonComment
	}

	// BlockComment is multi line comment
	BlockComment struct {
		BlockNode
		CommonComment
	}

	// Text

	// Text contains raw text
	Text struct {
		ValueNode
	}

	// Tag

	// CommonTag is the base structure for tags
	CommonTag struct {
		AttributedNode
		BlockNode

		SelfClosing bool  // SelfClosing if the tag is explicitly stated as self-closing
		IsInline    *bool // IsInline if the tag is defined as an inline tag as opposed to a block-level tag
	}

	// Tag is just a tag
	Tag struct {
		CommonTag
		Name string
	}

	// InterpolatedTag a tag which is created on the fly
	InterpolatedTag struct {
		CommonTag
		ExpressionNode
	}

	// Code

	// Code is a code block
	Code struct {
		BlockNode
		ValueNode

		Buffer     bool  // Buffer if the value of the piece of code is buffered in the template
		MustEscape bool  // MustEscape if the value must be HTML-escaped before being buffered
		IsInline   *bool // IsInline whether the node is the result of a string interpolation
	}

	// Code Helpers

	// Conditional is equivalent to `if`
	Conditional struct {
		Test       JavaScriptExpression
		Consequent Node
		Alternate  Node
	}

	// Each iterates over something
	Each struct {
		BlockNode

		Obj JavaScriptExpression
		Val JavaScriptIdentifier
		Key JavaScriptIdentifier
	}

	// Mixin

	// Mixin can be defined or called
	Mixin struct {
		AttributedNode
		BlockNode

		Name JavaScriptIdentifier
		Call bool
		Args string
	}

	// MixinBlock renders the mixin block
	MixinBlock struct {
		PlaceholderNode
	}

	// Case switch/case construct
	Case struct {
		BlockNode
		ExpressionNode
	}

	// When is a case in a case-construct
	When struct {
		BlockNode
		ExpressionNode
	}
)