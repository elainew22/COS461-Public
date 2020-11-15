/*****************************************************************************
 * http_proxy.go                                                                 
 * Names: Jana Sebaali and Elaine Wright
 * NetIds: jsebaali, ew22
 *****************************************************************************/

 // TODO: implement an HTTP proxy

 package main

 // Importing packages 
 import (
   "io"
   "os"
   "log"
   "net"
   "bufio"
   "bytes"
   "net/http"
   "io/ioutil"
   "strings"
   "fmt"
 )
 
 const SEND_BUFFER_SIZE = 2048
 
 func internalError(c net.Conn) {
	b := []byte("500 Internal Server Error")
	c.Write(b)
 }
 
 func badReqError(c net.Conn) {
	b := []byte("400 Bad Request Error")
	c.Write(b)
 }
 
 func checkError(err error, c net.Conn) (int){
	 if err != nil {
		  b := []byte("400 Bad Request Error")
		  c.Write(b)
          return -1
	 }
     return 0
 }
 
 
 func handleConnection(c net.Conn){
    fmt.Println("in handle")
 
     data := make([]byte, SEND_BUFFER_SIZE)
     var buffer bytes.Buffer
	 var err error 
 
	 // read request in newReader
	 for err != io.EOF {
	   _ , err := c.Read(data)
	   if checkError(err, c) == -1 {return}
	   
	   _, errW := buffer.Write(data)
      if checkError(errW, c) == -1 {return}
      fmt.Println("here")
	 }
 
	 // read request from reader
	 var reader io.Reader
     reader.Read(buffer.Bytes())
     readerBuf := bufio.NewReader(reader)
     
	 req, err := http.ReadRequest(readerBuf)
     if checkError(err, c) == -1 {return}
 
	 if req.Method != "GET" {
		badReqError(c)
        return
	 } else {
	 
	 // check if it is HTTP 1.1
		 if req.Proto != "HTTP/1.1" {
		   internalError(c)
           return
		 }
		 
		 // update the header to include close connection
		 if req.Header.Get("Connection") == "" {
			 req.Header.Add("Connection", "close")
		 } else {
		 req.Header.Set("Connection", "close")
		 }
		  
		
		 // get the host url
		 host := req.URL.Host
         if strings.LastIndex(host, ":") == -1 {
          host = host + ":80"
        }
       
		 
		 // connect to server
		 cServer, err := net.Dial("tcp", host)
         if checkError(err, c) == -1 {return}
		 
		 // send the request to the server
		 body, err := ioutil.ReadAll(req.Body)
         if checkError(err, c) == -1 {return}
		 cServer.Write(body)
		 
		 // read the response from the connection
         dataRes := make([]byte, SEND_BUFFER_SIZE)
		 var bufferRes bytes.Buffer
         
		 err = nil
		 for err != io.EOF {
		   _ , err := c.Read(dataRes)
		   if checkError(err, c) == -1 {return}
		   
		   _ , errW := bufferRes.Write(dataRes)
	       if checkError(errW, c) == -1 {return}
        fmt.Println("here")
		 }
 
		 // read response from reader
        var readerR io.Reader
        readerR.Read(bufferRes.Bytes())
        readerRBuf := bufio.NewReader(readerR)
	
		 
		 // send the response to the client
		 buf := new(bytes.Buffer)
		 buf.ReadFrom(readerRBuf)
		 c.Write(buf.Bytes())
	 
	 }
 }
 
 func proxy(port string) {
     fmt.Println("will start listening")
     addrPort := ":" + port
	 ln, err := net.Listen("tcp", addrPort)
	 if err != nil {
        return
     }
     fmt.Println("listening")
 
          
     fmt.Println("before handling...")
	 for {
		 conn, err := ln.Accept()
		 if err != nil {
			 // handle error
             return
 
		 }
         fmt.Println("handling...")
		 go handleConnection(conn)
	 }
 
 }
 
 // Main parses command-line arguments and calls client function
 func main() {
   if len(os.Args) != 2 {
	 log.Fatal("")
   }
   port := os.Args[1]
   proxy(port)
 
 }
 
