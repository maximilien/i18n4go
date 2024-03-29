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

package test_helpers

import (
	"github.com/onsi/gomega/gexec"

	. "github.com/onsi/gomega"
)

var I18n4goExec string

func BuildExecutable() {
	var err error
	I18n4goExec, err = gexec.Build("./../../i18n4go")
	Ω(err).ShouldNot(HaveOccurred())
}
