package tar2hdrs

import (
	"archive/tar"
	"errors"
	"io"
	"io/fs"
	"iter"
	"strings"
	"time"
)

var ErrInvalidTypeFlag error = errors.New("invalid type flag")

type TypeFlag byte

const (
	// A regular file.
	TypeFlagReg = tar.TypeReg

	// Header-only flags.
	TypeFlagLink    = tar.TypeLink
	TypeFlagSymlink = tar.TypeSymlink
	TypeFlagChar    = tar.TypeChar
	TypeFlagBlock   = tar.TypeBlock
	TypeFlagDir     = tar.TypeDir
	TypeFlagFifo    = tar.TypeFifo

	// Reserved.
	TypeFlagCont = tar.TypeCont

	// Non-global key-value records for PAX.
	TypeFlagXHeader = tar.TypeXHeader

	// Global key-value records for PAX.
	TypeFlagXGlobalHeader = tar.TypeXGlobalHeader

	// A sparse file.
	TypeFlagGNUSparse = tar.TypeGNUSparse

	// A long path.
	TypeFlagGNULongName = tar.TypeGNULongName

	// A long link name.
	TypeFlagGNULongLink = tar.TypeGNULongLink
)

var TypeFlagMap map[byte]TypeFlag = map[byte]TypeFlag{
	tar.TypeReg: TypeFlagReg,

	tar.TypeLink:    TypeFlagLink,
	tar.TypeSymlink: TypeFlagSymlink,
	tar.TypeChar:    TypeFlagChar,
	tar.TypeBlock:   TypeFlagBlock,
	tar.TypeDir:     TypeFlagDir,
	tar.TypeFifo:    TypeFlagFifo,

	tar.TypeCont: TypeFlagCont,

	tar.TypeXHeader:       TypeFlagXHeader,
	tar.TypeXGlobalHeader: TypeFlagXGlobalHeader,

	tar.TypeGNUSparse:   TypeFlagGNUSparse,
	tar.TypeGNULongName: TypeFlagGNULongName,
	tar.TypeGNULongLink: TypeFlagGNULongLink,
}

var TypeFlagDisplayMap map[TypeFlag]string = map[TypeFlag]string{
	TypeFlagReg: "Regular File",

	TypeFlagLink:    "Hard Link",
	TypeFlagSymlink: "Symbolic Link",
	TypeFlagChar:    "Character Device",
	TypeFlagBlock:   "Block Device",
	TypeFlagDir:     "Directory",
	TypeFlagFifo:    "FIFO",

	TypeFlagCont: "(reserved)",

	TypeFlagXHeader:       "Non-global Key/Val records(PAX)",
	TypeFlagXGlobalHeader: "Global Key/Val records(PAX)",

	TypeFlagGNUSparse:   "Sparse File",
	TypeFlagGNULongName: "Long Path",
	TypeFlagGNULongLink: "Long Link Name",
}

func (f TypeFlag) String() string { return TypeFlagDisplayMap[f] }

func TypeFlagFromByte(b byte) (TypeFlag, error) {
	val, found := TypeFlagMap[b]
	switch found {
	case true:
		return val, nil
	default:
		return 0, ErrInvalidTypeFlag
	}
}

type Format int

const (
	FormatUnknown Format = Format(tar.FormatUnknown)

	FormatUSTAR Format = Format(tar.FormatUSTAR)
	FormatPAX   Format = Format(tar.FormatPAX)
	FormatGNU   Format = Format(tar.FormatGNU)
)

var FormatMap map[tar.Format]Format = map[tar.Format]Format{
	tar.FormatUnknown: FormatUnknown,

	tar.FormatUSTAR: FormatUSTAR,
	tar.FormatPAX:   FormatPAX,
	tar.FormatGNU:   FormatGNU,
}

var FormatDisplayMap map[Format]string = map[Format]string{
	FormatUnknown: "Unknown",

	FormatUSTAR: "USTAR",
	FormatPAX:   "PAX",
	FormatGNU:   "GNU",
}

func (f Format) String() string { return FormatDisplayMap[f] }

func FormatNew(f tar.Format) Format {
	val, found := FormatMap[f]
	switch found {
	case true:
		return val
	default:
		return FormatUnknown
	}
}

type Header struct {
	TypeFlag   `json:"type_flag"`
	TypeString string `json:"type_string"`

	Name     string `json:"name"`
	Linkname string `json:"link_name"`
	Basename string `json:"basename"`
	Hidden   bool   `json:"hidden"`

	Size  int64  `json:"size"`
	Mode  int64  `json:"mode"`
	Uid   int    `json:"uid"`
	Gid   int    `json:"gid"`
	Uname string `json:"uname"`
	Gname string `json:"gname"`

	fs.FileMode    `json:"file_mode"`
	FileModeString string `json:"file_mode_string"`

	ModTime    time.Time `json:"mod_time"`
	AccessTime time.Time `json:"access_time"`
	ChangeTime time.Time `json:"change_time"`

	Devmajor int64 `json:"dev_major"`
	Devminor int64 `json:"dev_minor"`

	PAXRecords map[string]string `json:"pax_records"`

	Format `json:"format"`
}

func HeaderNew(t *tar.Header) (Header, error) {
	var ret Header

	typ, e := TypeFlagFromByte(t.Typeflag)
	if nil != e {
		return ret, e
	}
	ret.TypeFlag = typ

	ret.TypeString = typ.String()

	var finfo fs.FileInfo = t.FileInfo()
	var fmode fs.FileMode = finfo.Mode()

	ret.FileMode = fmode
	ret.FileModeString = fmode.String()

	ret.Name = t.Name
	ret.Linkname = t.Linkname
	ret.Basename = finfo.Name()

	ret.Hidden = strings.HasPrefix(ret.Basename, ".")

	ret.Size = t.Size
	ret.Mode = t.Mode
	ret.Uid = t.Uid
	ret.Gid = t.Gid
	ret.Uname = t.Uname
	ret.Gname = t.Gname

	ret.ModTime = t.ModTime
	ret.AccessTime = t.AccessTime
	ret.ChangeTime = t.ChangeTime

	ret.Devmajor = t.Devmajor
	ret.Devminor = t.Devminor

	ret.PAXRecords = t.PAXRecords

	ret.Format = FormatNew(t.Format)

	return ret, nil
}

func TarReaderToHeaders(r *tar.Reader) iter.Seq2[Header, error] {
	return func(yield func(Header, error) bool) {
		var empty Header

		for {
			hdr, e := r.Next()
			if e == io.EOF {
				return
			}

			if nil != e {
				yield(empty, e)
				return
			}

			converted, e := HeaderNew(hdr)
			if !yield(converted, e) {
				return
			}
		}
	}
}

func ReaderToHeaders(rdr io.Reader) iter.Seq2[Header, error] {
	var tr *tar.Reader = tar.NewReader(rdr)
	return TarReaderToHeaders(tr)
}
