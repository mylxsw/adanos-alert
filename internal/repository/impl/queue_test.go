package impl_test

import (
	"testing"
	"time"

	"github.com/mylxsw/adanos-alert/internal/repository"
	"github.com/mylxsw/adanos-alert/internal/repository/impl"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
)

type QueueTestSuit struct {
	suite.Suite
	repo repository.QueueRepo
}

func (q *QueueTestSuit) TearDownTest() {
	q.NoError(q.repo.Delete(bson.M{}))
}

func (q *QueueTestSuit) SetupTest() {
	db, err := Database()
	q.NoError(err)

	q.repo = impl.NewQueueRepo(db)
}

func (q *QueueTestSuit) TestEnqueueDequeue() {
	// test empty queue
	_, err := q.repo.Dequeue()
	q.Error(err)
	q.Equal(repository.ErrNotFound, err)

	// test enqueue
	item := repository.QueueJob{
		Name:    "action",
		Payload: "{}",
	}

	// add a item to queue
	insertID, err := q.repo.Enqueue(item)
	q.NoError(err)
	q.NotEmpty(insertID)

	{
		time.Sleep(10 * time.Millisecond)

		// test dequeue one item
		item2, err := q.repo.Dequeue()
		q.NoError(err)
		q.EqualValues(item.Name, item2.Name)
		q.EqualValues(repository.QueueItemStatusRunning, item2.Status)

		// test empty queue
		{
			_, err := q.repo.Dequeue()
			q.Error(err)
			q.Equal(repository.ErrNotFound, err)
		}

		// test item's status changed to running after dequeue
		item21, err := q.repo.Get(item2.ID)
		q.NoError(err)
		q.EqualValues(repository.QueueItemStatusRunning, item21.Status)

		// test queue item count
		c21, err := q.repo.Count(bson.M{"status": repository.QueueItemStatusWait})
		q.NoError(err)
		q.EqualValues(0, c21)
	}

}

func TestQueueRepo(t *testing.T) {
	suite.Run(t, new(QueueTestSuit))
}
