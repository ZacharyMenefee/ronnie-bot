# ronnie-bot

the ronnie-bot, named after [ronnie cramer](https://en.wikipedia.org/wiki/Ronnie_Cramer), the world's most influential director, at your service.

## development

clone the repo and do `make run`. this runs a new dev server that will respond to messages on the `bgsdigital_tech` channel and ignores others. this means when multiple developers are running at the same time both instances respond to messages on that channel but whatever.

## deployment

to take a break from the cloud-native containerized deployents and all that noise I elected to do this as dumb as possible. we have a single digital ocean VM running this as a systemd service. when we merge into master a github action ssh's onto the machine, removes the repo, and restarts the service. ¯\_(ツ)_/¯