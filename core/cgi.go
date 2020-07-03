package core

func ExecuteCGIScript(client *Client, request Request) Error {
	return client.Conn().WriteBytes([]byte("iEXECUTING CGI SCRIPT HERE\tnull.host\t0\r\n"))
}
