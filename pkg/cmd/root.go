// Copyright © 2020 The Tekton Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"context"
	"fmt"

	"github.com/open-policy-agent/opa/rego"
	"github.com/spf13/cobra"
	"github.com/tektoncd/catlin/pkg/app"
	"github.com/tektoncd/catlin/pkg/cmd/bump"
	"github.com/tektoncd/catlin/pkg/cmd/linter"
	"github.com/tektoncd/catlin/pkg/cmd/validate"
)

func Root(cli app.CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:          "catlin",
		Short:        "Lints Tekton Resources and Catalogs",
		Long:         ``,
		SilenceUsage: true,
	}

	cmd.AddCommand(
		eval(cli),
		validate.Command(cli),
		linter.Command(cli),
		bump.Command(cli),
	)

	return cmd
}

func eval(cli app.CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use: "eval",
		RunE: func(cmd *cobra.Command, args []string) error {
			/*
				https://pkg.go.dev/github.com/open-policy-agent/opa/rego#example-Rego.Eval-Storage
				https://github.com/enterprise-contract/ec-policies/blob/main/policy/lib/bundles.rego
				https://github.com/enterprise-contract/ec-cli/blob/main/cmd/validate/input.go
			*/
			ctx := context.Background()

			// Create query that returns a single boolean value.
			rego := rego.New(
				rego.Query("data.authz.allow"),
				rego.Module("example.rego",
					`package authz

default allow = false
allow {
	input.open == "sesame"
}`,
				),
				rego.Input(map[string]interface{}{"open": "bar", "foo": "baz"}),
			)

			// Run evaluation.
			rs, err := rego.Eval(ctx)
			if err != nil {
				panic(err)
			}

			// Inspect result.
			fmt.Println("allowed:", rs.Allowed())
			return nil
		},
	}

	return cmd
}
