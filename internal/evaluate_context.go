package internal

import "github.com/jeftadlvw/git-nest/models"

/*
EvaluateContext is a wrapper function for internal.CreateContext that also performs automatic validation on
the context's configuration.
*/
func EvaluateContext() (models.NestContext, error) {
	context, err := CreateContextFromCurrentWorkingDir()
	if err != nil {
		return models.NestContext{}, err
	}

	err = context.Config.Validate()
	if err != nil {
		return models.NestContext{}, err
	}

	return context, nil
}
