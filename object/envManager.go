package object

type EnvironmentManager struct {
	environments map[string]*Environment // A map of environments for each file
}

func NewEnvironmentManager() *EnvironmentManager {
	return &EnvironmentManager{
		environments: make(map[string]*Environment),
	}
}

func (em *EnvironmentManager) Get(fileName string) (*Environment, bool) {
	env, ok := em.environments[fileName]
	if !ok {
		return em.createFileEnv(fileName), false
	}

	return env, true
}

func (em *EnvironmentManager) createFileEnv(fileName string) *Environment {
	env := NewEnvironment()
	em.environments[fileName] = env

	return env
}
