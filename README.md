Bug Stats
===================

Analyzes git repositories for files and methods which are fixed in bugfix commits.
Very much not finished yet, but the goal is to be able to identify problematic files and methods. Taking it even further the stats could be used to predict which methods are likely to be buggy. That'd be dope indeed :)


To start the program run
shell> $(./setup-gopath.sh)
shell> go run src/main.go


To update the dependencies run
shell> ./update-dependencies.sh
