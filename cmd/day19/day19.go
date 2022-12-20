package main

import (
	"fmt"
	"strings"

	"github.com/Olegas/advent-of-code-2022/internal/util"
	"github.com/Olegas/goaocd"
)

// Only part A is ok. Part B is OK for test data but for real input it gives to small result =(

type Blueprint struct {
	ore, clay int
	obsidian  [2]int
	geode     [2]int
}

type Resources struct {
	materials [4]int
	robots    [4]int
}

type State struct {
	Resources
	minutes int
}

func (s State) Clone() State {
	ret := State{minutes: s.minutes}
	copy(ret.materials[:], s.materials[:])
	copy(ret.robots[:], s.robots[:])
	return ret
}

func sample() []string {
	d := `Blueprint 1: Each ore robot costs 4 ore. Each clay robot costs 2 ore. Each obsidian robot costs 3 ore and 14 clay. Each geode robot costs 2 ore and 7 obsidian.  
Blueprint 2: Each ore robot costs 2 ore. Each clay robot costs 3 ore. Each obsidian robot costs 3 ore and 8 clay. Each geode robot costs 3 ore and 12 obsidian.`
	return strings.Split(d, "\n")
}

func partA(lines []string, maxTurns int, pt2 bool) int {
	done := goaocd.Duration("Part A")
	defer done()

	res := 0
	if pt2 {
		res = 1
	}
	bps := loadBlueprints(lines)
	for idx, bp := range bps {

		fmt.Printf("simulating blueprint %d\n", idx)
		r := simulateBlueprint(bp, maxTurns)
		fmt.Printf("Blueprint %d, result %d\n", idx, r)
		if !pt2 {
			res += r * (idx + 1)
		} else {
			res *= r
		}

	}

	return res
}

func partB() int {
	done := goaocd.Duration("Part B")
	defer done()

	return 0
}

func loadBlueprints(lines []string) []Blueprint {
	res := []Blueprint{}
	for _, line := range lines {
		var bp, ore, clay, obs_ore, obs_clay, geo_ore, geo_obs int
		n, err := fmt.Sscanf(line,
			"Blueprint %d: Each ore robot costs %d ore. Each clay robot costs %d ore. Each obsidian robot costs %d ore and %d clay. Each geode robot costs %d ore and %d obsidian.",
			&bp, &ore, &clay, &obs_ore, &obs_clay, &geo_ore, &geo_obs)
		if err != nil || n != 7 {
			panic(fmt.Sprintf("Failed to parse line %s", line))
		}
		res = append(res, Blueprint{
			ore:      ore,
			clay:     clay,
			obsidian: [2]int{obs_ore, obs_clay},
			geode:    [2]int{geo_ore, geo_obs},
		})
	}
	return res
}

func spendingVariants(s *State, b *Blueprint) [][4]int {
	res := [][4]int{}

	// geode: ore + obsidian
	if s.materials[0] >= b.geode[0] && s.materials[2] >= b.geode[1] {
		res = append(res, [4]int{0, 0, 0, 1})
		return res
	}

	// obsidian: ore + clay
	if s.materials[0] >= b.obsidian[0] && s.materials[1] >= b.obsidian[1] {
		res = append(res, [4]int{0, 0, 1, 0})
	}

	// ore: just ore
	oreRobots := s.materials[0] >= b.ore
	if oreRobots {
		res = append(res, [4]int{1, 0, 0, 0})
	}

	// clay: just ore
	clayRobots := s.materials[0] >= b.clay
	if clayRobots {
		res = append(res, [4]int{0, 1, 0, 0})
	}

	res = append(res, [4]int{0, 0, 0, 0})

	return res
}

func simulateBlueprint(b Blueprint, maxTurns int) int {
	states := []State{{Resources: Resources{robots: [4]int{1, 0, 0, 0}}}}
	maxOpenGeodes := 0
	for {
		newItemsWasAdded := false
		seenAtMinute := map[Resources]bool{}
		for ok := len(states) > 0; ok; ok = len(states) > 0 {
			s := states[0]
			states = states[1:]
			if s.minutes < maxTurns {
				variants := spendingVariants(&s, &b)
				// Multiple variants of spendings of this step
				for _, variant := range variants {
					// Next state
					newS := s.Clone()
					// collect material
					for idx, r := range s.robots {
						newS.materials[idx] += r
					}

					// get new robots
					// variant has number of new robots of each kind
					for idx, diff := range variant {
						if diff == 0 {
							continue
						}
						newS.robots[idx] += diff
						switch idx {
						case 0:
							// ore robot costs ore
							newS.materials[0] -= diff * b.ore
						case 1:
							// clay robot costs ore
							newS.materials[0] -= diff * b.clay
						case 2:
							// obsidian robot cost ore and clay
							newS.materials[0] -= diff * b.obsidian[0]
							newS.materials[1] -= diff * b.obsidian[1]
						case 3:
							// geode robot costs ore and obisdian
							newS.materials[0] -= diff * b.geode[0]
							newS.materials[2] -= diff * b.geode[1]
						}

					}
					_, ok := seenAtMinute[newS.Resources]
					if ok {
						continue
					}
					seenAtMinute[newS.Resources] = true
					newS.minutes++
					newItemsWasAdded = true
					states = append(states, newS)
				}
			} else {
				maxOpenGeodes = util.Max(maxOpenGeodes, s.materials[3])
			}
		}
		if !newItemsWasAdded {
			break
		}
	}
	return maxOpenGeodes
}

func main() {
	lines := sample()
	lines = goaocd.Lines()

	// fmt.Printf("Part A: %d\n", partA(lines, 24))
	fmt.Printf("Part A: %d\n", partA(lines[:], 32, true))
	// fmt.Printf("Part B: %d\n", partB())
}
