// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package functions

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/GehirnInc/crypt/sha512_crypt"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var saltPattern = regexp.MustCompile(`^[./0-9A-Za-z]{1,16}$`)

type SHA512Crypt struct{}

func NewSHA512Crypt() function.Function {
	return &SHA512Crypt{}
}

func (f *SHA512Crypt) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "sha512_crypt"
}

func (f *SHA512Crypt) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary: "Create a deterministic Linux SHA-512 crypt password hash.",
		Description: "Returns a $6$ password hash using an explicit salt and round count. " +
			"Use 5000 rounds for compatibility with openssl passwd -6 without a rounds option. " +
			"The rounds field is omitted at 5000 to preserve that exact output. " +
			"No random salt is generated. The password and round limits bound computation; " +
			"the full crypt round range is not supported. 5000 rounds is a compatibility setting, " +
			"not a security recommendation. Protect both the password and the hash with sensitive(); " +
			"sensitive values can still be persisted by Terraform. No appliance connection is required.",
		Parameters: []function.Parameter{
			function.StringParameter{Name: "password", Description: "Password as UTF-8 bytes, at most 4096 bytes; NUL bytes are rejected. Empty passwords are accepted."},
			function.StringParameter{Name: "salt", Description: "Explicit salt: 1 to 16 characters from ./0-9A-Za-z."},
			function.Int64Parameter{
				Name: "rounds", Description: "Iteration count from 1000 to 1000000; use 5000 for OpenSSL compatibility.",
			},
		},
		Return: function.StringReturn{},
	}
}

func (f *SHA512Crypt) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var passwordValue, saltValue types.String
	var roundsValue types.Int64
	resp.Error = req.Arguments.Get(ctx, &passwordValue, &saltValue, &roundsValue)
	if resp.Error != nil {
		return
	}
	if passwordValue.IsUnknown() || saltValue.IsUnknown() || roundsValue.IsUnknown() {
		resp.Error = resp.Result.Set(ctx, types.StringUnknown())

		return
	}
	if passwordValue.IsNull() || saltValue.IsNull() || roundsValue.IsNull() {
		resp.Error = function.NewFuncError("Password, salt, and rounds must not be null.")

		return
	}
	password, salt, rounds := passwordValue.ValueString(), saltValue.ValueString(), roundsValue.ValueInt64()
	if strings.ContainsRune(password, 0) || len(password) > 4096 {
		resp.Error = function.NewArgumentFuncError(0, "Password must not contain NUL bytes or exceed 4096 bytes.")

		return
	}
	if !saltPattern.MatchString(salt) {
		resp.Error = function.NewArgumentFuncError(1, "Salt must contain 1 to 16 characters from ./0-9A-Za-z.")

		return
	}
	if rounds < 1000 || rounds > 1000000 {
		resp.Error = function.NewArgumentFuncError(2, "Rounds must be between 1000 and 1000000.")

		return
	}
	if ctx.Err() != nil {
		resp.Error = function.NewFuncError("Password hashing was cancelled.")

		return
	}
	prefix := "$6$"
	if rounds != 5000 {
		prefix += fmt.Sprintf("rounds=%d$", rounds)
	}
	encoded, err := sha512_crypt.New().Generate([]byte(password), []byte(prefix+salt))
	if err != nil {
		resp.Error = function.NewFuncError("Unable to derive password hash.")

		return
	}
	resp.Error = resp.Result.Set(ctx, encoded)
}
