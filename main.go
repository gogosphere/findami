package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type byCreationDate []*ec2.Image

func (a byCreationDate) Len() int      { return len(a) }
func (a byCreationDate) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a byCreationDate) Less(i, j int) bool {
	layout := "2006-01-02T15:04:05.000Z"
	t1, err := time.Parse(layout, *a[i].CreationDate)
	if err != nil {
		fmt.Println(err)
	}
	t2, err := time.Parse(layout, *a[j].CreationDate)
	if err != nil {
		fmt.Println(err)
	}
	return t1.After(t2)
}

func main() {
	name := flag.String("n", "ubuntu/images/hvm-ssd/ubuntu-trusty-14.04-amd64-server", "The name of the ami. * will be appended to enable a wildcard search.")
	owner := flag.String("o", "099720109477", "The owner of the ami.")
	region := flag.String("r", "us-east-1", "The AWS region you want to search.")
	version := flag.Bool("version", false, "Display version on stderr")
	flag.Parse()

	if *version {
		fmt.Fprintf(os.Stderr, "findami version: %s %s %s", Version, VersionPrerelease, GitCommit)
	} else {
		svc := ec2.New(session.New(), aws.NewConfig().WithRegion(*region).WithCredentialsChainVerboseErrors(true))
		findAndPrintAMI(svc, *name, *owner)
	}
}

func findAndPrintAMI(svc *ec2.EC2, name string, owner string) {
	params := &ec2.DescribeImagesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("name"),
				Values: []*string{aws.String(name + "*")},
			},
		},
		Owners: []*string{aws.String(owner)},
	}
	resp, err := svc.DescribeImages(params)

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	sort.Sort(byCreationDate(resp.Images))
	fmt.Println(*resp.Images[0].ImageId)
}
