package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

const basePath = "/home/tiago/fileservertest/"

type Handler struct {
	homeDir string
}
type lister []os.FileInfo

type FileInfo struct {
	name string
	size int64
	dir  bool
	mode os.FileMode
}

func (f FileInfo) Name() string       { return f.name }
func (f FileInfo) Size() int64        { return f.size }
func (f FileInfo) Dir() bool          { return f.dir }
func (f FileInfo) Mode() os.FileMode  { return f.mode }
func (f FileInfo) ModTime() time.Time { return time.Now() }
func (f FileInfo) IsDir() bool        { return f.dir }

func (f *FileInfo) SetDir(dir bool) {
	f.dir = dir
}
func (f FileInfo) Sys() any { return nil }

func (l lister) ListAt(dest []os.FileInfo, off int64) (int, error) {
	if off >= int64(len(l)) {
		return 0, io.EOF
	}
	n := copy(dest, l[off:])
	return n, nil
}

func main() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Println("couild not read .env")
	}

	config := &ssh.ServerConfig{
		PasswordCallback: func(c ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error) {
			users := map[string]string{
				"testuser": "tiger",
				"tiago":    os.Getenv("TIAGOPW"),
			}

			dir := map[string]string{
				"testuser": "/",
				"tiago":    "/tiago",
			}

			if want, ok := users[c.User()]; ok && string(pass) == want {
				return &ssh.Permissions{
					Extensions: map[string]string{
						"home": dir[c.User()],
					},
				}, nil
			}

			return nil, fmt.Errorf("password rejected for %q", c.User())
		},
	}

	privateBytes, err := os.ReadFile("/home/tiago/.ssh/id_rsa")
	if err != nil {
		log.Fatal("Failed to load private key", err)
	}

	private, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		log.Fatal("Failed to parse private key", err)
	}

	config.AddHostKey(private)

	listener, err := net.Listen("tcp", "0.0.0.0:3022")
	if err != nil {
		log.Fatal("failed to listen for connection", err)
	}
	fmt.Printf("Listening on %v\n", listener.Addr())

	nConn, err := listener.Accept()
	if err != nil {
		log.Println("failed to accept incoming connection", err)
		return
	}

	HandleConn(nConn, config)
}

func HandleConn(nConn net.Conn, config *ssh.ServerConfig) {
	defer nConn.Close()

	sshConn, chans, reqs, err := ssh.NewServerConn(nConn, config)
	if err != nil {
		log.Println(err)
		return
	}
	defer sshConn.Close()

	go ssh.DiscardRequests(reqs)
	for newChannel := range chans {
		if newChannel.ChannelType() != "session" {
			newChannel.Reject(ssh.UnknownChannelType, "unknown channel")
			continue
		}

		channel, requests, err := newChannel.Accept()
		if err != nil {
			continue
		}
		go HandleSession(channel, requests, sshConn, sshConn.Permissions.Extensions["home"])
	}
}
func HandleSession(channel ssh.Channel, requests <-chan *ssh.Request, conn *ssh.ServerConn, home string) {
	defer channel.Close()

	go func(in <-chan *ssh.Request) {
		for req := range in {
			if req.Type != "subsystem" || string(req.Payload[4:]) != "sftp" {
				req.Reply(false, nil)
				return
			}
			req.Reply(true, nil)
		}
	}(requests)

	h := &Handler{homeDir: home}
	handlers := sftp.Handlers{FileGet: h, FileList: h, FilePut: h, FileCmd: h}
	server := sftp.NewRequestServer(channel, handlers)
	err := server.Serve()
	if err != nil && err != io.EOF {
		log.Println(err)
	}

	server.Close()
}
func (h *Handler) Filelist(request *sftp.Request) (sftp.ListerAt, error) {
	p := path.Join(basePath, h.homeDir, request.Filepath)

	switch request.Method {
	case "List":
		dir, err := os.ReadDir(p)
		if err != nil {
			log.Println("Could not read the dir")
			return nil, err
		}
		var res []os.FileInfo

		for _, file := range dir {
			inf, err := file.Info()
			if err != nil {
				log.Println("Could not read fileinfo")
				return nil, err
			}
			res = append(res, inf)
		}
		return lister(res), nil
	case "Stat":
		var res os.FileInfo
		var err error

		res, err = os.Stat(p)
		if err != nil {
			log.Println("Could not read Stat of file")
			return nil, err
		}

		return lister([]os.FileInfo{res}), nil
	}

	return nil, errors.New("internal error")
}
func (h *Handler) Fileread(req *sftp.Request) (io.ReaderAt, error) {
	file, err := os.ReadFile(path.Join(basePath, h.homeDir, strings.ReplaceAll(req.Filepath, "../", "/")))
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(file), nil
}

func (h *Handler) Filewrite(req *sftp.Request) (io.WriterAt, error) {
	// TBD (open file (O_CREATE), write data)
	file, err := os.Create(path.Join(basePath, h.homeDir, req.Filepath))
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (h *Handler) Filecmd(req *sftp.Request) error {
	switch req.Method {
	case "Rename":
		return os.Rename(path.Join(basePath, h.homeDir, req.Filepath), path.Join(basePath, h.homeDir, req.Target))
	case "Remove":
		return os.Remove(path.Join(basePath, h.homeDir, req.Filepath))
	case "Mkdir":
		return os.Mkdir(path.Join(basePath, h.homeDir, req.Filepath), 0755)
	case "Rmdir":
		return os.Remove(path.Join(basePath, h.homeDir, req.Filepath))
	}
	return errors.New("Internal error")
}
