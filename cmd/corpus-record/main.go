package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/rui-ktei/socrates/internal/corpus"
	"github.com/rui-ktei/socrates/internal/corpus/runs"
	"github.com/rui-ktei/socrates/internal/runstore"
	"github.com/rui-ktei/socrates/internal/state"
)

const dateLayout = "2006-01-02"

func main() {
	if err := record(); err != nil {
		fmt.Fprintln(os.Stderr, failurePrefix+err.Error())
		os.Exit(1)
	}
}

func record() error {
	request := runs.Request{}
	flag.StringVar(&request.Run, "run", "", runFlag)
	flag.IntVar(&request.Round, "round", 1, roundFlag)
	flag.StringVar(&request.Name, "name", "", nameFlag)
	flag.StringVar(&request.Subject, "subject", "", subjectFlag)
	flag.StringVar(&request.On, "on", time.Now().UTC().Format(dateLayout), onFlag)
	flag.Parse()

	if request.Run == "" || request.Name == "" || request.Subject == "" {
		return fmt.Errorf(nothingToRecord)
	}

	root, err := os.Getwd()
	if err != nil {
		return err
	}

	home, err := state.Runs()
	if err != nil {
		return err
	}

	answered, err := runs.Answered(runstore.New(home), request)
	if err != nil {
		return err
	}

	at, err := corpus.Record(root, answered.Repository, answered.Entry)
	if err != nil {
		return err
	}

	fmt.Printf(recorded+"\n", answered.Entry.Name, answered.Entry.Model, answered.Entry.Effort, answered.Entry.Round, at)
	return nil
}
