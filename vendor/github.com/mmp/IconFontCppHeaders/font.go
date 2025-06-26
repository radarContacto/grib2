package IconFontCppHeaders

type Font struct {
	Filenames       [][2]string
	Min, Max16, Max int
	Icons           map[string]string
}
