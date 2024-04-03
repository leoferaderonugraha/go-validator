package validator

import (
    "testing"
)

func runTestCases(fn func(any) bool, cases []string, cond bool) (string, bool) {
    for i := 0; i < len(cases); i++ {
        if fn(cases[i]) != cond {
            return cases[i], false
        }
    }

    return "", true
}

type requiredInput struct {
    Required string `json:"required" validate:"required"`
}

func requiredTest(input any) bool {
    _, ok := Validate(requiredInput{
        input.(string),
    })

    return ok
}

func Test_Required_Success(t *testing.T) {
    notEmptyValue, ok := runTestCases(
        requiredTest,
        []string{
            "some data",
            " ",  // should we consider this as valid?
        }, true);

    if !ok {
        t.Error("Valid input detected as empty!")
        t.Errorf("Input: %s\n", notEmptyValue)
    }
}

func Test_Required_Fail(t *testing.T) {
    notEmptyValue, ok := runTestCases(
        requiredTest,
        []string{
            "",
        }, false);

    if !ok {
        t.Error("Invalid input detected as non-empty!")
        t.Errorf("Input: %s\n", notEmptyValue)
    }
}


type emailInput struct {
    Email string `json:"email" validate:"email"`
}

func emailTest(email any) bool {
    _, ok := Validate(emailInput{
        email.(string),
    })

    return ok
}

func Test_Email_Success(t *testing.T) {
    invalidEmail, ok := runTestCases(
        emailTest,
        []string{
            "valid@address.com",
            "valid@address.va.li.ds",  // email with sub domain
        }, true);

    if !ok {
        t.Error("This email should be valid!")
        t.Errorf("Input: %s\n", invalidEmail)
    }
}

func Test_Email_Fail(t *testing.T) {
    validEmail, ok := runTestCases(
        emailTest,
        []string{
            "invalid_address",
            "invalid@address",
            "inv.al.id@address",
            // TODO:
            // - This should be invalid
            // "invalid@address.inv..alid",
        }, false)

    if !ok {
        t.Error("Invalid email address should contains some errors!")
        t.Errorf("Input: %s\n", validEmail)
    }
}

// TODO:
// - Add more test cases here to reach 90% coverage
