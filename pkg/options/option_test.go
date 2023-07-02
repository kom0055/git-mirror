package options

import (
	"context"
	"testing"
)

func TestClone(t *testing.T) {

	var (
		gitlabOpt = BasicOpt{
			EcdsaPemFile:       "/Users/myuserr/.ssh/id_ecdsa",
			EcdsaPemFilePasswd: "",
			RemoteGitlabAddr:   "http://gitlab.mydomain.com/",
			User:               "user1",
			Token:              "xxxx",
			Proto:              "ssh",
			GroupName:          "group1",
		}

		githubOpt = BasicOpt{
			EcdsaPemFile:       "/Users/myuser/.ssh/id_ecdsa",
			EcdsaPemFilePasswd: "",
			User:               "user12",
			Token:              "xxxx",
			Proto:              "ssh",
			GroupName:          "group2",
		}
	)
	o := Option{
		Source: gitlabOpt,
		Dest:   githubOpt,
	}
	ctx := context.Background()
	if err := o.Mirror(ctx); err != nil {
		t.Fatal(err)
	}
}

func TestClone2(t *testing.T) {

	var (
		gitlabOpt = BasicOpt{
			EcdsaPemFile:       "",
			EcdsaPemFilePasswd: "",
			RemoteGitlabAddr:   "http://gitlab.mydomain.com/",
			User:               "xxx",
			Token:              "xxx",
			Proto:              "http",
			GroupName:          "mygroup1",
		}

		githubOpt = BasicOpt{
			EcdsaPemFile:       "/Users/myuser/.ssh/id_ecdsa",
			EcdsaPemFilePasswd: "",
			User:               "myuser",
			Token:              "xxx",
			Proto:              "ssh",
			GroupName:          "mygroup2",
		}
	)
	o := Option{
		Source: gitlabOpt,
		Dest:   githubOpt,
	}
	ctx := context.Background()
	if err := o.Mirror(ctx); err != nil {
		t.Fatal(err)
	}
}
