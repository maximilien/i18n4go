# Introducing i18n4go - Go language tool for internationalization (i18n)

## Abstract

In this post we will give an introductory overview of the i18n4go tool which allows one to easily globalize any Go language program. The i18n4go toolset was extracted while globalizing the CLoudFoundry command line interface and is an example of the CloudFoundry community usage of the Go language as well as contributing back to that community, in effect cross pollinating both communities, as was done with the Ginko and Gomega toolsets.

## Introduction
As the CloudFoundry project gets an increasingly global audience (users and developers) there is increasingly a need to globalize the parts of the system interfacing the main users. In particular, as IBM's Bluemix public installation of CF was announced, one of the gapping hole was to globalize the primary interface to CF from all users, the command line interface or CLI. Since the CF CLI is conversational by nature, it is a good candidate for internationalization (i18n), even if the majority of developers speak English, having the CLI converse in your own native tongue could make the whole experience more familiar and natural.

So it was with these motivation that the IBM CF community team decided to take on the major task of converting the entire CLI code base for globalization. Working closely with our Pivotal colleagues, what transpired in the span of about two months is a complete update of the CLI to enable any human language and a release of the CLI in the default English as well as French along with a call to action to the community to submit new translations. The following blog post does not chronicle this exercise, but rather gives an introduction to the resulting tooling that was develop in the process. This tool is a general purpose Go language tool (written in Go) to help in globalizing any Go program. We call it: i18n4go.

### Organization

* related tools
* architecture and design
* applying tool to small Go language project
* typical workflow
* conclusion and future

## Related Tools

* i18n tooling for Java and Ruby
* existing i18n tooling for Go language

## Architecture and Design

* i18n problem
* solving with tooling
* similar tools
* taking advantage of Go language's features

## Applying Tool

* using workflow to CF CLI
* using workflow to a small Go language program

## Typical Workflow

* applying to Go projects
* extracting strings
* merging strings
* rewriting code
* creating translations
* verifying translations
* maintaining

## Conclusion and Future

* recap of tool features
* using on CF CLI
* what is missing?
* immediate next steps

## References