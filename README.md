# Friendly Errors for Go Validator

## Purpose
This project is a simple library intended to enhance the error messages provided by the [Go Validator](https://github.com/go-playground/validator) library. It provides user-friendly and more informative error messages, making it easier for developers to identify and resolve validation issues.

## Usage Examples
Here are some examples of how to use the Friendly Errors library in your Go project:

```go
package main

import (
	"fmt"

	friendlyerrors "github.com/ProImpact/friendly-errors"
)

// User struct which has validation tags
type User struct {
	Email    string `validate:"required,email" json:"email,omitempty"`
	Password string `validate:"required,min=8" json:"password,omitempty"`
}

func main() {
	user := &User{
		Email:    "invalid-email",
		Password: "short",
	}

	// Validate the struct
	err := friendlyerrors.ValidateAny(user)
	if len(err) > 0 {
		for k, v := range err {
			fmt.Printf("Error: field -> %s, message -> %s\n", k, v)
		}
	}
}
```

## Collaboration Guide
We welcome contributions to improve this library. Here's how you can collaborate:

1. **Fork the repository**: Create a personal fork of the project on GitHub to work with.
2. **Clone the fork**: Clone your fork of the repository to your local development environment.
3. **Create a branch**: Always create a new branch for your work using `git checkout -b feature-branch-name`.
4. **Make changes**: Implement your changes, ensuring the code is well-documented and tested.
5. **Commit your changes**: Use clear and descriptive commit messages.
6. **Push to GitHub**: Push your changes to your forked repository on GitHub.
7. **Create a Pull Request**: Navigate to the original repository and create a pull request, describing your changes and the problem they solve.
8. **Review**: Engage in the review process, address any feedback or questions.

Thank you for considering contributing to the project! For any questions, please open an issue in the GitHub repository.

