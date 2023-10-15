# Utils Server

## User Space

Data is only accessible to the user who created it.

### Authentication

Authentication is done via basic auth header `Authorization: Basic <base64 encoded username:password>`.

- username: GitHub username (not email)
- password: GitHub personal access token.

## Terraform HTTP Backend

This can be used as a backend for terraform.

URL format: `<protocol>://<host>:<port>/v1/tfstate/<workspace>`

```hcl
terraform {
  backend "http" {
    address        = "http://localhost:8080/v1/tfstate/test"
    lock_address   = "http://localhost:8080/v1/tfstate/test"
    unlock_address = "http://localhost:8080/v1/tfstate/test"
    username       = "arpanrec"
  }
}
```

## File Server

This can be used as a file server.

URL format: `<protocol>://<host>:<port>/v1/files/<path>`

## Configuration

Configuration is done via config json file.

File location can be set via environment variable `SECURE_SERVER_CONFIG_FILE_PATH`.

```json
{
  "encryption": {
    "private_key_path": "Path to GPG private key",
    "public_key_path": "Path to GPG public key",
    "private_key_password_path": "Password file for GPG private key",
    "delete_key_files_after_startup": "Boolean, delete key files after startup"
  },
  "storage": {
    "type": "Storage type, currently only supports: file, s3",
    "config": "Storage config"
  },
  "users": {
    "githubusername1": {},
    "githubusername2": {}
  }
}
```

### Configuration: File Storage

```json
{
  "storage": {
    "type": "file",
    "config": {
      "path": "Path to storage directory"
    }
  }
}
```

## Deployment

### Deployment: GitLab Runner

Upload the `.env` file to GitLab Secure Files. (GitLab Project -> Settings -> CI/CD -> Secure Files -> Upload `.env` File)

<details>
  <summary>GitLab Runner Installation</summary>

Deployment is done via [gitlab-runner](https://docs.gitlab.com/runner/install/linux-repository.html).
Add the Server as gitlab-runner with shell executor, also make sure gitlab runner has root access.

- Please check the [gitlab-runner](https://docs.gitlab.com/runner/install/linux-repository.html) for the latest installation instructions.

```bash
echo "gitlab-runner ALL=(ALL) NOPASSWD: ALL" | sudo tee /etc/sudoers.d/010-gitlab-runner >/dev/null
sudo curl -L --output /usr/local/bin/gitlab-runner "https://gitlab-runner-downloads.s3.amazonaws.com/latest/binaries/gitlab-runner-linux-$(dpkg --print-architecture)"
sudo chmod +x /usr/local/bin/gitlab-runner
sudo useradd --comment 'GitLab Runner' --create-home gitlab-runner --shell /bin/bash
sudo gitlab-runner install --user=gitlab-runner --working-directory=/home/gitlab-runner
sudo gitlab-runner start
sudo gitlab-runner status
sudo rm -rf /home/gitlab-runner/.bash_logout
```

- Issue with shell executor, [check this](https://docs.gitlab.com/runner/shells/index.html#shell-profile-loading).

- Register gitlab-runner with shell executor

Settings -> CI/CD -> Runners -> Expand -> `Enable shared runners for this project`: False -> Save changes

Settings -> CI/CD -> Runners -> Expand -> `New Project Runner`

```markdown
Operating systems: Linux
Tags: secureserver
Run untagged jobs: False
Details: Secure Server
Configuration (optional):
  - Paused: False
  - Protected: False
  - Lock to current projects: True
Maximum job timeout: 600
```

```bash
sudo gitlab-runner register \
  --non-interactive \
  --name secureserver \
  --url "https://gitlab.com" \
  --token "${TOKEN}" \
  --executor "shell"
```

- Remove gitlab-runner

```bash
sudo gitlab-runner uninstall
sudo rm -rf /usr/local/bin/gitlab-runner
sudo userdel -r gitlab-runner
sudo rm -rf /home/gitlab-runner/
sudo rm -rf /etc/gitlab-runner
```

</details>

Deployment is locked with branch name `main`, and when this is not a scheduled job.

### Deployment: GitHub Actions

Upload the base64 encoded `.env` file to GitHub Secrets as `ENVIRONMENT_FILE`. (GitHub Project -> Settings -> Secrets -> New repository secret)

<details>
  <summary>Github Actions Self Hosted Runner</summary>

Deployment is done via [GitHub Actions Self Hosted Runner](https://docs.github.com/en/actions/hosting-your-own-runners/about-self-hosted-runners). 
Make sure GitHub Actions Self Hosted Runner has NOPASSWD root access.

- Install GitHub Actions Self Hosted Runner

```bash
sudo useradd -m -s /bin/bash actions-runner
echo "actions-runner ALL=(ALL) NOPASSWD: ALL" | sudo tee /etc/sudoers.d/010-actions-runner >/dev/null
sudo su - actions-runner
cd ~
curl -o actions-runner-linux-x64-2.309.0.tar.gz \
  -L https://github.com/actions/runner/releases/download/v2.309.0/actions-runner-linux-x64-2.309.0.tar.gz
echo "2974243bab2a282349ac833475d241d5273605d3628f0685bd07fb5530f9bb1a  actions-runner-linux-x64-2.309.0.tar.gz" | shasum -a 256 -c
tar xzf ./actions-runner-linux-x64-2.309.0.tar.gz
./config.sh --url https://github.com/arpanrec/secureserver --token "${TOKEN}" --name secureserver --work _work --labels secureserver --unattended
sudo ./svc.sh install
sudo ./svc.sh start
```

- Remove GitHub Actions Self Hosted Runner

```bash
sudo ./svc.sh stop
sudo ./svc.sh uninstall
./config.sh remove --token "${TOKEN}"
sudo userdel -r actions-runner
```

</details>

## Backup

Download a copy of secure working directory from secure server and upload it to offshore storage.

### Backup: GitLab Runner

Upload the `.env` file to GitLab Secure Files. (GitLab Project -> Settings -> CI/CD -> Secure Files -> Upload `.env` File)

- For runner installation, [check this](#deployment-gitlab-runner).

Backup is locked with branch name `main`.

A timer job is scheduled to run every 6 hours via [Scheduled pipelines](https://docs.gitlab.com/ee/ci/pipelines/schedules.html).

```markdown
Description: Backup Secure Server Working Directory
Interval Pattern: 0 */6 * * *
Cron timezone: [UTC+5.5] Kolkata
Select target branch or tag: main
Activated: true
```

### Backup: GitHub Actions

Upload the base64 encoded `.env` file to GitHub Secrets as `ENVIRONMENT_FILE`. (GitHub Project -> Settings -> Secrets -> New repository secret)

- For runner installation, [check this](#deployment-github-actions).

Backup is locked with branch name `main`.

A timer job is scheduled to run every 6 hours via [Scheduled Events](https://docs.github.com/en/actions/reference/events-that-trigger-workflows#scheduled-events).
