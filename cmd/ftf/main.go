package main

import (
	"fmt"
	"github.com/134130/ftf/internal/model"
	"github.com/134130/ftf/pkg/config"
	"github.com/134130/ftf/pkg/terminal"
	"github.com/134130/ftf/pkg/tree"
	"github.com/134130/ftf/pkg/view"
	"github.com/fatih/color"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"strings"
)

func main() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	tr := getTree()
	s := terminal.State{
		Root:   tr,
		Cursor: tr,
	}
	treeView := view.NewTreeView(config.DefaultGraphics, &s)
	previewView := view.NewPreviewView(config.DefaultGraphics, &s)
	searchbarView := view.NewSearchbarView(config.DefaultGraphics, &s)
	views := []terminal.ViewRenderer{searchbarView, treeView, previewView}

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

	log.Debug().Msg("bye")
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

	c1prev := strings.Builder{}
	c1prev.WriteString(color.BlueString(" [MySQL] aurora-cluster-1\n\n"))
	c1prev.WriteString(color.New().Sprint("  host    : mysql.querypie.io\n"))
	c1prev.WriteString(color.New().Sprint("  port    : 3306\n"))
	c1prev.WriteString(color.New().Sprint("  username: querypie\n"))
	c1prev.WriteString(color.New().Sprint("  password: querypie\n"))

	c1 := model.Cluster{
		UUID:    "5",
		Name:    "aurora-cluster-1",
		Preview: c1prev.String(),
		Parent:  &cg1,
	}
	cg1.Children = append(cg1.Children, &c1)

	c2 := model.Cluster{
		UUID:   "6",
		Name:   "aurora-cluster-2",
		Parent: &cg1,
	}
	cg1.Children = append(cg1.Children, &c2)

	cg2 := model.ClusterGroup{
		UUID:   "3",
		Name:   "qa-querypie-metastore-cluster",
		Parent: &cp,
	}
	cp.Children = append(cp.Children, &cg2)

	c3 := model.Cluster{
		UUID:   "7",
		Name:   "metastore-cluster-1",
		Parent: &cg2,
	}
	cg2.Children = append(cg2.Children, &c3)

	c4 := model.Cluster{
		UUID:   "8",
		Name:   "metastore-cluster-2",
		Parent: &cg2,
	}
	cg2.Children = append(cg2.Children, &c4)

	cg3 := model.ClusterGroup{
		UUID:   "4",
		Name:   "eks-zerocoke-dev-cluster",
		Parent: &cp,
	}
	cp.Children = append(cp.Children, &cg3)

	c5 := model.Cluster{
		UUID:   "9",
		Name:   "dev-cluster-1",
		Parent: &cg3,
	}
	cg3.Children = append(cg3.Children, &c5)

	c6 := model.Cluster{
		UUID:   "10",
		Name:   "dev-cluster-2",
		Parent: &cg3,
	}
	cg3.Children = append(cg3.Children, &c6)

	return &cp
}
