# Yum Get Repo MetaData Utility

This shim fetches a yum repomd from a given repo and mirror list.  It verifies the PGP
signature and checksums of each file to ensure integrity.

# Example usage:
```bash
./yum-get-repomd -mirrors mirrorlist_centos.txt -repo "/7/os/x86_64" -keyring keys/ -output test
```

and the output looks like:
```
$ ./yum-get-repomd -output test
2022/03/23 09:13:28 Reading in file keys/RPM-GPG-KEY-CentOS-7.gpg
  1) Loaded KeyID: 0x24C6A8A7F4A80EB5
2022/03/23 09:13:28 Reading in file keys/RPM-GPG-KEY-EPEL-7.gpg
  1) Loaded KeyID: 0x6A2FAEA2352C64E5
2022/03/23 09:13:28 Reading in file keys/openresty-package.gpg
  1) Loaded KeyID: 0x97DB7443D5EDEB74
2022/03/23 09:13:28 0 Fetching http://mirror.umd.edu/centos/7/os/x86_64/repodata/repomd.xml
2022/03/23 09:13:28 Fetching signature file: http://mirror.umd.edu/centos/7/os/x86_64/repodata/repomd.xml.asc
Verifying http://mirror.umd.edu/centos/7/os/x86_64/repodata/repomd.xml.asc has been signed by 0x24C6A8A7F4A80EB5 at 2020-11-12 11:20:09 -0500 EST...
GPG Verified!
2022/03/23 09:13:28 1 Fetching http://mirror.mia11.us.leaseweb.net/centos/7/os/x86_64/repodata/repomd.xml
2022/03/23 09:13:28 2 Fetching http://mirrors.cmich.edu/centos/7/os/x86_64/repodata/repomd.xml
2022/03/23 09:13:28 3 Fetching http://mirror.dal10.us.leaseweb.net/centos/7/os/x86_64/repodata/repomd.xml
2022/03/23 09:13:28 4 Fetching http://mirror.math.princeton.edu/pub/centos/7/os/x86_64/repodata/repomd.xml
2022/03/23 09:13:28 5 Fetching http://linux.cc.lehigh.edu/centos/7/os/x86_64/repodata/repomd.xml
2022/03/23 09:13:28 6 Fetching http://mirror.chpc.utah.edu/pub/centos/7/os/x86_64/repodata/repomd.xml
2022/03/23 09:13:28 7 Fetching http://centos.mirrors.hoobly.com/7/os/x86_64/repodata/repomd.xml
2022/03/23 09:13:28 8 Fetching http://sjc.edge.kernel.org/centos/7/os/x86_64/repodata/repomd.xml
2022/03/23 09:13:28 9 Fetching http://mirror.den01.meanservers.net/centos/7/os/x86_64/repodata/repomd.xml
2022/03/23 09:13:28 10 Fetching http://mirror.wdc1.us.leaseweb.net/centos/7/os/x86_64/repodata/repomd.xml
2022/03/23 09:13:28 11 Fetching http://mirror.vtti.vt.edu/centos/7/os/x86_64/repodata/repomd.xml
2022/03/23 09:13:28 12 Fetching http://mirror.pit.teraswitch.com/centos/7/os/x86_64/repodata/repomd.xml
2022/03/23 09:13:28 13 Fetching http://mirror.centos.iad1.serverforge.org/7/os/x86_64/repodata/repomd.xml
2022/03/23 09:13:28 14 Fetching http://mirror.grid.uchicago.edu/pub/linux/centos/7/os/x86_64/repodata/repomd.xml
2022/03/23 09:13:28 15 Fetching http://mirror.hackingand.coffee/centos/7/os/x86_64/repodata/repomd.xml
2022/03/23 09:13:29 16 Fetching http://mirror.cs.pitt.edu/centos/7/os/x86_64/repodata/repomd.xml
2022/03/23 09:13:29 17 Fetching http://distro.ibiblio.org/centos/7/os/x86_64/repodata/repomd.xml
2022/03/23 09:13:29 18 Fetching http://bay.uchicago.edu/centos/7/os/x86_64/repodata/repomd.xml
2022/03/23 09:13:29 19 Fetching http://mirror.web-ster.com/centos/7/os/x86_64/repodata/repomd.xml
2022/03/23 09:13:29 20 Fetching http://mirror.nodesdirect.com/centos/7/os/x86_64/repodata/repomd.xml
2022/03/23 09:13:29 21 Fetching http://mirror.nodespace.net/centos/7/os/x86_64/repodata/repomd.xml
2022/03/23 09:13:29 22 Fetching http://mirrors.greenmountainaccess.net/centos/7/os/x86_64/repodata/repomd.xml
2022/03/23 09:13:29 23 Fetching http://mirror.cs.uwp.edu/pub/centos/7/os/x86_64/repodata/repomd.xml
getting http://mirror.umd.edu/centos/7/os/x86_64/repodata/cca56f3cffa18f1e52302dbfcf2f0250a94c8a37acd8347ed6317cb52c8369dc-c7-x86_64-comps.xml
getting http://mirror.umd.edu/centos/7/os/x86_64/repodata/5319616dde574d636861a6e632939f617466a371e59b555cf816cf1f52f3e873-filelists.xml.gz
getting http://mirror.umd.edu/centos/7/os/x86_64/repodata/a4e2b46586aa556c3b6f814dad5b16db5a669984d66b68e873586cd7c7253301-c7-x86_64-comps.xml.gz
getting http://mirror.umd.edu/centos/7/os/x86_64/repodata/2b479c0f3efa73f75b7fb76c82687744275fff78e4a138b5b3efba95f91e099e-primary.xml.gz
getting http://mirror.umd.edu/centos/7/os/x86_64/repodata/6d0c3a488c282fe537794b5946b01e28c7f44db79097bb06826e1c0c88bad5ef-primary.sqlite.bz2
getting http://mirror.umd.edu/centos/7/os/x86_64/repodata/ecaab5cc3b9c10fefe6be2ecbf6f9fcb437231dac3e82cab8d9d2cf70e99644d-other.sqlite.bz2
getting http://mirror.umd.edu/centos/7/os/x86_64/repodata/845e42288d3b73a069e781b4307caba890fc168327baba20ce2d78a7507eb2af-other.xml.gz
getting http://mirror.umd.edu/centos/7/os/x86_64/repodata/d6d94c7d406fe7ad4902a97104b39a0d8299451832a97f31d71653ba982c955b-filelists.sqlite.bz2

$ ls test
2b479c0f3efa73f75b7fb76c82687744275fff78e4a138b5b3efba95f91e099e-primary.xml.gz
5319616dde574d636861a6e632939f617466a371e59b555cf816cf1f52f3e873-filelists.xml.gz
6d0c3a488c282fe537794b5946b01e28c7f44db79097bb06826e1c0c88bad5ef-primary.sqlite.bz2
845e42288d3b73a069e781b4307caba890fc168327baba20ce2d78a7507eb2af-other.xml.gz
a4e2b46586aa556c3b6f814dad5b16db5a669984d66b68e873586cd7c7253301-c7-x86_64-comps.xml.gz
cca56f3cffa18f1e52302dbfcf2f0250a94c8a37acd8347ed6317cb52c8369dc-c7-x86_64-comps.xml
d6d94c7d406fe7ad4902a97104b39a0d8299451832a97f31d71653ba982c955b-filelists.sqlite.bz2
ecaab5cc3b9c10fefe6be2ecbf6f9fcb437231dac3e82cab8d9d2cf70e99644d-other.sqlite.bz2
repomd.xml
repomd.xml.asc
```

# Example Salt
```bash
$ ./yum-get-repomd -mirrors mirrorlist_salt.txt -repo  7/x86_64/latest/ -output salt-test -insecure -debug
```

# Example Puppet
```bash
$ yum-get-repomd -mirrors mirrorlist_puppet.txt -repo el/7/x86_64 -output puppet -keyring keys/
```
# Example Microsoft
```bash
./yum-get-repomd -mirrors mirrorlist_microsoft.txt -repo "7/prod" -keyring keys/ -output microsoft
```

# Usage help:
```bash
$ ./yum-get-repomd -h
Yum Get RepoMD,  Version: 0.1.2...

Usage: ./yum-get-repomd [options...]

  -insecure
        Skip signature checks
  -keyring string
        Use keyring for verifying, keyring.gpg or keys/ directory (default "keys/")
  -mirrors string
        Mirror / directory list of prefixes to use (default "mirrorlist.txt")
  -output string
        Path to put the repodata files (default ".")
  -repo string
        Repo path to use in fetching (default "/7/os/x86_64")
  -timeout duration
        HTTP Client Timeout (default 5s)
```




