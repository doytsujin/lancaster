// main
package main

import (
	"encoding/hex"
	"fmt"
	"net"
	"time"

	"github.com/urfave/cli"
)

func main() {
	netInterfaceName := ""
	netInterface := (*net.Interface)(nil)
	address := ""
	datagramSize := 1500
	ttl := 8
	loopbackEnable := false

	app := cli.NewApp()

	app.Name = "lancaster"
	app.Description = "UDP multicast file transfer client and server"
	app.Version = "1.0.0"
	app.Authors = []cli.Author{
		{Name: "James Dunne", Email: "james.jdunne@gmail.com"},
	}

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "interface,i",
			Value:       "",
			Usage:       "Interface name to bind to",
			Destination: &netInterfaceName,
		},
		cli.StringFlag{
			Name:        "group,g",
			Value:       "236.0.0.100:1360",
			Usage:       "UDP multicast group for transfers",
			Destination: &address,
		},
		cli.IntFlag{
			Name:        "datagram size,s",
			Value:       1500,
			Destination: &datagramSize,
		},
		cli.IntFlag{
			Name:        "ttl,t",
			Value:       8,
			Destination: &ttl,
		},
		cli.BoolFlag{
			Name:        "loopback enable,l",
			Destination: &loopbackEnable,
		},
	}
	app.Before = func(c *cli.Context) error {
		// Find network interface by name:
		if netInterfaceName != "" {
			var err error
			netInterface, err = net.InterfaceByName(netInterfaceName)
			if err != nil {
				return err
			}
		}
		return nil
	}
	app.Commands = []cli.Command{
		cli.Command{
			Name:    "download",
			Aliases: []string{"d"},
			Usage:   "download files from a multicast group locally",
			Action: func(c *cli.Context) error {
				m, err := NewMulticastListener(address, netInterface)
				if err != nil {
					return err
				}

				m.SetDatagramSize(datagramSize)
				if err != nil {
					return err
				}
				m.SetTTL(ttl)
				if err != nil {
					return err
				}
				m.SetLoopback(loopbackEnable)
				if err != nil {
					return err
				}
				//local := c.Args().First()

				buf := make([]byte, datagramSize)

				// Read UDP messages from multicast:
				for {
					// TODO: use second parameter *net.UDPAddr to authenticate source?
					n, _, err := m.conn.ReadFromUDP(buf)
					if err != nil {
						return err
					}
					msg := buf[:n]
					fmt.Printf("%s", hex.Dump(msg))
				}

				err = m.conn.Close()
				return err
			},
		},
		cli.Command{
			Name:    "serve",
			Aliases: []string{"s"},
			Usage:   "server files to a multicast group",
			Action: func(c *cli.Context) error {
				//local := c.Args().First()
				m, err := NewMulticastSender(address, netInterface)
				if err != nil {
					return err
				}

				m.SetDatagramSize(datagramSize)
				if err != nil {
					return err
				}
				m.SetTTL(ttl)
				if err != nil {
					return err
				}
				m.SetLoopback(loopbackEnable)
				if err != nil {
					return err
				}
				msg := []byte("hello, world!\n")

				// Write UDP messages to multicast:
				for {
					_, err := m.conn.Write(msg)
					if err != nil {
						return err
					}
					fmt.Printf("%s", hex.Dump(msg))
					time.Sleep(1 * time.Second)
				}

				err = m.conn.Close()
				return err
			},
		},
	}

	app.RunAndExitOnError()
	return
}
