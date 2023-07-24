package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"reflect"

	"github.com/alecthomas/kingpin/v2"
	"gopkg.in/yaml.v3"
)

var (
	path1 = kingpin.Arg("file1", "first file to compare").Required().String()
	path2 = kingpin.Arg("file2", "second file to compare").Required().String()
)

func readYaml(file string) []map[string]interface{} {
	dynamic := []map[string]interface{}{}
	cfgdata, err := os.ReadFile(file)
	if err != nil {
		panic(err)
	}

	i := 0
	dec := yaml.NewDecoder(bytes.NewReader(cfgdata))
	for {
		d := make(map[string]interface{})
		err = dec.Decode(&d)
		if err == io.EOF {
			break
		} else if err != nil {
			if len(dynamic) > 0 {
				panic(fmt.Errorf("an error occured in the %d place of %s", i, file))
			} else {
				panic(err)
			}
		}

		if !reflect.ValueOf(d).IsZero() {
			dynamic = append(dynamic, d)
			i++
		}
	}
	return dynamic
}

func main() {
	kingpin.Version("v0.1.1")
	kingpin.Parse()
	data1 := readYaml(*path1)
	data2 := readYaml(*path2)

	if len(data1) > len(data2) {
		for i := range data1 {
			if i < len(data2) {
				diff := getDifferences(data1[i], data2[i], *path1, *path2)
				d, _ := yaml.Marshal(diff)
				fmt.Println(string(d))
				if i < len(data1)-1 {
					fmt.Println("---")
				}
			} else {
				diff := getDifferences(data1[i], nil, *path1, *path2)
				d, _ := yaml.Marshal(diff)
				fmt.Println(string(d))
				if i < len(data1)-1 {
					fmt.Println("---")
				}
			}
		}
	} else if len(data2) > len(data1) {
		for i := range data2 {
			if i < len(data1) {
				diff := getDifferences(data1[i], data2[i], *path1, *path2)
				d, _ := yaml.Marshal(diff)
				fmt.Println(string(d))
				if i < len(data1)-1 {
					fmt.Println("---")
				}
			} else {
				diff := getDifferences(nil, data2[i], *path1, *path2)
				d, _ := yaml.Marshal(diff)
				fmt.Println(string(d))
				if i < len(data1)-1 {
					fmt.Println("---")
				}
			}
		}
	} else {
		for i := range data1 {
			diff := getDifferences(data1[i], data2[i], *path1, *path2)
			d, _ := yaml.Marshal(diff)
			fmt.Println(string(d))
			if i < len(data1)-1 {
				fmt.Println("---")
			}
		}
	}

}

func getDifferences(data1, data2 map[string]interface{}, filename1, filename2 string) map[string]interface{} {
	differences := make(map[string]interface{})

	// Check keys in data1 that are not present in data2
	for key, val1 := range data1 {
		val2, exists := data2[key]

		if !exists {
			continue
		} else {
			// Recursive check for nested maps
			if val1Map, ok1 := val1.(map[string]interface{}); ok1 {
				if val2Map, ok2 := val2.(map[string]interface{}); ok2 {
					nestedDiff := getDifferences(val1Map, val2Map, filename1, filename2)
					if len(nestedDiff) > 0 {
						differences[key] = nestedDiff
					}
				} else if val1 != val2 {
					// Different value for the same key
					differences[key] = val2
				}
			} else if val1 != val2 {
				// Different value for the same key
				differences[key] = val2
			}
		}
	}

	// Check keys in data2 that are not present in data1
	for key, val2 := range data2 {
		if _, exists := data1[key]; !exists {
			differences[key] = val2
		}
	}

	return differences
}
