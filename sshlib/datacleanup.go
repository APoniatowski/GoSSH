package sshlib

// defaulter defaults all empty fields in the parsed pool entry and aborts if too
// many values are missing, eg both password and key_path.
func (pp *ParsedPool) defaulter() {
	if pp.Password == "" && pp.KeyPath == "" {
		panic("Both 'Password' and 'Key_Path' fields are empty... Aborting.\n")
	}
	if pp.Username == "" {
		pp.Username = "root"
	}
	if pp.Port == 0 {
		pp.Port = 22
	}
}
