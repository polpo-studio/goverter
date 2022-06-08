package builder

import (
	"github.com/dave/jennifer/jen"
	"github.com/polpo-studio/goverter/xtype"
)

// List handles array / slice types.
type List struct{}

// Matches returns true, if the builder can create handle the given types.
func (*List) Matches(source, target *xtype.Type) bool {
	return source.List && target.List && !target.ListFixed
}

// Build creates conversion source code for the given source and target type.
func (*List) Build(gen Generator, ctx *MethodContext, sourceID *xtype.JenID, source, target *xtype.Type) ([]jen.Code, *xtype.JenID, *Error) {
	targetSlice := ctx.Name(target.ID())
	index := ctx.Index()

	indexedSource := xtype.VariableID(sourceID.Code.Clone().Index(jen.Id(index)))

	newStmt, newID, err := gen.Build(ctx, indexedSource, source.ListInner, target.ListInner)
	if err != nil {
		return nil, nil, err.Lift(&Path{
			SourceID:   "[]",
			SourceType: source.ListInner.T.String(),
			TargetID:   "[]",
			TargetType: target.ListInner.T.String(),
		})
	}
	newStmt = append(newStmt, jen.Id(targetSlice).Index(jen.Id(index)).Op("=").Add(newID.Code))

	stmt := []jen.Code{
		jen.Id(targetSlice).Op(":=").Make(target.TypeAsJen(), jen.Len(sourceID.Code.Clone())),
		jen.For(jen.Id(index).Op(":=").Lit(0), jen.Id(index).Op("<").Len(sourceID.Code.Clone()), jen.Id(index).Op("++")).
			Block(newStmt...),
	}

	return stmt, xtype.VariableID(jen.Id(targetSlice)), nil
}
