package main

import (
	"flag"
	"log"

	"github.com/liut/staffio/pkg/backends"
	"github.com/liut/staffio/pkg/settings"
	"github.com/wealthworks/go-tencent-api/exmail"
	"github.com/wealthworks/go-tencent-api/exwechat"
)

var (
	uid = flag.String("uid", "", "uid: username | email")
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	settings.Parse()
	backends.Prepare()
	if *uid == "" {
		flag.PrintDefaults()
		return
	}

	log.Printf("fetching user %s", *uid)

	staff, err := backends.GetStaff(*uid)
	if err != nil {
		log.Printf("get staff err: %s", err)
		return
	}
	log.Printf("staff: %v", staff)
	if staff == nil {
		return
	}

	wechat := exwechat.NewAPI()
	user, err := wechat.GetUser(*uid)
	if err != nil {
		log.Print(err)
		return
	}
	log.Printf("WxWork User: %v", user)

	alias := backends.GetEmailAddress(*uid)

	count, err := exmail.CountNewMail(alias)
	if err != nil {
		log.Print(err)
		return
	}
	log.Printf("new mail: %d", count)

	url, err := exmail.GetLoginURL(alias)
	if err != nil {
		log.Print(err)
		return
	}
	log.Printf("open %s", url)
}
