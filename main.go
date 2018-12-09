package main

import (
	"fmt"
	"ssher"
	"strings"
	"testing"
	"time"

	"github.com/sclevine/agouti"
)

const (
	NULL            = "\uE000"
	CANCEL          = "\uE001"
	HELP            = "\uE002"
	BACK_SPACE      = "\uE003"
	TAB             = "\uE004"
	CLEAR           = "\uE005"
	RETURN          = "\uE006"
	ENTER           = "\uE007"
	SHIFT           = "\uE008"
	LEFT_SHIFT      = "\uE008"
	CONTROL         = "\uE009"
	LEFT_CONTROL    = "\uE009"
	ALT             = "\uE00A"
	LEFT_ALT        = "\uE00A"
	PAUSE           = "\uE00B"
	ESCAPE          = "\uE00C"
	SPACE           = "\uE00D"
	PAGE_UP         = "\uE00E"
	PAGE_DOWN       = "\uE00F"
	END             = "\uE010"
	HOME            = "\uE011"
	LEFT            = "\uE012"
	ARROW_LEFT      = "\uE012"
	UP              = "\uE013"
	ARROW_UP        = "\uE013"
	RIGHT           = "\uE014"
	ARROW_RIGHT     = "\uE014"
	DOWN            = "\uE015"
	ARROW_DOWN      = "\uE015"
	INSERT          = "\uE016"
	DELETE          = "\uE017"
	SEMICOLON       = "\uE018"
	EQUALS          = "\uE019"
	NUMPAD0         = "\uE01A"
	NUMPAD1         = "\uE01B"
	NUMPAD2         = "\uE01C"
	NUMPAD3         = "\uE01D"
	NUMPAD4         = "\uE01E"
	NUMPAD5         = "\uE01F"
	NUMPAD6         = "\uE020"
	NUMPAD7         = "\uE021"
	NUMPAD8         = "\uE022"
	NUMPAD9         = "\uE023"
	MULTIPLY        = "\uE024"
	ADD             = "\uE025"
	SEPARATOR       = "\uE026"
	SUBTRACT        = "\uE027"
	DECIMAL         = "\uE028"
	DIVIDE          = "\uE029"
	F1              = "\uE031"
	F2              = "\uE032"
	F3              = "\uE033"
	F4              = "\uE034"
	F5              = "\uE035"
	F6              = "\uE036"
	F7              = "\uE037"
	F8              = "\uE038"
	F9              = "\uE039"
	F10             = "\uE03A"
	F11             = "\uE03B"
	F12             = "\uE03C"
	META            = "\uE03D"
	COMMAND         = "\uE03D"
	ZENKAKU_HANKAKU = "\uE040"
)

func ErrorChecker(err error, place string) {
	if err != nil {
		fmt.Printf("test died at\t"+place+"\n%s\n", err)
	}
}

func PageRefresher(linku string, Driver *agouti.WebDriver) *agouti.Page {
	var t *testing.T
	page, err := Driver.NewPage(agouti.Browser("chrome"))
	if err != nil {
		t.Fatal("Failed to open page:", err)
	}

	if err := page.Navigate(linku); err != nil {
		t.Fatal("Failed to navigate:", err)
	}
	return page
}

func Clicker(button string, page *agouti.Page) (*agouti.Selection, error) {
	time.Sleep(2 * time.Second)
	element := page.FindByXPath(button)
	err := element.Click()
	return element, err
}

func CrmPrimitiveChecker(ip string) {
	b := ssher.SSH("root", ip, "echo `crm configure show`", "default")
	lines := strings.Split(b, "\\")
	for i := 0; i < len(lines); i++ {
		if strings.Contains(lines[i], "cool_primitive anything") {
			if strings.Contains(lines[i+1], "binfile=file") {
				fmt.Println("Binfile set correctly [for cool_primitive] -- PASSED")
			}
			if strings.Contains(lines[i+2], "start timeout=35s") {
				fmt.Println("start set correctly  --  PASSED")
			}
			if strings.Contains(lines[i+3], "stop timeout=15s on-fail=stop") {
				fmt.Println("stop set correctly  --  PASSED")
			}
			if strings.Contains(lines[i+4], "monitor timeout=9s interval=13s") {
				fmt.Println("monitor set correctly  --  PASSED")
			}
			if strings.Contains(lines[i+5], "meta target-role=Started") {
				fmt.Println("target-role set correctly  --  PASSED")
			}
		}
	}
}

//--- logging into the ha-cluster
func Login(linku string, page *agouti.Page) {
	element, err := Clicker("//*[@id=\"session_username\"]", page)
	place := "user login"
	ErrorChecker(err, place)

	err = element.Fill("hacluster")
	place = "typing user name"
	ErrorChecker(err, place)

	element, err = Clicker("//*[@id=\"session_password\"]", page)
	place = "password login"
	ErrorChecker(err, place)

	err = element.Fill("linux")
	place = "typing password"
	ErrorChecker(err, place)

	element, err = Clicker("//*[@id=\"new_session\"]/input[3]", page)
	place = "clicking \"LOGIN\" "
	ErrorChecker(err, place)
}

func Cluster_Troubler(linku string, page *agouti.Page, ip string) {

	//--- setting stonith-sbd maintenance state ON/OFF
	_, err := Clicker("//*[@id=\"resources\"]/div[1]/div[2]/div[2]/table/tbody/tr/td[6]/div/div", page)
	place := "Clicking cascade menu next to stonith-sbd resource"
	ErrorChecker(err, place)

	_, err = Clicker("//*[@id=\"resources\"]/div[1]/div[2]/div[2]/table/tbody/tr/td[6]/div/div/ul/li[1]/a", page)
	place = "Clicking maintenance mode in the cascade menu of stonit-sbd"
	ErrorChecker(err, place)

	time.Sleep(5 * time.Second)
	_, err = Clicker("//*[@id=\"confirmationDialog\"]/div/div/form/div[3]/button[2]", page)
	place = "Clicking \"OK\" in the pop-up menu window"
	ErrorChecker(err, place)

	//--- waiting few seconds for stonith to enter maintenance mode
	time.Sleep(4 * time.Second)

	//ssher.SSH("root", ip, "crm status") checking if stonith is "unmanaged"
	if ssher.SSH("root", ip, "crm status | grep -i stonith | grep -i unmanaged", "default") != "" {
		fmt.Println("stonith switched to maintenance mode --   PASSED")
	}

	time.Sleep(5 * time.Second)

	//_, err = Clicker("//*[@id=\"resources\"]/div[1]/div[2]/div[2]/table/tbody/tr/td[6]/div/div/button", page)
	//place = "Clicking on sbd properties cascade menu"
	//ErrorChecker(err, place)

	_, err = Clicker("//*[@id=\"resources\"]/div[1]/div[2]/div[2]/table/tbody/tr/td[6]/div/a[1]", page)
	place = "Clicking maintenance trigger next to stonith [to switch maintenance OFF]"
	ErrorChecker(err, place)

	_, err = Clicker("//*[@id=\"confirmationDialog\"]/div/div/form/div[3]/button[2]", page)
	place = "Clicking \"OK\" in the pop-up menu window"
	ErrorChecker(err, place)

	time.Sleep(7 * time.Second)

	//--- cleaning state of stonith-sbd
	_, err = Clicker("//*[@id=\"resources\"]/div[1]/div[2]/div[2]/table/tbody/tr/td[6]/div/div/button", page)
	place = "Clicking the cascade menu for stonith-sbd resource"
	ErrorChecker(err, place)

	//*[@id="resources"]/div[1]/div[2]/div[2]/table/tbody/tr/td[6]/div/div/button

	_, err = Clicker("//*[@id=\"resources\"]/div[1]/div[2]/div[2]/table/tbody/tr/td[6]/div/div/button", page)
	place = "Selecting \"Clear State\" [of the stonith] in the cascade menu"
	ErrorChecker(err, place)

	time.Sleep(5 * time.Second) //-- waiting, because cleaning stonith-sbd takes some seconds

	//--- navigating to "Nodes" section
	_, err = Clicker("//*[@id=\"middle\"]/div[2]/div[2]/div/div[1]/ul/li[2]/a", page)
	place = "Clicking \"Nodes\" Section"
	ErrorChecker(err, place)

	//--- checking the first listed node's details
	_, err = Clicker("//*[@id=\"nodes\"]/div[1]/div[2]/div[2]/table/tbody/tr[1]/td[5]/div/a[2]", page)
	place = "Clicking 1st listed node's details"
	ErrorChecker(err, place)

	time.Sleep(7 * time.Second) // -- checking first node details opens another windon, takes time

	_, err = Clicker("//*[@id=\"modal\"]/div/div/div[3]/button", page)
	place = "Pressing \"OK\" on node info pop-up"
	ErrorChecker(err, place)

	//--- clearing state of first listed node
	_, err = Clicker("//*[@id=\"nodes\"]/div[1]/div[2]/div[2]/table/tbody/tr[1]/td[5]/div/div/button/i", page)
	place = "Pressing cascade menu next to 1st listed node"
	ErrorChecker(err, place)

	_, err = Clicker("//*[@id=\"nodes\"]/div[1]/div[2]/div[2]/table/tbody/tr[1]/td[5]/div/div/ul/li[2]/a", page)
	place = "Chosing \"Clear State\" option @node cascade menu"
	ErrorChecker(err, place)

	time.Sleep(3 * time.Second)

	_, err = Clicker("//*[@id=\"confirmationDialog\"]/div/div/form/div[3]/button[2]", page)
	place = "Pressing \"OK\" in the pop-up menu"
	ErrorChecker(err, place)

	time.Sleep(7 * time.Second)

	//--- setting first listed node to maintenance mode
	_, err = Clicker("//*[@id=\"nodes\"]/div[1]/div[2]/div[2]/table/tbody/tr[1]/td[3]/a", page)
	place = "Clicking the Maintenance trigger [to turn it ON]"
	ErrorChecker(err, place)

	_, err = Clicker("//*[@id=\"confirmationDialog\"]/div/div/form/div[3]/button[2]", page)
	place = "Clicking OK on the pop-up dialogue window"
	ErrorChecker(err, place)

	time.Sleep(5 * time.Second)

	if ssher.SSH("root", ip, "crm status | grep -i node | grep -i maintenance", "default") != "" {
		fmt.Println("1st listed node switched to maintenance mode --   PASSED")
	}

	//--- setting first listed node's maintenance off
	_, err = Clicker("//*[@id=\"nodes\"]/div[1]/div[2]/div[2]/table/tbody/tr[1]/td[3]/a/i", page)
	place = "Clicking the Maintenance trigger [to turn it OFF]"
	ErrorChecker(err, place)

	_, err = Clicker("//*[@id=\"confirmationDialog\"]/div/div/form/div[3]/button[2]", page)
	place = "Clicking OK on the pop-up dialogue window"
	ErrorChecker(err, place)

	//-----------SWITCHING TO DASHBOARD--------------
	_, err = Clicker("//*[@id=\"monitoringMenu\"]/li[2]/a", page)
	place = "Clicking on \"Dashboard\""
	ErrorChecker(err, place)

	//----------SWITCHING TO TROUBLESHOOTING---------
	_, err = Clicker("//*[@id=\"accordion\"]/li[3]/a", page)
	place = "Clicking on \"Troubleshooting\""
	ErrorChecker(err, place)

	//-------Clicking History ----------------------
	_, err = Clicker("//*[@id=\"troubleshootMenu\"]/li[1]/a", page)
	place = "Clicking \"History\""
	ErrorChecker(err, place)

	_, err = Clicker("//*[@id=\"generate\"]/form/div/div[2]", page)
	place = "Clicking \"generate\""
	ErrorChecker(err, place)
	time.Sleep(20 * time.Second)

	//---------Clicking Command Log--------------------
	_, err = Clicker("//*[@id=\"troubleshootMenu\"]/li[2]/a", page)
	place = "Clicking \"Command Log\""
	ErrorChecker(err, place)
	//-----------HERE!!!
	//---------CLICKING CONFIGURATION--------------------
	_, err = Clicker("//*[@id=\"accordion\"]/li[4]/a", page)
	place = "Clicking \"Configuration\""
	ErrorChecker(err, place)

	//----------ADDING A NEW PRIMITIVE...----------------------|
	//--------Adding Resource---------------------------
	_, err = Clicker("//*[@id=\"configurationMenu\"]/li[1]/a", page)
	place = "Clicking \"Add Resource\""
	ErrorChecker(err, place)

	//--------Selecting New Primitive----------------------
	_, err = Clicker("//*[@id=\"middle\"]/div[2]/div[2]/ul/li[1]/a", page)
	place = "Clicking \"New Primitive\""
	ErrorChecker(err, place)

	//---------Selecting Primitive Name Field--------------
	element, err := Clicker("//*[@id=\"primitive_id\"]", page)
	place = "Clicking on primitive ID field"
	ErrorChecker(err, place)

	err = element.Fill("cool_primitive")
	place = "Filling Primitive ID"
	ErrorChecker(err, place)

	_, err = Clicker("//*[@id=\"primitive_clazz\"]", page)
	place = "Clicking \"Primitive Class\" cascade"
	ErrorChecker(err, place)

	_, err = Clicker("//*[@id=\"primitive_clazz\"]/option[3]", page)
	place = "selecting ocf"
	ErrorChecker(err, place)

	_, err = Clicker("//*[@id=\"primitive_clazz\"]", page)
	place = "Clicking \"Primitive Class\" cascade"
	ErrorChecker(err, place)

	_, err = Clicker("//*[@id=\"primitive_type\"]", page)
	place = "Clicking \"Type\" cascade"
	ErrorChecker(err, place)

	_, err = Clicker("//*[@id=\"primitive_type\"]/option[2]", page)
	place = "Clicking \"Type\" cascade"
	ErrorChecker(err, place)

	_, err = Clicker("//*[@id=\"primitive_type\"]", page)
	place = "Clicking \"Type\" cascade"
	ErrorChecker(err, place)

	//--writing binfile name
	element, err = Clicker("//*[@id=\"binfile\"]", page)
	place = "Clicking \"binfile\" form"
	ErrorChecker(err, place)

	err = element.Fill("file")
	place = "Filling \"binfile\" form"
	ErrorChecker(err, place)

	//---setting start stop monitoring params
	//--- start part
	_, err = Clicker("//*[@id=\"oplist\"]/fieldset/div/div[1]/div[1]/div[2]/div/div/a[1]", page)
	place = "Clicking \"Edit\" at \"start\" param"
	ErrorChecker(err, place)

	element, err = Clicker("//*[@id=\"modal\"]/div/div/form/div[2]/fieldset/div/div[1]/div/div", page)
	place = "Clicking the time form [for start param of cool_primitive]"
	ErrorChecker(err, place)

	//#modal > div > div > form > div.modal-body > fieldset > div > div:nth-child(11) > div > div
	time.Sleep(3 * time.Second)

	err = page.Find("#timeout").SendKeys(BACK_SPACE + BACK_SPACE + BACK_SPACE)
	place = "3 * backspace"
	ErrorChecker(err, place)

	err = page.Find("#timeout").Fill("35s")
	place = "Filling 35 seconds for timeout"
	ErrorChecker(err, place)

	_, err = Clicker("//*[@id=\"modal\"]/div/div/form/div[3]/input", page)
	place = "Submitting start param [35s]"
	ErrorChecker(err, place)

	//--- stop part
	_, err = Clicker("//*[@id=\"oplist\"]/fieldset/div/div[1]/div[2]/div[2]/div/div/a[1]", page)
	place = "Clicking \"Edit\" at \"stop\" param [of cool_primitive]"
	ErrorChecker(err, place)

	element, err = Clicker("//*[@id=\"modal\"]/div/div/form/div[2]/fieldset/div/div[1]/div/div", page)
	place = "Clicking the time form [for stop param of cool_primitive]"
	ErrorChecker(err, place)

	err = page.Find("#timeout").SendKeys(BACK_SPACE + BACK_SPACE + BACK_SPACE)
	place = "3 * backspace"
	ErrorChecker(err, place)

	err = page.Find("#timeout").Fill("15s")
	place = "Filling 15 seconds timeout"
	ErrorChecker(err, place)

	_, err = Clicker("//*[@id=\"modal\"]/div/div/form/div[2]/fieldset/div/div[2]/div/div/select/option[6]", page)
	place = "Setting \"on-fail\" for stop param [of cool_primitive]"
	ErrorChecker(err, place)

	_, err = Clicker("//*[@id=\"modal\"]/div/div/form/div[3]/input", page)
	place = "Submitting stop param [15s, stop - on-fail]"
	ErrorChecker(err, place)

	//--- monitoring part
	_, err = Clicker("//*[@id=\"oplist\"]/fieldset/div/div[1]/div[3]/div[2]/div/div/a[1]", page)
	place = "Clicking \"Edit\" at \"monitoring\" param [of cool_primitive]"
	ErrorChecker(err, place)

	element, err = Clicker("//*[@id=\"modal\"]/div/div/form/div[2]/fieldset/div/div[1]/div/div", page)
	place = "Clicking on \"timeout\" form [for monitoring of cool_primitive]"
	ErrorChecker(err, place)

	err = page.Find("#timeout").SendKeys(BACK_SPACE + BACK_SPACE + BACK_SPACE)
	place = "3 * backspace"
	ErrorChecker(err, place)

	err = page.Find("#timeout").Fill("9s")
	place = "Filling 9 seconds timeout"
	ErrorChecker(err, place)

	err = page.Find("#interval").SendKeys(BACK_SPACE + BACK_SPACE + BACK_SPACE)
	place = "3 * backspace"
	ErrorChecker(err, place)

	err = page.Find("#interval").Fill("13s")
	place = "Filling 13 seconds interval"
	ErrorChecker(err, place)

	_, err = Clicker("//*[@id=\"modal\"]/div/div/form/div[3]/input", page)
	place = "Submitting stop param [15s, stop - on-fail]"
	ErrorChecker(err, place)

	//-----setting target role of cool_primitive
	_, err = Clicker("//*[@id=\"target-role\"]/option[2]", page)
	place = "Selecting \"Started\" [role]"
	ErrorChecker(err, place)

	_, err = Clicker("//*[@id=\"target-role\"]", page)
	place = "Clicking \"target-role\""
	ErrorChecker(err, place)

	//---Pressing "CREATE" [Primitive] ...
	_, err = Clicker("//*[@id=\"new_primitive\"]/span/span/input", page)
	place = "Clicking \"Create\" [primitive]"
	ErrorChecker(err, place)

	time.Sleep(10 * time.Second)

	CrmPrimitiveChecker(ip)

	//---------CLICKING MONITORING----------------------|
	//_, err = Clicker("//*[@id=\"accordion\"]/li[2]/a", page)
	//place = "Clicking MONITORING"
	//ErrorChecker(err, place)

	_, err = Clicker("//*[@id=\"monitoringMenu\"]/li[1]/a", page)
	place = "Clicking Status"
	ErrorChecker(err, place)

	//--------deleting the Cool_Primitive---------------
	_, err = Clicker("//*[@id=\"resources\"]/div[1]/div[2]/div[2]/table/tbody/tr[1]/td[6]/div/div/button", page)
	place = "Clicking cascade next to cool_primitive resource"
	ErrorChecker(err, place)

	_, err = Clicker("//*[@id=\"resources\"]/div[1]/div[2]/div[2]/table/tbody/tr[1]/td[6]/div/div/ul/li[7]/a", page)
	place = "Clicking \"Edit\" in the cascade [of cool_primitive]"
	ErrorChecker(err, place)

	_, err = Clicker("//*[@id=\"middle\"]/div[2]/div[1]/h1/div/div/a[3]", page)
	place = "Clicking \"DELETE\" button"
	ErrorChecker(err, place)

	_, err = Clicker("//*[@id=\"confirmationDialog\"]/div/div/form/div[3]/button[2]", page)
	place = "Clicking \"OK\" in the pop-up window [to delete cool_primitive]"
	ErrorChecker(err, place)

	time.Sleep(20 * time.Second)

	if !strings.Contains(ssher.SSH("root", ip, "crm resource list", "default"), "ocf::") {
		fmt.Println("Deleting the primitive --   PASSED")
	}

	fmt.Println("test finished!")
}

func main() {
	var t *testing.T
	ip := "10.160.64.255"
	linku := "https://" + ip + ":7630"
	Driver := agouti.ChromeDriver()
	if err := Driver.Start(); err != nil {
		t.Fatal("Failed to start Selenium:", err)
	}

	//ssher.SSH("root", ip, "crm status | egrep -i \"stonith|unmanaged\"")
	//time.Sleep(1000 * time.Second)

	page := PageRefresher(linku, Driver)
	Login(linku, page)

	Cluster_Troubler(linku, page, ip)
}
