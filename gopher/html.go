package gopher

// generateHTMLRedirect takes a URL string and generates an HTML redirect page bytes
func generateHTMLRedirect(url string) []byte {
	content :=
		"<html>\n" +
			"<head>\n" +
			"<meta http-equiv=\"refresh\" content=\"1;URL=" + url + "\">" +
			"</head>\n" +
			"<body>\n" +
			"You are following an external link to a web site.\n" +
			"You will be automatically taken to the site shortly.\n" +
			"If you do not get sent there, please click <A HREF=\"" + url + "\">here</A> to go to the web site.\n" +
			"<p>\n" +
			"The URL linked is <A HREF=\"" + url + "\">" + url + "</A>\n" +
			"<p>\n" +
			"Thanks for using Gophor!\n" +
			"</body>\n" +
			"</html>\n"

	return []byte(content)
}
