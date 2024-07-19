// Copyright © 2015-2023 The Knative Authors
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

package common

import (
	"errors"
	"fmt"

	"go/ast"

	"github.com/maximilien/i18n4go/i18n4go/i18n"
)

func ImportsForASTFile(astFile *ast.File) (*ast.GenDecl, error) {
	for _, declaration := range astFile.Decls {
		decl, ok := declaration.(*ast.GenDecl)
		if !ok || len(decl.Specs) == 0 {
			continue
		}

		if _, ok = decl.Specs[0].(*ast.ImportSpec); ok {
			return decl, nil
		}
	}

	return nil, errors.New(fmt.Sprintf(i18n.T("Could not find imports for root node:\n\t{{.Arg0}}v\n", map[string]interface{}{"Arg0": astFile})))
}
