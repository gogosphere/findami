This tool can be used to find the latest AMI for things like packer builds.

Ex. By default it will get you the latest ubuntu ami.
```
$ findami
ami-663a6e0c
```

Ex. Or can do explicit searches
```
$ findami -n ubuntu/images/hvm-ssd/ubuntu-trusty-14.04-amd64-server -o 099720109477
ami-663a6e0c
```
