package constants

import "github.com/jeftadlvw/git-nest/models"

/*
Context holds the git-nest runtime context struct.

It will be evaluated at startup on non-help subcommands.
*/
var Context models.NestContext = models.NestContext{}
