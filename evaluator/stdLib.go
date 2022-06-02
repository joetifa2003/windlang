package evaluator

var initiatedLibraries map[string]*Environment

func GetStdlib(filePath string) (*Environment, bool) {
	switch filePath {
	case "math":
		return getLibrary("math", stdLibMath), true

	case "request":
		return getLibrary("request", stdLibReq), true
	}

	return nil, false
}

func getLibrary(name string, initiator func() *Environment) *Environment {
	lib, ok := initiatedLibraries[name]
	if !ok {
		lib = initiator()
		initiatedLibraries[name] = lib
	}

	return lib
}

func init() {
	initiatedLibraries = make(map[string]*Environment)
}
