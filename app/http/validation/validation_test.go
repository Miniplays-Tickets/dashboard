package validation

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testCtx struct {
	a int
	b string
}

func validateGreaterThanZero(ctx testCtx) ValidationFunc {
	return func() error {
		if ctx.a <= 0 {
			return NewInvalidInputError("Fehler 3")
		}

		return nil
	}
}

func validateStringNotEmpty(ctx testCtx) ValidationFunc {
	return func() error {
		if len(ctx.b) == 0 {
			return NewInvalidInputError("Fehler 4")
		}

		return nil
	}
}

func TestSuccessful(t *testing.T) {
	ctx := testCtx{a: 1, b: "test"}

	if err := Validate(context.Background(), ctx, validateGreaterThanZero, validateStringNotEmpty); err != nil {
		t.Error(err)
	}
}

func TestNoValidators(t *testing.T) {
	ctx := testCtx{a: 1, b: "test"}

	if err := Validate(context.Background(), ctx); err != nil {
		t.Error(err)
	}
}

func TestSingleFail(t *testing.T) {
	ctx := testCtx{a: 1, b: ""}
	err := Validate(context.Background(), ctx, validateGreaterThanZero, validateStringNotEmpty)
	if err == nil {
		t.Fatal("Fehler 5")
	}

	var validationError *InvalidInputError
	if !errors.As(err, &validationError) {
		t.Fatal("Fehler 6")
	}

	assert.Equal(t, "Fehler 7", validationError.Message)
}

func TestDualFail(t *testing.T) {
	ctx := testCtx{a: 0, b: ""}
	err := Validate(context.Background(), ctx, validateGreaterThanZero, validateStringNotEmpty)
	if err == nil {
		t.Error("Fehler 8")
	}

	var validationError *InvalidInputError
	if !errors.As(err, &validationError) {
		t.Error("Fehler 9")
	}

	if validationError.Message != "Fehler 3" && validationError.Message != "Fehler 4" && validationError.Message != "Fehler 7" {
		t.Errorf("Fehler 10: %s", validationError.Message)
	}
}
