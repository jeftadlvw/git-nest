package internal

import "github.com/jeftadlvw/git-nest/models"

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
