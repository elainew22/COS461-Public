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
	 dataTwo := make([]byte, 0)
     var buffer bytes.Buffer
	 var err error 
	 var length int

	 fmt.Println(buffer)
	 fmt.Println(dataTwo)

	 bufferTwo := bytes.NewBuffer(dataTwo)

 
	 // read request in newReader
	 for err != io.EOF {
	   length, err = c.Read(data)
	   if checkError(err, c) == -1 {return}
	
	   _, errW := buffer.Write(data[:length])
	  bufferTwo.Write(data[:length])
      if checkError(errW, c) == -1 {return}
	  fmt.Println("here - reading and writing")
	  fmt.Println(data)
	  fmt.Println(buffer)
	  fmt.Println(bufferTwo)
	  fmt.Println(bufferTwo.Bytes())
	 }

	 fmt.Println("finished reading and writing")
 
	 // read request from reader
	 var reader io.Reader
	 //reader.Read(buffer.Bytes())
	 reader.Read(bufferTwo.Bytes())
	 readerBuf := bufio.NewReader(reader)
	 
	 fmt.Println("made readerBuf")
     
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
			 fmt.Println("added Connection close header")
		 } else {
		 req.Header.Set("Connection", "close")
		 fmt.Println("set Connection close header")
		 }
		  
		
		 // get the host url
		 host := req.URL.Host
         if strings.LastIndex(host, ":") == -1 {
          host = host + ":80"
        }
       
		 
		 // connect to server
		 fmt.Println("dialing server...")
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
        	fmt.Println("reading and writing response")
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
