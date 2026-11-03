package knowledge

type (
	CreateKnowledgeBaseOption struct {
		Id      string
		Name    string
		Desc    string
		modelId string
	}

	DeleteKnowledgeBaseOption struct {
		Id         string
		SoftDelete bool
	}

	UpdateKnowledgeBaseOption struct {
		Id      string
		Name    *string
		Desc    *string
		modelId *string
	}

	KnowledgeStore struct {
		Id        string
		Name      string
		Desc      string
		modelId   string
		createdAt string
		updatedAt string
		deletedAt string
	}

	KnowledgeBase interface {
		CreateKnowledgeBase(option CreateKnowledgeBaseOption)
		GetKnowledgeBase(id string) KnowledgeStore
		DeleteKnowledgeBase(option DeleteKnowledgeBaseOption)
		UpdateKnowledgeBase(option UpdateKnowledgeBaseOption)
	}
)
