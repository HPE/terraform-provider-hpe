// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package functions_test

import (
	"context"
	"strings"
	"testing"

	"github.com/GehirnInc/crypt/sha512_crypt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HPE/terraform-provider-hpe/morpheus/functions"
)

func runHash(password, salt string, rounds int64) function.RunResponse {
	response := function.RunResponse{Result: function.NewResultData(types.StringUnknown())}
	functions.NewSHA512Crypt().Run(context.Background(), function.RunRequest{
		Arguments: function.NewArgumentsData([]attr.Value{
			types.StringValue(password), types.StringValue(salt), types.Int64Value(rounds),
		}),
	}, &response)

	return response
}

func TestUnitSHA512CryptOpenSSL(t *testing.T) {
	t.Parallel()
	const expected = "$6$saltsalt$qFmFH.bQmmtXzyBY0s9v7Oicd2z4XSIecDzlB5KiA2/jctKu9YterLp8wwnSq.qc.eoxqOmSuNp2xS0ktL3nh/"
	for range 2 {
		response := runHash("password", "saltsalt", 5000)
		if response.Error != nil {
			t.Fatal(response.Error)
		}
		if !response.Result.Equal(function.NewResultData(types.StringValue(expected))) {
			t.Fatal("result differs from OpenSSL vector")
		}
	}
}

func TestUnitSHA512CryptOpenSSLAdditionalVectors(t *testing.T) {
	t.Parallel()
	cases := []struct {
		password string
		salt     string
		rounds   int64
		expected string
	}{
		{strings.Repeat("a", 100), "saltsalt", 1000, "$6$rounds=1000$saltsalt$55YPoNFGMPdjADROQ0RqHcOANj0Bl3F3pKaD9IHMaD/aY0o66Skq8V8w1voV/.njxIg7tiTVdJr0./JbAWJ7E/"},
		{"unicode-ß-密码", "./0123456789abcd", 10000, "$6$rounds=10000$./0123456789abcd$w.8nArWVG5Levv4wHGVLsw.addXNM/jEQjtAApbwJS3BPFNXPEOIWwX7CBc9w6GG0PCkUviU/KlqNQvxDRnTv0"},
	}
	for _, test := range cases {
		response := runHash(test.password, test.salt, test.rounds)
		if response.Error != nil {
			t.Fatal(response.Error)
		}
		if !response.Result.Equal(function.NewResultData(types.StringValue(test.expected))) {
			t.Fatal("result differs from OpenSSL vector")
		}
	}
}

func TestUnitSHA512CryptBounds(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		password string
		salt     string
		rounds   int64
	}{
		{strings.Repeat("a", 4096), "0123456789abcdef", 1000},
		{"password", ".", 1000000},
	} {
		response := runHash(test.password, test.salt, test.rounds)
		if response.Error != nil {
			t.Fatal(response.Error)
		}
		value := response.Result.Value().(types.String).ValueString()
		if test.rounds == 1000000 {
			if value != "$6$rounds=1000000$.$FMY0VZjRNp763icn/llFNmh3nB2OCBEiVYB93x.3eyWR8zZMAKcqT3ZhbWoNAQTpGxWA2kMg.szyEk3.MwFIB." {
				t.Fatal("result differs from OpenSSL boundary vector")
			}
		} else if err := sha512_crypt.New().Verify(value, []byte(test.password)); err != nil {
			t.Fatal(err)
		}
	}
}

func TestUnitSHA512CryptCancelled(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	response := function.RunResponse{Result: function.NewResultData(types.StringUnknown())}
	functions.NewSHA512Crypt().Run(ctx, function.RunRequest{
		Arguments: function.NewArgumentsData([]attr.Value{
			types.StringValue("secret"), types.StringValue("salt"), types.Int64Value(5000),
		}),
	}, &response)
	if response.Error == nil || !strings.Contains(response.Error.Text, "cancelled") {
		t.Fatal("expected cancellation error")
	}
}

func TestUnitSHA512CryptParameters(t *testing.T) {
	t.Parallel()
	for _, password := range []string{"", "unicode-ß-密码", strings.Repeat("a", 100)} {
		for _, rounds := range []int64{1000, 5000, 10000} {
			response := runHash(password, "./0123456789abcd", rounds)
			if response.Error != nil {
				t.Fatal(response.Error)
			}
			value := response.Result.Value().(types.String).ValueString()
			if err := sha512_crypt.New().Verify(value, []byte(password)); err != nil {
				t.Fatal(err)
			}
			actualRounds, err := sha512_crypt.New().Cost(value)
			if err != nil || int64(actualRounds) != rounds {
				t.Fatalf("round count mismatch: %v", err)
			}
		}
	}
}

func TestUnitSHA512CryptInvalidInputs(t *testing.T) {
	t.Parallel()
	cases := []struct {
		password string
		salt     string
		rounds   int64
		argument int64
	}{
		{"secret\x00value", "salt", 5000, 0},
		{strings.Repeat("x", 4097), "salt", 5000, 0},
		{"secret", "", 5000, 1},
		{"secret", strings.Repeat("x", 17), 5000, 1},
		{"secret", "bad$salt", 5000, 1},
		{"secret", "bad salt", 5000, 1},
		{"secret", "salt", 999, 2},
		{"secret", "salt", 1000001, 2},
		{"secret", "salt", -1, 2},
	}
	for _, test := range cases {
		response := runHash(test.password, test.salt, test.rounds)
		if response.Error == nil || response.Error.FunctionArgument == nil || *response.Error.FunctionArgument != test.argument {
			t.Fatalf("expected error on argument %d", test.argument)
		}
		if strings.Contains(response.Error.Text, "secret") {
			t.Fatal("error discloses password")
		}
	}
}
