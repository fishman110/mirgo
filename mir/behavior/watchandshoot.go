package behavior

import (
	"time"

	"github.com/yenkeia/mirgo/common"
	"github.com/yenkeia/mirgo/mir"
	"github.com/yenkeia/mirgo/proto/server"
)

//
type WatchAndShootNode struct {
	Node
}

func WatchAndShoot() INode {
	return &WatchAndShootNode{}
}

func (n *WatchAndShootNode) Visit(c *BT) {

	if n.Status() == READY {
		if c.Monster.Target == nil || c.Monster.Target.IsDead() {
			n.status = FAILED
			return
		} else {
			n.status = RUNNING
		}
	}

	if n.Status() == RUNNING {
		if c.Monster.Target == nil || c.Monster.Target.IsDead() {
			n.status = SUCCESS
			c.Monster.Target = nil
		} else if c.Monster.CanAttack() {

			const AttackRange = 10
			m := c.Monster

			if mir.InRange(m.CurrentLocation, m.Target.GetPoint(), AttackRange) {
				n.status = FAILED
				return
			}

			m.CurrentDirection = mir.DirectionFromPoint(m.CurrentLocation, m.Target.GetPoint())
			m.Broadcast(&server.ObjectRangeAttack{
				ObjectID:  m.GetID(),
				Direction: m.CurrentDirection,
				Location:  m.CurrentLocation,
				TargetID:  m.Target.GetID(),
			})

			damage := m.GetAttackPower(int(m.MaxDC), int(m.MaxDC))
			if damage == 0 {
				return
			}

			now := time.Now()
			m.AttackTime = now.Add(time.Duration(m.AttackSpeed) * time.Millisecond)

			delay := mir.MaxDistance(m.CurrentLocation, m.Target.GetPoint())*50 + 500
			target := m.Target

			m.ActionList.PushDelayAction(mir.DelayedTypeDamage, delay, func() {
				m.CompleteAttack(target, damage, common.DefenceTypeACAgility)
			})
		} else {
			n.status = FAILED
		}
	}
}
