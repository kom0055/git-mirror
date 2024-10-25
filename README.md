# git-mirror an easy tool to mirror all you repos between Gitlab and Github

If there are too many repos on your Gitlab, it's too difficult to migrate or mirror those repos.

I was in the same situation before, to solve this, I wrote an easy tool to mirror all you visible repos between gitlab and github.
For example, it will mirror-clone all you visible repos on Gitlab, create repos on your Github account, and mirror-push them.

Gitee is not supported now.

## Build

```sh
git clone https://github.com/kom0055/git-mirror.git
cd git-mirror
make
```

## Usage

```
$ git-mirror --help
to clone repo: gclone https://github.com/kom0055/git-mirror
to sync all available repos: gclone --sync-from-remote --remote-type=gh https://github.com

Usage:
  gclone [flags]

Flags:
      --dest-ecdsa string                  ecdsa pem file
      --dest-ecdsa-passwd string           ecdsa pem file passwd
      --dest-proto string                  proto: git, ssh, http or https
      --dest-remote-gitlab-addr string     remote gitlab addr
      --dest-token string                  token or private key
      --dest-user string                   user name
  -h, --help                               help for gclone
      --source-ecdsa string                ecdsa pem file
      --source-ecdsa-passwd string         ecdsa pem file passwd
      --source-proto string                proto: git, ssh, http or https
      --source-remote-gitlab-addr string   remote gitlab addr
      --source-token string                token or private key
      --source-user string                 user name
      --worker int                         worker num (default 8)
```
