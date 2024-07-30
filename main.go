package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"net"
	"os"
	"time"

	"github.com/jackpal/bencode-go"
)

type bencodeInfo struct {
	Pieces      string `bencode:"pieces"`
	PieceLength int    `bencode:"piece length"`
	Length      int    `bencode:"length"`
	Name        string `bencode:"name"`
}

type bencodeTorrent struct {
	Announce string      `bencode:"announce"`
	Info     bencodeInfo `bencode:"info"`
}

type torrentInfo struct {
	Announce    string
	Pieces      string
	PieceLength int
	Length      int
	Name        string
	InfoHash    [20]byte
}

func main() {
	dat, err := os.ReadFile("test.torrent")
	if err != nil {
		panic(err)
	}

	var x bencodeTorrent
	err = bencode.Unmarshal(bytes.NewReader(dat), &x)
	if err != nil {
		panic(err)
	}
	var info = torrentInfo{Announce: x.Announce, Pieces: x.Info.Pieces, PieceLength: x.Info.PieceLength, Length: x.Info.Length, Name: x.Info.Name}
	fmt.Println(info.Announce[6:])
	udpAdress, err := net.ResolveUDPAddr("udp", info.Announce[6:])
	if err != nil {
		panic(err)
	}
	fmt.Println(udpAdress.IP)
	connection, err := net.ListenUDP("udp", udpAdress)
	if err != nil {
		fmt.Println("Error making connection")
		panic(err)
	}
	writebuf := make([]byte, 0)
	buf := make([]byte, 2048)
	writebuf = binary.BigEndian.AppendUint64(writebuf, 0x41727101980)
	writebuf = binary.BigEndian.AppendUint32(writebuf, 0)
	writebuf = binary.BigEndian.AppendUint32(writebuf, rand.Uint32())
	fmt.Println(writebuf)
	n, err := connection.Write(writebuf)
	if err != nil {
		panic(err)
	}
	fmt.Println(n)
	exp := 0.0
	reader := bufio.NewReader(connection)
	connection.SetReadDeadline(time.Now().Add(15 * time.Second * time.Duration((math.Pow(2, exp)))))
	data, err := reader.Read(buf)
	if errors.Is(err, os.ErrDeadlineExceeded) {
		exp += 1.0
		connection.SetReadDeadline(time.Now().Add(15 * time.Second * time.Duration((math.Pow(2, exp)))))
		data, err = reader.Read(buf)
	}
	fmt.Println(data)
	connection.Close()
}
