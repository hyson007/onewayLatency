package pingpong

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"sort"
	"time"

	"github.com/spf13/cobra"
)

type Packet struct {
	Timestamp int64
}

func RunPinger(cmd *cobra.Command) {
	target, _ := cmd.Flags().GetString("target")
	numPackets, _ := cmd.Flags().GetInt("numPackets")
	numPorts, _ := cmd.Flags().GetInt("numPorts")
	protocol, _ := cmd.Flags().GetString("protocol")

	rand.Seed(time.Now().UnixNano())

	latencies := make([]int64, 0, numPackets*numPorts)

	for i := 0; i < numPorts; i++ {
		sourcePort := rand.Intn(65535-1024) + 1024

		for _, latency := range latencies {
			if latency == int64(sourcePort) {
				continue
			}
		}

		for j := 0; j < numPackets; j++ {
			latency, err := sendPacket(target, sourcePort, protocol)
			if err != nil {
				fmt.Printf("Error sending packet: %v\n", err)
				continue
			}
			latencies = append(latencies, latency)
		}
	}

	sort.Slice(latencies, func(i, j int) bool { return latencies[i] < latencies[j] })

	fmt.Printf("10%% percentile: %d ns\n", latencies[len(latencies)*10/100])
	fmt.Printf("30%% percentile: %d ns\n", latencies[len(latencies)*30/100])
	fmt.Printf("50%% percentile: %d ns\n", latencies[len(latencies)*50/100])
	fmt.Printf("70%% percentile: %d ns\n", latencies[len(latencies)*70/100])
	fmt.Printf("90%% percentile: %d ns\n", latencies[len(latencies)*90/100])
}

func RunPonger(cmd *cobra.Command) {
	port, _ := cmd.Flags().GetInt("port")
	protocol, _ := cmd.Flags().GetString("protocol")

	servePackets(port, protocol)
}

func sendPacket(target string, sourcePort int, protocol string) (int64, error) {
	packet := Packet{Timestamp: time.Now().UnixNano()}
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(packet.Timestamp))

	var conn net.Conn
	var err error

	if protocol == "tcp" {
		conn, err = net.DialTimeout("tcp", target, 1*time.Second)
	} else {
		conn, err = net.DialTimeout("udp", target, 1*time.Second)
	}

	if err != nil {
		return 0, err
	}
	defer conn.Close()

	lAddr := conn.LocalAddr().(*net.UDPAddr)
	lAddr.Port = sourcePort
	conn.(*net.UDPConn).SetWriteDeadline(time.Now().Add(1 * time.Second))

	_, err = conn.Write(buf)
	if err != nil {
		return 0, err
	}

	conn.(*net.UDPConn).SetReadDeadline(time.Now().Add(1 * time.Second))
	n, err := conn.Read(buf)
	if err != nil || n != 8 {
		return 0, errors.New("invalid response from ponger")
	}
	t3 := time.Now().UnixNano()
	t2 := int64(binary.BigEndian.Uint64(buf))
	

	return t3 - t2, nil
}


func servePackets(port int, protocol string) {
	addr := net.TCPAddr{
		Port: port,
		IP:   net.ParseIP("0.0.0.0"),
	}

	if protocol == "tcp" {
		listener, err := net.ListenTCP("tcp", &addr)
		if err != nil {
			fmt.Printf("Error starting TCP listener: %v\n", err)
			return
		}
		defer listener.Close()

		fmt.Printf("Running in ponger mode on %s\n", listener.Addr())

		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Printf("Error accepting TCP connection: %v\n", err)
				continue
			}

			go handleTCPConnection(conn)
		}

	} else {
		udpAddr := net.UDPAddr{
			Port: port,
			IP:   net.ParseIP("0.0.0.0"),
		}

		conn, err := net.ListenUDP("udp", &udpAddr)
		if err != nil {
			fmt.Printf("Error starting UDP listener: %v\n", err)
			return
		}
		defer conn.Close()

		fmt.Printf("Running in ponger mode on %s\n", conn.LocalAddr())

		buf := make([]byte, 8)

		for {
			n, remoteAddr, err := conn.ReadFrom(buf)
			if err != nil || n != 8 {
				continue
			}

			t2 := time.Now().UnixNano()
			binary.BigEndian.PutUint64(buf, uint64(t2))

			conn.SetWriteDeadline(time.Now().Add(1 * time.Second))
			_, err = conn.WriteTo(buf, remoteAddr)
			if err != nil {
				fmt.Printf("Error sending UDP response: %v\n", err)
			}
		}
	}
}

func handleTCPConnection(conn *net.TCPConn) {
	defer conn.Close()

	buf := make([]byte, 8)
	n, err := conn.Read(buf)
	if err != nil || n != 8 {
		return
	}

	t2 := time.Now().UnixNano()
	binary.BigEndian.PutUint64(buf, uint64(t2))

	conn.SetWriteDeadline(time.Now().Add(1 * time.Second))
	_, err = conn.Write(buf)
	if err != nil {
		fmt.Printf("Error sending TCP response: %v\n", err)
	}
}
