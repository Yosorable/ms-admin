package handler

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Yosorable/ms-admin/global"
	pb "github.com/Yosorable/ms-shared/protoc_gen/admin"
	"gorm.io/gorm"
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
		strings.TrimPrefix(adminServiceName, "ms-"),
		strings.TrimPrefix(callerServiceName, "ms-"),
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
	tableOption := req.GetTableOption()
	tableName, err := getTableName(tableOption)
	if err != nil {
		return nil, NewStatusError(err)
	}
	foreignIdName := tableOption.ForeignIdName

	orderByExp := ""
	if len(req.GetOrderByList()) > 0 {
		orderByExp = "ORDER BY "
		for _, ele := range req.GetOrderByList() {
			order := ""
			if ele.Order == pb.OrderBy_DESC {
				order = "DESC"
			}
			orderByExp += fmt.Sprintf("`%s` %s,", ele.OrderByFiledName, order)
		}
		orderByExp = strings.TrimSuffix(orderByExp, ",")
	}

	pagnationExp := ""
	if req.GetPageOption() != nil {
		pageOpt := req.GetPageOption()
		pagnationExp = fmt.Sprintf("LIMIT %d OFFSET %d", pageOpt.Limit, pageOpt.Offset)
	}
	db := global.DATABASE
	resSql := ""

	countSql := ""
	if len(req.GetForeignIdList()) > 0 {
		resSql = db.ToSQL(func(tx *gorm.DB) *gorm.DB {
			return db.Raw(`
                SELECT *
                    FROM `+packFieldForSQL(tableName)+`
                    WHERE `+packFieldForSQL("user_id")+` = ?
                    AND `+packFieldForSQL(foreignIdName)+` IN ?
                    `+orderByExp+`
                    `+pagnationExp+`
                `, req.GetUserId(), req.GetForeignIdList())

		})
		countSql = db.ToSQL(func(tx *gorm.DB) *gorm.DB {
			return db.Raw(`
                SELECT count(*) total
                    FROM `+packFieldForSQL(tableName)+`
                    WHERE `+packFieldForSQL("user_id")+` = ?
                    AND `+packFieldForSQL(foreignIdName)+` IN ?
                `, req.GetUserId(), req.GetForeignIdList())

		})
	} else {
		resSql = db.ToSQL(func(tx *gorm.DB) *gorm.DB {
			return db.Raw(`
                SELECT *
                    FROM `+packFieldForSQL(tableName)+`
                    WHERE `+packFieldForSQL("user_id")+` = ?
                    `+orderByExp+`
                    `+pagnationExp+`
                `, req.GetUserId())

		})
		countSql = db.ToSQL(func(tx *gorm.DB) *gorm.DB {
			return db.Raw(`
                SELECT count(*) total
                    FROM `+packFieldForSQL(tableName)+`
                    WHERE `+packFieldForSQL("user_id")+` = ?
                `, req.GetUserId())
		})
	}

	execSql := `
        SELECT * 
            FROM 
            (` + resSql + `) res,
            (` + countSql + `) cnt
    `
	resMap := []map[string]any{}
	err = db.Raw(execSql).Scan(&resMap).Error
	if err != nil {
		return nil, NewStatusError(err)
	}
	res := &pb.QueryUserRecordReply{
		Total:   0,
		Records: []*pb.UserRecord{},
	}
	if len(resMap) > 0 {
		res.Total = int32(resMap[0]["total"].(int64))
	}
	for _, ele := range resMap {
		record := &pb.UserRecord{
			UserId: req.GetUserId(),
		}
		dataTransformErr := NewStatusError(errors.New("data transform error"))
		if foreignIdVal, success := ele[foreignIdName].(int32); success {
			record.ForeignId = int32(foreignIdVal)
		} else {
			return nil, dataTransformErr
		}
		if createdAt, success := ele["created_at"].(time.Time); success {
			record.CreatedAt = createdAt.Format("2006-01-02 15:04:05")
		} else {
			return nil, dataTransformErr
		}
		if updatedAt, success := ele["updated_at"].(time.Time); success {
			record.UpdatedAt = updatedAt.Format("2006-01-02 15:04:05")
		} else {
			return nil, dataTransformErr
		}
		res.Records = append(res.Records, record)
	}
	return res, nil
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
