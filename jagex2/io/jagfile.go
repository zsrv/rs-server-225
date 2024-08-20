package io

import (
	"errors"
	"maps"
	"slices"
	"strings"

	"github.com/zsrv/rs-server-225/jagex2/packet"
)

func init() {
	for i := range knownNames {
		//knownHashes[i] = genHash(knownNames[i])
		nameToHash[knownNames[i]] = genHash(knownNames[i])
	}
}

func genHash(name string) uint32 {
	hash := uint32(0)
	name = strings.ToUpper(name)
	for _, v := range name {
		hash = (hash*61 + uint32(v) - 32) | 0
	}
	return hash
}

type JagQueueFile struct {
	Hash    uint32
	Name    string
	Data    []uint8
	Write   bool
	Delete  bool
	Rename  bool
	NewName string
	NewHash uint32
}

type Jagfile struct {
	Data             []uint8
	FileCount        int
	FileHash         []uint32
	FileName         []string
	FileUnpackedSize []uint32
	FilePackedSize   []uint32
	FilePos          []int
	Unpacked         bool

	FileQueue []JagQueueFile
	FileWrite [][]uint8
}

func (jf *Jagfile) Get(index int) (*packet.Packet, error) {
	if index < 0 || index >= jf.FileCount {
		return nil, errors.New("index out of range")
	}

	if jf.Data == nil {
		return nil, errors.New("data is nil")
	}

	src := jf.Data[jf.FilePos[index] : jf.FilePos[index]+int(jf.FilePackedSize[index])]
	if jf.Unpacked {
		return packet.NewPacket(src), nil
	} else {
		decompressed, err := BZip2Decompress(src, int(jf.FileUnpackedSize[index]), true, false)
		if err != nil {
			return nil, err
		}
		return packet.NewPacket(decompressed), nil
	}
}

func (jf *Jagfile) Read(name string) (*packet.Packet, error) {
	hash := genHash(name)

	for i := range jf.FileCount {
		if jf.FileHash[i] == hash {
			return jf.Get(i)
		}
	}

	return nil, errors.New("file not found")
}

func (jf *Jagfile) Write(name string, data *packet.Packet) {
	hash := genHash(name)

	jf.FileQueue = append(jf.FileQueue, JagQueueFile{
		Hash:  hash,
		Name:  name,
		Data:  data.Buf[:data.Pos],
		Write: true,
	})
}

func (jf *Jagfile) Delete(name string) {
	hash := genHash(name)

	jf.FileQueue = append(jf.FileQueue, JagQueueFile{
		Hash:   hash,
		Name:   name,
		Delete: true,
	})
}

func (jf *Jagfile) Rename(oldName string, newName string) {
	oldHash := genHash(oldName)
	newHash := genHash(newName)

	jf.FileQueue = append(jf.FileQueue, JagQueueFile{
		Hash:    oldHash,
		Name:    oldName,
		Rename:  true,
		NewName: newName,
		NewHash: newHash,
	})
}

func (jf *Jagfile) Save(path string, doNotCompressWhole bool) error {
	buf := packet.AllocPacket(5)

	for i := range jf.FileQueue {
		queued := jf.FileQueue[i]
		index := slices.Index(jf.FileHash, queued.Hash)

		if queued.Write {
			if index == -1 {
				index = jf.FileCount
				jf.FileCount++

				jf.FileHash[index] = queued.Hash
				jf.FileName[index] = queued.Name
			}

			if queued.Data == nil {
				return errors.New("data is nil")
			}

			jf.FileUnpackedSize[index] = uint32(len(queued.Data))
			jf.FilePackedSize[index] = uint32(len(queued.Data))
			jf.FilePos[index] = -1
			jf.FileWrite[index] = queued.Data
		}

		if queued.Delete && index != -1 {
			jf.FileHash = slices.Delete(jf.FileHash, index, index+1)
			jf.FileName = slices.Delete(jf.FileName, index, index+1)
			jf.FileUnpackedSize = slices.Delete(jf.FileUnpackedSize, index, index+1)
			jf.FilePackedSize = slices.Delete(jf.FilePackedSize, index, index+1)
			jf.FilePos = slices.Delete(jf.FilePos, index, index+1)
			jf.FileCount--
		}

		if queued.Rename && index != -1 {
			if queued.NewHash == 0 {
				return errors.New("new hash is zero")
			}

			if queued.NewName == "" {
				return errors.New("new name is zero")
			}

			jf.FileHash[index] = queued.NewHash
			jf.FileName[index] = queued.NewName
		}

		jf.FileQueue = slices.Delete(jf.FileQueue, i, i+1)
		i--
	}

	var compressWhole bool
	if jf.FileCount == 1 {
		compressWhole = true
	}

	if doNotCompressWhole && compressWhole {
		compressWhole = false
	}

	// write header
	buf.P2(uint16(jf.FileCount))
	for i := range jf.FileCount {
		buf.P4(jf.FileHash[i])
		buf.P3(jf.FileUnpackedSize[i])

		if jf.FileWrite[i] != nil && !compressWhole {
			var err error
			jf.FileWrite[i], err = BZip2Compress(jf.FileWrite[i], false, true, 1, 0)
			if err != nil {
				return err
			}

			jf.FilePackedSize[i] = uint32(len(jf.FileWrite[i]))
		}

		buf.P3(jf.FilePackedSize[i])
	}

	// write files
	for i := range jf.FileCount {
		data := jf.FileWrite[i]
		buf.PData(data, len(data))
	}

	jag := packet.AllocPacket(5)
	jag.P3(uint32(buf.Pos))

	if compressWhole {
		b, err := BZip2Compress(buf.Buf, false, true, 1, 0)
		if err != nil {
			return err
		}
		buf = packet.NewPacket(b)
	}

	if compressWhole {
		jag.P3(uint32(buf.Len()))
		jag.PData(buf.Bytes(), buf.Len())
	} else {
		jag.P3(uint32(buf.Pos))
		jag.PData(buf.Bytes(), buf.Pos)
	}

	if err := jag.Save(path, jag.Pos, 0); err != nil {
		return err
	}
	buf.Release()
	jag.Release()

	return nil
}

func (jf *Jagfile) Deconstruct(name string) (uint16, []int, []int, []uint32, error) {
	dat, err := jf.Read(name + ".dat")
	if err != nil {
		return 0, nil, nil, nil, err
	}

	idx, err := jf.Read(name + ".idx")
	if err != nil {
		return 0, nil, nil, nil, err
	}

	count := idx.G2()

	sizes := make([]int, count)
	offsets := make([]int, count)

	offset := 2
	for i := range count {
		sizes[i] = int(idx.G2())
		offsets[i] = offset
		offset += sizes[i]
	}

	checksums := make([]uint32, count)
	for i := range count {
		dat.Pos = offsets[i]
		checksums[i] = packet.GetCRC(dat.Bytes(), offset, sizes[i])
	}

	return count, sizes, offsets, checksums, nil
}

func NewJagfile(src *packet.Packet) (*Jagfile, error) {
	if src == nil {
		return nil, errors.New("src cannot be nil")
	}

	jf := &Jagfile{}

	unpackedSize := src.G3()
	packedSize := src.G3()

	if unpackedSize == packedSize {
		jf.Data = src.Buf
		jf.Unpacked = false
	} else {
		var err error
		jf.Data, err = BZip2Decompress(src.Buf, int(unpackedSize), true, false)
		if err != nil {
			return nil, err
		}

		src = packet.NewPacket(jf.Data)
		jf.Unpacked = true
	}

	jf.FileCount = int(src.G2())

	// initialize
	jf.FileHash = make([]uint32, jf.FileCount)
	jf.FileName = make([]string, jf.FileCount)
	jf.FileUnpackedSize = make([]uint32, jf.FileCount)
	jf.FilePackedSize = make([]uint32, jf.FileCount)
	jf.FilePos = make([]int, jf.FileCount)

	pos := uint32(src.Pos + jf.FileCount*10)
	for i := range jf.FileCount {
		jf.FileHash[i] = src.G4()

		for k, v := range maps.All(nameToHash) {
			if v == jf.FileHash[i] {
				jf.FileName[i] = k
				break
			}
		}

		jf.FileUnpackedSize[i] = src.G3()
		jf.FilePackedSize[i] = src.G3()

		jf.FilePos[i] = int(pos)
		pos += jf.FilePackedSize[i]
	}

	return jf, nil
}

func LoadJagfile(path string) (*Jagfile, error) {
	p, err := packet.Load(path, false)
	if err != nil {
		return nil, err
	}
	jf, err := NewJagfile(p)
	if err != nil {
		return nil, err
	}
	return jf, nil
}

var knownNames = []string{
	// title
	"index.dat",
	"logo.dat",
	"p11.dat",
	"p12.dat",
	"b12.dat",
	"q8.dat",
	"runes.dat",
	"title.dat",
	"titlebox.dat",
	"titlebutton.dat",
	// seen in 274
	"p11_full.dat",
	"p12_full.dat",
	"b12_full.dat",
	"q8_full.dat",

	// config
	"flo.dat",
	"flo.idx",
	"idk.dat",
	"idk.idx",
	"loc.dat",
	"loc.idx",
	"npc.dat",
	"npc.idx",
	"obj.dat",
	"obj.idx",
	"seq.dat",
	"seq.idx",
	"spotanim.dat",
	"spotanim.idx",
	"varp.dat",
	"varp.idx",
	// seen in 254
	"varbit.dat",
	"varbit.idx",
	// seen in 274
	"mesanim.dat",
	"mesanim.idx",
	"mes.dat",
	"mes.idx",
	"param.dat",
	"param.idx",
	"hunt.dat",
	"hunt.idx",

	// interface
	"data",

	// media
	"backbase1.dat",
	"backbase2.dat",
	"backhmid1.dat",
	"backhmid2.dat",
	"backleft1.dat",
	"backleft2.dat",
	"backright1.dat",
	"backright2.dat",
	"backtop1.dat",
	"backtop2.dat",
	"backvmid1.dat",
	"backvmid2.dat",
	"backvmid3.dat",
	"chatback.dat",
	"combatboxes.dat",
	"combaticons.dat",
	"combaticons2.dat",
	"combaticons3.dat",
	"compass.dat",
	"cross.dat",
	"gnomeball_buttons.dat",
	"headicons.dat",
	"hitmarks.dat",
	// index.dat
	"invback.dat",
	"leftarrow.dat",
	"magicoff.dat",
	"magicoff2.dat",
	"magicon.dat",
	"magicon2.dat",
	"mapback.dat",
	"mapdots.dat",
	"mapflag.dat",
	"mapfunction.dat",
	"mapscene.dat",
	"miscgraphics.dat",
	"miscgraphics2.dat",
	"miscgraphics3.dat",
	"prayerglow.dat",
	"prayeroff.dat",
	"prayeron.dat",
	"redstone1.dat",
	"redstone2.dat",
	"redstone3.dat",
	"rightarrow.dat",
	"scrollbar.dat",
	"sideicons.dat",
	"staticons.dat",
	"staticons2.dat",
	"steelborder.dat",
	"steelborder2.dat",
	"sworddecor.dat",
	"tradebacking.dat",
	"wornicons.dat",
	// seen in 254
	"mapmarker.dat",
	"mod_icons.dat",
	"mapedge.dat",
	// seen in 336
	"blackmark.dat",
	"button_brown.dat",
	"button_brown_big.dat",
	"button_red.dat",
	"chest.dat",
	"coins.dat",
	"headicons_hint.dat",
	"headicons_pk.dat",
	"headicons_prayer.dat",
	"key.dat",
	"keys.dat",
	"leftarrow_small.dat",
	"letter.dat",
	"number_button.dat",
	"overlay_duel.dat",
	"overlay_multiway.dat",
	"pen.dat",
	"rightarrow_small.dat",
	"startgame.dat",
	"tex_brown.dat",
	"tex_red.dat",
	"titlescroll.dat",

	// models (225 and before)
	"base_head.dat",
	"base_label.dat",
	"base_type.dat",
	"frame_del.dat",
	"frame_head.dat",
	"frame_tran1.dat",
	"frame_tran2.dat",
	"ob_axis.dat",
	"ob_face1.dat",
	"ob_face2.dat",
	"ob_face3.dat",
	"ob_face4.dat",
	"ob_face5.dat",
	"ob_head.dat",
	"ob_point1.dat",
	"ob_point2.dat",
	"ob_point3.dat",
	"ob_point4.dat",
	"ob_point5.dat",
	"ob_vertex1.dat",
	"ob_vertex2.dat",

	// versionlist (introduced in 234)
	"anim_crc",
	"anim_index",
	"anim_version",
	"map_crc",
	"map_index",
	"map_version",
	"midi_crc",
	"midi_index",
	"midi_version",
	"model_crc",
	"model_index",
	"model_version",

	// textures
	// index.dat
	"0.dat",
	"1.dat",
	"2.dat",
	"3.dat",
	"4.dat",
	"5.dat",
	"6.dat",
	"7.dat",
	"8.dat",
	"9.dat",
	"10.dat",
	"11.dat",
	"12.dat",
	"13.dat",
	"14.dat",
	"15.dat",
	"16.dat",
	"17.dat",
	"18.dat",
	"19.dat",
	"20.dat",
	"21.dat",
	"22.dat",
	"23.dat",
	"24.dat",
	"25.dat",
	"26.dat",
	"27.dat",
	"28.dat",
	"29.dat",
	"30.dat",
	"31.dat",
	"32.dat",
	"33.dat",
	"34.dat",
	"35.dat",
	"36.dat",
	"37.dat",
	"38.dat",
	"39.dat",
	"40.dat",
	"41.dat",
	"42.dat",
	"43.dat",
	"44.dat",
	"45.dat",
	"46.dat",
	"47.dat",
	"48.dat",
	"49.dat",

	// wordenc
	"badenc.txt",
	"domainenc.txt",
	"fragmentsenc.txt",
	"tldlist.txt",

	// sounds
	"sounds.dat",

	// worldmap
	"labels.dat",
	"floorcol.dat",
	"underlay.dat",
	"overlay.dat",
	"size.dat", // added later
}

var nameToHash = make(map[string]uint32, len(knownNames))
