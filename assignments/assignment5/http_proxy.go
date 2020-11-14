/*****************************************************************************
 * http_proxy.go                                                                 
 * Names: 
 * NetIds:
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
)

const SEND_BUFFER_SIZE := 2048

func handleConnection(c Conn){

    // check for a properly-formatted HTTP request
    

    var data [SEND_BUFFER_SIZE]byte
    length, err := c.Read(data)
    

    if err != nil && err != io.EOF {
    // handle error
    }

    bytesReader := bytes.NewReader(data)
    reader := bufio.NewReader(bytesReader)

    req, err := reader.ReadRequest()
    
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


 
