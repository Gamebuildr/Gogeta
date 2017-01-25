# Gogeta
Gogeta is an open source source control micro service. It allows you to clone new repositories, update repositories on user demand, and automatically detect source code changes to pull down any new code in the repositories that it monitors.

It currently supports Amazon SQS as its messaging queue service and git. There is future plans to extend the source control service support and currently nothing in the pipeline for adding more queue services.

The service is tailored for use with [Gamebuildr](http://www.gamebuildr.io) and we're always working on improving the solutions to make everything more generic.

## Quick Install

First make sure to install Go and have all the tools configured to your path.

Note: Gogeta will only work with v23 of the libgit2 repository.

Next Gogeta relies on [libgit2](https://libgit2.github.com) for running core git commands.

To install first clone the libgit2 repo. Next look for the latest tag patched for version v23 and checkout that tag into a new branch.

Then compile the library manually by running:

```
mkdir build && cd build
```
```
cmake .. -DCMAKE_INSTALL_PREFIX=/usr/local
```
```
cmake --build . --target install
```

Note: Sometimes you'll get an error when running the project and libgit2 can't find a specific lib file. If this is the case then set a new environment variable ```export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/usr/local/lib```

Then run ```go get github.com/Gamebuildr/Gogeta```

This should install the main project and dependencies. After that running the ```Gogeta``` command from anywhere will run the system on PORT 9000 (unless otherwise specified in your env variables)

### System Configs

At the moment the configs are documented inside the project itself. We'll get more documentation up soon to help configure these.
