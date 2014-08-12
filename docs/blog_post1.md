# Introducing i18n4go - Go language tooling for internationalization (i18n)

## Abstract

In this post we will give an introductory overview of the i18n4go tool which allows one to easily globalize any Go language program. The i18n4go toolset was extracted while globalizing the CLoudFoundry command line interface and is an example of the CloudFoundry community usage of the Go language as well as contributing back to that community, in effect cross pollinating both communities, as was done with the Ginko and Gomega toolsets.

## Introduction
As the CloudFoundry project gets an increasingly global audience (users and developers) there is increasingly a need to globalize the parts of the system interfacing the main users. In particular, as IBM's Bluemix public installation of CF was announced, one of the gapping hole was to globalize the primary interface to CF from all users, the command line interface or CLI. Since the CF CLI is conversational by nature, it is a good candidate for internationalization (i18n), even if the majority of developers speak English, having the CLI converse in your own native tongue could make the whole experience more familiar and natural.

So it was with these motivation that the IBM CF community team decided to take on the major task of converting the entire CLI code base for globalization. Working closely with our Pivotal colleagues, what transpired in the span of about two months is a complete update of the CLI to enable any human language and a release of the CLI in the default English as well as French along with a call to action to the community to submit new translations. The following blog post does not chronicle this exercise, but rather gives an introduction to the resulting tooling that was develop in the process. This tool is a general purpose Go language tool (written in Go) to help in globalizing any Go program. We call it: i18n4go.

### Organization

The rest of this post is organized as follows. First we take a look at what related tools exist in the Go language community that are helpful in globalizing Go programs. Next we will give an overview of our approach, which is completely tools-driven. After that we will use a public OSS Go program that could make use of globalization and walk you through the steps to convert it for internationalization. We then complete with a summary of the typical workflow of a developer using i18n4go and briefly touch on future works.

## Related Tools

Since the Go language community is rather young, there
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
* modifying strings
* maintaining strings

## Conclusion and Future

* recap of tool features
* using tool on CF CLI
* what is missing?
* immediate next steps
* future next steps

## References
