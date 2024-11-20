package main
import (
	"fmt"
	"net/http"
	"net/smtp"
	"time"
)

var websites = map[string]int{
	"https://facebook.com": 200,
	"https://google.com":   200,
	"https://twitter.com":  200,
	"http://localhost":     200,
}

const checkInterval int = 5
const reminderInterval int = 1

type webStatus struct {
	web string
	status string
	lastFailure time.Time
}


var from string = "21pt12@psgtech.ac.in"
var password string = "24SEP2003$"

//Reciever email address
var to = []string{
	"maanasarajaraman@vidhyaniketanschool.com",
}

var smtpHost string = "smtp.gmail.com"
var smtpPort string = "587"

func main() {
	webStatusSlice := []webStatus{}
	for {
		if len(webStatusSlice) == 0 {
			fmt.Println("All websites are up")
		}

		for web, expectedStatusCode := range websites {
			res, err := http.Get(web)
			if err != nil {
				//website connection refused
				alertUser(web, err, &webStatusSlice)
				continue
			} else {
				if res.StatusCode != expectedStatusCode {
					errmsg := fmt.Errorf("%v is down", web)
					alertUser(web, errmsg, &webStatusSlice)
				}
			}
		}
		fmt.Printf("Sleep for %v secs\n", checkInterval)
		time.Sleep(time.Duration(checkInterval) * time.Second)

	}
}

func alertUser (web string, err error, webStatusSlice *[]webStatus) {
	//information for status slice
	downWebInfo := webStatus{web, "down", time.Now()}

	if len(*webStatusSlice) > 0 {
		prevAlert := checkforPrev(webStatusSlice, web)

		if !prevAlert {
			fmt.Printf("%v added to list\n", web)
			*webStatusSlice = append(*webStatusSlice, downWebInfo)
			triggerEmail(web)
		} else {
			fmt.Printf("%v already in alert list\n", web)
			triggerAnother := checkforReminderInterval(webStatusSlice, web)
		
			if triggerAnother {
				triggerEmail(web)
			}
		}
	} else {
		fmt.Printf("%v added to list\n", web)
		*webStatusSlice = append(*webStatusSlice, downWebInfo)
		triggerEmail(web)
	}
}


func triggerEmail(web string) {

	message := []byte("Subject : Web Monitor Alert \r\n\r\n" + web + " - Website is down\r\n")
	auth := smtp.PlainAuth("", from, password, smtpHost)

	//Sending error message email
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Email sent successfully!")
}


func checkforPrev(webStatusSlice *[]webStatus, web string) bool {
	alreadyDown := false
	for _, webStatusInfo := range *webStatusSlice {
		if webStatusInfo.web == web {
			alreadyDown = true
		}
	}

	if !alreadyDown {
		//not in alert list
		return false
	} else {
		//already in list
		return true
	}
}
 

func checkforReminderInterval(webStatusSlice *[]webStatus, web string) bool {
	triggerAnother := false
	for i, webStatusInfo := range *webStatusSlice {
		if webStatusInfo.web == web {
			lastFailPlusRem :=  webStatusInfo.lastFailure.Add(time.Duration(reminderInterval) * time.Minute)
			if lastFailPlusRem.Before(time.Now()) {
				triggerAnother = true
				//updating to current time
				(*webStatusSlice)[i] = webStatus{web, "down", time.Now()}
				fmt.Printf("%v : Time for new alert", web)
			} else {
				fmt.Printf("%v : New alert will be send after reminder time finishes\n",web)
			}
		}
	}
	return triggerAnother
}