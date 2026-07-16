package commonvalueobjects

import (
	"strconv"

	commondomainerrors "github.com/PurpleSavage/monekai-server/modules/shared/common/domain/errors"
)

type PaginationVO struct {
	Page     	int 
	Limit 		int 
}

func CreatePaginationVO(
	page string,
	limit string,
) (*PaginationVO, error){
	pageNumber,err:= strconv.Atoi(page)
	if err!=nil{
		return nil, commondomainerrors.NewValidationError(
			"page",
			"page must be a number",
		)
	}
	limitNumber, err := strconv.Atoi(limit)
	if err != nil {
		return nil, commondomainerrors.NewValidationError(
			"limit",
			"limit must be a number",
		)
	}
	return &PaginationVO{
		Page:  pageNumber,
		Limit: limitNumber,
	} , nil
}
