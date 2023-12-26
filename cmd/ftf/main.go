package main

import (
	"fmt"
	"github.com/134130/ftf/internal/model"
	"github.com/134130/ftf/pkg/config"
	"github.com/134130/ftf/pkg/terminal"
	"github.com/134130/ftf/pkg/tree"
	"github.com/134130/ftf/pkg/view"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.SetGlobalLevel(zerolog.NoLevel)

	tr := getTree()
	s := terminal.State{
		Root:   tr,
		Cursor: tr,
	}
	treeView := view.NewTreeView(config.DefaultGraphics, &s)
	previewView := view.NewPreviewView(config.DefaultGraphics, &s)
	views := []terminal.ViewRenderer{treeView, previewView}

	t, err := terminal.OpenTerm(&terminal.Config{
		Height: 1.0,
	})
	if err != nil {
		panic(err)
	}
	defer t.Close()

	flag, err := t.StartLoop(config.DefaultKeyBindings, views)
	if err != nil {
		panic(err)
	}
	t.Close()

	if flag == terminal.FlagPrint {
		for _, sel := range s.Selection {
			fmt.Println(sel.GetName())
		}
	}

	log.Info().Msg("bye")
}

func getTree() tree.TreeHandler {
	cp := model.CloudProvider{
		UUID: "1",
		Name: "dev-infra",
	}

	cg1 := model.ClusterGroup{
		UUID:   "2",
		Name:   "alpha-querypie-aurora-cluster",
		Parent: &cp,
	}
	cp.Children = append(cp.Children, &cg1)

	cg2 := model.ClusterGroup{
		UUID:   "3",
		Name:   "qa-querypie-metastore-cluster",
		Parent: &cp,
	}
	cp.Children = append(cp.Children, &cg2)

	cg3 := model.ClusterGroup{
		UUID:   "4",
		Name:   "eks-zerocoke-dev-cluster",
		Parent: &cp,
	}
	cp.Children = append(cp.Children, &cg3)

	return &cp
}
