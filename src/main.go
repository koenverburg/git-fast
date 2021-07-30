package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	// "github.com/go-git/go-git/v5/plumbing/object"
	"github.com/koenverburg/git-fast/types"
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

  // Refactor this to be more generic eg: sanitizeCharacters(["??", "A"])
	s1 := strings.ReplaceAll(status.String(), "?? ", "")
	s2 := strings.ReplaceAll(s1, "A  ", "")
	s3 := strings.ReplaceAll(s2, " M ", "")
	s4 := strings.ReplaceAll(s3, "M  ", "")
	s5 := strings.ReplaceAll(s4, "MM ", "")

	return strings.Split(s5, "\n"), *worktree
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

func genericInputPrompt(label string, emptyAllowed bool) string {
	validate := func(input string) error {
    if emptyAllowed {
      return nil
    } else {
      if len(input) == 0 {
        return errors.New("no input given")
      } else {
        return nil
      }
    }
	}

	prompt := promptui.Prompt{
		Label:    label,
		Validate: validate,
	}

	result, err := prompt.Run()
	utils.CheckIfError(err)

	return result
}

func useSegment(value string) string {
	if !utils.IsEmpty(value) && value != "none" {
		return value
	}
	return ""
}

func includeInMessage(segments []types.Segment) string {
	var sb string
	for _, segment := range segments {
		switch {
		case segment.Part == "type":
			sb += useSegment(segment.Value)

		case segment.Part == "scope":
			if !utils.IsEmpty(segment.Value) {
				sb += fmt.Sprintf("(%s): ", segment.Value)
			}

		case segment.Part == "commit":
			sb += useSegment(segment.Value)
			sb += " "

		case segment.Part == "tag":
			sb += useSegment(segment.Value)
			sb += " "

		case segment.Part == "ticket":
			if !utils.IsEmpty(segment.Value) && segment.Value != "x" {
				sb += fmt.Sprintf("#%s ", segment.Value)
			}
		}
	}
	return sb
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
	typeSegement := utils.CreateSegment(typeString, "type")

	scope := genericInputPrompt("Scope", true)
	scopeSegement := utils.CreateSegment(scope, "scope")

	commit := genericInputPrompt("Commit message", true)
	commitSegement := utils.CreateSegment(commit, "commit")

	tags := genericSelectPrompt("Select tags", []string{
		"WIP",
		"[skip ci]",
	},
	)
	tagSegement := utils.CreateSegment(tags, "tag")

	ticket := genericInputPrompt("ticket number", false)
	ticketSegement := utils.CreateSegment(ticket, "ticket")

	result := includeInMessage([]types.Segment{
		ticketSegement,
		typeSegement,
		scopeSegement,
		commitSegement,
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

	// _, commitErr := worktree.Commit(msg, &git.CommitOptions{
	// Author: &object.Signature{
	// Name:  "Koen Verburg",
	// Email: "creativekoen@gmail.com",
	// When:  time.Now(),
	// },
	// })
	// utils.CheckIfError(commitErr)

	fmt.Println(msg)
}
