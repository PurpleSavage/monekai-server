package commonports

import commonentities "github.com/PurpleSavage/monekai-server/modules/shared/common/domain/entities"

type ObserverBucketPort interface {
	AddObserver(observer Observer, eventname string)
	RemoveObserver(observer Observer, eventname string)
	NotifyObservers(event commonentities.Event, eventname string)
}
