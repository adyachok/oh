package remote

import (
	"fmt"
	"encoding/json"
	"strings"
	"log"

	"golang.org/x/crypto/ssh"

	"github.com/tmc/scp"	
	"github.com/adyachok/oh/models"
)

// Executes command on remote host
type RemoteFilesMonitor struct {
	user string
	password string
	host string //host:port
	command string
	sshConfig *ssh.ClientConfig
	sshCh <-chan *models.Command
	destination string
}

func NewRemoteFilesMonitor(user, password, host, dir string, ch <-chan *models.Command) *RemoteFilesMonitor {
	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{ssh.Password(password)},
	}
	sshConfig.HostKeyCallback = ssh.InsecureIgnoreHostKey()

	command_find := fmt.Sprintf("for i in `find %v -maxdepth 1 -not -type d -not -type l`;", dir)
	command_stat := "do stat -c '{ \"name\": \"%n\", \"created\": %W, \"last_access\": %X, \"last_modified\": %Y, \"last_changed\": %Z }' $i; done;"	
	command := command_find + command_stat

	log.Println("[COMMAND]: " + command)

	return &RemoteFilesMonitor{
		user: user,
		password: password,
		host: host,
		command: command,
		sshConfig: sshConfig,
		sshCh: ch,
		destination: dir, 
	}
}

// SSH conection constructor
func (rfm *RemoteFilesMonitor) connect() (*ssh.Client, *ssh.Session, error) {
	client, err := ssh.Dial("tcp", rfm.host, rfm.sshConfig)
	if err != nil {
		return nil, nil, err
	}
	session, err := client.NewSession()
	if err != nil {
		client.Close()
		return nil, nil, err
	}
	return client, session, nil
}

// Exracts information about all files in a specified forlder
func (rfm *RemoteFilesMonitor) listDir() ([]FileInfo, error){
	client, session, err := rfm.connect()
	if err != nil {
		log.Panic(err)
	}
	out, err := session.CombinedOutput(rfm.command)
	if err != nil {
		log.Panic(err)
	}
	var files []FileInfo
	for _, val := range strings.Split(string(out), "\n") {
		if len(val) > 0 {
			fi := FileInfo{}
			err = json.Unmarshal([]byte(val), &fi)
			if err != nil {
				log.Panic(err)
			}
			files = append(files, fi)
		}
		
	}
	
	client.Close()
	return files, err
}

// Copies file on remote host
func (rfm *RemoteFilesMonitor) Copy(src string) {
	_, session, err := rfm.connect()
	if err != nil {
		log.Panic(err)
	}
	err = scp.CopyPath(src, rfm.destination, session)
	if err != nil {
		log.Panic(err)
	}
}

// Stores information about discovered file
type FileInfo struct {
    Name    	string 		`json:"name"`
    CreateTime 	int64		`json:"created"`
    ModTime 	int64 		`json:"last_modified"`
    AccessTime 	int64		`json:"last_access"`
    ChangeTime	int64		`json:"last_changed"`
}

// FIXME: no need for processing logic here - move to the main
func (rfm *RemoteFilesMonitor) Run() {	
	for {
		select {
		case data := <- rfm.sshCh:
			rfm.Copy(data.Filename)
		}
	}	
}


//rm := NewRemoteFilesMonitor("root", "root", "0.0.0.0:32768", "/home")