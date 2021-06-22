package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/koenverburg/git-fast/utils"
	"github.com/manifoldco/promptui"
)

func getUntrackedFiles() ([]string, git.Worktree) {
  dir, err := filepath.Abs(filepath.Dir(os.Args[1]))
  utils.CheckIfError(err)
  fmt.Println(fmt.Sprint("folder %s", dir))

  repo, err := git.PlainOpen(dir)
  utils.CheckIfError(err)

  worktree, err := repo.Worktree()
  utils.CheckIfError(err)

  status, err := worktree.Status()
  utils.CheckIfError(err)

  s1 := strings.ReplaceAll(status.String(), "?? ", "")
  s2 := strings.ReplaceAll(s1, "A  ", "")
  s3 := strings.ReplaceAll(s2, " M ", "")

  return strings.Split(s3, "\n"), *worktree
}

func genericSelectPrompt(label string, items []string) string {
  prompt := promptui.Select{
		Label: label,
    Items: utils.FilterEmptyString(append(items, "none")),
  }

	_, result, err := prompt.Run()
  utils.CheckIfError(err)

  return result
}

func showUntrackedList(files []string) string {
  prompt := promptui.Select{
		Label: "Select files to commit",
		Items: utils.FilterEmptyString(append(files, "done")),
	}

	_, result, err := prompt.Run()
  utils.CheckIfError(err)

  return result
}

func stageSelectedFiles(files []string, worktree git.Worktree) {
  for _, v := range files {
    worktree.Add(v)
  }
}

func genericInputPrompt(label string) string {
  validate := func(input string) error {
    if len(input) == 0 {
      return errors.New("no input given")
    }
		return nil
	}

	prompt := promptui.Prompt{
		Label:    label,
		Validate: validate,
	}

	result, err := prompt.Run()
  utils.CheckIfError(err)

  return result
}

type Segment struct {
  value string
  part string
}

func isEmpty(str string) bool {
  if str == "" {
    return true
  }
  return false
}

func useSegment(value string) string {
  if !isEmpty(value) && value != "none" {
    return value
  }
  return ""
}

func includeInMessage(segments []Segment) string {
  var sb string
  for _, segment := range segments {
    switch {
    case segment.part == "type":
      sb += useSegment(segment.value)

    case segment.part == "scope":
      if !isEmpty(segment.value) {
        sb += fmt.Sprintf("(%s): ", segment.value)
      }

    case segment.part == "commit":
      sb += useSegment(segment.value)
      sb += " "

    case segment.part == "tag":
      sb += useSegment(segment.value)
      sb += " "

    case segment.part == "ticket":
      if !isEmpty(segment.value) && segment.value != "x" {
        sb += fmt.Sprintf("#%s ", segment.value)
      }
    }
  }
  return sb
}

func createSegment(value string, part string) Segment {
  var s Segment
  s.value = value
  s.part = part
  return s
}

func commitMessageWizard() string {
  typeString := genericSelectPrompt("Select the type of change", []string{
      "feat",
      "fix",
      "docs",
      "style",
      "refactor",
      "perf",
      "test",
      "chore",
    },
  )
  typeSegement := createSegment(typeString, "type")

  scope := genericInputPrompt("Scope")
  scopeSegement := createSegment(scope, "scope")

  commit := genericInputPrompt("Commit message")
  commitSegement := createSegment(commit, "commit")

  tags := genericSelectPrompt("Select tags", []string{
      "WIP",
      "[skip ci]",
    },
  )
  tagSegement := createSegment(tags, "tag")

  ticket := genericInputPrompt("ticket number")
  ticketSegement := createSegment(ticket, "ticket")

  result := includeInMessage([]Segment{
    typeSegement,
    scopeSegement,
    commitSegement,
    ticketSegement,
    tagSegement,
  })

  return result 
}

func main() {
  files, worktree := getUntrackedFiles()

  var selection []string

  for {
    selected := showUntrackedList(files)
    if selected == "done" {
      break
    } else {
      selection = append(selection, selected)
    }
  }

  stageSelectedFiles(selection, worktree)

  msg := commitMessageWizard()

  _, commitErr := worktree.Commit(msg, &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Koen Verburg",
			Email: "creativekoen@gmail.com",
			When:  time.Now(),
		},
	})
  utils.CheckIfError(commitErr)

  fmt.Println(msg)
}
