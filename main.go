// Written by Paul Schou (paulschou.com) March 2022
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
	"bytes"
	"flag"
	"fmt"
	"hash"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	"golang.org/x/crypto/openpgp/packet"
)

var version = "test"
var debug *bool

// Main is a function to fetch the HTTP repodata from a URL to get the latest
// package list for a repo
func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Yum Get RepoMD,  Version: %s\n\nUsage: %s [options...]\n\n", version, os.Args[0])
		flag.PrintDefaults()
	}

	var inRepoPath = flag.String("repo", "/7/os/x86_64", "Repo path to use in fetching")
	var mirrorList = flag.String("mirrors", "mirrorlist.txt", "Mirror / directory list of prefixes to use")
	var outputPath = flag.String("output", ".", "Path to put the repodata files")
	var insecure = flag.Bool("insecure", false, "Skip signature checks")
	var timeout = flag.Duration("timeout", 5*time.Second, "HTTP Client Timeout")
	var keyringFile = flag.String("keyring", "keys/", "Use keyring for verifying, keyring.gpg or keys/ directory")
	debug = flag.Bool("debug", false, "Turn on debug, more verbose")
	flag.Parse()

	client.Timeout = *timeout
	mirrors := readMirrors(*mirrorList)
	repoPath := strings.TrimSuffix(strings.TrimPrefix(*inRepoPath, "/"), "/")

	var latestRepomd Repomd
	var latestRepomdTime float64
	var keyring openpgp.EntityList
	if !*insecure {
		var err error
		if _, ok := isDirectory(*keyringFile); ok {
			//keyring = openpgp.EntityList{}
			for _, file := range getFiles(*keyringFile, ".gpg") {
				//fmt.Println("loading key", file)
				gpgFile := readFile(file)
				fileKeys, err := loadKeys(gpgFile)
				if err != nil {
					log.Fatal("Error loading keyring file", err)
				}
				//fmt.Println("  found", len(fileKeys), "keys")
				keyring = append(keyring, fileKeys...)
			}
		} else {
			gpgFile := readFile(*keyringFile)
			keyring, err = loadKeys(gpgFile)
			if err != nil {
				log.Fatal("Error loading keyring file", err)
			}
		}
		if len(keyring) == 0 {
			log.Fatal("no keys loaded")
		}
	}

	var mu sync.Mutex
	var wg sync.WaitGroup

	for j, mm := range mirrors {
		i := j
		m := mm

		wg.Add(1)
		go func() {
			defer wg.Done()
			repoPathSlash := m + "/" + repoPath + "/"
			repomdPath := repoPathSlash + "repodata/repomd.xml"
			repomdPathGPG := repoPathSlash + "repodata/repomd.xml.asc"
			log.Println(i, "Fetching", repomdPath)

			dat := readRepomdFile(repomdPath)
			mu.Lock()
			defer mu.Unlock()
			if dat != nil {
				for _, elem := range dat.Data {
					if elem.Timestamp > latestRepomdTime {
						if *insecure {
							//fmt.Printf("elem: %#v\n", elem.Timestamp)
							latestRepomdTime = elem.Timestamp
						} else {
							// Verify gpg signature file
							log.Println("Fetching signature file:", repomdPathGPG)
							gpgFile := readFile(repomdPathGPG)
							signature_block, err := armor.Decode(strings.NewReader(gpgFile))
							if len(gpgFile) > 100 && err == io.EOF {
								// This file may not be armored / encoded, go ahead and try the raw version
								signature_block = &armor.Block{Body: bytes.NewReader([]byte(gpgFile))}
								err = nil
							}
							if err != nil {
								//log.Println("Error decoding armored asc file:", err)
								log.Println("Signature file missing or unable to decode signature, maybe try with '-insecure' flag?")
								continue
							}
							p, err := packet.Read(signature_block.Body)
							if err != nil {
								log.Println("Unable parse signature file")
								continue
							}
							var signed_at time.Time
							var issuerKeyId uint64
							var hash hash.Hash

							switch sig := p.(type) {
							case *packet.Signature:
								issuerKeyId = *sig.IssuerKeyId
								signed_at = sig.CreationTime
								if hash == nil {
									hash = sig.Hash.New()
								}
							case *packet.SignatureV3:
								issuerKeyId = sig.IssuerKeyId
								signed_at = sig.CreationTime
								if hash == nil {
									hash = sig.Hash.New()
								}
							default:
								fmt.Println("Signature block is invalid")
								continue
							}

							if issuerKeyId == 0 {
								fmt.Println("Signature doesn't have an issuer")
								continue
							}

							if keyring == nil {
								fmt.Printf("  %s - Signed by 0x%02X at %v\n", repomdPathGPG, issuerKeyId, signed_at)
								os.Exit(1)
							} else {
								fmt.Printf("Verifying %s has been signed by 0x%02X at %v...\n", repomdPathGPG, issuerKeyId, signed_at)
							}
							keys := keyring.KeysByIdUsage(issuerKeyId, packet.KeyFlagSign)

							if len(keys) == 0 {
								fmt.Println("error: No matching public key found to verify")
								continue
							}
							if len(keys) > 1 {
								fmt.Println("warning: More than one public key found matching KeyID")
							}

							dat.ascFileContents = gpgFile
							fmt.Println("GPG Verified!")
						}
						if latestRepomdTime != 0 {
							log.Println("found newer")
						}
						dat.path = repoPathSlash
						dat.mirror = m
						latestRepomd = *dat
						latestRepomdTime = elem.Timestamp
					}
				}
			}
		}()
		time.Sleep(70 * time.Millisecond)
	}
	wg.Wait()

	//log.Printf("latest: %+v", latestRepomd)
	trylist := []string{latestRepomd.mirror}
	trylist = append(trylist, mirrors...)

	// Create the directory if needed
	err := ensureDir(*outputPath)
	if err != nil {
		log.Fatal("Cannot create output path:", err)
	}

	// Write out the repomd file into the path
	{
		outFile := path.Join(*outputPath, "repomd.xml")
		f, err := os.Create(outFile)
		if err != nil {
			log.Fatal("Error creating xml file", err)
		}
		_, err = f.Write(latestRepomd.fileContents)
		f.Close()
		timestamp := time.Unix(int64(latestRepomdTime), 0)
		//fmt.Printf("ts: %s\n", timestamp)
		os.Chtimes(outFile, timestamp, timestamp)
		if err != nil {
			log.Fatal("Cannot write repomd.xml", err)
		}
	}

	// If we have a signature file, write it out
	if len(latestRepomd.ascFileContents) > 0 {
		outFile := path.Join(*outputPath, "repomd.xml.asc")
		f, err := os.Create(outFile)
		if err != nil {
			log.Fatal("Error creating xml asc file", err)
		}
		_, err = f.Write([]byte(latestRepomd.ascFileContents))
		f.Close()
		timestamp := time.Unix(int64(latestRepomdTime), 0)
		os.Chtimes(outFile, timestamp, timestamp)
		if err != nil {
			log.Fatal("Cannot write repomd.xml.asc", err)
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
				outFile := path.Join(*outputPath, file)
				f, err := os.Create(outFile)
				if err != nil {
					log.Fatal(err)
				}
				_, err = f.Write(*fileData)
				f.Close()
				timestamp := time.Unix(int64(filePath.Timestamp), 0)
				os.Chtimes(outFile, timestamp, timestamp)

				if err == nil {
					continue RepoMdFile
				}
			}
		}
	}

}

func check(e error) {
	if e != nil {
		//panic(e)
		log.Fatal(e)
	}
}

// isDirectory determines if a file represented
// by `path` is a directory or not
func isDirectory(path string) (exist bool, isdir bool) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, false
	}
	return true, fileInfo.IsDir()
}

func getFiles(walkdir, suffix string) []string {
	ret := []string{}
	err := filepath.Walk(walkdir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				fmt.Println(err)
				return err
			}
			if !info.IsDir() && strings.HasSuffix(path, suffix) {
				ret = append(ret, path)
			}
			return nil
		})
	if err != nil {
		log.Fatal(err)
	}
	return ret
}
