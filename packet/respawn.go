package packet

type Respawn struct {
  DimensionType       string
	DimensionName       string
  HashedSeed          int64
  GameMode            uint8
	PreviousGameMode    int8
  IsDebug             bool
	IsFlat              bool
  DataKept            uint8
  DeathDimensionName  string
  DeathLocation       uint64
  PartialCooldown     int32
}

func (g Respawn) ID() int32 {
	return 0x41
}

func (g *Respawn) Decode(r *Reader) error {
	panic("implement") //todo implement decode join game packet
	return nil
}

func (g Respawn) Encode(w Writer) error {
  w.String(g.DimensionType)
  w.String(g.DimensionName)
  w.Int64(g.HashedSeed)
  w.Uint8(g.GameMode)
  w.Int8(g.PreviousGameMode)
  w.Bool(g.IsDebug)
  w.Bool(g.IsFlat)
  w.Uint8(g.DataKept)
  if g.DeathDimensionName != "" {
    w.Bool(true)
    w.String(g.DeathDimensionName)
    w.Uint64(g.DeathLocation)
  } else {
    w.Bool(false)
  }
  return w.VarInt(g.PartialCooldown)
}
