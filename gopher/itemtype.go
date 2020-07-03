package gopher

import "strings"

// ItemType specifies a gopher item type char
type ItemType byte

// RFC 1436 Standard item types
const (
	typeFile       = ItemType('0') /* Regular file (text) */
	typeDirectory  = ItemType('1') /* Directory (menu) */
	typeDatabase   = ItemType('2') /* CCSO flat db; other db */
	typeError      = ItemType('3') /* Error message */
	typeMacBinHex  = ItemType('4') /* Macintosh BinHex file */
	typeBinArchive = ItemType('5') /* Binary archive (zip, rar, 7zip, tar, gzip, etc), CLIENT MUST READ UNTIL TCP CLOSE */
	typeUUEncoded  = ItemType('6') /* UUEncoded archive */
	typeSearch     = ItemType('7') /* Query search engine or CGI script */
	typeTelnet     = ItemType('8') /* Telnet to: VT100 series server */
	typeBin        = ItemType('9') /* Binary file (see also, 5), CLIENT MUST READ UNTIL TCP CLOSE */
	typeTn3270     = ItemType('T') /* Telnet to: tn3270 series server */
	typeGif        = ItemType('g') /* GIF format image file (just use I) */
	typeImage      = ItemType('I') /* Any format image file */
	typeRedundant  = ItemType('+') /* Redundant (indicates mirror of previous item) */
)

// GopherII Standard item types
const (
	typeCalendar = ItemType('c') /* Calendar file */
	typeDoc      = ItemType('d') /* Word-processing document; PDF document */
	typeHTML     = ItemType('h') /* HTML document */
	typeInfo     = ItemType('i') /* Informational text (not selectable) */
	typeMarkup   = ItemType('p') /* Page layout or markup document (plain text w/ ASCII tags) */
	typeMail     = ItemType('M') /* Email repository (MBOX) */
	typeAudio    = ItemType('s') /* Audio recordings */
	typeXML      = ItemType('x') /* eXtensible Markup Language document */
	typeVideo    = ItemType(';') /* Video files */
)

// Commonly Used item types
const (
	typeTitle        = ItemType('!') /* [SERVER ONLY] Menu title (set title ONCE per gophermap) */
	typeComment      = ItemType('#') /* [SERVER ONLY] Comment, rest of line is ignored */
	typeHiddenFile   = ItemType('-') /* [SERVER ONLY] Hide file/directory from directory listing */
	typeEnd          = ItemType('.') /* [SERVER ONLY] Last line -- stop processing gophermap default */
	typeSubGophermap = ItemType('=') /* [SERVER ONLY] Include subgophermap / regular file here. */
	typeEndBeginList = ItemType('*') /* [SERVER ONLY] Last line + directory listing -- stop processing gophermap and end on directory listing */
)

// Internal item types
const (
	typeDefault       = typeBin
	typeInfoNotStated = ItemType('I')
	typeUnknown       = ItemType('?')
)

// fileExtMap specifies mapping of file extensions to gopher item types
var fileExtMap = map[string]ItemType{
	".out": typeBin,
	".a":   typeBin,
	".o":   typeBin,
	".ko":  typeBin, /* Kernel extensions... WHY ARE YOU GIVING ACCESS TO DIRECTORIES WITH THIS */

	".gophermap": typeDirectory,

	".lz":  typeBinArchive,
	".gz":  typeBinArchive,
	".bz2": typeBinArchive,
	".7z":  typeBinArchive,
	".zip": typeBinArchive,

	".gitignore":    typeFile,
	".txt":          typeFile,
	".json":         typeFile,
	".yaml":         typeFile,
	".ocaml":        typeFile,
	".s":            typeFile,
	".c":            typeFile,
	".py":           typeFile,
	".h":            typeFile,
	".go":           typeFile,
	".fs":           typeFile,
	".odin":         typeFile,
	".nanorc":       typeFile,
	".bashrc":       typeFile,
	".mkshrc":       typeFile,
	".vimrc":        typeFile,
	".vim":          typeFile,
	".viminfo":      typeFile,
	".sh":           typeFile,
	".conf":         typeFile,
	".xinitrc":      typeFile,
	".jstarrc":      typeFile,
	".joerc":        typeFile,
	".jpicorc":      typeFile,
	".profile":      typeFile,
	".bash_profile": typeFile,
	".bash_logout":  typeFile,
	".log":          typeFile,
	".ovpn":         typeFile,

	".md": typeMarkup,

	".xml": typeXML,

	".doc":  typeDoc,
	".docx": typeDoc,
	".pdf":  typeDoc,

	".jpg":  typeImage,
	".jpeg": typeImage,
	".png":  typeImage,
	".gif":  typeImage,

	".html": typeHTML,
	".htm":  typeHTML,

	".ogg":  typeAudio,
	".mp3":  typeAudio,
	".wav":  typeAudio,
	".mod":  typeAudio,
	".it":   typeAudio,
	".xm":   typeAudio,
	".mid":  typeAudio,
	".vgm":  typeAudio,
	".opus": typeAudio,
	".m4a":  typeAudio,
	".aac":  typeAudio,

	".mp4":  typeVideo,
	".mkv":  typeVideo,
	".webm": typeVideo,
	".avi":  typeVideo,
}

// getItemType is an internal function to get an ItemType for a file name string
func getItemType(name string) ItemType {
	// Split, name MUST be lower
	split := strings.Split(strings.ToLower(name), ".")

	// First we look at how many '.' in name string
	splitLen := len(split)
	switch splitLen {
	case 0:
		// Always return typeDefault, we can't tell
		return typeDefault

	default:
		// get index of str after last '.', look up in fileExtMap
		fileType, ok := fileExtMap["."+split[splitLen-1]]
		if ok {
			return fileType
		}
		return typeDefault
	}
}

// parseLineType parses a gophermap's line type based on first char and contents
func parseLineType(line string) ItemType {
	lineLen := len(line)

	if lineLen == 0 {
		return typeInfoNotStated
	}

	// Get ItemType for first char
	t := ItemType(line[0])

	if lineLen == 1 {
		// The only accepted types for length 1 line below:
		t := ItemType(line[0])
		if t == typeEnd ||
			t == typeEndBeginList ||
			t == typeComment ||
			t == typeInfo ||
			t == typeTitle {
			return t
		}
		return typeUnknown
	} else if !strings.Contains(line, "\t") {
		// The only accepted types for length >= 1 and with a tab
		if t == typeComment ||
			t == typeTitle ||
			t == typeInfo ||
			t == typeHiddenFile ||
			t == typeSubGophermap {
			return t
		}
		return typeInfoNotStated
	}

	return t
}
