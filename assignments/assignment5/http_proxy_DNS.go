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
   "net/http"
   "strings"
   "golang.org/x/net/html"
   
 )
 
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
    
     reqR := bufio.NewReader(c)
     req, err := http.ReadRequest(reqR)
     if checkError(err, c) == -1 {return}
 
     if req.Method != "GET" {
        badReqError(c)
        return
     } else {
     
        // check if it is HTTP 1.1
         if req.Proto != "HTTP/1.1" {
           req.Proto = "HTTP/1.1"
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
         
        req.Write(cServer)
        resR := bufio.NewReader(cServer)
        go dnsFetching(resR, c)
        
        res, err := http.ReadResponse(resR, nil)
        if checkError(err, c) == -1 {return}
        
        res.Write(c)
        defer c.Close()
     
     }
 }
 
  func proxy(port string) {
     addrPort := ":" + port
     ln, err := net.Listen("tcp", addrPort)
     if err != nil {
        return
     }
 
     for {
         conn, err := ln.Accept()
         if err != nil {
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
