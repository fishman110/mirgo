package mir

import (
	"time"

	"github.com/yenkeia/mirgo/common"
	"github.com/yenkeia/mirgo/proto/server"
)

type BaseStats struct {
	MinAC    uint16
	MaxAC    uint16
	MinMAC   uint16
	MaxMAC   uint16
	MinDC    uint16
	MaxDC    uint16
	MinMC    uint16
	MaxMC    uint16
	MinSC    uint16
	MaxSC    uint16
	Accuracy uint8
	Agility  uint8
}

type IDObject interface {
	GetID() uint32
}

type ISimpleMapObject interface {
	IDObject
	GetMap() *Map
	GetRace() common.ObjectType
	Broadcast(interface{})
}

type IProcessObject interface {
	IDObject
	Process(dt time.Duration)
}

type IMapObject interface {
	ISimpleMapObject
	GetName() string
	GetLevel() int
	GetPoint() common.Point
	GetCell() *Cell
	BroadcastHealthChange()
	BroadcastInfo()
	Spawned()
	GetDirection() common.MirDirection
	GetBaseStats() BaseStats
	IsAttackTarget(IMapObject) bool
	IsFriendlyTarget(IMapObject) bool
	IsDead() bool
	IsUndead() bool
	IsBlocking() bool
	AttackMode() common.AttackMode
	AddBuff(*Buff)
	ApplyPoison(*Poison, IMapObject)
	AddPlayerCount(n int)
	GetPlayerCount() int
	Attacked(attacker IMapObject, damage int, dtype common.DefenceType, damageWeapon bool) int
	GetMapObject() *MapObject
}

type ILifeObject interface {
	ISimpleMapObject
	GetHP() int
	GetMaxHP() int
	SetHP(uint32)
	ChangeHP(int)
}

type MapObject struct {
	ID               uint32
	Name             string
	NameColor        common.Color
	Map              *Map
	CurrentLocation  common.Point
	CurrentDirection common.MirDirection
	Dead             bool
	PlayerCount      int // 记录在DataRange内有多少个玩家
	InSafeZone       bool
}

func (m *MapObject) GetMapObject() *MapObject {
	return m
}
func (m *MapObject) UpdateInSafeZone() {
	if m.Map != nil {
		m.InSafeZone = m.Map.GetSafeZone(m.CurrentLocation) != nil
	} else {
		m.InSafeZone = false
	}
}

func IMapObject_Spawned(m IMapObject) {
	m.GetMapObject().UpdateInSafeZone()
	m.BroadcastInfo()
	m.BroadcastHealthChange()
}

func IMapObject_BroadcastHealthChange(m ILifeObject) {
	if m.GetRace() != common.ObjectTypePlayer && m.GetRace() != common.ObjectTypeMonster {
		return
	}

	// TODO RevTime

	percent := byte(float32(m.GetHP()) / float32(m.GetMaxHP()) * 100)

	msg := &server.ObjectHealth{
		ObjectID: m.GetID(),
		Percent:  percent,
		Expire:   5,
	}

	m.Broadcast(msg)
}
