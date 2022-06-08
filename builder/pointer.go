package builder

import (
	"github.com/dave/jennifer/jen"
	"github.com/polpo-studio/goverter/xtype"
)

// Pointer handles pointer types.
type Pointer struct{}

// Matches returns true, if the builder can create handle the given types.
func (*Pointer) Matches(source, target *xtype.Type) bool {
	return source.Pointer && target.Pointer
}

// Build creates conversion source code for the given source and target type.
func (*Pointer) Build(gen Generator, ctx *MethodContext, sourceID *xtype.JenID, source, target *xtype.Type) ([]jen.Code, *xtype.JenID, *Error) {
	ctx.PointerChange = true

	outerVar := ctx.Name(target.ID())

	nextBlock, id, err := gen.Build(ctx, xtype.OtherID(jen.Op("*").Add(sourceID.Code.Clone())), source.PointerInner, target.PointerInner)
	if err != nil {
		return nil, nil, err.Lift(&Path{
			SourceID:   "*",
			SourceType: source.PointerInner.T.String(),
			TargetID:   "*",
			TargetType: target.PointerInner.T.String(),
		})
	}

	ifBlock := nextBlock
	if id.Variable {
		ifBlock = append(ifBlock, jen.Id(outerVar).Op("=").Op("&").Add(id.Code.Clone()))
	} else {
		tempName := ctx.Name(target.PointerInner.ID())
		ifBlock = append(ifBlock, jen.Id(tempName).Op(":=").Add(id.Code.Clone()))
		ifBlock = append(ifBlock, jen.Id(outerVar).Op("=").Op("&").Id(tempName))
	}

	stmt := []jen.Code{
		jen.Var().Id(outerVar).Add(target.TypeAsJen()),
		jen.If(sourceID.Code.Clone().Op("!=").Nil()).Block(ifBlock...),
	}
	return stmt, xtype.VariableID(jen.Id(outerVar)), err
}

// TargetPointer handles type were only the target is a pointer.
type TargetPointer struct{}

// Matches returns true, if the builder can create handle the given types.
func (*TargetPointer) Matches(source, target *xtype.Type) bool {
	return !source.Pointer && target.Pointer
}

// Build creates conversion source code for the given source and target type.
func (*TargetPointer) Build(gen Generator, ctx *MethodContext, sourceID *xtype.JenID, source, target *xtype.Type) ([]jen.Code, *xtype.JenID, *Error) {
	ctx.PointerChange = true

	stmt, id, err := gen.Build(ctx, sourceID, source, target.PointerInner)
	if err != nil {
		return nil, nil, err.Lift(&Path{
			SourceID:   "*",
			SourceType: source.T.String(),
			TargetID:   "*",
			TargetType: target.PointerInner.T.String(),
		})
	}
	if id.Variable {
		return stmt, xtype.OtherID(jen.Op("&").Add(id.Code)), nil
	}
	tempName := ctx.Name(target.PointerInner.ID())
	stmt = append(stmt, jen.Id(tempName).Op(":=").Add(id.Code))
	return stmt, xtype.OtherID(jen.Op("&").Id(tempName)), nil
}
