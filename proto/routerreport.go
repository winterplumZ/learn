package proto

import (
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"staticzeng.com/config"
	"staticzeng.com/packet"
	"staticzeng.com/tools"
)

var (
	db *sqlx.DB
)

type RouteInfo struct {
	Name    string `json:"router_name"`
	Version string `json:"router_version"`
	Sn      string `json:"router_sn"`
	RMac    string `json:"router_mac"`
	CMac    string `json:"mac"`
	User    string `json:"user"`
	CpuUse  int    `json:"router_cpu_use"`
	MemUse  int    `json:"router_mem_use"`
}

type Device struct {
	Name string `json:"devname"`
	Mac  string `json:"mac"`
	Ip   string `json:"ip"`
}

func (dev *Device) save(rmac string) {
	db.Exec(fmt.Sprintf("insert into device (rmac, name, mac, ip, time) values(\"%s\", \"%s\", \"%s\", \"%s\", now())",
		rmac,
		dev.Name,
		dev.Mac,
		dev.Ip))
}

type UrlStatus struct {
	Url           string  `json:"url"`
	NameLookup    float64 `json:"nameloopup"`
	StarTrans     float64 `json:"starttrans"`
	TotalTime     float64 `json:"totaltime"`
	DownloadSize  int     `json:"downloadsize"`
	DownloadSpeed float64 `json:"downloadspeed"`
	Code          int     `json:"curlcode"`
	HttpCode      int     `json:"httpcode"`
}

func (url *UrlStatus) save(rmac string) {
	db.Exec(fmt.Sprintf("insert into url (mac, url, namelookup, starttrans, totaltime, downloadsize, downloadspeed,code, httpcode, time) values(\"%s\", \"%s\", %f, %f, %f, %d, %f, %d, %d, now())",
		rmac,
		url.Url,
		url.NameLookup,
		url.StarTrans,
		url.TotalTime,
		url.DownloadSize,
		url.DownloadSpeed,
		url.Code,
		url.HttpCode))
}

type Salt struct {
	Ip    string `json:"ip"`
	Times int    `json:"times"`
}

func (sa *Salt) save(rmac string) {
	url := tools.IptoUrl(sa.Ip)
	if url == "" {
		return
	}
	db.Exec(fmt.Sprintf("insert into dpi (mac, url, times, time) values(\"%s\", \"%s\", %d, now())",
		rmac,
		url,
		sa.Times))
}

type RouterData struct {
	Route    RouteInfo   `json:"Router"`
	Devices  []Device    `json:"devices"`
	Urls     []UrlStatus `json:"urls"`
	Captures []Salt      `json:"capture"`
}

func (ri *RouteInfo) saveRI() (err error) {
	stmt, _ := db.Preparex("select mac from router where mac=?")
	defer stmt.Close()
	var obj string
	err = stmt.Get(&obj, ri.RMac)
	if obj != "" {
		_, err = db.Exec(fmt.Sprintf("update router set version=\"%s\", name=\"%s\", catmac=\"%s\", catuser=\"%s\", cpu=\"%d\", mem=\"%d\", time=now() where mac=\"%s\"",
			ri.Version,
			ri.Name,
			ri.CMac,
			ri.User,
			ri.CpuUse,
			ri.MemUse,
			ri.RMac))
	} else {
		_, err = db.Exec(fmt.Sprintf("insert into router values(\"%s\", \"%s\", \"%s\", \"%s\",\"%s\", \"%s\", %d, %d, now())",
			ri.RMac,
			ri.Sn,
			ri.Version,
			ri.Name,
			ri.CMac,
			ri.User,
			ri.CpuUse,
			ri.MemUse))
	}
	return
}

func checkErr(prefix string, err error) {
	if err != nil {
		fmt.Println(prefix, err)
	}
}

func (rd *RouterData) savedb() (err error) {
	err = (&rd.Route).saveRI()
	checkErr("save RI", err)
	for _, d := range rd.Devices {
		(&d).save(rd.Route.RMac)
	}
	for _, u := range rd.Urls {
		(&u).save(rd.Route.RMac)
	}
	for _, cap := range rd.Captures {
		(&cap).save(rd.Route.RMac)
	}
	return
}

type RouterProto struct {
}

func init() {
	cfg := config.LoadConfig()
	fmt.Println(cfg)
	db, _ = sqlx.Connect("mysql", fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8",
		cfg.DB.User,
		cfg.DB.Pwd,
		cfg.DB.Ip,
		cfg.DB.Port,
		cfg.DB.DbName))
	//	db.SetMaxIdleConns(5)
	//	db.SetMaxOpenConns(20)
	registeProto(&RouterProto{})
}

func (*RouterProto) Handle(pkt *packet.Packet) (err error) {
	var routerdata *RouterData
	err = json.Unmarshal(pkt.Body, &routerdata)
	if err != nil {
		return
	}
	s, _ := json.Marshal(&routerdata)
	fmt.Println(string(s))
	routerdata.savedb()
	return nil
}

func (*RouterProto) Code() (code uint8) {
	code = 0XC7
	return
}
