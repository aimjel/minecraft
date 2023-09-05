package packet

import (
	"fmt"
	"unsafe"
)

type DeclareCommands struct {
	Nodes []Node

	RootIndex int32
}

func (d *DeclareCommands) ID() int32 {
	return 0x10
}

func (d *DeclareCommands) Decode(r *Reader) error {
	var length int32
	if err := r.VarInt(&length); err != nil {
		return err
	}

	d.Nodes = make([]Node, length)
	for i := int32(0); i < length; i++ {
		n := &d.Nodes[i]

		_ = r.Uint8(&n.Flags)

		var childLength int32
		if err := r.VarInt(&childLength); err != nil {
			return err
		}

		n.Children = make([]int32, childLength)
		for j := int32(0); j < childLength; j++ {
			if err := r.VarInt(&n.Children[j]); err != nil {
				return err
			}
		}

		if n.Flags&0x08 == 0x08 {
			if err := r.VarInt(&n.RedirectNode); err != nil {
				return err
			}
		}

		switch n.Flags & 0x03 {

		//root node type
		case 0:

		//literal node type
		case 1:
			_ = r.String(&n.Name)

		//argument node type
		case 2:
			_ = r.String(&n.Name)
			_ = r.String(&n.ParserID)
			switch n.ParserID {

			case "brigadier:float":
				_ = r.Uint8(&n.Properties.Flags)
				if n.Properties.Flags&0x01 == 1 {
					_ = r.Float32((*float32)(unsafe.Pointer(&n.Properties.Min)))
				}
				if n.Properties.Flags&0x02 == 2 {
					_ = r.Float32((*float32)(unsafe.Pointer(&n.Properties.Max)))
				}

			case "brigadier:double":
				_ = r.Uint8(&n.Properties.Flags)
				if n.Properties.Flags&0x01 == 1 {
					_ = r.Float64((*float64)(unsafe.Pointer(&n.Properties.Min)))
				}
				if n.Properties.Flags&0x02 == 2 {
					_ = r.Float64((*float64)(unsafe.Pointer(&n.Properties.Max)))
				}

			case "brigadier:integer":
				_ = r.Uint8(&n.Properties.Flags)
				fmt.Printf("brigadier:integer flags %08b\n", n.Properties.Flags)
				if n.Properties.Flags&0x01 == 1 {
					_ = r.Int32((*int32)(unsafe.Pointer(&n.Properties.Min)))
				}
				if n.Properties.Flags&0x02 == 2 {
					_ = r.Int32((*int32)(unsafe.Pointer(&n.Properties.Max)))
				}

			case "brigadier:long":
				_ = r.Uint8(&n.Properties.Flags)
				if n.Properties.Flags&0x01 == 1 {
					_ = r.Int64((*int64)(unsafe.Pointer(&n.Properties.Min)))
				}
				if n.Properties.Flags&0x02 == 2 {
					_ = r.Int64((*int64)(unsafe.Pointer(&n.Properties.Max)))
				}

			case "brigadier:string":
				_ = r.Uint8(&n.Properties.Flags) //suppose to be var-int type but the max value is 2 bits

			case "minecraft:entity":
				_ = r.Uint8(&n.Properties.Flags)

			case "minecraft:score_holder":
				_ = r.Uint8(&n.Properties.Flags)
			}
		}

		if n.Flags&0x10 == 0x10 {
			_ = r.String(&n.Properties.Identifier)
		}
	}

	return r.VarInt(&d.RootIndex)
}

func (d *DeclareCommands) Encode(w Writer) error {
	_ = w.VarInt(int32(len(d.Nodes)))
	for _, n := range d.Nodes {
		_ = w.Uint8(n.Flags)

		_ = w.VarIntArray(n.Children)

		if n.Flags&0x08 == 0x08 {
			_ = w.VarInt(n.RedirectNode)
		}

		switch n.Flags & 0x03 {

		//root
		case 0:

			//literal
		case 1:
			_ = w.String(n.Name)

		case 2:
			_ = w.String(n.Name)
			_ = w.String(n.ParserID)
			switch n.ParserID {

			case "brigadier:float":
				_ = w.Uint8(n.Properties.Flags)
				if n.Properties.Flags&0x01 == 1 {
					_ = w.Float32(*(*float32)(unsafe.Pointer(&n.Properties.Min)))
				}
				if n.Properties.Flags&0x02 == 2 {
					_ = w.Float32(*(*float32)(unsafe.Pointer(&n.Properties.Max)))
				}

			case "brigadier:double":
				_ = w.Uint8(n.Properties.Flags)
				if n.Properties.Flags&0x01 == 1 {
					_ = w.Float64(*(*float64)(unsafe.Pointer(&n.Properties.Min)))
				}
				if n.Properties.Flags&0x02 == 2 {
					_ = w.Float64(*(*float64)(unsafe.Pointer(&n.Properties.Max)))
				}

			case "brigadier:integer":
				_ = w.Uint8(n.Properties.Flags)
				fmt.Printf("brigadier:integer flags %08b\n", n.Properties.Flags)
				if n.Properties.Flags&0x01 == 1 {
					_ = w.Int32(*(*int32)(unsafe.Pointer(&n.Properties.Min)))
				}
				if n.Properties.Flags&0x02 == 2 {
					_ = w.Int32(*(*int32)(unsafe.Pointer(&n.Properties.Max)))
				}

			case "brigadier:long":
				_ = w.Uint8(n.Properties.Flags)
				if n.Properties.Flags&0x01 == 1 {
					_ = w.Int64(*(*int64)(unsafe.Pointer(&n.Properties.Min)))
				}
				if n.Properties.Flags&0x02 == 2 {
					_ = w.Int64(*(*int64)(unsafe.Pointer(&n.Properties.Max)))
				}

			case "brigadier:string":
				_ = w.Uint8(n.Properties.Flags) //suppose to be var-int type but the max value is 2 bits

			case "minecraft:entity":
				_ = w.Uint8(n.Properties.Flags)

			case "minecraft:score_holder":
				_ = w.Uint8(n.Properties.Flags)
			}
		}

		if n.Flags&0x10 == 0x10 {
			_ = w.String(n.Properties.Identifier)
		}
	}

	return w.VarInt(d.RootIndex)
}

type Node struct {
	Flags        uint8
	Children     []int32 `json:",omitempty"`
	RedirectNode int32   `json:",omitempty"`
	Name         string  `json:",omitempty"`
	ParserID     string  `json:",omitempty"`
	Properties   struct {
		Flags      uint8  `json:",omitempty"`
		Min, Max   uint64 `json:",omitempty"` //min and max stores the bits of the actual data type
		Identifier string `json:",omitempty"`
	}
	SuggestionsType string `json:",omitempty"`
}
