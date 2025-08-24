# personal-vibesite

## Overview
Simple personal website and a Go server written with some serioue vibe coding.
I've procrastinated making a personal website for so long I finally set aside my software engineering skills and pulled out my prompt engineering skills.

## Setup

Running the site:
```sh
git clone https://github.com/Matthew-Berthoud/personal-vibesite.git
go install github.com/air-verse/air@latest
air
```

### Github Token
If you want to not be quite so rate-limited on your github API requests, you should make authenticated requests.
Create a file called `.env`, and put a [github personal access token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens) for your account in there.
```sh:.env
export GITHUB_TOKEN="YOUR_TOKEN_HERE"

```

