package fonts

type Font8x8 struct {
	Height int
	Width  int
	Bitmap [][]byte
}

func (f *Font8x8) Load() {

	f.Height = 8
	f.Width = 8
	f.Bitmap = [][]byte{
		{
			0x00,
			0x00,
			0x00,
			0x00,
			0x00,
			0x00,
			0x00,
			0x00,
		},
		{
			0x00,
			0x00,
			0x00,
			0x00,
			0x00,
			0x00,
			0x00,
			0x00,
		},
	}
}
