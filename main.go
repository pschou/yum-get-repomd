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
	"flag"
	"fmt"
	"hash"
	"log"
	"os"
	"path"
	"strings"
	"time"

	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	"golang.org/x/crypto/openpgp/packet"
)

var version = "test"

// HelloGet is an HTTP Cloud Function.
func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Yum Get RepoMD,  Version: %s\n\nUsage: %s [options...]\n\n", version, os.Args[0])
		flag.PrintDefaults()
	}

	var inRepoPath = flag.String("repo", "/7/os/x86_64", "Repo path to use in fetching")
	var mirrorList = flag.String("mirrors", "mirrorlist.txt", "Mirror / directory list of prefixes to use")
	var outputPath = flag.String("output", ".", "Path to put the repodata files")
	var insecure = flag.Bool("insecure", false, "Skip signature checks")
	var keyringFile = flag.String("keyring", "keyring.gpg", "Use keyring for verifying signature")
	flag.Parse()

	mirrors := readMirrors(*mirrorList)
	repoPath := strings.TrimSuffix(strings.TrimPrefix(*inRepoPath, "/"), "/")

	var latestRepomd Repomd
	var latestRepomdTime int
	var keyring openpgp.EntityList
	if !*insecure {
		var err error
		gpgFile := readFile(*keyringFile)
		keyring, err = loadKeys(gpgFile)
		if err != nil {
			log.Fatal("Error loading keyring file", err)
		}
	}

	for i, m := range mirrors {
		repoPath := m + "/" + repoPath + "/"
		repomdPath := repoPath + "repodata/repomd.xml"
		repomdPathGPG := repoPath + "repodata/repomd.xml.asc"
		log.Println(i, "Fetching", repomdPath)

		dat := readRepomdFile(repomdPath)
		if dat != nil {
			for _, elem := range dat.Data {
				if elem.Timestamp > latestRepomdTime {
					if !*insecure {
						// Verify gpg signature file
						log.Println("Fetching signature file:", repomdPathGPG)
						gpgFile := readFile(repomdPathGPG)
						signature_block, err := armor.Decode(strings.NewReader(gpgFile))
						if err != nil {
							log.Println("Unable decode signature")
							continue
						}
						p, err := packet.Read(signature_block.Body)
						if err != nil {
							log.Println("Unable parse signature")
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
							continue
						}

						fmt.Println("Verified!")
					}
					if latestRepomdTime != 0 {
						log.Println("found newer")
					}
					readFile(repomdPathGPG)
					dat.path = repoPath
					dat.mirror = m
					latestRepomd = *dat
					latestRepomdTime = elem.Timestamp
				}
			}
		}
	}

	//log.Printf("latest: %+v", latestRepomd)
	trylist := []string{latestRepomd.mirror}
	trylist = append(trylist, mirrors...)

	err := ensureDir(*outputPath)
	if err != nil {
		log.Fatal(err)
	}

	{
		f, err := os.Create(path.Join(*outputPath, "repomd.xml"))
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
				f, err := os.Create(path.Join(*outputPath, file))
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

}

func check(e error) {
	if e != nil {
		//panic(e)
		log.Fatal(e)
	}
}
