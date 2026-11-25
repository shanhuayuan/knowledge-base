package knowledge

import (
	"app/global"
	"app/model"
	"context"
	"strings"

	"github.com/shanhuayuan/object"
)

func kbId2collectionName(kbId string) string {
	return "kb_" + strings.ReplaceAll(kbId, "-", "_")
}

func CreateKnowledgeBase(kb *model.KnowledgeBase) error {
	return global.OverallRuntime.DataStore.CreateEntry(kb)
}

func DeleteKnowledgeBase(ctx context.Context, kb *model.KnowledgeBase) error {
	// 1. 从数据库删除知识库记录
	if err := global.OverallRuntime.DataStore.DeleteEntry(kb); err != nil {
		return err
	}

	// 2. 从数据库删除知识库文件记录
	document := &model.Document{KbId: kb.Id}

	if err := global.OverallRuntime.DataStore.DeleteEntry(document); err != nil {
		return err
	}

	// 3. 从向量数据库中删除知识库对应的集合
	name := kbId2collectionName(kb.Id)

	if err := global.OverallRuntime.VectorStore.DropCollection(ctx, name); err != nil {
		return err
	}

	// 4. 从对象存储中删除对应的文件
	option := object.Object{BucketName: "knowledge", ObjectKey: kb.Id}
	global.OverallRuntime.ObjectStore.DeleteObjects(ctx, option)

	return nil
}

func ListKnowledgeBase(kbs *[]model.KnowledgeBase) error {
	return global.OverallRuntime.DataStore.FindEntry(&model.KnowledgeBase{}, kbs)
}

func IsKnowledgeBaseExist(kb *model.KnowledgeBase) (bool, error) {
	return global.OverallRuntime.DataStore.IsEntryExist(kb)
}

func GetKnowledgeBase(kb *model.KnowledgeBase) error {
	return global.OverallRuntime.DataStore.TakeEntry(kb)
}
