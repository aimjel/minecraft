package packet

import _ "embed"

//go:embed internal/data/dimension.nbt
var dimensions []byte

type JoinGame struct {
	EntityID         int32
	IsHardcore       bool
	GameMode         uint8
	PreviousGameMode int8
	DimensionNames   []string

	//Dimension Codec & Dimension

	DimensionName       string
	HashedSeed          int64
	ViewDistance        int32
	ReducedDebugInfo    bool
	IsDebug             bool
	IsFlat              bool
}

func (g JoinGame) ID() int32 {
	return 0x26
}

func (g *JoinGame) Decode(r *Reader) error {
	panic("implement") //todo implement decode join game packet
	return nil
}

func (g JoinGame) Encode(w Writer) error {
	_ = w.Int32(g.EntityID)
	_ = w.Bool(g.IsHardcore)
	_ = w.Uint8(g.GameMode)
	_ = w.Int8(g.PreviousGameMode)
	_ = w.VarInt(int32(len(g.DimensionNames)))
	for _, world := range g.DimensionNames {
		_ = w.String(world)
	}
	_ = w.Nbt(dimensions)
	_ = w.String(g.DimensionName)
	_ = w.Int64(g.HashedSeed)
	_ = w.VarInt(g.MaxPlayers)
	_ = w.VarInt(g.ViewDistance)
	_ = w.Bool(g.ReducedDebugInfo)
	_ = w.Bool(g.EnableRespawnScreen)
	_ = w.Bool(g.IsDebug)
	return w.Bool(g.IsFlat)
}
