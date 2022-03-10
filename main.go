//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
)

var version = "test"

// HelloGet is an HTTP Cloud Function.
func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Yum Get RepoMD,  Version: %s\n\nUsage: %s [options...]\n\n", version, os.Args[0])
		flag.PrintDefaults()
	}

	var inRepoPath = flag.String("repo", "/7/os/x86_64", "Which package to get")
	var mirrorList = flag.String("mirrors", "mirrorlist.txt", "Mirror / directory list of prefixes to use")
	var outputPath = flag.String("output", ".", "Path to put the repodata files")
	flag.Parse()

	mirrors := readMirrors(*mirrorList)
	repoPath := strings.TrimSuffix(strings.TrimPrefix(*inRepoPath, "/"), "/")

	var latestRepomd Repomd
	var latestRepomdTime int
	for i, m := range mirrors {
		repoPath := m + "/" + repoPath + "/"
		repomdPath := repoPath + "repodata/repomd.xml"
		log.Println(i, "Fetching", repomdPath)

		dat := readRepomdFile(repomdPath)
		if dat != nil {
			for _, elem := range dat.Data {
				if elem.Timestamp > latestRepomdTime {
					dat.path = repoPath
					dat.mirror = m
					latestRepomd = *dat
					latestRepomdTime = elem.Timestamp
					log.Println("found newer")
				}
			}
		}
	}

	//log.Printf("latest: %+v", latestRepomd)
	trylist := []string{latestRepomd.mirror}
	trylist = append(trylist, mirrors...)

	{
		f, err := os.Create("repomd.xml")
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		_, err = f.Write(latestRepomd.fileContents)
		if err != nil {
			log.Fatal("Cannot write repomd.xml", err)
		}
	}

RepoMdFile:
	for _, filePath := range latestRepomd.Data {
		for _, tryMirror := range trylist {
			fileURL := tryMirror + "/" + repoPath + "/" + strings.TrimPrefix(filePath.Location.Href, "/")
			fmt.Println("getting", fileURL)
			fileData := readWithChecksum(fileURL, filePath.Checksum.Text, filePath.Checksum.Type)
			if fileData != nil {
				//fmt.Println("length", len(*fileData))
				//u, err := url.Parse(fileURL)
				//if err != nil {
				//	continue
				//}
				_, file := path.Split(fileURL)
				f, err := os.Create(file)
				if err != nil {
					log.Fatal(err)
				}
				defer f.Close()
				_, err = f.Write(*fileData)
				if err == nil {
					continue RepoMdFile
				}
			}
		}
	}

	err := ensureDir(*outputPath)
	if err != nil {
		log.Fatal(err)
	}
}

func check(e error) {
	if e != nil {
		//panic(e)
		log.Fatal(e)
	}
}
