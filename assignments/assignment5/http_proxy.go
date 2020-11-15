/*****************************************************************************
 * http_proxy.go                                                                 
 * Names: Jana Sebaali and Elaine Wright
 * NetIds: jsebaali, ew22
 *****************************************************************************/

 // TODO: implement an HTTP proxy

package main

// Importing packages 
import (
 "fmt"
  "io"
  "os"
  "log"
  "net"
  "bufio"
  "bytes"
  "strings"
  "net/http"
  "net/url"
  "io/ioutil"
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

func checkError(err error, c net.Conn){
    if err != nil {
         b := []byte("400 Bad Request Error")
         c.Write(b)
    }
}


func handleConnection(c net.Conn){

    var data [SEND_BUFFER_SIZE]byte
    var builder strings.Builder

    // read request in newReader
    err := nil
    for err != io.EOF {
      length, err := c.Read(data)
      checkError(err, c)
      
      lengthData, errW := builder.write(data)
      checkError(errW, c) // TODO need to return from method if error
    }

    // read request from reader
    reader := bufio.NewReader(builder.String())
    req, err := http.ReadRequest(reader)

    if req.method != "GET" {
       badReqError(c)
    } else {
    
    // check if it is HTTP 1.1
        if req.Proto != "HTTP/1.1" {
          internalError(c)
        }
        
        // update the header to include close connection
        if req.Header.Get("Connection") == "" {
            req.Header.Add("Connection", "close")
        } else {
        req.Header.Set("Connection", "close")
        }
         
       
        // get the host url
        host := req.URL.Host
        u, err := url.Parse(host)
        checkError(err, c)
        
        // if the port is not available, append the default port
        port := u.Port()
        if port==""{
         host = host + ":80"
         port = "80"
        }
        
        // connect to server
        cServer, err := net.Dial("tcp", host)
        checkError(err, c)
        
        // send the request to the server
        body, err := ioutil.ReadAll(req.Body)
        cServer.Write(string(body))
        
        // read the response from the connection
        var dataRes [SEND_BUFFER_SIZE]byte
        var builderRes strings.builder
        err := nil
        for err != io.EOF {
          _ , err := c.Read(dataRes)
          checkError(err, c)
          
          _ , errW := builderRes.write(dataRes)
          checkError(errW, c)
        }

        // read response from reader
        readerRes := bufio.NewReader(builderRes.String())
        
        // send the response to the client
        buf := new(bytes.Buffer)
        buf.ReadFrom(readerRes)
        c.Write(buf.Bytes())
    
    }
}

func proxy(port string) {
    addrPort := ":" + port
    ln, err := net.Listen("tcp", addrPort)
    if err != nil {
	    // handle error

    }

    for {
	    conn, err := ln.Accept()
	    if err != nil {
		    // handle error

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


 

