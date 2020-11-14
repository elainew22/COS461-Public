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
)

const SEND_BUFFER_SIZE := 2048

//func internalError() {
 // Error(w http.ResponseWriter, "Internal Server Error", 500)
//}

func handleConnection(c Conn){

    // check for a properly-formatted HTTP request
    

    var data [SEND_BUFFER_SIZE]byte
    var builder strings.builder


    length, err := c.Read(data)
    if err != nil {
      // return error on Read() somehow
      return err
    }
    
    //err != nil && 
    while err != io.EOF {
      lengthData, err := builder.write(data) 
      length, err := c.Read(data)
      if err != nil {
        // return error on Read() somehow
        return err
      }
    }

    // write EOF part
    lengthData, err := builder.write(data) 


    reader := bufio.newReader(builder.String())

  //  bytesReader := bytes.NewReader(data)
  //  reader := bufio.NewReader(bytesReader)

    req, err := reader.ReadRequest()

    // make HTTP response
    func ReadResponse(r *bufio.Reader, req *Request) (*Response, error)

    response, err := ReadResponse(reader, )


    if req.method != 'GET' {
      // return response w 400 error 
      resp, err := http.get('http://www.princeton.edu/%') 
      // send to client
      // func (c *TCPConn) Write(b []byte) (int, error)
      // c.Write("400 Bad Request")
      writer := bufio.newWriter()
      err := resp.Write(writer)
      if err != nil {
        // there was a problem

      }
      // write array
      // turn this into func later?
      var writeByte [writer.Size()]byte
      len, err := writer.Write(writeByte)
      err := c.Write(writeByte)
    }
    

    // else talk to server and do GET

    // fix headers
    
    resp, err := http.get(req)
    
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


 

