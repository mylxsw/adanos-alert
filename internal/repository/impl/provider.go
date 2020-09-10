package impl

import (
	"time"

	"github.com/mylxsw/adanos-alert/configs"
	"github.com/mylxsw/adanos-alert/internal/repository"
	"github.com/mylxsw/asteria/log"
	"github.com/mylxsw/coll"
	"github.com/mylxsw/container"
	"github.com/mylxsw/glacier/cron"
	"github.com/mylxsw/glacier/infra"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ServiceProvider struct{}

func (s ServiceProvider) Register(app container.Container) {
	app.MustSingleton(NewSequenceRepo)
	app.MustSingleton(NewKVRepo)
	app.MustSingleton(NewMessageRepo)
	app.MustSingleton(NewMessageGroupRepo)
	app.MustSingleton(NewUserRepo)
	app.MustSingleton(NewRuleRepo)
	app.MustSingleton(NewQueueRepo)
	app.MustSingleton(NewTemplateRepo)
	app.MustSingleton(NewDingdingRobotRepo)
	app.MustSingleton(NewLockRepo)
	app.MustSingleton(NewAgentRepo)
	app.MustSingleton(NewAuditLogRepo)
	app.MustSingleton(NewRecoveryRepo)
}

func (s ServiceProvider) Boot(app infra.Glacier) {
	app.Cron(func(cr cron.Manager, cc container.Container) error {
		return cc.Resolve(func(
			kvRepo repository.KVRepo,
			groupRepo repository.MessageGroupRepo,
			msgRepo repository.MessageRepo,
			auditRepo repository.AuditLogRepo,
			conf *configs.Config,
		) {
			_ = cr.Add("kv_repository_gc", "@hourly", func() {
				if err := kvRepo.GC(); err != nil {
					log.Errorf("kv kvRepo gc failed: %v", err)
				}
			})

			_ = cr.Add("remove_expired_audit_log", "@midnight", func() {
				deadLineDate := time.Now().AddDate(0, 0, -7*2)
				log.Infof("clear expired audit logs before %v", deadLineDate)

				if err := auditRepo.Delete(bson.M{"created_at": bson.M{"$lt": deadLineDate}}); err != nil {
					log.Errorf("clear expired audit logs before %v failed: %v", deadLineDate, err)
				}
			})

			if conf.KeepPeriod > 0 {
				_ = cr.Add("remove_expired_messages", "@midnight", func() {
					expiredMessageGC(conf, msgRepo, groupRepo)
				})

				// 每次重启服务时，自动触发一次GC
				expiredMessageGC(conf, msgRepo, groupRepo)
			}
		})
	})
}

// expiredMessageGC 清理过期的message/message_group
func expiredMessageGC(conf *configs.Config, msgRepo repository.MessageRepo, groupRepo repository.MessageGroupRepo) {
	deadLineDate := time.Now().AddDate(0, 0, -conf.KeepPeriod)
	log.Infof("clear expired/canceled messages and groups before %v", deadLineDate)

	// 删除过期或者取消发送的 messages
	if err := msgRepo.Delete(bson.M{
		"status":     bson.M{"$in": []repository.MessageStatus{repository.MessageStatusCanceled, repository.MessageStatusExpired}},
		"created_at": bson.M{"$lt": deadLineDate},
	}); err != nil {
		log.Errorf("remove expired/canceled messages before %v failed: %v", deadLineDate, err)
	}

	// 删除过期的 groups
	// 1. 查询过期的 groups
	// 2. 删除过期分组关联的所有 messages
	// 3. 删除过期分组
	groups, err := groupRepo.Find(bson.M{"created_at": bson.M{"$lt": deadLineDate}})
	if err != nil {
		log.Errorf("query expired message groups before %v failed: %v", deadLineDate, err)
		return
	}

	groupIds := coll.MustNew(groups).Map(func(grp repository.MessageGroup) primitive.ObjectID {
		return grp.ID
	}).Items()
	if err := msgRepo.Delete(bson.M{"group_ids": bson.M{"$in": groupIds}}); err != nil {
		log.Errorf("remove messages in group_ids before %v failed: %v", deadLineDate, err)
		return
	}

	if err := groupRepo.Delete(bson.M{"_id": bson.M{"$in": groupIds}}); err != nil {
		log.Errorf("remove messages in group_ids before %v failed: %v", deadLineDate, err)
		return
	}
}
