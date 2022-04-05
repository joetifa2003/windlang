package evaluator

func GetStdlib(filePath string) (*Environment, bool) {
	switch filePath {
	case "stdlib/strings":
		return StdlibStrings(), true
	}

	return nil, false
}
