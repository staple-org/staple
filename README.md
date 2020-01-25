[![CircleCI](https://circleci.com/gh/staple-org/staple.svg?style=svg)](https://circleci.com/gh/staple-org/staple)
[![codecov](https://codecov.io/gh/staple-org/staple/branch/master/graph/badge.svg)](https://codecov.io/gh/staple-org/staple)

# Staple

The backend of staple

# Testing localhost https

```bash
mkcert -key-file key.pem -cert-file cert.pem localhost
```

# Deploying

All settings are through command line options. These options are defined through vault or
kubernetes secret storage.