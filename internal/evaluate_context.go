package internal

import "github.com/jeftadlvw/git-nest/models"

/*
EvaluateContext is a wrapper function for internal.CreateContext that also performs automatic validation on
the context's configuration.
*/
func EvaluateContext() (models.NestContext, error) {
	context, err := CreateContext()
	if err != nil {
		return context, err
	}

	err = context.Config.Validate()
	if err != nil {
		return context, err
	}

	return context, nil
}
