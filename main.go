package main

import (
	"fmt"
	"net"
	"net/smtp"
	"os"
	"strings"
)

var (
	hostlist     = make([]string, 0)
	result       = ""
	invalidemail = "Email id is invalid"
	err          error
)

func main() {
	//https://play.golang.org/p/8_FCJsRmg2
	var fromemail, toemail string = "bharatsewani@gmail.com", "bharatsewani@gmail.com"

	//get the domain name from email recipient
	email_arr := strings.Split(toemail, "@")
	domain := email_arr[1]
	fmt.Println("Host name is==>", domain)

	//get mx record of domain
	mxhosts, err := net.LookupMX(domain)
	if err != nil {
		//panic(err)
	}

	//if mx record not found get dns record of domain
	if len(mxhosts) < 1 {
		dnshosts, err := net.LookupNS(domain)
		if err != nil {
			//panic(err)
		}
		if len(dnshosts) > 0 {
			for i := 0; i < len(dnshosts); i++ {
				hostlist = append(hostlist, dnshosts[i].Host)
			}
		}
	} else {
		//found mxhosts
		for i := 0; i < len(mxhosts); i++ {
			hostlist = append(hostlist, mxhosts[i].Host)
		}
	}

	//check if host list is not empty
	if len(hostlist) > 0 {
		ip := hostlist[0]
		//connect to smtp server
		c, err := smtp.Dial(ip + ":25")
		if err != nil {
			//log.Fatal(err)
		}
		//send hello to smtp server
		err = c.Hello(ip)
		if err != nil {
			//fmt.Println("%v", err)
			fmt.Println(invalidemail)
			os.Exit(3)
		} else {
			result = result + "> HELO " + ip + "\n"
		}
		//send from command to server
		err = c.Mail(fromemail)
		if err != nil {
			//fmt.Println("%v", err)
			fmt.Println(invalidemail)
			os.Exit(3)
		} else {
			result = result + "> MAIL FROM: <" + fromemail + ">\n"
		}

		//sent to command to server.
		if err = c.Rcpt(toemail); err != nil {
			//fmt.Println("%v",err)
			fmt.Println(invalidemail)
			os.Exit(3)
		} else {
			result = result + "> Mail To: <" + toemail + ">\n"
		}
		err = c.Quit()
		//fmt.Println("command executed",result)
		fmt.Println("Email id is valid")
	} else {
		fmt.Println("Email id is invalid")
		fmt.Println("No suitable mx record found!")
	}

}
