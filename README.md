[![CircleCI](https://circleci.com/gh/staple-org/staple.svg?style=svg)](https://circleci.com/gh/staple-org/staple)
[![codecov](https://codecov.io/gh/staple-org/staple/branch/master/graph/badge.svg)](https://codecov.io/gh/staple-org/staple)

# Staple

The backend of Staple. This REST api allows for creation and management of staples.

# What are staples

Staples are chronologically ordered bookmarks which you wish to read later. Staple differs from other Read Later
app in that it doesn't allow more then a configured number of read later entries and one twist...

You can only read your entries in the order when you created them. This restriction is aiming to make it easier
to catch up on things that you wish to read later. Because you can only read in a FIFO (first in first out) manner
you are forced to go through your list in the priority in which you created it.

This will make it easier to actually read the things you wanted to read and not have things laying around somehwere
for eternity... unread. Never to be known.

Staples also provides a frontend written with React located here: [Staple Frontend](https://github.com/staple-org/frontend).

But it's totally okay to use staple through the API only. 

# Testing localhost https

```bash
mkcert -key-file key.pem -cert-file cert.pem localhost
```

# Deploying

All settings are through command line options. These options are defined through vault or
kubernetes secret storage.