package waitlist

import (
	"option-dance/cmd/config"
	"option-dance/core"
	"option-dance/pkg/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

var mutex sync.Mutex

func New(waitListStore core.WaitListStore) core.WaitListService {
	return &service{waitListStore: waitListStore}
}

type service struct {
	waitListStore core.WaitListStore
}

func (s *service) Create(waitlist core.WaitList, c *gin.Context) {
	mutex.Lock()
	defer mutex.Unlock()
	var oldWaitList core.WaitList
	//mixin environment
	if waitlist.Type == 0 {
		if len(waitlist.MixinID) > 0 {
			config.Db.Model(core.WaitList{}).Where("mixin_id=?", waitlist.MixinID).Find(&oldWaitList)
			if oldWaitList.ID > 0 {
				http.OkWithMessage("已加入waitlist", c)
				return
			} else {
				waitlist.Wid = s.genWid()
				waitlist.IP = c.ClientIP()
				waitlist.CreatedAt = time.Now()
				waitlist.UpdatedAt = time.Now()
				if err := s.waitListStore.Create(waitlist); err != nil {
					http.FailWithMsg(err.Error(), c)
					return
				} else {
					if waitlist.InviterWid > 0 {
						s.IncrInvitedCount(waitlist.InviterWid)
					}
					http.Ok(c)
					return
				}
			}
		} else {
			http.FailWithMsg("未授权，无法加入waitlist", c)
			return
		}
	}
	//Non-mixin environment
	if waitlist.Type == 1 {
		if len(waitlist.Email) > 0 {
			config.Db.Model(core.WaitList{}).Where("email=?", waitlist.Email).Find(&oldWaitList)
			if oldWaitList.ID > 0 {
				http.OkWithMessage("已加入waitlist", c)
				return
			} else {
				waitlist.Wid = s.genWid()
				waitlist.IP = c.ClientIP()
				waitlist.CreatedAt = time.Now()
				waitlist.UpdatedAt = time.Now()
				if err := config.Db.Model(core.WaitList{}).Create(&waitlist).Error; err != nil {
					http.FailWithMsg(err.Error(), c)
					return
				} else {
					if waitlist.InviterWid > 0 {
						s.IncrInvitedCount(waitlist.InviterWid)
					}
					http.Ok(c)
					return
				}
			}
		} else {
			http.FailWithMsg("邮箱不能为空", c)
			return
		}
	}
}

func (s *service) genWid() (wid int64) {
	var w core.WaitList
	config.Db.Model(core.WaitList{}).Order("created_at desc").First(&w)
	if w.ID > 0 && w.Wid > 0 {
		wid = w.Wid + 42
	} else {
		wid = 30000
	}
	return
}

func (s *service) IncrInvitedCount(wid int64) {
	_ = s.waitListStore.IncrInvitedCount(wid)
	return
}

func (s *service) GetRankInfoByMid(mid string) (info core.WaitListRankInfo) {
	var records []core.WaitList
	config.Db.Model(core.WaitList{}).Raw("select mixin_id,invite_count from waitlist order by invite_count desc, created_at asc").Find(&records)
	for i, e := range records {
		if e.MixinID == mid {
			info.Rank = int64(i + 1)
			info.InviteCount = e.InviteCount
			break
		}
	}
	return
}

func (s *service) GetWid(mid string, email string) (wid int64, err error) {
	var w *core.WaitList
	if len(mid) > 0 {
		w, err = s.waitListStore.GetByMixinId(mid)
		if err != nil {
			return 0, err
		}
	} else {
		w, err = s.waitListStore.GetByEmail(email)
		if err != nil {
			return 0, err
		}
	}
	if w.ID > 0 {
		wid = w.Wid
	} else {
		wid = 0
	}
	return
}
