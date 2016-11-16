package ctags

import "arduino.cc/builder/types"

type CTags struct{}

func (c *CTags) Run(ctx *types.Context) error {
	parser := &CTagsParser{}
	ctx.CTagsOfPreprocessedSource = parser.Parse(ctx.CTagsOutput)
	protos, line := parser.GeneratePrototypes()
	if line != -1 {
		ctx.PrototypesLineWhereToInsert = line
	}
	ctx.Prototypes = protos
	return nil
}
