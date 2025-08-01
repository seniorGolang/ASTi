// @asti name="CyclicStructures" version=1.5
package cyclicstructures

import (
	"context"
	"time"
)

// @asti type=Node validation=strict
type Node struct {
	// @asti field=ID required=true
	ID string `json:"id"`

	// @asti field=Data validation=required
	Data map[string]interface{} `json:"data"`

	// Циклическая ссылка на родителя
	Parent *Node `json:"parent,omitempty"`

	// Слайс циклических ссылок на детей
	Children []*Node `json:"children,omitempty"`

	// Карта циклических ссылок
	Neighbors map[string]*Node `json:"neighbors,omitempty"`

	// @asti field=CreatedAt validation=timestamp
	CreatedAt time.Time `json:"created_at"`

	// @asti field=UpdatedAt validation=timestamp
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

// @asti type=ComplexGraph validation=strict
type ComplexGraph struct {
	// @asti field=Root required=true
	Root *Node `json:"root"`

	// Слайс слайсов узлов (матрица)
	Matrix [][]*Node `json:"matrix,omitempty"`

	// Карта карт узлов
	AdjacencyMap map[string]map[string]*Node `json:"adjacency_map,omitempty"`

	// Слайс карт
	EdgeList []map[string]*Node `json:"edge_list,omitempty"`

	// Карта слайсов
	NodeGroups map[string][]*Node `json:"node_groups,omitempty"`
}

// @asti type=RecursiveType validation=strict
type RecursiveType struct {
	// @asti field=Name required=true
	Name string `json:"name"`

	// Рекурсивная ссылка на себя
	Self *RecursiveType `json:"self,omitempty"`

	// Слайс рекурсивных ссылок
	Children []*RecursiveType `json:"children,omitempty"`

	// Карта рекурсивных ссылок
	Variants map[string]*RecursiveType `json:"variants,omitempty"`
}

// @asti name="CyclicStructureService" timeout=120
type CyclicStructureService interface {
	// @asti method=CreateNode retry=3
	CreateNode(ctx context.Context, data map[string]interface{}) (node *Node, err error)

	// @asti method=BuildGraph timeout=60
	BuildGraph(ctx context.Context, nodes []*Node) (graph *ComplexGraph, err error)

	// @asti method=TraverseGraph retry=5
	TraverseGraph(ctx context.Context, graph *ComplexGraph) (result []*Node, err error)

	// @asti method=CreateRecursive timeout=30
	CreateRecursive(ctx context.Context, name string) (recursive *RecursiveType, err error)

	// @asti method=ProcessMatrix timeout=45
	ProcessMatrix(ctx context.Context, matrix [][]*Node) (result [][]*Node, err error)

	// @asti method=UpdateAdjacencyMap retry=2
	UpdateAdjacencyMap(ctx context.Context, adjacencyMap map[string]map[string]*Node) (result map[string]map[string]*Node, err error)
}
