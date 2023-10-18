package packet

import (
	"fmt"
	"github.com/aimjel/minecraft/nbt"
	"io"
	"math"
)

type Writer struct {
	wr interface {
		io.Writer
		io.ByteWriter
	}
}

func NewWriter(w interface {
	io.Writer
	io.ByteWriter
}) Writer {
	return Writer{wr: w}
}

func (w *Writer) Bool(x bool) error {
	if x {
		return w.wr.WriteByte(0x01)
	}

	return w.wr.WriteByte(0x00)
}

func (w *Writer) Uint8(x uint8) error {
	return w.wr.WriteByte(x)
}

func (w *Writer) Int8(x int8) error {
	return w.wr.WriteByte(byte(x))
}

func (w *Writer) Int16(x int16) error {
	return w.Uint16(uint16(x))
}

func (w *Writer) Uint16(x uint16) error {
	if err := w.wr.WriteByte(byte(x >> 8)); err != nil {
		return err
	}

	return w.wr.WriteByte(byte(x))
}

func (w *Writer) Int32(x int32) error {
	_, err := w.wr.Write([]byte{byte(x >> 24), byte(x >> 16), byte(x >> 8), byte(x)})
	return err
}

func (w *Writer) Int64(x int64) error {
	_, err := w.wr.Write([]byte{byte(x >> 56), byte(x >> 48), byte(x >> 40), byte(x >> 32), byte(x >> 24), byte(x >> 16), byte(x >> 8), byte(x)})
	return err
}

func (w *Writer) Float32(x float32) error {
	f := math.Float32bits(x)
	_, err := w.wr.Write([]byte{byte(f >> 24), byte(f >> 16), byte(f >> 8), byte(f)})
	return err
}

func (w *Writer) Float64(x float64) error {
	f := math.Float64bits(x)
	_, err := w.wr.Write([]byte{byte(f >> 56), byte(f >> 48), byte(f >> 40), byte(f >> 32), byte(f >> 24), byte(f >> 16), byte(f >> 8), byte(f)})
	return err
}

func (w *Writer) VarInt(x int32) error {
	ux := uint32(x)

	for ux >= 0x80 {

		if err := w.wr.WriteByte(byte(ux&0x7F) | 0x80); err != nil {
			return err
		}

		ux >>= 7
	}

	return w.wr.WriteByte(byte(ux))
}

func (w *Writer) String(x string) error {
	if err := w.VarInt(int32(len(x))); err != nil {
		return fmt.Errorf("%v wrintng string length", err)
	}

	_, err := w.wr.Write([]byte(x))
	return err
}

func (w *Writer) ByteArray(x []byte) error {
	if err := w.VarInt(int32(len(x))); err != nil {
		fmt.Println(len(x))
		return fmt.Errorf("%v wrintng byte array length", err)
	}

	_, err := w.wr.Write(x)
	return err
}

func (w *Writer) FixedByteArray(x []byte) error {
	_, err := w.wr.Write(x)
	return err
}

func (w *Writer) VarIntArray(x []int32) error {
	if err := w.VarInt(int32(len(x))); err != nil {
		return fmt.Errorf("%v wrintng varint32 array length", err)
	}

	for _, v := range x {

		if err := w.VarInt(v); err != nil {
			return fmt.Errorf("%v writng varint32 array value", err)
		}
	}

	return nil
}

func (w *Writer) Int64Array(x []int64) error {
	if err := w.VarInt(int32(len(x))); err != nil {
		return fmt.Errorf("%v wrintng varint32 array length", err)
	}

	for _, v := range x {

		if err := w.Int64(v); err != nil {
			return fmt.Errorf("%v writng int64 array value", err)
		}
	}

	return nil
}

func (w *Writer) UUID(x [16]byte) error {
	_, err := w.wr.Write(x[:])
	return err
}

func (w *Writer) Nbt(x []byte) error {
	_, err := w.wr.Write(x)
	return err
}

func (w *Writer) Nbt2(x any) error {
	return nbt.NewEncoder(w.wr).Encode(x)
}
