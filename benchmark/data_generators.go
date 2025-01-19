package benchmark

import "fmt"

type ProjectSize struct {
	numAreas      int
	numModules    int
	numComponents int
	numFiles      int
}

var (
	SmallSize = ProjectSize{
		numAreas:      2,
		numModules:    5,
		numComponents: 2,
		numFiles:      5,
	} // 100 nodes (80 files + 20 dirs)

	MediumSize = ProjectSize{
		numAreas:      2,
		numModules:    10,
		numComponents: 5,
		numFiles:      5,
	} // 500 nodes (400 files + 100 dirs)

	LargeSize = ProjectSize{
		numAreas:      3,
		numModules:    10,
		numComponents: 6,
		numFiles:      5,
	} // 1000 nodes (810 files + 190 dirs)

	ExtraLargeSize = ProjectSize{
		numAreas:      4,
		numModules:    15,
		numComponents: 8,
		numFiles:      10,
	} // 5000 nodes (4320 files + 680 dirs)
)

// generateTreeStructure creates an ASCII tree with the given parameters
func generateTreeStructure(size ProjectSize) string {
	tree := "monorepo\n"
	mainAreas := []string{"frontend", "backend", "infrastructure", "tools", "docs"}
	mainAreas = mainAreas[:size.numAreas] // Limit areas based on size

	for areaIndex, area := range mainAreas {
		isLast := areaIndex == len(mainAreas)-1
		prefix := "├──"
		if isLast {
			prefix = "└──"
		}

		tree += fmt.Sprintf("%s %s\n", prefix, area)

		// Add modules for each area
		for i := 0; i < size.numModules; i++ {
			subPrefix := "│   ├──"
			if isLast {
				subPrefix = "    ├──"
			}
			if i == size.numModules-1 {
				subPrefix = "│   └──"
				if isLast {
					subPrefix = "    └──"
				}
			}

			tree += fmt.Sprintf("%s module%d\n", subPrefix, i+1)

			// Add components for each module
			for j := 0; j < size.numComponents; j++ {
				componentPrefix := "│   │   ├──"
				if isLast {
					componentPrefix = "    │   ├──"
				}
				if i == size.numModules-1 {
					componentPrefix = "    │   ├──"
				}
				if j == size.numComponents-1 {
					componentPrefix = "│   │   └──"
					if isLast || i == size.numModules-1 {
						componentPrefix = "    │   └──"
					}
				}

				tree += fmt.Sprintf("%s component%d\n", componentPrefix, j+1)

				// Add files for each component
				for k := 0; k < size.numFiles; k++ {
					filePrefix := "│   │   │   ├──"
					if isLast || i == size.numModules-1 {
						filePrefix = "    │   │   ├──"
					}
					if j == size.numComponents-1 {
						filePrefix = "│   │   │   ├──"
						if isLast || i == size.numModules-1 {
							filePrefix = "    │   │   ├──"
						}
					}
					if k == size.numFiles-1 {
						filePrefix = "│   │   │   └──"
						if isLast || i == size.numModules-1 {
							filePrefix = "    │   │   └──"
						}
					}

					extension := getExtensionForArea(areaIndex)
					tree += fmt.Sprintf("%s file%d%s\n", filePrefix, k+1, extension)
				}
			}
		}
	}

	return tree
}

// generateJSONStructure creates a JSON structure with the given parameters
func generateJSONStructure(size ProjectSize) []interface{} {
	mainAreas := []string{"frontend", "backend", "infrastructure", "tools", "docs"}
	mainAreas = mainAreas[:size.numAreas]

	// Calculate total counts for the report
	totalDirs := 1 + // root
		size.numAreas + // main areas
		(size.numAreas * size.numModules) + // modules
		(size.numAreas * size.numModules * size.numComponents) // components

	totalFiles := size.numAreas * size.numModules * size.numComponents * size.numFiles

	// Generate structure
	areaContents := make([]interface{}, len(mainAreas))

	for areaIdx, areaName := range mainAreas {
		moduleContents := make([]interface{}, size.numModules)

		for moduleIdx := 0; moduleIdx < size.numModules; moduleIdx++ {
			componentContents := make([]interface{}, size.numComponents)

			for compIdx := 0; compIdx < size.numComponents; compIdx++ {
				fileContents := make([]interface{}, size.numFiles)

				extension := getExtensionForArea(areaIdx)

				for fileIdx := 0; fileIdx < size.numFiles; fileIdx++ {
					fileContents[fileIdx] = map[string]interface{}{
						"type": "file",
						"name": fmt.Sprintf("file%d%s", fileIdx+1, extension),
					}
				}

				componentContents[compIdx] = map[string]interface{}{
					"type":     "directory",
					"name":     fmt.Sprintf("component%d", compIdx+1),
					"contents": fileContents,
				}
			}

			moduleContents[moduleIdx] = map[string]interface{}{
				"type":     "directory",
				"name":     fmt.Sprintf("module%d", moduleIdx+1),
				"contents": componentContents,
			}
		}

		areaContents[areaIdx] = map[string]interface{}{
			"type":     "directory",
			"name":     areaName,
			"contents": moduleContents,
		}
	}

	return []interface{}{
		map[string]interface{}{
			"type":     "directory",
			"name":     "monorepo",
			"contents": areaContents,
		},
		map[string]interface{}{
			"type":        "report",
			"directories": totalDirs,
			"files":       totalFiles,
		},
	}
}

func getExtensionForArea(areaIndex int) string {
	switch areaIndex {
	case 0:
		return ".tsx"
	case 1:
		return ".go"
	case 2:
		return ".tf"
	case 3:
		return ".sh"
	default:
		return ".md"
	}
}
