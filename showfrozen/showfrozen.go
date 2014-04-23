package showfrozen

// Copyright 2013-2014 Vubeology, Inc.

import "fmt"

import "github.com/vube/depman/dep"
import "github.com/vube/depman/util"
import "github.com/vube/depman/colors"

//Read - get top-level frozen dependencies
func Read(deps dep.DependencyMap) (result string) {
	var err error
	var resultMap = make(map[string]*dep.Dependency)

	util.Print(colors.Yellow("NOTE: This will not reflect the state of the remote unless you have just run `depman install`."))

	for k, v := range deps.Map {
		if v.Type == dep.TypeGitClone && v.Alias == "" {
			util.PrintIndent(colors.Red("Error: Repo '" + k + "' Type '" + v.Type + "' requires 'alias' field (defined in " + deps.Path + ")"))
			continue
		}

		v.Version, err = v.VCS.GetHead(v)
		if err != nil {
			util.Fatal(err)
		}

		resultMap[k] = v
	}

	//not changing the logic in the loop because we might want to change the print format later
	for _, v := range resultMap {
		result += fmt.Sprintf("%s %s\n", v.Repo, v.Version)
	}

	return
}

//ReadRecursively - get frozen dependencies recursively
func ReadRecursively(deps dep.DependencyMap, set map[string]string) (result string) {
	var err error

	if set == nil {
		util.Print(colors.Yellow("NOTE: This will not reflect the state of the remote unless you have just run `depman install`."))

		set = make(map[string]string)
	}

	for name, d := range deps.Map {
		var subPath string
		var depsFile string
		var subDeps dep.DependencyMap

		if _, ok := set[d.Repo]; ok {
			continue
		}

		if d.Type == dep.TypeGitClone && d.Alias == "" {
			util.PrintIndent(colors.Red("Error: Repo '" + name + "' Type '" + d.Type + "' requires 'alias' field (defined in " + deps.Path + ")"))
			continue
		}

		{
			var temp string

			temp, err = d.VCS.GetHead(d)
			if err != nil {
				util.Fatal(err)
			}

			set[d.Repo] = temp
			result += fmt.Sprintf("%s %s\n", d.Repo, temp)
		}

		subPath = d.Path()

		// Recursive
		depsFile = util.UpwardFind(subPath, dep.DepsFile)
		if depsFile != "" {
			subDeps, err = dep.Read(depsFile)
			if err == nil {
				result += ReadRecursively(subDeps, set)
			} else {
				util.Print(colors.Yellow("Error reading deps from '" + subDeps.Path + "': " + err.Error()))
			}
		}
	}
	return
}
