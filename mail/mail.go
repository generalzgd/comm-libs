/**
 * @version: 1.0.0
 * @author: zhangguodong:general_zgd
 * @license: LGPL v3
 * @contact: general_zgd@163.com
 * @site: github.com/generalzgd
 * @software: Gogland
 * @file: mail.go
 * @time: 2017/10/20 19:34
 */
package mail

/*func SendMail(user, password, host, subject, to string, attaches []string) error {
	hp := strings.Split(host, ":")
	auth := smtp.PlainAuth("", user, password, hp[0])

	buffer := bytes.NewBuffer(nil)
	boudary := "~~~~~~~~~~~~~~~~~~~"

	header := fmt.Sprintf("To:%s\r\n" +
		"From:%s\r\n" +
		"Subject:%s\r\n" +
		"Content-Type:multipart/mixed;Boundary=\"%s\"\r\n" +
		"Mime-Version:1.0\r\n" +
		"Date:%s\r\n", to, user, subject,boudary,time.Now().String())

	buffer.WriteString(header)
	fmt.Println(header)

	msg1 := "\r\n\r\n--" + boudary + "\r\n"+"Content-Type:text/plain;charset-utf-0\r\n\r\n这是正文\r\n"

	buffer.WriteString(msg1)
	fmt.Println(msg1)

	msg2 := fmt.Sprintf("\r\n--%s\r\n" +
		"Content-Transfer-Encoding: base64\r\n" +
		"Content-Disposition: attachment;\r\n" +
		"Content-Type:image/jpg;name=\"test.jpg\"\r\n", boudary)

	buffer.WriteString(msg2)
	fmt.Println(msg2)

	attachmentBytes, err := ioutil.ReadFile("./test.jpg")
	if err != nil {
		fmt.Println("ReadFile Error. ", err.Error())
		return err
	}

	b := make([]byte, base64.StdEncoding.EncodedLen(len(attachmentBytes)))
	base64.StdEncoding.Encode(b, attachmentBytes)
	buffer.WriteString("\r\n")
	for i,l:=0,len(b);i<l;i++{
		buffer.WriteByte(b[i])
		if (i+1)%76 == 0 {
			buffer.WriteString("\r\n")
		}
	}

	buffer.WriteString("\r\n--"+boudary+"--")

	sendTo := strings.Split(to, ";")
	err = smtp.SendMail(host, auth, user, sendTo, buffer.Bytes())

	return err
}*/
