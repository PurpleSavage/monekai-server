package commonports

import commonentities "github.com/PurpleSavage/monekai-server/modules/shared/common/domain/entities"
type Observer interface {
	ReceiveEvent(event commonentities.Event[any])
}