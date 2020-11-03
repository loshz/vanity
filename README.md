# vanity [![Build Status](https://github.com/syscll/vanity/workflows/build/badge.svg)](https://github.com/syscll/vanity/actions)
Go vanity URL service.

## Usage
- Build a runable Docker image: `docker build -t syscll/vanity .`
- Run the docker image using the `vanity` command.

### Environment Variables
- `VANITY_VCS` defaults to git if not set.
- `VANITY_VCS_URL` e.g https://github.com/syscll
