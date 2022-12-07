package main

import (
	"container/list"
	"fmt"
	"math"
	"strings"

	"github.com/Olegas/advent-of-code-2022/internal/util"
	"github.com/Olegas/goaocd"
)

const (
	typeFile = "file"
	typeDir  = "dir"
)

type File struct {
	Name     string
	Size     int
	Type     string
	Parent   *File
	Children *list.List
}

func sample() []string {
	data := `$ cd /
$ ls
dir a
14848514 b.txt
8504156 c.dat
dir d
$ cd a
$ ls
dir e
29116 f
2557 g
62596 h.lst
$ cd e
$ ls
584 i
$ cd ..
$ cd ..
$ cd d
$ ls
4060174 j
8033020 d.log
5626152 d.ext
7214296 k`
	return strings.Split(strings.TrimRight(data, "\n"), "\n")
}

func loadStructure() *File {
	var root = File{Name: "/", Type: typeDir, Children: list.New()}
	lines := goaocd.Lines()
	currentDir := &root
	waitForDirEntries := false
	for _, line := range lines {
		if line[0] == '$' {
			line := line[2:]
			// Command terminates waiting for dir entries
			waitForDirEntries = false
			cmdAndArgs := strings.Split(line, " ")
			switch cmdAndArgs[0] {
			case "cd":
				if cmdAndArgs[1] == ".." {
					currentDir = currentDir.Parent
				} else if cmdAndArgs[1] == "/" {
					currentDir = &root
				} else {
					elt := currentDir.Children.Front()
					done := false
					for {
						item := elt.Value.(*File)
						if item.Type == typeDir && item.Name == cmdAndArgs[1] {
							currentDir = item
							done = true
							break
						}
						elt = elt.Next()
						if elt == nil {
							break
						}
					}
					if !done {
						panic(fmt.Sprintf("Dir %s not found", cmdAndArgs[1]))
					}
				}
			case "ls":
				waitForDirEntries = true
			}
		} else {
			if !waitForDirEntries {
				panic("Unexpected state. Dir entry but waiting for command")
			}
			// Dir entry
			if util.IsDigit(line[0]) {
				var name string
				var size int
				n, err := fmt.Sscanf(line, "%d %s", &size, &name)
				if err != nil {
					panic(err)
				}
				if n != 2 {
					panic(fmt.Sprintf("Not enought data to parse %s", line))
				}
				elt := &File{Name: name, Size: size, Type: typeFile, Parent: currentDir}
				currentDir.Children.PushBack(elt)
			} else {
				var name string
				n, err := fmt.Sscanf(line, "dir %s", &name)
				if err != nil {
					panic(err)
				}
				if n != 1 {
					panic(fmt.Sprintf("Not enought data to parse %s", line))
				}
				elt := &File{Name: name, Type: typeDir, Parent: currentDir, Children: list.New()}
				currentDir.Children.PushBack(elt)
			}
		}
	}
	return &root
}

func countSize(dir *File) {
	elt := dir.Children.Front()
	for ok := elt != nil; ok; ok = elt != nil {
		file := elt.Value.(*File)
		if file.Type == typeDir {
			countSize(file)
		}
		dir.Size += file.Size
		elt = elt.Next()
	}
}

func traverse(from *File, visitor func(node *File)) {
	visitor(from)
	if from.Children != nil {
		elt := from.Children.Front()
		for ok := elt != nil; ok; ok = elt != nil {
			item := elt.Value.(*File)
			traverse(item, visitor)
			elt = elt.Next()
		}
	}
}

func solveA(root *File) int {
	accu := 0
	traverse(root, func(node *File) {
		if node.Type == typeDir {
			if node.Size <= 100000 {
				accu += node.Size
			}
		}
	})
	return accu
}

func solveB(root *File) int {
	totalDiskSize := 70000000
	freeSpaceNeeded := 30000000
	spaceAvail := totalDiskSize - root.Size
	accu := math.MaxUint32
	traverse(root, func(node *File) {
		if node.Type == typeDir && spaceAvail+node.Size > freeSpaceNeeded && node.Size < accu {
			accu = node.Size
		}
	})
	return accu
}

func main() {
	root := loadStructure()
	countSize(root)
	fmt.Printf("Root size: %d\n", root.Size)
	fmt.Printf("Part A: %d\n", solveA(root))
	fmt.Printf("Part B: %d\n", solveB(root))
}
