/*****************************************************************************
 * http_proxy_DNS.go
 * Names: Jana Sebaali and Elaine Wright
 * NetIds: jsebaali, ew22
 *****************************************************************************/

 // TODO: implement an HTTP proxy

 package main

 // Importing packages
 import (
   "os"
   "io"
   "log"
   "net"
   "bufio"
   "bytes"
   "net/http"
   "strings"
  // "fmt"
   "golang.org/x/net/html"
   
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
 
 func dnsFetching(r io.Reader, c net.Conn){
    _, err := html.Parse(r)
            if err != nil{
       // fmt.Println("parsing response")
        }
    if checkError(err, c) == -1 {return}
    
    var f func(*html.Node)
    f = func(n *html.Node) {
        if n.Type == html.ElementNode && n.Data == "a" {
            for _, a := range n.Attr {
                if a.Key == "href" {
                    go net.LookupHost(a.Val)
                    break
                }
            }
        }
        for c := n.FirstChild; c != nil; c = c.NextSibling {
            f(c)
        }
    }
}
 func handleConnection(c net.Conn){
   // fmt.Println("in handle")
 
   ////  data := make([]byte, SEND_BUFFER_SIZE)
   //  dataTwo := make([]byte, 0)
   //  buffer := bytes.NewBuffer(dataTwo)
   //  var err error
    // var length int

      test := bufio.NewReader(c)
   //  if checkError(err, c) == -1 {return}
  
  //  _, errW := buffer.Write(data[:length])
   // if checkError(errW, c) == -1 {return}

 //   reader := bytes.NewReader(buffer.Bytes())
  //  readerBuf := bufio.NewReader(reader)

     
     req, err := http.ReadRequest(test)

    

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
         //    fmt.Println("added Connection close header")
         } else {
         req.Header.Set("Connection", "close")
        // fmt.Println("set Connection close header")
         }
          //  fmt.Println("")
       //   fmt.Println(req)
             
         // get the host url
         host := req.URL.Host
         if strings.LastIndex(host, ":") == -1 {
          host = host + ":80"
        }
       
         // connect to server
         cServer, err := net.Dial("tcp", host)
         if checkError(err, c) == -1 {return}
         
     //    fmt.Println("sending request...")
        
         req.Write(cServer)
         
         // read the response from the connection
      //  fmt.Println("reading response...")
         dataRes := make([]byte, SEND_BUFFER_SIZE)
         bufferRes := bytes.NewBuffer(dataRes)
         
         err = nil
        length, err := cServer.Read(dataRes)
        if checkError(err, c) == -1 {return}
        
        _ , errW := bufferRes.Write(dataRes[length:])
        if checkError(errW, c) == -1 {return}
 
         // read response from reader
        readerR := bytes.NewReader(bufferRes.Bytes())
        go dnsFetching(readerR, c)
        readerRBuf := bufio.NewReader(readerR)
        res, _ := http.ReadResponse(readerRBuf, nil)
              
    //    fmt.Println("response:----")
    //    fmt.Println(res)
      //  if checkError(err, c) == -1 {return}
        
    //    fmt.Println("sending response")
        res.Write(c)
        c.Close()
    
     }
 }
 
 func proxy(port string) {
     addrPort := ":" + port
     ln, err := net.Listen("tcp", addrPort)
     if err != nil {
        return
     }
  //   fmt.Println("listening")
 
     for {
         conn, err := ln.Accept()
         if err != nil {
             // handle error
             return
 
         }
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
