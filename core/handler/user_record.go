package handler

import (
	"context"
	"fmt"
	"strings"

	"github.com/Yosorable/ms-admin/global"
	pb "github.com/Yosorable/ms-shared/protoc_gen/admin"
)

func packFiledForSQL(field string) string {
	return "`" + field + "`"
}

func CreateUserRecordTableIfNotExist(ctx context.Context, req *pb.CreateUserRecordTableIfNotExistRequest) (*pb.CreateUserRecordTableIfNotExistReply, error) {
	if req.GetTableOption() == nil {
		return nil, NotNULLError("table_option")
	}
	tableOption := req.GetTableOption()
	tableName := fmt.Sprintf(
		"%s_%s_%s",
		strings.TrimPrefix(global.CONFIG.ServiceName, "ms_"),
		strings.TrimPrefix(tableOption.ServiceName, "ms_"),
		tableOption.TableTag,
	)
	foreignIdName := tableOption.ForeignIdName

	err := global.DATABASE.Exec(`
        CREATE TABLE IF NOT EXISTS ` + packFiledForSQL(tableName) + ` (
            ` + packFiledForSQL("user_id") + ` int(0) NOT NULL,
            ` + packFiledForSQL(foreignIdName) + ` int(0) NOT NULL,
            ` + packFiledForSQL("created_at") + ` datetime(0) DEFAULT CURRENT_TIMESTAMP,
            ` + packFiledForSQL("updated_at") + ` datetime(0) DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
            PRIMARY KEY (` + packFiledForSQL("user_id") + `, ` + packFiledForSQL(foreignIdName) + `) USING BTREE
        ) ENGINE = InnoDB CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;
    `).Error
	if err != nil {
		return nil, NewStatusError(err)
	}
	return &pb.CreateUserRecordTableIfNotExistReply{}, nil
}

func QueryUserRecord(ctx context.Context, req *pb.QueryUserRecordRequest) (*pb.QueryUserRecordReply, error) {
	return nil, nil
}

func CreateOrUpdateUserRecord(ctx context.Context, req *pb.CreateOrUpdateUserRecordRequest) (*pb.CreateOrUpdateUserRecordReply, error) {
	return nil, nil
}

func DeleteUserRecord(ctx context.Context, req *pb.DeleteUserRecordRequest) (*pb.DeleteUserRecordReply, error) {
	return nil, nil
}
