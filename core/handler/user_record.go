package handler

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Yosorable/ms-admin/global"
	pb "github.com/Yosorable/ms-shared/protoc_gen/admin"
)

func packFieldForSQL(field string) string {
	return "`" + field + "`"
}

func getTableName(tableOption *pb.UserRecordTableOption) (tableName string, err error) {
	if tableOption == nil {
		return "", errors.New("UserRecordTableOption is null")
	}
	adminServiceName := global.CONFIG.ServiceName
	callerServiceName := tableOption.ServiceName
	tableTag := tableOption.TableTag
	if callerServiceName == "" || tableTag == "" {
		return "", errors.New("ServiceName or TableTag is empty")
	}
	return fmt.Sprintf(
		"%s_%s_%s",
		strings.TrimPrefix(adminServiceName, "ms_"),
		strings.TrimPrefix(callerServiceName, "ms_"),
		tableTag,
	), nil
}

func CreateUserRecordTableIfNotExist(ctx context.Context, req *pb.CreateUserRecordTableIfNotExistRequest) (*pb.CreateUserRecordTableIfNotExistReply, error) {
	tableOption := req.GetTableOption()
	tableName, err := getTableName(tableOption)
	if err != nil {
		return nil, NewStatusError(err)
	}
	foreignIdName := tableOption.ForeignIdName

	err = global.DATABASE.Exec(`
        CREATE TABLE IF NOT EXISTS ` + packFieldForSQL(tableName) + ` (
            ` + packFieldForSQL("user_id") + ` int(0) NOT NULL,
            ` + packFieldForSQL(foreignIdName) + ` int(0) NOT NULL,
            ` + packFieldForSQL("created_at") + ` datetime(0) DEFAULT CURRENT_TIMESTAMP,
            ` + packFieldForSQL("updated_at") + ` datetime(0) DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
            PRIMARY KEY (` + packFieldForSQL("user_id") + `, ` + packFieldForSQL(foreignIdName) + `) USING BTREE
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
	tableOption := req.GetTableOption()
	tableName, err := getTableName(tableOption)
	if err != nil {
		return nil, NewStatusError(err)
	}
	foreignIdName := tableOption.ForeignIdName

	err = global.DATABASE.Exec(`
    INSERT INTO `+packFieldForSQL(tableName)+` 
        (`+packFieldForSQL("user_id")+`, `+packFieldForSQL(foreignIdName)+`) 
        values (?, ?)
        ON DUPLICATE KEY 
        UPDATE updated_at = CURRENT_TIMESTAMP();
    `, req.UserId, req.ForeignItemId).Error
	if err != nil {
		return nil, NewStatusError(err)
	}

	return &pb.CreateOrUpdateUserRecordReply{}, nil
}

func DeleteUserRecord(ctx context.Context, req *pb.DeleteUserRecordRequest) (*pb.DeleteUserRecordReply, error) {
	tableOption := req.GetTableOption()
	tableName, err := getTableName(tableOption)
	if err != nil {
		return nil, NewStatusError(err)
	}
	foreignIdName := tableOption.ForeignIdName

	err = global.DATABASE.Exec(`
    DELETE FROM `+packFieldForSQL(tableName)+`
        WHERE `+packFieldForSQL("user_id")+` = ?
        AND `+packFieldForSQL(foreignIdName)+` = ?
    `, req.UserId, req.ForeignItemId).Error
	if err != nil {
		return nil, NewStatusError(err)
	}

	return &pb.DeleteUserRecordReply{}, nil
}
