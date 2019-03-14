package ast

import (
	"strings"
)

type DefinitionKind string

const (
	Scalar      DefinitionKind = "SCALAR"
	Object      DefinitionKind = "OBJECT"
	Interface   DefinitionKind = "INTERFACE"
	Union       DefinitionKind = "UNION"
	Enum        DefinitionKind = "ENUM"
	InputObject DefinitionKind = "INPUT_OBJECT"
)

// ObjectDefinition is the core type definition object, it includes all of the definable types
// but does *not* cover schema or directives.
//
// @vektah: Javascript implementation has different types for all of these, but they are
// more similar than different and don't define any behaviour. I think this style of
// "some hot" struct works better, at least for go.
//
// Type extensions are also represented by this same struct.
type Definition struct {
	Kind        DefinitionKind
	Description string
	Name        string
	Directives  DirectiveList
	Interfaces  []string      // object and input object
	Fields      FieldList     // object and input object
	Types       []string      // union
	EnumValues  EnumValueList // enum

	Position *Position `dump:"-"`
	BuiltIn  bool      `dump:"-"`
}

func (d *Definition) String() string {

	sb := new(strings.Builder)
	if d.Kind == Object {
		indent := 0
		sb.WriteString("type ")
		sb.WriteString(d.Name)
		sb.WriteString(" {\n")
		indent++
		for _, f := range d.Fields {
			sb.WriteString(strings.Repeat("  ", indent))
			sb.WriteString(strings.ToLower(f.Name)) // todo camel
			sb.WriteString(": ")
			sb.WriteString(f.Type.Name())
			if f.Type.NonNull {
				sb.WriteString("!")
			}
			sb.WriteString("\n")
		}
		indent--
		sb.WriteString("}\n")

		indent = 0
		sb.WriteString("type Query")
		sb.WriteString(" {\n")
		indent++

		sb.WriteString(strings.Repeat("  ", indent))
		sb.WriteString(strings.ToLower(d.Name) + "s") // HACK OF HACKS
		sb.WriteString(": ")
		sb.WriteString("[" + d.Name + "]")
		sb.WriteString("\n")
		indent--
		sb.WriteString("}\n")

	}
	return sb.String()
}

func (d *Definition) IsLeafType() bool {
	return d.Kind == Enum || d.Kind == Scalar
}

func (d *Definition) IsAbstractType() bool {
	return d.Kind == Interface || d.Kind == Union
}

func (d *Definition) IsCompositeType() bool {
	return d.Kind == Object || d.Kind == Interface || d.Kind == Union
}

func (d *Definition) IsInputType() bool {
	return d.Kind == Scalar || d.Kind == Enum || d.Kind == InputObject
}

func (d *Definition) OneOf(types ...string) bool {
	for _, t := range types {
		if d.Name == t {
			return true
		}
	}
	return false
}

type FieldDefinition struct {
	Description  string
	Name         string
	Arguments    ArgumentDefinitionList // only for objects
	DefaultValue *Value                 // only for input objects
	Type         *Type
	Directives   DirectiveList
	Position     *Position `dump:"-"`
}

type ArgumentDefinition struct {
	Description  string
	Name         string
	DefaultValue *Value
	Type         *Type
	Directives   DirectiveList
	Position     *Position `dump:"-"`
}

type EnumValueDefinition struct {
	Description string
	Name        string
	Directives  DirectiveList
	Position    *Position `dump:"-"`
}

type DirectiveDefinition struct {
	Description string
	Name        string
	Arguments   ArgumentDefinitionList
	Locations   []DirectiveLocation
	Position    *Position `dump:"-"`
}
