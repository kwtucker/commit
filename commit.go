package main

func GetCommits(modifiedFiles []string) []string {
	out := []string{}
	// Parse Files
	for _, filename := range modifiedFiles {
		UnStageFile(filename)
		out = append(out, ReadFile(filename)...)
		StageFile(filename)
	}

	return out
}
