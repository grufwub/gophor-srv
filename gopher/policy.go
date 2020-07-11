package gopher

import "gophor/core"

func generatePolicyHeader(name string) string {
	text := "# This is an automatically generated" + "\r\n"
	text += "# server policy file: " + name + "\r\n"
	text += "#" + "\r\n"
	text += "# BlackLivesMatter" + "\r\n"
	return text
}

func generateCapsTxt(desc, admin, geo string) []byte {
	text := "CAPS" + "\r\n"
	text += "\r\n"
	text += generatePolicyHeader("caps.txt")
	text += "\r\n"
	text += "CapsVersion=1" + "\r\n"
	text += "ExpireCapsAfter=1800" + "\r\n"
	text += "\r\n"
	text += "PathDelimeter=/" + "\r\n"
	text += "PathIdentity=." + "\r\n"
	text += "PathParent=.." + "\r\n"
	text += "PathParentDouble=FALSE" + "\r\n"
	text += "PathEscapeCharacter=\\" + "\r\n"
	text += "PathKeepPreDelimeter=FALSE" + "\r\n"
	text += "\r\n"
	text += "ServerSoftware=Gophor" + "\r\n"
	text += "ServerSoftwareVersion=" + core.Version + "\r\n"
	text += "ServerDescription=" + desc + "\r\n"
	text += "ServerGeolocationString=" + geo + "\r\n"
	text += "ServerDefaultEncoding=utf-8" + "\r\n"
	text += "\r\n"
	text += "ServerAdmin=" + admin + "\r\n"
	return []byte(text)
}

func generateRobotsTxt() []byte {
	text := generatePolicyHeader("robots.txt")
	text += "\r\n"
	text += "Usage-agent: *" + "\r\n"
	text += "Disallow: *" + "\r\n"
	text += "\r\n"
	text += "Crawl-delay: 99999" + "\r\n"
	text += "\r\n"
	text += "# This server does not support scraping" + "\r\n"
	return []byte(text)
}
