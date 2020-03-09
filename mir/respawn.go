package mir

import (
	"errors"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/yenkeia/mirgo/common"
	"github.com/yenkeia/mirgo/ut"
)

type RouteInfo struct {
	Location common.Point
	Delay    int
}

func RouteInfoFromText(text string) (*RouteInfo, error) {
	arr := strings.Split(text, ",")
	if len(arr) != 2 && len(arr) != 3 {
		return nil, errors.New("error route text:" + text)
	}

	x, err := strconv.Atoi(arr[0])
	if err != nil {
		return nil, errors.New("error route text x:" + text)
	}
	y, err := strconv.Atoi(arr[1])
	if err != nil {
		return nil, errors.New("error route text y:" + text)
	}

	delay := 0
	if len(arr) == 3 {
		delay, err = strconv.Atoi(arr[2])
		if err != nil {
			return nil, errors.New("error route text delay:" + text)
		}
	}

	return &RouteInfo{Location: common.NewPoint(x, y), Delay: delay}, nil
}

type Respawn struct {
	Info     *common.RespawnInfo
	Monster  *common.MonsterInfo
	Routes   []*RouteInfo
	Count    int
	Map      *Map
	Interval time.Duration
	Elapsed  time.Duration
}

func NewRespawn(m *Map, info *common.RespawnInfo) (*Respawn, error) {
	r := &Respawn{}
	r.Map = m
	r.Info = info
	r.Monster = data.GetMonsterInfoByID(info.MonsterID)
	r.Interval = time.Minute

	err := r.LoadRoutes()
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (r *Respawn) LoadRoutes() error {
	if r.Info.RoutePath == "" {
		return nil
	}

	filename := filepath.Join(settings.RoutePath, r.Info.RoutePath+".txt")
	if !ut.IsFile(filename) {
		return errors.New("Route文件不存在:" + r.Info.RoutePath)
	}

	lines, err := ut.ReadLines(filename)
	if err != nil {
		return errors.New("Route文件读取失败:" + err.Error())
	}

	r.Routes = []*RouteInfo{}
	for _, line := range lines {
		route, err := RouteInfoFromText(line)
		if err != nil {
			return errors.New("Route文件解析失败:" + err.Error())
		}
		r.Routes = append(r.Routes, route)
	}

	return nil
}

func (r *Respawn) Process(dt time.Duration) {
	r.Elapsed += dt
	if r.Elapsed < r.Interval {
		return
	}

	r.Elapsed = 0
	r.Spawn()
}

func (r *Respawn) Spawn() {
	for i := r.Count; i < r.Info.Count; i++ {
		r.SpawnOne()
	}
}

func (r *Respawn) SpawnOne() bool {
	for i := 0; i < 10; i++ {
		x := r.Info.LocationX + ut.RandomInt(-r.Info.Spread, r.Info.Spread)
		y := r.Info.LocationY + ut.RandomInt(-r.Info.Spread, r.Info.Spread)

		if !r.Map.ValidPointXY(x, y) {
			continue
		}

		m := NewMonster(r.Map, common.NewPoint(x, y), r.Monster)
		m.CurrentDirection = common.MirDirection(r.Info.Direction)
		r.Map.AddObject(m)

		m.BroadcastInfo()
		m.BroadcastHealthChange()

		r.Count++

		return true
	}

	return false
}
