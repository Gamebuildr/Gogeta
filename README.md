# Gogeta
A open source SCM micro service for monitoring and updating repositories.

 * Runs on Heroku
 * Supports Git
 * Repos are saved to Amazon S3
 * Future planned support for SVN
 * Future planned support for Perforce

Checkout the milestone list to see where this project is heading and how far along we are in development.


## Quick API

### Install

To run the app make sure to add all the go tools to your $PATH. Then cd into the cloned directory, run go install then gogeta. This will run the micro service and make the API available for debugging. Make sure port 9000 is open.

### Git shallow clone

```javascript
curl -XPOST -d'{"usr":"you","repo":"your.repo.url","project":"your.project"}' localhost:9000/0/gitclone
```
