package ssh_lib

import ("")  //found the ssh package online... lost the exact name of it, but will add later
// need to add some args to pass into this func later
func sendCommand() {
config := &ssh.ClientConfig {
	User: "username",
	Auth: []ssh.AuthMethod{ 
	  publickey("mykey")
	},
	HostKeyCallback: ssh.InsecureIgnoreHostKey(),
  }func publicKey(path string) ssh.AuthMethod {
   key, err := ioutil.ReadFile(path)
   check(err)
   signer, err := ssh.ParsePrivateKey(key)
   check(err)
   return ssh.PublicKeys(signer)
  }
}