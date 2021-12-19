package main

import (
	"bufio"
	"bytes"
	"context"
	"crypto/md5"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	core "./service/line"
	sqlogin "./service/secondaryqrcodeloginservice"
	thrift "./service/thrift"
	"github.com/imroc/req"
	"github.com/kardianos/osext"
	"github.com/leekchan/timeutil"
	"github.com/shirou/gopsutil/mem"
)

type (
	tagdata struct {
		S string `json:"S"`
		E string `json:"E"`
		M string `json:"M"`
	}
	mentions struct {
		MENTIONEES []struct {
			Start string `json:"S"`
			End   string `json:"E"`
			Mid   string `json:"M"`
		} `json:"MENTIONEES"`
	}
	LINE struct {
		qrLogin   LoginRequest
		AuthToken string
		AppName   string
		UserAgent string
		Host      string
		MID       string
		Limited   bool
		Akick         int
		Ainvite       int
		Acancel      int
		Ckick         int
		Cinvite       int
		Ccancel     int
		SHani           int
		Count     int32
		Revision  int64
		GRevision int64
		IRevision int64
		Squads    []string
		Backup    []string
	}
	DATA struct {
		SelfToken   string              `json:"SelfToken"`
		SelfStatus  bool                `json:"SelfStatus"`
		Prefix      string              `json:"Prefix"`
		Setkey      string              `json:"Setkey"`
		Rname       string              `json:"Rname"`
		Authoken    []string            `json:"Authoken"`
		Bot         []string            `json:"Bot"`
		Buyer       []string            `json:"Buyer"`
		Owner       []string            `json:"Owner"`
		Master      []string            `json:"Master"`
		Admin       []string            `json:"Admin"`
		Blacklist   []string            `json:"Blacklist"`
		Fucklist    []string            `json:"Fucklist"`
		Premlist    []string            `json:"Premlist"`
		Whitelist    []string            `json:"Whitelist"`
		SquadBots []string            `json:"Backuplist"`
		ForceInvite bool                `json:"ForceQr"`
		ForceJoinqr bool                `json:"ForceInv"`
		VictimMode    bool                `json:"VictimMode"`
		FastMode    bool                `json:"FastMode"`
		QrMode    bool                `json:"QrMode"`
		MixMode    bool              `json:"MixMode"`
		KillMode    bool                 `json:"KillMode"`
		NukeJoin    bool                `json:"NukeJoin"`
		AutoPurge   bool              `json:"AutoPurge"`
		AutoPro    bool                   `json:"AutoPro"`
		AntiTag       bool                `json:"AntiTag"`
		Identict    bool                   `json:"Identict"`
		StayAjs   map[string][]string `json:"AjsStay"`
		StayGroup   map[string][]string `json:"BotStay"`
		GroupOwn    map[string][]string `json:"OwnStay"`
		GroupAdm    map[string][]string `json:"AdmStay"`
		GroupBan    map[string][]string `json:"BanGroup"`
		GroupName map[string]string `json:"NameGroup"`
		ProName map[string]int `json:"ProName"`
		ProQr       []string            `json:"ProQr"`
		ProKick     []string            `json:"ProKick"`
		ProInvite   []string            `json:"ProInvite"`
		ProCancel   []string            `json:"ProCancel"`
		ProJoin   []string            `json:"ProJoin"`
		Message     struct {
			Welcome  string `json:"MWelcome"`
			Respon string `json:"MRespon"`
			Ban    string `json:"MUnban"`
			Bye    string `json:"MBye"`
			Flag   string `json:"MFlag"`
			Sider  string `json:"Msider"`
			Fresh  string `json:"MFfresh"`
			Limit  string `json:"MLimit"`
		} `json:"Message"`
		Command     struct {
			Banlist   string `json:"CBanlist"`
			Clearban   string `json:"CClearban"`
			Count   string `json:"CCount"`
			Kick   string `json:"CKick"`
			Leave   string `json:"CLeave"`
			Outall   string `json:"COutall"`
			Out   string `json:"COut"`
			Respon   string `json:"CRespon"`
			Setting   string `json:"CSetting"`
			Set   string `json:"CSet"`
			Speed   string `json:"CSpeed"`
			Status   string `json:"CStatus"`
			Unsend   string `json:"CUnsend"`
		} `json:"Command"`
	}
)

type ProfileCoverStruct struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Result  struct {
		HomeId       string `json:"homeId"`
		HomeType     string `json:"homeType"`
		HasNewPost   bool   `json:"hasNewPost"`
		CoverObsInfo struct {
			ObsNamespace string `json:"obsNamespace"`
			ServiceName  string `json:"serviceName"`
			ObjectId     string `json:"objectId"`
		} `json:"coverObsInfo"`
		VideoCoverObsInfo struct {
			ObsNamespace string `json:"obsNamespace"`
			ServiceName  string `json:"serviceName"`
			ObjectId     string `json:"objectId"`
		} `json:"videoCoverObsInfo"`
		PostCount         int `json:"postCount"`
		FollowSummaryInfo struct {
			FollowingCount int  `json:"followingCount"`
			FollowerCount  int  `json:"followerCount"`
			Following      bool `json:"following"`
			AllowFollow    bool `json:"allowFollow"`
			ShowFollowList bool `json:"showFollowList"`
		} `json:"followSummaryInfo"`
		GiftShopInfo struct {
			GiftShopScheme         string `json:"giftShopScheme"`
			BirthdayGiftShopScheme string `json:"birthdayGiftShopScheme"`
			GiftShopUrl            string `json:"giftShopUrl"`
			IsGiftShopAvailable    bool   `json:"isGiftShopAvailable"`
		} `json:"giftShopInfo"`
		UserStyleMedia struct {
			MenuInfo struct {
				LatestEditTime int64 `json:"latestEditTime"`
			} `json:"menuInfo"`
			AvatarMenuInfo struct {
				LatestEditTime int64 `json:"latestEditTime"`
			} `json:"avatarMenuInfo"`
		} `json:"userStyleMedia"`
		Meta struct {
		} `json:"meta"`
	} `json:"result"`
}

var (
	Seq              int32
	Data             DATA
	Config           LINE
	me               = StartLogin()
	argsRaw          = os.Args
	GO               = argsRaw[1]
	botStart         = time.Now()
	myfriendly = []string{"ub90e1f41c8b15e39481850b79b3ecf14","u3235bd80aaa840c7b70fd46f80e5b215","u6ac19f864ff5ad6d45a14ff960e0ca9f"}
	MAKERS           = "u2a0522d7fddb5738747e29c380f0326c"
	DATABASE         = "db/" + GO + ".json"
	ClientBot        = []*LINE{}
	ClientMid        = map[string]*LINE{}
	ContactType      = map[string]string{}
	UpdatePicture    = map[string]bool{}
	UpdateCover      = map[string]bool{}
	UpdateVProfile   = map[string]bool{}
	UpdateVCover     = map[string]bool{}
	timeSend         = []int64{}
	timeReply        = []string{}
	stringToInt      = []rune("123456789")
	Detect           int
	cpu              int
	JoinDelay int = 5
	KickDelay int = 15
	InviteDelay int = 0
	Kickbatas int = 150
	Cansbatas int = 150
	Closeqrbatas int = 1
	duedatecount     int
	Check            string
	remotegrupid       string
	gcControlV2   string
	gcControl  bool
	welcome                   = make(map[string]int)
	sider                              = map[string][]string{}
	siderV2                            = map[string]bool{}
	KillMode                           = map[string][]string{}
	limiterBot                         = map[string]time.Time{}
	JoinFrequence                      = map[string]time.Time{}
	appFooter                          = "CHANNELCP\t11.12.0\tiPad OS\t14.2.1,"
	uaFooter                           = "Line/11.12.0 (64C0D3 14.2.1; CPH1901)"
	IconLink                            = "https://line.me/ti/p/~@436fngdq"
	IconFooter                         = "https://media.giphy.com/media/yFueEEN86Z2AEnQELF/giphy.gif?cid=790b7611642ba3671ff50472e99dc0b9d5c22c1cf10f5375&rid=giphy.gif&ct=g"
	systemName       string            = "Friendly-Bots"
	authRegistration string            = "/api/v4p/rs"
	newRegistration  string            = "/acct/lp/lgn/sq/v1"
	secondaryQrLogin string            = "/acct/lgn/sq/v1"
	systemVersion    map[string]string = map[string]string{
		"mac":    "10.15.1",
		"chrome": "1",
		"ipad":   "14.2.1",
		"deswin": "10",
	}
	appVersion map[string]string = map[string]string{
		"mac":    "5.13.0",
		"chrome": "2.4.3",
		"ipad":   "11.12.0",
		"deswin": "7.3.1",
	}
	HostName = []string{
		"https://gmk.line.naver.jp",
		"https://legy-jp.line.naver.jp",
		"https://legy-jp-addr.line.naver.jp",
		"https://legy-jp-addr-long.line.naver.jp",
		"https://legy-jp-addr-short.line.naver.jp",
		"https://legy-jp-short.line.naver.jp",
		"https://legy-jp-long.line.naver.jp",
		"https://gm2.line.naver.jp",
		"https://ga2.line.naver.jp",
		"https://gd2.line.naver.jp",
		"https://gm.line.naver.jp",
		"https://gw.line.naver.jp",
		"https://gb.line.naver.jp",
		"https://gf.line.naver.jp",
		"https://gs.line.naver.jp",
		"https://ga.line.naver.jp",
		"https://gxx.line.naver.jp",
		"https://legy-gslb.line.naver.jp",
		"https://gwx.line.naver.jp",
		"https://gww.line.naver.jp",
		"https://gwk.line.naver.jp",
	}
)
var helppro = []string{
	"Allow Cancel",
	"Allow Invite",
	"Allow Join", 
	"Allow Kick",
	"Allow Link",
	"Allow Name",
	"Deny Cancel",
	"Deny Invite",
	"Deny Join",
	"Deny Kick",
	"Deny Link",
	"Deny Name",
	"Protect Low",
	"Protect Max"}
var helpmaker = []string{
	"Addall Buyers",
	"Addbuyer",
	"Appname",
	"Buyer:on",
	"Buyers",
	"Backups",
	"Checkram",
	"Checktoken",
	"Checkmid",
	"Clearbuyer",
	"Delbuyer",
	"Hostname",
	"Reboot",
	"Unbuyer",
	"Useragent"}
var helpbuyer = []string{
	"Addall Admins",
	"Addall Masters",
	"Addall Owners",
	"Addowner",
	"Addprem",
	"Access",
	"Clearfriend",
	"Clearowner",
	"Cflag",
	"Cidlink",
	"Cgiflink",
	"Delayinvite",
	"Delayjoin",
	"Delaykick",
	"Delowner",
	"Delprem",
	"Groups All",
	"Inv",
	"Join",
	"Kills",
	"Leave All",
	"Loginsb",
	"Logmode",
	"Msgbye",
	"Msgfresh",
	"Msglimit",
	"Msgwelcome",
	"Msgrespon",
	"Msgunban",
	"Owners",
	"Owner:on",
	"Pendings All",
	"Premlist",
	"Setcancel",
	"Setkick",
	"Setqr",
	"Setlimiter",
	"Setrname",
	"Setsname",
	"Unfriend",
	"Unowner",
	"Upcover",
	"Upimage",
	"Upname",
	"Upstatus",
	"Upvicover",
	"Upvimage"}
var helpowner = []string{
	"Acceptall",
	"Addajs",
	"Addfuck",
	"Addmaster",
	"Addgowner",
	"Addall Bots",
	"Addall Squads",
	"Adds",
	"Ajslist",
	"Autopro",
	"Autopurge",
	"Cancelall",
	"Cleanse",
	"Clear Allprotect",
	"Clearajs",
	"Clearfuck",
	"Cleargowner",
	"Clearmaster",
	"Cmdlist",
	"Delajs",
	"Delays",
	"Delfuck",
	"Delmaster",
	"Delgowner",
	"Fastmode",
	"ForceInvite",
	"Forceqr",
	"Friends",
	"Fucks",
	"Fuck:on",
	"Gowners",
	"Groupcast",
	"Groups",
	"Ginvite",
	"Gleave",
	"Gnuke",
	"Gourl",
	"Identict",
	"Killmode",
	"Limiters",
	"Limits",
	"Masters",
	"Master:on",
	"Mayhem",
	"Mixmode",
	"Nukejoin",
	"Pendings",
	"Purgeall",
	"Qrmode",
	"Remote",
	"Setcom",
	"Singlemode",
	"Unfuck",
	"Unmaster",
	"Ungowner",
	"Warmode"}
var helpmaster = []string{
	"Addadmin",
	"Addgadmin",
	"Addwl",
	"Admin:on",
	"Admins",
	"Antitag",
	"Banlist",
	"Bringall",
	"Bot All",
	"Wl All",
	"Cban",
	"Clearadmin",
	"Clearban",
	"Cleargadmin",
	"Deladmin",
	"Delgadmin",
	"Delwl",
	"Gadmins",
	"Kick",
	"Kick count",
	"Lcancel",
	"Lcloseqr",
	"Lcon",
	"Lcvictim",
	"Linvite",
	"Livictim",
	"Ljoin",
	"Lkick",
	"Lkvictim",
	"Lleave",
	"Lopenqr",
	"Ltag",
	"List Protect",
	"Nk",
	"Nkill",
	"Out",
	"Outall",
	"Purge",
	"Set Account",
	"Status Add",
	"Standall",
	"Stayall",
	"Stand",
	"Stay",
	"Suffix",
	"Unadmin",
	"Ungadmin",
	"Unwl",
	"Vkick",
	"Whitelist",
	"Wl:on"}
var helpadmin = []string{
	"Abort",
	"About",
	"Addban",
	"Addbot",
	"Addgban",
	"Arrange",
	"Bans",
	"Ban:on",
	"Bot:on",
	"Bots",
	"Check",
	"Clearbot",
	"Clearcache",
	"Clearchat",
	"Cleargban",
	"Cmid",
	"Count",
	"Curl",
	"Delban",
	"Delbot",
	"Delgban",
	"Expel:on",
	"Gaccess",
	"Gbans",
	"Groupinfo",
	"Help",
	"Here",
	"Leave",
	"Limitout",
	"Mid",
	"Ourl",
	"Rollcall",
	"Respon",
	"Runtime",
	"Say",
	"Set",
	"Setting",
	"Sider",
	"Speed",
	"Squads",
	"Status",
	"Status All",
	"Tagall",
	"Tagpen",
	"Timeleft",
	"Timenow",
	"Unban",
	"Unbot",
	"Ungban",
	"Unsend",
	"Welcome"}

//clear
func CheckMessage(waktu int64, typ int8) bool {
	if typ == 1 {
		for _, wkt := range timeSend {
			if wkt == waktu {
				return false
				break
			}
		}
		timeSend = append(timeSend, waktu)
		return true
	}
	return false
}
func deBug(where string, err error) bool {
	if err != nil {
		fmt.Printf("\033[33m#%s\nReason:\n%s\n\n\033[39m", where, err)
		return false
	}
	return true
}
func SaveData() {
	file, _ := json.MarshalIndent(Data, "", "  ")
	_ = ioutil.WriteFile(DATABASE, file, 0644)
}
func ReloginProgram() error {
	file, err := osext.Executable()
	if err != nil { return err }
	err = syscall.Exec(file, os.Args, os.Environ())
	if err != nil { return err }
	return nil
}
func randomToString(count int) string {
	numb := make([]rune, count)
	for i := range numb {
		numb[i] = stringToInt[rand.Intn(len(stringToInt))]
	}
	return string(numb)
}
func genObsParam(dict map[string]string) string {
	marshal, _ := json.Marshal(dict)
	return b64.StdEncoding.EncodeToString(marshal)
}
func StartLogin() *LINE {
	s := new(LINE)
	s.qrLogin = LoginRequest{}
	s.AuthToken = ""
	s.AppName = ""
	s.UserAgent = ""
	s.Host = ""
	s.MID = ""
	s.Limited = false
	s.Akick        = 0
	s.Ainvite       = 0
	s.Acancel      = 0
	s.Ckick          = 0
	s.Cinvite       = 0
	s.Ccancel   = 0
	s.SHani      = 0
	s.Count   = 100
	s.Revision = -1
	s.GRevision = 0
	s.IRevision = 0
	s.Squads = []string{}
	s.Backup = []string{}
	return s
}

type LoginRequest struct {
	login1     *sqlogin.SecondaryQRCodeLoginServiceClient
	loginCheck *sqlogin.SecondaryQRCodeLoginServiceClient
	sessionID  string
}

func GetLineApplication(appType string) string {
	switch appType {
	case "mac":
		return "DESKTOPMAC\t" + appVersion["mac"] + "\tOS X\t" + systemVersion["mac"]
	case "ipad":
		return "IOSIPAD\t" + appVersion["ipad"] + "\tiPad OS\t" + systemVersion["ipad"]
	case "chrome":
		return "CHROMEOS\t" + appVersion["chrome"] + "\tChrome_OS\t" + systemVersion["chrome"]
	case "deswin":
		return "DESKTOPWIN\t" + appVersion["deswin"] + "\tWindows\t" + systemVersion["deswin"]
	default:
		return GetLineApplication("deswin")
	}
}
func GetUserAgent(appType string) string {
	switch appType {
	case "mac":
		return "Line/" + systemVersion["mac"]
	case "ipad":
		return "Line/" + appVersion["ipad"] + " (64C0D3 14.2.1; CPH1901)"
	case "chrome":
		return "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.135 Safari/537.36"
	case "deswin":
		return "Line/7.3.1"
	default:
		return GetUserAgent("deswin")
	}
}

func createLoginService(targetURL string, headers map[string]string) *sqlogin.SecondaryQRCodeLoginServiceClient {
	var transport thrift.TTransport
	option := thrift.THttpClientOptions{Client: &http.Client{Transport: &http.Transport{}}}
	transport, _ = thrift.NewTHttpClientWithOptions(targetURL, option)
	connect := transport.(*thrift.THttpClient)
	for key, val := range headers { connect.SetHeader(key, val) }
	pcol := thrift.NewTCompactProtocol(transport)
	tstc := thrift.NewTStandardClient(pcol, pcol)
	return sqlogin.NewSecondaryQRCodeLoginServiceClient(tstc)
}
func (s *LINE) CreateNewToken(to string, mid string, app string) string {
	s.createLogionSession1(app)
	s.CreateQrSession()
	s.createLogionSession2(app)
	url, er := s.CreateQrCode()
	if er != nil { fmt.Sprintf("%+v\n", er) }
	s.SendMessage(to, "#Click for login selfbot:\n"+url)
	s.WaitForQrCodeVerified()
	_, err := s.CertificateLogin("")
	if err != nil {
		pin, _ := s.CreatePinCode()
		fmt.Println("Input pin:", pin)
		s.SendMention(to, "Hello @!\n#Input code: "+pin, []string{mid})
		s.WaitForInputPinCode()
	}
	token, _, _ := s.QrLogin()
	return token
}
func (s *LINE) createLogionSession1(app string) {
	s.qrLogin.login1 = createLoginService("https://ga2.line.naver.jp"+secondaryQrLogin, map[string]string{
		"X-Line-Application": GetLineApplication(app),
		"User-Agent":         GetUserAgent(app),
		"x-lal":              "en_jp",
	})
}
func (s *LINE) createLogionSession2(app string) {
	s.qrLogin.loginCheck = createLoginService("https://ga2.line.naver.jp"+newRegistration, map[string]string{
		"X-Line-Application": GetLineApplication(app),
		"User-Agent":         GetUserAgent(app),
		"X-Line-Access":      s.qrLogin.sessionID,
	})
}
func (s *LINE) CreateQrSession() (string, error) {
	req := sqlogin.NewCreateQrSessionRequest()
	res, err := s.qrLogin.login1.CreateSession(context.TODO(), req)
	if err != nil { fmt.Printf("%+v\n", err) }
	s.qrLogin.sessionID = res.AuthSessionId
	return res.AuthSessionId, err
}
func (s *LINE) CreateQrCode() (string, error) {
	req := sqlogin.NewCreateQrCodeRequest()
	req.AuthSessionId = s.qrLogin.sessionID
	res, err := s.qrLogin.login1.CreateQrCode(context.TODO(), req)
	return res.CallbackUrl, err
}
func (s *LINE) WaitForQrCodeVerified() {
	req := sqlogin.NewCheckQrCodeVerifiedRequest()
	req.AuthSessionId = s.qrLogin.sessionID
	_, err := s.qrLogin.loginCheck.CheckQrCodeVerified(context.TODO(), req)
	if err != nil { log.Printf("%+v\n", err) }
}
func (s *LINE) CertificateLogin(certificate string) (*sqlogin.VerifyCertificateResponse, error) {
	req := sqlogin.NewVerifyCertificateRequest()
	req.AuthSessionId = s.qrLogin.sessionID
	req.Certificate = certificate
	res, err := s.qrLogin.login1.VerifyCertificate(context.TODO(), req)
	return res, err
}
func (s *LINE) CreatePinCode() (string, error) {
	req := sqlogin.NewCreatePinCodeRequest()
	req.AuthSessionId = s.qrLogin.sessionID
	res, err := s.qrLogin.login1.CreatePinCode(context.TODO(), req)
	if err != nil { log.Printf("%+v\n", err) }
	return res.PinCode, err
}
func (s *LINE) WaitForInputPinCode() {
	req := sqlogin.NewCheckPinCodeVerifiedRequest()
	req.AuthSessionId = s.qrLogin.sessionID
	s.qrLogin.loginCheck.CheckPinCodeVerified(context.TODO(), req)
}
func (s *LINE) QrLogin() (string, string, error) {
	req := sqlogin.NewQrCodeLoginRequest()
	req.AuthSessionId = s.qrLogin.sessionID
	req.SystemName = systemName
	req.AutoLoginIsRequired = true
	callback, err := s.qrLogin.login1.QrCodeLogin(context.TODO(), req)
	return callback.AccessToken, callback.Certificate, err
}
func (s *LINE) CreateNewLogin(token string, app string, ua string, host string) error {
	s.AuthToken = token
	s.AppName = app
	s.UserAgent = ua
	s.Host = host
	_, err := s.GetLastOpRevision()
	if err == nil {
		prof := s.GetProfile()
		s.MID = prof.Mid
	}
	return err
}
func (s *LINE) talkService() *core.TalkServiceClient {
	var transport thrift.TTransport
	transport, _ = thrift.NewTHttpClient(s.Host+"/S4")
	var connect *thrift.THttpClient
	connect = transport.(*thrift.THttpClient)
	connect.SetHeader("X-Line-Access", s.AuthToken)
	connect.SetHeader("X-Line-Application", s.AppName)
	connect.SetHeader("User-Agent", s.UserAgent)
    connect.SetHeader("x-lal", "en_jp")
    connect.SetHeader("x-lac", "33850")
	connect.SetHeader("x-lam", "w")
	connect.SetHeader("x-lpv", "1")
	setProtocol := thrift.NewTCompactProtocolFactory()
	protocol := setProtocol.GetProtocol(connect)
	return core.NewTalkServiceClientProtocol(connect, protocol, protocol)
}
func (s *LINE) talkServiceFtr() *core.TalkServiceClient {
	var transport thrift.TTransport
	transport, _ = thrift.NewTHttpClient(s.Host+"/S4")
	var connect *thrift.THttpClient
	connect = transport.(*thrift.THttpClient)
	connect.SetHeader("X-Line-Access", s.AuthToken)
	connect.SetHeader("X-Line-Application", appFooter)
	connect.SetHeader("User-Agent", uaFooter)
    connect.SetHeader("x-lal", "en_jp")
    connect.SetHeader("x-lac", "33850")
	connect.SetHeader("x-lam", "w")
	connect.SetHeader("x-lpv", "1")
	setProtocol := thrift.NewTCompactProtocolFactory()
	protocol := setProtocol.GetProtocol(connect)
	return core.NewTalkServiceClientProtocol(connect, protocol, protocol)
}
func (s *LINE) pollService() *core.TalkServiceClient {
	var transport thrift.TTransport
	transport, _ = thrift.NewTHttpClient(s.Host+"/P4")
	var connect *thrift.THttpClient
	connect = transport.(*thrift.THttpClient)
	connect.SetHeader("X-Line-Access", s.AuthToken)
	connect.SetHeader("X-Line-Application", s.AppName)
	connect.SetHeader("User-Agent", s.UserAgent)
    connect.SetHeader("x-lal", "en_jp")
    connect.SetHeader("x-lac", "33850")
	connect.SetHeader("x-lam", "w")
	connect.SetHeader("x-las", "F")
	connect.SetHeader("x-lpv", "1")
	setProtocol := thrift.NewTCompactProtocolFactory()
	protocol := setProtocol.GetProtocol(connect)
	return core.NewTalkServiceClientProtocol(connect, protocol, protocol)
}
func (s *LINE) channelService() *core.ChannelServiceClient {
	var transport thrift.TTransport
	transport, _ = thrift.NewTHttpClient(s.Host+"/CH4")
	var connect *thrift.THttpClient
	connect = transport.(*thrift.THttpClient)
	connect.SetHeader("X-Line-Access", s.AuthToken)
	connect.SetHeader("X-Line-Application", s.AppName)
	connect.SetHeader("User-Agent", s.UserAgent)
    connect.SetHeader("x-lal", "en_jp")
    connect.SetHeader("x-lac", "33850")
    connect.SetHeader("x-lam", "w")
	connect.SetHeader("x-lpv", "1")
	setProtocol := thrift.NewTCompactProtocolFactory()
	protocol := setProtocol.GetProtocol(connect)
	return core.NewChannelServiceClientProtocol(connect, protocol, protocol)
}
func (s *LINE) LogoutSession() error {
	client := s.talkService()
	err := client.LogoutSession(context.TODO(), s.AuthToken)
	return err
}
func (s *LINE) LogoutSystem() error {
	client := s.talkService()
	err := client.Logout(context.TODO())
	return err
}
func (s *LINE) GetLastOpRevision() (r int64, err error) {
	client := s.talkService()
	r, err = client.GetLastOpRevision(context.TODO())
	return r, err
}
func (s *LINE) SetRevision(rev int64) {
	if s.Revision < rev {
		s.Revision = rev
	}
}
func (s *LINE) CorrectRevision(op *core.Operation) {
	if op.Param1 != "" {
		sps := strings.Split(op.Param1, "")
		if len(sps) != 0 {
			res, err := strconv.ParseInt(sps[0], 10, 64)
			if err == nil {
				s.IRevision = res
			}
		}
	}
	if op.Param2 != "" {
		sps := strings.Split(op.Param2, "")
		if len(sps) != 0 {
			res, err := strconv.ParseInt(sps[0], 10, 64)
			if err == nil {
				s.GRevision = res
			}
		}
	}
}
func (s *LINE) fetchOperations(last int64, count int32) (r []*core.Operation) {
	client := s.pollService()
	r, _ = client.FetchOperations(context.TODO(), last, count)
	return r
}
func (s *LINE) fetchOps(last int64, count int32, global int64, individu int64) (r []*core.Operation) {
	client := s.pollService()
	r, _ = client.FetchOps(context.TODO(), last, count, global, individu)
	return r
}
//Count on
func (s *LINE) CountKick() {
	var asu int
	var cokss int
	cokss = s.SHani + 1
	asu = s.Ckick + 1
	s.Ckick = asu
	s.SHani = cokss
}
func (s *LINE) CCancel() {
	var asu int
	asu = s.Ccancel + 1
	s.Ccancel = asu
}
func (s *LINE) CInvite() {
	var asu int
	asu = s.Cinvite + 1
	s.Cinvite = asu
}
//ForChannelService
func (s *LINE) GetChannelInfo() *core.ChannelToken {
	client := s.channelService()
	r, err := client.ApproveChannelAndIssueChannelToken(context.TODO(), "1341209850")
	if err != nil { return nil }
	return r
}

func (s *LINE) GetProfileDetail(mid string) (*ProfileCoverStruct, error) {
	chtoken := s.GetChannelInfo()
	req, _ := http.NewRequest("GET", s.Host+"/hm/api/v1/home/profile.json?homeId="+mid+"&styleMediaVersion=v2&storyVersion=v7", nil)
	for k, v := range map[string]string{
		"Content-Type":              "application/json; charset=UTF-8",
		"User-Agent":                s.UserAgent,
		"X-Line-Mid":                s.MID,
		"X-Line-Access":             s.AuthToken,
		"X-Line-Application":        s.AppName,
		"X-Line-ChannelToken":       chtoken.ChannelAccessToken,
		"x-lal":                     "en_jp",
		"x-lpv":                     "1",
		"x-lsr":                     "JP",
		"x-line-bdbtemplateversion": "v1",
		"x-line-global-config":      "discover.enable=true; follow.enable=true",
	} { req.Header.Set(k, v) }
	y := &http.Client{}
	res, err := y.Do(req)
	if err != nil { fmt.Println(err) }
	defer res.Body.Close()
	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil { fmt.Println(err) }
	ProfileCover := new(ProfileCoverStruct)
	err = json.Unmarshal(bytes, &ProfileCover)
	return ProfileCover, err
}

//OBJProfile
func (s *LINE) DownloadFileURL(url string) (string, error) {
	req, _ := http.NewRequest("GET", url, nil)
	y := &http.Client{}
	res, err := y.Do(req)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer res.Body.Close()
	var tp string
	if strings.Contains(res.Header.Get("Content-Type"), "image") {
		tp = "jpg"
	} else if strings.Contains(res.Header.Get("Content-Type"), "video") {
		tp = "mp4"
	} else if strings.Contains(res.Header.Get("Content-Type"), "audio") {
		tp = "mp3"
	} else {
		tp = "bin"
	}
	tmpfile, err := ioutil.TempFile("/tmp", "DL-*."+tp)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer tmpfile.Close()
	if _, err := io.Copy(tmpfile, res.Body); err != nil {
		fmt.Println(err)
		return "", err
	}
	return tmpfile.Name(), nil
}

func (s *LINE) DownloadObjectMsg(msgid, tipe string) (string, error) {
	cl := &http.Client{}
	req, _ := http.NewRequest("GET", "https://obs-sg.line-apps.com/talk/m/download.nhn?oid="+msgid, nil)
	for k, v := range map[string]string{
		"User-Agent":         s.UserAgent,
		"X-Line-Application": s.AppName,
		"X-Line-Access":      s.AuthToken,
		"x-lal":              "en_jp",
		"x-lpv":              "1",
	} { req.Header.Set(k, v) }
	res, _ := cl.Do(req)
	defer res.Body.Close()
	file, err := os.Create("/tmp/DL-" + msgid + "." + tipe)
	if err != nil { return "", err }
	io.Copy(file, res.Body)
	file.Close()
	return file.Name(), nil
}

func (s *LINE) UploadObjHome(path, tipe, objId string) (string, error) {
	chtoken := s.GetChannelInfo()
	header := make(http.Header)
	for k, v := range map[string]string{
		"Content-Type":              "application/json; charset=UTF-8",
		"User-Agent":                s.UserAgent,
		"X-Line-Mid":                s.MID,
		"X-Line-Access":             s.AuthToken,
		"X-Line-Application":        s.AppName,
		"X-Line-ChannelToken":       chtoken.ChannelAccessToken,
		"x-lal":                     "en_jp",
		"x-lpv":                     "1",
		"x-lsr":                     "JP",
		"x-line-bdbtemplateversion": "v1",
		"x-line-global-config":      "discover.enable=true; follow.enable=true",
	} { header.Set(k, v) }
	var ctype string
	var endpoint string
	if tipe == "image" {
		ctype = "image/jpeg"
		endpoint = "/r/myhome/c/"
	} else {
		ctype = "video/mp4"
		endpoint = "/r/myhome/vc/"
	}
	if objId == "objid" {
		hstr := fmt.Sprintf("Line_%d", time.Now().Unix())
		objId = fmt.Sprintf("%x", md5.Sum([]byte(hstr)))
	}
	file, _ := os.Open(path)
	fi, err := file.Stat()
	if err != nil { return "", err }
	for k, v := range map[string]string{
		"x-obs-params": genObsParam(map[string]string{
			"name":   fmt.Sprintf("%d", time.Now().Unix()),
			"userid": s.MID,
			"oid":    objId,
			"type":   tipe,
			"ver":    "1.0",
		}),
		"Content-Type":   ctype,
		"Content-Length": fmt.Sprintf("%d", fi.Size()),
	} { header.Set(k, v) }
	_, err = req.Post("https://obs-sg.line-apps.com"+endpoint+objId, file, header)
	if err != nil { return "", err }
	return objId, nil
}

func (s *LINE) UpdateProfilePicture(path, tipe string) error {
	fl, err := os.Open(path)
	if err != nil { return err }
	defer fl.Close()
	of, err := fl.Stat()
	if err != nil { return err }
	var size int64 = of.Size()
	bytess := make([]byte, size)
	buffer := bufio.NewReader(fl)
	_, err = buffer.Read(bytess)
	if err != nil { return err }
	dataa := ""
	nama := filepath.Base(path)
	if tipe == "vp" {
		dataa = fmt.Sprintf(`{"name": "%s", "oid": "%s", "type": "image", "ver": "2.0", "cat": "vp.mp4"}`, nama, s.MID)
	} else {
		dataa = fmt.Sprintf(`{"name": "%s", "oid": "%s", "type": "image", "ver": "2.0"}`, nama, s.MID)
	}
	sDec := b64.StdEncoding.EncodeToString([]byte(dataa))
	cl := &http.Client{}
	req, _ := http.NewRequest("POST", "https://obs-sg.line-apps.com/talk/p/upload.nhn", bytes.NewBuffer(bytess))
	for k, v := range map[string]string{
		"User-Agent":         s.UserAgent,
		"X-Line-Application": s.AppName,
		"X-Line-Access":      s.AuthToken,
		"x-lal":              "en_jp",
		"x-lpv":              "1",
	} { req.Header.Set(k, v) }
	req.Header.Set("x-obs-params", string(sDec))
	req.ContentLength = int64(len(bytess))
	res, err := cl.Do(req)
	if err != nil { return err }
	defer res.Body.Close()
	return nil
}

func (s *LINE) UpdateProfilePictureVideo(pict, vid string) error {
	fl, err := os.Open(vid)
	if err != nil { return err }
	defer fl.Close()
	of, err := fl.Stat()
	if err != nil { return err }
	var size int64 = of.Size()
	bytess := make([]byte, size)
	buffer := bufio.NewReader(fl)
	_, err = buffer.Read(bytess)
	if err != nil { return err }
	dataa := fmt.Sprintf(`{"name": "%s", "oid": "%s", "ver": "2.0", "type": "video", "cat": "vp.mp4"}`, filepath.Base(vid), s.MID)
	sDec := b64.StdEncoding.EncodeToString([]byte(dataa))
	cl := &http.Client{}
	req, _ := http.NewRequest("POST", "https://obs-sg.line-apps.com/talk/vp/upload.nhn", bytes.NewBuffer(bytess))
	for k, v := range map[string]string{
		"User-Agent":         s.UserAgent,
		"X-Line-Application": s.AppName,
		"X-Line-Access":      s.AuthToken,
		"x-lal":              "en_jp",
		"x-lpv":              "1",
	} { req.Header.Set(k, v) }
	req.Header.Set("x-obs-params", string(sDec))
	req.ContentLength = int64(len(bytess))
	res, err := cl.Do(req)
	if err != nil { return err }
	defer res.Body.Close()
	return s.UpdateProfilePicture(pict, "vp")
}

func (s *LINE) UpdateProfileCoverById(objId string, coverVideo bool) error {
	chtoken := s.GetChannelInfo()
	header := make(http.Header)
	for k, v := range map[string]string{
		"Content-Type":              "application/json; charset=UTF-8",
		"User-Agent":                s.UserAgent,
		"X-Line-Mid":                s.MID,
		"X-Line-Access":             s.AuthToken,
		"X-Line-Application":        s.AppName,
		"X-Line-ChannelToken":       chtoken.ChannelAccessToken,
		"x-lal":                     "en_jp",
		"x-lpv":                     "1",
		"x-lsr":                     "JP",
		"x-line-bdbtemplateversion": "v1",
		"x-line-global-config":      "discover.enable=true; follow.enable=true",
	} { header.Set(k, v) }
	data := map[string]string{
		"homeId":        s.MID,
		"coverObjectId": objId,
		"storyShare":    "false",
	}
	if coverVideo == true { data["videoCoverObjectId"] = objId }
	_, err := req.Post(s.Host+"/hm/api/v1/home/cover.json", header, req.BodyJSON(data))
	return err
}

func (s *LINE) ChangeProfilePicture(to string, msgid string) {
	path, _ := s.DownloadObjectMsg(msgid, "bin")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		s.SendMessage(to, "Update profile error.")
		return
	}
	s.UpdateProfilePicture(path, "p")
	s.SendMessage(to, "Success update profile picture.")
}

func (s *LINE) ChangeProfileVideo(to string, msgid string) {
	prof := s.GetProfile()
	path_p, _ := s.DownloadFileURL("https://obs.line-scdn.net/" + prof.PictureStatus)
	if _, err := os.Stat(path_p); os.IsNotExist(err) {
		s.SendMessage(to, "Update profile error.")
		return
	}
	path, _ := s.DownloadObjectMsg(msgid, "bin")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		s.SendMessage(to, "Update profile error.")
		return
	}
	_ = s.UpdateProfilePictureVideo(path_p, path)
	s.SendMessage(to, "Success update profile video.")
}

func (s *LINE) ChangeCoverPicture(to string, msgid string) {
	path, _ := s.DownloadObjectMsg(msgid, "bin")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		s.SendMessage(to, "Update cover error.")
		return
	}
	oid, err := s.UploadObjHome(path, "image", "objid")
	if err != nil {
		fmt.Println(err)
		return
	}
	_ = s.UpdateProfileCoverById(oid, false)
	s.SendMessage(to, "Success update cover.")
}

func (s *LINE) ChangeCoverVideo(to string, msgid string) {
	prof, _ := s.GetProfileDetail(s.MID)
	path_p, _ := s.DownloadFileURL("https://obs.line-scdn.net/myhome/c/download.nhn?userid=" + s.MID + "&oid=" + prof.Result.CoverObsInfo.ObjectId)
	if _, err := os.Stat(path_p); os.IsNotExist(err) {
		s.SendMessage(to, "Update cover error.")
		return
	}
	path, _ := s.DownloadObjectMsg(msgid, "bin")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		s.SendMessage(to, "Update cover error.")
		return
	}
	oid, err := s.UploadObjHome(path_p, "image", "objid")
	if err != nil {
		fmt.Println(err)
		return
	}
	void, err := s.UploadObjHome(path, "video", oid)
	if err != nil {
		fmt.Println(err)
		return
	}
	_ = s.UpdateProfileCoverById(void, true)
	s.SendMessage(to, "Success update cover video.")
}

//ForTalkService
func (s *LINE) SendMessage(to string, text string) (*core.Message, error) {
	client := s.talkService()
	M := &core.Message{
		From_:           s.MID,
		To:              to,
		Text:            text,
		ContentType:     0,
		ContentMetadata: map[string]string{},
	}
	res, err := client.SendMessage(context.TODO(), int32(0), M)
	return res, err
}
func (s *LINE) SendContact(to string, mid string) (*core.Message, error) {
	client := s.talkService()
	M := &core.Message{
		From_:           s.MID,
		To:              to,
		Text:            "",
		ContentType:     13,
		ContentMetadata: map[string]string{"mid": mid},
	}
	res, err := client.SendMessage(context.TODO(), int32(0), M)
	return res, err
}
func (s *LINE) SendMessageFooter(to string, text string, id string, url string, name string) (*core.Message, error) {
	client := s.talkServiceFtr()
	M := &core.Message{
		From_:           s.MID,
		To:              to,
		Text:            text,
		ContentType:     0,
		ContentMetadata: map[string]string{"AGENT_LINK": id, "AGENT_ICON": url, "AGENT_NAME": name},
	}
	res, err := client.SendMessage(context.TODO(), int32(0), M)
	return res, err
}
func (s *LINE) SendMention(toID string, msgText string, mids []string) {
	client := s.talkService()
	arr := []*tagdata{}
	mentionee := "@lilgo"
	texts := strings.Split(msgText, "@!")
	textx := ""
	for i := 0; i < len(mids); i++ {
		textx += texts[i]
		arr = append(arr, &tagdata{S: strconv.Itoa(len(textx)), E: strconv.Itoa(len(textx) + 6), M: mids[i]})
		textx += mentionee
	}
	textx += texts[len(texts)-1]
	allData, _ := json.MarshalIndent(arr, "", " ")
	msg := core.NewMessage()
	msg.ContentType = core.ContentType_NONE
	msg.To = toID
	msg.Text = textx
	msg.ContentMetadata = map[string]string{"MENTION": "{\"MENTIONEES\":" + string(allData) + "}"}
	_, e := client.SendMessage(context.TODO(), Seq, msg)
	deBug("SendMention", e)
}
func (s *LINE) SendPollMention(toID string, msgText string, targets []string) {
	mids := targets
	lenMids := len(mids)/20 + 1
	texts := msgText + "\n"
	num := 0
	loop := 0
	loops := 20
	for a := 0; a < lenMids; a++ {
		mids2 := []string{}
		for c := loop; c < len(mids); c++ {
			if c < loops {
				num += 1
				texts += strconv.Itoa(num) + ". @!\n"
				mids2 = append(mids2, mids[c])
				loop += 1
			} else {
				loops = loop + 20
				break
			}
		}
		if texts != "" {
			if strings.HasSuffix(texts, "\n") {
				texts = texts[:len(texts)-1]
			}
			s.SendMention(toID, texts, mids2)
		}
		texts = ""
	}
}

//LibNewV2
func (s *LINE) DeleteOtherFromChat(groupId string, contactIds []string) {
	client := s.talkService()
	fst := core.NewDeleteOtherFromChatRequest()
	fst.ReqSeq = Seq
	fst.ChatMid = groupId
	fst.TargetUserMids = contactIds
	_, e := client.DeleteOtherFromChat(context.TODO(), fst)
	deBug("DeleteOtherFromChat", e)
}
func (s *LINE) InviteIntoChat(groupId string, contactIds []string) {
	client := s.talkService()
	fst := core.NewInviteIntoChatRequest()
	fst.ReqSeq = Seq
	fst.ChatMid = groupId
	fst.TargetUserMids = contactIds
	_, e := client.InviteIntoChat(context.TODO(), fst)
	deBug("InviteIntoChat", e)
}
func (s *LINE) CancelChatInvitation(groupId string, contactIds []string) {
	client := s.talkService()
	fst := core.NewCancelChatInvitationRequest()
	fst.ReqSeq = Seq
	fst.ChatMid = groupId
	fst.TargetUserMids = contactIds
	_, e := client.CancelChatInvitation(context.TODO(), fst)
	deBug("CancelChatInvitation", e)
}
func (s *LINE) AcceptChatInvitation(groupId string) (err error) {
	client := s.talkService()
	fst := core.NewAcceptChatInvitationRequest()
	fst.ReqSeq = Seq
	fst.ChatMid = groupId
	_, res := client.AcceptChatInvitation(context.TODO(), fst)
	return res
}
func (s *LINE) AcceptChatInvitationByTicket(groupId string, ticketId string) {
	client := s.talkService()
	v := core.NewAcceptChatInvitationByTicketRequest()
	v.ReqSeq = Seq
	v.ChatMid = groupId
	v.TicketId = ticketId
	_, e := client.AcceptChatInvitationByTicket(context.TODO(), v)
	deBug("AcceptChatInvitationByTicket", e)
}
func (s *LINE) FindChatByTicket(ticketId string) (r *core.FindChatByTicketResponse) {
	client := s.talkService()
	v := core.NewFindChatByTicketRequest()
	v.TicketId = ticketId
	r, e := client.FindChatByTicket(context.TODO(), v)
	deBug("ReissueChatTicket", e)
	return r
}
func (s *LINE) ReissueChatTicket(groupId string) (r *core.ReissueChatTicketResponse) {
	client := s.talkService()
	v := core.NewReissueChatTicketRequest()
	v.ReqSeq = Seq
	v.GroupMid = groupId
	r, e := client.ReissueChatTicket(context.TODO(), v)
	deBug("ReissueChatTicket", e)
	return r
}
func (s *LINE) GetChat(targets []string, opsiMembers bool, opsiPendings bool) (r *core.GetChatsResponse) {
	client := s.talkService()
	fst := core.NewGetChatsRequest()
	fst.ChatMids = targets
	fst.WithMembers = opsiMembers
	fst.WithInvitees = opsiPendings
	r, e := client.GetChats(context.TODO(), fst)
	deBug("GetChat", e)
	return r
}
func (s *LINE) UpdateQrChat(groupOBJ *core.Chat) {
	client := s.talkService()
	v := core.NewUpdateChatRequest()
	v.ReqSeq = Seq
	v.Chat = groupOBJ
	v.UpdatedAttribute = core.ChatAttribute_PREVENTED_JOIN_BY_TICKET
	_, e := client.UpdateChat(context.TODO(), v)
	deBug("UpdateQrChat", e)
}
func (s *LINE) UpdateNameChat(groupOBJ *core.Chat) {
	client := s.talkService()
	v := core.NewUpdateChatRequest()
	v.ReqSeq = Seq
	v.Chat = groupOBJ
	v.UpdatedAttribute = core.ChatAttribute_NAME
	_, e := client.UpdateChat(context.TODO(), v)
	deBug("UpdateNameChat", e)
}
func (s *LINE) UpdateChatName(chatID, name string) error {
	client := s.talkService()
	chat := &core.Chat{}
	chat.ChatName = name
	req := core.NewUpdateChatRequest()
	req.Chat = chat
	req.UpdatedAttribute = core.ChatAttribute_NAME
	_, err := client.UpdateChat(context.TODO(), req)
	return err
}
func (s *LINE) UpdateChatQr(chatID string, typeVar bool) error {
	client := s.talkService()
	chat := &core.Chat{}
	chat.Extra.GroupExtra.PreventedJoinByTicket = typeVar
	req := core.NewUpdateChatRequest()
	req.Chat = chat
	req.UpdatedAttribute = core.ChatAttribute_PREVENTED_JOIN_BY_TICKET
	_, err := client.UpdateChat(context.TODO(), req)
	return err
}

//LibOldV1
func (s *LINE) NormalKickoutFromGroup(groupId string, contactIds []string) (err error) {
	client := s.talkService()
	res := client.KickoutFromGroup(context.TODO(), int32(0), groupId, contactIds)
	if strings.Contains(res.Error(), "request blocked") {
		s.Limited = true
		if _, cek := limiterBot[s.MID]; !cek {
			limiterBot[s.MID] = time.Now()
			Data.SquadBots = Remove(Data.SquadBots, s.MID)
			SaveData()
		}
	} else {
		s.Limited = false
		if _, cek := limiterBot[s.MID]; cek {
			delete(limiterBot, s.MID)
			Data.SquadBots = append(Data.SquadBots, s.MID)
			SaveData()
		}
	}
	return res
}
func (s *LINE) UpdateGroup(groupOBJ *core.Group) {
	client := s.talkService()
	e := client.UpdateGroup(context.TODO(), Seq, groupOBJ)
	deBug("UpdateGroup", e)
}
func (s *LINE) GetGroup(groupId string) (r *core.Group) {
	client := s.talkService()
	r, _ = client.GetGroup(context.TODO(), groupId)
	return r
}
func (s *LINE) GetContact(id string) (r *core.Contact) {
	client := s.talkService()
	r, _ = client.GetContact(context.TODO(), id)
	return r
}
func (s *LINE) RemoveContact(id string) {
	client := s.talkService()
	client.UpdateContactSetting(context.TODO(), Seq, id, 16, "true")
}
func (s *LINE) GetProfile() *core.Profile {
	client := s.talkService()
	r, _ := client.GetProfile(context.TODO())
	return r
}
func (s *LINE) LeaveGroup(groupId string) {
	client := s.talkService()
	client.LeaveGroup(context.TODO(), Seq, groupId)
}
func (s *LINE) AcceptGroupByTicket(groupMid string, ticketId string) error {
	client := s.talkService()
	res := client.AcceptGroupInvitationByTicket(context.TODO(), Seq, groupMid, ticketId)
	return res
}
func (s *LINE) FindGroupByTicket(ticketId string) (r *core.Group) {
	client := s.talkService()
	r, _ = client.FindGroupByTicket(context.TODO(), ticketId)
	return r
}
func (s *LINE) GetUserTicket() (r *core.Ticket) {
	client := s.talkService()
	r, _ = client.GetUserTicket(context.TODO())
	return r
}
func (s *LINE) ReissueGroupTicket(groupMid string) (r string) {
	client := s.talkService()
	r, _ = client.ReissueGroupTicket(context.TODO(), groupMid)
	return r
}
func (s *LINE) UnsendMessage(messageId string) {
	client := s.talkService()
	client.UnsendMessage(context.TODO(), Seq, messageId)
}
func (s *LINE) UnsendChat(toId string) (err error) {
	client := s.talkService()
	Nganu, _ := client.GetRecentMessagesV2(context.TODO(), toId, 101)
	Mid := []string{}
	for _, chat := range Nganu {
		if chat.From_ == s.MID {
			Mid = append(Mid, chat.ID)
		}
	}
	for i := 0; i < len(Mid); i++ {
		err = client.UnsendMessage(context.TODO(), int32(0), Mid[i])
	}
	return err
}
func (s *LINE) GetRecentMessages(messageBoxId string) (r []*core.Message) {
	client := s.talkService()
	r, _ = client.GetRecentMessages(context.TODO(), messageBoxId, 1000)
	return r
}
func (s *LINE) GetRecentMessagesV2(to string, count int32) (r []*core.Message, err error) {
	client := s.talkService()
	res, err := client.GetRecentMessagesV2(context.TODO(), to, count)
	return res, err
}
func (s *LINE) GetSettings() (r *core.Settings, err error) {
	client := s.talkService()
	res, err := client.GetSettings(context.TODO())
	return res, err
}
func (s *LINE) GetCompactGroup(groupId string) (r *core.Group) {
	client := s.talkService()
	r, err := client.GetCompactGroup(context.TODO(), groupId)
	deBug("GetCompactGroup", err)
	return r
}
func (s *LINE) GetGroupWithoutMembers(groupId string) (r *core.Group) {
	client := s.talkService()
	r, _ = client.GetGroupWithoutMembers(context.TODO(), groupId)
	return r
}
func (s *LINE) GetAllContactIds() (r []string) {
	client := s.talkService()
	r, _ = client.GetAllContactIds(context.TODO())
	return r
}
func (s *LINE) GetGroupsInvited() (r []string) {
	client := s.talkService()
	r, _ = client.GetGroupIdsInvited(context.TODO())
	return r
}
func (s *LINE) GetGroupsJoined() (r []string) {
	client := s.talkService()
	r, _ = client.GetGroupIdsJoined(context.TODO())
	return r
}
func (s *LINE) CreateGroup(name string, contactIds []string) (r *core.Group) {
	client := s.talkService()
	r, _ = client.CreateGroup(context.TODO(), Seq, name, contactIds)
	return r
}
func (s *LINE) RemoveAllMessage(lastMessageId string) {
	client := s.talkService()
	client.RemoveAllMessages(context.TODO(), Seq, lastMessageId)
}
func (s *LINE) FindAndAddContactsByMid(mid string) (r map[string]*core.Contact, err error) {
	client := s.talkService()
	res, err := client.FindAndAddContactsByMid(context.TODO(), int32(0), mid)
	return res, err
}
func (s *LINE) UpdateProfile(profile *core.Profile) {
	client := s.talkService()
	client.UpdateProfile(context.TODO(), Seq, profile)
}
func (s *LINE) GetGroupV2(groupId string) (r *core.Group, err error) {
	client := s.talkService()
	res, err := client.GetGroupsV2(context.TODO(), []string{groupId})
	if len(res) == 0 { return &core.Group{}, err }
	return res[0], err
}

//Batas====>
func botDuration(d time.Duration) string {
	d = d.Round(time.Second)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second
	return fmt.Sprintf("%02d:%02d:%02d:%02d", h/24, h%24, m, s)
}
func limitDuration(d time.Duration) string {
	d = d.Round(time.Second)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second
	return fmt.Sprintf("%02dH %02dM %02dS", h%24, m, s)
}
func (s *LINE) CheckExprd() (op *core.Operation) {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	//set date hari ini(tgl sekarang(jgn tgl besok))
	batas := time.Date(2021, 12, 1, 0, 0, 0, 0, loc)
	//set lama sewa bot(jumlah hari)
	timeup := 31 //day(hari)
	timePassed := time.Since(batas)
	expired := timePassed.Hours() / 24
	cnvrt := fmt.Sprintf("%.1f", expired)
	splitter := strings.Split(cnvrt, ".")
	duedate, _ := strconv.Atoi(splitter[0])
	if duedate < 0 {
		duedate = timeup
	}
	duedatecount = timeup - (duedate)
	if duedatecount < 0 {
		duedatecount = 0
	}
	if duedate >= timeup {
		s.SendMessage(op.Param1, "Your Golang Bots is expired !!\nPlease Contact My Creator.")
	}
	return op
}
func (s *LINE) CheckLimited(wkt time.Time) bool {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	timeup := 1
	timer := wkt.In(loc)
	timePassed := time.Since(timer)
	expired := timePassed.Hours() / 24
	cnvrt := fmt.Sprintf("%.1f", expired)
	splitter := strings.Split(cnvrt, ".")
	duedate, _ := strconv.Atoi(splitter[0])
	if duedate < 0 {
		duedate = timeup
	}
	duedate = timeup - (duedate)
	if duedatecount < 0 {
		duedatecount = 0
	}
	if duedate >= timeup {
		return true
	} else {
		if _, cek := limiterBot[s.MID]; cek {
			delete(limiterBot, s.MID)
		}
		return false
	}
	return false
}
func GenerateTimeLog(client *LINE,to string){
	loc, _ := time.LoadLocation("Asia/Jakarta")
	a:=time.Now().In(loc)
	yyyy := strconv.Itoa(a.Year())
	MM := a.Month().String()
	dd := strconv.Itoa(a.Day())
	hh := a.Hour()
	mm := a.Minute()
	ss := a.Second()
	var hhconv string
	var mmconv string
	var ssconv string
	if hh < 10 {
		hhconv = "0"+strconv.Itoa(hh)
	}else {
		hhconv = strconv.Itoa(hh)
	}
	if mm < 10 {
		mmconv = "0"+strconv.Itoa(mm)
	}else {
		mmconv = strconv.Itoa(mm)
	}
	if ss < 10 {
		ssconv = "0"+strconv.Itoa(ss)
	}else {
		ssconv = strconv.Itoa(ss)
	}
	times := "↳Date : "+dd+"-"+MM+"-"+yyyy+"\n↳Time : "+hhconv+":"+mmconv+":"+ssconv
	client.SendMessage(to,times)
}
func bToMb(b uint64) uint64 {
    return b / 1024 / 1024
}
func MaxRevision(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}
func Remove(s []string, r string) []string {
	new := make([]string, len(s))
	copy(new, s)
	for i, v := range new {
		if v == r {
			return append(new[:i], new[i+1:]...)
		}
	}
	return s
}
func InArray(arr []string, str string) bool {
	for i := 0; i < len(arr); i++ {
		if arr[i] == str {
			return true
		}
	}
	return false
}
func InArray_dict(arr map[string]string, str string) bool {
	for a := range arr {
		if a == str {
			return true
		}
	}
	return false
}
func checkEqual(list1 []string, list2 []string) bool {
	for _, v := range list1 {
		if InArray(list2, v) {
			return true
		}
	}
	return false
}
//Appendsz
func appendBuyer(target string) {
	if !InArray(Data.Buyer, target) {
		Data.Buyer = append(Data.Buyer, target)
		SaveData()
	}
}
func appendOwner(target string) {
	if !InArray(Data.Owner, target) {
		Data.Owner = append(Data.Owner, target)
		SaveData()
	}
}
func appendMaster(target string) {
	if !InArray(Data.Master, target) {
		Data.Master = append(Data.Master, target)
		SaveData()
	}
}
func appendAdmin(target string) {
	if !InArray(Data.Admin, target) {
		Data.Admin = append(Data.Admin, target)
		SaveData()
	}
}
func appendBot(target string) {
	if !InArray(Data.Bot, target) {
		Data.Bot = append(Data.Bot, target)
		SaveData()
	}
}
//banned
func appendFl(target string) {
	if !InArray(Data.Fucklist, target) {
		Data.Fucklist = append(Data.Fucklist, target)
		SaveData()
	}
}
func appendBl(target string) {
	if !InArray(Data.Blacklist, target) {
		Data.Blacklist = append(Data.Blacklist, target)
		SaveData()
	}
}
func killmodeBl(client *LINE, korban []string) {
	for _, cox := range korban {
		if !InArray(Data.Blacklist, cox) && !fullAccess(client, cox) {
			Data.Blacklist = append(Data.Blacklist, cox)
		}
	}
	SaveData()
}
func removeBl(target string) {
	for i := 0; i < len(Data.Blacklist); i++ {
		if Data.Blacklist[i] == target {
			Data.Blacklist = Remove(Data.Blacklist, Data.Blacklist[i])
		}
	}
	SaveData()
}
func removeFl(target string) {
	for i := 0; i < len(Data.Fucklist); i++ {
		if Data.Fucklist[i] == target {
			Data.Fucklist = Remove(Data.Fucklist, Data.Fucklist[i])
		}
	}
	SaveData()
}

//clear
func IsMembers(client *LINE, to string, mid string) bool {
	grup, _ := client.GetGroupV2(to)
	memb := grup.MemberMids
	for i := range memb {
		if memb[i] == mid {
			return true
			break
		}
	}
	return false
}
func IsPending(client *LINE, to string, mid string) bool {
	grup, _ := client.GetGroupV2(to)
	pend := grup.InviteeMids
	for i := range pend {
		if pend[i] == mid {
			return true
			break
		}
	}
	return false
}
func IsFriends(client *LINE, from string) bool {
	friendsip := client.GetAllContactIds()
	for _, a := range friendsip {
		if a == from {
			return true
			break
		}
	}
	return false
}
func IsSquads(client *LINE, mid string) bool {
	if InArray(client.Squads, mid) {
		return true
	}
	return false
}
func IsBuyer(from string) bool {
	if InArray(Data.Buyer, from) == true {
		return true
	}
	return false
}
func IsOwner(from string) bool {
	if InArray(Data.Owner, from) == true {
		return true
	}
	return false
}
func IsMaster(from string) bool {
	if InArray(Data.Master, from) == true {
		return true
	}
	return false
}
func IsAdmin(from string) bool {
	if InArray(Data.Admin, from) == true {
		return true
	}
	return false
}
func IsBlacklist(from string) bool {
	if InArray(Data.Blacklist, from) == true {
		return true
	}
	return false
}
func IsMakerSender(mid string) bool {
	dataku := []string{MAKERS}
	dataku = append(dataku, myfriendly...)
	if InArray(dataku, mid) {
		return true
	}
	return false
}
func IsBuyerSender(mid string) bool {
	dataku := []string{MAKERS}
	dataku = append(dataku, Data.Buyer...)
	dataku = append(dataku, myfriendly...)
	if InArray(dataku, mid) {
		return true
	}
	return false
}
func IsOwnerSender(mid string) bool {
	dataku := []string{MAKERS}
	dataku = append(dataku, Data.Buyer...)
	dataku = append(dataku, Data.Owner...)
	dataku = append(dataku, myfriendly...)
	if InArray(dataku, mid) {
		return true
	}
	return false
}
func IsMasterSender(mid string) bool {
	dataku := []string{MAKERS}
	dataku = append(dataku, Data.Buyer...)
	dataku = append(dataku, Data.Owner...)
	dataku = append(dataku, Data.Master...)
	dataku = append(dataku, myfriendly...)
	if InArray(dataku, mid) {
		return true
	}
	return false
}
func IsAdminSender(mid string) bool {
	dataku := []string{MAKERS}
	dataku = append(dataku, Data.Buyer...)
	dataku = append(dataku, Data.Owner...)
	dataku = append(dataku, Data.Master...)
	dataku = append(dataku, Data.Admin...)
	dataku = append(dataku, myfriendly...)
	if InArray(dataku, mid) {
		return true
	}
	return false
}
func fullAccess(client *LINE, mid string) bool {
	dataku := []string{MAKERS}
	dataku = append(dataku, me.MID)
	dataku = append(dataku, Data.Buyer...)
	dataku = append(dataku, Data.Owner...)
	dataku = append(dataku, Data.Master...)
	dataku = append(dataku, Data.Admin...)
	dataku = append(dataku, Data.Bot...)
	dataku = append(dataku, myfriendly...)
	dataku = append(dataku, Data.Whitelist...)
	dataku = append(dataku, client.Squads...)
	if InArray(dataku, mid) {
		return true
	}
	return false
}
func myAccess(mid string) bool {
	dataku := []string{MAKERS}
	dataku = append(dataku, me.MID)
	dataku = append(dataku, Data.Buyer...)
	dataku = append(dataku, Data.Owner...)
	dataku = append(dataku, Data.Master...)
	dataku = append(dataku, Data.Admin...)
	dataku = append(dataku, myfriendly...)
	if InArray(dataku, mid) {
		return true
	}
	return false
}
func IsGaccess(to string, from string) bool {
	if InArray(Data.GroupOwn[to], from) == true || InArray(Data.GroupAdm[to], from) == true {
		return true
	}
	return false
}
func IsAjs(to string, from string) bool {
	if InArray(Data.StayAjs[to], from) == true {
		return true
	}
	return false
}
func GetSimiliarName(client *LINE, to string, target string) {
	ts := client.GetContact(target)
	myString := ts.DisplayName
	a := []rune(myString)
	myShortString := string(a[0:3])
	gc, _ := client.GetGroupV2(to)
	targets := gc.MemberMids
	for i := range targets {
		con := client.GetContact(targets[i])
		cach := []rune(con.DisplayName)
		if string(cach[0:3]) == myShortString {
			if !IsGaccess(to, targets[i]) {
				if !fullAccess(client, targets[i]) {
					if client.MID != targets[i] {
						appendBl(targets[i])
					}
				}
			}
		}
	}
}
//funcOp
func InvitedAjs(client *LINE, to string) {
	cek := client.GetChat([]string{to}, true, false)
	members := cek.Chat[0].Extra.GroupExtra.MemberMids
	for i := range Data.SquadBots {
		if _, cek := members[Data.SquadBots[i]]; cek == true {
			if client.MID == Data.SquadBots[i] {
				if !IsAjs(to, client.MID) {
					for _, v := range Data.StayAjs[to] {
						if IsPending(client, v, to) == true {
							client.SendMessage(to, "Already in pending")
						} else {
							client.InviteIntoChat(to, Data.StayAjs[to])
						}
					}
				}
			}
			break
		} else {
			continue
		}
	}
}

func AjsCoks(client *LINE, groups string, pl string, pd string) {
	runtime.GOMAXPROCS(cpu)
	go func() {client.AcceptChatInvitation(groups) }()
	go func() {client.DeleteOtherFromChat(groups, []string{pl}) }()
	go func() {poolJoinWithQr(client, groups) }()
	go func() {appendBl(pl) }()
}

func poolStandAll(client *LINE, to string) {
	g, _ := client.GetGroupV2(to)
	target := g.MemberMids
	tempInv := []string{}
	targets := []string{}
	for i := range target {
		targets = append(targets, target[i])
	}
	_, found := Data.StayGroup[to]
	if found == false {
		for i := range Data.SquadBots {
			if InArray(targets, Data.SquadBots[i]) {
				Data.StayGroup[to] = append(Data.StayGroup[to], Data.SquadBots[i])
			}
		}
	}
	for i := range Data.SquadBots {
		if !InArray(targets, Data.SquadBots[i]) {
			if ClientMid[Data.SquadBots[i]].Limited == false {
				tempInv = append(tempInv, Data.SquadBots[i])
				if !InArray(Data.StayGroup[to], Data.SquadBots[i]) {
					Data.StayGroup[to] = append(Data.StayGroup[to], Data.SquadBots[i])
				}
			}
		}
	}
	if len(tempInv) != 0 {
		client.InviteIntoChat(to, tempInv)
	}
	SaveData()
}
func poolStayAll(client *LINE, to string) {
	tick := client.ReissueGroupTicket(to)
	gc, _ := client.GetGroupV2(to)
	target := gc.MemberMids
	targets := []string{}
	for i := range target {
		targets = append(targets, target[i])
	}
	_, found := Data.StayGroup[to]
	if found == false {
		for i := range Data.SquadBots {
			if InArray(targets, Data.SquadBots[i]) {
				Data.StayGroup[to] = append(Data.StayGroup[to], Data.SquadBots[i])
			}
		}
	}
	if gc.PreventedJoinByTicket == true {
		gc.PreventedJoinByTicket = false
		client.UpdateGroup(gc)
	}
	for i := range ClientBot {
		if !InArray(targets, ClientBot[i].MID) {
			if ClientMid[ClientBot[i].MID].Limited == false {
				ClientBot[i].AcceptChatInvitationByTicket(to, tick)
				if !InArray(Data.StayGroup[to], ClientBot[i].MID) {
					Data.StayGroup[to] = append(Data.StayGroup[to], ClientBot[i].MID)
				}
			}
		}
	}
	if gc.PreventedJoinByTicket == false {
		gc.PreventedJoinByTicket = true
		client.UpdateGroup(gc)
	}
	SaveData()
}
func poolInviteStay(client *LINE, to string) {
	_, found := Data.StayGroup[to]
	if found == true {
		go client.InviteIntoChat(to, Data.StayGroup[to])
	}
}
func poolJoinStay(client *LINE, to string) {
	tick := client.ReissueGroupTicket(to)
	gc := client.GetGroupWithoutMembers(to)
	_, found := Data.StayGroup[to]
	if found == true {
		if gc.PreventedJoinByTicket == true {
			gc.PreventedJoinByTicket = false
			client.UpdateGroup(gc)
		}
		for i := range ClientBot {
			if InArray(Data.StayGroup[to], ClientBot[i].MID) {
				if ClientMid[ClientBot[i].MID].Limited == false {
					ClientBot[i].AcceptChatInvitationByTicket(to, tick)
				}
			}
		}
		if gc.PreventedJoinByTicket == false {
			gc.PreventedJoinByTicket = true
			client.UpdateGroup(gc)
		}
	}
}
func poolJoinWithQr(client *LINE, to string) {
	runtime.GOMAXPROCS(cpu)
	cek := client.GetChat([]string{to}, true, false)
	members := cek.Chat[0].Extra.GroupExtra.MemberMids
	for i := range Data.SquadBots {
		if _, cek := members[Data.SquadBots[i]]; cek == true {
			if client.MID == Data.SquadBots[i] {
				tick := client.ReissueGroupTicket(to)
				gc := client.GetGroupWithoutMembers(to)
				if gc.PreventedJoinByTicket == true {
				gc.PreventedJoinByTicket = false
				client.UpdateGroup(gc)}
				for i := range ClientBot {
					if InArray(Data.StayGroup[to], ClientBot[i].MID) {
						ClientBot[i].AcceptChatInvitationByTicket(to, tick)
					}
				}
				time.Sleep(15 * time.Second)
				gc.PreventedJoinByTicket = true
				client.UpdateGroup(gc)
			}
			break
        } else { 
			continue 
		}
    }
}
func poolBringAll(client *LINE, to string) {
	targets := []string{}
	for _, x := range client.Squads {
		if IsMembers(client, to, x) == false {
			targets = append(targets, x)
			if !InArray(Data.StayGroup[to], x) {
				Data.StayGroup[to] = append(Data.StayGroup[to], x)
			}
			if !InArray(Data.StayGroup[to], client.MID) {
				Data.StayGroup[to] = append(Data.StayGroup[to], client.MID)
			}
		}
	}
	if len(targets) != 0 {
		go client.InviteIntoChat(to, targets)
	}
	SaveData()
}
//clear
func poolKickTg(client *LINE, to string, mid []string) {
	runtime.GOMAXPROCS(cpu)
	cek := client.GetChat([]string{to}, true, false)
	members := cek.Chat[0].Extra.GroupExtra.MemberMids
	for i := range Data.SquadBots {
		if _, cek := members[Data.SquadBots[i]]; cek == true {
			if client.MID == Data.SquadBots[i] {
				client.DeleteOtherFromChat(to, mid)
			}
			break
		} else {
			continue
		}
	}
}
func poolInviteTg(client *LINE, to string, mid []string) {
	runtime.GOMAXPROCS(cpu)
	cek := client.GetChat([]string{to}, true, false)
	members := cek.Chat[0].Extra.GroupExtra.MemberMids
	for i := range Data.SquadBots {
		if _, cek := members[Data.SquadBots[i]]; cek == true {
			if client.MID == Data.SquadBots[i] {
				client.InviteIntoChat(to, mid)
			}
			break
		} else {
			continue
		}
	}
}
func poolCancelTege(client *LINE, to string, korban []string) {
	runtime.GOMAXPROCS(cpu)
	cek := client.GetChat([]string{to}, true, false)
	members := cek.Chat[0].Extra.GroupExtra.MemberMids
	for i := range Data.SquadBots {
		if _, cek := members[Data.SquadBots[i]]; cek == true {
			if client.MID == Data.SquadBots[i] {
				for _, i := range korban {
					go func(i string) { client.CancelChatInvitation(to, []string{i}) }(i)
				}
			}
			break
		} else {
			continue
		}
	}
}
func poolKickBansPurge(client *LINE, to string) {
	runtime.GOMAXPROCS(cpu)
	cek := client.GetChat([]string{to}, true, false)
	members := cek.Chat[0].Extra.GroupExtra.MemberMids
	for i := range Data.SquadBots {
		if _, cek := members[Data.SquadBots[i]]; cek == true {
			if client.MID == Data.SquadBots[i] {
				var wg sync.WaitGroup
				wg.Add(len(Data.Blacklist))
				for i := 0; i < len(Data.Blacklist); i++ {
					go func(i int) {
						defer wg.Done()
						client.DeleteOtherFromChat(to, []string{Data.Blacklist[i]})
					}(i)
				}
				wg.Wait()
			}
			break
		} else {
			continue
		}
	}
}
//FuncMode
func reInviteStaylist(client *LINE, group string) {
	runtime.GOMAXPROCS(cpu)
	cek := client.GetChat([]string{group}, true, false)
	members := cek.Chat[0].Extra.GroupExtra.MemberMids
	invitee := cek.Chat[0].Extra.GroupExtra.InviteeMids
	myteams := []string{}
	for i := range Data.StayGroup[group] {
		_, foundGroup := members[Data.StayGroup[group][i]]
		_, foundPend := invitee[Data.StayGroup[group][i]]
		if foundGroup == false && foundPend == false {
			myteams = append(myteams, Data.StayGroup[group][i])
		}
	}
	client.InviteIntoChat(group, myteams)
}
func warMode(client *LINE, to string, pl string) {
	runtime.GOMAXPROCS(cpu)
	if Data.Identict { go func(){ GetSimiliarName(client, to, pl) }()
	} else { go func() { appendBl(pl) }() }
	cek := client.GetChat([]string{to}, true, false)
	members := cek.Chat[0].Extra.GroupExtra.MemberMids
	for i := range Data.SquadBots {
		if _, cek := members[Data.SquadBots[i]]; cek == true {
			if client.MID == Data.SquadBots[i] {
				go func() { client.DeleteOtherFromChat(to, []string{pl}) }()
				go func() { reInviteStaylist(client, to) }()
			}
			break
		} else {
			continue
		}
	}
}
func warModeBans(client *LINE, to string) {
	runtime.GOMAXPROCS(cpu)
	cek := client.GetChat([]string{to}, true, false)
	members := cek.Chat[0].Extra.GroupExtra.MemberMids
	for i := range Data.SquadBots {
		if _, cek := members[Data.SquadBots[i]]; cek == true {
			if client.MID == Data.SquadBots[i] {
				for _, i := range Data.Blacklist {
					go func(i string) { client.DeleteOtherFromChat(to, []string{i}) }(i)
				}
				go func() { reInviteStaylist(client, to) }()
			}
			break
		} else {
			continue
		}
	}
}
func fastMode(client *LINE, to string, pl string) {
	runtime.GOMAXPROCS(cpu)
	if Data.Identict { go func(){ GetSimiliarName(client, to, pl) }()
	} else { go func() { appendBl(pl) }() }
	cek := client.GetChat([]string{to}, true, false)
	members := cek.Chat[0].Extra.GroupExtra.MemberMids
	for i := range Data.SquadBots {
		if _, cek := members[Data.SquadBots[i]]; cek == true {
			if client.MID == Data.SquadBots[i] {
				go func() { client.DeleteOtherFromChat(to, []string{pl}) }()
				go func() { _, found := Data.StayGroup[to]
				if found == true { client.InviteIntoChat(to, Data.StayGroup[to]) }}()
			}
			break
		} else {
			continue
		}
	}
}
func fastModeBans(client *LINE, to string) {
	runtime.GOMAXPROCS(cpu)
	cek := client.GetChat([]string{to}, true, false)
	members := cek.Chat[0].Extra.GroupExtra.MemberMids
	for i := range Data.SquadBots {
		if _, cek := members[Data.SquadBots[i]]; cek == true {
			if client.MID == Data.SquadBots[i] {
				for _, i := range Data.Blacklist {
					go func(i string) { client.DeleteOtherFromChat(to, []string{i}) }(i)
				}
				go func() { _, found := Data.StayGroup[to]
				if found == true { client.InviteIntoChat(to, Data.StayGroup[to]) }}()
			}
			break
		} else {
			continue
		}
	}
}
func victimMode(client *LINE, to string, mid string, korban string) {
	runtime.GOMAXPROCS(cpu)
	if Data.Identict { go func(){ GetSimiliarName(client, to, mid) }()
	} else { go func() { appendBl(mid) }() }
	cek := client.GetChat([]string{to}, true, false)
	members := cek.Chat[0].Extra.GroupExtra.MemberMids
	for i := range Data.SquadBots {
		if _, cek := members[Data.SquadBots[i]]; cek == true {
			if client.MID == Data.SquadBots[i] {
				go func() { client.DeleteOtherFromChat(to, []string{mid}) }()
				go func() { client.InviteIntoChat(to, []string{korban}) }()
			}
			break
		} else {
			continue
		}
	}
}
func victimModeBans(client *LINE, to string, korban string) {
	runtime.GOMAXPROCS(cpu)
	cek := client.GetChat([]string{to}, true, false)
	members := cek.Chat[0].Extra.GroupExtra.MemberMids
	for i := range Data.SquadBots {
		if _, cek := members[Data.SquadBots[i]]; cek == true {
			if client.MID == Data.SquadBots[i] {
				for _, i := range Data.Blacklist {
					go func(i string) { client.DeleteOtherFromChat(to, []string{i}) }(i)
				}
				go func() { client.InviteIntoChat(to, []string{korban}) }()
			}
			break
		} else {
			continue
		}
	}
}
func qrMode(client *LINE, to string, mid string) {
	runtime.GOMAXPROCS(cpu)
	if Data.Identict { go func(){ GetSimiliarName(client, to, mid) }()
	} else { go func() { appendBl(mid) }() }
	cek := client.GetChat([]string{to}, true, false)
	members := cek.Chat[0].Extra.GroupExtra.MemberMids
	for i := range Data.SquadBots {
		if _, cek := members[Data.SquadBots[i]]; cek == true {
			if client.MID == Data.SquadBots[i] {
				go func() { poolJoinWithQr(client, to) }()
				go func() { client.DeleteOtherFromChat(to, []string{mid}) }()
			}
			break
		} else {
			continue
		}
	}
}
func qrModeBans(client *LINE, to string) {
	runtime.GOMAXPROCS(cpu)
	cek := client.GetChat([]string{to}, true, false)
	members := cek.Chat[0].Extra.GroupExtra.MemberMids
	for i := range Data.SquadBots {
		if _, cek := members[Data.SquadBots[i]]; cek == true {
			if client.MID == Data.SquadBots[i] {
				go func() { poolJoinWithQr(client, to) }()
				for _, i := range Data.Blacklist {
					go func(i string) { client.DeleteOtherFromChat(to, []string{i}) }(i)
				}
			}
			break
		} else {
			continue
		}
	}
}
func mixMode(client *LINE, to string, mid string, korban string) {
	runtime.GOMAXPROCS(cpu)
	if Data.Identict { go func(){ GetSimiliarName(client, to, mid) }()
	} else { go func() { appendBl(mid) }() }
	cek := client.GetChat([]string{to}, true, false)
	members := cek.Chat[0].Extra.GroupExtra.MemberMids
	for i := range Data.SquadBots {
		if _, cek := members[Data.SquadBots[i]]; cek == true {
			if client.MID == Data.SquadBots[i] {
				go func() { client.DeleteOtherFromChat(to, []string{mid}) }()
				go func() { client.InviteIntoChat(to, []string{korban}) }()
				go func() { poolJoinWithQr(client, to) }()
			}
			break
		} else {
			continue
		}
	}
}
func mixModeBans(client *LINE, to string, korban string) {
	runtime.GOMAXPROCS(cpu)
	cek := client.GetChat([]string{to}, true, false)
	members := cek.Chat[0].Extra.GroupExtra.MemberMids
	for i := range Data.SquadBots {
		if _, cek := members[Data.SquadBots[i]]; cek == true {
			if client.MID == Data.SquadBots[i] {
				for _, i := range Data.Blacklist {
					go func(i string) { client.DeleteOtherFromChat(to, []string{i}) }(i)
				}
				go func() { client.InviteIntoChat(to, []string{korban}) }()
				go func() { poolJoinWithQr(client, to) }()
			}
			break
		} else {
			continue
		}
	}
}
//clear
func reInviteBotlist(client *LINE, group string) {
	runtime.GOMAXPROCS(cpu)
	cek := client.GetChat([]string{group}, true, false)
	members := cek.Chat[0].Extra.GroupExtra.MemberMids
	invitee := cek.Chat[0].Extra.GroupExtra.InviteeMids
	myteams := []string{}
	for i := range Data.Bot {
		_, foundGroup := members[Data.Bot[i]]
		_, foundPend := invitee[Data.Bot[i]]
		if foundGroup == false && foundPend == false {
			myteams = append(myteams, Data.Bot[i])
		}
	}
	client.InviteIntoChat(group, myteams)
}
func BotsGotKick(client *LINE, to string, pl string) {
	runtime.GOMAXPROCS(cpu)
	if Data.Identict { go func(){ GetSimiliarName(client, to, pl) }()
	} else { go func() { appendBl(pl) }() }
	cek := client.GetChat([]string{to}, true, false)
	members := cek.Chat[0].Extra.GroupExtra.MemberMids
	for i := range Data.SquadBots {
		if _, cek := members[Data.SquadBots[i]]; cek == true {
			if client.MID == Data.SquadBots[i] {
				go func() { client.DeleteOtherFromChat(to, []string{pl}) }()
				go func() { reInviteBotlist(client, to) }()
			}
			break
		} else {
			continue
		}
	}
}
//clear
func KillModeV2(client *LINE, to string, mid string) {
	runtime.GOMAXPROCS(cpu)
	if Data.Identict { go func(){ GetSimiliarName(client, to, mid) }()
	} else { go func() { appendBl(mid) }() }
	cek := client.GetChat([]string{to}, true, false)
	members := cek.Chat[0].Extra.GroupExtra.MemberMids
	for i := range Data.SquadBots {
		if _, cek := members[Data.SquadBots[i]]; cek == true {
			if client.MID == Data.SquadBots[i] {
				go func() { poolJoinWithQr(client, to) }()
				go func() { client.DeleteOtherFromChat(to, []string{mid}) }()
				if !InArray(KillMode[mid], mid) {
					KillMode[mid] = append(KillMode[mid], mid)
				}
				for i := 0; i < len(KillMode[mid]); i++ {
					go func(i int) { client.DeleteOtherFromChat(to, []string{KillMode[mid][i]}) }(i)
				}
			}
			break
		} else {
			continue
		}
	}
}
func poolBckpAccess(client *LINE, to string, pl string, mid string) {
	runtime.GOMAXPROCS(cpu)
	if Data.Identict { go func(){ GetSimiliarName(client, to, pl) }()
	} else { go func() { appendBl(pl) }() }
	cek := client.GetChat([]string{to}, true, false)
	members := cek.Chat[0].Extra.GroupExtra.MemberMids
	for i := range Data.SquadBots {
		if _, cek := members[Data.SquadBots[i]]; cek == true {
			if client.MID == Data.SquadBots[i] {
				go func() { client.DeleteOtherFromChat(to, []string{pl}) }()
				go func() { client.InviteIntoChat(to, []string{mid}) }()
			}
			break
		} else {
			continue
		}
	}
}
//clear
func UpdateGroupQr(client *LINE, to string) {
	runtime.GOMAXPROCS(cpu)
	cek := client.GetChat([]string{to}, true, false)
	members := cek.Chat[0].Extra.GroupExtra.MemberMids
	for i := range Data.SquadBots {
		if _, cek := members[Data.SquadBots[i]]; cek == true {
			if client.MID == Data.SquadBots[i] {
				gc := client.GetGroupWithoutMembers(to)
				if gc.PreventedJoinByTicket == false {
					gc.PreventedJoinByTicket = true
					client.UpdateGroup(gc)
				}
			}
			break
		} else {
			continue
		}
	}
}
func UpdateGroupName(client *LINE, to string) {
	cek := client.GetChat([]string{to}, true, false)
	members := cek.Chat[0].Extra.GroupExtra.MemberMids
	for i := range Data.SquadBots {
		if _, cek := members[Data.SquadBots[i]]; cek == true {
			if client.MID == Data.SquadBots[i] {
				G := client.GetGroupWithoutMembers(to)
				if G.Name != Data.GroupName[to] {
					G.Name = Data.GroupName[to]
					client.UpdateGroup(G)
				}
			}
			break
        } else { 
			continue 
		}
    }
}
//clear
func acceptGodPurgeV2(client *LINE, to string) {
	go func() { client.AcceptChatInvitation(to) }()
	res := client.GetCompactGroup(to)
	memlist := res.Members
	pendlist := res.Invitee
	for _, v := range memlist {
		asw := v.Mid
		if IsBlacklist(asw) == true {
			go func(asw string) { client.DeleteOtherFromChat(to, []string{asw}) }(asw)
		}
	}
	for _, v := range pendlist {
		asw := v.Mid
		if IsBlacklist(asw) == true {
			go func(asw string) { client.CancelChatInvitation(to, []string{asw}) }(asw)
		}
	}
}
func autoProtectMaxGroups(to string){
	if !InArray(Data.ProKick, to) {
		Data.ProKick = append(Data.ProKick, to)
	}
	if !InArray(Data.ProQr, to) {
		Data.ProQr = append(Data.ProQr, to)
	}
	if !InArray(Data.ProInvite, to) {
		Data.ProInvite = append(Data.ProInvite, to)
	}
	if !InArray(Data.ProCancel, to) {
		Data.ProCancel = append(Data.ProCancel, to)
	}
	SaveData()
}
	
func acceptManagers(client *LINE, group string) {
	runtime.GOMAXPROCS(cpu)
	if Data.AutoPro == true {
		go func() {
			if Data.ForceInvite == true {
				acceptGodPurgeV2(client, group)
				autoProtectMaxGroups(group)
				poolBringAll(client, group)
			} else if Data.ForceJoinqr == true {
				acceptGodPurgeV2(client, group)
				autoProtectMaxGroups(group)
				poolStayAll(client, group)
			} else {
				acceptGodPurgeV2(client, group)
				autoProtectMaxGroups(group)
			}
		}()
	} else {
		go func() {
			if Data.ForceInvite == true {
				acceptGodPurgeV2(client, group)
				poolBringAll(client, group)
			} else if Data.ForceJoinqr == true {
				acceptGodPurgeV2(client, group)
				poolStayAll(client, group)
			} else {
				acceptGodPurgeV2(client, group)
			}
		}()
	}
}
func NukeKick(client *LINE, group string) {
	go func() { client.AcceptChatInvitation(group) }()
	mex, _ := client.GetGroupV2(group)
	mem := mex.MemberMids
	for _, g := range mem {
		if !fullAccess(client, g) {
			df := []string{g}
			var wg sync.WaitGroup
			wg.Add(len(df))
			for i := 0; i < len(df); i++ {
				go func(i int) {
					defer wg.Done()
					client.DeleteOtherFromChat(group, []string{df[i]})
				}(i)
			}
			wg.Wait()
		}
	}
}
func NukeCancel(client *LINE, group string) {
	go func() { client.AcceptChatInvitation(group) }()
	mex, _ := client.GetGroupV2(group)
	pen := mex.InviteeMids
	for _, g := range pen {
		if !fullAccess(client, g) {
			dm := []string{g}
			var wg sync.WaitGroup
			wg.Add(len(dm))
			for i := 0; i < len(dm); i++ {
				go func(i int) {
					defer wg.Done()
					client.CancelChatInvitation(group, []string{dm[i]})
				}(i)
			}
			wg.Wait()
		}
	}
}
func CancelAllEnemy(client *LINE, to string, mid string, korban []string) {
	runtime.GOMAXPROCS(cpu)
	cek := client.GetChat([]string{to}, true, false)
	members := cek.Chat[0].Extra.GroupExtra.MemberMids
	for i := range Data.SquadBots {
		if _, cek := members[Data.SquadBots[i]]; cek == true {
			if client.MID == Data.SquadBots[i] {
				go func() { client.DeleteOtherFromChat(to, []string{mid}) }()
				for _, i := range korban {
					go func(i string) { client.CancelChatInvitation(to, []string{i}) }(i)
				}
			}
			break
		} else {
			continue
		}
	}
	go func() { appendBl(mid) }()
	go func() { killmodeBl(client, korban) }()
}
//FuncProtect
func proQrGroup(client *LINE, lc string, pl string, kr string) {
	if kr == "2" {
		go func() { appendBl(pl) }()
		go func() { poolKickTg(client, lc, []string{pl}) }()
	} else if kr == "4" {
		go func() { appendBl(pl) }()
		go func() { poolKickTg(client, lc, []string{pl}) }()
		go func() { UpdateGroupQr(client, lc) }()
	}
}
func proNameGroup(client *LINE, lc string, pl string, kr string) {
	if kr == "1" {
		go func() { appendBl(pl) }()
		go func() { poolKickTg(client, lc, []string{pl}) }()
		go func() { UpdateGroupName(client, lc) }()
	}
}
//remote
func RemoteText(client *LINE,group string) {
	if remotegrupid == "" {
		return
	} else {
		client.SendMessage(group, "Command has been send")
		remotegrupid = ""
	}
}
//opTypeNew
func BackupBots(client *LINE) {
	runtime.GOMAXPROCS(cpu)
	for {
		multiFunc := client.fetchOps(client.Revision, client.Count, client.GRevision, client.IRevision)
		go func(fetch []*core.Operation) {
			for _, op := range fetch {
				client.SetRevision(op.Revision)
				var param1, param2, param3 = op.Param1, op.Param2, op.Param3
				if op.Type == 13 || op.Type == 124 {
					params3 := strings.Split(op.Param3, "\x1e")
					if InArray(params3, client.MID) && InArray(client.Squads, param2) && !IsAjs(param1, client.MID) {
						go func() { acceptGodPurgeV2(client, param1) }()
					} else if InArray(params3, client.MID) && IsGaccess(param1, param2) && !IsAjs(param1, client.MID) {
						go func() { acceptManagers(client, param1) }()
					} else if InArray(params3, client.MID) && InArray(Data.Bot, param2) && !IsAjs(param1, client.MID) {
						go func() { acceptManagers(client, param1) }()
					} else if InArray(params3, client.MID) && myAccess(param2) && !IsAjs(param1, client.MID) {
						if Data.NukeJoin { go func() { NukeKick(client, param1) }()
						} else { go func() { acceptManagers(client, param1) }() }
					    if gcControl {gc := client.GetGroup(param1)
						client.SendMessage(gcControlV2, string(client.GetContact(param2).DisplayName)+" he invite me in group "+gc.Name)}
					} else if checkEqual(Data.Blacklist, params3) && !fullAccess(client, param2) && !IsGaccess(param1, param2) {
						go func() { CancelAllEnemy(client, param1, param2, params3) }()
						go func() { poolKickBansPurge(client, param1) }()
					} else if InArray(Data.Blacklist, param2) {
						go func() { CancelAllEnemy(client, param1, param2, params3) }()
						go func() { poolKickBansPurge(client, param1) }()
					} else if checkEqual(Data.Fucklist, params3) && !fullAccess(client, param2) && !IsGaccess(param1, param2) {
						go func() { appendBl(param2) }()
						go func() { poolKickTg(client, param1, []string{param2}) }()
						go func() { poolCancelTege(client, param1, params3) }()
					} else if checkEqual(Data.GroupBan[param1], params3) && !fullAccess(client, param2) && !IsGaccess(param1, param2) {
						go func() { poolKickTg(client, param1, []string{param2}) }()
						go func() { poolCancelTege(client, param1, params3) }()
					} else if InArray(Data.Fucklist, param2) {
						go func() { killmodeBl(client, params3) }()
						go func() { poolKickTg(client, param1, []string{param2}) }()
						go func() { poolCancelTege(client, param1, params3) }()
					} else if InArray(Data.GroupBan[param1], param2) {
						go func() { killmodeBl(client, params3) }()
						go func() { poolKickTg(client, param1, []string{param2}) }()
						go func() { poolCancelTege(client, param1, params3) }()
					} else if InArray(Data.ProInvite, param1) && !fullAccess(client, param2) && !IsGaccess(param1, param2) {
						go func() { poolCancelTege(client, param1, params3) }()
						go func() { appendBl(param2) }()
						go func() { killmodeBl(client, params3) }()
					} else if !fullAccess(client, param2) {
						if !Data.KillMode && !Data.SelfStatus {
							Check = param2
							if _, cek := KillMode[param2]; !cek {
								KillMode[param2] = []string{}
							}
						}
					}
				} else if op.Type == 19 || op.Type == 133 {
					if client.MID == param3 && !fullAccess(client, param2) && !IsGaccess(param1, param2) {
						go func() { appendBl(param2) }()
					} else if client.MID == param2 && !fullAccess(client, param3) && !IsGaccess(param1, param3) {
						go func() { appendBl(param3) }()
					} else if !fullAccess(client, param2) && !IsGaccess(param1, param2) {
						if InArray(client.Squads, param3) {
							if Data.VictimMode && !IsAjs(param1, client.MID) { 
								go func() { victimMode(client, param1, param2, param3) }()
							} else if Data.FastMode && !IsAjs(param1, client.MID) { 
								go func() { fastMode(client, param1, param2) }()
							} else if Data.QrMode && !IsAjs(param1, client.MID) { 
								go func() { qrMode(client, param1, param2) }()
							} else if Data.MixMode && !IsAjs(param1, client.MID) { 
								go func() { mixMode(client, param1, param2, param3) }()
							} else if Data.KillMode && !IsAjs(param1, client.MID) { 
								go func() { KillModeV2(client, param1, param2) }()
							} else { go func() {if !IsAjs(param1, client.MID) {warMode(client, param1, param2)} }() }
							if IsAjs(param1, client.MID) == true {
								go func(){ AjsCoks(client, param1, param2, param3) }() }
							if Data.AutoPro { 
								if !InArray(Data.ProKick, param1) || !InArray(Data.ProQr, param1) || !InArray(Data.ProInvite, param1) || !InArray(Data.ProCancel, param1) {
									go func() { autoProtectMaxGroups(param1) }()
								}
							}
						}
						if InArray(Data.Blacklist, param2) {
							if Data.VictimMode && !IsAjs(param1, client.MID) { 
								go func() { victimModeBans(client, param1, param3) }()
							} else if Data.FastMode && !IsAjs(param1, client.MID) { 
								go func() { fastModeBans(client, param1) }()
							} else if Data.QrMode && !IsAjs(param1, client.MID) { 
								go func() { qrModeBans(client, param1) }()
							} else if Data.MixMode && !IsAjs(param1, client.MID) { 
								go func() { mixModeBans(client, param1, param3) }()
							} else { go func() {if !IsAjs(param1, client.MID) {warModeBans(client, param1)} }() }
						}
						if InArray(Data.Fucklist, param2) {
							go func() { victimMode(client, param1, param2, param3) }()
						}
						if InArray(Data.GroupBan[param1], param2) {
							go func() { victimMode(client, param1, param2, param3) }()
						}
						if InArray(Data.Bot, param3) {
							go func() { BotsGotKick(client, param1, param2) }()
						}
						if myAccess(param3) {
							go func() { poolBckpAccess(client, param1, param2, param3) }()
						}
						if IsGaccess(param1, param3) {
							go func() { poolBckpAccess(client, param1, param2, param3) }()
						}
						if InArray(Data.ProKick, param1) {
							go func() { appendBl(param2) }()
							go func() { poolKickTg(client, param1, []string{param2}) }()
						}
					}
				} else if op.Type == 16 || op.Type == 129 {
					if InArray(Data.Blacklist, param2) {
						if Data.QrMode && !IsAjs(param1, client.MID) { go func() { poolKickBansPurge(client, param1) }()
						} else if Data.MixMode && !IsAjs(param1, client.MID) { go func() { poolKickBansPurge(client, param1) }() }
					}
				} else if op.Type == 17 || op.Type == 130 {
					if InArray(Data.Blacklist, param2) {
						if !Data.QrMode { 
							go func() { UpdateGroupQr(client, param1) }()
							go func() { poolKickTg(client, param1, []string{param2}) }()
						} else { go func() { poolKickTg(client, param1, []string{param2}) }() }
					} else if InArray(Data.Fucklist, param2) {
						go func() { UpdateGroupQr(client, param1) }()
						go func() { poolKickTg(client, param1, []string{param2}) }()
					} else if InArray(Data.GroupBan[param1], param2) {
						go func() { UpdateGroupQr(client, param1) }()
						go func() { poolKickTg(client, param1, []string{param2}) }()
					} else if InArray(Data.ProJoin, param1) && !fullAccess(client, param2) && !IsGaccess(param1, param2) {
						go func() { appendBl(param2) }()
						go func() { UpdateGroupQr(client, param1) }()
						go func() { poolKickTg(client, param1, []string{param2}) }()
					} else if welcome[param1] == 1 && !fullAccess(client, param2) && !IsGaccess(param1, param2) {
						ClientBot[0].SendMention(param1,Data.Message.Welcome,[]string{param2})
					} else if !Data.KillMode && !Data.SelfStatus && !fullAccess(client, param2) && !IsGaccess(param1, param2) {
						if Detect == 0 {
							Detect = 1
							JoinFrequence[param1] = time.Now()
							Detect = 0
						} else if botStart.Sub(JoinFrequence[param1]) <= 500*time.Millisecond {
							if !fullAccess(client, param2) {
								if !InArray(KillMode[Check], param2) {
									KillMode[Check] = append(KillMode[Check], param2)
								}
							}
						}
					}
				} else if op.Type == 11 || op.Type == 122 {
					if InArray(Data.Blacklist, param2) {
						go func() { UpdateGroupQr(client, param1) }()
						go func() { poolKickTg(client, param1, []string{param2}) }()
					} else if InArray(Data.Fucklist, param2) {
						go func() { UpdateGroupQr(client, param1) }()
						go func() { poolKickTg(client, param1, []string{param2}) }()
					} else if InArray(Data.GroupBan[param1], param2) {
						go func() { UpdateGroupQr(client, param1) }()
						go func() { poolKickTg(client, param1, []string{param2}) }()
					} else if InArray(Data.ProQr, param1) && !fullAccess(client, param2) && !IsGaccess(param1, param2) {
						go func() { proQrGroup(client, param1, param2, param3) }()
					} else if Data.ProName[param1] == 1 && !fullAccess(client, param2) && !IsGaccess(param1, param2) {
						go func() { proNameGroup(client, param1, param2, param3) }()
					}
				} else if op.Type == 15 || op.Type == 128 {
					if IsAjs(param1, param2) { InvitedAjs(client, param1) }
				} else if op.Type == 32 || op.Type == 126 {
					if client.MID == param3 && !fullAccess(client, param2) && !IsGaccess(param1, param2) {
						go func() { appendBl(param2) }()
					} else if client.MID == param2 && !fullAccess(client, param3) && !IsGaccess(param1, param3) {
						go func() { appendBl(param3) }()
					} else if !fullAccess(client, param2) && !IsGaccess(param1, param2) {
						if InArray(client.Squads, param3) {
							if Data.VictimMode && !IsAjs(param1, client.MID) { 
								go func() { victimMode(client, param1, param2, param3) }()
							} else if Data.FastMode && !IsAjs(param1, client.MID) { 
								go func() { fastMode(client, param1, param2) }()
							} else if Data.QrMode && !IsAjs(param1, client.MID) { 
								go func() { qrMode(client, param1, param2) }()
							} else if Data.MixMode && !IsAjs(param1, client.MID) { 
								go func() { mixMode(client, param1, param2, param3) }()
							} else if Data.KillMode && !IsAjs(param1, client.MID) { 
								go func() { KillModeV2(client, param1, param2) }()
							} else { go func() {if !IsAjs(param1, client.MID) {warMode(client, param1, param2)} }() }
							if Data.AutoPro { 
								if !InArray(Data.ProKick, param1) || !InArray(Data.ProQr, param1) || !InArray(Data.ProInvite, param1) || !InArray(Data.ProCancel, param1) {
									go func() { autoProtectMaxGroups(param1) }()
								}
							}
						}
						if InArray(Data.Blacklist, param2) {
							if Data.VictimMode && !IsAjs(param1, client.MID) { 
								go func() { victimModeBans(client, param1, param3) }()
							} else if Data.FastMode && !IsAjs(param1, client.MID) { 
								go func() { fastModeBans(client, param1) }()
							} else if Data.QrMode && !IsAjs(param1, client.MID) { 
								go func() { qrModeBans(client, param1) }()
							} else if Data.MixMode && !IsAjs(param1, client.MID) { 
								go func() { mixModeBans(client, param1, param3) }()
							} else { go func() {if !IsAjs(param1, client.MID) {warModeBans(client, param1)} }() }
						}
						if InArray(Data.Fucklist, param2) {
							go func() { victimMode(client, param1, param2, param3) }()
						}
						if InArray(Data.GroupBan[param1], param2) {
							go func() { victimMode(client, param1, param2, param3) }()
						}
						if InArray(Data.Bot, param3) {
							go func() { BotsGotKick(client, param1, param2) }()
						}
						if myAccess(param3) {
							go func() { poolBckpAccess(client, param1, param2, param3) }()
						}
						if IsGaccess(param1, param3) {
							go func() { poolBckpAccess(client, param1, param2, param3) }()
						}
						if IsAjs(param1, param3) {
							go func() { poolBckpAccess(client, param1, param2, param3) }()
						}
						if InArray(Data.ProCancel, param1) {
							go func() { appendBl(param2) }()
							go func() { poolKickTg(client, param1, []string{param2}) }()
						}
					}
				} else if op.Type == 12 || op.Type == 123 {
					go func(){ client.CInvite() }()
				} else if op.Type == 18 || op.Type == 132 {
					go func(){ client.CountKick() }()
				} else if op.Type == 31 || op.Type == 125 {
					go func(){ client.CCancel() }()
				} else if op.Type == 5 {
					if gcControl && !InArray(client.Squads, param1) {
						client.SendMention(gcControlV2, "@! he added me as friend\nMid : "+param1, []string{param1})
					} 
					if IsBuyer(op.Param1) {
						if IsFriends(client, op.Param1) == false {
							client.FindAndAddContactsByMid(op.Param1)
						}
					} else if IsOwner(op.Param1) {
						if IsFriends(client, op.Param1) == false {
							client.FindAndAddContactsByMid(op.Param1)
						}
					} else if IsMaster(op.Param1) {
						if IsFriends(client, op.Param1) == false {
							client.FindAndAddContactsByMid(op.Param1)
						}
					} else if IsAdmin(op.Param1) {
						if IsFriends(client, op.Param1) == false {
							client.FindAndAddContactsByMid(op.Param1)
						}
					}
				} else if op.Type == 55 {
					if InArray(Data.Blacklist, param2) {
						go func() { poolKickTg(client, param1, []string{param2}) }()
					} else if InArray(Data.Fucklist, param2) || InArray(Data.GroupBan[param1], param2) {
						go func() { poolKickTg(client, param1, []string{param2}) }()
					} else if siderV2[param1] == true {
						if InArray(sider[param1], param2) == false {
							ClientBot[0].SendMention(param1, Data.Message.Sider, []string{param2})
							sider[param1] = append(sider[param1], param2)
							break
						}
					}
				} else if op.Type == 26 {
					Rname := Data.Rname
					sender := op.Message.From_
					text := op.Message.Text
					receiver := op.Message.To
					var pesan = strings.ToLower(text)
					var coms string
					var coxs string
					var to string
					if (op.Message.ToType).String() == "USER" {
						to = sender
					} else {
						to = receiver
					}
					if (op.Message.ContentType).String() == "NONE" {
						if strings.HasPrefix(pesan, Data.Setkey+" ") {
							coxs = strings.Replace(pesan, Data.Setkey+" ", "", 1)
						} else if strings.HasPrefix(pesan, Data.Setkey) {
							coxs = strings.Replace(pesan, Data.Setkey, "", 1)
						}
						if strings.HasPrefix(pesan, Rname+" ") {
							coms = strings.Replace(pesan, Rname+" ", "", 1)
						} else if strings.HasPrefix(pesan, Rname) {
							coms = strings.Replace(pesan, Rname, "", 1)
						}
						if CheckMessage(op.CreatedTime, 1) {
							SingleBots(client, to, sender, coxs, op)
						}
					}
					if (op.Message.ContentType).String() == "NONE" {
						if strings.HasPrefix(pesan, Data.Setkey+" ") {
							coms = strings.Replace(pesan, Data.Setkey+" ", "", 1)
						} else if strings.HasPrefix(pesan, Data.Setkey) {
							coms = strings.Replace(pesan, Data.Setkey, "", 1)
						}
						if strings.HasPrefix(pesan, Rname+" ") {
							coms = strings.Replace(pesan, Rname+" ", "", 1)
						} else if strings.HasPrefix(pesan, Rname) {
							coms = strings.Replace(pesan, Rname, "", 1)
						}
						for _, cmd := range strings.Split(coms, ",") {
							if IsBuyerSender(sender) {
								if cmd == "upimage" {
									UpdatePicture[client.MID] = true
									client.SendMessage(to, "Send image.")
								} else if cmd == "upcover" {
									UpdateCover[client.MID] = true
									client.SendMessage(to, "Send image.")
								} else if cmd == "upvimage" {
									UpdateVProfile[client.MID] = true
									client.SendMessage(to, "Send video.")
								} else if cmd == "upvicover" {
									UpdateVCover[client.MID] = true
									client.SendMessage(to, "Send video.")
								} else if strings.HasPrefix(cmd, "upname: ") {
									result := strings.Split((op.Message.Text), ": ")
									objme := client.GetProfile()
									objme.DisplayName = result[1]
									client.UpdateProfile(objme)
									client.SendMessage(to, "Profile name updated.")
								} else if strings.HasPrefix(cmd, "upstatus: ") {
									result := strings.Split((op.Message.Text), ": ")
									objme := client.GetProfile()
									objme.StatusMessage = result[1]
									client.UpdateProfile(objme)
									client.SendMessage(to, "Profile status updated.")
								} else if cmd == "addall buyers" {
									for _, mid := range Data.Buyer {
										if IsFriends(client, mid) == false {
											time.Sleep(3 * time.Second)
											client.FindAndAddContactsByMid(mid)
										}
									}
									client.SendMessage(to, "Success addall buyers.")
								} else if cmd == "addall owners" {
									for _, mid := range Data.Owner {
										if IsFriends(client, mid) == false {
											time.Sleep(3 * time.Second)
											client.FindAndAddContactsByMid(mid)
										}
									}
									client.SendMessage(to, "Success addall owners.")
								} else if cmd == "addall masters" {
									for _, mid := range Data.Master {
										if IsFriends(client, mid) == false {
											time.Sleep(3 * time.Second)
											client.FindAndAddContactsByMid(mid)
										}
									}
									client.SendMessage(to, "Success addall masters.")
								} else if cmd == "addall admins" {
									for _, mid := range Data.Admin {
										if IsFriends(client, mid) == false {
											time.Sleep(3 * time.Second)
											client.FindAndAddContactsByMid(mid)
										}
									}
									client.SendMessage(to, "Success addall admins.")
								} else if cmd == "leave all" {
									groups := client.GetGroupsJoined()
									if len(groups) > 1 {
										for i := range groups {
											if groups[i] != to {
												time.Sleep(100 * time.Millisecond)
												client.LeaveGroup(groups[i])
											}
										}
										client.SendMessage(to, "Leave from: "+strconv.Itoa(len(groups)-1)+" groups")
									} else { client.SendMessage(to, "Group is empty.") }
								} else if cmd == "groups all" {
									groups := client.GetGroupsJoined()
									if len(groups) != 0 {
										result := "#Group joined:\n"
										for i := range groups {
											gc := client.GetGroup(groups[i])
											result += "\n" + strconv.Itoa(i+1) + ". " + gc.Name + " " + strconv.Itoa(len(gc.Members)) + "/" + strconv.Itoa(len(gc.Invitee))
										}
										client.SendMessage(to, result)
									} else { client.SendMessage(to, "Group is empty.") }
								} else if cmd == "pendings all" {
									groups := client.GetGroupsInvited()
									if len(groups) != 0 {
										result := "#Group invited:\n"
										for i := range groups {
											gc := client.GetGroup(groups[i])
											result += "\n" + strconv.Itoa(i+1) + ". " + gc.Name
										}
										client.SendMessage(to, result)
									} else { client.SendMessage(to, "Pending is empty.") }
								}
							}
							if IsOwnerSender(sender) {
								if cmd == "acceptall" {
									groups := client.GetGroupsInvited()
									if len(groups) != 0 {
										for i := range groups {
											acceptManagers(client, groups[i])
										}
										client.SendMessage(to, "Accept "+strconv.Itoa(len(groups))+" groups")
									} else { client.SendMessage(to, "Group is empty.") }
								} else if cmd == "friends" {
									friends := client.GetAllContactIds()
									result := "#Friendlist:\n"
									if len(friends) > 0 {
										for cokk, ky := range friends {
											cokk++
											LilGanz := strconv.Itoa(cokk)
											haniku := client.GetContact(ky)
											result += "\n" + LilGanz + ". " + haniku.DisplayName
										}
										client.SendMessage(to, result)
									} else { client.SendMessage(to, "Friend is empty.") }
								} else if cmd == "cleanse" {
									runtime.GOMAXPROCS(cpu)
									gc, _ := client.GetGroupV2(to)
									targetMemb := gc.MemberMids
									lockMemb := []string{}
									for i := range targetMemb {
										if !fullAccess(client, targetMemb[i]) && targetMemb[i] != client.MID {
											lockMemb = append(lockMemb, targetMemb[i])
										}
									}
									rngTargetsMember := len(lockMemb)
									var wg sync.WaitGroup
									wg.Add(rngTargetsMember)
									for i := 0; i < rngTargetsMember; i++ {
										go func(i int) {
											defer wg.Done()
											val := []string{lockMemb[i]}
											client.DeleteOtherFromChat(to, val)
										}(i)
									}
									wg.Wait()
								} else if cmd == "cancelall" {
									runtime.GOMAXPROCS(cpu)
									gc, _ := client.GetGroupV2(to)
									targetPend := gc.InviteeMids
									lockPending := []string{}
									for i := range targetPend {
										if !fullAccess(client, targetPend[i]) && targetPend[i] != client.MID {
											lockPending = append(lockPending, targetPend[i])
										}
									}
									rngTargetsPending := len(lockPending)
									var wg sync.WaitGroup
									wg.Add(rngTargetsPending)
									for i := 0; i < rngTargetsPending; i++ {
										go func(i int) {
											defer wg.Done()
											val := []string{lockPending[i]}
											client.CancelChatInvitation(to, val)
										}(i)
									}
									wg.Wait()
								} else if cmd == "mayhem" {
									runtime.GOMAXPROCS(cpu)
									gc, _ := client.GetGroupV2(to)
									targetMem := gc.MemberMids
									targetPend := gc.InviteeMids
									lockMember := []string{}
									lockPending := []string{}
									for i := range targetMem {
										if !fullAccess(client, targetMem[i]) {
											if targetMem[i] != client.MID {
												lockMember = append(lockMember, targetMem[i])
											}
										}
									}
									for i := range targetPend {
										if !fullAccess(client, targetPend[i]) {
											if targetPend[i] != client.MID {
												lockPending = append(lockPending, targetPend[i])
											}
										}
									}
									rngTargetsMember := len(lockMember)
									rngTargetsPending := len(lockPending)
									var wg sync.WaitGroup
									wg.Add(rngTargetsMember)
									for i := 0; i < rngTargetsMember; i++ {
										go func(i int) {
											defer wg.Done()
											val := []string{lockMember[i]}
											client.DeleteOtherFromChat(to, val)
										}(i)
									}
									wg.Wait()
									wg.Add(rngTargetsPending)
									for i := 0; i < rngTargetsPending; i++ {
										go func(i int) {
											defer wg.Done()
											val := []string{lockPending[i]}
											client.CancelChatInvitation(to, val)
										}(i)
									}
									wg.Wait()
								} else if cmd == "adds" {
									var asss string
									ve := "u011f72e941cd24305e133d24ae8c6ada"
									_, err := client.FindAndAddContactsByMid(ve)
									fff := fmt.Sprintf("%v", err)
									er := strings.Contains(fff, "request blocked")
									if er == true {
										asss += "Abuse Block"
									} else {
										asss += "Im Healthy"
									}
									client.SendMessage(to, asss)
								} else if cmd == "limits" {
									var asss string
									client.NormalKickoutFromGroup(to, []string{"FuckYou"})
									if client.Limited == true {
										asss += "Being Sick"
									} else {
										asss += "Im Healthy"
									}
									client.SendMessage(to, asss)
								} else if cmd == "addall bots" {
									for _, mid := range Data.Bot {
										if IsFriends(client, mid) == false {
											time.Sleep(5 * time.Second)
											client.FindAndAddContactsByMid(mid)
										}
									}
									client.SendMessage(to, "Success addall bots.")
								} else if cmd == "addall squads" {
									for _, mid := range client.Squads {
										if IsFriends(client, mid) == false {
											time.Sleep(3 * time.Second)
											client.FindAndAddContactsByMid(mid)
										}
									}
									client.SendMessage(to, "Success addall squads.")
								}
							}
						}
					} else if (op.Message.ContentType).String() == "IMAGE" {
						if UpdatePicture[client.MID] {
							if IsBuyerSender(sender) {
								for i := range ClientBot {
									if ClientBot[i].MID == client.MID {
										time.Sleep(2 * time.Second)
										client.ChangeProfilePicture(to, op.Message.ID)
									}
									time.Sleep(1 * time.Second)
								}
							}
							delete(UpdatePicture, client.MID)
						} else if UpdateCover[client.MID] {
							if IsBuyerSender(sender) {
								for i := range ClientBot {
									if ClientBot[i].MID == client.MID {
										time.Sleep(2 * time.Second)
										client.ChangeCoverPicture(to, op.Message.ID)
									}
									time.Sleep(1 * time.Second)
								}
							}
							delete(UpdateCover, client.MID)
						}
					} else if (op.Message.ContentType).String() == "VIDEO" {
						if UpdateVProfile[client.MID] {
							if IsBuyerSender(sender) {
								for i := range ClientBot {
									if ClientBot[i].MID == client.MID {
										time.Sleep(3 * time.Second)
										client.ChangeProfileVideo(to, op.Message.ID)
									}
									time.Sleep(2 * time.Second)
								}
							}
							delete(UpdateVProfile, client.MID)
						} else if UpdateVCover[client.MID] {
							if IsBuyerSender(sender) {
								for i := range ClientBot {
									if ClientBot[i].MID == client.MID {
										time.Sleep(3 * time.Second)
										client.ChangeCoverVideo(to, op.Message.ID)
									}
									time.Sleep(2 * time.Second)
								}
							}
							delete(UpdateVCover, client.MID)
						}
					} else if (op.Message.ContentType).String() == "CONTACT" {
						midd := op.Message.ContentMetadata["mid"]
						if _, cek := ContactType[sender]; cek {
							if IsMakerSender(sender) {
								if ContactType[sender] == "buyer" {
									if !InArray(Data.Buyer, midd) && !fullAccess(client, midd) {
										Data.Buyer = append(Data.Buyer, midd)
										SaveData()
										client.SendPollMention(to, "Added @! to buyers", []string{midd})
									}
								} else if ContactType[sender] == "expel" {
									if InArray(Data.Buyer, midd) {
										Data.Buyer = Remove(Data.Buyer, midd)
										SaveData()
										client.SendPollMention(to, "Delete @! from buyers", []string{midd})
									}
								}
							}
							if IsBuyerSender(sender) {
								if ContactType[sender] == "owner" {
									if !InArray(Data.Owner, midd) && !fullAccess(client, midd) {
										Data.Owner = append(Data.Owner, midd)
										SaveData()
										client.SendPollMention(to, "Added @! to owners", []string{midd})
									}
								} else if ContactType[sender] == "expel" {
									if InArray(Data.Owner, midd) {
										Data.Owner = Remove(Data.Owner, midd)
										SaveData()
										client.SendPollMention(to, "Delete @! from owners", []string{midd})
									} else if InArray(Data.Premlist, midd) {
										Data.Premlist = Remove(Data.Premlist, midd)
										SaveData()
										client.SendPollMention(to, "Delete @! from premlist", []string{midd})
									}
								}
							}
							if IsOwnerSender(sender) {
								if ContactType[sender] == "fuck" {
									if !InArray(Data.Fucklist, midd) && !fullAccess(client, midd) {
										Data.Fucklist = append(Data.Fucklist, midd)
										SaveData()
										client.SendPollMention(to, "Added @! to fucks", []string{midd})
									}
								} else if ContactType[sender] == "master" {
									if !InArray(Data.Master, midd) && !fullAccess(client, midd) {
										Data.Master = append(Data.Master, midd)
										SaveData()
										client.SendPollMention(to, "Added @! to masters", []string{midd})
									}
								} else if ContactType[sender] == "expel" {
									if InArray(Data.Fucklist, midd) {
										Data.Fucklist = Remove(Data.Fucklist, midd)
										SaveData()
										client.SendPollMention(to, "Delete @! from fucks", []string{midd})
									} else if InArray(Data.Master, midd) {
										Data.Master = Remove(Data.Master, midd)
										SaveData()
										client.SendPollMention(to, "Delete @! from masters", []string{midd})
									}
								}
							}
							if IsMasterSender(sender) {
								if ContactType[sender] == "admin" {
									if !InArray(Data.Admin, midd) && !fullAccess(client, midd) {
										Data.Admin = append(Data.Admin, midd)
										SaveData()
										client.SendPollMention(to, "Added @! to admins", []string{midd})
									}
								} else if ContactType[sender] == "wl" {
									if !InArray(Data.Whitelist, midd) && !fullAccess(client, midd) {
										Data.Whitelist = append(Data.Whitelist, midd)
										SaveData()
										client.SendPollMention(to, "Added @! to whitelist", []string{midd})
									}
								} else if ContactType[sender] == "expel" {
									if InArray(Data.Admin, midd) {
										Data.Admin = Remove(Data.Admin, midd)
										SaveData()
										client.SendPollMention(to, "Delete @! from admins", []string{midd})
									} else if InArray(Data.Whitelist, midd) {
										Data.Whitelist = Remove(Data.Whitelist, midd)
										SaveData()
										client.SendPollMention(to, "Delete @! from whitelist", []string{midd})
									}
								}
							}
							if IsAdminSender(sender) {
								if ContactType[sender] == "bot" {
									if !InArray(Data.Bot, midd) && !fullAccess(client, midd) {
										Data.Bot = append(Data.Bot, midd)
										SaveData()
										client.SendPollMention(to, "Added @! to bots", []string{midd})
									}
								} else if ContactType[sender] == "ban" {
									if !InArray(Data.Blacklist, midd) && !fullAccess(client, midd) {
										Data.Blacklist = append(Data.Blacklist, midd)
										SaveData()
										client.SendPollMention(to, "Added @! to bans", []string{midd})
									}
								} else if ContactType[sender] == "expel" {
									if !InArray(Data.Master, sender) || !InArray(Data.Admin, sender) || InArray(Data.Premlist, sender) {
										if InArray(Data.Blacklist, midd) {
											Data.Blacklist = Remove(Data.Blacklist, midd)
											SaveData()
											client.SendPollMention(to, "Delete @! from bans", []string{midd})
										}
									}
									if InArray(Data.Bot, midd) {
										Data.Bot = Remove(Data.Bot, midd)
										SaveData()
										client.SendPollMention(to, "Delete @! from bots", []string{midd})
									}
								}
							}
						}
					}
					mentions := mentions{}
					json.Unmarshal([]byte(op.Message.ContentMetadata["MENTION"]), &mentions)
					for _, mention := range mentions.MENTIONEES {
						if mention.Mid == client.MID {
							if IsBuyerSender(sender) {
								if strings.Contains(text, "upimage") {
									UpdatePicture[client.MID] = true
									client.SendMessage(to, "Send image.")
								} else if strings.Contains(text, "upcover") {
									UpdateCover[client.MID] = true
									client.SendMessage(to, "Send image.")
								} else if strings.Contains(text, "upvimage") {
									UpdateVProfile[client.MID] = true
									client.SendMessage(to, "Send video.")
								} else if strings.Contains(text, "upvicover") {
									UpdateVCover[client.MID] = true
									client.SendMessage(to, "Send video.")
								} else if strings.Contains(text, "upname: ") {
									result := strings.Split((text), ": ")
									objme := client.GetProfile()
									objme.DisplayName = result[1]
									client.UpdateProfile(objme)
									client.SendMessage(to, "Profile name updated.")
								} else if strings.Contains(text, "upstatus: ") {
									result := strings.Split((text), ": ")
									objme := client.GetProfile()
									objme.StatusMessage = result[1]
									client.UpdateProfile(objme)
									client.SendMessage(to, "Profile status updated.")
								}
							}
							if IsOwnerSender(sender) {
								if strings.Contains(text, "ginvite: ") {
									result := strings.Split((text), ": ")
									num, _ := strconv.Atoi(result[1])
									groups := client.GetGroupsJoined()
									if num > 0 && num <= len(groups) {
										if IsFriends(client, sender) == false {
											client.FindAndAddContactsByMid(sender)
										}
										client.InviteIntoChat(groups[num-1], []string{sender})
									}
									client.SendMessage(to, "success invited you.")
								} else if strings.Contains(text,"gleave: "){
									result := strings.Split((text),": ")
									num, _ := strconv.Atoi(result[1])
									groups := client.GetGroupsJoined()
									if num > 0&&num <= len(groups){
										delete(sider,groups[num-1])
										delete(Data.StayGroup,groups[num-1])
										SaveData()
										client.LeaveGroup(groups[num-1])
									}
									client.SendMessage(to, "success leave!")
								} else if strings.Contains(text, "gourl: ") {
									result := strings.Split((text), ": ")
									num, _ := strconv.Atoi(result[1])
									groups := client.GetGroupsJoined()
									if num > 0 && num <= len(groups) {
										gc := client.GetGroup(groups[num-1])
										if gc.PreventedJoinByTicket == true {
											gc.PreventedJoinByTicket = false
											client.UpdateGroup(gc)
										}
										tick := client.ReissueGroupTicket(groups[num-1])
										client.SendMessage(to, "https://line.me/R/ti/g/"+tick)
									} else { 
										client.SendMessage(to, "out of range.") 
										break 
									}
								}
							}
							if Data.AntiTag == true && !fullAccess(client,sender) && !IsGaccess(to, sender) {
								client.DeleteOtherFromChat(to, []string{sender})
								appendBl(sender)
								break
							}
						}
					}
				} else {
					client.CorrectRevision(op)
				}
			}
		}(multiFunc)
		for _, ops := range multiFunc {
			if ops.Revision != -1 {
				client.SetRevision(ops.Revision)
			} else {
				client.CorrectRevision(ops)
			}
		}
	}
}
func SingleBots(client *LINE, to string, sender string, coms string, op *core.Operation) {
	if IsOwnerSender(sender) {
		if !InArray(client.Squads, sender) {
			if remotegrupid != "" {
				fucek := to
				to = remotegrupid
				RemoteText(client,fucek)
			}
		}
	}
	text := op.Message.Text
	pesan := strings.ToLower(text)
	for _, cmd := range strings.Split(coms, ",") {
//MAKERS
		if IsMakerSender(sender) {
			if strings.HasPrefix(cmd, "addbuyer") {
				mentions := mentions{}
				targets := []string{}
				json.Unmarshal([]byte(op.Message.ContentMetadata["MENTION"]), &mentions)
				for _, mention := range mentions.MENTIONEES {
					if !InArray(Data.Buyer, string(mention.Mid)) && client.MID != mention.Mid {
						if !fullAccess(client, mention.Mid) && mention.Mid != client.MID {
							Data.Buyer = append(Data.Buyer, mention.Mid)
							targets = append(targets, mention.Mid)
						}
					}
				}
				SaveData()
				if len(targets) != 0 {
					list := "#Addbuyer:\n"
					for i := range targets {
						list += "\n" + strconv.Itoa(i+1) + ". " + "@!"
					}
					client.SendPollMention(to, list, targets)
				}
			} else if cmd == "backups" {
				if len(Data.SquadBots) > 0 {
					listsq := "#Backups:\n"
					for i := range Data.SquadBots {
						listsq += "\n" + strconv.Itoa(i+1) + ". " + "@!"
						fmt.Sprintf("... %v", i)
					}
					client.SendPollMention(to, listsq, Data.SquadBots)
				} else { client.SendMessage(to, "Backup is empty.") }
			} else if cmd == "buyer:on" {
				ContactType[sender] = "buyer"
				client.SendMessage(to, "Send contact.")
			} else if cmd == "buyers" {
				if len(Data.Buyer) > 0 {
					list := "#Buyerlist:\n"
					for i := range Data.Buyer {
						list += "\n" + strconv.Itoa(i+1) + ". " + "@!"
						fmt.Sprintf("... %v", i)
					}
					client.SendPollMention(to, list, Data.Buyer)
				} else { client.SendMessage(to, "Buyer is empty.") }
			} else if strings.HasPrefix(cmd,"buyer"+" "){
				fl := strings.Split(cmd, " ")
				typec := strings.Replace(cmd, fl[0]+" ", "", 1)
				msg, _ := client.GetRecentMessagesV2(to, 300)
				for i := 0; i < len(msg); i++ {
					item := msg[i]
					if item.ContentType == 13 {
						mid := item.ContentMetadata["mid"]
						if typec == "lcon" && !fullAccess(client, mid) {
							appendBuyer(mid)
							for i := range ClientBot{ClientBot[i].FindAndAddContactsByMid(mid)}
							client.SendMessage(to, "added to buyers.")
							break
						}
					}
				}
			}else if cmd == "checkram"{
				v, _ := mem.VirtualMemory()
				r := fmt.Sprintf("  ↳Cpu : %v core\n  ↳Ram : %v mb\n  ↳Free : %v mb\n  ↳Cache : %v mb\n  ↳UsedPercent : %f %%",cpu,bToMb(v.Used + v.Free + v.Buffers + v.Cached), bToMb(v.Free), bToMb(v.Buffers + v.Cached), v.UsedPercent)
				client.SendMessage(to, r)
			} else if cmd == "clearbuyer" {
				if len(Data.Buyer) != 0 {
					client.SendMessage(to, fmt.Sprintf("Cleared %v buyerlist", len(Data.Buyer)))
					Data.Buyer = []string{}
					SaveData()
				} else { client.SendMessage(to, "Buyer is empty.") }
			} else if strings.HasPrefix(cmd, "delbuyer") {
				mentions := mentions{}
				targets := []string{}
				json.Unmarshal([]byte(op.Message.ContentMetadata["MENTION"]), &mentions)
				for _, mention := range mentions.MENTIONEES {
					if InArray(Data.Buyer, string(mention.Mid)) {
						Data.Buyer = Remove(Data.Buyer, mention.Mid)
						targets = append(targets, mention.Mid)
					}
				}
				SaveData()
				if len(targets) != 0 {
					list := "#Delbuyer:\n"
					for i := range targets {
						list += "\n" + strconv.Itoa(i+1) + ". " + "@!"
					}
					client.SendPollMention(to, list, targets)
				}
			} else if cmd == "loginsb" {
				tok := client.CreateNewToken(to, sender, "deswin")
				Data.SelfToken = tok
				Data.SelfStatus = true
				SaveData()
				client.SendMessage(to, "Success login selfbot")
				ReloginProgram()
			} else if cmd == "help maker" {
				res := "「✿ Menu Maker ✿」"
				res += "\n"
				for a, x := range helpmaker {
					res += fmt.Sprintf("\n%v› %s %s", a+1, Data.Setkey, x)
				}
				client.SendMessage(to, res)
			} else if cmd == "appname" {
				for i := range ClientBot {
					ClientBot[i].SendMessage(to, string(ClientBot[i].AppName))
				}
			} else if cmd == "useragent" {
				for i := range ClientBot {
					ClientBot[i].SendMessage(to, string(ClientBot[i].UserAgent))
				}
			} else if cmd == "hostname" {
				for i := range ClientBot {
					ClientBot[i].SendMessage(to, string(ClientBot[i].Host))
				}
			} else if cmd == "checktoken" {
				for i := range ClientBot {
					ClientBot[i].SendMessage(to, string(ClientBot[i].AuthToken))
				}
			} else if cmd == "checkmid" {
				for i := range ClientBot {
					ClientBot[i].SendMessage(to, string(ClientBot[i].MID))
				}
			} else if strings.HasPrefix(cmd, "unbuyer ") {
				result := strings.Split((cmd), " ")
				if result[1] != "0" {
					result2, err := strconv.Atoi(result[1])
					if err != nil { client.SendMessage(to, "invalid number.")
					} else {
						for i := range Data.Buyer {
							if result2 > 0 && result2-1 < len(Data.Buyer) {
								if i == result2-1 {
									kura := Data.Buyer[i]
									Data.Buyer = Remove(Data.Buyer, kura)
									client.SendMention(to, "success delbuyer @!", []string{kura})
									SaveData()
									break
								}
							} else { 
								client.SendMessage(to, "out of range.") 
								break 
							}
						}
					}
				} else { client.SendMessage(to, "invalid range.") }
			}
		}
//BUYERS
		if IsBuyerSender(sender) {
			if strings.HasPrefix(cmd, "addowner") {
				mentions := mentions{}
				targets := []string{}
				json.Unmarshal([]byte(op.Message.ContentMetadata["MENTION"]), &mentions)
				for _, mention := range mentions.MENTIONEES {
					if !InArray(Data.Owner, string(mention.Mid)) && client.MID != mention.Mid {
						if !fullAccess(client, mention.Mid) && mention.Mid != client.MID {
							Data.Owner = append(Data.Owner, mention.Mid)
							targets = append(targets, mention.Mid)
						}
					}
				}
				SaveData()
				if len(targets) != 0 {
					list := "#Addowner:\n"
					for i := range targets {
						list += "\n" + strconv.Itoa(i+1) + ". " + "@!"
					}
					client.SendPollMention(to, list, targets)
				}
			} else if strings.HasPrefix(cmd, "addprem") {
				mentions := mentions{}
				targets := []string{}
				json.Unmarshal([]byte(op.Message.ContentMetadata["MENTION"]), &mentions)
				for _, mention := range mentions.MENTIONEES {
					if !InArray(Data.Premlist, string(mention.Mid)) && client.MID != mention.Mid {
						Data.Premlist = append(Data.Premlist, mention.Mid)
						targets = append(targets, mention.Mid)
					}
				}
				SaveData()
				if len(targets) != 0 {
					list := "#Addprem:\n"
					for i := range targets {
						list += "\n" + strconv.Itoa(i+1) + ". " + "@!"
					}
					client.SendPollMention(to, list, targets)
				}
			}else if cmd == "access"{
				allmanagers := []string{}
				nourut := 1
				listadm := " *** BUYER *** "
				if len(Data.Buyer) > 0{
					for i:= range Data.Buyer{
						allmanagers = append(allmanagers,Data.Buyer[i])
						listadm += "\n"+strconv.Itoa(i+nourut) + ". @!"
					}
					nourut = len(Data.Buyer)+1
				}
				if len(Data.Owner) > 0{
					listadm += "\n\n *** OWNER *** "
					for i:= range Data.Owner{
						allmanagers = append(allmanagers,Data.Owner[i])
						listadm += "\n"+strconv.Itoa(i+nourut) + ". @!"
					}
					nourut = len(Data.Owner)+2
				}else{listadm += "\n\n *** OWNER *** "}
				if len(Data.Master) > 0{
					listadm += "\n\n *** MASTER *** "
					for i:= range Data.Master{
						allmanagers = append(allmanagers,Data.Master[i])
						listadm += "\n"+strconv.Itoa(i+nourut) + ". @!"
					}
					nourut = len(Data.Master)+3
				}else{listadm += "\n\n *** MASTER *** "}
				if len(Data.Admin) > 0{
					listadm += "\n\n *** ADMIN *** "
					for i:= range Data.Admin{
						allmanagers = append(allmanagers,Data.Admin[i])
						listadm += "\n"+strconv.Itoa(i+nourut) + ". @!"
					}
					nourut = len(Data.Admin)+4
				}else{listadm += "\n\n *** ADMIN *** \n"}
				if len(allmanagers) != 0 {
					client.SendPollMention(to,listadm,allmanagers)
				}else{client.SendMessage(to, "Access is empty.")}
			} else if strings.HasPrefix(cmd, "delowner") {
				mentions := mentions{}
				targets := []string{}
				json.Unmarshal([]byte(op.Message.ContentMetadata["MENTION"]), &mentions)
				for _, mention := range mentions.MENTIONEES {
					if InArray(Data.Owner, string(mention.Mid)) {
						Data.Owner = Remove(Data.Owner, mention.Mid)
						targets = append(targets, mention.Mid)
					}
				}
				SaveData()
				if len(targets) != 0 {
					list := "#Delowner:\n"
					for i := range targets {
						list += "\n" + strconv.Itoa(i+1) + ". " + "@!"
					}
					client.SendPollMention(to, list, targets)
				}
			} else if strings.HasPrefix(cmd, "delprem") {
				mentions := mentions{}
				targets := []string{}
				json.Unmarshal([]byte(op.Message.ContentMetadata["MENTION"]), &mentions)
				for _, mention := range mentions.MENTIONEES {
					if InArray(Data.Premlist, string(mention.Mid)) {
						Data.Premlist = Remove(Data.Premlist, mention.Mid)
						targets = append(targets, mention.Mid)
					}
				}
				SaveData()
				if len(targets) != 0 {
					list := "#Delprem:\n"
					for i := range targets {
						list += "\n" + strconv.Itoa(i+1) + ". " + "@!"
					}
					client.SendPollMention(to, list, targets)
				}
			} else if cmd == "clearfriend" {
				for i := range ClientBot {
					targets := []string{}
					friends := ClientBot[i].GetAllContactIds()
					for _, x := range friends {
						if !fullAccess(ClientBot[i], x) {
							ClientBot[i].RemoveContact(x)
							targets = append(targets, x)
						}
					}
					ClientBot[i].SendMessage(to, "Reset: "+strconv.Itoa(len(targets))+" contacts.")
				}
			} else if cmd == "clearowner" {
				if len(Data.Owner) != 0 {
					client.SendMessage(to, fmt.Sprintf("Cleared %v ownerlist", len(Data.Owner)))
					Data.Owner = []string{}
					SaveData()
				} else { client.SendMessage(to, "Owner is empty.") }
			} else if strings.HasPrefix(cmd, "cflag ") {
				result := strings.Replace(cmd, "cflag ", "", 1)
				Data.Message.Flag = result
				SaveData()
				client.SendMessage(to, "Flag change to: "+result)
			} else if strings.HasPrefix(cmd, "cidlink ") {
				result := strings.Replace(cmd, "cidlink ", "", 1)
				IconFooter = result
				SaveData()
				client.SendMessage(to, "Id link change to: "+result)
			} else if strings.HasPrefix(cmd, "cgiflink ") {
				result := strings.Replace(cmd, "cgiflink ", "", 1)
				IconLink = result
				SaveData()
				client.SendMessage(to, "Gif link change to: "+result)
			} else if cmd == "help buyer" {
				res := "「✿ Menu Buyer ✿」"
				res += "\n"
				for a, x := range helpbuyer {
					res += fmt.Sprintf("\n%v› %s %s", a+1, Data.Setkey, x)
				}
				client.SendMessage(to, res)
			} else if strings.HasPrefix(cmd, "logmode ") {
				result := strings.Split((cmd), " ")
				r := strings.Replace(cmd, result[0]+" ", "", 1)
				if r == "on" {
					gcControl = true
					gc := client.GetGroup(to)
					gcControlV2 = gc.ID
					client.SendMessage(to, "group control actived")
				} else if r == "off" {
					gcControl = false
					gcControlV2 = ""
					client.SendMessage(to, "group control deactived")
				} else { client.SendMessage(to, "invalid command") }
			} else if cmd == "owners" {
				if len(Data.Owner) > 0 {
					list := "#Ownerlist:\n"
					for i := range Data.Owner {
						list += "\n" + strconv.Itoa(i+1) + ". " + "@!"
						fmt.Sprintf("... %v", i)
					}
					client.SendPollMention(to, list, Data.Owner)
				} else { client.SendMessage(to, "Owner is empty.") }
			} else if strings.HasPrefix(cmd,"owner"+" "){
				fl := strings.Split(cmd, " ")
				typec := strings.Replace(cmd, fl[0]+" ", "", 1)
				msg, _ := client.GetRecentMessagesV2(to, 300)
				for i := 0; i < len(msg); i++ {
					item := msg[i]
					if item.ContentType == 13 {
						mid := item.ContentMetadata["mid"]
						if typec == "lcon" && !fullAccess(client, mid) {
							appendOwner(mid)
							for i := range ClientBot{ClientBot[i].FindAndAddContactsByMid(mid)}
							client.SendMessage(to, "added to owners.")
							break
						}
					}
				}
			} else if cmd == "owner:on" {
				ContactType[sender] = "owner"
				client.SendMessage(to, "Send contact.")
			} else if strings.HasPrefix(cmd, "setsname ") {
				result := strings.Replace(cmd, "setsname ", "", 1)
				Data.Setkey = result
				SaveData()
				client.SendMessage(to, "Sname set to: "+result)
			} else if strings.HasPrefix(cmd, "setrname ") {
				result := strings.Replace(cmd, "setrname ", "", 1)
				Data.Rname = result
				SaveData()
				client.SendMessage(to, "Rname set to: "+result)
			} else if strings.HasPrefix(cmd, "msgrespon ") {
				result := strings.Replace(cmd, "msgrespon ", "", 1)
				Data.Message.Respon = result
				SaveData()
				client.SendMessage(to, "Message respon set to: "+result)
			} else if strings.HasPrefix(cmd, "msgunban ") {
				result := strings.Replace(cmd, "msgunban ", "", 1)
				Data.Message.Ban = result
				SaveData()
				client.SendMessage(to, "Message unban set to: "+result)
			} else if strings.HasPrefix(cmd, "msgbye ") {
				result := strings.Replace(cmd, "msgbye ", "", 1)
				Data.Message.Bye = result
				SaveData()
				client.SendMessage(to, "Message bye set to: "+result)
			} else if strings.HasPrefix(cmd, "msgsider ") {
				result := strings.Replace(cmd, "msgsider ", "", 1)
				Data.Message.Sider = result
				SaveData()
				client.SendMessage(to, "Message sider set to: "+result)
			} else if strings.HasPrefix(cmd, "msgwelcome ") {
				result := strings.Replace(cmd, "msgwelcome ", "", 1)
				Data.Message.Welcome = result
				SaveData()
				client.SendMessage(to, "Message sider set to: "+result)
			} else if strings.HasPrefix(cmd, "msgfresh ") {
				result := strings.Replace(cmd, "msgfresh ", "", 1)
				Data.Message.Fresh = result
				SaveData()
				client.SendMessage(to, "Message fresh set to: "+result)
			} else if strings.HasPrefix(cmd, "msglimit ") {
				result := strings.Replace(cmd, "msglimit ", "", 1)
				Data.Message.Limit = result
				SaveData()
				client.SendMessage(to, "Message limit set to: "+result)
			}else if strings.HasPrefix(cmd, "setkick ") {
				anjay := strings.Split((cmd), " ")
				num, err := strconv.Atoi(anjay[1])
				if err != nil {
					client.SendMessage(to, "Please use number!")
				} else {
					Kickbatas = num
					client.SendMessage(to, "Limiter kick set to "+anjay[1])
				}
			}else if strings.HasPrefix(cmd, "setcancel ") {
				anjay := strings.Split((cmd), " ")
				num, err := strconv.Atoi(anjay[1])
				if err != nil {
					client.SendMessage(to, "Please use number!")
				} else {
					Cansbatas = num
					client.SendMessage(to, "Limiter cancel set to "+anjay[1])
				}
			}else if strings.HasPrefix(cmd, "setqr ") {
				anjay := strings.Split((cmd), " ")
				num, err := strconv.Atoi(anjay[1])
				if err != nil {
					client.SendMessage(to, "Please use number!")
				} else {
					Closeqrbatas = num
					client.SendMessage(to, "Limiter spamqr set to "+anjay[1])
				}
			}else if strings.HasPrefix(cmd,"setlimiter "){
				result := strings.Split((cmd)," ")
				no,err := strconv.Atoi(result[1])
				if err != nil {
					client.SendMessage(to, "Please use number!")
				}else{
					Kickbatas = no
					Cansbatas = no
					client.SendMessage(to, "Limiter successs set to "+result[1])
				}
			}else if strings.HasPrefix(cmd, "delaykick ") {
				anjay := strings.Split((cmd), " ")
				num, err := strconv.Atoi(anjay[1])
				if err != nil {
					client.SendMessage(to, "Please use number!")
				} else {
					KickDelay = num
					client.SendMessage(to, "Delay kick set to "+anjay[1]+"ms")
				}
			}else if strings.HasPrefix(cmd, "delayinvite ") {
				anjay := strings.Split((cmd), " ")
				num, err := strconv.Atoi(anjay[1])
				if err != nil {
					client.SendMessage(to, "Please use number!")
				} else {
					InviteDelay = num
					client.SendMessage(to, "Delay invite set to "+anjay[1]+"ms")
				}
			}else if strings.HasPrefix(cmd, "delayjoin ") {
				anjay := strings.Split((cmd), " ")
				num, err := strconv.Atoi(anjay[1])
				if err != nil {
					client.SendMessage(to, "Please use number!")
				} else {
					JoinDelay = num
					client.SendMessage(to, "Delay join set to "+anjay[1]+"ms")
				} 
			}else if cmd == "delays"{
				delays := fmt.Sprintf("Command\n  Delaykick <num>\n  Delayinvite <num>\n  Delayjoin <num>\n\nDelays bots\n  Kick: %vms\n  Invite: %vms\n  Join: %vms", KickDelay, InviteDelay, JoinDelay)
				client.SendMessage(to, delays)
			} else if cmd == "inv" {
				poolInviteStay(client, to)
			} else if cmd == "join" {
				poolJoinStay(client, to)
			} else if cmd == "kills" {
				if len(KillMode) > 0 {
					listsq := "#Listkill:\n"
					var no = 0
					for i := range KillMode {
						if i != "" {
							con := client.GetContact(i)
							listsq += fmt.Sprintf("\n%v. %s", no+1, con.DisplayName)
							for _, a := range KillMode[i] {
								conn := client.GetContact(a)
								listsq += fmt.Sprintf("\n    - %s", conn.DisplayName)
							}
							no++
						}
					}
					client.SendMessage(to, listsq)
				}
			}else if cmd == "limiters"{
				limiter := fmt.Sprintf("Command\n  Setlimiter <num>\n  Setspamqr <num>\n  Setkick <num>\n  Setcancel <num>\n\nLimiter bots\n  Spamqr: %v\n  Kick: %v\n  Cancel: %v", Closeqrbatas, Kickbatas, Cansbatas)
				client.SendMessage(to, limiter)
			} else if cmd == "premlist" {
				if len(Data.Premlist) > 0 {
					list := "#Premlist:\n"
					for i := range Data.Premlist {
						list += "\n" + strconv.Itoa(i+1) + ". " + "@!"
						fmt.Sprintf("... %v", i)
					}
					client.SendPollMention(to, list, Data.Premlist)
				} else { client.SendMessage(to, "Premlist is empty.") }
			} else if cmd == "reboot" {
				client.SendMessage(to, "Waiting Restarted.")
				ReloginProgram()
			} else if strings.HasPrefix(cmd, "unfriend ") {
				result := strings.Split((cmd), " ")
				if result[1] != "0" {
					result2, err := strconv.Atoi(result[1])
					if err != nil { client.SendMessage(to, "invalid number.")
					} else {
						friends := client.GetAllContactIds()
						for i := range friends {
							if result2 > 0 && result2-1 < len(friends) {
								if i == result2-1 {
									kura := friends[i]
									friends = Remove(friends, kura)
									client.SendMention(to, "success delfriend @!", []string{kura})
									SaveData()
									break
								}
							} else { 
								client.SendMessage(to, "out of range.") 
								break 
							}
						}
					}
				} else { client.SendMessage(to, "invalid range.") }
			} else if strings.HasPrefix(cmd, "unowner ") {
				result := strings.Split((cmd), " ")
				if result[1] != "0" {
					result2, err := strconv.Atoi(result[1])
					if err != nil { client.SendMessage(to, "invalid number.")
					} else {
						for i := range Data.Owner {
							if result2 > 0 && result2-1 < len(Data.Owner) {
								if i == result2-1 {
									kura := Data.Owner[i]
									Data.Owner = Remove(Data.Owner, kura)
									client.SendMention(to, "success delowner @!", []string{kura})
									SaveData()
									break
								}
							} else { 
								client.SendMessage(to, "out of range.") 
								break 
							}
						}
					}
				} else { client.SendMessage(to, "invalid range.") }
			}
		}
//OWNERS
		if IsOwnerSender(sender) {
			if strings.HasPrefix(cmd, "addfuck") {
				mentions := mentions{}
				targets := []string{}
				json.Unmarshal([]byte(op.Message.ContentMetadata["MENTION"]), &mentions)
				for _, mention := range mentions.MENTIONEES {
					if !InArray(Data.Fucklist, string(mention.Mid)) && client.MID != mention.Mid {
						if !fullAccess(client, mention.Mid) && mention.Mid != client.MID {
							Data.Fucklist = append(Data.Fucklist, mention.Mid)
							targets = append(targets, mention.Mid)
						}
					}
				}
				SaveData()
				if len(targets) != 0 {
					list := "#Addfuck:\n"
					for i := range targets {
						list += "\n" + strconv.Itoa(i+1) + ". " + "@!"
					}
					client.SendPollMention(to, list, targets)
				}
			} else if strings.HasPrefix(cmd,"addajs"){
				mentions := mentions{}
				json.Unmarshal([]byte(op.Message.ContentMetadata["MENTION"]), &mentions)
				for _, mention := range mentions.MENTIONEES{
					if InArray(client.Squads, mention.Mid) && !myAccess(mention.Mid){
						if !InArray(Data.StayAjs[to], mention.Mid){
							if IsAjs(to, mention.Mid) == false {
								Data.StayAjs[to] = append(Data.StayAjs[to], mention.Mid)
							}
							time.Sleep(time.Second * 2)
							if len(Data.StayAjs[to]) != 0 {
								for i := range ClientBot {
									if IsAjs(to, ClientBot[i].MID) == true {
										ClientBot[i].LeaveGroup(to)
									}
								}
								client.SendMessage(to,"Added antijs success")
							}
						}else{client.SendMessage(to,"Already on antijs")}
					}
				}
				SaveData()
			} else if strings.HasPrefix(cmd, "delajs") {
				mentions := mentions{}
				targets := []string{}
				json.Unmarshal([]byte(op.Message.ContentMetadata["MENTION"]), &mentions)
				for _, mention := range mentions.MENTIONEES {
					if InArray(Data.StayAjs[to], string(mention.Mid)) {
						Data.StayAjs[to] = Remove(Data.StayAjs[to], mention.Mid)
						targets = append(targets, mention.Mid)
					}
				}
				SaveData()
				if len(targets) != 0 {
					list := "#Delajs:\n"
					for i := range targets {
						list += "\n" + strconv.Itoa(i+1) + ". " + "@!"
					}
					client.SendPollMention(to, list, targets)
				}
			} else if cmd == "ajs:on"{
				for i := range ClientBot {
					if IsAjs(to, ClientBot[i].MID) == true {
						ClientBot[i].LeaveGroup(to)
					} else {
						InvitedAjs(client, to)
					}
				}
			} else if cmd == "clearajs"{
				if len(Data.StayAjs[to]) == 0{
					client.SendMessage(to, "dont have ajslist")
				}else{
					client.SendMessage(to, "Cleared "+strconv.Itoa(len(Data.StayAjs[to]))+" antijs.")
					tick := client.ReissueGroupTicket(to)
					gc, _ := client.GetGroupV2(to)
					_, found := Data.StayAjs[to]
					if found == true {
						if gc.PreventedJoinByTicket == true {
							gc.PreventedJoinByTicket = false
							client.UpdateGroup(gc)
						}
						for i := range ClientBot {
							if InArray(Data.StayAjs[to], ClientBot[i].MID) {
								if ClientMid[ClientBot[i].MID].Limited == false {
									ClientBot[i].AcceptChatInvitationByTicket(to, tick)
								}
							}
						}
						if gc.PreventedJoinByTicket == false {
							gc.PreventedJoinByTicket = true
							time.Sleep(500 * time.Millisecond)
							client.UpdateGroup(gc)
						}
					}
					delete(Data.StayAjs,to)
					SaveData()
				}
			} else if cmd == "ajslist" {
				if len(Data.StayAjs[to]) > 0 {
					list := "#Ajslist:\n"
					for i := range Data.StayAjs[to] {
						list += "\n" + strconv.Itoa(i+1) + ". " + "@!"
						fmt.Sprintf("... %v", i)
					}
					client.SendPollMention(to, list, Data.StayAjs[to])
				} else { client.SendMessage(to, "Antijs is empty.") }
			} else if strings.HasPrefix(cmd, "addmaster") {
				mentions := mentions{}
				targets := []string{}
				json.Unmarshal([]byte(op.Message.ContentMetadata["MENTION"]), &mentions)
				for _, mention := range mentions.MENTIONEES {
					if !InArray(Data.Master, string(mention.Mid)) && client.MID != mention.Mid {
						if !fullAccess(client, mention.Mid) && mention.Mid != client.MID {
							Data.Master = append(Data.Master, mention.Mid)
							targets = append(targets, mention.Mid)
						}
					}
				}
				SaveData()
				if len(targets) != 0 {
					list := "#Addmaster:\n"
					for i := range targets {
						list += "\n" + strconv.Itoa(i+1) + ". " + "@!"
					}
					client.SendPollMention(to, list, targets)
				}
			} else if strings.HasPrefix(cmd, "addgowner") {
				mentions := mentions{}
				targets := []string{}
				json.Unmarshal([]byte(op.Message.ContentMetadata["MENTION"]), &mentions)
				for _, mention := range mentions.MENTIONEES {
					if !InArray(Data.GroupOwn[to], string(mention.Mid)) && client.MID != mention.Mid {
						if !fullAccess(client, mention.Mid) && mention.Mid != client.MID {
							Data.GroupOwn[to] = append(Data.GroupOwn[to], mention.Mid)
							targets = append(targets, mention.Mid)
							if IsFriends(client, mention.Mid) == false {
								client.FindAndAddContactsByMid(mention.Mid)
							}
						}
					}
				}
				SaveData()
				if len(targets) != 0 {
					list := "#Addgowner:\n"
					for i := range targets {
						list += "\n" + strconv.Itoa(i+1) + ". " + "@!"
					}
					client.SendPollMention(to, list, targets)
				}
			}else if cmd == "warmode"{
				Data.FastMode = false
				Data.QrMode = false
				Data.MixMode = false
				Data.VictimMode = false
				Data.KillMode = false
				SaveData()
				client.SendMessage(to, "war is enabled.")
			}else if cmd == "fastmode"{
				Data.QrMode = false
				Data.MixMode = false
				Data.VictimMode = false
				Data.KillMode = false
				Data.FastMode = true
				SaveData()
				client.SendMessage(to, "fast is enabled.")
			}else if cmd == "singlemode"{
				Data.FastMode = false
				Data.QrMode = false
				Data.MixMode = false
				Data.KillMode = false
				Data.VictimMode = true
				SaveData()
				client.SendMessage(to, "single is enabled.")
			}else if cmd == "qrmode"{
				Data.FastMode = false
				Data.MixMode = false
				Data.VictimMode = false
				Data.KillMode = false
				Data.QrMode = true
				SaveData()
				client.SendMessage(to, "qr is enabled.")
			}else if cmd == "mixmode"{
				Data.FastMode = false
				Data.QrMode = false
				Data.VictimMode = false
				Data.KillMode = false
				Data.MixMode = true
				SaveData()
				client.SendMessage(to, "mix is enabled.")
			} else if cmd == "cleargowner" {
				if len(Data.GroupOwn[to]) != 0 {
					client.SendMessage(to, fmt.Sprintf("Cleared %v gownerlist", len(Data.GroupOwn[to])))
					delete(Data.GroupOwn, to)
					SaveData()
				} else { client.SendMessage(to, "Gowner is empty.") }
			} else if cmd == "clearmaster" {
				if len(Data.Master) != 0 {
					client.SendMessage(to, fmt.Sprintf("Cleared %v masterlist", len(Data.Master)))
					Data.Master = []string{}
					SaveData()
				} else { client.SendMessage(to, "Master is empty.") }
			} else if cmd == "clearfuck" {
				if len(Data.Fucklist) != 0 {
					client.SendMessage(to, fmt.Sprintf("Cleared %v fucklist", len(Data.Fucklist)))
					Data.Fucklist = []string{}
					SaveData()
				} else { client.SendMessage(to, "Fucklist is empty.") }
			} else if cmd == "clear allprotect" {
				Data.ProKick = []string{}
				Data.ProQr = []string{}
				Data.ProInvite = []string{}
				Data.ProCancel = []string{}
				SaveData()
				client.SendMessage(to, "Cleared allprotected.")
			} else if strings.HasPrefix(cmd, "delfuck") {
				mentions := mentions{}
				targets := []string{}
				json.Unmarshal([]byte(op.Message.ContentMetadata["MENTION"]), &mentions)
				for _, mention := range mentions.MENTIONEES {
					if InArray(Data.Fucklist, string(mention.Mid)) {
						Data.Fucklist = Remove(Data.Fucklist, mention.Mid)
						targets = append(targets, mention.Mid)
					}
				}
				SaveData()
				if len(targets) != 0 {
					list := "#Delfuck:\n"
					for i := range targets {
						list += "\n" + strconv.Itoa(i+1) + ". " + "@!"
					}
					client.SendPollMention(to, list, targets)
				}
			} else if strings.HasPrefix(cmd, "delmaster") {
				mentions := mentions{}
				targets := []string{}
				json.Unmarshal([]byte(op.Message.ContentMetadata["MENTION"]), &mentions)
				for _, mention := range mentions.MENTIONEES {
					if InArray(Data.Master, string(mention.Mid)) {
						Data.Master = Remove(Data.Master, mention.Mid)
						targets = append(targets, mention.Mid)
					}
				}
				SaveData()
				if len(targets) != 0 {
					list := "#Delmaster:\n"
					for i := range targets {
						list += "\n" + strconv.Itoa(i+1) + ". " + "@!"
					}
					client.SendPollMention(to, list, targets)
				}
			} else if strings.HasPrefix(cmd, "delgowner") {
				mentions := mentions{}
				targets := []string{}
				json.Unmarshal([]byte(op.Message.ContentMetadata["MENTION"]), &mentions)
				for _, mention := range mentions.MENTIONEES {
					if InArray(Data.GroupOwn[to], string(mention.Mid)) {
						Data.GroupOwn[to] = Remove(Data.GroupOwn[to], mention.Mid)
						targets = append(targets, mention.Mid)
					}
				}
				SaveData()
				if len(targets) != 0 {
					list := "#Delgowner:\n"
					for i := range targets {
						list += "\n" + strconv.Itoa(i+1) + ". " + "@!"
					}
					client.SendPollMention(to, list, targets)
				}
			} else if strings.HasPrefix(cmd, "forceinvite ") {
				spl := strings.Replace(cmd, "forceinvite ", "", 1)
				if spl == "on" {
					Data.ForceInvite = true
					SaveData()
					client.SendMessage(to, "forceinvite enabled.")
				} else if spl == "off" {
					Data.ForceInvite = false
					SaveData()
					client.SendMessage(to, "forceinvite disabled.")
				}
			} else if strings.HasPrefix(cmd, "forceqr ") {
				spl := strings.Replace(cmd, "forceqr ", "", 1)
				if spl == "on" {
					Data.ForceJoinqr = true
					SaveData()
					client.SendMessage(to, "forceinvite enabled.")
				} else if spl == "off" {
					Data.ForceJoinqr = false
					SaveData()
					client.SendMessage(to, "forceinvite disabled.")
				}
			} else if cmd == "fucks" {
				if len(Data.Fucklist) > 0 {
					list := "#Fucklist:\n"
					for i := range Data.Fucklist {
						list += "\n" + strconv.Itoa(i+1) + ". " + "@!"
						fmt.Sprintf("... %v", i)
					}
					client.SendPollMention(to, list, Data.Fucklist)
				} else { client.SendMessage(to, "Fucklist is empty.") }
			} else if cmd == "fuck:on" {
				ContactType[sender] = "fuck"
				client.SendMessage(to, "Send contact.")
			} else if cmd == "gowners" {
				if len(Data.GroupOwn[to]) > 0 {
					list := "#Gownerlist:\n"
					for i := range Data.GroupOwn[to] {
						list += "\n" + strconv.Itoa(i+1) + ". " + "@!"
						fmt.Sprintf("... %v", i)
					}
					client.SendPollMention(to, list, Data.GroupOwn[to])
				} else { client.SendMessage(to, "Gowner is empty.") }
			} else if strings.HasPrefix(cmd, "groupcast ") { //groupcast
				result := strings.Split((cmd), " ")
				r := strings.Replace(cmd, result[0]+" ", "", 1)
				groups := client.GetGroupsJoined()
				for i := range groups {
					client.SendMessage(groups[i], r)
				}
				client.SendMessage(to, "Success broadcast to "+strconv.Itoa(len(groups))+" group")
			} else if cmd == "help owner" {
				res := "「✿ Menu Owner ✿」"
				res += "\n"
				for a, x := range helpowner {
					res += fmt.Sprintf("\n%v› %s %s", a+1, Data.Setkey, x)
				}
				client.SendMessage(to, res)
			} else if strings.HasPrefix(cmd, "ginvite ") {
				result := strings.Split((cmd), " ")
				num, _ := strconv.Atoi(result[1])
				groups := client.GetGroupsJoined()
				if num > 0 && num <= len(groups) {
					if IsFriends(client, sender) == false {
						client.FindAndAddContactsByMid(sender)
					}
					client.InviteIntoChat(groups[num-1], []string{sender})
				}
				client.SendMessage(to, "success invited you.")
			}else if strings.HasPrefix(cmd,"gleave "){
				result := strings.Split((cmd)," ")
				num, _ := strconv.Atoi(result[1])
				groups := client.GetGroupsJoined()
				if num > 0&&num <= len(groups){
					delete(sider,groups[num-1])
					delete(Data.StayGroup,groups[num-1])
					SaveData()
					client.LeaveGroup(groups[num-1])
				}
				client.SendMessage(to, "success leave!")
			} else if strings.HasPrefix(cmd, "gourl ") {
				result := strings.Split((cmd), " ")
				groups := client.GetGroupsJoined()
				num, _ := strconv.Atoi(result[1])
				if num > 0 && num <= len(groups) {
					gc, _ := client.GetGroupV2(groups[num-1])
					if gc.PreventedJoinByTicket == true {
						gc.PreventedJoinByTicket = false
						client.UpdateGroup(gc)
					}
					tick := client.ReissueGroupTicket(groups[num-1])
					client.SendMessage(to, "https://line.me/R/ti/g/"+tick)
				} else { 
								client.SendMessage(to, "out of range.") 
								break 
							}
			} else if strings.HasPrefix(cmd, "gnuke ") {
				runtime.GOMAXPROCS(cpu)
				result := strings.Split((text), " ")
				groups := client.GetGroupsJoined()
				num, _ := strconv.Atoi(result[1])
				if num > 0 && num <= len(groups) {
					gc, _ := client.GetGroupV2(groups[num-1])
					if gc.PreventedJoinByTicket == false {
						gc.PreventedJoinByTicket = true
						client.UpdateGroup(gc)
					}
					target := gc.MemberMids
					alltargets := []string{}
					for i := range target {
						if !fullAccess(client, target[i]) {
							alltargets = append(alltargets, target[i])
						}
					}
					tl := len(alltargets)
					var wg sync.WaitGroup
					wg.Add(tl)
					for i := 0; i < tl; i++ {
						go func(i int) {
							defer wg.Done()
							val := []string{alltargets[i]}
							client.DeleteOtherFromChat(to, val)
						}(i)
					}
					wg.Wait()
				}
				client.SendMessage(to, "Success nuke!!")
			} else if cmd == "groups" {
				groups := client.GetGroupsJoined()
				if len(groups) != 0 {
					result := "#Group joined:\n"
					for i := range groups {
						gc := client.GetGroup(groups[i])
						result += "\n" + strconv.Itoa(i+1) + ". " + gc.Name + " " + strconv.Itoa(len(gc.Members)) + "/" + strconv.Itoa(len(gc.Invitee))
					}
					client.SendMessage(to, result)
				} else { client.SendMessage(to, "Group is empty.") }
			} else if strings.HasPrefix(cmd, "killmode ") {
				spl := strings.Replace(cmd, "killmode ", "", 1)
				if spl == "on" {
					Data.FastMode = false
					Data.QrMode = false
					Data.MixMode = false
					Data.VictimMode = false
					Data.KillMode = true
					SaveData()
					client.SendMessage(to, "Killmode is enabled.")
				} else if spl == "off" {
					Data.KillMode = false
					SaveData()
					client.SendMessage(to, "Killmode is disabled.")
				}
			} else if strings.HasPrefix(cmd, "nukejoin ") {
				spl := strings.Replace(cmd, "nukejoin ", "", 1)
				if spl == "on" {
					Data.NukeJoin = true
					SaveData()
					client.SendMessage(to, "Nukejoin is enabled.")
				} else if spl == "off" {
					Data.NukeJoin = false
					SaveData()
					client.SendMessage(to, "Nukejoin is disabled.")
				}
			} else if strings.HasPrefix(cmd, "identict ") {
				spl := strings.Replace(cmd, "identict ", "", 1)
				if spl == "on" {
					Data.Identict = true
					SaveData()
					client.SendMessage(to, "Identict is enabled.")
				} else if spl == "off" {
					Data.Identict = false
					SaveData()
					client.SendMessage(to, "Identict is disabled.")
				}
			} else if strings.HasPrefix(cmd, "autopro ") {
				spl := strings.Replace(cmd, "autopro ", "", 1)
				if spl == "on" {
					Data.AutoPro = true
					SaveData()
					client.SendMessage(to, "Autopro is enabled.")
				} else if spl == "off" {
					Data.AutoPro = false
					SaveData()
					client.SendMessage(to, "Autopro is disabled.")
				}
			} else if strings.HasPrefix(cmd, "autopurge ") {
				spl := strings.Replace(cmd, "autopurge ", "", 1)
				if spl == "on" {
					Data.AutoPurge = true
					SaveData()
					client.SendMessage(to, "Autopurge is enabled.")
				} else if spl == "off" {
					Data.AutoPurge = false
					SaveData()
					client.SendMessage(to, "Autopurge is disabled.")
				}
			} else if cmd == "pendings" {
				groups := client.GetGroupsInvited()
				if len(groups) != 0 {
					result := "#Group invited:\n"
					for i := range groups {
						gc := client.GetGroup(groups[i])
						result += "\n" + strconv.Itoa(i+1) + ". " + gc.Name
					}
					client.SendMessage(to, result)
				} else { client.SendMessage(to, "Pending is empty.") }
			} else if cmd == "purgeall" {
				for x := range ClientBot {
					groups := ClientBot[x].GetGroupsJoined()
					if len(groups) != 0 {
						for i := range groups {
							ClientBot[x].NormalKickoutFromGroup(to, []string{"FuckYou"})
							if ClientMid[ClientBot[x].MID].Limited == false {
								gc := ClientBot[x].GetGroup(groups[i])
								memlist := gc.Members
								for _, v := range memlist {
									asw := v.Mid
									if IsBlacklist(asw) == true {
										go func(asw string){ClientBot[x].DeleteOtherFromChat(groups[i],[]string{asw})}(asw)
									}
								}
							}
						}
					}
				}
				client.SendMessage(to, "Success purgeall shitlist.")
			} else if cmd == "masters" {
				if len(Data.Master) > 0 {
					list := "#Masterlist:\n"
					for i := range Data.Master {
						list += "\n" + strconv.Itoa(i+1) + ". " + "@!"
						fmt.Sprintf("... %v", i)
					}
					client.SendPollMention(to, list, Data.Master)
				} else { client.SendMessage(to, "Master is empty.") }
			} else if strings.HasPrefix(cmd,"master"+" "){
				fl := strings.Split(cmd, " ")
				typec := strings.Replace(cmd, fl[0]+" ", "", 1)
				msg, _ := client.GetRecentMessagesV2(to, 300)
				for i := 0; i < len(msg); i++ {
					item := msg[i]
					if item.ContentType == 13 {
						mid := item.ContentMetadata["mid"]
						if typec == "lcon" && !fullAccess(client, mid) {
							appendMaster(mid)
							for i := range ClientBot{ClientBot[i].FindAndAddContactsByMid(mid)}
							client.SendMessage(to, "added to masters.")
							break
						}
					}
				}
			} else if cmd == "master:on" {
				ContactType[sender] = "master"
				client.SendMessage(to, "Send contact.")
			} else if strings.HasPrefix(cmd, "remote ") {
				result := strings.Split((cmd), " ")
				if result[1] != "0" {
					result2, err := strconv.Atoi(result[1])
					if err != nil { client.SendMessage(to, "invalid number.")
					} else {
						gr := client.GetGroupsJoined()
						for i := range gr {
							if result2 > 0 && result2-1 < len(gr) {
								if i == result2-1 {
									gid := gr[i]
									remotegrupid = string(gid)
									client.SendMessage(to, "Your command?")
									break
								}
							} else { 
								client.SendMessage(to, "out of range.") 
								break 
							}
						}
					}
				} else { client.SendMessage(to, "invalid range.") }
			} else if strings.HasPrefix(cmd,"setcom ") {
				ditha := strings.ReplaceAll(cmd, "setcom ", "")
				cmdLil := strings.Split(ditha, " ")
				if cmdLil[0] == "banlist" {
					Data.Command.Banlist = cmdLil[1]
				}
				if cmdLil[0] == "clearban" {
					Data.Command.Clearban = cmdLil[1]
				}
				if cmdLil[0] == "count" {
					Data.Command.Count = cmdLil[1]
				}
				if cmdLil[0] == "kick" {
					Data.Command.Kick = cmdLil[1]
				}
				if cmdLil[0] == "leave" {
					Data.Command.Leave = cmdLil[1]
				}
				if cmdLil[0] == "outall" {
					Data.Command.Outall = cmdLil[1]
				}
				if cmdLil[0] == "out" {
					Data.Command.Out = cmdLil[1]
				}
				if cmdLil[0] == "respon" {
					Data.Command.Respon = cmdLil[1]
				}
				if cmdLil[0] == "setting" {
					Data.Command.Setting = cmdLil[1]
				}
				if cmdLil[0] == "set" {
					Data.Command.Set = cmdLil[1]
				}
				if cmdLil[0] == "speed" {
					Data.Command.Speed = cmdLil[1]
				}
				if cmdLil[0] == "status" {
					Data.Command.Status = cmdLil[1]
				}
				if cmdLil[0] == "unsend" {
					Data.Command.Unsend = cmdLil[1]
				}
				SaveData()
				kowe := cmdLil[0]
				jancuk := cmdLil[1]
				client.SendMessage(to, "Changed cmd: "+kowe+" to "+jancuk)
			} else if cmd == "cmdlist"{
				rst := "=== 𝗟𝗶𝘀𝘁 𝗖𝗼𝗺𝗺𝗮𝗻𝗱 ===\n"
				rst += "\n   Cmd Banlist: "+Data.Command.Banlist+""
				rst += "\n   Cmd Clearban: "+Data.Command.Clearban+""
				rst += "\n   Cmd Count: "+Data.Command.Count+""
				rst += "\n   Cmd Kick: "+Data.Command.Kick+""
				rst += "\n   Cmd Leave: "+Data.Command.Leave+""
				rst += "\n   Cmd Outall: "+Data.Command.Outall+""
				rst += "\n   Cmd Out: " +Data.Command.Out+""
				rst += "\n   Cmd Respon: "+Data.Command.Respon+""
				rst += "\n   Cmd Setting: "+Data.Command.Setting+""
				rst += "\n   Cmd Set: "+Data.Command.Set+""
				rst += "\n   Cmd Speed: "+Data.Command.Speed+""
				rst += "\n   Cmd Status: "+Data.Command.Status+""
				rst += "\n   Cmd Unsend: "+Data.Command.Unsend+""
				client.SendMessage(to, rst)
			}else if strings.HasPrefix(cmd,"tagall:"){
				result := strings.Split((cmd),":")
				num, _ := strconv.Atoi(result[1])
				groups := client.GetGroupsJoined()
				if num > 0&&num <= len(groups){
					gc,_ := client.GetGroupV2(groups[num-1])
					target := gc.MemberMids
					targets:= []string{}
					for i:= range target{
						targets = append(targets,target[i])
					}
					client.SendPollMention(to,"Mentions member:\n",targets)
				}else{client.SendMessage(to, "out of range")}
			}else if strings.HasPrefix(cmd,"tagpen:"){
				result := strings.Split((cmd),":")
				num, _ := strconv.Atoi(result[1])
				groups := client.GetGroupsJoined()
				if num > 0&&num <= len(groups){
					gc,_ := client.GetGroupV2(groups[num-1])
					target := gc.InviteeMids
					targets:= []string{}
					if len(target) != 0 {
						for i:= range target{
							targets = append(targets,target[i])
						}
						client.SendPollMention(to,"Mentions Pending:\n",targets)
					}else{client.SendMessage(to, "Pending is empty.")}
				}else{client.SendMessage(to, "out of range")}
			} else if strings.HasPrefix(cmd, "unfuck ") {
				result := strings.Split((cmd), " ")
				if result[1] != "0" {
					result2, err := strconv.Atoi(result[1])
					if err != nil { client.SendMessage(to, "invalid number.")
					} else {
						for i := range Data.Fucklist {
							if result2 > 0 && result2-1 < len(Data.Fucklist) {
								if i == result2-1 {
									kura := Data.Fucklist[i]
									Data.Fucklist = Remove(Data.Fucklist, kura)
									client.SendMention(to, "success delfuck @!", []string{kura})
									SaveData()
									break
								}
							} else { 
								client.SendMessage(to, "out of range.") 
								break 
							}
						}
					}
				} else { client.SendMessage(to, "invalid range.") }
			} else if strings.HasPrefix(cmd, "ungowner ") {
				result := strings.Split((cmd), " ")
				if result[1] != "0" {
					result2, err := strconv.Atoi(result[1])
					if err != nil { client.SendMessage(to, "invalid number.")
					} else {
						for i := range Data.GroupOwn[to] {
							if result2 > 0 && result2-1 < len(Data.GroupOwn[to]) {
								if i == result2-1 {
									kura := Data.GroupOwn[to][i]
									Data.GroupOwn[to] = Remove(Data.GroupOwn[to], kura)
									client.SendMention(to, "success delgowner @!", []string{kura})
									SaveData()
									break
								}
							} else { 
								client.SendMessage(to, "out of range.") 
								break 
							}
						}
					}
				} else { client.SendMessage(to, "invalid range.") }
			} else if strings.HasPrefix(cmd, "unmaster ") {
				result := strings.Split((cmd), " ")
				if result[1] != "0" {
					result2, err := strconv.Atoi(result[1])
					if err != nil { client.SendMessage(to, "invalid number.")
					} else {
						for i := range Data.Master {
							if result2 > 0 && result2-1 < len(Data.Master) {
								if i == result2-1 {
									kura := Data.Master[i]
									Data.Master = Remove(Data.Master, kura)
									client.SendMention(to, "success delmaster @!", []string{kura})
									SaveData()
									break
								}
							} else { 
								client.SendMessage(to, "out of range.") 
								break 
							}
						}
					}
				} else { client.SendMessage(to, "invalid range.") }
			}
		}
//MASTERS
		if IsMasterSender(sender) || InArray(Data.GroupOwn[to], sender) {
			if cmd == "admin:on" {
				if !InArray(Data.GroupOwn[to], sender) {
				ContactType[sender] = "admin"
				client.SendMessage(to, "Send contact")
				}
			} else if cmd == "wl:on" {
				if !InArray(Data.GroupOwn[to], sender) {
				ContactType[sender] = "wl"
				client.SendMessage(to, "Send contact")
				}
			} else if strings.HasPrefix(cmd, "addadmin") {
				if !InArray(Data.GroupOwn[to], sender) {
				mentions := mentions{}
				targets := []string{}
				json.Unmarshal([]byte(op.Message.ContentMetadata["MENTION"]), &mentions)
				for _, mention := range mentions.MENTIONEES {
					if !InArray(Data.Admin, string(mention.Mid)) && client.MID != mention.Mid {
						if !fullAccess(client, mention.Mid) && mention.Mid != client.MID {
							Data.Admin = append(Data.Admin, mention.Mid)
							targets = append(targets, mention.Mid)
						}
					}
				}
				SaveData()
				if len(targets) != 0 {
					list := "#Addadmin:\n"
					for i := range targets {
						list += "\n" + strconv.Itoa(i+1) + ". " + "@!"
					}
					client.SendPollMention(to, list, targets)
				}
				}
			} else if strings.HasPrefix(cmd, "addgadmin") {
				mentions := mentions{}
				targets := []string{}
				json.Unmarshal([]byte(op.Message.ContentMetadata["MENTION"]), &mentions)
				for _, mention := range mentions.MENTIONEES {
					if !InArray(Data.GroupAdm[to], string(mention.Mid)) && client.MID != mention.Mid {
						if !fullAccess(client, mention.Mid) && mention.Mid != client.MID {
							Data.GroupAdm[to] = append(Data.GroupAdm[to], mention.Mid)
							targets = append(targets, mention.Mid)
							if IsFriends(client, mention.Mid) == false {
								client.FindAndAddContactsByMid(mention.Mid)
							}
						}
					}
				}
				SaveData()
				if len(targets) != 0 {
					list := "#Addgadmin:\n"
					for i := range targets {
						list += "\n" + strconv.Itoa(i+1) + ". " + "@!"
					}
					client.SendPollMention(to, list, targets)
				}
			} else if strings.HasPrefix(cmd, "addwl") {
				if !InArray(Data.GroupOwn[to], sender) {
				mentions := mentions{}
				targets := []string{}
				json.Unmarshal([]byte(op.Message.ContentMetadata["MENTION"]), &mentions)
				for _, mention := range mentions.MENTIONEES {
					if !InArray(Data.Whitelist, string(mention.Mid)) && client.MID != mention.Mid {
						if !fullAccess(client, mention.Mid) && mention.Mid != client.MID {
							Data.Whitelist = append(Data.Whitelist, mention.Mid)
							targets = append(targets, mention.Mid)
						}
					}
				}
				SaveData()
				if len(targets) != 0 {
					list := "#Addwl:\n"
					for i := range targets {
						list += "\n" + strconv.Itoa(i+1) + ". " + "@!"
					}
					client.SendPollMention(to, list, targets)
				}
				}
			} else if cmd == "whitelist" {
				if !InArray(Data.GroupOwn[to], sender) {
				if len(Data.Whitelist) > 0 {
					list := "#Whitelist:\n"
					for i := range Data.Whitelist {
						list += "\n" + strconv.Itoa(i+1) + ". " + "@!"
						fmt.Sprintf("... %v", i)
					}
					client.SendPollMention(to, list, Data.Whitelist)
				} else { client.SendMessage(to, "Whitelist is empty.") }
				}
			} else if cmd == "admins" {
				if !InArray(Data.GroupOwn[to], sender) {
				if len(Data.Admin) > 0 {
					list := "#Adminlist:\n"
					for i := range Data.Admin {
						list += "\n" + strconv.Itoa(i+1) + ". " + "@!"
						fmt.Sprintf("... %v", i)
					}
					client.SendPollMention(to, list, Data.Admin)
				} else { client.SendMessage(to, "Admin is empty.") }
				}
			} else if strings.HasPrefix(cmd,"admin"+" "){
				if !InArray(Data.GroupOwn[to], sender) {
				fl := strings.Split(cmd, " ")
				typec := strings.Replace(cmd, fl[0]+" ", "", 1)
				msg, _ := client.GetRecentMessagesV2(to, 300)
				for i := 0; i < len(msg); i++ {
					item := msg[i]
					if item.ContentType == 13 {
						mid := item.ContentMetadata["mid"]
						if typec == "lcon" && !fullAccess(client, mid) {
							appendAdmin(mid)
							for i := range ClientBot{ClientBot[i].FindAndAddContactsByMid(mid)}
							client.SendMessage(to, "added to admins.")
							break
						}
					}
				}}
			} else if strings.HasPrefix(cmd, "antitag ") {
				if !InArray(Data.GroupOwn[to], sender) {
				spl := strings.Replace(cmd, "antitag ", "", 1)
				if spl == "on" {
					Data.AntiTag = true
					SaveData()
					client.SendMessage(to, "antitag enabled.")
				} else if spl == "off" {
					Data.AntiTag = false
					SaveData()
					client.SendMessage(to, "antitag disabled.")
				}
				}
			} else if cmd == "banlist" || cmd == Data.Command.Banlist && Data.Command.Banlist != "" {
				if len(Data.Blacklist) > 0 {
					listbl := "#Blacklist:\n"
					for i := range Data.Blacklist {
						listbl += "\n" + strconv.Itoa(i+1) + ". " + "@!"
						fmt.Sprintf("... %v", i)
					}
					client.SendPollMention(to, listbl, Data.Blacklist)
				} else { client.SendMessage(to, "Banlist is empty.") }
			}else if cmd == "bringall"{
				poolBringAll(client, to)
			}else if cmd == "bot all"{
				if !InArray(Data.GroupOwn[to], sender) {
				gc,_ := client.GetGroupV2(to)
				targets := gc.MemberMids
				target := []string{}
				for i:= range targets{
					if !IsGaccess(to,targets[i]){
						if !fullAccess(client,targets[i]){
							if client.MID != targets[i]{
								Data.Bot = append(Data.Bot, targets[i])
								target = append(target, targets[i])
							}
						}
					}
				}
				result := "Bot Allmember:\n"
				if len(target) > 0{
					for i := range target{
						result += "\n"+strconv.Itoa(i+1) + ". @!"
					}
					client.SendPollMention(to,result,target)
				}else{client.SendMessage(to, "have all botlist")}
				}
			} else if cmd == "cban" {
				if !InArray(Data.Master, sender) || !InArray(Data.GroupOwn[to], sender) || InArray(Data.Premlist, sender) {
					if len(Data.Blacklist) > 0 {
						for _, x := range Data.Blacklist {
							if x == "" {
								Data.Blacklist = Remove(Data.Blacklist, x)
								SaveData()
							}
						}
						client.SendPollMention(to, "#Banlist user:\n", Data.Blacklist)
						Data.Blacklist = []string{}
						time.Sleep(1 * time.Second)
						Data.Blacklist = []string{}
						SaveData()
					} else {
						client.SendMessage(to, "Banlist is empty.")
						Data.Blacklist = []string{}
						SaveData()
					}
				}
			} else if cmd == "clearadmin" {
				if !InArray(Data.GroupOwn[to], sender) {
				if len(Data.Admin) != 0 {
					client.SendMessage(to, fmt.Sprintf("Cleared %v adminlist", len(Data.Admin)))
					Data.Admin = []string{}
					SaveData()
				} else { client.SendMessage(to, "Admin is empty.") }
				}
			} else if cmd == "clearban" || cmd == Data.Command.Clearban && Data.Command.Clearban != "" {
				if !InArray(Data.Master, sender) || !InArray(Data.GroupOwn[to], sender) || InArray(Data.Premlist, sender) {
					client.RemoveAllMessage(string(op.Param2))
					if len(Data.Blacklist) != 0 {
						msgcbn := fmt.Sprintf(Data.Message.Ban, len(Data.Blacklist))
						client.SendMessage(to, msgcbn)
						Data.Blacklist = []string{}
						SaveData()
					} else { client.SendMessage(to, "Banlist is empty.") }
				}
			} else if cmd == "cleargadmin" {
				if len(Data.GroupAdm[to]) != 0 {
					client.SendMessage(to, fmt.Sprintf("Cleared %v gadminlist", len(Data.GroupAdm[to])))
					delete(Data.GroupAdm, to)
					SaveData()
				} else { client.SendMessage(to, "Gadmin is empty.") }
			} else if strings.HasPrefix(cmd, "deladmin") {
				if !InArray(Data.GroupOwn[to], sender) {
				mentions := mentions{}
				targets := []string{}
				json.Unmarshal([]byte(op.Message.ContentMetadata["MENTION"]), &mentions)
				for _, mention := range mentions.MENTIONEES {
					if InArray(Data.Admin, string(mention.Mid)) {
						Data.Admin = Remove(Data.Admin, mention.Mid)
						targets = append(targets, mention.Mid)
					}
				}
				SaveData()
				if len(targets) != 0 {
					list := "#Deladmin:\n"
					for i := range targets {
						list += "\n" + strconv.Itoa(i+1) + ". " + "@!"
					}
					client.SendPollMention(to, list, targets)
				}
				}
			} else if strings.HasPrefix(cmd, "delwl") {
				if !InArray(Data.GroupOwn[to], sender) {
				mentions := mentions{}
				targets := []string{}
				json.Unmarshal([]byte(op.Message.ContentMetadata["MENTION"]), &mentions)
				for _, mention := range mentions.MENTIONEES {
					if InArray(Data.Whitelist, string(mention.Mid)) {
						Data.Whitelist = Remove(Data.Whitelist, mention.Mid)
						targets = append(targets, mention.Mid)
					}
				}
				SaveData()
				if len(targets) != 0 {
					list := "#Delwl:\n"
					for i := range targets {
						list += "\n" + strconv.Itoa(i+1) + ". " + "@!"
					}
					client.SendPollMention(to, list, targets)
				}
				}
			} else if strings.HasPrefix(cmd, "delgadmin") {
				mentions := mentions{}
				targets := []string{}
				json.Unmarshal([]byte(op.Message.ContentMetadata["MENTION"]), &mentions)
				for _, mention := range mentions.MENTIONEES {
					if InArray(Data.GroupAdm[to], string(mention.Mid)) {
						Data.GroupAdm[to] = Remove(Data.GroupAdm[to], mention.Mid)
						targets = append(targets, mention.Mid)
					}
				}
				SaveData()
				if len(targets) != 0 {
					list := "#Delgadmin:\n"
					for i := range targets {
						list += "\n" + strconv.Itoa(i+1) + ". " + "@!"
					}
					client.SendPollMention(to, list, targets)
				}
			} else if cmd == "gadmins" {
				if len(Data.GroupAdm[to]) > 0 {
					list := "#Gadminlist:\n"
					for i := range Data.GroupAdm[to] {
						list += "\n" + strconv.Itoa(i+1) + ". " + "@!"
						fmt.Sprintf("... %v", i)
					}
					client.SendPollMention(to, list, Data.GroupAdm[to])
				} else { client.SendMessage(to, "Gadmin is empty.") }
			} else if cmd == "help master" {
				if !InArray(Data.GroupOwn[to], sender) {
				res := "「✿ Menu Master ✿」"
				res += "\n"
				for i, x := range helpmaster {
					res += fmt.Sprintf("\n%v› %s %s", i+1, Data.Setkey, x)
				}
				client.SendMessage(to, res)
				}
			} else if cmd == "list protect" {
				res := " ✿ 𝗟𝗶𝘀𝘁 𝗣𝗿𝗼𝘁𝗲𝗰𝘁 ✿ "
				res += "\n\n↳𝗜𝗻𝘃𝗶𝘁𝗲:"
				for a, x := range Data.ProInvite {
					a += 1
					grup := client.GetGroup(x)
					res += fmt.Sprintf("\n%v. %s", a, grup.Name)
				}
				res += "\n\n↳𝗞𝗶𝗰𝗸:"
				for a, x := range Data.ProKick {
					a += 1
					grup := client.GetGroup(x)
					res += fmt.Sprintf("\n%v. %s", a, grup.Name)
				}
				res += "\n\n↳𝗟𝗶𝗻𝗸:"
				for a, x := range Data.ProQr {
					a += 1
					grup := client.GetGroup(x)
					res += fmt.Sprintf("\n%v. %s", a, grup.Name)
				}
				res += "\n\n↳𝗖𝗮𝗻𝗰𝗲𝗹:"
				for a, x := range Data.ProCancel {
					a += 1
					grup := client.GetGroup(x)
					res += fmt.Sprintf("\n%v. %s", a, grup.Name)
				}
				client.SendMessage(to, res)
			} else if strings.HasPrefix(cmd, "nk:") {
				if !InArray(Data.GroupOwn[to], sender) {
				result := strings.Split((cmd), ":")
				gc, _ := client.GetGroupV2(to)
				targets := gc.MemberMids
				poolTargets := []string{}
				for i := range targets {
					con := client.GetContact(targets[i])
					if strings.Contains(strings.ToLower(con.DisplayName), strings.ToLower(result[1])) {
						if !fullAccess(client, targets[i]) {
							if client.MID != targets[i] {
								poolTargets = append(poolTargets, targets[i])
								appendBl(targets[i])
							}
						}
						runtime.GOMAXPROCS(cpu)
						tl := len(poolTargets)
						var wg sync.WaitGroup
						wg.Add(tl)
						for i := 0; i < tl; i++ {
							go func(i int) {
								defer wg.Done()
								val := []string{poolTargets[i]}
								go func() { client.DeleteOtherFromChat(to, val) }()
							}(i)
						}
						wg.Wait()
					}
				}
				}
			} else if strings.HasPrefix(cmd, "nkill") {
				targets := []string{}
				mentions := mentions{}
				json.Unmarshal([]byte(op.Message.ContentMetadata["MENTION"]), &mentions)
				for _, mention := range mentions.MENTIONEES {
					if !fullAccess(client, string(mention.Mid)) {
						appendBl(mention.Mid)
						GetSimiliarName(client, to, mention.Mid)
						targets = append(targets, mention.Mid)
					}
				}
				SaveData()
				runtime.GOMAXPROCS(cpu)
				tl := len(Data.Blacklist)
				var wg sync.WaitGroup
				wg.Add(tl)
				for i := 0; i < tl; i++ {
					go func(i int) {
						defer wg.Done()
						val := []string{Data.Blacklist[i]}
						go func() { client.DeleteOtherFromChat(to, val) }()
					}(i)
				}
				wg.Wait()
			} else if cmd == "out" || cmd == Data.Command.Out && Data.Command.Out != "" {
				if !InArray(Data.GroupOwn[to], sender) {
				delete(Data.StayGroup, to)
				client.SendMessage(to, Data.Message.Bye)
				client.LeaveGroup(to)
				}
			} else if cmd == "outall" || cmd == Data.Command.Outall && Data.Command.Outall != "" {
				if !InArray(Data.GroupOwn[to], sender) {
				delete(Data.StayGroup, to)
				for i := range ClientBot {
					if ClientBot[i] != client {
						ClientBot[i].LeaveGroup(to)
					}
				}
				}
			} else if cmd == "purge" {
				runtime.GOMAXPROCS(cpu)
				var wg sync.WaitGroup
				wg.Add(len(Data.Blacklist))
				for i:=0;i<len(Data.Blacklist);i++ {
					go func(i int) {
						defer wg.Done()
						client.DeleteOtherFromChat(to,[]string{Data.Blacklist[i]})
			    	}(i)
				}
				wg.Wait()
			} else if cmd == "standall" {
				poolStandAll(client, to)
			} else if cmd == "stayall" {
				poolStayAll(client, to)
			} else if strings.HasPrefix(cmd, "stand ") {
				str := strings.Replace(cmd, "stand ", "", 1)
				result2, _ := strconv.Atoi(str)
				if result2 > 0 && result2 <= len(Data.SquadBots) {
					grup, _ := client.GetGroupV2(to)
					target := grup.MemberMids
					targets := []string{}
					tempInv := []string{}
					batastim := 0
					batastim2 := 0
					delete(Data.StayGroup, to)
					for i := range ClientBot {
						ClientBot[i].NormalKickoutFromGroup(to, []string{"FuckYou"})
						if ClientBot[i].Limited { ClientBot[i].LeaveGroup(to) }
					}
					for i := range target {
						targets = append(targets, target[i])
					}
					_, found := Data.StayGroup[to]
					if found == false {
						for i := range Data.SquadBots {
							if InArray(targets, Data.SquadBots[i]) {
								Data.StayGroup[to] = append(Data.StayGroup[to], Data.SquadBots[i])
							}
						}
					}
					for i := range targets {
						if InArray(Data.StayGroup[to], targets[i]) {
							if batastim < result2 {
								batastim = batastim + 1
							} else {
								ClientMid[targets[i]].LeaveGroup(to)
								Data.StayGroup[to] = Remove(Data.StayGroup[to], targets[i])
							}
						}
					}
					for io := range Data.SquadBots {
						if InArray(targets, Data.SquadBots[io]) {
							batastim2 = batastim2 + 1
						}
					}
					for i := range Data.SquadBots {
						if batastim2 < result2 {
							if !InArray(targets, Data.SquadBots[i]) {
								if ClientMid[Data.SquadBots[i]].Limited == false {
									batastim2 = batastim2 + 1
									tempInv = append(tempInv, Data.SquadBots[i])
									if !InArray(Data.StayGroup[to], Data.SquadBots[i]) {
										Data.StayGroup[to] = append(Data.StayGroup[to], Data.SquadBots[i])
									}
								}
							}
						}
					}
					if len(tempInv) != 0 {
						client.InviteIntoChat(to, tempInv)
					}
					SaveData()
				} else { 
								client.SendMessage(to, "out of range.") 
								break 
							}
			} else if strings.HasPrefix(cmd, "stay ") {
				str := strings.Replace(cmd, "stay ", "", 1)
				ticket := client.ReissueGroupTicket(to)
				result2, _ := strconv.Atoi(str)
				if result2 > 0 && result2 <= len(Data.SquadBots) {
					getmem, _ := client.GetGroupV2(to)
					target := getmem.MemberMids
					targets := []string{}
					batastim := 0
					batastim2 := 0
					delete(Data.StayGroup, to)
					for i := range ClientBot {
						ClientBot[i].NormalKickoutFromGroup(to, []string{"FuckYou"})
						if ClientBot[i].Limited { ClientBot[i].LeaveGroup(to) }
					}
					for i := range target {
						targets = append(targets, target[i])
					}
					_, found := Data.StayGroup[to]
					if found == false {
						for i := range Data.SquadBots {
							if InArray(targets, Data.SquadBots[i]) {
								Data.StayGroup[to] = append(Data.StayGroup[to], Data.SquadBots[i])
							}
						}
					}
					for i := range targets {
						if InArray(Data.StayGroup[to], targets[i]) {
							if batastim < result2 {
								batastim = batastim + 1
							} else {
								ClientMid[targets[i]].LeaveGroup(to)
								Data.StayGroup[to] = Remove(Data.StayGroup[to], targets[i])
							}
						}
					}
					for io := range Data.SquadBots {
						if InArray(targets, Data.SquadBots[io]) {
							batastim2 = batastim2 + 1
						}
					}
					if batastim2 < result2 {
						if getmem.PreventedJoinByTicket == true {
							getmem.PreventedJoinByTicket = false
							client.UpdateGroup(getmem)
						}
					}
					for i := range ClientBot {
						if batastim2 < result2 {
							if !InArray(targets, ClientBot[i].MID) {
								if ClientMid[ClientBot[i].MID].Limited == false {
									err := ClientBot[i].AcceptGroupByTicket(to, ticket)
									if err == nil { batastim2 = batastim2 + 1 }
									if !InArray(Data.StayGroup[to], ClientBot[i].MID) {
										Data.StayGroup[to] = append(Data.StayGroup[to], ClientBot[i].MID)
									}
								}
							}
						}
					}
					if batastim2 == result2 {
						if getmem.PreventedJoinByTicket == false {
							getmem.PreventedJoinByTicket = true
							client.UpdateGroup(getmem)
						}
					}
					SaveData()
				} else { 
								client.SendMessage(to, "out of range.") 
								break 
							}
			} else if cmd == "suffix" {
				syzy := "kick/invite/cancel"
				a := " --- * 𝗦𝘂𝗳𝗳𝗶𝘅 𝗖𝗼𝗺𝗺𝗮𝗻𝗱 * --- "
				a += "\n\n " + syzy + " lkick"
				a += "\n " + syzy + " lkvictim"
				a += "\n " + syzy + " lcancel"
				a += "\n " + syzy + " lcvictim"
				a += "\n " + syzy + " linvite"
				a += "\n " + syzy + " livictim"
				a += "\n " + syzy + " lcloseqr"
				a += "\n " + syzy + " lopenqr"
				a += "\n " + syzy + " ljoin"
				a += "\n " + syzy + " lleave"
				a += "\n " + syzy + " lcon"
				a += "\n " + syzy + " ltag"
				client.SendMessage(to, a)
			} else if strings.HasPrefix(cmd, "unadmin ") {
				if !InArray(Data.GroupOwn[to], sender) {
				result := strings.Split((cmd), " ")
				if result[1] != "0" {
					result2, err := strconv.Atoi(result[1])
					if err != nil { client.SendMessage(to, "invalid number.")
					} else {
						for i := range Data.Admin {
							if result2 > 0 && result2-1 < len(Data.Admin) {
								if i == result2-1 {
									kura := Data.Admin[i]
									Data.Admin = Remove(Data.Admin, kura)
									client.SendMention(to, "success deladmin @!", []string{kura})
									SaveData()
									break
								}
							} else { 
								client.SendMessage(to, "out of range.") 
								break 
							}
						}
					}
				} else { client.SendMessage(to, "invalid range.") }
				}
			} else if strings.HasPrefix(cmd, "unwl ") {
				if !InArray(Data.GroupOwn[to], sender) {
				result := strings.Split((cmd), " ")
				if result[1] != "0" {
					result2, err := strconv.Atoi(result[1])
					if err != nil { client.SendMessage(to, "invalid number.")
					} else {
						for i := range Data.Whitelist {
							if result2 > 0 && result2-1 < len(Data.Whitelist) {
								if i == result2-1 {
									kura := Data.Whitelist[i]
									Data.Whitelist = Remove(Data.Whitelist, kura)
									client.SendMention(to, "success delwl @!", []string{kura})
									SaveData()
									break
								}
							} else { 
								client.SendMessage(to, "out of range.") 
								break 
							}
						}
					}
				} else { client.SendMessage(to, "invalid range.") }
				}
			} else if strings.HasPrefix(cmd, "ungadmin ") {
				result := strings.Split((cmd), " ")
				if result[1] != "0" {
					result2, err := strconv.Atoi(result[1])
					if err != nil { client.SendMessage(to, "invalid number.")
					} else {
						for i := range Data.GroupAdm[to] {
							if result2 > 0 && result2-1 < len(Data.GroupAdm[to]) {
								if i == result2-1 {
									kura := Data.GroupAdm[to][i]
									Data.GroupAdm[to] = Remove(Data.GroupAdm[to], kura)
									client.SendMention(to, "success delgadmin @!", []string{kura})
									SaveData()
									break
								}
							} else { client.SendMessage(to, "out of range.")
							}
						}
					}
				} else { client.SendMessage(to, "invalid range.") }
			} else if strings.HasPrefix(cmd, "vkick ") {
				mentions := mentions{}
				json.Unmarshal([]byte(op.Message.ContentMetadata["MENTION"]), &mentions)
				for _, mention := range mentions.MENTIONEES {
					if IsFriends(client, mention.Mid) == false {
						client.FindAndAddContactsByMid(mention.Mid)
					}
					go func() { client.DeleteOtherFromChat(to, []string{mention.Mid}) }()
					go func() { client.InviteIntoChat(to, []string{mention.Mid}) }()
					go func() { client.CancelChatInvitation(to, []string{mention.Mid}) }()
				}
			} else if strings.HasPrefix(cmd, "cancel ") {
				fl := strings.Split(cmd, " ")
				typec := strings.Replace(cmd, fl[0]+" ", "", 1)
				msg, _ := client.GetRecentMessagesV2(to, 300)
				for i := 0; i < len(msg); i++ {
					item := msg[i]
					if item.ContentType == 18 {
						mids := strings.Split(item.ContentMetadata["LOC_ARGS"], "#4672a8")
						mid := mids[0][:33]
						if typec == "lk" || typec == "lkick" {
							if item.ContentMetadata["LOC_KEY"] == "C_MR" {
								client.CancelChatInvitation(to, []string{mid})
								break
							}
						} else if typec == "lkv" || typec == "lkvictim" {
							if item.ContentMetadata["LOC_KEY"] == "C_MR" {
								mid = strings.Split(mids[0], "\x1e")[1]
								client.CancelChatInvitation(to, []string{mid})
								break
							}
						} else if typec == "lc" || typec == "lcancel" {
							if item.ContentMetadata["LOC_KEY"] == "C_IC" {
								client.CancelChatInvitation(to, []string{mid})
								break
							}
						} else if typec == "lcv" || typec == "lcvictim" {
							if item.ContentMetadata["LOC_KEY"] == "C_IC" {
								mid = strings.Split(mids[0], "\x1e")[1]
								client.CancelChatInvitation(to, []string{mid})
								break
							}
						} else if typec == "li" || typec == "linvite" {
							if item.ContentMetadata["LOC_KEY"] == "C_MI" {
								client.CancelChatInvitation(to, []string{mid})
								break
							}
						} else if typec == "liv" || typec == "livictim" {
							if item.ContentMetadata["LOC_KEY"] == "C_MI" {
								mid = strings.Split(mids[0], "\x1e")[1]
								if !fullAccess(client, mid) {
									appendBl(mid)
								}
								client.CancelChatInvitation(to, []string{mid})
								break
							}
						} else if typec == "lcqr" || typec == "lcloseqr" {
							if item.ContentMetadata["LOC_KEY"] == "C_SP" {
								client.CancelChatInvitation(to, []string{mid})
								break
							}
						} else if typec == "loqr" || typec == "lopenqr" {
							if item.ContentMetadata["LOC_KEY"] == "C_SN" {
								client.CancelChatInvitation(to, []string{mid})
								break
							}
						} else if typec == "lj" || typec == "ljoin" {
							if item.ContentMetadata["LOC_KEY"] == "C_MJ" {
								client.CancelChatInvitation(to, []string{mid})
								break
							}
						} else if typec == "ll" || typec == "lleave" {
							if item.ContentMetadata["LOC_KEY"] == "C_ML" {
								client.CancelChatInvitation(to, []string{mid})
								break
							}
						}
					} else if item.ContentType == 13 {
						mid := item.ContentMetadata["mid"]
						if typec == "lcon" || typec == "lcontact" {
							client.CancelChatInvitation(to, []string{mid})
							break
						}
					} else if item.ContentType == 0 {
						if InArray_dict(item.ContentMetadata, "MENTION") {
							var saodd []string
							mentions := mentions{}
							if typec == "lt" || typec == "ltag" {
								json.Unmarshal([]byte(item.ContentMetadata["MENTION"]), &mentions)
								for _, mention := range mentions.MENTIONEES {
									saodd = append(saodd, mention.Mid)
								}
								for _, miz := range saodd {
									client.CancelChatInvitation(to, []string{miz})
								}
								break
							}
						}
					}
				}
			} else if strings.HasPrefix(cmd, "invite ") {
				fl := strings.Split(cmd, " ")
				typec := strings.Replace(cmd, fl[0]+" ", "", 1)
				msg, _ := client.GetRecentMessagesV2(to, 300)
				for i := 0; i < len(msg); i++ {
					item := msg[i]
					if item.ContentType == 18 {
						mids := strings.Split(item.ContentMetadata["LOC_ARGS"], "#4672a8")
						mid := mids[0][:33]
						if typec == "lk" || typec == "lkick" {
							if item.ContentMetadata["LOC_KEY"] == "C_MR" {
								client.FindAndAddContactsByMid(mid)
								client.InviteIntoChat(to, []string{mid})
								break
							}
						} else if typec == "lkv" || typec == "lkvictim" {
							if item.ContentMetadata["LOC_KEY"] == "C_MR" {
								mid = strings.Split(mids[0], "\x1e")[1]
								client.FindAndAddContactsByMid(mid)
								client.InviteIntoChat(to, []string{mid})
								break
							}
						} else if typec == "lc" || typec == "lcancel" {
							if item.ContentMetadata["LOC_KEY"] == "C_IC" {
								client.FindAndAddContactsByMid(mid)
								client.InviteIntoChat(to, []string{mid})
								break
							}
						} else if typec == "lcv" || typec == "lcvictim" {
							if item.ContentMetadata["LOC_KEY"] == "C_IC" {
								mid = strings.Split(mids[0], "\x1e")[1]
								client.FindAndAddContactsByMid(mid)
								client.InviteIntoChat(to, []string{mid})
								break
							}
						} else if typec == "li" || typec == "linvite" {
							if item.ContentMetadata["LOC_KEY"] == "C_MI" {
								client.FindAndAddContactsByMid(mid)
								client.InviteIntoChat(to, []string{mid})
								break
							}
						} else if typec == "liv" || typec == "livictim" {
							if item.ContentMetadata["LOC_KEY"] == "C_MI" {
								mid = strings.Split(mids[0], "\x1e")[1]
								client.FindAndAddContactsByMid(mid)
								client.InviteIntoChat(to, []string{mid})
								break
							}
						} else if typec == "lcqr" || typec == "lcloseqr" {
							if item.ContentMetadata["LOC_KEY"] == "C_SP" {
								client.FindAndAddContactsByMid(mid)
								client.InviteIntoChat(to, []string{mid})
								break
							}
						} else if typec == "loqr" || typec == "lopenqr" {
							if item.ContentMetadata["LOC_KEY"] == "C_SN" {
								client.FindAndAddContactsByMid(mid)
								client.InviteIntoChat(to, []string{mid})
								break
							}
						} else if typec == "lj" || typec == "ljoin" {
							if item.ContentMetadata["LOC_KEY"] == "C_MJ" {
								client.FindAndAddContactsByMid(mid)
								client.InviteIntoChat(to, []string{mid})
								break
							}
						} else if typec == "ll" || typec == "lleave" {
							if item.ContentMetadata["LOC_KEY"] == "C_ML" {
								client.FindAndAddContactsByMid(mid)
								client.InviteIntoChat(to, []string{mid})
								break
							}
						}
					} else if item.ContentType == 13 {
						mid := item.ContentMetadata["mid"]
						if typec == "lcon" || typec == "lcontact" {
							client.FindAndAddContactsByMid(mid)
							client.InviteIntoChat(to, []string{mid})
							break
						}
					} else if item.ContentType == 0 {
						if InArray_dict(item.ContentMetadata, "MENTION") {
							var saodd []string
							mentions := mentions{}
							if typec == "lt" || typec == "ltag" {
								json.Unmarshal([]byte(item.ContentMetadata["MENTION"]), &mentions)
								for _, mention := range mentions.MENTIONEES {
									saodd = append(saodd, mention.Mid)
								}
								for _, miz := range saodd {
									client.FindAndAddContactsByMid(miz)
								}
								client.InviteIntoChat(to, saodd)
								break
							}
						}
					}
				}
			} else if strings.HasPrefix(cmd, "kick ") || strings.HasPrefix(cmd, Data.Command.Kick) && Data.Command.Kick != "" {
				targets := []string{}
				mentions := mentions{}
				json.Unmarshal([]byte(op.Message.ContentMetadata["MENTION"]), &mentions)
				for _, mention := range mentions.MENTIONEES {
					targets = append(targets, mention.Mid)
					if !fullAccess(client, mention.Mid) {
						appendBl(mention.Mid)
					}
				}
				for i := range ClientBot {
					ClientBot[i].NormalKickoutFromGroup(to, []string{"FuckYou"})
				}
				var wg sync.WaitGroup
				wg.Add(len(targets))
				for i := 0; i < len(targets); i++ {
					go func(i int) {
						defer wg.Done()
						client.DeleteOtherFromChat(to, []string{targets[i]})
					}(i)
				}
				wg.Wait()
				fl := strings.Split(cmd, " ")
				typec := strings.Replace(cmd, fl[0]+" ", "", 1)
				msg, _ := client.GetRecentMessagesV2(to, 300)
				for i := 0; i < len(msg); i++ {
					item := msg[i]
					if item.ContentType == 18 {
						mids := strings.Split(item.ContentMetadata["LOC_ARGS"], "#4672a8")
						mid := mids[0][:33]
						if typec == "lk" || typec == "lkick" {
							if item.ContentMetadata["LOC_KEY"] == "C_MR" {
								if !fullAccess(client, mid) {
									appendBl(mid)
								}
								client.DeleteOtherFromChat(to, []string{mid})
								break
							}
						} else if typec == "lkv" || typec == "lkvictim" {
							if item.ContentMetadata["LOC_KEY"] == "C_MR" {
								mid = strings.Split(mids[0], "\x1e")[1]
								if !fullAccess(client, mid) {
									appendBl(mid)
								}
								client.DeleteOtherFromChat(to, []string{mid})
								break
							}
						} else if typec == "lc" || typec == "lcancel" {
							if item.ContentMetadata["LOC_KEY"] == "C_IC" {
								if !fullAccess(client, mid) {
									appendBl(mid)
								}
								client.DeleteOtherFromChat(to, []string{mid})
								break
							}
						} else if typec == "lcv" || typec == "lcvictim" {
							if item.ContentMetadata["LOC_KEY"] == "C_IC" {
								mid = strings.Split(mids[0], "\x1e")[1]
								if !fullAccess(client, mid) {
									appendBl(mid)
								}
								client.DeleteOtherFromChat(to, []string{mid})
								break
							}
						} else if typec == "li" || typec == "linvite" {
							if item.ContentMetadata["LOC_KEY"] == "C_MI" {
								if !fullAccess(client, mid) {
									appendBl(mid)
								}
								client.DeleteOtherFromChat(to, []string{mid})
								break
							}
						} else if typec == "liv" || typec == "livictim" {
							if item.ContentMetadata["LOC_KEY"] == "C_MI" {
								mid = strings.Split(mids[0], "\x1e")[1]
								if !fullAccess(client, mid) {
									appendBl(mid)
								}
								client.DeleteOtherFromChat(to, []string{mid})
								break
							}
						} else if typec == "lcqr" || typec == "lcloseqr" {
							if item.ContentMetadata["LOC_KEY"] == "C_SP" {
								if !fullAccess(client, mid) {
									appendBl(mid)
								}
								client.DeleteOtherFromChat(to, []string{mid})
								break
							}
						} else if typec == "loqr" || typec == "lopenqr" {
							if item.ContentMetadata["LOC_KEY"] == "C_SN" {
								if !fullAccess(client, mid) {
									appendBl(mid)
								}
								client.DeleteOtherFromChat(to, []string{mid})
								break
							}
						} else if typec == "lj" || typec == "ljoin" {
							if item.ContentMetadata["LOC_KEY"] == "C_MJ" {
								if !fullAccess(client, mid) {
									appendBl(mid)
								}
								client.DeleteOtherFromChat(to, []string{mid})
								break
							}
						} else if typec == "ll" || typec == "lleave" {
							if item.ContentMetadata["LOC_KEY"] == "C_ML" {
								if !fullAccess(client, mid) {
									appendBl(mid)
								}
								client.DeleteOtherFromChat(to, []string{mid})
								break
							}
						}
					} else if item.ContentType == 13 {
						mid := item.ContentMetadata["mid"]
						if typec == "lcon" || typec == "lcontact" {
							if !fullAccess(client, mid) {
								appendBl(mid)
							}
							client.DeleteOtherFromChat(to, []string{mid})
							break
						}
					} else if item.ContentType == 0 {
						if InArray_dict(item.ContentMetadata, "MENTION") {
							var saodd []string
							if typec == "lt" || typec == "ltag" {
								json.Unmarshal([]byte(item.ContentMetadata["MENTION"]), &mentions)
								for _, mention := range mentions.MENTIONEES {
									saodd = append(saodd, mention.Mid)
								}
								for _, miz := range saodd {
									if !fullAccess(client, miz) {
										appendBl(miz)
									}
									client.DeleteOtherFromChat(to, []string{miz})
								}
								break
							}
						}
					}
				}
				if typec == "count" {
					X, _ := client.GetGroupV2(to)
					memb := X.MemberMids
					var a = 0
					ret := "#Historycount:"
					ret += "\n"
					for i := range ClientBot {
						if InArray(memb, ClientBot[i].MID) {
							Ki := fmt.Sprintf("\n    -kick %v, -inv %v, -cancel %v", ClientBot[i].Ckick, ClientBot[i].Cinvite, ClientBot[i].Ccancel)
							a = a + 1
							ret += fmt.Sprintf("\nAssist%v: %s", a, Ki)
						}
					}
					client.SendMessage(to, ret)
					break
				}
			} else if strings.HasPrefix(cmd, "ban") {
				fl := strings.Split(cmd, " ")
				typec := strings.Replace(cmd, fl[0]+" ", "", 1)
				msg, _ := client.GetRecentMessagesV2(to, 300)
				targets := []string{}
				for i := 0; i < len(msg); i++ {
					item := msg[i]
					if item.ContentType == 18 {
						mids := strings.Split(item.ContentMetadata["LOC_ARGS"], "#4672a8")
						mid := mids[0][:33]
						if typec == "lk" || typec == "lkick" {
							if item.ContentMetadata["LOC_KEY"] == "C_MR" {
								if !fullAccess(client, mid) {
									targets = append(targets, mid)
									appendBl(mid)
								}
								break
							}
						} else if typec == "lkv" || typec == "lkvictim" {
							if item.ContentMetadata["LOC_KEY"] == "C_MR" {
								mid = strings.Split(mids[0], "\x1e")[1]
								if !fullAccess(client, mid) {
									targets = append(targets, mid)
									appendBl(mid)
								}
								break
							}
						} else if typec == "lc" || typec == "lcancel" {
							if item.ContentMetadata["LOC_KEY"] == "C_IC" {
								if !fullAccess(client, mid) {
									targets = append(targets, mid)
									appendBl(mid)
								}
								break
							}
						} else if typec == "lcv" || typec == "lcvictim" {
							if item.ContentMetadata["LOC_KEY"] == "C_IC" {
								mid = strings.Split(mids[0], "\x1e")[1]
								if !fullAccess(client, mid) {
									targets = append(targets, mid)
									appendBl(mid)
								}
								break
							}
						} else if typec == "li" || typec == "linvite" {
							if item.ContentMetadata["LOC_KEY"] == "C_MI" {
								if !fullAccess(client, mid) {
									targets = append(targets, mid)
									appendBl(mid)
								}
								break
							}
						} else if typec == "liv" || typec == "livictim" {
							if item.ContentMetadata["LOC_KEY"] == "C_MI" {
								mid = strings.Split(mids[0], "\x1e")[1]
								if !fullAccess(client, mid) {
									targets = append(targets, mid)
									appendBl(mid)
								}
								break
							}
						} else if typec == "lcqr" || typec == "lcloseqr" {
							if item.ContentMetadata["LOC_KEY"] == "C_SP" {
								if !fullAccess(client, mid) {
									targets = append(targets, mid)
									appendBl(mid)
								}
								break
							}
						} else if typec == "loqr" || typec == "lopenqr" {
							if item.ContentMetadata["LOC_KEY"] == "C_SN" {
								if !fullAccess(client, mid) {
									targets = append(targets, mid)
									appendBl(mid)
								}
								break
							}
						} else if typec == "lj" || typec == "ljoin" {
							if item.ContentMetadata["LOC_KEY"] == "C_MJ" {
								if !fullAccess(client, mid) {
									targets = append(targets, mid)
									appendBl(mid)
								}
								break
							}
						} else if typec == "ll" || typec == "lleave" {
							if item.ContentMetadata["LOC_KEY"] == "C_ML" {
								if !fullAccess(client, mid) {
									targets = append(targets, mid)
									appendBl(mid)
								}
								break
							}
						}
					} else if item.ContentType == 13 {
						mid := item.ContentMetadata["mid"]
						if typec == "lcon" || typec == "lcontact" {
							if !fullAccess(client, mid) {
								targets = append(targets, mid)
								appendBl(mid)
							}
							break
						}
					} else if item.ContentType == 0 {
						if InArray_dict(item.ContentMetadata, "MENTION") {
							var saodd []string
							mentions := mentions{}
							if typec == "lt" || typec == "ltag" {
								json.Unmarshal([]byte(item.ContentMetadata["MENTION"]), &mentions)
								for _, mention := range mentions.MENTIONEES {
									targets = append(targets, mention.Mid)
									saodd = append(saodd, mention.Mid)
								}
								for _, miz := range saodd {
									if !fullAccess(client, miz) {
										appendBl(miz)
									}
								}
								break
							}
						}
					}
				}
				if len(targets) != 0 {
					list := "#Addban:\n"
					for i := range targets {
						list += "\n" + strconv.Itoa(i+1) + ". " + "@!"
					}
					client.SendPollMention(to, list, targets)
				}
			} else if strings.HasPrefix(cmd, "unban ") {
				if !InArray(Data.Master, sender) || !InArray(Data.GroupOwn[to], sender) || InArray(Data.Premlist, sender) {
					result := strings.Split((cmd), " ")
					if result[1] != "0" {
						result2, err := strconv.Atoi(result[1])
						if err != nil {
							typec := strings.ToLower(result[1])
							msg, _ := client.GetRecentMessagesV2(to, 300)
							targets := []string{}
							for i := 0; i < len(msg); i++ {
								item := msg[i]
								if item.ContentType == 18 {
									mids := strings.Split(item.ContentMetadata["LOC_ARGS"], "#4672a8")
									mid := mids[0][:33]
									if typec == "lk" || typec == "lkick" {
										if item.ContentMetadata["LOC_KEY"] == "C_MR" {
											targets = append(targets, mid)
											removeBl(mid)
											break
										}
									} else if typec == "lc" || typec == "lcancel" {
										if item.ContentMetadata["LOC_KEY"] == "C_IC" {
											targets = append(targets, mid)
											removeBl(mid)
											break
										}
									} else if typec == "lkv" || typec == "lkvictim" {
										if item.ContentMetadata["LOC_KEY"] == "C_MR" {
											mid = strings.Split(mids[0], "\x1e")[1]
											targets = append(targets, mid)
											removeBl(mid)
											break
										}
									} else if typec == "lcv" || typec == "lcvictim" {
										if item.ContentMetadata["LOC_KEY"] == "C_IC" {
											mid = strings.Split(mids[0], "\x1e")[1]
											targets = append(targets, mid)
											removeBl(mid)
											break
										}
									} else if typec == "li" || typec == "linvite" {
										if item.ContentMetadata["LOC_KEY"] == "C_MI" {
											targets = append(targets, mid)
											removeBl(mid)
											break
										}
									} else if typec == "lcqr" || typec == "lcloseqr" {
										if item.ContentMetadata["LOC_KEY"] == "C_SP" {
											targets = append(targets, mid)
											removeBl(mid)
											break
										}
									} else if typec == "loqr" || typec == "lopenqr" {
										if item.ContentMetadata["LOC_KEY"] == "C_SN" {
											targets = append(targets, mid)
											removeBl(mid)
											break
										}
									} else if typec == "lj" || typec == "ljoin" {
										if item.ContentMetadata["LOC_KEY"] == "C_MJ" {
											targets = append(targets, mid)
											removeBl(mid)
											break
										}
									} else if typec == "ll" || typec == "lleave" {
										if item.ContentMetadata["LOC_KEY"] == "C_ML" {
											targets = append(targets, mid)
											removeBl(mid)
											break
										}
									}
								} else if item.ContentType == 13 {
									mid := item.ContentMetadata["mid"]
									if typec == "lcon" || typec == "lcontact" {
										targets = append(targets, mid)
										removeBl(mid)
										break
									}
								} else if item.ContentType == 0 {
									if InArray_dict(item.ContentMetadata, "MENTION") {
										var saodd []string
										mentions := mentions{}
										if typec == "lt" || typec == "ltag" {
											json.Unmarshal([]byte(item.ContentMetadata["MENTION"]), &mentions)
											for _, mention := range mentions.MENTIONEES {
												targets = append(targets, mention.Mid)
												saodd = append(saodd, mention.Mid)
											}
											for _, miz := range saodd {
												removeBl(miz)
											}
											break
										}
									}
								}
							}
							if len(targets) != 0 {
								list := "#Delban:\n"
								for i := range targets {
									list += "\n" + strconv.Itoa(i+1) + ". " + "@!"
								}
								client.SendPollMention(to, list, targets)
							}
						} else {
							for i := range Data.Blacklist {
								if result2 > 0 && result2-1 < len(Data.Blacklist) {
									if i == result2-1 {
										kura := Data.Blacklist[i]
										Data.Blacklist = Remove(Data.Blacklist, kura)
										client.SendMention(to, "success delban @!", []string{kura})
										SaveData()
										break
									}
								} else { 
								client.SendMessage(to, "out of range.") 
								break 
							}
							}
						}
					} else { client.SendMessage(to, "invalid range.") }
				}
			} else if strings.HasPrefix(cmd, "deleteban ") {
				if !InArray(Data.Master, sender) || !InArray(Data.GroupOwn[to], sender) || InArray(Data.Premlist, sender) {
					str := strings.Replace(cmd, "deleteban ", "", 1)
					spl := strings.Split(str, "-")
					target := []string{}
					no, _ := strconv.Atoi(spl[0])
					num, _ := strconv.Atoi(spl[1])
					var a = no - 1
					if no > 0 && num <= len(Data.Blacklist) {
						for x, mid := range Data.Blacklist {
							if a == x {
								Data.Blacklist = Remove(Data.Blacklist, mid)
								SaveData()
								a += 1
								target = append(target, mid)
								if a == num {
									break
								}
							}
						}
					}
					if len(target) != 0 {
						list := "#Deleteban:\n"
						for i := range target {
							list += "\n" + strconv.Itoa(i+1) + ". " + "@!"
						}
						client.SendPollMention(to, list, target)
					}
				}
			}
		}
//ADMINS
		if IsAdminSender(sender) || IsGaccess(to, sender) {
			if pesan == "sname" {
				client.SendMessage(to, Data.Setkey)
			} else if pesan == "rname" {
				client.SendMessage(to, Data.Rname)
			} else if pesan == Data.Setkey {
				client.SendMessage(to, Data.Message.Respon)
			} else if pesan == Data.Rname {
				client.SendMessage(to, Data.Message.Respon)
			}
			if cmd == "abort" {
				if _, cek := ContactType[sender]; cek {
					delete(ContactType, sender)
				}
				client.SendMessage(to, "Abort done.")
			} else if cmd == "about" {
				loc, _ := time.LoadLocation("Asia/Jakarta")
				groups := client.GetGroupsJoined()
				a := time.Now().In(loc)
				base := time.Date(a.Year(), a.Month(), a.Day(), a.Hour(), a.Minute(), a.Second(), 0, loc)
				td := timeutil.Timedelta{Days: time.Duration(duedatecount)}
				exp := base.Add(td.Duration())
				rst := " ☘️ | 𝗕𝗢𝗧 𝗜𝗡𝗙𝗢"
				rst += fmt.Sprintf("\n\n  ↳System: %s", runtime.GOOS)
				rst += fmt.Sprintf("\n  ↳Version: %s", runtime.Version())
				rst += "\n  ↳Type: Multi Selfbot"
				rst += "\n  ↳Base: Line/11.18.2"
				rst += "\n  ↳Lang: Go"
				rst += "\n\n ×̾ | 𝗔𝗖𝗖𝗢𝗨𝗡𝗧 𝗜𝗡𝗙𝗢"
				rst += "\n ⚬  Fetch bots: " + string(client.GetContact(client.MID).DisplayName)
				rst += fmt.Sprintf("\n ⚬  AllBots: %v", len(client.Squads)+1)
				rst += "\n ⚬  Groups: " + strconv.Itoa(len(groups))
				rst += "\n ⚬  Prefix: " + Data.Setkey
				rst += "\n ⚬  Expired: " + (exp).String()[:10]
				rst += "\n ⚬  Updates: 2021-12-01"
				rst += "\n\n  ᵖᵒʷᵉʳᵉᵈ ᵇʸ\n         Friendly™ ᴮᴼᵀ"
				client.SendMessageFooter(to, rst, IconLink, IconFooter, Data.Message.Flag)
			} else if strings.HasPrefix(cmd, "addban") {
				mentions := mentions{}
				targets := []string{}
				json.Unmarshal([]byte(op.Message.ContentMetadata["MENTION"]), &mentions)
				for _, mention := range mentions.MENTIONEES {
					if !InArray(Data.Blacklist, string(mention.Mid)) && !fullAccess(client, mention.Mid) {
						if !fullAccess(client, mention.Mid) && mention.Mid != client.MID {
							Data.Blacklist = append(Data.Blacklist, mention.Mid)
							targets = append(targets, mention.Mid)
						}
					}
				}
				SaveData()
				if len(targets) != 0 {
					list := "#Addban:\n"
					for i := range targets {
						list += "\n" + strconv.Itoa(i+1) + ". " + "@!"
					}
					client.SendPollMention(to, list, targets)
				}
			} else if strings.HasPrefix(cmd, "addbot") {
				if !InArray(Data.Master, sender) || !InArray(Data.GroupOwn[to], sender) {
				mentions := mentions{}
				targets := []string{}
				json.Unmarshal([]byte(op.Message.ContentMetadata["MENTION"]), &mentions)
				for _, mention := range mentions.MENTIONEES {
					if !InArray(Data.Bot, string(mention.Mid)) && client.MID != mention.Mid {
						if !fullAccess(client, mention.Mid) && mention.Mid != client.MID {
							Data.Bot = append(Data.Bot, mention.Mid)
							targets = append(targets, mention.Mid)
						}
					}
				}
				SaveData()
				if len(targets) != 0 {
					list := "#Addbot:\n"
					for i := range targets {
						list += "\n" + strconv.Itoa(i+1) + ". " + "@!"
					}
					client.SendPollMention(to, list, targets)
				}
				}
			} else if strings.HasPrefix(cmd, "addgban") {
				mentions := mentions{}
				targets := []string{}
				json.Unmarshal([]byte(op.Message.ContentMetadata["MENTION"]), &mentions)
				for _, mention := range mentions.MENTIONEES {
					if !InArray(Data.GroupBan[to], string(mention.Mid)) && client.MID != mention.Mid {
						if !fullAccess(client, mention.Mid) && mention.Mid != client.MID {
							Data.GroupBan[to] = append(Data.GroupBan[to], mention.Mid)
							targets = append(targets, mention.Mid)
						}
					}
				}
				SaveData()
				if len(targets) != 0 {
					list := "#Addgban:\n"
					for i := range targets {
						list += "\n" + strconv.Itoa(i+1) + ". " + "@!"
					}
					client.SendPollMention(to, list, targets)
				}
			} else if cmd == "arrange" {
				SaveData()
				client.SendMessage(to, "done.")
			} else if cmd == "bans" {
				limiter := []string{}
				for i := range ClientBot {
					ClientBot[i].NormalKickoutFromGroup(to, []string{"FuckYou"})
					if _, cek := limiterBot[ClientBot[i].MID]; cek {
						limiter = append(limiter, ClientBot[i].MID)
					}
				}
				if len(limiter) != 0 {
					var no = 1
					ret := fmt.Sprintf("%v/%v bots on limits.", len(limiter), len(client.Squads)+1)
					ret += "\n"
					for _, x := range limiter {
						wkt := limiterBot[x]
						waktu := time.Since(wkt)
						if client.CheckLimited(wkt) {
							limit := limitDuration(waktu)
							con := client.GetContact(x)
							client.SendContact(to, x)
							ret += fmt.Sprintf("\n%v. %s", no, con.DisplayName)
							ret += fmt.Sprintf("\nLimits have: %s\n", limit)
							no++
						}
					}
					client.SendMessage(to, ret)
				} else { client.SendMessage(to, "Allbots fix.") }
			} else if cmd == "ban:on" {
				ContactType[sender] = "ban"
				client.SendMessage(to, "Send contact.")
			} else if cmd == "bot:on" {
				if !InArray(Data.Master, sender) || !InArray(Data.GroupOwn[to], sender) {
				ContactType[sender] = "bot"
				client.SendMessage(to, "Send contact")
				}
			} else if cmd == "bots" {
				if !InArray(Data.Master, sender) || !InArray(Data.GroupOwn[to], sender) {
				if len(Data.Bot) > 0 {
					list := "#Botlist:\n"
					for i := range Data.Bot {
						list += "\n" + strconv.Itoa(i+1) + ". " + "@!"
						fmt.Sprintf("... %v", i)
					}
					client.SendPollMention(to, list, Data.Bot)
				} else { client.SendMessage(to, "Bot is empty.") }
				}
			} else if strings.HasPrefix(cmd,"bot"+" "){
				if !InArray(Data.Master, sender) || !InArray(Data.GroupOwn[to], sender) {
				fl := strings.Split(cmd, " ")
				typec := strings.Replace(cmd, fl[0]+" ", "", 1)
				msg, _ := client.GetRecentMessagesV2(to, 300)
				for i := 0; i < len(msg); i++ {
					item := msg[i]
					if item.ContentType == 13 {
						mid := item.ContentMetadata["mid"]
						if typec == "lcon" && !fullAccess(client, mid) {
							appendBot(mid)
							for i := range ClientBot{ClientBot[i].FindAndAddContactsByMid(mid)}
							client.SendMessage(to, "added to admins.")
							break
						}
					}
				}}
			} else if cmd == "clearbot" {
				if !InArray(Data.Master, sender) || !InArray(Data.GroupOwn[to], sender) {
				if len(Data.Bot) != 0 {
					client.SendMessage(to, fmt.Sprintf("Cleared %v botlist", len(Data.Bot)))
					Data.Bot = []string{}
					SaveData()
				} else { client.SendMessage(to, "Bot is empty.") }
				}
			} else if cmd == "cleargban" {
				if len(Data.GroupBan[to]) != 0 {
					client.SendMessage(to, fmt.Sprintf("Cleared %v gbanlist", len(Data.GroupBan[to])))
					delete(Data.GroupBan, to)
					SaveData()
				} else { client.SendMessage(to, "Gban is empty.") }
			} else if cmd == "clearchat" {
				for i := range ClientBot {
					ClientBot[i].RemoveAllMessage(string(op.Param2))
				}
				client.SendMessage(to, "Cleared all message.")
			}else if cmd == "clearcache"{
				exec.Command("bash","-c","sudo systemd-resolve --flush-caches").Output()
				exec.Command("bash","-c","echo 3 > /proc/sys/vm/drop_caches&&swapoff -a&&swapon -a").Output()
				client.SendMessage(to, "Cleared all cache.")
			} else if cmd == "check" {
				temprs := []string{}
				nomor := 1
				nomorx := 1
				r := "Assist inside:\n"
				_, found := Data.StayGroup[to]
				if found == false {
					r += "\n\nAssist outside:\n"
					for ii := range ClientBot {
						temprs = append(temprs, ClientBot[ii].MID)
						r += fmt.Sprintf("\n"+"%v. @!", nomorx)
						nomorx += 1
					}
					client.SendPollMention(to, r, temprs)
				} else {
					if len(Data.StayGroup[to]) > 0 {
						for i := range Data.StayGroup[to] {
							temprs = append(temprs, Data.StayGroup[to][i])
							r += fmt.Sprintf("\n"+"%v. @!", nomor)
							nomor += 1
						}
						r += "\n\nAssist outside:\n"
						for ii := range ClientBot {
							if !InArray(Data.StayGroup[to], ClientBot[ii].MID) {
								temprs = append(temprs, ClientBot[ii].MID)
								r += fmt.Sprintf("\n"+"%v. @!", nomorx)
								nomorx += 1
							}
						}
						client.SendPollMention(to, r, temprs)
					} else {
						r += "\n\nAssist outside:\n"
						for ii := range ClientBot {
							temprs = append(temprs, ClientBot[ii].MID)
							r += fmt.Sprintf("\n"+"%v. @!", nomorx)
							nomorx += 1
						}
						client.SendPollMention(to, r, temprs)
					}
				}
			} else if strings.HasPrefix(cmd,"cmid"){
				result := strings.Split((cmd)," ")
				r := strings.Replace(cmd,result[0]+" ", "", 1)
				client.SendContact(to, r)
			} else if cmd == "count" || cmd == Data.Command.Count && Data.Command.Count != "" {
				X, _ := client.GetGroupV2(to)
				memb := X.MemberMids
				var a = 0
				for i := range ClientBot {
					if InArray(memb, ClientBot[i].MID) {
						a = a + 1
						ClientBot[i].SendMessage(to, fmt.Sprintf("%v", a))
					}
				}
			} else if strings.HasPrefix(cmd, "delban") {
				if !InArray(Data.Master, sender) || !InArray(Data.GroupOwn[to], sender) || !InArray(Data.GroupAdm[to], sender) || InArray(Data.Premlist, sender) {
					mentions := mentions{}
					targets := []string{}
					json.Unmarshal([]byte(op.Message.ContentMetadata["MENTION"]), &mentions)
					for _, mention := range mentions.MENTIONEES {
						if InArray(Data.Blacklist, string(mention.Mid)) {
							Data.Blacklist = Remove(Data.Blacklist, mention.Mid)
							targets = append(targets, mention.Mid)
						}
					}
					SaveData()
					if len(targets) != 0 {
						list := "#Delban:\n"
						for i := range targets {
							list += "\n" + strconv.Itoa(i+1) + ". " + "@!"
						}
						client.SendPollMention(to, list, targets)
					}
				}
			} else if strings.HasPrefix(cmd, "delbot") {
				if !InArray(Data.Master, sender) || !InArray(Data.GroupOwn[to], sender) {
				mentions := mentions{}
				targets := []string{}
				json.Unmarshal([]byte(op.Message.ContentMetadata["MENTION"]), &mentions)
				for _, mention := range mentions.MENTIONEES {
					if InArray(Data.Bot, string(mention.Mid)) {
						Data.Bot = Remove(Data.Bot, mention.Mid)
						targets = append(targets, mention.Mid)
					}
				}
				SaveData()
				if len(targets) != 0 {
					list := "#Delbot:\n"
					for i := range targets {
						list += "\n" + strconv.Itoa(i+1) + ". " + "@!"
					}
					client.SendPollMention(to, list, targets)
				}
				}
			} else if strings.HasPrefix(cmd, "delgban") {
				mentions := mentions{}
				targets := []string{}
				json.Unmarshal([]byte(op.Message.ContentMetadata["MENTION"]), &mentions)
				for _, mention := range mentions.MENTIONEES {
					if InArray(Data.GroupBan[to], string(mention.Mid)) {
						Data.GroupBan[to] = Remove(Data.GroupBan[to], mention.Mid)
						targets = append(targets, mention.Mid)
					}
				}
				SaveData()
				if len(targets) != 0 {
					list := "#Delgban:\n"
					for i := range targets {
						list += "\n" + strconv.Itoa(i+1) + ". " + "@!"
					}
					client.SendPollMention(to, list, targets)
				}
			} else if cmd == "expel:on" {
				if !InArray(Data.Master, sender) || !InArray(Data.GroupOwn[to], sender) {
				ContactType[sender] = "expel"
				client.SendMessage(to, "Send contact.")
				}
			} else if cmd == "gaccess"{
				if !InArray(Data.Master, sender) || !InArray(Data.GroupOwn[to], sender) {
				allmanagers := []string{}
				nourut := 1
				listadm := " *** GOWNER *** "
				if len(Data.GroupOwn[to]) > 0{
					for i:= range Data.GroupOwn[to]{
						allmanagers = append(allmanagers,Data.GroupOwn[to][i])
						listadm += "\n"+strconv.Itoa(i+nourut) + ". @!"
					}
					nourut = len(Data.Buyer)+1
				}
				if len(Data.GroupAdm[to]) > 0{
					listadm += "\n\n *** GADMIN *** "
					for i:= range Data.GroupAdm[to]{
						allmanagers = append(allmanagers,Data.GroupAdm[to][i])
						listadm += "\n"+strconv.Itoa(i+nourut) + ". @!"
					}
					nourut = len(Data.GroupAdm[to])+2
				}else{listadm += "\n\n *** GADMIN *** "}
				if len(Data.GroupBan[to]) > 0{
					listadm += "\n\n *** GBANS *** "
					for i:= range Data.GroupBan[to]{
						allmanagers = append(allmanagers,Data.GroupBan[to][i])
						listadm += "\n"+strconv.Itoa(i+nourut) + ". @!"
					}
					nourut = len(Data.GroupBan[to])+3
				}else{listadm += "\n\n *** GBANS *** "}
				if len(allmanagers) != 0 {
					client.SendPollMention(to,listadm,allmanagers)
				}else{client.SendMessage(to, "Gaccess is empty.")}
				}
			} else if cmd == "gbans" {
				if len(Data.GroupBan[to]) > 0 {
					list := "#Gbanlist:\n"
					for i := range Data.GroupBan[to] {
						list += "\n" + strconv.Itoa(i+1) + ". " + "@!"
						fmt.Sprintf("... %v", i)
					}
					client.SendPollMention(to, list, Data.GroupBan[to])
				} else { client.SendMessage(to, "Gban is empty.") }
			} else if cmd == "groupinfo" {
				var gQr string
				var gTicket string
				gc, _ := client.GetGroupV2(to)
				tick := client.ReissueGroupTicket(to)
				i := time.Unix(gc.CreatedTime/1000, 0)
				if gc.PreventedJoinByTicket == true {
					gQr += "Closed"
					gTicket += "Not Found"
				} else {
					gQr += "Opened"
					gTicket += "https://line.me/R/ti/g/" + tick
				}
				ret_ := "╭─「 Grup Info 」─"
				ret_ += "\n├≽ Group Name : " + gc.Name
				ret_ += "\n├≽ ID : " + gc.ID
				ret_ += "\n├≽ Made On : " + i.String()[:19]
				ret_ += "\n├≽ Group Member : " + strconv.Itoa(len(gc.MemberMids)) + " Memb"
				ret_ += "\n├≽ Group Pending : " + strconv.Itoa(len(gc.InviteeMids)) + " Pend"
				ret_ += "\n├≽ Group Qr : " + gQr
				ret_ += "\n├≽ Ticket : " + gTicket
				ret_ += "\n╰─「 Done 」"
				client.SendMessageFooter(to, ret_, IconLink, "https://obs.line-scdn.net/"+gc.PictureStatus, gc.Name)
				client.SendContact(to, gc.Creator.Mid)
			} else if cmd == "help" {
				res := "Visit All Command:"
				res += "\n	Help All"
				res += "\n"
				res += "\nVisit Protection:"
				res += "\n	Help Pro"
				res += "\n"
				res += "\nVisit Permission:"
				res += "\n	Help Maker"
				res += "\n	Help Buyer"
				res += "\n	Help Owner"
				res += "\n	Help Master"
				res += "\n	Help Admin"
				client.SendMessage(to, res)
			} else if cmd == "help all" {
				res := "「✿ Menu Help ✿」"
				res += "\n"
				res += "\n	✿ Protection ✿"
				res += "\n"
				for a, x := range helppro {
					res += fmt.Sprintf("\n%v› %s %s", a+1, Data.Setkey, x)
				}
				res += "\n"
				res += "\n	✿ Owner ✿"
				res += "\n"
				for a, x := range helpowner {
					res += fmt.Sprintf("\n%v› %s %s", a+1, Data.Setkey, x)
				}
				res += "\n"
				res += "\n	✿ Master ✿"
				res += "\n"
				for i, x := range helpmaster {
					res += fmt.Sprintf("\n%v› %s %s", i+1, Data.Setkey, x)
				}
				res += "\n"
				res += "\n	✿ Admin ✿"
				res += "\n"
				for a, x := range helpadmin {
					res += fmt.Sprintf("\n%v› %s %s", a+1, Data.Setkey, x)
				}
				res += "\n\nᵖᵒʷᵉʳᵉᵈ ᵇʸ\n		" + Data.Message.Flag + ""
				client.SendMessage(to, res)
			} else if cmd == "help admin" {
				if !InArray(Data.Master, sender) || !InArray(Data.GroupOwn[to], sender) {
				res := "「✿ Menu admin ✿」"
				res += "\n"
				for a, x := range helpadmin {
					res += fmt.Sprintf("\n%v› %s %s", a+1, Data.Setkey, x)
				}
				client.SendMessage(to, res)
				}
			} else if cmd == "here" {
				targets := []string{client.MID}
				for _, x := range client.Squads {
					if IsMembers(client, to, x) == true {
						targets = append(targets, x)
					}
				}
				client.SendMessage(to, fmt.Sprintf("%v/%v bots here.", len(targets), len(client.Squads)+1))
			} else if cmd == "leave" || cmd == Data.Command.Leave && Data.Command.Leave != "" {
				if !InArray(Data.Master, sender) || !InArray(Data.GroupOwn[to], sender) {
				_, found := Data.StayGroup[to]
				if found == true {
					delete(Data.StayGroup, to)
				}
				for i := range ClientBot {
					ClientBot[i].LeaveGroup(to)
				}
				}
			} else if cmd == "limitout" {
				for i := range ClientBot {
					ClientBot[i].NormalKickoutFromGroup(to, []string{"FuckYou"})
					if ClientBot[i].Limited == true {
						ClientBot[i].LeaveGroup(to)
						if InArray(Data.StayGroup[to], ClientBot[i].MID) {
							Data.StayGroup[to] = Remove(Data.StayGroup[to], ClientBot[i].MID)
						}
					}
				}
				SaveData()
			}else if strings.HasPrefix(cmd,"mid "){
				mentions := mentions{}
				json.Unmarshal([]byte(op.Message.ContentMetadata["MENTION"]), &mentions)
				for _, mention := range mentions.MENTIONEES{ client.SendMessage(to, mention.Mid) }
				fl := strings.Split(cmd, " ")
				typec := strings.Replace(cmd, fl[0]+" ", "", 1)
				msg, _ := client.GetRecentMessagesV2(to, 300)
				for i := 0; i < len(msg); i++ {
					item := msg[i] 
					if item.ContentType == 13 {
						mid := item.ContentMetadata["mid"]
						if typec == "lcon" || typec == "lcontact" {
							client.SendMessage(to, mid)
							break
						}
					}
				}
			} else if cmd == "ourl" {
				gc := client.GetGroup(to)
				tick := client.ReissueGroupTicket(to)
				if gc.PreventedJoinByTicket == true {
					gc.PreventedJoinByTicket = false
					client.UpdateGroup(gc)
				}
				client.SendMessage(to, "https://line.me/R/ti/g/"+tick)
			} else if cmd == "curl" {
				gc := client.GetGroup(to)
				if gc.PreventedJoinByTicket == false {
					gc.PreventedJoinByTicket = true
					client.UpdateGroup(gc)
				}
			} else if cmd == "rollcall" {
				for i := range ClientBot {
					ClientBot[i].SendMessage(to, string(client.GetContact(ClientBot[i].MID).DisplayName))
				}
			} else if cmd == "respon" || cmd == Data.Command.Respon && Data.Command.Respon != "" {
				for i := range ClientBot {
					ClientBot[i].SendMessage(to, Data.Message.Respon)
				}
			} else if strings.HasPrefix(cmd, "say ") {
				str := strings.Replace(cmd, "say ", "", 1)
				for i := range ClientBot {
					ClientBot[i].SendMessage(to, str)
				}
			}else if cmd == "timeleft"{
				client.CheckExprd()
				client.SendMessage(to, "Expired in: "+strconv.Itoa(duedatecount)+" days")
			}else if cmd == "timenow"{
				GenerateTimeLog(client,to)
			}else if cmd == "runtime"{
				elapsed := time.Since(botStart)
				client.SendMessage(to, "Running Time:\n"+botDuration(elapsed))
			} else if cmd == "set" || cmd == Data.Command.Set && Data.Command.Set != "" {
				elapsed := time.Since(botStart)
				Lilgo := fmt.Sprintf("%v", Kickbatas)
				Sys := fmt.Sprintf("%s", runtime.GOOS)
				ret := "Setting Bots:"
				ret += "\n"
				if Data.AntiTag {
					ret += "\n 🟢 Antitag"
				} else {
					ret += "\n 🔴 Antitag"
				}
				if Data.AutoPro {
					ret += "\n 🟢 Autopro"
				} else {
					ret += "\n 🔴 Autopro"
				}
				if Data.AutoPurge {
					ret += "\n 🟢 Autopurge"
				} else {
					ret += "\n 🔴 Autopurge"
				}
				if Data.ForceInvite {
					ret += "\n 🟢 Forceinvite"
				} else {
					ret += "\n 🔴 Forceinvite"
				}
				if Data.ForceJoinqr {
					ret += "\n 🟢 Forceqr"
				} else {
					ret += "\n 🔴 Forceqr"
				}
				if Data.Identict {
					ret += "\n 🟢 Identict"
				} else {
					ret += "\n 🔴 Identict"
				}
				if Data.NukeJoin {
					ret += "\n 🟢 Nukejoin"
				} else {
					ret += "\n 🔴 Nukejoin"
				}
				if Data.VictimMode {
					ret += "\n 🟢 Singlemode"
				} else {
					ret += "\n 🔴 Singlemode"
				}
				if Data.FastMode {
					ret += "\n 🟢 Fastmode"
				} else {
					ret += "\n 🔴 Fastmode"
				}
				if Data.FastMode || Data.VictimMode || Data.QrMode || Data.MixMode || Data.KillMode {
					ret += "\n 🔴 Warmode"
				} else {
					ret += "\n 🟢 Warmode"
				}
				if Data.MixMode {
					ret += "\n 🟢 Mixmode"
				} else {
					ret += "\n 🔴 Mixmode"
				}
				if Data.KillMode {
					ret += "\n 🟢 Killmode"
				} else {
					ret += "\n 🔴 Killmode"
				}
				if Data.QrMode {
					ret += "\n 🟢 Qrmode"
				} else {
					ret += "\n 🔴 Qrmode"
				}
				ret += "\n"
				ret += "\n↳System: " + Sys
				ret += "\n↳Limiterset: "+Lilgo
				ret += "\n● Active: " + botDuration(elapsed)
				ret += "\n\nᵖᵒʷᵉʳᵉᵈ ᵇʸ\n	Friendly™ ᴮᴼᵀ"
				client.SendMessageFooter(to, ret, IconLink, IconFooter, Data.Message.Flag)
			} else if cmd == "setting" || cmd == Data.Command.Setting && Data.Command.Setting != "" {
				ver := fmt.Sprintf("%s", runtime.Version())
				targets := []string{client.MID}
				for _, x := range client.Squads {
					if IsMembers(client, to, x) == true {
						targets = append(targets, x)
					}
				}
				ret := "Setting Group:"
				if (op.Message.ToType).String() == "GROUP" {
					ret += "\n"
					ret += "\n✿ Protect Bots:"
					ret += "\n"
					if InArray(Data.ProQr, to) {
						ret += "\n	🟢 Pro QR"
					} else {
						ret += "\n	🔴 Pro QR"
					}
					if InArray(Data.ProKick, to) {
						ret += "\n	🟢 Pro Kick"
					} else {
						ret += "\n	🔴 Pro Kick"
					}
					if InArray(Data.ProInvite, to) {
						ret += "\n	🟢 Pro Invite"
					} else {
						ret += "\n	🔴 Pro Invite"
					}
					if InArray(Data.ProCancel, to) {
						ret += "\n	🟢 Pro Cancel"
					} else {
						ret += "\n	🔴 ProCancel"
					}
					if InArray(Data.ProJoin, to) {
						ret += "\n	🟢 Pro Join"
					} else {
						ret += "\n	🔴 Pro Join"
					}
					if Data.ProName[to] == 1 {
						ret += "\n	🟢 Pro Name"
					} else {
						ret += "\n	🔴 Pro Name"
					}
					if len(Data.StayAjs[to]) != 0 {
						for _, v := range Data.StayAjs[to] {
							if IsPending(client, to, v) == true {
								ret += "\n	🟢 Pro Ajs"
							}else{
								ret += "\n	🔴 Pro Ajs"
							}
						}
					}else{ret += "\n	🔴 Pro Ajs"}
					ret += "\n"
					ret += "\n✿ General Bots:"
					ret += "\n"
					if siderV2[to] == true {
						ret += "\n	🟢 Lurking"
					} else {
						ret += "\n	🔴 Lurking"
					}
					if welcome[to] == 1 {
						ret += "\n	🟢 Welcome"
					} else {
						ret += "\n	🔴 Welcome"
					}
					if gcControlV2 == to {
						ret += "\n	🟢 Logmode"
					} else {
						ret += "\n	🔴 Logmode"
					}
				}
				ret += "\n"
				ret += "\n↳Bots here: "+strconv.Itoa(len(targets))+"/"+strconv.Itoa(len(client.Squads)+1)
				ret += "\n↳Version: " + ver
				ret += "\n● Timeleft: - "+strconv.Itoa(duedatecount)+" days"
				ret += "\n\nᵖᵒʷᵉʳᵉᵈ ᵇʸ\n	Friendly™ ᴮᴼᵀ"
				client.SendMessageFooter(to, ret, IconLink, IconFooter, Data.Message.Flag)
			} else if cmd == "sider on" {
				sider[to] = []string{}
				siderV2[to] = true
				client.SendMessage(to, "group sider activated")
			} else if cmd == "sider off" {
				delete(sider, to)
				siderV2[to] = false
				client.SendMessage(to, "group sider deactivated")
			}else if cmd == "welcome on"{
				welcome[to] = 1
				client.SendMessage(to,"welcome is activated")
			}else if cmd == "welcome off"{
				delete(welcome,to)
				client.SendMessage(to,"welcome is deactivated")
			} else if cmd == "speed" || cmd == Data.Command.Speed && Data.Command.Speed != "" {
				X, _ := client.GetGroupV2(to)
				memb := X.MemberMids
				var a = 0
				ret := "Speed Bot:"
				ret += "\n"
				for i := range ClientBot {
					if InArray(memb, ClientBot[i].MID) {
						a = a + 1
						start := time.Now()
						ClientBot[i].GetProfile()
						elapsed := time.Since(start)
						stringTime := elapsed.String()
						asu := fmt.Sprintf("%v", start.Nanosecond())
						ClientBot[i].SendMessage(to, "Benchmark: "+asu[0:2])
						ret += fmt.Sprintf("\nAssist%v: %s", a, stringTime[0:4])
					}
				}
				client.SendMessage(to, ret)
			} else if cmd == "squads" {
				if len(client.Squads) > 0 {
					listsq := "#Squads:\n"
					for i := range client.Squads {
						listsq += "\n" + strconv.Itoa(i+1) + ". " + "@!"
						fmt.Sprintf("... %v", i)
					}
					client.SendPollMention(to, listsq, client.Squads)
				} else { client.SendMessage(to, "Squad is empty.") }
			} else if cmd == "status" || cmd == Data.Command.Status && Data.Command.Status != "" {
				X, _ := client.GetGroupV2(to)
				memb := X.MemberMids
				var a = 0
				ret := "Status Bot:"
				ret += "\n"
				for i := range ClientBot {
					if InArray(memb, ClientBot[i].MID) {
						a = a + 1
						ClientBot[i].NormalKickoutFromGroup(to, []string{"FuckYou"})
						if ClientBot[i].Limited == true {
							ret += fmt.Sprintf("\nAssist%v: %s", a, Data.Message.Limit)
						} else {
							ret += fmt.Sprintf("\nAssist%v: %s", a, Data.Message.Fresh)
						}
					}
				}
				client.SendMessage(to, ret)
			} else if cmd == "status all" {
				ret := "Status Allbot:"
				ret += "\n"
				for i := range ClientBot {
					ClientBot[i].NormalKickoutFromGroup(to, []string{"FuckYou"})
					if ClientBot[i].Limited == true {
						ret += fmt.Sprintf("\nAssist%v: %s", i+1, Data.Message.Limit)
					} else {
						ret += fmt.Sprintf("\nAssist%v: %s", i+1, Data.Message.Fresh)
					}
				}
				client.SendMessage(to, ret)
			} else if cmd == "status add" {
				ret := "Status Add:"
				ret += "\n"
				for i := range ClientBot {
					ve := "u011f72e941cd24305e133d24ae8c6ada"
					_, err := ClientBot[i].FindAndAddContactsByMid(ve)
					fff := fmt.Sprintf("%v",err)
					er := strings.Contains(fff, "request blocked")
					if er == true {
						ret += fmt.Sprintf("\nAssist%v: %s", i+1, Data.Message.Limit)
					}else {
						ret += fmt.Sprintf("\nAssist%v: %s", i+1, Data.Message.Fresh)
					}
				}
				client.SendMessage(to, ret)
			} else if cmd == "set account"{
				X, _ := client.GetGroupV2(to)
				memb := X.MemberMids
				var a = 0
				ret := "Set Account:\n"
				for i := range ClientBot {
					if InArray(memb, ClientBot[i].MID) {
						cokk, _ := ClientBot[i].GetSettings()
						Ka := fmt.Sprintf("\n")
						Ki := fmt.Sprintf("\n")
						Ku := fmt.Sprintf("\n")
						Ke := fmt.Sprintf("\n")
						if cokk.PrivacyReceiveMessagesFromNotFriend == true {
							Ka += fmt.Sprintf("   ✓   Filter")
						} else {
							Ka += fmt.Sprintf("   ✘   Filter")
						}
						if cokk.EmailConfirmationStatus == 3 {
							Ki += fmt.Sprintf("   ✓   Email")
						} else {
							Ki += fmt.Sprintf("   ✘   Email")
						}
						if cokk.E2eeEnable == true {
							Ku += fmt.Sprintf("   ✓   Lsealing")
						} else {
							Ku += fmt.Sprintf("   ✘   Lsealing")
						}
						if cokk.PrivacyAllowSecondaryDeviceLogin == true {
							Ke += fmt.Sprintf("   ✓   Secondary\n")
						} else {
							Ke += fmt.Sprintf("   ✘   Secondary\n")
						}
						a = a + 1
						ret += fmt.Sprintf("\nAssist%v: %s", a, Ka, Ki, Ku, Ke)
					}
				}
				client.SendMessage(to, ret)
			} else if cmd == "tagall" {
				gc, _ := client.GetGroupV2(to)
				target := gc.MemberMids
				targets := []string{}
				for i := range target {
					targets = append(targets, target[i])
				}
				client.SendPollMention(to, "Mentions member:\n", targets)
			} else if cmd == "tagpen" {
				gc, _ := client.GetGroupV2(to)
				target := gc.InviteeMids
				targets := []string{}
				if len(target) != 0 {
					for i:= range target{
						targets = append(targets,target[i])
					}
					client.SendPollMention(to,"Mentions Pending:\n",targets)
				}else{client.SendMessage(to, "Pending is empty.")}
			} else if strings.HasPrefix(cmd, "unbot ") {
				if !InArray(Data.Master, sender) || !InArray(Data.GroupOwn[to], sender) {
				result := strings.Split((cmd), " ")
				if result[1] != "0" {
					result2, err := strconv.Atoi(result[1])
					if err != nil { client.SendMessage(to, "invalid number.")
					} else {
						for i := range Data.Bot {
							if result2 > 0 && result2-1 < len(Data.Bot) {
								if i == result2-1 {
									kura := Data.Bot[i]
									Data.Bot = Remove(Data.Bot, kura)
									client.SendMention(to, "success delbot @!", []string{kura})
									SaveData()
									break
								}
							} else { 
								client.SendMessage(to, "out of range.") 
								break 
							}
						}
					}
				} else { client.SendMessage(to, "invalid range.") }
				}
			} else if strings.HasPrefix(cmd, "ungban ") {
				result := strings.Split((cmd), " ")
				if result[1] != "0" {
					result2, err := strconv.Atoi(result[1])
					if err != nil { client.SendMessage(to, "invalid number.")
					} else {
						for i := range Data.GroupBan[to] {
							if result2 > 0 && result2-1 < len(Data.GroupBan[to]) {
								if i == result2-1 {
									kura := Data.GroupBan[to][i]
									Data.GroupBan[to] = Remove(Data.GroupBan[to], kura)
									client.SendMention(to, "success delgban @!", []string{kura})
									SaveData()
									break
								}
							} else { 
								client.SendMessage(to, "out of range.") 
								break 
							}
						}
					}
				} else { client.SendMessage(to, "invalid range.") }
			} else if cmd == "unsend" || cmd == Data.Command.Unsend && Data.Command.Unsend != "" {
				for i := range ClientBot {
					ClientBot[i].UnsendChat(to)
				}
				//PROTECTS
			} else if cmd == "help pro" {
				res := "	✿ Protection ✿"
				res += "\n"
				for a, x := range helppro {
					res += fmt.Sprintf("\n%v› %s %s", a+1, Data.Setkey, x)
				}
				client.SendMessage(to, res)
			} else if cmd == "deny kick" {
				if !InArray(Data.ProKick, to) {
					Data.ProKick = append(Data.ProKick, to)
					client.SendMessage(to, "Deny kick enabled.")
					SaveData()
				} else { client.SendMessage(to, "Already enabled.") }
			} else if cmd == "allow kick" {
				if InArray(Data.ProKick, to) {
					Data.ProKick = Remove(Data.ProKick, to)
					client.SendMessage(to, "Allow kick enabled.")
					SaveData()
				} else { client.SendMessage(to, "Already enabled.") }
			} else if cmd == "deny link" {
				if !InArray(Data.ProQr, to) {
					Data.ProQr = append(Data.ProQr, to)
					client.SendMessage(to, "Deny link enabled.")
					SaveData()
				} else { client.SendMessage(to, "Already enabled.") }
			} else if cmd == "allow link" {
				if InArray(Data.ProQr, to) {
					Data.ProQr = Remove(Data.ProQr, to)
					client.SendMessage(to, "Allow link enabled.")
					SaveData()
				} else { client.SendMessage(to, "Already enabled.") }
			} else if cmd == "deny invite" {
				if !InArray(Data.ProInvite, to) {
					Data.ProInvite = append(Data.ProInvite, to)
					client.SendMessage(to, "Deny invite enabled.")
					SaveData()
				} else { client.SendMessage(to, "Already enabled.") }
			} else if cmd == "allow invite" {
				if InArray(Data.ProInvite, to) {
					Data.ProInvite = Remove(Data.ProInvite, to)
					client.SendMessage(to, "Allow invite enabled.")
					SaveData()
				} else { client.SendMessage(to, "Already enabled.") }
			} else if cmd == "deny cancel" {
				if !InArray(Data.ProCancel, to) {
					Data.ProCancel = append(Data.ProCancel, to)
					client.SendMessage(to, "Deny cancel enabled.")
					SaveData()
				} else { client.SendMessage(to, "Already enabled.") }
			} else if cmd == "allow cancel" {
				if InArray(Data.ProCancel, to) {
					Data.ProCancel = Remove(Data.ProCancel, to)
					client.SendMessage(to, "Allow cancel enabled.")
					SaveData()
				} else { client.SendMessage(to, "Already enabled.") }
			} else if cmd == "deny join" {
				if !InArray(Data.ProJoin, to) {
					Data.ProJoin = append(Data.ProJoin, to)
					client.SendMessage(to, "Deny join enabled.")
					SaveData()
				} else { client.SendMessage(to, "Already enabled.") }
			} else if cmd == "allow join" {
				if InArray(Data.ProJoin, to) {
					Data.ProJoin = Remove(Data.ProJoin, to)
					client.SendMessage(to, "Allow join enabled.")
					SaveData()
				} else { client.SendMessage(to, "Already enabled.") }
			} else if cmd == "deny name" {
				if Data.ProName[to] != 1 {
					Data.ProName[to] = 1
					G := client.GetCompactGroup(to)
					Data.GroupName[to] = G.Name
					client.SendMessage(to, "Deny name enabled.")
					SaveData()
				} else { client.SendMessage(to, "Already enabled.") }
			} else if cmd == "allow name" {
				if Data.ProName[to] == 1 {
					delete(Data.GroupName,to)
					delete(Data.ProName,to)
					client.SendMessage(to, "Allow name enabled.")
					SaveData()
				} else { client.SendMessage(to, "Already enabled.") }
			} else if strings.HasPrefix(cmd, "protect ") {
				result := strings.Split((cmd), " ")
				spl := strings.Replace(cmd, result[0]+" ", "", 1)
				if spl == "max" {
					if !InArray(Data.ProKick, to) {
						Data.ProKick = append(Data.ProKick, to)
					}
					if !InArray(Data.ProQr, to) {
						Data.ProQr = append(Data.ProQr, to)
					}
					if !InArray(Data.ProInvite, to) {
						Data.ProInvite = append(Data.ProInvite, to)
					}
					if !InArray(Data.ProCancel, to) {
						Data.ProCancel = append(Data.ProCancel, to)
					}
					SaveData()
					client.SendMessage(to, "Max protect enabled.")
				} else if spl == "low" {
					if InArray(Data.ProKick, to) {
						Data.ProKick = Remove(Data.ProKick, to)
					}
					if InArray(Data.ProQr, to) {
						Data.ProQr = Remove(Data.ProQr, to)
					}
					if InArray(Data.ProInvite, to) {
						Data.ProInvite = Remove(Data.ProInvite, to)
					}
					if InArray(Data.ProCancel, to) {
						Data.ProCancel = Remove(Data.ProCancel, to)
					}
					SaveData()
					client.SendMessage(to, "Max protect disabled.")
				}
				SaveData()
			}
		}
	}
}
func SelfBots(client *LINE) {
	runtime.GOMAXPROCS(cpu)
	for {
		multiFunc := client.fetchOps(client.Revision, client.Count, client.GRevision, client.IRevision)
		client.Revision = -1
		go func(fetch []*core.Operation) {
			for _, op := range fetch {
				client.SetRevision(op.Revision)
				var param1, param2, param3 = op.Param1, op.Param2, op.Param3
				if op.Type == 13 || op.Type == 124 {
					params3 := strings.Split(param3, "\x1e")
					if InArray(params3, client.MID) && fullAccess(client, param2) {
						go func() { client.AcceptChatInvitation(param1) }()
					} else if !fullAccess(client, param2) {
						Check = param2
						if _, cek := KillMode[param2]; !cek {
							KillMode[param2] = []string{}
						}
					}
				} else if op.Type == 17 || op.Type == 130 {
					if !fullAccess(client, param2) && !IsGaccess(param1, param2) {
						if Detect == 0 {
							Detect = 1
							JoinFrequence[param1] = time.Now()
							Detect = 0
						}
						if botStart.Sub(JoinFrequence[param1]) <= 500*time.Millisecond {
							if !fullAccess(client, param2) {
								if !InArray(KillMode[Check], param2) {
									KillMode[Check] = append(KillMode[Check], param2)
								}
							}
						}
					}
				} else if op.Type == 25 {
					var msg, sender, text, receiver = op.Message, op.Message.From_, op.Message.Text, op.Message.To
					var to string
					var coms string
					if (msg.ToType).String() == "USER" {
						to = sender
					} else {
						to = receiver
					}
					Setkey := Data.Prefix
					if (msg.ContentType).String() == "NONE" {
						if strings.HasPrefix(text, Setkey) {
							coms = strings.Replace(text, Setkey, "", 1)
						}
						if strings.HasPrefix(text, Setkey+" ") {
							coms = strings.Replace(text, Setkey+" ", "", 1)
						}
						for _, cmd := range strings.Split(coms, ", ") {
							if strings.HasPrefix(cmd, "selfbot ") {
								spl := strings.Replace(cmd, "selfbot ", "", 1)
								if spl == "on" {
									Data.SelfStatus = true
									SaveData()
									client.SendMessage(to, "Selfbot enabled")
								} else if spl == "off" {
									Data.SelfStatus = false
									SaveData()
									client.SendMessage(to, "Selfbot disabled")
								}
							}
							if Data.SelfStatus {
								if cmd == "me" {
									client.SendContact(to, client.MID)
								} else if cmd == "kill" {
									if len(KillMode) > 0 {
										listsq := "#KillMode:\n"
										var no = 0
										for i := range KillMode {
											if i != "" {
												con := client.GetContact(i)
												listsq += fmt.Sprintf("\n%v. %s", no+1, con.DisplayName)
												for _, a := range KillMode[i] {
													conn := client.GetContact(a)
													listsq += fmt.Sprintf("\n    - %s", conn.DisplayName)
												}
												no++
											}
										}
										client.SendMessage(to, listsq)
									}
								} else if cmd == "mymid" {
									client.SendMessage(to, client.MID)
								} else if cmd == "speed" {
									start := time.Now()
									client.GetProfile()
									elapsed := time.Since(start)
									stringTime := elapsed.String()
									client.SendMessage(to, "Myspeed: "+stringTime[0:4]+" ms")
								} else if strings.HasPrefix(cmd, "reinv") {
									mentions := mentions{}
									json.Unmarshal([]byte(op.Message.ContentMetadata["MENTION"]), &mentions)
									for _, mention := range mentions.MENTIONEES {
										if IsFriends(client, mention.Mid) == false {
											client.FindAndAddContactsByMid(mention.Mid)
											client.InviteIntoChat(to, []string{mention.Mid})
										}
									}
								} else if strings.HasPrefix(cmd, "kick") {
									targets := []string{}
									mentions := mentions{}
									json.Unmarshal([]byte(op.Message.ContentMetadata["MENTION"]), &mentions)
									for _, mention := range mentions.MENTIONEES {
										targets = append(targets, mention.Mid)
										if !fullAccess(client, mention.Mid) {
											appendBl(mention.Mid)
										}
									}
									var wg sync.WaitGroup
									wg.Add(len(targets))
									for i := 0; i < len(targets); i++ {
										go func(i int) {
											defer wg.Done()
											client.DeleteOtherFromChat(to, []string{targets[i]})
										}(i)
									}
									wg.Wait()
								} else if strings.HasPrefix(cmd, "invite ") {
									str := strings.Replace(cmd, "invite ", "", 1)
									count := strings.Split(str, " ")
									list, _ := strconv.Atoi(count[0])
									if list > 0 && list <= len(ClientBot) {
										numb := list - 1
										grup, _ := client.GetGroupV2(to)
										member := grup.MemberMids
										if !InArray(member, ClientBot[numb].MID) {
											if IsFriends(client, ClientBot[numb].MID) == false {
												client.FindAndAddContactsByMid(ClientBot[numb].MID)
												client.InviteIntoChat(to, []string{ClientBot[numb].MID})
											}
										} else {
											client.SendMessage(to, "Already in groups.")
										}
									} else {
										client.SendMessage(to, "Out of range.")
									}
								} else if cmd == "here" {
									targets := []string{}
									for _, x := range Data.SquadBots {
										if IsMembers(client, to, x) == true {
											targets = append(targets, x)
										}
									}
									client.SendMessage(to, fmt.Sprintf("%v/%v bots here.", len(targets), len(Data.SquadBots)))
								} else if cmd == "byeall" {
									go func() {
										for i := range ClientBot {
											ClientBot[i].LeaveGroup(to)
										}
									}()
								} else if cmd == "clearban" {
									if len(Data.Blacklist) == 0 {
										client.SendMessage(to, "Banlist is empty.")
									} else {
										for _, x := range Data.Blacklist {
											if x == "" {
												Data.Blacklist = Remove(Data.Blacklist, x)
												SaveData()
											}
										}
										msgcbn := fmt.Sprintf(Data.Message.Ban, len(Data.Blacklist))
										client.SendMessage(to, msgcbn)
										Data.Blacklist = []string{}
										SaveData()
									}
									Data.Blacklist = []string{}
									SaveData()
								} else if cmd == "cban" {
									go func(to string) {
										if len(Data.Blacklist) > 0 {
											for _, x := range Data.Blacklist {
												if x == "" {
													Data.Blacklist = Remove(Data.Blacklist, x)
													SaveData()
												}
											}
											client.SendPollMention(to, "#Banlist:\n", Data.Blacklist)
											Data.Blacklist = []string{}
											time.Sleep(1 * time.Second)
											Data.Blacklist = []string{}
											SaveData()
										} else {
											client.SendMessage(to, "Banlist is empty.")
											Data.Blacklist = []string{}
											SaveData()
										}
									}(to)
								} else if strings.HasPrefix(cmd, "stand ") {
									str := strings.Replace(cmd, "stand ", "", 1)
									result2, _ := strconv.Atoi(str)
									if result2 > 0 && result2 <= len(Data.SquadBots) {
										grup, _ := client.GetGroupV2(to)
										target := grup.MemberMids
										targets := []string{}
										tempInv := []string{}
										batastim := 0
										batastim2 := 0
										for i := range ClientBot {
											ClientBot[i].NormalKickoutFromGroup(to, []string{"FuckYou"})
											if ClientBot[i].Limited {
												ClientBot[i].LeaveGroup(to)
											}
											_, found := Data.StayGroup[to]
											if found == true {
												Data.StayGroup[to] = []string{}
											}
										}
										for i := range target {
											targets = append(targets, target[i])
										}
										_, found := Data.StayGroup[to]
										if found == false {
											for i := range Data.SquadBots {
												if InArray(targets, Data.SquadBots[i]) {
													Data.StayGroup[to] = append(Data.StayGroup[to], Data.SquadBots[i])
												}
											}
										}
										for i := range targets {
											if InArray(Data.StayGroup[to], targets[i]) {
												if batastim < result2 {
													batastim = batastim + 1
												} else {
													ClientMid[targets[i]].LeaveGroup(to)
													Data.StayGroup[to] = Remove(Data.StayGroup[to], targets[i])
												}
											}
										}
										for io := range Data.SquadBots {
											if InArray(targets, Data.SquadBots[io]) {
												batastim2 = batastim2 + 1
											}
										}
										for i := range Data.SquadBots {
											if batastim2 < result2 {
												if !InArray(targets, Data.SquadBots[i]) {
													if ClientMid[Data.SquadBots[i]].Limited == false {
														batastim2 = batastim2 + 1
														tempInv = append(tempInv, Data.SquadBots[i])
														if !InArray(Data.StayGroup[to], Data.SquadBots[i]) {
															Data.StayGroup[to] = append(Data.StayGroup[to], Data.SquadBots[i])
														}
													}
												}
											}
										}
										if len(tempInv) != 0 {
											client.InviteIntoChat(to, tempInv)
										}
										SaveData()
									} else {
										client.SendMessage(to, "Out of range.")
									}
								} else if strings.HasPrefix(cmd, "stay ") {
									str := strings.Replace(cmd, "stay ", "", 1)
									ticket := client.ReissueGroupTicket(to)
									result2, _ := strconv.Atoi(str)
									if result2 > 0 && result2 <= len(Data.SquadBots) {
										getmem, _ := client.GetGroupV2(to)
										target := getmem.MemberMids
										targets := []string{}
										batastim := 0
										batastim2 := 0
										for i := range ClientBot {
											ClientBot[i].NormalKickoutFromGroup(to, []string{"FuckYou"})
											if ClientBot[i].Limited {
												ClientBot[i].LeaveGroup(to)
											}
											_, found := Data.StayGroup[to]
											if found == true {
												Data.StayGroup[to] = []string{}
											}
										}
										for i := range target {
											targets = append(targets, target[i])
										}
										_, found := Data.StayGroup[to]
										if found == false {
											for i := range Data.SquadBots {
												if InArray(targets, Data.SquadBots[i]) {
													Data.StayGroup[to] = append(Data.StayGroup[to], Data.SquadBots[i])
												}
											}
										}
										for i := range targets {
											if InArray(Data.StayGroup[to], targets[i]) {
												if batastim < result2 {
													batastim = batastim + 1
												} else {
													ClientMid[targets[i]].LeaveGroup(to)
													Data.StayGroup[to] = Remove(Data.StayGroup[to], targets[i])
												}
											}
										}
										for io := range Data.SquadBots {
											if InArray(targets, Data.SquadBots[io]) {
												batastim2 = batastim2 + 1
											}
										}
										if batastim2 < result2 {
											if getmem.PreventedJoinByTicket == true {
												getmem.PreventedJoinByTicket = false
												client.UpdateGroup(getmem)
											}
										}
										for i := range ClientBot {
											if batastim2 < result2 {
												if !InArray(targets, Data.SquadBots[i]) {
													if ClientMid[ClientBot[i].MID].Limited == false {
														err := ClientBot[i].AcceptGroupByTicket(to, ticket)
														if err == nil { batastim2 = batastim2 + 1 }
														if !InArray(Data.StayGroup[to], Data.SquadBots[i]) {
															Data.StayGroup[to] = append(Data.StayGroup[to], Data.SquadBots[i])
														}
													}
												}
											}
										}
										if batastim2 == result2 {
											if getmem.PreventedJoinByTicket == false {
												getmem.PreventedJoinByTicket = true
												client.UpdateGroup(getmem)
											}
										}
										SaveData()
									} else {
										client.SendMessage(to, "Out of range.")
									}
								}
							}
						}
					}
				} else {
					client.CorrectRevision(op)
				}
			}
		}(multiFunc)
		for _, ops := range multiFunc {
			if ops.Revision != -1 {
				client.SetRevision(ops.Revision)
			} else {
				client.CorrectRevision(ops)
			}
		}
	}
}

/**************
RUNNING BOTS
***************/
func main() {
	cpu = runtime.NumCPU()
	jsonFile, err := os.Open(DATABASE)
	if err != nil { fmt.Println(err) }
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &Data)
	Data.SquadBots = []string{}
	Data.SelfStatus = false
	SaveData()
	fmt.Println("\n_Started Login_:")
	if Data.SelfToken != "" {
		xy := randomToString(4)
		var app = fmt.Sprintf("DESKTOPWIN\t7.3.1\tWindows\t10.%v", xy)
		var ua = fmt.Sprintf("Line/7.3.1 Windows 10.%v", xy)
		erx := me.CreateNewLogin(Data.SelfToken, app, ua, HostName[0])
		if erx == nil {
			go func() { SelfBots(me) }()
			Data.SelfStatus = true
			SaveData()
			prof := me.GetProfile()
			if me.MID != MAKERS { me.FindAndAddContactsByMid(MAKERS) }
			fmt.Println("\n\n  ↳ DisplayName : " + prof.DisplayName + "\n  ↳ Mid : " + me.MID + "\n  ↳ AppName : " + me.AppName + "\n  ↳ UserAgent : " + me.UserAgent + "\n  ↳ Type: SelfBots")
		} else {
			logs := fmt.Sprintf("\n\n▪︎ ERROR: %s", erx)
			fmt.Println(logs)
		}
	}
	for no, tok := range Data.Authoken {
		time.Sleep(250 * time.Millisecond)
		cl := StartLogin()
		xy := randomToString(4)
		var app = fmt.Sprintf("ANDROIDLITE\t2.17.1\tAndroid OS\t10.%v", xy)
		var ua = fmt.Sprintf("LLA/2.17.1 64C0D3 10.%v", xy)
		err := cl.CreateNewLogin(tok, app, ua, HostName[no+1])
		if err == nil {
			prof := cl.GetProfile()
			go func() { BackupBots(cl) }()
			if IsFriends(cl, MAKERS) == false { cl.FindAndAddContactsByMid(MAKERS) }
			fmt.Println("\n\n  ↳ DisplayName : " + prof.DisplayName + "\n  ↳ Mid : " + cl.MID + "\n  ↳ AppName : " + cl.AppName + "\n  ↳ UserAgent : " + cl.UserAgent + "\n  ↳ Bots No: " + fmt.Sprintf("%v", no+1))
			Data.SquadBots = append(Data.SquadBots, cl.MID)
			Data.Message.Welcome = "oh hi, welcome @!"
			Data.Message.Sider = "oh hi, i see u @!"
			ClientBot = append(ClientBot, cl)
			ClientMid[cl.MID] = cl
			Data.Rname = "bot"
			cl.CheckExprd()
		} else {
			logs := fmt.Sprintf("\n\n▪︎ No: %v ERROR: %s", no+1, err)
			fmt.Println(logs)
		}
	}
	SaveData()
	fmt.Println("\n\n___All Login Success___\n")
	if len(ClientBot) != 0 {
		ClientBot[0].SendMessage(MAKERS, "Im fetcher.")
		for i := range ClientBot {
			ClientBot[i].Backup = Data.SquadBots
			for _, x := range Data.SquadBots {
				if !InArray(ClientBot[i].Squads, x) && x != ClientBot[i].MID {
					ClientBot[i].Squads = append(ClientBot[i].Squads, x)
				}
			}
		}
		for i := range ClientBot {
			group := ClientBot[i].GetGroupsJoined()
			for _, x := range group {
				grup, _ := ClientBot[i].GetGroupV2(x)
				memb := grup.MemberMids
				for _, a := range Data.SquadBots {
					if InArray(memb, a) && !InArray(Data.StayGroup[x], a) {
						Data.StayGroup[x] = append(Data.StayGroup[x], a)
					}
				}
			}
		}
	}
	SaveData()
	ch := make(chan int, len(Data.Authoken))
	for v := range ch {
		if v == 51 {
			break
		}
	}
	fmt.Println("__GOOD_LUCK__")
}