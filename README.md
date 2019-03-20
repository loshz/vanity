# vanity [![Build Status](https://travis-ci.org/syscll/vanity.svg?branch=master)](https://travis-ci.org/syscll/vanity)
Go vanity URL service.

## Usage

- Build a runable Docker image: `docker build -t username/vanity .`
- Run the docker image using the `vanity` command.

### Environment Variables

- `VANITY_VCS` defaults to git if not set.
- `VANITY_VCS_URL` e.g https://github.com/syscll
