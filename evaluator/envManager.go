package evaluator

type EnvironmentManager struct {
	environments map[string]*Environment // A map of environments for each file
}

func NewEnvironmentManager() *EnvironmentManager {
	return &EnvironmentManager{
		environments: make(map[string]*Environment),
	}
}

// Get returns the environment for the given file name. and whither it's evaluated or not
func (em *EnvironmentManager) Get(fileName string) (*Environment, bool) {
	env, ok := GetStdlib(fileName)
	if ok {
		em.environments[fileName] = env
		return env, true
	}

	env, ok = em.environments[fileName]
	if ok {
		return env, true
	}

	return em.createFileEnv(fileName), false
}

func (em *EnvironmentManager) createFileEnv(fileName string) *Environment {
	env := NewEnvironment()
	em.environments[fileName] = env

	return env
}
