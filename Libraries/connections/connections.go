package connections

import(
  "fmt"
  "bufio"
  "net"
)

type peerData struct
{
  //int
  //string
  //[]bytes
}

func openConn()
{
  //call recvThread()
  //call sendThread()
}

func sendThread()
{//include ip address in parameters
  //send peerData
}

func recvThread()
{
  //recv peerData
  //open new recv thread
  //this.thread.close()
  dstream, err := net.Listen("tcp",":8080")

  if err != nil
  {
    fmt.Println(err)
    return
  }
  defer dstream.Close()

  for
  {
    con, err := dstream.Accept()
    if err != nil
    {
      fmt.Println(err)
      return
    }

    go handle(con)
  }
}

func handle(con net.Conn)
{
  for
  {
    data, err := bufio.NewReader(con).ReadStream('\n')

    if err != nil
    {
      fmt.Println(err)
      return
    }

    fmt.Println(data)
  }

  con.Close()
}
