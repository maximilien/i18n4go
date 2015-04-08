package verify_strings_test

import (
	"os"
	"path/filepath"

	. "github.com/maximilien/i18n4go/integration/test_helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("verify-strings -f fileName --languages \"[lang,?]+\"", func() {
	var (
		inputFilesPath    string
		expectedFilesPath string
	)

	BeforeEach(func() {
		fixturesPath := filepath.Join("..", "..", "test_fixtures", "verify_strings")
		inputFilesPath = filepath.Join(fixturesPath, "f_option", "input_files")
		expectedFilesPath = filepath.Join(fixturesPath, "f_option", "expected_output")
	})

	Context("valid input file provided", func() {
		Context("using --source-language", func() {
			Context("passes verifications", func() {
				BeforeEach(func() {
					session := Runi18n("-c", "verify-strings", "-v", "-f", filepath.Join(inputFilesPath, "quota.go.en.json"), "--languages", "\"fr,zh_CN\"", "-o", expectedFilesPath, "--source-language", "en")
					Ω(session.ExitCode()).Should(Equal(0))
				})

				It("with language file with valid keys", func() {
					_, err := os.Stat(GetFilePath(inputFilesPath, "quota.go.en.json.missing.diff"))
					Ω(os.IsNotExist(err)).Should(Equal(true))

					_, err = os.Stat(GetFilePath(inputFilesPath, "quota.go.en.json.extra.diff"))
					Ω(os.IsNotExist(err)).Should(Equal(true))
				})
			})
		})

		Context("not using --source-language", func() {
			Context("passes verifications", func() {
				BeforeEach(func() {
					session := Runi18n("-c", "verify-strings", "-v", "-f", filepath.Join(inputFilesPath, "quota.go.en.json"), "--languages", "\"fr,zh_CN\"", "-o", expectedFilesPath)
					Ω(session.ExitCode()).Should(Equal(0))
				})

				It("with language file with valid keys", func() {
					_, err := os.Stat(GetFilePath(inputFilesPath, "quota.go.fr.json.missing.diff.json"))
					Ω(os.IsNotExist(err)).Should(Equal(true))

					_, err = os.Stat(GetFilePath(inputFilesPath, "quota.go.fr.json.missing.diff.json"))
					Ω(os.IsNotExist(err)).Should(Equal(true))

					_, err = os.Stat(GetFilePath(inputFilesPath, "quota.go.zh_CN.json.extra.diff.json"))
					Ω(os.IsNotExist(err)).Should(Equal(true))

					_, err = os.Stat(GetFilePath(inputFilesPath, "quota.go.zh_CN.json.extra.diff.json"))
					Ω(os.IsNotExist(err)).Should(Equal(true))
				})
			})
		})

		Context("fails verification", func() {
			Context("with one language file", func() {
				Context("with missing keys", func() {
					BeforeEach(func() {
						session := Runi18n("-c", "verify-strings", "-v", "-f", filepath.Join(inputFilesPath, "quota.go.en.json"), "--languages", "\"de\"", "-o", expectedFilesPath, "--source-language", "en")
						Ω(session.ExitCode()).Should(Equal(1))
					})

					AfterEach(func() {
						RemoveAllFiles(
							GetFilePath(expectedFilesPath, "quota.go.de.json.missing.diff.json"),
						)
					})

					It("generates missing diff file", func() {
						fileInfo, err := os.Stat(GetFilePath(expectedFilesPath, "quota.go.de.json.missing.diff.json"))
						Ω(err).Should(BeNil())
						Ω(fileInfo.Name()).Should(Equal("quota.go.de.json.missing.diff.json"))
					})
				})

				Context("with missing and extra keys", func() {
					BeforeEach(func() {
						session := Runi18n("-c", "verify-strings", "-v", "-f", filepath.Join(inputFilesPath, "quota.go.en.json"), "--languages", "\"af\"", "-o", expectedFilesPath, "--source-language", "en")
						Ω(session.ExitCode()).Should(Equal(1))
					})

					AfterEach(func() {
						RemoveAllFiles(
							GetFilePath(expectedFilesPath, "quota.go.af.json.missing.diff.json"),
							GetFilePath(expectedFilesPath, "quota.go.af.json.extra.diff.json"),
						)
					})

					It("generates missing and extra diff file", func() {
						fileInfo, err := os.Stat(GetFilePath(expectedFilesPath, "quota.go.af.json.missing.diff.json"))
						Ω(err).Should(BeNil())
						Ω(fileInfo.Name()).Should(Equal("quota.go.af.json.missing.diff.json"))

						fileInfo, err = os.Stat(GetFilePath(expectedFilesPath, "quota.go.af.json.extra.diff.json"))
						Ω(err).Should(BeNil())
						Ω(fileInfo.Name()).Should(Equal("quota.go.af.json.extra.diff.json"))
					})
				})

				Context("with templated keys whose translation does not contain same arguments", func() {
					BeforeEach(func() {
						session := Runi18n("-c", "verify-strings", "-v", "-f", filepath.Join(inputFilesPath, "quota.go.en.json"), "--languages", "\"es\"", "-o", expectedFilesPath, "--source-language", "en")
						Ω(session.ExitCode()).Should(Equal(1))
					})

					AfterEach(func() {
						RemoveAllFiles(
							GetFilePath(expectedFilesPath, "quota.go.es.json.invalid.diff.json"),
						)
					})

					It("generates invalid diff file", func() {
						fileInfo, err := os.Stat(GetFilePath(expectedFilesPath, "quota.go.es.json.invalid.diff.json"))
						Ω(err).Should(BeNil())
						Ω(fileInfo.Name()).Should(Equal("quota.go.es.json.invalid.diff.json"))
					})
				})
			})

			Context("with multiple language files", func() {
				BeforeEach(func() {
					session := Runi18n("-c", "verify-strings", "-v", "-f", filepath.Join(inputFilesPath, "quota.go.en.json"), "--languages", "\"de,it\"", "-o", expectedFilesPath, "--source-language", "en")
					Ω(session.ExitCode()).Should(Equal(1))
				})

				AfterEach(func() {
					RemoveAllFiles(
						GetFilePath(expectedFilesPath, "quota.go.de.json.missing.diff.json"),
						GetFilePath(expectedFilesPath, "quota.go.it.json.missing.diff.json"),
					)
				})

				It("with invalid keys", func() {
					fileInfo, err := os.Stat(GetFilePath(expectedFilesPath, "quota.go.de.json.missing.diff.json"))
					Ω(err).Should(BeNil())
					Ω(fileInfo.Name()).Should(Equal("quota.go.de.json.missing.diff.json"))

					fileInfo, err = os.Stat(GetFilePath(expectedFilesPath, "quota.go.it.json.missing.diff.json"))
					Ω(err).Should(BeNil())
					Ω(fileInfo.Name()).Should(Equal("quota.go.it.json.missing.diff.json"))
				})
			})
		})

		Context("with language file", func() {
			BeforeEach(func() {
				session := Runi18n("-c", "verify-strings", "-v", "-f", filepath.Join(inputFilesPath, "quota.go.en.json"), "--languages", "\"ja\"", "-o", expectedFilesPath)
				Ω(session.ExitCode()).Should(Equal(1))
			})

			AfterEach(func() {
				RemoveAllFiles(
					GetFilePath(expectedFilesPath, "quota.go.ja.json.extra.diff.json"),
				)
			})

			It("with additional keys", func() {
				fileInfo, err := os.Stat(GetFilePath(expectedFilesPath, "quota.go.ja.json.extra.diff.json"))
				Ω(err).Should(BeNil())
				Ω(fileInfo.Name()).Should(Equal("quota.go.ja.json.extra.diff.json"))
			})
		})

		Context("with multiple language file", func() {
			BeforeEach(func() {
				session := Runi18n("-c", "verify-strings", "-v", "-f", filepath.Join(inputFilesPath, "quota.go.en.json"), "--languages", "\"ja,cs\"", "-o", expectedFilesPath)
				Ω(session.ExitCode()).Should(Equal(1))
			})

			AfterEach(func() {
				RemoveAllFiles(
					GetFilePath(expectedFilesPath, "quota.go.ja.json.extra.diff.json"),
					GetFilePath(expectedFilesPath, "quota.go.cs.json.extra.diff.json"),
				)
			})

			It("with additional keys", func() {
				fileInfo, err := os.Stat(GetFilePath(expectedFilesPath, "quota.go.ja.json.extra.diff.json"))
				Ω(err).Should(BeNil())
				Ω(fileInfo.Name()).Should(Equal("quota.go.ja.json.extra.diff.json"))

				fileInfo, err = os.Stat(GetFilePath(expectedFilesPath, "quota.go.cs.json.extra.diff.json"))
				Ω(err).Should(BeNil())
				Ω(fileInfo.Name()).Should(Equal("quota.go.cs.json.extra.diff.json"))
			})
		})

		Context("when missing a language file", func() {
			BeforeEach(func() {
				session := Runi18n("-c", "verify-strings", "-v", "-f", filepath.Join(inputFilesPath, "quota.go.en.json"), "--languages", "\"ja,ht\"", "-o", expectedFilesPath)
				Ω(session.ExitCode()).Should(Equal(1))
			})

			AfterEach(func() {
				RemoveAllFiles(
					GetFilePath(expectedFilesPath, "quota.go.ja.json.extra.diff.json"),
				)
			})

			It("with additional keys", func() {
				fileInfo, err := os.Stat(GetFilePath(expectedFilesPath, "quota.go.ja.json.extra.diff.json"))
				Ω(err).Should(BeNil())
				Ω(fileInfo.Name()).Should(Equal("quota.go.ja.json.extra.diff.json"))
			})
		})
	})

	Context("invalid input file provided", func() {
		Context("does not exist", func() {
			BeforeEach(func() {
				session := Runi18n("-c", "verify-strings", "-v", "-f", filepath.Join(inputFilesPath, "quota.go.ht.json"), "--languages", "\"fr\"", "-o", expectedFilesPath, "--source-language", "en")
				Ω(session.ExitCode()).Should(Equal(1))
			})

			It("fails verification", func() {
				_, err := os.Stat(GetFilePath(expectedFilesPath, "quota.go.ht.json"))
				Ω(os.IsNotExist(err)).Should(Equal(true))
			})
		})

		Context("does not have any keys", func() {
			BeforeEach(func() {
				session := Runi18n("-c", "verify-strings", "-v", "-f", filepath.Join(inputFilesPath, "quota.go.vi.json"), "--languages", "\"fr\"", "-o", expectedFilesPath, "--source-language", "en")
				Ω(session.ExitCode()).Should(Equal(1))
			})

			It("fails verification", func() {
				fileInfo, err := os.Stat(GetFilePath(inputFilesPath, "quota.go.vi.json"))
				Ω(err).Should(BeNil())
				Ω(fileInfo.Name()).Should(Equal("quota.go.vi.json"))
			})
		})
	})
})
