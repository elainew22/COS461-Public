/*****************************************************************************
 * http_proxy.go
 * Names: Jana Sebaali and Elaine Wright
 * NetIds: jsebaali, ew22
 *****************************************************************************/

 // TODO: implement an HTTP proxy

 package main

 // Importing packages
 import (
   "os"
   "log"
   "net"
   "bufio"
   "bytes"
   "net/http"
   "strings"
 //  "fmt"
 )
 // io/ioutil
 
 const SEND_BUFFER_SIZE = 20480
 
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
  //  fmt.Println("in handle")
 
   ////  data := make([]byte, SEND_BUFFER_SIZE)
  ////   dataTwo := make([]byte, 0)
    // var buffer bytes.Buffer
  ////  buffer := bytes.NewBuffer(dataTwo)
     var err error
     var length int

   //  fmt.Println(buffer)
   //  fmt.Println(dataTwo)
    test := bufio.NewReader(c)
    

 ////    length, err = c.Read(data)
   ////  if checkError(err, c) == -1 {return}
  
  ////  _, errW := buffer.Write(data[:length])
    //bufferTwo.Write(data[:length])
  ////  if checkError(errW, c) == -1 {return}
  //  fmt.Println("here - reading and writing")

 
    
  //   fmt.Println("finished reading and writing")
 
     // read request from reader
   // HERE reader := bytes.NewReader(buffer.Bytes())
  //HERE  readerBuf := bufio.NewReader(reader)
     
     req, err := http.ReadRequest(test)
     if checkError(err, c) == -1 {return}
 
     if req.Method != "GET" {
        badReqError(c)
        return
     } else {
     
     // check if it is HTTP 1.1
         if req.Proto != "HTTP/1.1" {
           req.Proto = "HTTP/1.1"
          // internalError(c)
          // return
         }
         
         // update the header to include close connection
         if req.Header.Get("Connection") == "" {
             req.Header.Add("Connection", "close")
         } else {
         req.Header.Set("Connection", "close")
         }
      
      //  fmt.Println(req)
         // get the host url
         host := req.URL.Host
         if strings.LastIndex(host, ":") == -1 {
          host = host + ":80"
        }
       
         
         // connect to server
         cServer, err := net.Dial("tcp", host)
    //     fmt.Println(err)
         if checkError(err, c) == -1 {return}
         
        // fmt.Println("sending request...")
         req.Write(cServer)
         
         // read the response from the connection
       // fmt.Println("reading response...")
         dataRes := make([]byte, SEND_BUFFER_SIZE)
         bufferRes := bytes.NewBuffer(dataRes)
         
         err = nil

        length, err = cServer.Read(dataRes)
        if checkError(err, c) == -1 {return}
        
        _ , errW := bufferRes.Write(dataRes[length:])
        if checkError(errW, c) == -1 {return}
 
         // read response from reader
        readerR := bytes.NewReader(bufferRes.Bytes())
        readerRBuf := bufio.NewReader(readerR)
        
        res, err := http.ReadResponse(readerRBuf, nil)
   //     fmt.Println("response:----")
     //   fmt.Println(res)
     
        if checkError(err, c) == -1 {return}
        res.Write(c)
        defer c.Close()
    
     
     }
 }
 
 func proxy(port string) {
   //  fmt.Println("will start listening")
     addrPort := ":" + port
     ln, err := net.Listen("tcp", addrPort)
     if err != nil {
        return
     }
   //  fmt.Println("listening")
 
          
  //   fmt.Println("before handling...")
     for {
         conn, err := ln.Accept()
         if err != nil {
             // handle error
             return
 
         }
       //  fmt.Println("handling...")
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
